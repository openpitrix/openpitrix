// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package pilot

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/golang/protobuf/proto"

	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pb/frontgate"
	"openpitrix.io/openpitrix/pkg/pb/pilot"
	"openpitrix.io/openpitrix/pkg/pb/types"
	"openpitrix.io/openpitrix/pkg/service/pilot/pilotutil"
)

var (
	_ pbpilot.PilotServiceServer = (*Server)(nil)
)

func (p *Server) GetPilotConfig(context.Context, *pbtypes.Empty) (*pbtypes.PilotConfig, error) {
	return proto.Clone(p.cfg).(*pbtypes.PilotConfig), nil
}

func (p *Server) GetFrontgateList(context.Context, *pbtypes.Empty) (*pbtypes.FrontgateIdList, error) {
	return nil, fmt.Errorf("TODO")
}

func (p *Server) GetFrontgateConfig(ctx context.Context, arg *pbtypes.FrontgateId) (*pbtypes.FrontgateConfig, error) {
	client, err := p.fgClientMgr.GetClient(arg.Id)
	if err != nil {
		return nil, err
	}

	reply, err := client.GetFrontgateConfig(&pbtypes.Empty{})
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func (p *Server) SetFrontgateConfig(ctx context.Context, arg *pbtypes.FrontgateConfig) (*pbtypes.Empty, error) {
	client, err := p.fgClientMgr.GetClient(arg.Id)
	if err != nil {
		return nil, err
	}

	reply, err := client.SetFrontgateConfig(arg)
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func (p *Server) GetDroneConfig(ctx context.Context, arg *pbtypes.DroneEndpoint) (*pbtypes.DroneConfig, error) {
	client, err := p.fgClientMgr.GetClient(arg.FrontgateId)
	if err != nil {
		return nil, err
	}

	reply, err := client.GetDroneConfig(arg)
	if err != nil {
		return nil, err
	}

	return reply, nil
}
func (p *Server) SetDroneConfig(ctx context.Context, arg *pbtypes.SetDroneConfigRequest) (*pbtypes.Empty, error) {
	client, err := p.fgClientMgr.GetClient(arg.Endpoint.FrontgateId)
	if err != nil {
		return nil, err
	}

	_, err = client.SetDroneConfig(arg)
	if err != nil {
		return nil, err
	}

	return &pbtypes.Empty{}, nil
}

func (p *Server) FrontgateChannel(ch pbpilot.PilotService_FrontgateChannelServer) error {
	c := pbfrontgate.NewFrontgateServiceClient(
		pilotutil.NewFrontgateChannelFromServer(ch),
	)

	info, err := c.GetFrontgateConfig(&pbtypes.Empty{})
	if err != nil {
		logger.Debug("%+v", err)
		return err
	}

	// if return, the channel will be closed
	<-p.fgClientMgr.PutClient(c, info)
	return nil
}

func (p *Server) GetConfdConfig(ctx context.Context, arg *pbtypes.ConfdEndpoint) (*pbtypes.ConfdConfig, error) {
	client, err := p.fgClientMgr.GetClient(arg.FrontgateId)
	if err != nil {
		return nil, err
	}

	reply, err := client.GetConfdConfig(&pbtypes.ConfdEndpoint{
		DroneIp:   arg.DroneIp,
		DronePort: arg.DronePort,
	})
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func (p *Server) IsConfdRunning(ctx context.Context, arg *pbtypes.DroneEndpoint) (*pbtypes.Bool, error) {
	client, err := p.fgClientMgr.GetClient(arg.FrontgateId)
	if err != nil {
		return nil, err
	}

	reply, err := client.IsConfdRunning(&pbtypes.ConfdEndpoint{
		DroneIp:   arg.DroneIp,
		DronePort: arg.DronePort,
	})
	if err != nil {
		return nil, err
	}

	return &pbtypes.Bool{Value: reply.GetValue()}, nil
}

func (p *Server) StartConfd(ctx context.Context, arg *pbtypes.DroneEndpoint) (*pbtypes.Empty, error) {
	client, err := p.fgClientMgr.GetClient(arg.FrontgateId)
	if err != nil {
		return nil, err
	}

	_, err = client.StartConfd(&pbtypes.ConfdEndpoint{
		DroneIp:   arg.DroneIp,
		DronePort: arg.DronePort,
	})
	if err != nil {
		return nil, err
	}

	return &pbtypes.Empty{}, nil
}

func (p *Server) StopConfd(ctx context.Context, arg *pbtypes.DroneEndpoint) (*pbtypes.Empty, error) {
	client, err := p.fgClientMgr.GetClient(arg.FrontgateId)
	if err != nil {
		return nil, err
	}

	_, err = client.StopConfd(&pbtypes.ConfdEndpoint{
		DroneIp:   arg.DroneIp,
		DronePort: arg.DronePort,
	})
	if err != nil {
		return nil, err
	}

	return &pbtypes.Empty{}, nil
}

func (p *Server) RegisterMetadata(ctx context.Context, arg *pbtypes.SubTask_RegisterMetadata) (*pbtypes.Empty, error) {
	client, err := p.fgClientMgr.GetClient(arg.FrontgateId)
	if err != nil {
		return nil, err
	}

	_, err = client.RegisterMetadata(arg)
	if err != nil {
		return nil, err
	}

	return &pbtypes.Empty{}, nil
}

func (p *Server) DeregisterMetadata(ctx context.Context, arg *pbtypes.SubTask_DeregisterMetadata) (*pbtypes.Empty, error) {
	client, err := p.fgClientMgr.GetClient(arg.FrontgateId)
	if err != nil {
		return nil, err
	}

	_, err = client.DeregisterMetadata(arg)
	if err != nil {
		return nil, err
	}

	return &pbtypes.Empty{}, nil
}

func (p *Server) RegisterCmd(ctx context.Context, arg *pbtypes.SubTask_RegisterCmd) (*pbtypes.Empty, error) {
	client, err := p.fgClientMgr.GetClient(arg.FrontgateId)
	if err != nil {
		return nil, err
	}

	_, err = client.RegisterCmd(arg)
	if err != nil {
		return nil, err
	}

	return &pbtypes.Empty{}, nil
}

func (p *Server) DeregisterCmd(ctx context.Context, arg *pbtypes.SubTask_DeregisterCmd) (*pbtypes.Empty, error) {
	client, err := p.fgClientMgr.GetClient(arg.FrontgateId)
	if err != nil {
		return nil, err
	}

	_, err = client.DeregisterCmd(arg)
	if err != nil {
		return nil, err
	}

	return &pbtypes.Empty{}, nil
}

func (p *Server) ReportSubTaskStatus(ctx context.Context, arg *pbtypes.SubTaskStatus) (*pbtypes.Empty, error) {
	p.taskStatusMgr.PutStatus(*arg)
	return &pbtypes.Empty{}, nil
}

func (p *Server) GetSubtaskStatus(ctx context.Context, arg *pbtypes.SubTaskId) (*pbtypes.SubTaskStatus, error) {
	if s, ok := p.taskStatusMgr.GetStatus(arg.TaskId); ok {
		return &s, nil
	}

	return nil, fmt.Errorf("pilot: not found, taskId = %s", arg.TaskId)
}

func (p *Server) HandleSubtask(ctx context.Context, msg *pbtypes.SubTaskMessage) (*pbtypes.Empty, error) {
	switch msg.Action {
	case pbtypes.SubTaskAction_StartConfd.String():
		var x pbtypes.SubTask_StartConfd
		err := json.Unmarshal([]byte(msg.Directive), &x)
		if err != nil {
			return nil, err
		}

		x.Action = msg.Action
		x.TaskId = msg.TaskId

		return p.StartConfd(ctx, &pbtypes.DroneEndpoint{
			FrontgateId: x.FrontgateId,
			DroneIp:     x.DroneIp,
			DronePort:   9112,
		})

	case pbtypes.SubTaskAction_StopConfd.String():
		var x pbtypes.SubTask_StopConfd
		err := json.Unmarshal([]byte(msg.Directive), &x)
		if err != nil {
			return nil, err
		}

		x.Action = msg.Action
		x.TaskId = msg.TaskId

		return p.StopConfd(ctx, &pbtypes.DroneEndpoint{
			FrontgateId: x.FrontgateId,
			DroneIp:     x.DroneIp,
			DronePort:   9112,
		})

	case pbtypes.SubTaskAction_RegisterMetadata.String():

		var x pbtypes.SubTask_RegisterMetadata
		err := json.Unmarshal([]byte(msg.Directive), &x)
		if err != nil {
			return nil, err
		}

		x.Action = msg.Action
		x.TaskId = msg.TaskId

		return p.RegisterMetadata(ctx, &x)

	case pbtypes.SubTaskAction_DeregisterMetadata.String():

		var x pbtypes.SubTask_DeregisterMetadata
		err := json.Unmarshal([]byte(msg.Directive), &x)
		if err != nil {
			return nil, err
		}

		x.Action = msg.Action
		x.TaskId = msg.TaskId

		return p.DeregisterMetadata(ctx, &x)

	case pbtypes.SubTaskAction_RegisterCmd.String():

		var x pbtypes.SubTask_RegisterCmd
		err := json.Unmarshal([]byte(msg.Directive), &x)
		if err != nil {
			return nil, err
		}

		x.Action = msg.Action
		x.TaskId = msg.TaskId

		return p.RegisterCmd(ctx, &x)

	case pbtypes.SubTaskAction_DeregisterCmd.String():

		var x pbtypes.SubTask_DeregisterCmd
		err := json.Unmarshal([]byte(msg.Directive), &x)
		if err != nil {
			return nil, err
		}

		x.Action = msg.Action
		x.TaskId = msg.TaskId

		return p.DeregisterCmd(ctx, &x)

	case pbtypes.SubTaskAction_GetTaskStatus.String():

		var x pbtypes.SubTask_GetTaskStatus
		err := json.Unmarshal([]byte(msg.Directive), &x)
		if err != nil {
			return nil, err
		}

		x.Action = msg.Action
		x.TaskId = msg.TaskId

		_, err = p.GetSubtaskStatus(ctx, &pbtypes.SubTaskId{})
		if err != nil {
			return nil, err
		}

		return &pbtypes.Empty{}, err

	default:
		return nil, fmt.Errorf("pilot: unknown action: %s", msg.Action)
	}
}

func (p *Server) PingPilot(ctx context.Context, arg *pbtypes.Empty) (*pbtypes.Empty, error) {
	return &pbtypes.Empty{}, nil
}

func (p *Server) PingFrontgate(ctx context.Context, arg *pbtypes.FrontgateId) (*pbtypes.Empty, error) {
	client, err := p.fgClientMgr.GetClient(arg.Id)
	if err != nil {
		return nil, err
	}

	_, err = client.PingFrontgate(&pbtypes.Empty{})
	if err != nil {
		return nil, err
	}

	return &pbtypes.Empty{}, nil
}

func (p *Server) PingFrontgateNode(ctx context.Context, arg *pbtypes.FrontgateNodeId) (*pbtypes.Empty, error) {
	client, err := p.fgClientMgr.GetNodeClient(arg.Id, arg.NodeId)
	if err != nil {
		return nil, err
	}

	_, err = client.PingFrontgateNode(&pbtypes.Empty{})
	if err != nil {
		return nil, err
	}

	return &pbtypes.Empty{}, nil
}

func (p *Server) PingDrone(ctx context.Context, arg *pbtypes.DroneEndpoint) (*pbtypes.Empty, error) {
	client, err := p.fgClientMgr.GetClient(arg.FrontgateId)
	if err != nil {
		return nil, err
	}

	_, err = client.PingDrone(&pbtypes.DroneEndpoint{
		FrontgateId: arg.FrontgateId,
		DroneIp:     arg.DroneIp,
		DronePort:   arg.DronePort,
	})
	if err != nil {
		return nil, err
	}

	return &pbtypes.Empty{}, nil
}

func (p *Server) RunCommandOnFrontgateNode(ctx context.Context, arg *pbtypes.RunCommandOnFrontgateRequest) (*pbtypes.String, error) {
	client, err := p.fgClientMgr.GetNodeClient(
		arg.GetEndpoint().GetFrontgateId(),
		arg.GetEndpoint().GetFrontgateNodeId(),
	)
	if err != nil {
		return nil, err
	}

	reply, err := client.RunCommand(arg)
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func (p *Server) RunCommandOnDrone(ctx context.Context, arg *pbtypes.RunCommandOnDroneRequest) (*pbtypes.String, error) {
	client, err := p.fgClientMgr.GetClient(arg.GetEndpoint().GetFrontgateId())
	if err != nil {
		return nil, err
	}

	reply, err := client.RunCommandOnDrone(arg)
	if err != nil {
		return nil, err
	}

	return reply, nil
}
