// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package helm

import (
	"bufio"
	"context"
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
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/engine"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	_ "k8s.io/kubernetes/pkg/apis/apps/install"
	_ "k8s.io/kubernetes/pkg/apis/extensions/install"

	clientutil "openpitrix.io/openpitrix/pkg/client"
	appclient "openpitrix.io/openpitrix/pkg/client/app"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
)

type Parser struct {
	ctx       context.Context
	Chart     *chart.Chart
	Conf      string
	VersionId string
	RuntimeId string
}

func (p *Parser) getAppVersion() (*pb.AppVersion, error) {
	ctx := clientutil.SetSystemUserToContext(p.ctx)
	appManagerClient, err := appclient.NewAppManagerClient()
	if err != nil {
		return nil, err
	}

	req := &pb.DescribeAppVersionsRequest{
		VersionId: []string{p.VersionId},
	}

	resp, err := appManagerClient.DescribeAppVersions(ctx, req)
	if err != nil {
		return nil, err
	}

	if len(resp.AppVersionSet) == 0 {
		return nil, fmt.Errorf("app version [%s] not found", p.VersionId)
	}

	appVersion := resp.AppVersionSet[0]
	return appVersion, nil
}

func (p *Parser) parseCluster(name string, description string) (*models.Cluster, error) {
	appVersion, err := p.getAppVersion()
	if err != nil {
		return nil, err
	}

	cluster := &models.Cluster{
		Name:        name,
		Description: description,
		AppId:       appVersion.AppId.GetValue(),
		VersionId:   p.VersionId,
		Status:      constants.StatusPending,
		CreateTime:  time.Now(),
		StatusTime:  time.Now(),
	}

	return cluster, nil
}

func (p *Parser) parseClusterRolesAndClusterCommons(vals map[string]interface{}, customVals map[string]interface{}) (
	map[string]*models.ClusterRole,
	map[string]*models.ClusterCommon,
	error,
) {
	env := jsonutil.ToString(customVals)

	renderer := engine.New()
	out, err := renderer.Render(p.Chart, vals)
	if err != nil {
		return nil, nil, err
	}

	if len(out) == 0 {
		return nil, nil, fmt.Errorf("this chart has no resources defined")
	}

	var apiVersions []string
	decode := legacyscheme.Codecs.UniversalDeserializer().Decode

	clusterRoles := map[string]*models.ClusterRole{}
	clusterCommons := map[string]*models.ClusterCommon{}
	for k, v := range out {
		if filepath.Ext(k) != ".yaml" {
			continue
		}

		if len(strings.TrimSpace(v)) == 0 {
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
				logger.Error(p.ctx, "Decode file [%s] in chart failed, %+v", k, err)
				return nil, nil, err
			}
			obj, groupVersionKind, err := decode(doc, nil, nil)

			if err != nil {
				logger.Error(p.ctx, "Decode file [%s] in chart failed, %+v", k, err)
				return nil, nil, err
			}
			logger.Debug(p.ctx, "Yaml content: %+v", obj)
			logger.Debug(p.ctx, "Group version: %+v", groupVersionKind.GroupVersion().String())

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

	kubeHandler := GetKubeHandler(p.ctx, p.RuntimeId)
	err = kubeHandler.CheckApiVersionsSupported(apiVersions)
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

func (p *Parser) Parse() (*models.ClusterWrapper, error) {
	customVals, name, description, err := p.parseCustomValues()
	if err != nil {
		return nil, err
	}

	vals, err := p.parseValues(customVals, name)
	if err != nil {
		return nil, err
	}

	cluster, err := p.parseCluster(name, description)
	if err != nil {
		return nil, err
	}

	clusterRoles, clusterCommons, err := p.parseClusterRolesAndClusterCommons(vals, customVals)
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

func (p *Parser) parseCustomValues() (map[string]interface{}, string, string, error) {
	customVals, err := chartutil.ReadValues([]byte(p.Conf))
	if err != nil {
		return nil, "", "", err
	}
	name, ok := GetStringFromValues(customVals, "Name")
	if !ok {
		return nil, "", "", fmt.Errorf("config [Name] is missing")
	}
	desc, _ := GetStringFromValues(customVals, "Description")

	return customVals, name, desc, nil
}

func (p *Parser) parseValues(customVals map[string]interface{}, name string) (map[string]interface{}, error) {
	// Get and merge values
	chartVals, err := chartutil.ReadValues([]byte(p.Chart.Values.GetRaw()))
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

	vals, err := chartutil.ToRenderValues(p.Chart, config, options)
	if err != nil {
		return nil, err
	}

	return vals, nil
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
