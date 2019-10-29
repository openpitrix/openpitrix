// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package pilot

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/golang/protobuf/proto"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	pbfrontgate "openpitrix.io/openpitrix/pkg/pb/metadata/frontgate"
	pbpilot "openpitrix.io/openpitrix/pkg/pb/metadata/pilot"
	pbtypes "openpitrix.io/openpitrix/pkg/pb/metadata/types"
	"openpitrix.io/openpitrix/pkg/service/metadata/pilot/pilotutil"
	"openpitrix.io/openpitrix/pkg/util/funcutil"
	"openpitrix.io/openpitrix/pkg/version"
)

var (
	_ pbpilot.PilotServiceServer = (*Server)(nil)
)

func (p *Server) GetPilotVersion(context.Context, *pbtypes.Empty) (*pbtypes.Version, error) {
	reply := &pbtypes.Version{
		ShortVersion:   version.ShortVersion,
		GitSha1Version: version.GitSha1Version,
		BuildDate:      version.BuildDate,
	}
	return reply, nil
}
func (p *Server) GetFrontgateVersion(context.Context, *pbtypes.FrontgateId) (*pbtypes.Version, error) {
	err := fmt.Errorf("TODO")
	logger.Warn(nil, "%+v", err)
	return nil, err
}
func (p *Server) GetDroneVersion(context.Context, *pbtypes.DroneEndpoint) (*pbtypes.Version, error) {
	err := fmt.Errorf("TODO")
	logger.Warn(nil, "%+v", err)
	return nil, err
}

func (p *Server) PingMetadataBackend(ctx context.Context, arg *pbtypes.FrontgateId) (*pbtypes.Empty, error) {
	logger.Info(nil, funcutil.CallerName(1))

	client, err := p.fgClientMgr.GetClient(arg.Id)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	defer func() {
		if err != nil && p.fgClientMgr.IsFrontgateShutdownError(err) {
			p.fgClientMgr.CloseClient(client.info.Id, client.info.NodeId)
		}
	}()

	_, err = client.PingMetadataBackend(&pbtypes.Empty{})
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	return &pbtypes.Empty{}, nil
}

func (p *Server) GetPilotConfig(context.Context, *pbtypes.Empty) (*pbtypes.PilotConfig, error) {
	logger.Info(nil, funcutil.CallerName(1))

	return proto.Clone(p.cfg).(*pbtypes.PilotConfig), nil
}

func (p *Server) GetPilotClientTLSConfig(context.Context, *pbtypes.Empty) (*pbtypes.PilotClientTLSConfig, error) {
	logger.Info(nil, funcutil.CallerName(1))

	cfg := &pbtypes.PilotClientTLSConfig{
		CaCrtData:       p.pbTlsCfg.CaCrtData,
		ClientCrtData:   p.pbTlsCfg.ClientCrtData,
		ClientKeyData:   p.pbTlsCfg.ClientKeyData,
		PilotServerName: p.pbTlsCfg.PilotServerName,
	}

	return cfg, nil
}

func (p *Server) GetFrontgateList(context.Context, *pbtypes.Empty) (*pbtypes.FrontgateIdList, error) {
	logger.Info(nil, funcutil.CallerName(1))

	return nil, fmt.Errorf("TODO")
}

func (p *Server) GetFrontgateConfig(ctx context.Context, arg *pbtypes.FrontgateId) (*pbtypes.FrontgateConfig, error) {
	logger.Info(nil, funcutil.CallerName(1))

	client, err := p.fgClientMgr.GetClient(arg.Id)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	defer func() {
		if err != nil && p.fgClientMgr.IsFrontgateShutdownError(err) {
			p.fgClientMgr.CloseClient(client.info.Id, client.info.NodeId)
		}
	}()

	reply, err := client.GetFrontgateConfig(&pbtypes.Empty{})
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	return reply, nil
}

func (p *Server) SetFrontgateConfig(ctx context.Context, arg *pbtypes.FrontgateConfig) (*pbtypes.Empty, error) {
	logger.Info(nil, funcutil.CallerName(1))

	client, err := p.fgClientMgr.GetClient(arg.Id)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	defer func() {
		if err != nil && p.fgClientMgr.IsFrontgateShutdownError(err) {
			p.fgClientMgr.CloseClient(client.info.Id, client.info.NodeId)
		}
	}()

	reply, err := client.SetFrontgateConfig(arg)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	return reply, nil
}

