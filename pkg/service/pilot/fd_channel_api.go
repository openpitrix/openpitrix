// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package pilot

import (
	"bytes"
	"context"
	"io"
	"sync"

	google_protobuf "github.com/golang/protobuf/ptypes/wrappers"
	grpc "google.golang.org/grpc"

	pb "openpitrix.io/openpitrix/pkg/pb"
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

func NewFrameChannel(
	recv func() ([]byte, error), send func(data []byte) error,
	close func() error,
) *FrameChannel {
	return &FrameChannel{
		RecvMsgFunc: recv, SendMsgFunc: send,
		CloseFunc: close,
	}
}

func DialFrontgateChannel(
	ctx context.Context, target string,
	opts ...grpc.DialOption,
) (
	ch *FrameChannel, conn *grpc.ClientConn,
	err error,
) {
	conn, err = grpc.Dial(target, opts...)
	if err != nil {
		return
	}

	channel, err := pb.NewPilotManagerClient(conn).FrontgateChannel(ctx)
	if err != nil {
		conn.Close()
		conn = nil
		return
	}

	ch = NewFrontgateChannelFromClient(channel)
	return
}

func NewFrontgateChannelFromServer(ch pb.PilotManager_FrontgateChannelServer) *FrameChannel {
	return &FrameChannel{
		RecvMsgFunc: func() ([]byte, error) {
			msg, err := ch.Recv()
			if err != nil {
				return nil, err
			}
			return msg.GetValue(), nil
		},
		SendMsgFunc: func(data []byte) error {
			return ch.Send(&google_protobuf.BytesValue{Value: data})
		},
		CloseFunc: func() error {
			return nil
		},
	}
}

func NewFrontgateChannelFromClient(ch pb.PilotManager_FrontgateChannelClient) *FrameChannel {
	return &FrameChannel{
		RecvMsgFunc: func() ([]byte, error) {
			msg, err := ch.Recv()
			if err != nil {
				return nil, err
			}
			return msg.GetValue(), nil
		},
		SendMsgFunc: func(data []byte) error {
			return ch.Send(&google_protobuf.BytesValue{Value: data})
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

	n, err = r.Read(data)
	return n, nil
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
