package helm

import (
	"bufio"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	"k8s.io/api/apps/v1beta1"
	"k8s.io/api/apps/v1beta2"
	exv1beta1 "k8s.io/api/extensions/v1beta1"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/engine"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	_ "k8s.io/kubernetes/pkg/apis/apps/install"
	_ "k8s.io/kubernetes/pkg/apis/extensions/install"

	clientutil "openpitrix.io/openpitrix/pkg/client"
	appclient "openpitrix.io/openpitrix/pkg/client/app"
	runtimeclient "openpitrix.io/openpitrix/pkg/client/runtime"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
	"openpitrix.io/openpitrix/pkg/util/stringutil"
)

type Parser struct {
	Logger *logger.Logger
}

func (p *Parser) ParseCluster(name string, description string, versionId string) (*models.Cluster, error) {
	ctx := clientutil.GetSystemUserContext()
	appManagerClient, err := appclient.NewAppManagerClient()
	if err != nil {
		return nil, err
	}

	req := &pb.DescribeAppVersionsRequest{
		VersionId: []string{versionId},
	}

	resp, err := appManagerClient.DescribeAppVersions(ctx, req)
	if err != nil {
		return nil, err
	}

	if len(resp.AppVersionSet) == 0 {
		return nil, fmt.Errorf("app version [%s] not found", versionId)
	}

	appVersion := resp.AppVersionSet[0]

	cluster := &models.Cluster{
		Name:        name,
		Description: description,
		AppId:       appVersion.AppId.GetValue(),
		VersionId:   versionId,
		CreateTime:  time.Now(),
		StatusTime:  time.Now(),
	}

	return cluster, nil
}

