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

	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/engine"
	"helm.sh/helm/v3/pkg/releaseutil"
	appsv1 "k8s.io/api/apps/v1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	corev1 "k8s.io/api/core/v1"
	exv1beta1 "k8s.io/api/extensions/v1beta1"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes/scheme"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
)

type Parser struct {
	ctx       context.Context
	Chart     *chart.Chart
	Conf      string
	VersionId string
	RuntimeId string
	Namespace string
}

func (p *Parser) parseCluster(name string, description string, additionalInfo string, customVals map[string]interface{}, appId string) (*models.Cluster, error) {
	env := jsonutil.ToString(customVals)

	cluster := &models.Cluster{
		Zone:           p.Namespace,
		Name:           name,
		Description:    description,
		AppId:          appId,
		VersionId:      p.VersionId,
		Status:         constants.StatusPending,
		RuntimeId:      p.RuntimeId,
		CreateTime:     time.Now(),
		StatusTime:     time.Now(),
		AdditionalInfo: additionalInfo,
		Env:            env,
	}

	return cluster, nil
}

func (p *Parser) parseClusterRolesAndClusterCommons(vals map[string]interface{}) (
	map[string]*models.ClusterRole,
	map[string]*models.ClusterCommon,
	string,
	error,
) {
	additionalInfo := map[string][]map[string]interface{}{
		"service":   {},
		"configmap": {},
		"secret":    {},
		"pvc":       {},
		"ingress":   {},
	}

	files, err := engine.Render(p.Chart, vals)
	if err != nil {
		return nil, nil, "", err
	}

	if len(files) == 0 {
		return nil, nil, "", fmt.Errorf("this chart has no resources defined")
	}

	var apiVersions []string
	decode := scheme.Codecs.UniversalDeserializer().Decode

	clusterRoles := map[string]*models.ClusterRole{}
	clusterCommons := map[string]*models.ClusterCommon{}
	for filePath, content := range files {
		if filepath.Ext(filePath) != ".yaml" && filepath.Ext(filePath) != ".yml" {
			continue
		}

		if len(strings.TrimSpace(content)) == 0 {
			continue
		}

		manifests := releaseutil.SplitManifests(content)

		for _, manifest := range manifests {
			b := bufio.NewReader(strings.NewReader(manifest))
			r := k8syaml.NewYAMLReader(b)
			for {
				doc, err := r.Read()
				if err == io.EOF {
					break
				}
				if err != nil {
					logger.Error(p.ctx, "Decode file [%s] in chart failed, %+v", filePath, err)
					return nil, nil, "", err
				}
				obj, groupVersionKind, err := decode(doc, nil, nil)

				if err != nil {
					logger.Error(p.ctx, "Decode file [%s] in chart failed, %+v", filePath, err)
					//if unregister do not deal error
					continue
				}
				logger.Debug(p.ctx, "Yaml content: %+v", obj)
				logger.Debug(p.ctx, "Group version: %+v", groupVersionKind.GroupVersion().String())

				apiVersions = append(apiVersions, groupVersionKind.GroupVersion().String())

				switch o := obj.(type) {
				case *appsv1.Deployment:
					clusterRole := &models.ClusterRole{
						Role:       fmt.Sprintf("%s-Deployment", o.GetObjectMeta().GetName()),
						ApiVersion: groupVersionKind.GroupVersion().String(),
					}

					if o.Spec.Replicas == nil {
						clusterRole.Replicas = 1
					} else {
						clusterRole.Replicas = uint32(*o.Spec.Replicas)
					}

					if len(o.Spec.Template.Spec.Containers) > 0 {
						clusterRole.Cpu = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().Value())
						clusterRole.Memory = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().Value() / 1024 / 1024 / 1024)
						clusterRole.StorageSize = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.StorageEphemeral().Value() / 1024 / 1024 / 1024)
					}

					clusterCommon := &models.ClusterCommon{
						Role:       clusterRole.Role,
						Hypervisor: "docker",
					}

					clusterRoles[clusterRole.Role] = clusterRole
					clusterCommons[clusterRole.Role] = clusterCommon
				case *appsv1beta2.Deployment:
					clusterRole := &models.ClusterRole{
						Role:       fmt.Sprintf("%s-Deployment", o.GetObjectMeta().GetName()),
						ApiVersion: groupVersionKind.GroupVersion().String(),
					}

					if o.Spec.Replicas == nil {
						clusterRole.Replicas = 1
					} else {
						clusterRole.Replicas = uint32(*o.Spec.Replicas)
					}

					if len(o.Spec.Template.Spec.Containers) > 0 {
						clusterRole.Cpu = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().Value())
						clusterRole.Memory = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().Value() / 1024 / 1024 / 1024)
						clusterRole.StorageSize = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.StorageEphemeral().Value() / 1024 / 1024 / 1024)
					}

					clusterCommon := &models.ClusterCommon{
						Role:       clusterRole.Role,
						Hypervisor: "docker",
					}

					clusterRoles[clusterRole.Role] = clusterRole
					clusterCommons[clusterRole.Role] = clusterCommon
				case *appsv1beta1.Deployment:
					clusterRole := &models.ClusterRole{
						Role:       fmt.Sprintf("%s-Deployment", o.GetObjectMeta().GetName()),
						ApiVersion: groupVersionKind.GroupVersion().String(),
					}

					if o.Spec.Replicas == nil {
						clusterRole.Replicas = 1
					} else {
						clusterRole.Replicas = uint32(*o.Spec.Replicas)
					}

					if len(o.Spec.Template.Spec.Containers) > 0 {
						clusterRole.Cpu = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().Value())
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
						Role:       fmt.Sprintf("%s-Deployment", o.GetObjectMeta().GetName()),
						ApiVersion: groupVersionKind.GroupVersion().String(),
					}

					if o.Spec.Replicas == nil {
						clusterRole.Replicas = 1
					} else {
						clusterRole.Replicas = uint32(*o.Spec.Replicas)
					}

					if len(o.Spec.Template.Spec.Containers) > 0 {
						clusterRole.Cpu = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().Value())
						clusterRole.Memory = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().Value() / 1024 / 1024 / 1024)
						clusterRole.StorageSize = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.StorageEphemeral().Value() / 1024 / 1024 / 1024)
					}

					clusterCommon := &models.ClusterCommon{
						Role:       clusterRole.Role,
						Hypervisor: "docker",
					}

					clusterRoles[clusterRole.Role] = clusterRole
					clusterCommons[clusterRole.Role] = clusterCommon
				case *appsv1.StatefulSet:
					clusterRole := &models.ClusterRole{
						Role:       fmt.Sprintf("%s-StatefulSet", o.GetObjectMeta().GetName()),
						ApiVersion: groupVersionKind.GroupVersion().String(),
					}

					if o.Spec.Replicas == nil {
						clusterRole.Replicas = 1
					} else {
						clusterRole.Replicas = uint32(*o.Spec.Replicas)
					}

					if len(o.Spec.Template.Spec.Containers) > 0 {
						clusterRole.Cpu = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().Value())
						clusterRole.Memory = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().Value() / 1024 / 1024 / 1024)
						clusterRole.StorageSize = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.StorageEphemeral().Value() / 1024 / 1024 / 1024)
					}

					clusterCommon := &models.ClusterCommon{
						Role:       clusterRole.Role,
						Hypervisor: "docker",
					}

					clusterRoles[clusterRole.Role] = clusterRole
					clusterCommons[clusterRole.Role] = clusterCommon
				case *appsv1beta2.StatefulSet:
					clusterRole := &models.ClusterRole{
						Role:       fmt.Sprintf("%s-StatefulSet", o.GetObjectMeta().GetName()),
						ApiVersion: groupVersionKind.GroupVersion().String(),
					}

					if o.Spec.Replicas == nil {
						clusterRole.Replicas = 1
					} else {
						clusterRole.Replicas = uint32(*o.Spec.Replicas)
					}

					if len(o.Spec.Template.Spec.Containers) > 0 {
						clusterRole.Cpu = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().Value())
						clusterRole.Memory = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().Value() / 1024 / 1024 / 1024)
						clusterRole.StorageSize = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.StorageEphemeral().Value() / 1024 / 1024 / 1024)
					}

					clusterCommon := &models.ClusterCommon{
						Role:       clusterRole.Role,
						Hypervisor: "docker",
					}

					clusterRoles[clusterRole.Role] = clusterRole
					clusterCommons[clusterRole.Role] = clusterCommon
				case *appsv1beta1.StatefulSet:
					clusterRole := &models.ClusterRole{
						Role:       fmt.Sprintf("%s-StatefulSet", o.GetObjectMeta().GetName()),
						ApiVersion: groupVersionKind.GroupVersion().String(),
					}

					if o.Spec.Replicas == nil {
						clusterRole.Replicas = 1
					} else {
						clusterRole.Replicas = uint32(*o.Spec.Replicas)
					}

					if len(o.Spec.Template.Spec.Containers) > 0 {
						clusterRole.Cpu = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().Value())
						clusterRole.Memory = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().Value() / 1024 / 1024 / 1024)
						clusterRole.StorageSize = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.StorageEphemeral().Value() / 1024 / 1024 / 1024)
					}

					clusterCommon := &models.ClusterCommon{
						Role:       clusterRole.Role,
						Hypervisor: "docker",
					}

					clusterRoles[clusterRole.Role] = clusterRole
					clusterCommons[clusterRole.Role] = clusterCommon
				case *appsv1.DaemonSet:
					clusterRole := &models.ClusterRole{
						Role:       fmt.Sprintf("%s-DaemonSet", o.GetObjectMeta().GetName()),
						ApiVersion: groupVersionKind.GroupVersion().String(),
						Replicas:   uint32(1),
					}

					if len(o.Spec.Template.Spec.Containers) > 0 {
						clusterRole.Cpu = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().Value())
						clusterRole.Memory = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().Value() / 1024 / 1024 / 1024)
						clusterRole.StorageSize = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.StorageEphemeral().Value() / 1024 / 1024 / 1024)
					}

					clusterCommon := &models.ClusterCommon{
						Role:       clusterRole.Role,
						Hypervisor: "docker",
					}

					clusterRoles[clusterRole.Role] = clusterRole
					clusterCommons[clusterRole.Role] = clusterCommon
				case *appsv1beta2.DaemonSet:
					clusterRole := &models.ClusterRole{
						Role:       fmt.Sprintf("%s-DaemonSet", o.GetObjectMeta().GetName()),
						ApiVersion: groupVersionKind.GroupVersion().String(),
						Replicas:   uint32(1),
					}

					if len(o.Spec.Template.Spec.Containers) > 0 {
						clusterRole.Cpu = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().Value())
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
						Role:       fmt.Sprintf("%s-DaemonSet", o.GetObjectMeta().GetName()),
						ApiVersion: groupVersionKind.GroupVersion().String(),
						Replicas:   uint32(1),
					}

					if len(o.Spec.Template.Spec.Containers) > 0 {
						clusterRole.Cpu = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().Value())
						clusterRole.Memory = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().Value() / 1024 / 1024 / 1024)
						clusterRole.StorageSize = uint32(o.Spec.Template.Spec.Containers[0].Resources.Requests.StorageEphemeral().Value() / 1024 / 1024 / 1024)
					}

					clusterCommon := &models.ClusterCommon{
						Role:       clusterRole.Role,
						Hypervisor: "docker",
					}

					clusterRoles[clusterRole.Role] = clusterRole
					clusterCommons[clusterRole.Role] = clusterCommon
				case *corev1.Service:
					additionalInfo["service"] = append(additionalInfo["service"], map[string]interface{}{
						"apiVersion": groupVersionKind.GroupVersion().String(),
						"name":       o.GetObjectMeta().GetName(),
					})
				case *corev1.ConfigMap:
					additionalInfo["configmap"] = append(additionalInfo["configmap"], map[string]interface{}{
						"apiVersion": groupVersionKind.GroupVersion().String(),
						"name":       o.GetObjectMeta().GetName(),
					})
				case *corev1.Secret:
					additionalInfo["secret"] = append(additionalInfo["secret"], map[string]interface{}{
						"apiVersion": groupVersionKind.GroupVersion().String(),
						"name":       o.GetObjectMeta().GetName(),
					})
				case *corev1.PersistentVolumeClaim:
					additionalInfo["pvc"] = append(additionalInfo["pvc"], map[string]interface{}{
						"apiVersion": groupVersionKind.GroupVersion().String(),
						"name":       o.GetObjectMeta().GetName(),
					})
				case *exv1beta1.Ingress:
					additionalInfo["ingress"] = append(additionalInfo["ingress"], map[string]interface{}{
						"apiVersion": groupVersionKind.GroupVersion().String(),
						"name":       o.GetObjectMeta().GetName(),
					})
				default:
					continue
				}
			}
		}
	}

	proxy := NewProxy(p.ctx, p.RuntimeId)
	err = proxy.CheckApiVersionsSupported(apiVersions)
	if err != nil {
		return nil, nil, "", err
	}

	return clusterRoles, clusterCommons, jsonutil.ToString(additionalInfo), nil
}

