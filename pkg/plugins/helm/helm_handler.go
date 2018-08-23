// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package helm

import (
	"context"
	"fmt"
	"regexp"

	"google.golang.org/grpc/transport"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/helm/portforwarder"
	"k8s.io/helm/pkg/proto/hapi/chart"
	rls "k8s.io/helm/pkg/proto/hapi/services"
	"k8s.io/helm/pkg/tiller/environment"

	runtimeclient "openpitrix.io/openpitrix/pkg/client/runtime"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/util/funcutil"
)

var (
	ClusterNameRegExp = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`)
)

type HelmHandler struct {
	ctx       context.Context
	RuntimeId string
}

func GetHelmHandler(ctx context.Context, runtimeId string) *HelmHandler {
	helmHandler := new(HelmHandler)
	helmHandler.ctx = ctx
	helmHandler.RuntimeId = runtimeId
	return helmHandler
}

func (p *HelmHandler) initKubeClient() (*kubernetes.Clientset, *rest.Config, error) {
	kubeconfigGetter := func() (*clientcmdapi.Config, error) {
		runtime, err := runtimeclient.NewRuntime(p.ctx, p.RuntimeId)
		if err != nil {
			return nil, err
		}

		credential := runtime.Credential

		return clientcmd.Load([]byte(credential))
	}

	config, err := clientcmd.BuildConfigFromKubeconfigGetter("", kubeconfigGetter)
	if err != nil {
		return nil, nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}
	return clientset, config, err
}

func (p *HelmHandler) initHelmClient() (*helm.Client, error) {
	client, clientConfig, err := p.initKubeClient()
	if err != nil {
		return nil, fmt.Errorf("could not get a kube client: %+v. ", err)
	}

	tunnel, err := portforwarder.New(environment.DefaultTillerNamespace, client, clientConfig)
	if err != nil {
		return nil, fmt.Errorf("could not get a connection to tiller: %+v. ", err)
	}

	hc := helm.NewClient(helm.Host(fmt.Sprintf("localhost:%d", tunnel.Local)))
	return hc, nil
}

func (p *HelmHandler) InstallReleaseFromChart(c *chart.Chart, ns string, rawVals []byte, releaseName string) error {
	hc, err := p.initHelmClient()
	if err != nil {
		return err
	}

	_, err = hc.InstallReleaseFromChart(c, ns, helm.ValueOverrides(rawVals), helm.ReleaseName(releaseName), helm.InstallWait(true))
	if err != nil {
		return err
	}
	return nil
}

func (p *HelmHandler) UpdateReleaseFromChart(releaseName string, c *chart.Chart, rawVals []byte) error {
	hc, err := p.initHelmClient()
	if err != nil {
		return err
	}

	_, err = hc.UpdateReleaseFromChart(releaseName, c, helm.UpdateValueOverrides(rawVals), helm.UpgradeWait(true))
	if err != nil {
		return err
	}
	return nil
}

func (p *HelmHandler) RollbackRelease(releaseName string) error {
	hc, err := p.initHelmClient()
	if err != nil {
		return err
	}

	_, err = hc.RollbackRelease(releaseName, helm.RollbackWait(true))
	if err != nil {
		return err
	}
	return nil
}

func (p *HelmHandler) DeleteRelease(releaseName string, purge bool) error {
	hc, err := p.initHelmClient()
	if err != nil {
		return err
	}

	_, err = hc.DeleteRelease(releaseName, helm.DeletePurge(purge))
	if err != nil {
		return err
	}
	return nil
}

func (p *HelmHandler) ReleaseStatus(releaseName string) (*rls.GetReleaseStatusResponse, error) {
	hc, err := p.initHelmClient()
	if err != nil {
		return nil, err
	}

	return hc.ReleaseStatus(releaseName)
}

func (p *HelmHandler) CheckClusterNameIsUnique(clusterName string) error {
	if clusterName == "" {
		return fmt.Errorf("cluster name must be provided")
	}

	if !ClusterNameRegExp.MatchString(clusterName) {
		return fmt.Errorf(`cluster name must match with regexp "[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*"`)
	}

	err := funcutil.WaitForSpecificOrError(func() (bool, error) {
		_, err := p.ReleaseStatus(clusterName)
		if err != nil {
			if _, ok := err.(transport.ConnectionError); ok {
				return false, nil
			}
			return true, nil
		}

		return true, gerr.New(p.ctx, gerr.PermissionDenied, gerr.ErrorHelmReleaseExists, clusterName)
	}, constants.DefaultServiceTimeout, constants.WaitTaskInterval)
	return err
}
