// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package pilot

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"

	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pb/frontgate"
	"openpitrix.io/openpitrix/pkg/pb/types"
)

type FrontgateClientManager struct {
	clientMap map[string][]*fgClient
	sync.Mutex
}

type fgClient struct {
	info   *pbtypes.FrontgateConfig
	client *pbfrontgate.FrontgateServiceClient
	closed chan bool
}

func NewFrontgateClientManager() *FrontgateClientManager {
	return &FrontgateClientManager{
		clientMap: make(map[string][]*fgClient),
	}
}

func (p *FrontgateClientManager) CheckAllClient() {
	// copy client map
	p.Lock()
	clientMap := make(map[string][]*fgClient)
	for id, cs := range p.clientMap {
		clientMap[id] = append([]*fgClient{}, cs...)
	}
	p.Unlock()

	var (
		validClinetMap    = map[string]*fgClient{}
		invalidClientMap  = map[string]*fgClient{}
		invalidClientKeys []string
	)

	// ping frontgate
	for id, cs := range clientMap {
		for _, c := range cs {
			if _, err := c.client.HeartBeat(&pbtypes.Empty{}); err != nil {
				invalidClientKeys = append(invalidClientKeys, c.info.Id+c.info.NodeId)
				invalidClientMap[c.info.Id+c.info.NodeId] = c
				p.CloseClient(id, c.info.NodeId)
			} else {
				validClinetMap[c.info.Id+c.info.NodeId] = c
			}
		}
	}

	sort.Strings(invalidClientKeys)
	for _, key := range invalidClientKeys {
		if _, ok := validClinetMap[key]; ok {
			continue
		}

		c := invalidClientMap[key]
		logger.Infof("fontgate(%s/%s) offline\n", c.info.Id, c.info.NodeId)
	}
}

func (p *FrontgateClientManager) GetClient(id string) (*pbfrontgate.FrontgateServiceClient, error) {
	p.Lock()
	defer p.Unlock()

	if cs := p.clientMap[id]; len(cs) > 0 {
		return cs[rand.Intn(len(cs))].client, nil
	}

	return nil, fmt.Errorf("not found")
}

func (p *FrontgateClientManager) PutClient(c *pbfrontgate.FrontgateServiceClient, info *pbtypes.FrontgateConfig) (closed chan bool) {
	p.Lock()
	defer p.Unlock()

	logger.Infof("fontgate(%s/%s) online\n", info.Id, info.NodeId)

	client := &fgClient{
		client: c,
		info:   info,
		closed: make(chan bool),
	}

	p.clientMap[info.Id] = append(p.clientMap[info.Id], client)
	return client.closed
}

func (p *FrontgateClientManager) CloseClient(id, nodeId string) {
	p.Lock()
	defer p.Unlock()

	cs, ok := p.clientMap[id]
	if !ok {
		return
	}

	for i, t := range cs {
		if t.info.Id == id && t.info.NodeId == nodeId {
			close(t.closed)
			cs[i] = cs[len(cs)-1]
			p.clientMap[id] = cs[:len(cs)-1]
			return
		}
	}
}
