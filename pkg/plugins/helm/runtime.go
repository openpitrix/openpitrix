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
	plugins.RegisterRuntimePlugin(constants.RuntimeKubernetes, new(Runtime))
}

type Runtime struct {
}

func (p *Runtime) getKubeClient() (clientset *kubernetes.Clientset, config *rest.Config, err error) {
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

func (p *Runtime) setupTillerConnection(client kubernetes.Interface, restClientConfig *rest.Config, namespace string) (*kube.Tunnel, error) {
	tunnel, err := portforwarder.New(namespace, client, restClientConfig)
	if err != nil {
		return nil, fmt.Errorf("Could not get a connection to tiller: %v. ", err)
	}

	return tunnel, err
}

func (p *Runtime) setupHelm(kubeClient *kubernetes.Clientset, restClientConfig *rest.Config, namespace string) (*helm.Client, error) {
	tunnel, err := p.setupTillerConnection(kubeClient, restClientConfig, namespace)
	if err != nil {
		return nil, err
	}

	return helm.NewClient(helm.Host(fmt.Sprintf("localhost:%d", tunnel.Local))), nil
}

func (p *Runtime) getHelmClient() (helmClient *helm.Client, err error) {
	client, clientConfig, err := p.getKubeClient()
	if err != nil {
		return nil, fmt.Errorf("Could not get a kube client: %v. ", err)
	}

	helmClient, err = p.setupHelm(client, clientConfig, environment.DefaultTillerNamespace)
	return
}

func (p *Runtime) ParseClusterConf(versionId, conf string) (*models.ClusterWrapper, error) {
	return nil, nil
}

func (p *Runtime) SplitJobIntoTasks(job *models.Job) (*models.TaskLayer, error) {
	return nil, nil
}
func (p *Runtime) HandleSubtask(task *models.Task) error {
	return nil
}
func (p *Runtime) WaitSubtask(taskId string, timeout time.Duration, waitInterval time.Duration) error {
	return nil
}