func (p *Server) GetDroneConfig(ctx context.Context, arg *pbtypes.DroneEndpoint) (*pbtypes.DroneConfig, error) {
	logger.Info(nil, funcutil.CallerName(1))

	client, err := p.fgClientMgr.GetClient(arg.FrontgateId)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	defer func() {
		if err != nil && p.fgClientMgr.IsFrontgateShutdownError(err) {
			p.fgClientMgr.CloseClient(client.info.Id, client.info.NodeId)
		}
	}()

	reply, err := client.GetDroneConfig(arg)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	return reply, nil
}
func (p *Server) SetDroneConfig(ctx context.Context, arg *pbtypes.SetDroneConfigRequest) (*pbtypes.Empty, error) {
	logger.Info(nil, funcutil.CallerName(1))

	client, err := p.fgClientMgr.GetClient(arg.Endpoint.FrontgateId)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	defer func() {
		if err != nil && p.fgClientMgr.IsFrontgateShutdownError(err) {
			p.fgClientMgr.CloseClient(client.info.Id, client.info.NodeId)
		}
	}()

	_, err = client.SetDroneConfig(arg)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	return &pbtypes.Empty{}, nil
}

func (p *Server) FrontgateChannel(ch pbpilot.PilotServiceForFrontgate_FrontgateChannelServer) error {
	logger.Info(nil, funcutil.CallerName(1))

	c := pbfrontgate.NewFrontgateServiceClient(
		pilotutil.NewFrontgateChannelFromServer(ch),
	)

	info, err := c.GetFrontgateConfig(&pbtypes.Empty{})
	if err != nil {
		logger.Error(nil, "Get frontgate config failed: %+v", err)
		return err
	}

	// if return, the channel will be closed
	<-p.fgClientMgr.PutClient(c, info)
	return nil
}

func (p *Server) GetConfdConfig(ctx context.Context, arg *pbtypes.ConfdEndpoint) (*pbtypes.ConfdConfig, error) {
	logger.Info(nil, funcutil.CallerName(1))

	client, err := p.fgClientMgr.GetClient(arg.FrontgateId)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	defer func() {
		if err != nil && p.fgClientMgr.IsFrontgateShutdownError(err) {
			p.fgClientMgr.CloseClient(client.info.Id, client.info.NodeId)
		}
	}()

	reply, err := client.GetConfdConfig(&pbtypes.ConfdEndpoint{
		DroneIp:   arg.DroneIp,
		DronePort: arg.DronePort,
	})
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	return reply, nil
}

func (p *Server) IsConfdRunning(ctx context.Context, arg *pbtypes.DroneEndpoint) (*pbtypes.Bool, error) {
	logger.Info(nil, funcutil.CallerName(1))

	client, err := p.fgClientMgr.GetClient(arg.FrontgateId)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	defer func() {
		if err != nil && p.fgClientMgr.IsFrontgateShutdownError(err) {
			p.fgClientMgr.CloseClient(client.info.Id, client.info.NodeId)
		}
	}()

	reply, err := client.IsConfdRunning(&pbtypes.ConfdEndpoint{
		DroneIp:   arg.DroneIp,
		DronePort: arg.DronePort,
	})
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	return &pbtypes.Bool{Value: reply.GetValue()}, nil
}

func (p *Server) StartConfd(ctx context.Context, arg *pbtypes.DroneEndpoint) (*pbtypes.Empty, error) {
	logger.Info(nil, funcutil.CallerName(1))

	client, err := p.fgClientMgr.GetClient(arg.FrontgateId)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	defer func() {
		if err != nil && p.fgClientMgr.IsFrontgateShutdownError(err) {
			p.fgClientMgr.CloseClient(client.info.Id, client.info.NodeId)
		}
	}()

	_, err = client.StartConfd(&pbtypes.ConfdEndpoint{
		DroneIp:   arg.DroneIp,
		DronePort: arg.DronePort,
	})
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	return &pbtypes.Empty{}, nil
}

