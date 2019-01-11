package runtime

import (
	"context"
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/ghodss/yaml"

	providerclient "openpitrix.io/openpitrix/pkg/client/runtime_provider"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/plugins"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

func ValidateURL(ctx context.Context, url string) error {
	if !govalidator.IsURL(url) {
		return gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorIllegalUrlFormat, url)
	}
	return nil
}

func ValidateRuntime(ctx context.Context, runtimeId, zone string, runtimeCredential *models.RuntimeCredential, needCreate bool) error {
	if plugins.IsVmbasedProviders(runtimeCredential.Provider) {
		err := ValidateURL(ctx, runtimeCredential.RuntimeUrl)
		if err != nil {
			return err
		}
		_, err = yaml.JSONToYAML([]byte(runtimeCredential.RuntimeCredentialContent))
		if err != nil {
			return err
		}
	} else if constants.ProviderKubernetes == runtimeCredential.Provider {
		_, err := yaml.YAMLToJSON([]byte(runtimeCredential.RuntimeCredentialContent))
		if err != nil {
			return err
		}
	}

	providerClient, err := providerclient.NewRuntimeProviderManagerClient()
	if err != nil {
		return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	response, err := providerClient.ValidateRuntime(ctx, &pb.ValidateRuntimeRequest{
		RuntimeId:         pbutil.ToProtoString(runtimeId),
		Zone:              pbutil.ToProtoString(zone),
		RuntimeCredential: models.RuntimeCredentialToPb(runtimeCredential),
		NeedCreate:        pbutil.ToProtoBool(needCreate),
	})
	if err != nil {
		return err
	}
	if !response.Ok.GetValue() {
		return fmt.Errorf("response is not ok")
	}
	return nil
}
