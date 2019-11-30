package push

import (
	"google.golang.org/grpc"

	"openpitrix.io/openpitrix/pkg/config"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
)

type Server struct {
}

func Serve(cfg *config.Config) {
	pi.SetGlobal(cfg)
	s := Server{}
	manager.NewGrpcServer("push-manager", constants.ServicePushPort).
		ShowErrorCause(cfg.Grpc.ShowErrorCause).
		Serve(func(server *grpc.Server) {
			pb.RegisterStreamServer(server, &s)
		})
}
