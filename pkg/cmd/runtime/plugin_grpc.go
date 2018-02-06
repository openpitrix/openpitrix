// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package runtime

import (
	"context"
	"fmt"
	"sync"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"

	pb "openpitrix.io/openpitrix/pkg/service.pb"
)

var (
	grpcPluginManager = new(grpcRuntimePluginManager)
)

type grpcRuntimePluginManager struct {
	infos   []*pb.AppRuntimePluginInfo
	clients []*grpcRuntimePluginClient
	sync.Mutex
}

func (p *grpcRuntimePluginManager) registerPlugin(args *pb.AppRuntimePluginInfo) error {
	p.Lock()
	defer p.Unlock()

	client, err := openGrpcRuntimePluginClient(args)
	if err != nil {
		return err
	}

	p.infos = append(p.infos, &pb.AppRuntimePluginInfo{
		Name: proto.String(args.GetName()),
		Host: proto.String(args.GetHost()),
		Port: proto.Int32(args.GetPort()),
	})
	p.clients = append(p.clients, client)

	return nil
}

func (p *grpcRuntimePluginManager) getRuntime(name string) RuntimeInterface {
	p.Lock()
	defer p.Unlock()

	for i, info := range p.infos {
		if info.GetName() == name {
			return p.clients[i]
		}
	}
	return nil
}

type grpcRuntimePluginClient struct {
	info   *pb.AppRuntimePluginInfo
	conn   *grpc.ClientConn
	client pb.AppRuntimePluginServiceClient
}

func openGrpcRuntimePluginClient(info *pb.AppRuntimePluginInfo) (*grpcRuntimePluginClient, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", info.GetHost(), info.GetPort()))
	if err != nil {
		return nil, err
	}

	p := &grpcRuntimePluginClient{
		info:   info,
		conn:   conn,
		client: pb.NewAppRuntimePluginServiceClient(conn),
	}

	return p, nil
}

func (p *grpcRuntimePluginClient) Close() error {
	return p.conn.Close()
}

func (p *grpcRuntimePluginClient) Name() string {
	return p.info.GetName()
}

func (p *grpcRuntimePluginClient) Run(app string, args ...string) error {
	return nil
}

func (p *grpcRuntimePluginClient) CreateCluster(appConf string, shouldWait bool, args ...string) (clusterId string, err error) {
	reply, err := p.client.CeaseClusters(
		context.Background(),
		&pb.AppRuntimePluginInput{
			AppConf:    proto.String(appConf),
			ClusterIds: proto.String(""),
			ShouldWait: proto.Bool(shouldWait),
			Args:       args,
		},
	)
	if err != nil {
		return "", err
	}

	clusterId = reply.GetClusterId()
	return
}
func (p *grpcRuntimePluginClient) StopClusters(clusterIds string, shouldWait bool, args ...string) error {
	_, err := p.client.StopClusters(
		context.Background(),
		&pb.AppRuntimePluginInput{
			AppConf:    proto.String(""),
			ClusterIds: proto.String(clusterIds),
			ShouldWait: proto.Bool(shouldWait),
			Args:       args,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
func (p *grpcRuntimePluginClient) StartClusters(clusterIds string, shouldWait bool, args ...string) error {
	_, err := p.client.StartClusters(
		context.Background(),
		&pb.AppRuntimePluginInput{
			AppConf:    proto.String(""),
			ClusterIds: proto.String(clusterIds),
			ShouldWait: proto.Bool(shouldWait),
			Args:       args,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
func (p *grpcRuntimePluginClient) DeleteClusters(clusterIds string, shouldWait bool, args ...string) error {
	_, err := p.client.DeleteClusters(
		context.Background(),
		&pb.AppRuntimePluginInput{
			AppConf:    proto.String(""),
			ClusterIds: proto.String(clusterIds),
			ShouldWait: proto.Bool(shouldWait),
			Args:       args,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
func (p *grpcRuntimePluginClient) RecoverClusters(clusterIds string, shouldWait bool, args ...string) error {
	_, err := p.client.RecoverClusters(
		context.Background(),
		&pb.AppRuntimePluginInput{
			AppConf:    proto.String(""),
			ClusterIds: proto.String(clusterIds),
			ShouldWait: proto.Bool(shouldWait),
			Args:       args,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
func (p *grpcRuntimePluginClient) CeaseClusters(clusterIds string, shouldWait bool, args ...string) error {
	_, err := p.client.CeaseClusters(
		context.Background(),
		&pb.AppRuntimePluginInput{
			AppConf:    proto.String(""),
			ClusterIds: proto.String(clusterIds),
			ShouldWait: proto.Bool(shouldWait),
			Args:       args,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