func (p *Server) StopConfd(ctx context.Context, arg *pbtypes.DroneEndpoint) (*pbtypes.Empty, error) {
	logger.Info(nil, funcutil.CallerName(1))

	client, err := p.fgClientMgr.GetClient(arg.FrontgateId)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	defer func() {
		if err != nil && p.fgClientMgr.IsFrontgateShutdownError(err) {
			p.fgClientMgr.CloseClient(client.info.Id, client.info.NodeId)
		}
	}()

	_, err = client.StopConfd(&pbtypes.ConfdEndpoint{
		DroneIp:   arg.DroneIp,
		DronePort: arg.DronePort,
	})
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	return &pbtypes.Empty{}, nil
}

func (p *Server) RegisterMetadata(ctx context.Context, arg *pbtypes.SubTask_RegisterMetadata) (*pbtypes.Empty, error) {
	logger.Info(nil, funcutil.CallerName(1))

	client, err := p.fgClientMgr.GetClient(arg.FrontgateId)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	defer func() {
		if err != nil && p.fgClientMgr.IsFrontgateShutdownError(err) {
			p.fgClientMgr.CloseClient(client.info.Id, client.info.NodeId)
		}
	}()

	_, err = client.RegisterMetadata(arg)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	return &pbtypes.Empty{}, nil
}

func (p *Server) DeregisterMetadata(ctx context.Context, arg *pbtypes.SubTask_DeregisterMetadata) (*pbtypes.Empty, error) {
	logger.Info(nil, funcutil.CallerName(1))

	client, err := p.fgClientMgr.GetClient(arg.FrontgateId)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	defer func() {
		if err != nil && p.fgClientMgr.IsFrontgateShutdownError(err) {
			p.fgClientMgr.CloseClient(client.info.Id, client.info.NodeId)
		}
	}()

	_, err = client.DeregisterMetadata(arg)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	return &pbtypes.Empty{}, nil
}

func (p *Server) RegisterMetadataMapping(ctx context.Context, arg *pbtypes.SubTask_RegisterMetadata) (*pbtypes.Empty, error) {
	logger.Info(nil, funcutil.CallerName(1))

	client, err := p.fgClientMgr.GetClient(arg.FrontgateId)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	defer func() {
		if err != nil && p.fgClientMgr.IsFrontgateShutdownError(err) {
			p.fgClientMgr.CloseClient(client.info.Id, client.info.NodeId)
		}
	}()

	_, err = client.RegisterMetadataMapping(arg)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	return &pbtypes.Empty{}, nil
}

func (p *Server) DeregisterMetadataMapping(ctx context.Context, arg *pbtypes.SubTask_DeregisterMetadata) (*pbtypes.Empty, error) {
	logger.Info(nil, funcutil.CallerName(1))

	client, err := p.fgClientMgr.GetClient(arg.FrontgateId)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	defer func() {
		if err != nil && p.fgClientMgr.IsFrontgateShutdownError(err) {
			p.fgClientMgr.CloseClient(client.info.Id, client.info.NodeId)
		}
	}()

	_, err = client.DeregisterMetadataMapping(arg)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	return &pbtypes.Empty{}, nil
}

func (p *Server) RegisterCmd(ctx context.Context, arg *pbtypes.SubTask_RegisterCmd) (*pbtypes.Empty, error) {
	logger.Info(nil, funcutil.CallerName(1))

	client, err := p.fgClientMgr.GetClient(arg.FrontgateId)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	defer func() {
		if err != nil && p.fgClientMgr.IsFrontgateShutdownError(err) {
			p.fgClientMgr.CloseClient(client.info.Id, client.info.NodeId)
		}
	}()

	_, err = client.RegisterCmd(arg)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	return &pbtypes.Empty{}, nil
}

func (p *Server) DeregisterCmd(ctx context.Context, arg *pbtypes.SubTask_DeregisterCmd) (*pbtypes.Empty, error) {
	logger.Info(nil, funcutil.CallerName(1))

	client, err := p.fgClientMgr.GetClient(arg.FrontgateId)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	defer func() {
		if err != nil && p.fgClientMgr.IsFrontgateShutdownError(err) {
			p.fgClientMgr.CloseClient(client.info.Id, client.info.NodeId)
		}
	}()

	_, err = client.DeregisterCmd(arg)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	return &pbtypes.Empty{}, nil
}

