package helm

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"github.com/ghodss/yaml"
	"k8s.io/api/apps/v1beta1"
	"k8s.io/api/apps/v1beta2"
	exv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/engine"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/kubernetes/pkg/api"
	_ "k8s.io/kubernetes/pkg/api/install"
	_ "k8s.io/kubernetes/pkg/apis/apps/install"
	_ "k8s.io/kubernetes/pkg/apis/extensions/install"

	appclient "openpitrix.io/openpitrix/pkg/client/app"
	runtimeclient "openpitrix.io/openpitrix/pkg/client/runtime"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
)

type Parser struct {
}

func (p *Parser) ParseCluster(vals map[string]interface{}, versionId string) (*models.Cluster, error) {
	clusterId, ok := p.getStringFromValues(vals, "ClusterId")
	if !ok {
		return nil, fmt.Errorf("ClusterId must be set and type is string")
	}

	name, ok := p.getStringFromValues(vals, "Name")
	if !ok {
		name = ""
	}

	desc, ok := p.getStringFromValues(vals, "Description")
	if !ok {
		desc = ""
	}

	ctx := context.Background()
	appManagerClient, err := appclient.NewAppManagerClient(ctx)
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
		return nil, fmt.Errorf("App version [%s] not found", versionId)
	}

	appVersion := resp.AppVersionSet[0]

	runtimeId, ok := p.getStringFromValues(vals, "RuntimeId")
	if !ok {
		return nil, fmt.Errorf("RuntimeId must be set and type is string")
	}

	cluster := &models.Cluster{
		ClusterId:   clusterId,
		Name:        name,
		Description: desc,
		AppId:       appVersion.AppId.GetValue(),
		VersionId:   versionId,
		Owner:       appVersion.Owner.GetValue(),
		RuntimeId:   runtimeId,
		CreateTime:  time.Now(),
		StatusTime:  time.Now(),
	}

	return cluster, nil
}

