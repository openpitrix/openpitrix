// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package helm

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/helm/portforwarder"
	"k8s.io/helm/pkg/kube"
	"k8s.io/helm/pkg/tiller/environment"

	"openpitrix.io/openpitrix/pkg/cmd/runtime"
	"openpitrix.io/openpitrix/pkg/utils"
)

func init() {
	runtime.RegisterRuntime(new(HelmRuntime))
}

type HelmRuntime struct{}

func (p *HelmRuntime) Name() string { return "helm" }

func (p *HelmRuntime) Run(app string, args ...string) error {
	return errors.New("TODO")
}

func (p *HelmRuntime) getKubeClient() (clientset *kubernetes.Clientset, config *rest.Config, err error) {
	config, err = clientcmd.BuildConfigFromFlags("", homedir.HomeDir()+"/.kube/config")
	if err != nil {
		return
	}
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return
	}
	return
}

func (p *HelmRuntime) setupTillerConnection(client kubernetes.Interface, restClientConfig *rest.Config, namespace string) (*kube.Tunnel, error) {
	tunnel, err := portforwarder.New(namespace, client, restClientConfig)
	if err != nil {
		return nil, fmt.Errorf("Could not get a connection to tiller: %v", err)
	}

	return tunnel, err
}

func (p *HelmRuntime) setupHelm(kubeClient *kubernetes.Clientset, restClientConfig *rest.Config, namespace string) (*helm.Client, error) {
	tunnel, err := p.setupTillerConnection(kubeClient, restClientConfig, namespace)
	if err != nil {
		return nil, err
	}

	return helm.NewClient(helm.Host(fmt.Sprintf("localhost:%d", tunnel.Local))), nil
}

func (p *HelmRuntime) getHelmClient() (helmClient *helm.Client, err error) {
	client, clientConfig, err := p.getKubeClient()
	if err != nil {
		return nil, fmt.Errorf("Could not get a kube client: %v", err)
	}

	helmClient, err = p.setupHelm(client, clientConfig, environment.DefaultTillerNamespace)
	return
}

func (p *HelmRuntime) CreateCluster(appConf string, shouldWait bool, args ...string) (clusterId string, err error) {
	appConf = strings.Replace(appConf, "~/", os.Getenv("HOME")+"/", 1)

	helmClient, err := p.getHelmClient()
	if err != nil {
		return
	}

	values := ""
	if len(args) >= 1 && args[0] != "" {
		values = args[0]
	}

	ops := []helm.InstallOption{
		helm.ValueOverrides([]byte(values)),
		helm.InstallWait(shouldWait),
		helm.ReleaseName(p.getClusterId()),
	}

	res, err := helmClient.InstallRelease(appConf, v1.NamespaceDefault, ops...)
	if err != nil {
		return
	}

	return res.GetRelease().Name, nil
}

func (p *HelmRuntime) StartClusters(clusterIds string, shouldWait bool, args ...string) error {
	return nil
}

func (p *HelmRuntime) RecoverClusters(clusterIds string, shouldWait bool, args ...string) error {
	return nil
}

func (p *HelmRuntime) DeleteClusters(clusterIds string, shouldWait bool, args ...string) error {
	helmClient, err := p.getHelmClient()
	if err != nil {
		return err
	}

	_, err = helmClient.DeleteRelease(clusterIds)
	return err
}

func (p *HelmRuntime) StopClusters(clusterIds string, shouldWait bool, args ...string) error {
	return nil
}

func (p *HelmRuntime) CeaseClusters(clusterIds string, shouldWait bool, args ...string) error {
	return nil
}

func (p *HelmRuntime) getClusterId() string {
	// todo: need to check db
	uuid := utils.GetLowerAndNumUuid(8)
	return "cl-" + uuid
}