func (p *Server) ReportSubTaskStatus(ctx context.Context, arg *pbtypes.SubTaskStatus) (*pbtypes.Empty, error) {
	logger.Info(nil, "%s taskId: %s", funcutil.CallerName(1), arg.TaskId)

	p.taskStatusMgr.PutStatus(*arg)
	return &pbtypes.Empty{}, nil
}

func (p *Server) GetSubtaskStatus(ctx context.Context, arg *pbtypes.SubTaskId) (*pbtypes.SubTaskStatus, error) {
	logger.Info(nil, funcutil.CallerName(1))

	if s, ok := p.taskStatusMgr.GetStatus(arg.TaskId); ok {
		return &s, nil
	}

	err := fmt.Errorf("pilot: not found, taskId = %s", arg.TaskId)
	logger.Warn(nil, "%+v", err)
	return nil, err
}

func (p *Server) HandleSubtask(ctx context.Context, msg *pbtypes.SubTaskMessage) (*pbtypes.Empty, error) {
	logger.Info(nil, funcutil.CallerName(1))

	switch msg.Action {
	case pbtypes.SubTaskAction_StartConfd.String():
		var x pbtypes.SubTask_StartConfd
		err := json.Unmarshal([]byte(msg.Directive), &x)
		if err != nil {
			logger.Warn(nil, "%+v", err)
			return nil, err
		}

		x.Action = msg.Action
		x.TaskId = msg.TaskId

		reply, err := p.StartConfd(ctx, &pbtypes.DroneEndpoint{
			FrontgateId: x.FrontgateId,
			DroneIp:     x.DroneIp,
			DronePort:   9112,
		})
		if err != nil {
			p.taskStatusMgr.PutStatus(pbtypes.SubTaskStatus{
				TaskId: x.TaskId,
				Status: constants.StatusFailed,
			})
			return nil, err
		}

		// start confd is async task,
		// the drone service should report the task status
		return reply, nil

	case pbtypes.SubTaskAction_StopConfd.String():
		var x pbtypes.SubTask_StopConfd
		err := json.Unmarshal([]byte(msg.Directive), &x)
		if err != nil {
			logger.Warn(nil, "%+v", err)
			return nil, err
		}

		x.Action = msg.Action
		x.TaskId = msg.TaskId

		reply, err := p.StopConfd(ctx, &pbtypes.DroneEndpoint{
			FrontgateId: x.FrontgateId,
			DroneIp:     x.DroneIp,
			DronePort:   9112,
		})

		if err != nil {
			p.taskStatusMgr.PutStatus(pbtypes.SubTaskStatus{
				TaskId: x.TaskId,
				Status: constants.StatusFailed,
			})
			return nil, err
		}

		p.taskStatusMgr.PutStatus(pbtypes.SubTaskStatus{
			TaskId: x.TaskId,
			Status: constants.StatusSuccessful,
		})
		return reply, nil

	case pbtypes.SubTaskAction_RegisterMetadata.String():

		var x pbtypes.SubTask_RegisterMetadata
		err := json.Unmarshal([]byte(msg.Directive), &x)
		if err != nil {
			logger.Warn(nil, "%+v", err)
			return nil, err
		}

		x.Action = msg.Action
		x.TaskId = msg.TaskId

		reply, err := p.RegisterMetadata(ctx, &x)

		if err != nil {
			p.taskStatusMgr.PutStatus(pbtypes.SubTaskStatus{
				TaskId: x.TaskId,
				Status: constants.StatusFailed,
			})
			return nil, err
		}

		p.taskStatusMgr.PutStatus(pbtypes.SubTaskStatus{
			TaskId: x.TaskId,
			Status: constants.StatusSuccessful,
		})
		return reply, nil

	case pbtypes.SubTaskAction_DeregisterMetadata.String():

		var x pbtypes.SubTask_DeregisterMetadata
		err := json.Unmarshal([]byte(msg.Directive), &x)
		if err != nil {
			logger.Warn(nil, "%+v", err)
			return nil, err
		}

		x.Action = msg.Action
		x.TaskId = msg.TaskId

		reply, err := p.DeregisterMetadata(ctx, &x)

		if err != nil {
			p.taskStatusMgr.PutStatus(pbtypes.SubTaskStatus{
				TaskId: x.TaskId,
				Status: constants.StatusFailed,
			})
			return nil, err
		}

		p.taskStatusMgr.PutStatus(pbtypes.SubTaskStatus{
			TaskId: x.TaskId,
			Status: constants.StatusSuccessful,
		})
		return reply, nil

	case pbtypes.SubTaskAction_RegisterMetadataMapping.String():

		var x pbtypes.SubTask_RegisterMetadata
		err := json.Unmarshal([]byte(msg.Directive), &x)
		if err != nil {
			logger.Warn(nil, "%+v", err)
			return nil, err
		}

		x.Action = msg.Action
		x.TaskId = msg.TaskId

		reply, err := p.RegisterMetadataMapping(ctx, &x)

		if err != nil {
			p.taskStatusMgr.PutStatus(pbtypes.SubTaskStatus{
				TaskId: x.TaskId,
				Status: constants.StatusFailed,
			})
			return nil, err
		}

		p.taskStatusMgr.PutStatus(pbtypes.SubTaskStatus{
			TaskId: x.TaskId,
			Status: constants.StatusSuccessful,
		})
		return reply, nil

	case pbtypes.SubTaskAction_DeregisterMetadataMapping.String():

		var x pbtypes.SubTask_DeregisterMetadata
		err := json.Unmarshal([]byte(msg.Directive), &x)
		if err != nil {
			logger.Warn(nil, "%+v", err)
			return nil, err
		}

		x.Action = msg.Action
		x.TaskId = msg.TaskId

		reply, err := p.DeregisterMetadataMapping(ctx, &x)

		if err != nil {
			p.taskStatusMgr.PutStatus(pbtypes.SubTaskStatus{
				TaskId: x.TaskId,
				Status: constants.StatusFailed,
			})
			return nil, err
		}

		p.taskStatusMgr.PutStatus(pbtypes.SubTaskStatus{
			TaskId: x.TaskId,
			Status: constants.StatusSuccessful,
		})
		return reply, nil

	case pbtypes.SubTaskAction_RegisterCmd.String():

		var x pbtypes.SubTask_RegisterCmd
		err := json.Unmarshal([]byte(msg.Directive), &x)
		if err != nil {
			logger.Warn(nil, "%+v", err)
			return nil, err
		}

		x.Action = msg.Action
		x.TaskId = msg.TaskId

		reply, err := p.RegisterCmd(ctx, &x)

		if err != nil {
			p.taskStatusMgr.PutStatus(pbtypes.SubTaskStatus{
				TaskId: x.TaskId,
				Status: constants.StatusFailed,
			})
			return nil, err
		}

		p.taskStatusMgr.PutStatus(pbtypes.SubTaskStatus{
			TaskId: x.TaskId,
			Status: constants.StatusSuccessful,
		})
		return reply, nil

	case pbtypes.SubTaskAction_DeregisterCmd.String():

		var x pbtypes.SubTask_DeregisterCmd
		err := json.Unmarshal([]byte(msg.Directive), &x)
		if err != nil {
			logger.Warn(nil, "%+v", err)
			return nil, err
		}

		x.Action = msg.Action
		x.TaskId = msg.TaskId

		reply, err := p.DeregisterCmd(ctx, &x)
		if err != nil {
			p.taskStatusMgr.PutStatus(pbtypes.SubTaskStatus{
				TaskId: x.TaskId,
				Status: constants.StatusFailed,
			})
			return nil, err
		}

		p.taskStatusMgr.PutStatus(pbtypes.SubTaskStatus{
			TaskId: x.TaskId,
			Status: constants.StatusSuccessful,
		})
		return reply, nil

	case pbtypes.SubTaskAction_GetTaskStatus.String():

		var x pbtypes.SubTask_GetTaskStatus
		err := json.Unmarshal([]byte(msg.Directive), &x)
		if err != nil {
			logger.Warn(nil, "%+v", err)
			return nil, err
		}

		x.Action = msg.Action
		x.TaskId = msg.TaskId

		_, err = p.GetSubtaskStatus(ctx, &pbtypes.SubTaskId{})
		if err != nil {
			logger.Warn(nil, "%+v", err)
			return nil, err
		}

		return &pbtypes.Empty{}, err

	default:
		err := fmt.Errorf("pilot: unknown action: %s", msg.Action)
		logger.Warn(nil, "%+v", err)
		return nil, err
	}
}

