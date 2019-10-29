// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package pilotutil

import (
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"openpitrix.io/openpitrix/pkg/logger"
	pbpilot "openpitrix.io/openpitrix/pkg/pb/metadata/pilot"
	pbtypes "openpitrix.io/openpitrix/pkg/pb/metadata/types"
)

var (
	_ io.ReadWriteCloser = (*FrameChannel)(nil)
)

type FrameChannel struct {
	RecvMsgFunc func() ([]byte, error)
	SendMsgFunc func(data []byte) error
	CloseFunc   func() error

	r  *bytes.Reader
	mu sync.Mutex // only for update r reader
}

func DialPilotChannelTLS(
	ctx context.Context,
	target string,
	tlsConfig *tls.Config,
	opts ...grpc.DialOption,
) (*FrameChannel, *grpc.ClientConn, error) {
	creds := credentials.NewTLS(tlsConfig)
	opts = append([]grpc.DialOption{
		grpc.WithTransportCredentials(creds),
	}, opts...)

	logger.Info(nil, "Dial [%s] begin.", target)
	conn, err := grpc.DialContext(ctx, target, opts...)
	if err != nil {
		logger.Error(nil, "Dial [%s] failed: %+v", target, err)
		return nil, conn, err
	}
	logger.Info(nil, "Dial [%s] finished, create channel begin.", target)

	channel, err := pbpilot.NewPilotServiceForFrontgateClient(conn).FrontgateChannel(ctx)
	if err != nil {
		logger.Error(nil, "Create channel failed: %+v", err)
		conn.Close()
		conn = nil
		return nil, conn, err
	}

	logger.Info(nil, "Create channel finished.")

	ch := NewFrontgateChannelFromClient(channel)
	return ch, conn, nil
}

func NewFrontgateChannelFromServer(ch pbpilot.PilotServiceForFrontgate_FrontgateChannelServer) *FrameChannel {
	return &FrameChannel{
		RecvMsgFunc: func() ([]byte, error) {
			msg, err := ch.Recv()
			if err != nil {
				logger.Error(nil, "Receive error is: %+v", err)
				return nil, err
			}
			return msg.GetValue(), nil
		},
		SendMsgFunc: func(data []byte) error {
			return ch.Send(&pbtypes.Bytes{Value: data})
		},
		CloseFunc: func() error {
			return nil
		},
	}
}

func NewFrontgateChannelFromClient(ch pbpilot.PilotServiceForFrontgate_FrontgateChannelClient) *FrameChannel {
	return &FrameChannel{
		RecvMsgFunc: func() ([]byte, error) {
			msg, err := ch.Recv()
			if err != nil {
				logger.Error(nil, "Receive error is: %+v", err)
				return nil, err
			}
			return msg.GetValue(), nil
		},
		SendMsgFunc: func(data []byte) error {
			return ch.Send(&pbtypes.Bytes{Value: data})
		},
		CloseFunc: func() error {
			return ch.CloseSend()
		},
	}
}

func (p *FrameChannel) nextReader() (*bytes.Reader, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.r == nil || p.r.Len() == 0 {
		msg, err := p.RecvMsgFunc()
		if err != nil {
			p.r = nil

			return nil, err
		}
		p.r = bytes.NewReader(msg)
	}

	return p.r, nil
}

func (p *FrameChannel) Read(data []byte) (n int, err error) {
	r, err := p.nextReader()
	if err != nil {
		return 0, err
	}

	return r.Read(data)
}

func (p *FrameChannel) Write(data []byte) (n int, err error) {
	if err = p.SendMsgFunc(data); err != nil {
		return 0, err
	}
	return len(data), nil
}

func (p *FrameChannel) Close() error {
	if p.CloseFunc != nil {
		return p.CloseFunc()
	}
	return nil
}
