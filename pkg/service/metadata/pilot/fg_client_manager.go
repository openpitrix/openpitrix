// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package pilot

import (
	"fmt"
	"math/rand"
	"net/rpc"
	"sort"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"openpitrix.io/openpitrix/pkg/logger"
	pbfrontgate "openpitrix.io/openpitrix/pkg/pb/metadata/frontgate"
	pbtypes "openpitrix.io/openpitrix/pkg/pb/metadata/types"
)

type FrontgateClientManager struct {
	clientMap map[string][]*FrontgateClient
	sync.Mutex
}

type FrontgateClient struct {
	info *pbtypes.FrontgateConfig
	*pbfrontgate.FrontgateServiceClient
	closed chan bool
}

func NewFrontgateClientManager() *FrontgateClientManager {
	return &FrontgateClientManager{
		clientMap: make(map[string][]*FrontgateClient),
	}
}

func (p *FrontgateClientManager) CheckAllClient() {
	// copy client map
	p.Lock()
	clientMap := make(map[string][]*FrontgateClient)
	for id, cs := range p.clientMap {
		clientMap[id] = append([]*FrontgateClient{}, cs...)
	}
	p.Unlock()

	var (
		validClinetMap    = map[string]*FrontgateClient{}
		invalidClientMap  = map[string]*FrontgateClient{}
		invalidClientKeys []string
	)

	// ping frontgate
	for id, cs := range clientMap {
		for _, c := range cs {
			if _, err := c.HeartBeat(&pbtypes.Empty{}); err != nil {
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
		logger.Info(nil, "frontgate (%s/%s) offline", c.info.Id, c.info.NodeId)
	}
}

func (p *FrontgateClientManager) GetClient(id string) (*FrontgateClient, error) {
	p.Lock()
	defer p.Unlock()

	if cs := p.clientMap[id]; len(cs) > 0 {
		return cs[rand.Intn(len(cs))], nil
	}

	return nil, fmt.Errorf("frontgate [%s] not found", id)
}

func (p *FrontgateClientManager) GetNodeClient(id, nodeId string) (*FrontgateClient, error) {
	p.Lock()
	defer p.Unlock()

	for _, cs := range p.clientMap[id] {
		if cs.info.GetNodeId() == nodeId {
			return cs, nil
		}
	}

	return nil, fmt.Errorf("frontgate [%s] node [%s] not found", id, nodeId)
}

func (p *FrontgateClientManager) PutClient(c *pbfrontgate.FrontgateServiceClient, info *pbtypes.FrontgateConfig) (closed chan bool) {
	p.Lock()
	defer p.Unlock()

	logger.Info(nil, "frontgate (%s/%s) online", info.Id, info.NodeId)

	client := &FrontgateClient{
		FrontgateServiceClient: c,
		info:                   info,
		closed:                 make(chan bool),
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

	if nodeId == "" {
		for _, t := range cs {
			if t.info.Id == id {
				close(t.closed)
			}
		}
		delete(p.clientMap, id)
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

func (_ *FrontgateClientManager) IsFrontgateShutdownError(err error) bool {
	if err == nil || status.Code(err) != codes.Unknown {
		return false
	}

	return err.Error() == rpc.ErrShutdown.Error()
}