func (p *Server) PingPilot(ctx context.Context, arg *pbtypes.Empty) (*pbtypes.Empty, error) {
	logger.Info(nil, funcutil.CallerName(1))

	return &pbtypes.Empty{}, nil
}

func (p *Server) PingFrontgate(ctx context.Context, arg *pbtypes.FrontgateId) (*pbtypes.Empty, error) {
	logger.Info(nil, funcutil.CallerName(1))

	client, err := p.fgClientMgr.GetClient(arg.Id)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	defer func() {
		if err != nil && p.fgClientMgr.IsFrontgateShutdownError(err) {
			p.fgClientMgr.CloseClient(client.info.Id, client.info.NodeId)
		}
	}()

	_, err = client.PingFrontgate(&pbtypes.Empty{})
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	return &pbtypes.Empty{}, nil
}

func (p *Server) PingFrontgateNode(ctx context.Context, arg *pbtypes.FrontgateNodeId) (*pbtypes.Empty, error) {
	logger.Info(nil, funcutil.CallerName(1))

	client, err := p.fgClientMgr.GetNodeClient(arg.Id, arg.NodeId)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	defer func() {
		if err != nil && p.fgClientMgr.IsFrontgateShutdownError(err) {
			p.fgClientMgr.CloseClient(client.info.Id, client.info.NodeId)
		}
	}()

	_, err = client.PingFrontgateNode(&pbtypes.Empty{})
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	return &pbtypes.Empty{}, nil
}