func (p *Parser) ParseClusterRoles(c *chart.Chart, vals map[string]interface{}) (map[string]*models.ClusterRole, error) {
	clusterId, ok := p.getStringFromValues(vals, "ClusterId")
	if !ok {
		return nil, fmt.Errorf("ClusterId must be set and type is string")
	}

	env, err := json.Marshal(vals)
	if err != nil {
		return nil, err
	}

	renderer := engine.New()
	out, err := renderer.Render(c, vals)
	if err != nil {
		return nil, err
	}

	decode := api.Codecs.UniversalDeserializer().Decode

	clusterRoles := map[string]*models.ClusterRole{}
	for k, v := range out {
		if filepath.Ext(k) != ".yaml" {
			continue
		}

		if len([]byte(v)) == 0 {
			continue
		}

		obj, _, err := decode([]byte(v), nil, nil)
		if err != nil {
			logger.Warnf("Decode file [%s] in chart failed, %+v", k, err)
			continue
		}

		switch o := obj.(type) {
		case *v1beta2.Deployment:
			clusterRole := &models.ClusterRole{
				ClusterId:    clusterId,
				Role:         o.GetObjectMeta().GetName(),
				InstanceSize: uint32(*o.Spec.Replicas),
				Env:          string(env),
			}

			if len(o.Spec.Template.Spec.Containers) > 0 {
				clusterRole.Cpu = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().Value())
				clusterRole.Gpu = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.NvidiaGPU().Value())
				clusterRole.Memory = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().Value())
				clusterRole.StorageSize = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.StorageEphemeral().Value())
			}

			clusterRoles[o.GetObjectMeta().GetName()] = clusterRole
		case *v1beta1.Deployment:
			clusterRole := &models.ClusterRole{
				ClusterId:    clusterId,
				Role:         o.GetObjectMeta().GetName(),
				InstanceSize: uint32(*o.Spec.Replicas),
				Env:          string(env),
			}

			if len(o.Spec.Template.Spec.Containers) > 0 {
				clusterRole.Cpu = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().Value())
				clusterRole.Gpu = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.NvidiaGPU().Value())
				clusterRole.Memory = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().Value())
				clusterRole.StorageSize = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.StorageEphemeral().Value())
			}

			clusterRoles[o.GetObjectMeta().GetName()] = clusterRole
		case *exv1beta1.Deployment:
			clusterRole := &models.ClusterRole{
				ClusterId:    clusterId,
				Role:         o.GetObjectMeta().GetName(),
				InstanceSize: uint32(*o.Spec.Replicas),
				Env:          string(env),
			}

			if len(o.Spec.Template.Spec.Containers) > 0 {
				clusterRole.Cpu = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().Value())
				clusterRole.Gpu = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.NvidiaGPU().Value())
				clusterRole.Memory = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().Value())
				clusterRole.StorageSize = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.StorageEphemeral().Value())
			}

			clusterRoles[o.GetObjectMeta().GetName()] = clusterRole
		case *v1beta2.StatefulSet:
			clusterRole := &models.ClusterRole{
				ClusterId:    clusterId,
				Role:         o.GetObjectMeta().GetName(),
				InstanceSize: uint32(*o.Spec.Replicas),
				Env:          string(env),
			}

			if len(o.Spec.Template.Spec.Containers) > 0 {
				clusterRole.Cpu = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().Value())
				clusterRole.Gpu = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.NvidiaGPU().Value())
				clusterRole.Memory = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().Value())
				clusterRole.StorageSize = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.StorageEphemeral().Value())
			}
		case *v1beta1.StatefulSet:
			clusterRole := &models.ClusterRole{
				ClusterId:    clusterId,
				Role:         o.GetObjectMeta().GetName(),
				InstanceSize: uint32(*o.Spec.Replicas),
				Env:          string(env),
			}

			if len(o.Spec.Template.Spec.Containers) > 0 {
				clusterRole.Cpu = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().Value())
				clusterRole.Gpu = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.NvidiaGPU().Value())
				clusterRole.Memory = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().Value())
				clusterRole.StorageSize = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.StorageEphemeral().Value())
			}
		case *v1beta2.DaemonSet:
			clusterRole := &models.ClusterRole{
				ClusterId:    clusterId,
				Role:         o.GetObjectMeta().GetName(),
				InstanceSize: uint32(1),
				Env:          string(env),
			}

			if len(o.Spec.Template.Spec.Containers) > 0 {
				clusterRole.Cpu = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().Value())
				clusterRole.Gpu = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.NvidiaGPU().Value())
				clusterRole.Memory = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().Value())
				clusterRole.StorageSize = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.StorageEphemeral().Value())
			}
		case *exv1beta1.DaemonSet:
			clusterRole := &models.ClusterRole{
				ClusterId:    clusterId,
				Role:         o.GetObjectMeta().GetName(),
				InstanceSize: uint32(1),
				Env:          string(env),
			}

			if len(o.Spec.Template.Spec.Containers) > 0 {
				clusterRole.Cpu = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().Value())
				clusterRole.Gpu = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.NvidiaGPU().Value())
				clusterRole.Memory = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().Value())
				clusterRole.StorageSize = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.StorageEphemeral().Value())
			}

		default:
			continue
		}
	}

	_, ok = clusterRoles[""]
	if !ok {
		clusterRoles[""] = &models.ClusterRole{
			ClusterId: clusterId,
			Role:      "",
			Env:       string(env),
		}
	}

	return clusterRoles, nil
}

func (p *Parser) Parse(c *chart.Chart, conf []byte, versionId string) (*models.ClusterWrapper, error) {
	vals, err := p.parseValues(c, conf)
	if err != nil {
		return nil, err
	}

	cluster, err := p.ParseCluster(vals, versionId)
	if err != nil {
		return nil, err
	}

	clusterRoles, err := p.ParseClusterRoles(c, vals)
	if err != nil {
		return nil, err
	}

	clusterWrapper := &models.ClusterWrapper{
		Cluster:      cluster,
		ClusterRoles: clusterRoles,
	}
	return clusterWrapper, nil
}

func (p *Parser) parseValues(c *chart.Chart, rawConf []byte) (map[string]interface{}, error) {
	// Get and merge values
	chartVals, err := chartutil.ReadValues([]byte(c.Values.GetRaw()))
	if err != nil {
		return nil, err
	}

	customVals, err := chartutil.ReadValues(rawConf)
	if err != nil {
		return nil, err
	}

	mergedVals := p.mergeValues(chartVals, customVals)

	rawVals, err := yaml.Marshal(mergedVals)
	if err != nil {
		return nil, err
	}
	config := &chart.Config{Raw: string(rawVals), Values: map[string]*chart.Value{}}

	// Get release option
	clusterId, ok := p.getStringFromValues(customVals, "ClusterId")
	if !ok {
		return nil, fmt.Errorf("ClusterId must be set and type is string")
	}

	runtimeId, ok := p.getStringFromValues(customVals, "RuntimeId")
	if !ok {
		return nil, fmt.Errorf("RuntimeId must be set and type is string")
	}

	runtime, err := runtimeclient.NewRuntime(runtimeId)
	if err != nil {
		return nil, err
	}

	namespace := runtime.Zone

	options := chartutil.ReleaseOptions{
		Name:      clusterId,
		Namespace: namespace,
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
