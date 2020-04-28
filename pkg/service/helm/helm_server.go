package helm

import (
	"google.golang.org/grpc"

	"openpitrix.io/openpitrix/pkg/config"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
)

type HelmServer struct {
}

func HelmServe(cfg *config.Config) {
	pi.SetGlobal(cfg)
	s := HelmServer{}
	manager.NewGrpcServer("openpitrix-release-manager", constants.ReleaseManagerPort).
		ShowErrorCause(cfg.Grpc.ShowErrorCause).
		Serve(func(server *grpc.Server) {
			pb.RegisterReleaseManagerServer(server, &s)
		})
}
