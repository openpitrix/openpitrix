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

	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/engine"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/releaseutil"
	appsv1 "k8s.io/api/apps/v1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	corev1 "k8s.io/api/core/v1"
	exv1beta1 "k8s.io/api/extensions/v1beta1"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes/scheme"

	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
)

type Render struct {
	Chart        *chart.Chart
	CustomValues string
	RuntimeId    string
	Namespace    string
}

func PrepareValues(c *chart.Chart, runtimeId, namespace, customValues string) (map[string]interface{}, error) {
	customVals, err := chartutil.ReadValues([]byte(customValues))
	if err != nil {
		return nil, err
	}
	name, ok := GetStringFromValues(customVals, "Name")
	if !ok {
		return nil, fmt.Errorf("config [Name] is missing or empty")
	}

	mergedVals := mergeValues(c.Values, customVals)

	// Get release option
	options := chartutil.ReleaseOptions{
		Name:      name,
		Namespace: namespace,
	}

	proxy := NewProxy(context.Background(), runtimeId)
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

	vals, err := chartutil.ToRenderValues(c, mergedVals, options, caps)
	if err != nil {
		return nil, err
	}

	return vals, nil
}

func RenderChartWithValues(c *chart.Chart, vals map[string]interface{}) (map[string]string, error) {
	files, err := engine.Render(c, vals)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("this chart has no resources defined")
	}

	return files, nil
}

func exactApis(runtimeId string, files map[string]string) ([]string, error) {
	var apiVersions []string
	decode := scheme.Codecs.UniversalDeserializer().Decode

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
					logger.Error(nil, "Decode file [%s] in chart failed, %+v", filePath, err)
					return nil, err
				}
				obj, groupVersionKind, err := decode(doc, nil, nil)

				if err != nil {
					logger.Error(nil, "Decode file [%s] in chart failed, %+v", filePath, err)
					return nil, err
				}
				logger.Debug(nil, "Yaml content: %+v", obj)
				logger.Debug(nil, "Group version: %+v", groupVersionKind.GroupVersion().String())

				apiVersions = append(apiVersions, groupVersionKind.GroupVersion().String())
			}
		}
	}
	proxy := NewProxy(context.Background(), runtimeId)
	err := proxy.CheckApiVersionsSupported(apiVersions)
	if err != nil {
		return nil, err
	}
	return apiVersions, nil
}

func ExactApis(c *chart.Chart, namespace, runtimeId, customValues string) ([]string, error) {
	vals, err := PrepareValues(c, runtimeId, namespace, customValues)
	if err != nil {
		return nil, err
	}

	files, err := RenderChartWithValues(c, vals)
	if err != nil {
		return nil, fmt.Errorf("render with values error")
	}
	return exactApis(runtimeId, files)
}

func ExactResources(rls *release.Release) (map[string]string, string, error) {
	additionalInfo := map[string][]map[string]interface{}{
		"service":   {},
		"configmap": {},
		"secret":    {},
		"pvc":       {},
		"ingress":   {},
	}
	decode := scheme.Codecs.UniversalDeserializer().Decode

	roles := map[string]string{}
	manifests := releaseutil.SplitManifests(rls.Manifest)

	for _, manifest := range manifests {
		b := bufio.NewReader(strings.NewReader(manifest))
		r := k8syaml.NewYAMLReader(b)
		doc, err := r.Read()

		if err != nil {
			return nil, "", err
		}
		obj, groupVersionKind, err := decode(doc, nil, nil)

		switch o := obj.(type) {
		case *appsv1.Deployment:
			roles["deployment"] = o.GetObjectMeta().GetName()
		case *appsv1beta2.Deployment:
			roles["deployment"] = o.GetObjectMeta().GetName()
		case *appsv1beta1.Deployment:
			roles["deployment"] = o.GetObjectMeta().GetName()
		case *exv1beta1.Deployment:
			roles["deployment"] = o.GetObjectMeta().GetName()
		case *appsv1.StatefulSet:
			roles["statefulset"] = o.GetObjectMeta().GetName()
		case *appsv1beta2.StatefulSet:
			roles["statefulset"] = o.GetObjectMeta().GetName()
		case *appsv1beta1.StatefulSet:
			roles["statefulset"] = o.GetObjectMeta().GetName()
		case *appsv1.DaemonSet:
			roles["daemonset"] = o.GetObjectMeta().GetName()
		case *appsv1beta2.DaemonSet:
			roles["daemonset"] = o.GetObjectMeta().GetName()
		case *exv1beta1.DaemonSet:
			roles["daemonset"] = o.GetObjectMeta().GetName()
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

	return roles, jsonutil.ToString(additionalInfo), nil
}

func (p *Render) parseCustomValues() (map[string]interface{}, string, string, error) {
	customVals, err := chartutil.ReadValues([]byte(p.CustomValues))
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

func (p *Render) parseValues(customVals map[string]interface{}, name string) (map[string]interface{}, error) {
	mergedVals := p.mergeValues(p.Chart.Values, customVals)

	// Get release option
	options := chartutil.ReleaseOptions{
		Name:      name,
		Namespace: p.Namespace,
	}

	proxy := NewProxy(context.Background(), p.RuntimeId)
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

func (p *Render) mergeValues(dest map[string]interface{}, src map[string]interface{}) map[string]interface{} {
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
