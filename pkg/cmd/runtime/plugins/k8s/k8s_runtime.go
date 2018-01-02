// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package k8s_runtime

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"openpitrix.io/openpitrix/pkg/cmd/runtime"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/utils"
)

func init() {
	runtime.RegisterRuntime(new(K8sRuntime))
}

type K8sRuntime struct{}

func (p *K8sRuntime) Name() string { return "k8s" }

func (p *K8sRuntime) Run(app string, args ...string) error {
	return errors.New("TODO")
}

func (p *K8sRuntime) getDefaultClient() (clientset *kubernetes.Clientset, err error) {
	config, err := clientcmd.BuildConfigFromFlags("", homedir.HomeDir()+"/.kube/config")
	if err != nil {
		return
	}
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return
	}
	return
}

func (p *K8sRuntime) CreateCluster(appConf string, shouldWait bool, args ...string) (clusterId string, err error) {
	if len(args) >= 1 {
		clusterId = args[0]
	} else {
		clusterId = p.getClusterId()
	}
	clientset, err := p.getDefaultClient()
	if err != nil {
		logger.Errorf("Failed to get client when create cluster: %v", err)
		return
	}

	decoder := scheme.Codecs.UniversalDeserializer()
	d := yaml.NewYAMLOrJSONDecoder(strings.NewReader(appConf), 4096)
	for {
		ext := apiruntime.RawExtension{}
		if err := d.Decode(&ext); err != nil {
			if err == io.EOF {
				return clusterId, nil
			}
			return clusterId, err
		}
		ext.Raw = bytes.TrimSpace(ext.Raw)
		if len(ext.Raw) == 0 || bytes.Equal(ext.Raw, []byte("null")) {
			continue
		}
		obj, _, err := decoder.Decode(ext.Raw, nil, nil)
		switch o := obj.(type) {
		case *v1beta1.Deployment:
			o.SetName(p.getName(clusterId, o.GetName()))
			deploymentsClient := clientset.ExtensionsV1beta1().Deployments(v1.NamespaceDefault)
			_, err = deploymentsClient.Create(o)
			if err != nil {
				logger.Errorf("Failed to create cluster: %v", err)
				return clusterId, err
			}
		default:
			logger.Errorf("Failed to create cluster with unknow type")
			return clusterId, fmt.Errorf("Failed to create cluster [%s] with unknow type", clusterId)
		}
	}

	return
}

func (p *K8sRuntime) StartClusters(clusterIds string, shouldWait bool, args ...string) error {
	if len(args) < 1 {
		logger.Errorf("Failed to get conf when start cluster")
		return fmt.Errorf("Failed to get conf when start cluster [%s]", clusterIds)
	}

	appConf := args[0]
	clusterId := strings.Split(clusterIds, ",")[0]
	_, err := p.CreateCluster(appConf, shouldWait, clusterId)
	return err
}

func (p *K8sRuntime) RecoverClusters(clusterIds string, shouldWait bool, args ...string) error {
	if len(args) < 1 {
		logger.Errorf("Failed to get conf when recover cluster")
		return fmt.Errorf("Failed to get conf when recover cluster [%s]", clusterIds)
	}

	appConf := args[0]
	clusterId := strings.Split(clusterIds, ",")[0]
	_, err := p.CreateCluster(appConf, shouldWait, clusterId)
	return err
}

func (p *K8sRuntime) DeleteClusters(clusterIds string, shouldWait bool, args ...string) error {
	clientset, err := p.getDefaultClient()
	if err != nil {
		logger.Errorf("Failed to get client when create cluster: %v", err)
		return err
	}
	if len(args) < 1 {
		logger.Errorf("Failed to get conf when create cluster")
		return fmt.Errorf("Failed to get conf when create cluster [%s]", clusterIds)
	}

	appConf := args[0]
	decoder := scheme.Codecs.UniversalDeserializer()
	d := yaml.NewYAMLOrJSONDecoder(strings.NewReader(appConf), 4096)
	for {
		ext := apiruntime.RawExtension{}
		if err := d.Decode(&ext); err != nil {
			if err == io.EOF {
				err = nil
				return err
			}
			return err
		}
		ext.Raw = bytes.TrimSpace(ext.Raw)
		if len(ext.Raw) == 0 || bytes.Equal(ext.Raw, []byte("null")) {
			continue
		}
		obj, _, err := decoder.Decode(ext.Raw, nil, nil)
		switch o := obj.(type) {
		case *v1beta1.Deployment:
			name := p.getName(clusterIds, o.GetName())
			deploymentsClient := clientset.ExtensionsV1beta1().Deployments(v1.NamespaceDefault)
			background := metav1.DeletePropagationBackground
			err = deploymentsClient.Delete(name, &metav1.DeleteOptions{PropagationPolicy: &background})
			if err != nil {
				logger.Errorf("Failed to delete cluster: %v", err)
				return err
			}
		default:
			logger.Errorf("Failed to delete cluster with unknow type")
			return fmt.Errorf("Failed to delete cluster [%s] with unknow type", clusterIds)
		}
	}

	return err
}

func (p *K8sRuntime) StopClusters(clusterIds string, shouldWait bool, args ...string) error {
	if len(args) < 1 {
		logger.Errorf("Failed to get conf when stop cluster")
		return fmt.Errorf("Failed to get conf when stop cluster [%s]", clusterIds)
	}

	appConf := args[0]
	return p.DeleteClusters(clusterIds, shouldWait, appConf)
}

func (p *K8sRuntime) CeaseClusters(clusterIds string, shouldWait bool, args ...string) error {
	if len(args) < 1 {
		logger.Errorf("Failed to get conf when cease cluster")
		return fmt.Errorf("Failed to get conf when cease cluster [%s]", clusterIds)
	}

	appConf := args[0]
	return p.DeleteClusters(clusterIds, shouldWait, appConf)
}

func (p *K8sRuntime) getName(clusterId, oriName string) string {
	return clusterId + "-" + oriName
}

func (p *K8sRuntime) getClusterId() string {
	// todo: need to check db
	uuid := utils.GetLowerAndNumUuid(8)
	return "cl-" + uuid
}

func (p *K8sRuntime) getClusterNodeId() string {
	// todo: need to check db
	uuid := utils.GetLowerAndNumUuid(8)
	return "cln-" + uuid
}