func (p *Parser) Parse(clusterWrapper *models.ClusterWrapper, appId string) error {
	customVals, name, description, err := p.parseCustomValues()
	if err != nil {
		return err
	}

	if clusterWrapper.Cluster != nil && clusterWrapper.Cluster.Name != "" {
		if name != clusterWrapper.Cluster.Name {
			return fmt.Errorf("cluster name can't be changed by config [Name]")
		}
	}

	vals, err := p.parseValues(customVals, name)
	if err != nil {
		return err
	}

	clusterRoles, clusterCommons, additionalInfo, err := p.parseClusterRolesAndClusterCommons(vals)
	if err != nil {
		return err
	}

	cluster, err := p.parseCluster(name, description, additionalInfo, customVals, appId)
	if err != nil {
		return err
	}

	clusterWrapper.Cluster = cluster
	clusterWrapper.ClusterRoles = clusterRoles
	clusterWrapper.ClusterCommons = clusterCommons
	return nil
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

	if name == "" {
		return nil, "", "", fmt.Errorf("config [Name] is empty")
	}

	desc, _ := GetStringFromValues(customVals, "Description")

	return customVals, name, desc, nil
}

func (p *Parser) parseValues(customVals map[string]interface{}, name string) (map[string]interface{}, error) {
	mergedVals := p.mergeValues(p.Chart.Values, customVals)

	// Get release option
	options := chartutil.ReleaseOptions{
		Name:      name,
		Namespace: p.Namespace,
	}

	proxy := NewProxy(p.ctx, p.RuntimeId)
	version, err := proxy.DescribeVersionInfo()
	if err != nil {
		return nil, err
	}

	kubeversion := chartutil.KubeVersion{
		Version: "v" + version.Minor + "." + version.Minor,
		Major:   version.Major,
		Minor:   version.Minor,
	}
	caps := &chartutil.Capabilities{APIVersions: chartutil.DefaultVersionSet, KubeVersion: kubeversion}

	vals, err := chartutil.ToRenderValues(p.Chart, mergedVals, options, caps)
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
