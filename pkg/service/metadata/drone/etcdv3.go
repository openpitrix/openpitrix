// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package drone

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"

	"openpitrix.io/libconfd"
	"openpitrix.io/openpitrix/pkg/logger"
)

var (
	_ libconfd.BackendClient = (*_EtcdClient)(nil)
)

const Etcdv3BackendType = "libconfd-backend-etcdv3"

func init() {
	libconfd.RegisterBackendClient(
		Etcdv3BackendType,
		func(cfg *libconfd.BackendConfig) (libconfd.BackendClient, error) {
			return NewEtcdClient(cfg)
		},
	)
}

type _EtcdClient struct {
	cfg clientv3.Config
}

func NewEtcdClient(cfg *libconfd.BackendConfig) (libconfd.BackendClient, error) {
	etcdConfig := clientv3.Config{
		Endpoints:            cfg.Host,
		DialTimeout:          5 * time.Second,
		DialKeepAliveTime:    10 * time.Second,
		DialKeepAliveTimeout: 3 * time.Second,
	}

	etcdConfig.Username = cfg.UserName
	etcdConfig.Password = cfg.Password

	tlsEnabled := false
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
	}

	if cfg.ClientCAKeys != "" {
		certBytes, err := ioutil.ReadFile(cfg.ClientCAKeys)
		if err != nil {
			return nil, err
		}

		caCertPool := x509.NewCertPool()
		ok := caCertPool.AppendCertsFromPEM(certBytes)

		if ok {
			tlsConfig.RootCAs = caCertPool
		}
		tlsEnabled = true
	}

	if cfg.ClientCert != "" && cfg.ClientKey != "" {
		tlsCert, err := tls.LoadX509KeyPair(cfg.ClientCert, cfg.ClientKey)
		if err != nil {
			return nil, err
		}
		tlsConfig.Certificates = []tls.Certificate{tlsCert}
		tlsEnabled = true
	}

	if tlsEnabled {
		etcdConfig.TLS = tlsConfig
	}

	return &_EtcdClient{etcdConfig}, nil
}

func (c *_EtcdClient) Type() string {
	return Etcdv3BackendType
}

func (c *_EtcdClient) WatchEnabled() bool {
	return true
}

// GetValues queries etcd for keys prefixed by prefix.
func (c *_EtcdClient) GetValues(keys []string) (map[string]string, error) {
	vars := make(map[string]string)

	client, err := clientv3.New(c.cfg)
	if err != nil {
		return vars, err
	}
	defer client.Close()

	for _, key := range keys {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
		resp, err := client.Get(ctx, key, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortDescend))
		cancel()
		if err != nil {
			return vars, err
		}
		for _, ev := range resp.Kvs {
			vars[string(ev.Key)] = string(ev.Value)
		}
	}
	return vars, nil
}

func (c *_EtcdClient) WatchPrefix(prefix string, keys []string, waitIndex uint64, stopChan chan bool) (uint64, error) {
	var err error

	// return something > 0 to trigger a key retrieval from the store
	if waitIndex == 0 {
		return 1, err
	}

	client, err := clientv3.New(c.cfg)
	if err != nil {
		return 1, err
	}
	defer client.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancelRoutine := make(chan bool)
	defer close(cancelRoutine)

	go func() {
		select {
		case <-stopChan:
			cancel()
		case <-cancelRoutine:
			return
		}
	}()

	rch := client.Watch(ctx, prefix, clientv3.WithPrefix())

	for wresp := range rch {
		for _, ev := range wresp.Events {
			logger.Debugf("Key updated %s", string(ev.Kv.Key))

			// Only return if we have a key prefix we care about.
			// This is not an exact match on the key so there is a chance
			// we will still pickup on false positives. The net win here
			// is reducing the scope of keys that can trigger updates.
			for _, k := range keys {
				if strings.HasPrefix(string(ev.Kv.Key), k) {
					return uint64(ev.Kv.Version), err
				}
			}
		}
	}

	return 0, err
}
