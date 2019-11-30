package push

import (
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/topic"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

func (s *Server) CreateServerStream(req *pb.CreateStreamRequest, stream pb.Stream_CreateServerStreamServer) error {
	for {
		msg := <-topic.MsgChan
		jsonMsg, err := jsonutil.ToJson(msg.Message).Bytes()
		if err != nil {
			return err
		}
		resp := &pb.CreateStreamResponse{
			UserID:  pbutil.ToProtoString(msg.UserId),
			Message: pbutil.ToProtoBytes(jsonMsg),
		}
		err = stream.Send(resp)
		if err != nil {
			return err
		}
	}
}
