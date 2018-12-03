// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// Copyright 2018 Yunify Inc. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// copy from metad/pkg/client/client.go

package backends

import (
	"container/ring"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"openpitrix.io/openpitrix/pkg/libconfd"
	"openpitrix.io/openpitrix/pkg/logger"
)

var (
	_ libconfd.BackendClient = (*MetadClient)(nil)
)

const MetadBackendType = "libconfd-backend-metad"

func init() {
	libconfd.RegisterBackendClient(
		MetadBackendType,
		func(cfg *libconfd.BackendConfig) (libconfd.BackendClient, error) {
			return NewMetadClient(cfg.Host)
		},
	)
}

type MetadClient struct {
	connections *ring.Ring
	current     *MetadConnection
}

type MetadConnection struct {
	url        string
	httpClient *http.Client
	waitIndex  uint64
	errTimes   uint32
}

func (c *MetadClient) Type() string       { return MetadBackendType }
func (c *MetadClient) WatchEnabled() bool { return true }
func (c *MetadClient) Close() error       { return nil }

func NewMetadClient(backendNodes []string) (*MetadClient, error) {
	connections := ring.New(len(backendNodes))
	for _, backendNode := range backendNodes {
		url := "http://" + backendNode
		connection := &MetadConnection{
			url: url,
			httpClient: &http.Client{
				Transport: &http.Transport{
					Proxy: http.ProxyFromEnvironment,
					DialContext: (&net.Dialer{
						KeepAlive: 1 * time.Second,
						DualStack: true,
					}).DialContext,
				},
			},
		}
		connections.Value = connection
		connections = connections.Next()
	}

	client := &MetadClient{
		connections: connections,
	}

	err := client.selectConnection()

	return client, err
}

func (c *MetadClient) selectConnection() error {
	maxTime := 15 * time.Second
	i := 1 * time.Second
	for ; i < maxTime; i *= time.Duration(2) {
		if conn, err := c.testConnection(); err == nil {
			//found available conn
			if c.current != nil {
				atomic.StoreUint32(&c.current.errTimes, 0)
			}
			c.current = conn
			break
		}
		time.Sleep(i)
	}
	if i >= maxTime {
		return fmt.Errorf("fail to connect any backend.")
	}
	logger.Info(nil, "Using Metad URL: "+c.current.url)
	return nil
}

func (c *MetadClient) testConnection() (*MetadConnection, error) {
	//random start
	if c.current == nil {
		rand.Seed(time.Now().Unix())
		r := rand.Intn(c.connections.Len())
		c.connections = c.connections.Move(r)
	}
	c.connections = c.connections.Next()
	conn := c.connections.Value.(*MetadConnection)
	startConn := conn
	_, err := conn.makeMetaDataRequest("/")
	for err != nil {
		logger.Error(nil, "connection to [%s], error: [%v]", conn.url, err)
		c.connections = c.connections.Next()
		conn = c.connections.Value.(*MetadConnection)
		if conn == startConn {
			break
		}
		_, err = conn.makeMetaDataRequest("/")
	}
	return conn, err
}

func (c *MetadClient) GetValues(keys []string) (map[string]string, error) {
	vars := map[string]string{}

	for _, key := range keys {
		body, err := c.current.makeMetaDataRequest(key)
		if err != nil {
			atomic.AddUint32(&c.current.errTimes, 1)
			return vars, err
		}

		var jsonResponse interface{}
		if err = json.Unmarshal(body, &jsonResponse); err != nil {
			return vars, err
		}

		if err = treeWalk(key, jsonResponse, vars); err != nil {
			return vars, err
		}
	}
	return vars, nil
}

func treeWalk(root string, val interface{}, vars map[string]string) error {
	switch val.(type) {
	case map[string]interface{}:
		for k := range val.(map[string]interface{}) {
			treeWalk(strings.Join([]string{root, k}, "/"), val.(map[string]interface{})[k], vars)
		}
	case []interface{}:
		for i, item := range val.([]interface{}) {
			idx := strconv.Itoa(i)
			if i, isMap := item.(map[string]interface{}); isMap {
				if name, exists := i["name"]; exists {
					idx = name.(string)
				}
			}

			treeWalk(strings.Join([]string{root, idx}, "/"), item, vars)
		}
	case bool:
		vars[root] = strconv.FormatBool(val.(bool))
	case string:
		vars[root] = val.(string)
	case float64:
		vars[root] = strconv.FormatFloat(val.(float64), 'f', -1, 64)
	case nil:
		vars[root] = "null"
	default:
		logger.Error(nil, "Unknown type: "+reflect.TypeOf(val).Name())
	}
	return nil
}

func (c *MetadClient) WatchPrefix(prefix string, keys []string, waitIndex uint64, stopChan chan bool) (uint64, error) {
	if c.current.errTimes >= 3 {
		c.selectConnection()
	}

	conn := c.current

	// return something > 0 to trigger a key retrieval from the store
	if waitIndex == 0 {
		conn.waitIndex = 1
		return conn.waitIndex, nil
	}
	// when switch to anther server, so set waitIndex 0, and let server response current version.
	if conn.waitIndex == 0 {
		waitIndex = 0
	}

	done := make(chan struct{})
	defer close(done)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		select {
		case <-stopChan:
			cancel()
		case <-done:
			return
		}
	}()
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s?wait=true&prev_version=%d", conn.url, prefix, waitIndex), nil)
	if err != nil {
		return conn.waitIndex, err
	}

	req.Header.Set("Accept", "application/json")
	req = req.WithContext(ctx)

	// just ignore resp, notify confd to reload metadata from metad
	resp, err := conn.httpClient.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		logger.Error(nil, "failed to connect to metad when watch prefix")
		atomic.AddUint32(&conn.errTimes, 1)
		return conn.waitIndex, err
	}
	if resp.StatusCode != 200 {
		return conn.waitIndex, errors.New(fmt.Sprintf("metad response status [%v], requestID: [%s]", resp.StatusCode, resp.Header.Get("X-Metad-RequestID")))
	}
	versionStr := resp.Header.Get("X-Metad-Version")
	if versionStr != "" {
		v, err := strconv.ParseUint(versionStr, 10, 64)
		if err != nil {
			logger.Error(nil, "Parse X-Metad-Version %s error:%s", versionStr, err.Error())
		}
		conn.waitIndex = v
	} else {
		logger.Warn(nil, "Metad response miss X-Metad-Version header.")
		conn.waitIndex = conn.waitIndex + 1
	}
	return conn.waitIndex, nil
}

func (c *MetadConnection) makeMetaDataRequest(path string) ([]byte, error) {
	req, err := http.NewRequest("GET", strings.Join([]string{c.url, path}, ""), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
