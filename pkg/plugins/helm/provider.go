// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package helm

import (
	"bytes"
	"context"
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

	clientutil "openpitrix.io/openpitrix/pkg/client"
	appclient "openpitrix.io/openpitrix/pkg/client/app"
	clusterclient "openpitrix.io/openpitrix/pkg/client/cluster"
	runtimeclient "openpitrix.io/openpitrix/pkg/client/runtime"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/funcutil"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

type Provider struct {
	Logger *logger.Logger
}

func NewProvider(l *logger.Logger) *Provider {
	return &Provider{
		Logger: l,
	}
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
		return nil, fmt.Errorf("could not get a connection to tiller: %+v. ", err)
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
		return nil, fmt.Errorf("could not get a kube client: %+v. ", err)
	}

	helmClient, err = p.setupHelm(client, clientConfig, environment.DefaultTillerNamespace)
	return
}

func (p *Provider) checkClusterNameIsUniqueInRuntime(clusterName, runtimeId string) (err error) {
	hc, err := p.getHelmClient(runtimeId)
	if err != nil {
		return err
	}

	err = funcutil.WaitForSpecificOrError(func() (bool, error) {
		_, err := hc.ReleaseStatus(clusterName)
		if err != nil {
			if _, ok := err.(transport.ConnectionError); ok {
				return false, nil
			}
			return true, nil
		}

		return true, gerr.New(gerr.PermissionDenied, gerr.ErrorHelmReleaseExists, clusterName)
	}, constants.DefaultServiceTimeout, constants.WaitTaskInterval)
	return
}