func (p *Parser) ParseClusterRolesAndClusterCommons(
	c *chart.Chart,
	vals map[string]interface{},
	customVals map[string]interface{},
	runtimeId string,
) (
	map[string]*models.ClusterRole,
	map[string]*models.ClusterCommon,
	error,
) {
	env := jsonutil.ToString(customVals)

	renderer := engine.New()
	out, err := renderer.Render(c, vals)
	if err != nil {
		return nil, nil, err
	}

	var apiVersions []string
	decode := legacyscheme.Codecs.UniversalDeserializer().Decode

	clusterRoles := map[string]*models.ClusterRole{}
	clusterCommons := map[string]*models.ClusterCommon{}
	for k, v := range out {
		if filepath.Ext(k) != ".yaml" {
			continue
		}

		if len([]byte(v)) == 0 {
			continue
		}
		b := bufio.NewReader(strings.NewReader(v))
		r := k8syaml.NewYAMLReader(b)
		for {
			doc, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				p.Logger.Error("Decode file [%s] in chart failed, %+v", k, err)
				return nil, nil, err
			}
			obj, groupVersionKind, err := decode(doc, nil, nil)

			if err != nil {
				p.Logger.Error("Decode file [%s] in chart failed, %+v", k, err)
				return nil, nil, err
			}
			p.Logger.Debug("Yaml content: %+v", obj)
			p.Logger.Debug("Group version: %+v", groupVersionKind.GroupVersion().String())

			apiVersions = append(apiVersions, groupVersionKind.GroupVersion().String())

			switch o := obj.(type) {
			case *v1beta2.Deployment:
				clusterRole := &models.ClusterRole{
					Role: fmt.Sprintf("%s-Deployment", o.GetObjectMeta().GetName()),
					Env:  string(env),
				}

				if o.Spec.Replicas == nil {
					clusterRole.InstanceSize = 1
				} else {
					clusterRole.InstanceSize = uint32(*o.Spec.Replicas)
				}

				if len(o.Spec.Template.Spec.Containers) > 0 {
					clusterRole.Cpu = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().Value())
					clusterRole.Gpu = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.NvidiaGPU().Value())
					clusterRole.Memory = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().Value() / 1024 / 1024 / 1024)
					clusterRole.StorageSize = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.StorageEphemeral().Value() / 1024 / 1024 / 1024)
				}

				clusterCommon := &models.ClusterCommon{
					Role:       clusterRole.Role,
					Hypervisor: "docker",
				}

				clusterRoles[clusterRole.Role] = clusterRole
				clusterCommons[clusterRole.Role] = clusterCommon
			case *v1beta1.Deployment:
				clusterRole := &models.ClusterRole{
					Role: fmt.Sprintf("%s-Deployment", o.GetObjectMeta().GetName()),
					Env:  string(env),
				}

				if o.Spec.Replicas == nil {
					clusterRole.InstanceSize = 1
				} else {
					clusterRole.InstanceSize = uint32(*o.Spec.Replicas)
				}

				if len(o.Spec.Template.Spec.Containers) > 0 {
					clusterRole.Cpu = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().Value())
					clusterRole.Gpu = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.NvidiaGPU().Value())
					clusterRole.Memory = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().Value() / 1024 / 1024 / 1024)
					clusterRole.StorageSize = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.StorageEphemeral().Value() / 1024 / 1024 / 1024)
				}

				clusterCommon := &models.ClusterCommon{
					Role:       clusterRole.Role,
					Hypervisor: "docker",
				}

				clusterRoles[clusterRole.Role] = clusterRole
				clusterCommons[clusterRole.Role] = clusterCommon
			case *exv1beta1.Deployment:
				clusterRole := &models.ClusterRole{
					Role: fmt.Sprintf("%s-Deployment", o.GetObjectMeta().GetName()),
					Env:  string(env),
				}

				if o.Spec.Replicas == nil {
					clusterRole.InstanceSize = 1
				} else {
					clusterRole.InstanceSize = uint32(*o.Spec.Replicas)
				}

				if len(o.Spec.Template.Spec.Containers) > 0 {
					clusterRole.Cpu = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().Value())
					clusterRole.Gpu = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.NvidiaGPU().Value())
					clusterRole.Memory = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().Value() / 1024 / 1024 / 1024)
					clusterRole.StorageSize = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.StorageEphemeral().Value() / 1024 / 1024 / 1024)
				}

				clusterCommon := &models.ClusterCommon{
					Role:       clusterRole.Role,
					Hypervisor: "docker",
				}

				clusterRoles[clusterRole.Role] = clusterRole
				clusterCommons[clusterRole.Role] = clusterCommon
			case *v1beta2.StatefulSet:
				clusterRole := &models.ClusterRole{
					Role: fmt.Sprintf("%s-StatefulSet", o.GetObjectMeta().GetName()),
					Env:  string(env),
				}

				if o.Spec.Replicas == nil {
					clusterRole.InstanceSize = 1
				} else {
					clusterRole.InstanceSize = uint32(*o.Spec.Replicas)
				}

				if len(o.Spec.Template.Spec.Containers) > 0 {
					clusterRole.Cpu = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().Value())
					clusterRole.Gpu = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.NvidiaGPU().Value())
					clusterRole.Memory = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().Value() / 1024 / 1024 / 1024)
					clusterRole.StorageSize = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.StorageEphemeral().Value() / 1024 / 1024 / 1024)
				}

				clusterCommon := &models.ClusterCommon{
					Role:       clusterRole.Role,
					Hypervisor: "docker",
				}

				clusterRoles[clusterRole.Role] = clusterRole
				clusterCommons[clusterRole.Role] = clusterCommon
			case *v1beta1.StatefulSet:
				clusterRole := &models.ClusterRole{
					Role: fmt.Sprintf("%s-StatefulSet", o.GetObjectMeta().GetName()),
					Env:  string(env),
				}

				if o.Spec.Replicas == nil {
					clusterRole.InstanceSize = 1
				} else {
					clusterRole.InstanceSize = uint32(*o.Spec.Replicas)
				}

				if len(o.Spec.Template.Spec.Containers) > 0 {
					clusterRole.Cpu = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().Value())
					clusterRole.Gpu = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.NvidiaGPU().Value())
					clusterRole.Memory = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().Value() / 1024 / 1024 / 1024)
					clusterRole.StorageSize = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.StorageEphemeral().Value() / 1024 / 1024 / 1024)
				}

				clusterCommon := &models.ClusterCommon{
					Role:       clusterRole.Role,
					Hypervisor: "docker",
				}

				clusterRoles[clusterRole.Role] = clusterRole
				clusterCommons[clusterRole.Role] = clusterCommon
			case *v1beta2.DaemonSet:
				clusterRole := &models.ClusterRole{
					Role:         fmt.Sprintf("%s-DaemonSet", o.GetObjectMeta().GetName()),
					InstanceSize: uint32(1),
					Env:          string(env),
				}

				if len(o.Spec.Template.Spec.Containers) > 0 {
					clusterRole.Cpu = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().Value())
					clusterRole.Gpu = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.NvidiaGPU().Value())
					clusterRole.Memory = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().Value() / 1024 / 1024 / 1024)
					clusterRole.StorageSize = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.StorageEphemeral().Value() / 1024 / 1024 / 1024)
				}

				clusterCommon := &models.ClusterCommon{
					Role:       clusterRole.Role,
					Hypervisor: "docker",
				}

				clusterRoles[clusterRole.Role] = clusterRole
				clusterCommons[clusterRole.Role] = clusterCommon
			case *exv1beta1.DaemonSet:
				clusterRole := &models.ClusterRole{
					Role:         fmt.Sprintf("%s-DaemonSet", o.GetObjectMeta().GetName()),
					InstanceSize: uint32(1),
					Env:          string(env),
				}

				if len(o.Spec.Template.Spec.Containers) > 0 {
					clusterRole.Cpu = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().Value())
					clusterRole.Gpu = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.NvidiaGPU().Value())
					clusterRole.Memory = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().Value() / 1024 / 1024 / 1024)
					clusterRole.StorageSize = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.StorageEphemeral().Value() / 1024 / 1024 / 1024)
				}

				clusterCommon := &models.ClusterCommon{
					Role:       clusterRole.Role,
					Hypervisor: "docker",
				}

				clusterRoles[clusterRole.Role] = clusterRole
				clusterCommons[clusterRole.Role] = clusterCommon
			default:
				continue
			}
		}
	}

	err = p.CheckApiVersionsSupported(runtimeId, apiVersions)
	if err != nil {
		return nil, nil, err
	}

	_, ok := clusterRoles[""]
	if !ok {
		clusterRoles[""] = &models.ClusterRole{
			Role: "",
			Env:  string(env),
		}
	}

	return clusterRoles, clusterCommons, nil
}

