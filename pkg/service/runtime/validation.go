package runtime

import (
	"context"
	"net/url"
	"regexp"

	"github.com/asaskevich/govalidator"
	"github.com/ghodss/yaml"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/plugins"
)

func ValidateName(ctx context.Context, name string) error {
	if !govalidator.StringLength(name, NameMinLength, NameMaxLength) {
		return gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorIllegalParameterLength, "name")
	}
	return nil
}

func ValidateURL(ctx context.Context, url string) error {
	if !govalidator.IsURL(url) {
		return gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorIllegalUrlFormat, url)
	}
	return nil
}

func ValidateCredential(ctx context.Context, provider, url, credential, zone string) error {
	if len(credential) < CredentialMinLength {
		return gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorIllegalParameterLength, "credential")
	}
	if plugins.IsVmbasedProviders(provider) {
		err := ValidateURL(ctx, url)
		if err != nil {
			return err
		}
		_, err = yaml.JSONToYAML([]byte(credential))
		if err != nil {
			return err
		}
	} else if constants.ProviderKubernetes == provider {
		_, err := yaml.YAMLToJSON([]byte(credential))
		if err != nil {
			return err
		}
	}

	providerInterface, err := plugins.GetProviderPlugin(ctx, provider)
	if err != nil {
		logger.Error(ctx, "No such provider [%s]. ", provider)
		return gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorProviderNotFound, provider)
	}
	return providerInterface.ValidateCredential(ctx, url, credential, zone)
}

func ValidateZone(ctx context.Context, zone string) error {
	if !govalidator.StringLength(zone, ZoneMinLength, ZoneMaxLength) {
		return gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorIllegalParameterLength, "zone")
	}
	return nil
}

func ValidateLabelString(ctx context.Context, labelString string) error {
	mapLabel, err := url.ParseQuery(labelString)
	if err != nil {
		return err
	}
	err = ValidateLabelMapFmt(ctx, mapLabel)
	if err != nil {
		return err
	}
	return nil
}

func ValidateSelectorString(ctx context.Context, selectorString string) error {
	selectorMap, err := url.ParseQuery(selectorString)
	if err != nil {
		return err
	}
	err = ValidateSelectorMapFmt(ctx, selectorMap)
	if err != nil {
		return err
	}
	return nil
}

func ValidateSelectorMapFmt(ctx context.Context, selectorMap map[string][]string) error {
	for sKey, sValues := range selectorMap {
		err := ValidateLabelKey(ctx, sKey)
		if err != nil {
			return err
		}
		for _, sValue := range sValues {
			err := ValidateLabelValue(ctx, sValue)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func ValidateLabelMapFmt(ctx context.Context, labelMap map[string][]string) error {
	for mKey, mValue := range labelMap {
		if len(mValue) != 1 {
			return gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorIllegalLabelFormat)
		}
		err := ValidateLabelKey(ctx, mKey)
		if err != nil {
			return err
		}
		err = ValidateLabelValue(ctx, mValue[0])
		if err != nil {
			return err
		}
	}
	return nil
}

var LabelNameRegexp = regexp.MustCompile(LabelKeyFmt)

func ValidateLabelKey(ctx context.Context, labelName string) error {
	if !govalidator.StringLength(labelName, LabelKeyMinLength, LabelKeyMaxLength) {
		return gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorIllegalParameterLength, "label_key")
	}
	if !LabelNameRegexp.Match([]byte(labelName)) {
		return gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorIllegalLabelFormat)
	}
	return nil
}

func ValidateLabelValue(ctx context.Context, labelValue string) error {
	if !govalidator.StringLength(labelValue, LabelValueMinLength, LabelValueMaxLength) {
		return gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorIllegalParameterLength, "label_value")
	}
	return nil
}

func validateCreateRuntimeRequest(ctx context.Context, req *pb.CreateRuntimeRequest) error {
	err := ValidateName(ctx, req.Name.GetValue())
	if err != nil {
		return err
	}
	err = ValidateLabelString(ctx, req.Labels.GetValue())
	if err != nil {
		return err
	}
	err = ValidateCredential(
		ctx,
		req.Provider.GetValue(),
		req.RuntimeUrl.GetValue(),
		req.RuntimeCredential.GetValue(),
		req.GetZone().GetValue())
	if err != nil {
		return err
	}
	err = ValidateZone(ctx, req.Zone.GetValue())
	if err != nil {
		return err
	}
	return nil
}

func validateModifyRuntimeRequest(ctx context.Context, req *pb.ModifyRuntimeRequest) error {
	err := ValidateLabelString(ctx, req.Labels.GetValue())
	if err != nil {
		return err
	}
	return nil
}

func validateDescribeRuntimesRequest(ctx context.Context, req *pb.DescribeRuntimesRequest) error {
	err := ValidateSelectorString(ctx, req.Label.GetValue())
	if err != nil {
		return err
	}
	return nil
}
