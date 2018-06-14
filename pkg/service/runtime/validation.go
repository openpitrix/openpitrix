package runtime

import (
	"net/url"
	"regexp"

	"github.com/asaskevich/govalidator"
	"github.com/ghodss/yaml"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/plugins"
	"openpitrix.io/openpitrix/pkg/util/stringutil"
)

func ValidateName(name string) error {
	if !govalidator.StringLength(name, NameMinLength, NameMaxLength) {
		return gerr.New(gerr.InvalidArgument, gerr.ErrorIllegalParameterLength, "name")
	}
	return nil
}

func ValidateProvider(provider string) error {
	if !govalidator.StringLength(provider, ProviderMinLength, ProviderMaxLength) {
		return gerr.New(gerr.InvalidArgument, gerr.ErrorIllegalParameterLength, "provider")
	}
	if i := stringutil.FindString(constants.VmBaseProviders, provider); i != -1 {
		return nil
	}
	if constants.ProviderKubernetes == provider {
		return nil
	}
	return gerr.New(gerr.InvalidArgument, gerr.ErrorUnsupportedParameterValue, "provider", provider)
}

func ValidateURL(url string) error {
	if !govalidator.IsURL(url) {
		return gerr.New(gerr.InvalidArgument, gerr.ErrorIllegalUrlFormat, url)
	}
	return nil
}

func ValidateCredential(provider, url, credential string) error {
	if len(credential) < CredentialMinLength {
		return gerr.New(gerr.InvalidArgument, gerr.ErrorIllegalParameterLength, "credential")
	}
	if i := stringutil.FindString(constants.VmBaseProviders, provider); i != -1 {
		err := ValidateURL(url)
		if err != nil {
			return err
		}
		_, err = yaml.JSONToYAML([]byte(credential))
		if err != nil {
			return err
		}
	}
	if constants.ProviderKubernetes == provider {
		_, err := yaml.YAMLToJSON([]byte(credential))
		if err != nil {
			return err
		}
	}

	providerInterface, err := plugins.GetProviderPlugin(provider, nil)
	if err != nil {
		logger.Error("No such provider [%s]. ", provider)
		return gerr.NewWithDetail(gerr.NotFound, err, gerr.ErrorProviderNotFound, provider)
	}
	return providerInterface.ValidateCredential(url, credential)
}

func ValidateZone(zone string) error {
	if !govalidator.StringLength(zone, ZoneMinLength, ZoneMaxLength) {
		return gerr.New(gerr.InvalidArgument, gerr.ErrorIllegalParameterLength, "zone")
	}
	return nil
}

func ValidateLabelString(labelString string) error {
	mapLabel, err := url.ParseQuery(labelString)
	if err != nil {
		return err
	}
	err = ValidateLabelMapFmt(mapLabel)
	if err != nil {
		return err
	}
	return nil
}

func ValidateSelectorString(selectorString string) error {
	selectorMap, err := url.ParseQuery(selectorString)
	if err != nil {
		return err
	}
	err = ValidateSelectorMapFmt(selectorMap)
	if err != nil {
		return err
	}
	return nil
}

func ValidateSelectorMapFmt(selectorMap map[string][]string) error {
	for sKey, sValues := range selectorMap {
		err := ValidateLabelKey(sKey)
		if err != nil {
			return err
		}
		for _, sValue := range sValues {
			err := ValidateLabelValue(sValue)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func ValidateLabelMapFmt(labelMap map[string][]string) error {
	for mKey, mValue := range labelMap {
		if len(mValue) != 1 {
			return gerr.New(gerr.InvalidArgument, gerr.ErrorIllegalLabelFormat)
		}
		err := ValidateLabelKey(mKey)
		if err != nil {
			return err
		}
		err = ValidateLabelValue(mValue[0])
		if err != nil {
			return err
		}
	}
	return nil
}

var LabelNameRegexp = regexp.MustCompile(LabelKeyFmt)

func ValidateLabelKey(labelName string) error {
	if !govalidator.StringLength(labelName, LabelKeyMinLength, LabelKeyMaxLength) {
		return gerr.New(gerr.InvalidArgument, gerr.ErrorIllegalParameterLength, "label_key")
	}
	if !LabelNameRegexp.Match([]byte(labelName)) {
		return gerr.New(gerr.InvalidArgument, gerr.ErrorIllegalLabelFormat)
	}
	return nil
}

func ValidateLabelValue(labelValue string) error {
	if !govalidator.StringLength(labelValue, LabelValueMinLength, LabelValueMaxLength) {
		return gerr.New(gerr.InvalidArgument, gerr.ErrorIllegalParameterLength, "label_value")
	}
	return nil
}

func validateCreateRuntimeRequest(req *pb.CreateRuntimeRequest) error {
	err := ValidateName(req.Name.GetValue())
	if err != nil {
		return err
	}
	err = ValidateLabelString(req.Labels.GetValue())
	if err != nil {
		return err
	}
	err = ValidateProvider(req.Provider.GetValue())
	if err != nil {
		return err
	}
	err = ValidateCredential(req.Provider.GetValue(),
		req.RuntimeUrl.GetValue(), req.RuntimeCredential.GetValue())
	if err != nil {
		return err
	}
	err = ValidateZone(req.Zone.GetValue())
	if err != nil {
		return err
	}
	return nil
}

func validateModifyRuntimeRequest(req *pb.ModifyRuntimeRequest) error {
	err := ValidateLabelString(req.Labels.GetValue())
	if err != nil {
		return err
	}
	return nil
}

func validateDescribeRuntimesRequest(req *pb.DescribeRuntimesRequest) error {
	err := ValidateSelectorString(req.Label.GetValue())
	if err != nil {
		return err
	}
	return nil
}
