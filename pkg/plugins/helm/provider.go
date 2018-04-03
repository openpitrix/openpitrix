// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package helm

import (
	"fmt"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/helm/portforwarder"
	"k8s.io/helm/pkg/kube"
	"k8s.io/helm/pkg/tiller/environment"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/plugins"
)

func init() {
	plugins.RegisterProviderPlugin(constants.ProviderKubernetes, new(Provider))
}

type Provider struct {
}

func (p *Provider) getKubeClient() (clientset *kubernetes.Clientset, config *rest.Config, err error) {
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

func (p *Provider) setupTillerConnection(client kubernetes.Interface, restClientConfig *rest.Config, namespace string) (*kube.Tunnel, error) {
	tunnel, err := portforwarder.New(namespace, client, restClientConfig)
	if err != nil {
		return nil, fmt.Errorf("Could not get a connection to tiller: %v. ", err)
	}

	return tunnel, err
}

func (p *Provider) setupHelm(kubeClient *kubernetes.Clientset, restClientConfig *rest.Config, namespace string) (*helm.Client, error) {
	tunnel, err := p.setupTillerConnection(kubeClient, restClientConfig, namespace)
	if err != nil {
		return nil, err
	}

	return helm.NewClient(helm.Host(fmt.Sprintf("localhost:%d", tunnel.Local))), nil
}

func (p *Provider) getHelmClient() (helmClient *helm.Client, err error) {
	client, clientConfig, err := p.getKubeClient()
	if err != nil {
		return nil, fmt.Errorf("Could not get a kube client: %v. ", err)
	}

	helmClient, err = p.setupHelm(client, clientConfig, environment.DefaultTillerNamespace)
	return
}

func (p *Provider) ParseClusterConf(versionId, conf string) (*models.ClusterWrapper, error) {
	return nil, nil
}

func (p *Provider) SplitJobIntoTasks(job *models.Job) (*models.TaskLayer, error) {
	return nil, nil
}
func (p *Provider) HandleSubtask(task *models.Task) error {
	return nil
}
func (p *Provider) WaitSubtask(task *models.Task, timeout time.Duration, waitInterval time.Duration) error {
	return nil
}
func (p *Provider) DescribeSubnet(subnetId string) (*models.Subnet, error) {
	return nil, nil
}
func (p *Provider) DescribeVpc(vpcId string) (*models.Vpc, error) {
	return nil, nil
}