func (p *Provider) ParseClusterConf(versionId, runtimeId, conf string) (*models.ClusterWrapper, error) {
	ctx := clientutil.GetSystemUserContext()
	appManagerClient, err := appclient.NewAppManagerClient()
	if err != nil {
		p.Logger.Error("Connect to app manager failed: %+v", err)
		return nil, err
	}

	req := &pb.GetAppVersionPackageRequest{
		VersionId: pbutil.ToProtoString(versionId),
	}

	resp, err := appManagerClient.GetAppVersionPackage(ctx, req)
	if err != nil {
		p.Logger.Error("Get app version [%s] package failed: %+v", versionId, err)
		return nil, err
	}

	pkg := resp.GetPackage()
	r := bytes.NewReader(pkg)

	c, err := chartutil.LoadArchive(r)
	if err != nil {
		p.Logger.Error("Load helm chart from app version [%s] failed: %+v", versionId, err)
		return nil, err
	}

	parser := Parser{Logger: p.Logger}
	clusterWrapper, err := parser.Parse(c, []byte(conf), versionId, runtimeId)
	if err != nil {
		p.Logger.Error("Parse app version [%s] failed: %+v", versionId, err)
		return nil, err
	}

	err = p.checkClusterNameIsUniqueInRuntime(clusterWrapper.Cluster.Name, runtimeId)
	if err != nil {
		p.Logger.Error("Check cluster name [%s] is unique in runtime [%s] failed: %+v", clusterWrapper.Cluster.Name, runtimeId, err)
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
		td := TaskDirective{
			VersionId:   job.VersionId,
			Namespace:   jodDirective.Namespace,
			RuntimeId:   jodDirective.RuntimeId,
			Values:      jodDirective.Values,
			ClusterName: jodDirective.ClusterName,
		}
		tdj := getTaskDirectiveJson(td)

		task := models.NewTask(constants.PlaceHolder, job.JobId, "", constants.ProviderKubernetes, constants.ActionCreateCluster, tdj, job.Owner, false)
		tl := models.TaskLayer{
			Tasks: []*models.Task{task},
			Child: nil,
		}

		return &tl, nil
	case constants.ActionUpgradeCluster:
		td := TaskDirective{
			VersionId:   job.VersionId,
			RuntimeId:   jodDirective.RuntimeId,
			Values:      jodDirective.Values,
			ClusterName: jodDirective.ClusterName,
		}
		tdj := getTaskDirectiveJson(td)

		task := models.NewTask(constants.PlaceHolder, job.JobId, "", constants.ProviderKubernetes, constants.ActionUpgradeCluster, tdj, job.Owner, false)
		tl := models.TaskLayer{
			Tasks: []*models.Task{task},
			Child: nil,
		}

		return &tl, nil
	case constants.ActionRollbackCluster:
		td := TaskDirective{ClusterName: jodDirective.ClusterName, RuntimeId: jodDirective.RuntimeId}
		tdj := getTaskDirectiveJson(td)

		task := models.NewTask(constants.PlaceHolder, job.JobId, "", constants.ProviderKubernetes, constants.ActionRollbackCluster, tdj, job.Owner, false)
		tl := models.TaskLayer{
			Tasks: []*models.Task{task},
			Child: nil,
		}

		return &tl, nil
	case constants.ActionDeleteClusters:
		td := TaskDirective{ClusterName: jodDirective.ClusterName, RuntimeId: jodDirective.RuntimeId}
		tdj := getTaskDirectiveJson(td)

		task := models.NewTask(constants.PlaceHolder, job.JobId, "", constants.ProviderKubernetes, constants.ActionDeleteClusters, tdj, job.Owner, false)
		tl := models.TaskLayer{
			Tasks: []*models.Task{task},
			Child: nil,
		}

		return &tl, nil
	case constants.ActionCeaseClusters:
		td := TaskDirective{ClusterName: jodDirective.ClusterName, RuntimeId: jodDirective.RuntimeId}
		tdj := getTaskDirectiveJson(td)

		task := models.NewTask(constants.PlaceHolder, job.JobId, "", constants.ProviderKubernetes, constants.ActionCeaseClusters, tdj, job.Owner, false)
		tl := models.TaskLayer{
			Tasks: []*models.Task{task},
			Child: nil,
		}

		return &tl, nil
	default:
		return nil, fmt.Errorf("the job action [%s] is not supported", job.JobAction)
	}
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

	ctx := clientutil.GetSystemUserContext()

	switch task.TaskAction {
	case constants.ActionCreateCluster:
		appc, err := appclient.NewAppManagerClient()
		if err != nil {
			return err
		}

		req := pb.GetAppVersionPackageRequest{
			VersionId: pbutil.ToProtoString(taskDirective.VersionId),
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
		err = jsonutil.Decode([]byte(taskDirective.Values), &v)
		if err != nil {
			return err
		}

		rawVals, err := yaml.Marshal(v)
		if err != nil {
			return err
		}

		_, err = hc.InstallReleaseFromChart(c, taskDirective.Namespace,
			helm.ValueOverrides(rawVals),
			helm.ReleaseName(taskDirective.ClusterName),
			helm.InstallWait(true))
		if err != nil {
			return err
		}
	case constants.ActionUpgradeCluster:
		appc, err := appclient.NewAppManagerClient()
		if err != nil {
			return err
		}

		req := pb.GetAppVersionPackageRequest{
			VersionId: pbutil.ToProtoString(taskDirective.VersionId),
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

		var v map[string]interface{}
		err = jsonutil.Decode([]byte(taskDirective.Values), &v)
		if err != nil {
			return err
		}

		rawVals, err := yaml.Marshal(v)
		if err != nil {
			return err
		}

		_, err = hc.UpdateReleaseFromChart(taskDirective.ClusterName, chart,
			helm.UpdateValueOverrides(rawVals),
			helm.UpgradeWait(true))
		if err != nil {
			return err
		}
	case constants.ActionRollbackCluster:
		_, err = hc.RollbackRelease(taskDirective.ClusterName, helm.RollbackWait(true))
		if err != nil {
			return err
		}
	case constants.ActionDeleteClusters:
		_, err = hc.DeleteRelease(taskDirective.ClusterName)
		if err != nil {
			return err
		}
	case constants.ActionCeaseClusters:
		_, err = hc.DeleteRelease(taskDirective.ClusterName, helm.DeletePurge(true))
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

	funcutil.WaitForSpecificOrError(func() (bool, error) {
		switch task.TaskAction {
		case constants.ActionCreateCluster:
			fallthrough
		case constants.ActionUpgradeCluster:
			fallthrough
		case constants.ActionRollbackCluster:
			resp, err := hc.ReleaseStatus(taskDirective.ClusterName)
			if err != nil {
				//network or api error, not considered task fail.
				return false, nil
			}

			if resp.Info.Status.Code == release.Status_DEPLOYED {
				return true, nil
			}
		case constants.ActionDeleteClusters:
			resp, err := hc.ReleaseStatus(taskDirective.ClusterName)
			if err != nil {
				//network or api error, not considered task fail.
				return false, nil
			}

			if resp.Info.Status.Code == release.Status_DELETED {
				return true, nil
			}
		case constants.ActionCeaseClusters:
			_, err := hc.ReleaseStatus(taskDirective.ClusterName)
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

func (p *Provider) DescribeSubnets(ctx context.Context, req *pb.DescribeSubnetsRequest) (*pb.DescribeSubnetsResponse, error) {
	return nil, nil
}

func (p *Provider) CheckResourceQuotas(ctx context.Context, clusterWrapper *models.ClusterWrapper) error {
	return nil
}

func (p *Provider) DescribeVpc(runtimeId, vpcId string) (*models.Vpc, error) {
	return nil, nil
}

func (p *Provider) getKubePodsAsClusterNodes(runtimeId, namespace, clusterId, clusterName, owner string) ([]*pb.ClusterNode, error) {
	kubeClient, _, err := p.getKubeClient(runtimeId)
	if err != nil {
		return nil, err
	}

	podsClient := kubeClient.CoreV1().Pods(namespace)
	list, err := podsClient.List(metav1.ListOptions{LabelSelector: fmt.Sprintf("release=%s", clusterName)})
	if err != nil {
		return nil, err
	}

	var pbClusterNodes []*pb.ClusterNode
	for _, pod := range list.Items {

		clusterNode := &models.ClusterNode{
			ClusterId:  clusterId,
			Name:       pod.GetName(),
			InstanceId: string(pod.GetUID()),
			PrivateIp:  pod.Status.PodIP,
			Status:     string(pod.Status.Phase),
			Owner:      owner,
			//GlobalServerId: pod.Spec.NodeName,
			CustomMetadata: getLabelString(pod.GetObjectMeta().GetLabels()),
			CreateTime:     pod.GetObjectMeta().GetCreationTimestamp().Time,
			StatusTime:     pod.GetObjectMeta().GetCreationTimestamp().Time,
		}

		if len(pod.OwnerReferences) != 0 {
			clusterNode.Role = fmt.Sprintf("%s-%s", pod.OwnerReferences[0].Name, pod.OwnerReferences[0].Kind)
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

	clusterName := clusterWrapper.Cluster.Name
	runtimeId := clusterWrapper.Cluster.RuntimeId

	runtime, err := runtimeclient.NewRuntime(runtimeId)
	if err != nil {
		return err
	}
	namespace := runtime.Zone

	pbClusterNodes, err := p.getKubePodsAsClusterNodes(runtimeId, namespace, job.ClusterId, clusterName, job.Owner)
	if err != nil {
		p.Logger.Error("Get kubernetes pods failed, %+v", err)
		return err
	}

	ctx := clientutil.GetSystemUserContext()
	clusterClient, err := clusterclient.NewClient()
	if err != nil {
		return err
	}

	// get all old node ids
	describeNodesRequest := &pb.DescribeClusterNodesRequest{
		ClusterId: pbutil.ToProtoString(job.ClusterId),
	}
	describeNodesResponse, err := clusterClient.DescribeClusterNodes(ctx, describeNodesRequest)
	if err != nil {
		p.Logger.Error("Get old nodes failed, %+v", err)
		return err
	}
	var nodeIds []string
	for _, clusterNode := range describeNodesResponse.ClusterNodeSet {
		nodeIds = append(nodeIds, clusterNode.GetNodeId().GetValue())
	}

	// delete old nodes from table
	deleteNodesRequest := &pb.DeleteTableClusterNodesRequest{
		NodeId: nodeIds,
	}
	clusterClient.DeleteTableClusterNodes(ctx, deleteNodesRequest)

	// add new nodes into table
	addNodesRequest := &pb.AddTableClusterNodesRequest{
		ClusterNodeSet: pbClusterNodes,
	}
	clusterClient.AddTableClusterNodes(ctx, addNodesRequest)

	return nil
}

func getLabelString(m map[string]string) string {
	b := new(bytes.Buffer)
	for k, v := range m {
		fmt.Fprintf(b, "%s=%s ", k, v)
	}
	return b.String()
}

func (p *Provider) ValidateCredential(url, credential, zone string) error {
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
	_, err = cli.Get(zone, metav1.GetOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (p *Provider) DescribeRuntimeProviderZones(url, credential string) ([]string, error) {
	return nil, nil
}
