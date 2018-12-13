package runtime

import (
	"context"

	"github.com/asaskevich/govalidator"
	"github.com/ghodss/yaml"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/plugins"
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

	providerInterface, err := plugins.GetProviderPlugin(ctx, runtimeCredential.Provider)
	if err != nil {
		logger.Error(ctx, "No such provider [%s]. ", runtimeCredential.Provider)
		return gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorProviderNotFound, runtimeCredential.Provider)
	}
	return providerInterface.ValidateRuntime(ctx, runtimeId, zone, runtimeCredential, needCreate)
}