func (p *Server) PingDrone(ctx context.Context, arg *pbtypes.DroneEndpoint) (*pbtypes.Empty, error) {
	logger.Info(nil, funcutil.CallerName(1))

	client, err := p.fgClientMgr.GetClient(arg.FrontgateId)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	defer func() {
		if err != nil && p.fgClientMgr.IsFrontgateShutdownError(err) {
			p.fgClientMgr.CloseClient(client.info.Id, client.info.NodeId)
		}
	}()

	_, err = client.PingDrone(&pbtypes.DroneEndpoint{
		FrontgateId: arg.FrontgateId,
		DroneIp:     arg.DroneIp,
		DronePort:   arg.DronePort,
	})
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	return &pbtypes.Empty{}, nil
}

func (p *Server) RunCommandOnFrontgateNode(ctx context.Context, arg *pbtypes.RunCommandOnFrontgateRequest) (*pbtypes.String, error) {
	logger.Info(nil, funcutil.CallerName(1))

	client, err := p.fgClientMgr.GetNodeClient(
		arg.GetEndpoint().GetFrontgateId(),
		arg.GetEndpoint().GetFrontgateNodeId(),
	)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	defer func() {
		if err != nil && p.fgClientMgr.IsFrontgateShutdownError(err) {
			p.fgClientMgr.CloseClient(client.info.Id, client.info.NodeId)
		}
	}()

	reply, err := client.RunCommand(arg)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	return reply, nil
}

func (p *Server) RunCommandOnDrone(ctx context.Context, arg *pbtypes.RunCommandOnDroneRequest) (*pbtypes.String, error) {
	logger.Info(nil, funcutil.CallerName(1))

	client, err := p.fgClientMgr.GetClient(arg.GetEndpoint().GetFrontgateId())
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	defer func() {
		if err != nil && p.fgClientMgr.IsFrontgateShutdownError(err) {
			p.fgClientMgr.CloseClient(client.info.Id, client.info.NodeId)
		}
	}()

	reply, err := client.RunCommandOnDrone(arg)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	return reply, nil
}
