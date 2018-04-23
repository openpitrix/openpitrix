// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package helm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ghodss/yaml"
	"google.golang.org/grpc/transport"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/helm/portforwarder"
	"k8s.io/helm/pkg/kube"
	"k8s.io/helm/pkg/proto/hapi/release"
	"k8s.io/helm/pkg/tiller/environment"

	appclient "openpitrix.io/openpitrix/pkg/client/app"
	clusterclient "openpitrix.io/openpitrix/pkg/client/cluster"
	runtimeclient "openpitrix.io/openpitrix/pkg/client/runtime"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/utils"
)

type Provider struct {
}

func (p *Provider) getKubeClient(runtimeId string) (clientset *kubernetes.Clientset, config *rest.Config, err error) {
	kubeconfigGetter := func() (*clientcmdapi.Config, error) {
		runtime, err := runtimeclient.NewRuntime(runtimeId)
		if err != nil {
			return nil, err
		}

		credential := runtime.Credential

		return clientcmd.Load([]byte(credential))
	}

	config, err = clientcmd.BuildConfigFromKubeconfigGetter("", kubeconfigGetter)
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
		return nil, fmt.Errorf("Could not get a connection to tiller: %+v. ", err)
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

func (p *Provider) getHelmClient(runtimeId string) (helmClient *helm.Client, err error) {
	client, clientConfig, err := p.getKubeClient(runtimeId)
	if err != nil {
		return nil, fmt.Errorf("Could not get a kube client: %+v. ", err)
	}

	helmClient, err = p.setupHelm(client, clientConfig, environment.DefaultTillerNamespace)
	return
}

func (p *Provider) ParseClusterConf(versionId, conf string) (*models.ClusterWrapper, error) {
	ctx := context.Background()
	appManagerClient, err := appclient.NewAppManagerClient(ctx)
	if err != nil {
		logger.Errorf("Connect to app manager failed: %+v", err)
		return nil, err
	}

	req := &pb.GetAppVersionPackageRequest{
		VersionId: utils.ToProtoString(versionId),
	}

	resp, err := appManagerClient.GetAppVersionPackage(ctx, req)
	if err != nil {
		logger.Errorf("Get app version [%s] package failed: %+v", versionId, err)
		return nil, err
	}

	pkg := resp.GetPackage()
	r := bytes.NewReader(pkg)

	c, err := chartutil.LoadArchive(r)
	if err != nil {
		return nil, err
	}

	parser := Parser{}
	clusterWrapper, err := parser.Parse(c, []byte(conf), versionId)
	if err != nil {
		return nil, err
	}
	return clusterWrapper, nil
}

func (p *Provider) SplitJobIntoTasks(job *models.Job) (*models.TaskLayer, error) {
	jodDirective, err := getJobDirective(job.Directive)
	if err != nil {
		return nil, err
	}

	switch job.JobAction {
	case constants.ActionCreateCluster:
		tasks := make([]*models.Task, 1)

		td := TaskDirective{
			VersionId: job.VersionId,
			Namespace: jodDirective.Namespace,
			ClusterId: job.ClusterId,
			RuntimeId: jodDirective.RuntimeId,
			Values:    jodDirective.Values,
		}
		tdj := getTaskDirectiveJson(td)

		task := models.NewTask("", job.JobId, "", constants.ProviderKubernetes, constants.ActionCreateCluster, tdj, job.Owner)
		tasks = append(tasks, task)
		tl := models.TaskLayer{
			Tasks: tasks,
			Child: nil,
		}

		return &tl, nil
	case constants.ActionUpgradeCluster:
		tasks := make([]*models.Task, 1)

		td := TaskDirective{
			VersionId: job.VersionId,
			ClusterId: job.ClusterId,
			RuntimeId: jodDirective.RuntimeId,
			Values:    jodDirective.Values,
		}
		tdj := getTaskDirectiveJson(td)

		task := models.NewTask("", job.JobId, "", constants.ProviderKubernetes, constants.ActionUpgradeCluster, tdj, job.Owner)
		tasks = append(tasks, task)
		tl := models.TaskLayer{
			Tasks: tasks,
			Child: nil,
		}

		return &tl, nil
	case constants.ActionRollbackCluster:
		tasks := make([]*models.Task, 1)

		td := TaskDirective{ClusterId: job.ClusterId, RuntimeId: jodDirective.RuntimeId}
		tdj := getTaskDirectiveJson(td)

		task := models.NewTask("", job.JobId, "", constants.ProviderKubernetes, constants.ActionRollbackCluster, tdj, job.Owner)
		tasks = append(tasks, task)
		tl := models.TaskLayer{
			Tasks: tasks,
			Child: nil,
		}

		return &tl, nil
	case constants.ActionDeleteClusters:
		tasks := make([]*models.Task, 1)

		td := TaskDirective{ClusterId: job.ClusterId, RuntimeId: jodDirective.RuntimeId}
		tdj := getTaskDirectiveJson(td)

		task := models.NewTask("", job.JobId, "", constants.ProviderKubernetes, constants.ActionDeleteClusters, tdj, job.Owner)
		tasks = append(tasks, task)
		tl := models.TaskLayer{
			Tasks: tasks,
			Child: nil,
		}

		return &tl, nil
	case constants.ActionCeaseClusters:
		tasks := make([]*models.Task, 1)

		td := TaskDirective{ClusterId: job.ClusterId, RuntimeId: jodDirective.RuntimeId}
		tdj := getTaskDirectiveJson(td)

		task := models.NewTask("", job.JobId, "", constants.ProviderKubernetes, constants.ActionCeaseClusters, tdj, job.Owner)
		tasks = append(tasks, task)
		tl := models.TaskLayer{
			Tasks: tasks,
			Child: nil,
		}

		return &tl, nil
	default:
		return nil, fmt.Errorf("the job action [%s] is not supported", job.JobAction)
	}
	return nil, nil
}

func (p *Provider) HandleSubtask(task *models.Task) error {
	taskDirective, err := getTaskDirective(task.Directive)
	if err != nil {
		return err
	}

	hc, err := p.getHelmClient(taskDirective.RuntimeId)
	if err != nil {
		return err
	}

	ctx := context.Background()

	switch task.TaskAction {
	case constants.ActionCreateCluster:
		appc, err := appclient.NewAppManagerClient(ctx)
		if err != nil {
			return err
		}

		req := pb.GetAppVersionPackageRequest{
			VersionId: utils.ToProtoString(taskDirective.VersionId),
		}

		resp, err := appc.GetAppVersionPackage(ctx, &req)
		if err != nil {
			return err
		}

		pkg := resp.GetPackage()
		r := bytes.NewReader(pkg)

		c, err := chartutil.LoadArchive(r)
		if err != nil {
			return err
		}

		var v map[string]interface{}
		err = json.Unmarshal([]byte(taskDirective.Values), &v)
		if err != nil {
			return err
		}

		rawVals, err := yaml.Marshal(v)
		if err != nil {
			return err
		}

		_, err = hc.InstallReleaseFromChart(c, taskDirective.Namespace,
			helm.ValueOverrides(rawVals),
			helm.ReleaseName(taskDirective.ClusterId),
			helm.InstallWait(false))
		if err != nil {
			return err
		}
	case constants.ActionUpgradeCluster:
		appc, err := appclient.NewAppManagerClient(ctx)
		if err != nil {
			return err
		}

		req := pb.GetAppVersionPackageRequest{
			VersionId: utils.ToProtoString(taskDirective.VersionId),
		}

		resp, err := appc.GetAppVersionPackage(ctx, &req)
		if err != nil {
			return err
		}

		p := resp.GetPackage()
		r := bytes.NewReader(p)

		chart, err := chartutil.LoadArchive(r)
		if err != nil {
			return err
		}

		_, err = hc.UpdateReleaseFromChart(taskDirective.ClusterId, chart, helm.UpgradeWait(false))
		if err != nil {
			return err
		}
	case constants.ActionRollbackCluster:
		_, err = hc.RollbackRelease(taskDirective.ClusterId, helm.RollbackWait(false))
		if err != nil {
			return err
		}
	case constants.ActionDeleteClusters:
		_, err = hc.DeleteRelease(taskDirective.ClusterId)
		if err != nil {
			return err
		}
	case constants.ActionCeaseClusters:
		_, err = hc.DeleteRelease(taskDirective.ClusterId, helm.DeletePurge(true))
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("the task action [%s] is not supported", task.TaskAction)
	}

	return nil
}

func (p *Provider) WaitSubtask(task *models.Task, timeout time.Duration, waitInterval time.Duration) error {
	taskDirective, err := getTaskDirective(task.Directive)
	if err != nil {
		return err
	}

	hc, err := p.getHelmClient(taskDirective.RuntimeId)
	if err != nil {
		return err
	}

	utils.WaitForSpecificOrError(func() (bool, error) {
		switch task.TaskAction {
		case constants.ActionCreateCluster:
			fallthrough
		case constants.ActionUpgradeCluster:
			fallthrough
		case constants.ActionRollbackCluster:
			resp, err := hc.ReleaseStatus(taskDirective.ClusterId)
			if err != nil {
				//network or api error, not considered task fail.
				return false, nil
			}

			if resp.Info.Status.Code == release.Status_DEPLOYED {
				return true, nil
			}
		case constants.ActionDeleteClusters:
			resp, err := hc.ReleaseStatus(taskDirective.ClusterId)
			if err != nil {
				//network or api error, not considered task fail.
				return false, nil
			}

			if resp.Info.Status.Code == release.Status_DELETED {
				return true, nil
			}
		case constants.ActionCeaseClusters:
			_, err := hc.ReleaseStatus(taskDirective.ClusterId)
			if err != nil {
				if _, ok := err.(transport.ConnectionError); ok {
					return false, nil
				}
				return true, nil
			}
		}
		return false, nil
	}, timeout, waitInterval)

	return nil
}

func (p *Provider) DescribeSubnet(runtimeId, subnetId string) (*models.Subnet, error) {
	return nil, nil
}
func (p *Provider) DescribeVpc(runtimeId, vpcId string) (*models.Vpc, error) {
	return nil, nil
}

func (p *Provider) getKubePodsAsClusterNodes(runtimeId, namespace, clusterId, owner string) ([]*pb.ClusterNode, error) {
	kubeClient, _, err := p.getKubeClient(runtimeId)
	if err != nil {
		return nil, err
	}

	podsClient := kubeClient.CoreV1().Pods(namespace)
	list, err := podsClient.List(metav1.ListOptions{LabelSelector: fmt.Sprintf("release=%s", clusterId)})
	if err != nil {
		return nil, err
	}

	var pbClusterNodes []*pb.ClusterNode
	for _, v := range list.Items {

		clusterNode := &models.ClusterNode{
			NodeId:     models.NewClusterNodeId(),
			ClusterId:  clusterId,
			Name:       v.GetName(),
			InstanceId: string(v.GetUID()),
			PrivateIp:  v.Status.PodIP,
			Status:     string(v.Status.Phase),
			Owner:      owner,
			//GlobalServerId: v.Spec.NodeName,
			CustomMetadata: getLabelString(v.GetObjectMeta().GetLabels()),
			CreateTime:     v.GetObjectMeta().GetCreationTimestamp().Time,
			StatusTime:     v.GetObjectMeta().GetCreationTimestamp().Time,
		}

		pbClusterNode := models.ClusterNodeToPb(clusterNode)
		pbClusterNodes = append(pbClusterNodes, pbClusterNode)
	}

	return pbClusterNodes, nil
}

func (p *Provider) UpdateClusterStatus(job *models.Job) error {
	clusterWrapper, err := models.NewClusterWrapper(job.Directive)
	if err != nil {
		return err
	}

	runtimeId := clusterWrapper.Cluster.RuntimeId

	runtime, err := runtimeclient.NewRuntime(runtimeId)
	if err != nil {
		return err
	}
	namespace := runtime.Zone

	pbClusterNodes, err := p.getKubePodsAsClusterNodes(runtimeId, namespace, job.ClusterId, job.Owner)
	if err != nil {
		return err
	}

	req2 := &pb.ModifyClusterRequest{
		Cluster: models.ClusterToPb(&models.Cluster{
			ClusterId: job.ClusterId,
		}),

		ClusterNodeSet: pbClusterNodes,
	}

	ctx := context.Background()
	clusterClient, err := clusterclient.NewClusterManagerClient(ctx)
	if err != nil {
		return err
	}
	_, err = clusterClient.ModifyCluster(ctx, req2)
	if err != nil {
		return err
	}

	return nil
}

func getLabelString(m map[string]string) string {
	b := new(bytes.Buffer)
	for k, v := range m {
		fmt.Fprintf(b, "%s=%s ", k, v)
	}
	return b.String()
}

func (p *Provider) ValidateCredential(url, credential string) error {
	kubeconfigGetter := func() (*clientcmdapi.Config, error) {
		return clientcmd.Load([]byte(credential))
	}

	config, err := clientcmd.BuildConfigFromKubeconfigGetter("", kubeconfigGetter)
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	cli := clientset.CoreV1().Namespaces()
	_, err = cli.List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (p *Provider) DescribeRuntimeProviderZones(url, credential string) []string {
	return nil
}
