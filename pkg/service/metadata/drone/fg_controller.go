// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package drone

import (
	"fmt"
	"math/rand"
	"reflect"
	"sync"

	"github.com/golang/protobuf/proto"

	"openpitrix.io/openpitrix/pkg/logger"
	pbfrontgate "openpitrix.io/openpitrix/pkg/pb/metadata/frontgate"
	pbtypes "openpitrix.io/openpitrix/pkg/pb/metadata/types"
	"openpitrix.io/openpitrix/pkg/util/funcutil"
)

type FrontgateController struct {
	mu        sync.Mutex
	cfg       *pbtypes.FrontgateConfig
	clientMap map[string]*pbfrontgate.FrontgateServiceClient
}

func NewFrontgateController() *FrontgateController {
	return &FrontgateController{}
}

func (p *FrontgateController) GetConfig() (*pbtypes.FrontgateConfig, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.cfg == nil {
		err := fmt.Errorf("drone: Frontgate config is empty")
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	cfg := proto.Clone(p.cfg).(*pbtypes.FrontgateConfig)
	return cfg, nil
}

func (p *FrontgateController) SetConfig(cfg *pbtypes.FrontgateConfig) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if reflect.DeepEqual(p.cfg, cfg) {
		logger.Info(nil, "SetConfig: not changed")
		return nil // not changed
	}

	p.cfg = proto.Clone(cfg).(*pbtypes.FrontgateConfig)
	p.closeAllClient()
	return nil
}

func (p *FrontgateController) ReportSubTaskStatus(in *pbtypes.SubTaskStatus) error {
	logger.Info(nil, "%s taskId: %s", funcutil.CallerName(1), in.TaskId)
	p.mu.Lock()
	defer p.mu.Unlock()

	client, err := p.getClient()
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}

	_, err = client.ReportSubTaskStatus(in)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}

	return nil
}

func (p *FrontgateController) getClient() (
	*pbfrontgate.FrontgateServiceClient,
	error,
) {
	if p.cfg == nil {
		err := fmt.Errorf("drone: Frontgate config is empty")
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	nodes := p.cfg.GetNodeList()
	if len(nodes) == 0 {
		err := fmt.Errorf("drone: no frontgate node")
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	i := rand.Intn(len(nodes))
	addr := fmt.Sprintf("%s:%d", nodes[i].NodeIp, nodes[i].NodePort)
	if x, ok := p.clientMap[addr]; ok {
		return x, nil
	}

	client, err := pbfrontgate.DialFrontgateService("tcp", addr)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}
	p.clientMap[addr] = client

	return client, nil
}

func (p *FrontgateController) closeAllClient() {
	for _, c := range p.clientMap {
		c.Close()
	}

	p.clientMap = make(map[string]*pbfrontgate.FrontgateServiceClient)
}
