package helm

import (
	providerclient "openpitrix.io/openpitrix/pkg/client/runtime_provider"
	"openpitrix.io/openpitrix/pkg/config"
	"openpitrix.io/openpitrix/pkg/manager"
	runtimeprovider "openpitrix.io/openpitrix/pkg/service/runtime_provider"

	"google.golang.org/grpc"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
)

type Server struct {
	runtimeprovider.Server
}

func Serve(cfg *config.Config) {
	pi.SetGlobal(cfg)
	err := providerclient.RegisterRuntimeProvider(Provider, ProviderConfig)
	if err != nil {
		logger.Critical(nil, "failed to register provider config: %+v", err)
	}
	s := Server{}
	manager.NewGrpcServer("openpitrix-rp-kubernetes", constants.KubernetesProviderPort).
		ShowErrorCause(cfg.Grpc.ShowErrorCause).
		Serve(func(server *grpc.Server) {
			pb.RegisterRuntimeProviderManagerServer(server, &s)
		})
}