func (p *Parser) CheckApiVersionsSupported(runtimeId string, apiVersions []string) error {
	if len(apiVersions) == 0 {
		return nil
	}
	kubeconfigGetter := func() (*clientcmdapi.Config, error) {
		runtime, err := runtimeclient.NewRuntime(runtimeId)
		if err != nil {
			return nil, err
		}

		credential := runtime.Credential

		return clientcmd.Load([]byte(credential))
	}

	config, err := clientcmd.BuildConfigFromKubeconfigGetter("", kubeconfigGetter)
	if err != nil {
		return err
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	apiGroupResources, err := discovery.GetAPIGroupResources(client)
	if err != nil {
		return err
	}
	var supportedVersions []string
	for _, group := range apiGroupResources {
		for _, version := range group.Group.Versions {
			supportedVersions = append(supportedVersions, version.GroupVersion)
		}
	}
	p.Logger.Debug("Get runtime [%s] supported versions [%+v]", runtimeId, supportedVersions)
	p.Logger.Debug("Check api versions [%+v]", apiVersions)
	for _, apiVersion := range apiVersions {
		if !stringutil.StringIn(apiVersion, supportedVersions) {
			return gerr.New(gerr.PermissionDenied, gerr.ErrorUnsupportedApiVersion, apiVersion)
		}
	}
	return nil
}

func (p *Parser) Parse(c *chart.Chart, conf []byte, versionId, runtimeId string) (*models.ClusterWrapper, error) {
	customVals, name, description, err := p.ParseCustomValues(conf)
	if err != nil {
		return nil, err
	}

	vals, err := p.parseValues(c, customVals, name)
	if err != nil {
		return nil, err
	}

	cluster, err := p.ParseCluster(name, description, versionId)
	if err != nil {
		return nil, err
	}

	clusterRoles, clusterCommons, err := p.ParseClusterRolesAndClusterCommons(c, vals, customVals, runtimeId)
	if err != nil {
		return nil, err
	}

	clusterWrapper := &models.ClusterWrapper{
		Cluster:        cluster,
		ClusterRoles:   clusterRoles,
		ClusterCommons: clusterCommons,
	}
	return clusterWrapper, nil
}

func (p *Parser) ParseCustomValues(conf []byte) (map[string]interface{}, string, string, error) {
	customVals, err := chartutil.ReadValues(conf)
	if err != nil {
		return nil, "", "", err
	}
	name, ok := p.getStringFromValues(customVals, "Name")
	if !ok {
		return nil, "", "", fmt.Errorf("config [Name] is missing")
	}
	desc, _ := p.getStringFromValues(customVals, "Description")

	delete(customVals, "Name")
	delete(customVals, "Description")

	return customVals, name, desc, nil
}

func (p *Parser) parseValues(c *chart.Chart, customVals map[string]interface{}, name string) (map[string]interface{}, error) {
	// Get and merge values
	chartVals, err := chartutil.ReadValues([]byte(c.Values.GetRaw()))
	if err != nil {
		return nil, err
	}

	mergedVals := p.mergeValues(chartVals, customVals)

	rawMergedVals, err := yaml.Marshal(mergedVals)
	if err != nil {
		return nil, err
	}

	config := &chart.Config{Raw: string(rawMergedVals), Values: map[string]*chart.Value{}}

	// Get release option
	options := chartutil.ReleaseOptions{
		Name:      name,
		Namespace: "",
	}

	vals, err := chartutil.ToRenderValues(c, config, options)
	if err != nil {
		return nil, err
	}

	return vals, nil
}

func (p *Parser) getStringFromValues(vals map[string]interface{}, key string) (string, bool) {
	v, ok := vals[key]
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	if !ok {
		return "", false
	}
	return s, true
}

func (p *Parser) mergeValues(dest map[string]interface{}, src map[string]interface{}) map[string]interface{} {
	for k, v := range src {
		// If the key doesn't exist already, then just set the key to that value
		if _, exists := dest[k]; !exists {
			dest[k] = v
			continue
		}
		nextMap, ok := v.(map[string]interface{})
		// If it isn't another map, overwrite the value
		if !ok {
			dest[k] = v
			continue
		}
		// Edge case: If the key exists in the destination, but isn't a map
		destMap, isMap := dest[k].(map[string]interface{})
		// If the source map has a map for this key, prefer it
		if !isMap {
			dest[k] = v
			continue
		}
		// If we got to this point, it is a map in both, so merge them
		dest[k] = p.mergeValues(destMap, nextMap)
	}
	return dest
}
