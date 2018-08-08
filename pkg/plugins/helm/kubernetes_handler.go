// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package helm

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	runtimeclient "openpitrix.io/openpitrix/pkg/client/runtime"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/funcutil"
	"openpitrix.io/openpitrix/pkg/util/stringutil"
)

var (
	NamespaceRegExp = regexp.MustCompile(`[a-z0-9]([-a-z0-9]*[a-z0-9])?`)
)

type KubeHandler struct {
	Logger    *logger.Logger
	RuntimeId string
}

func GetKubeHandler(Logger *logger.Logger, runtimeId string) *KubeHandler {
	kubeHandler := new(KubeHandler)
	if Logger == nil {
		kubeHandler.Logger = logger.NewLogger()
	} else {
		kubeHandler.Logger = Logger
	}

	kubeHandler.RuntimeId = runtimeId
	return kubeHandler
}

func (p *KubeHandler) initKubeClient() (*kubernetes.Clientset, *rest.Config, error) {
	kubeconfigGetter := func() (*clientcmdapi.Config, error) {
		runtime, err := runtimeclient.NewRuntime(p.RuntimeId)
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

func (p *KubeHandler) initKubeClientWithCredential(credential string) (*kubernetes.Clientset, *rest.Config, error) {
	kubeconfigGetter := func() (*clientcmdapi.Config, error) {
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

func (p *KubeHandler) CheckApiVersionsSupported(apiVersions []string) error {
	if len(apiVersions) == 0 {
		return nil
	}

	client, _, err := p.initKubeClient()
	if err != nil {
		return err
	}
	apiGroupResources, err := discovery.GetAPIGroupResources(client)
	if err != nil {
		return err
	}
	var supportedVersions []string
	for _, group := range apiGroupResources {
		for _, version := range group.Group.Versions {
			supportedVersions = append(supportedVersions, version.GroupVersion)
		}
	}
	p.Logger.Debug("Get runtime [%s] supported versions [%+v]", p.RuntimeId, supportedVersions)
	p.Logger.Debug("Check api versions [%+v]", apiVersions)
	for _, apiVersion := range apiVersions {
		if !stringutil.StringIn(apiVersion, supportedVersions) {
			return gerr.New(gerr.PermissionDenied, gerr.ErrorUnsupportedApiVersion, apiVersion)
		}
	}
	return nil
}

func (p *KubeHandler) WaitPodsRunning(runtimeId, namespace string, clusterRoles map[string]*models.ClusterRole, timeout time.Duration, waitInterval time.Duration) error {
	err := funcutil.WaitForSpecificOrError(func() (bool, error) {
		for _, clusterRole := range clusterRoles {
			pods, err := p.getPodsByClusterRole(namespace, clusterRole)
			if err != nil {
				return true, err
			}

			if pods == nil {
				continue
			}

			if !p.checkPodsCount(pods, clusterRole.InstanceSize) {
				return false, nil
			}

			if !p.checkPodsRunning(pods) {
				return false, nil
			}
		}

		return true, nil
	}, timeout, waitInterval)
	return err
}

func (p *KubeHandler) getPodsByClusterRole(namespace string, clusterRole *models.ClusterRole) (pods *corev1.PodList, err error) {
	kubeClient, _, err := p.initKubeClient()
	if err != nil {
		return nil, err
	}

	if strings.HasSuffix(clusterRole.Role, DeploymentFlag) {
		deploymentName := strings.TrimSuffix(clusterRole.Role, DeploymentFlag)
		deployment, err := kubeClient.AppsV1().Deployments(namespace).Get(deploymentName, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}

		labelSelector := labels.Set(deployment.Spec.Selector.MatchLabels).AsSelector().String()
		pods, err = kubeClient.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: labelSelector})
		if err != nil {
			return nil, err
		}
	} else if strings.HasSuffix(clusterRole.Role, StatefulSetFlag) {
		statefulSetName := strings.TrimSuffix(clusterRole.Role, StatefulSetFlag)
		statefulSet, err := kubeClient.AppsV1().StatefulSets(namespace).Get(statefulSetName, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}

		labelSelector := labels.Set(statefulSet.Spec.Selector.MatchLabels).AsSelector().String()
		pods, err = kubeClient.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: labelSelector})
		if err != nil {
			return nil, err
		}

	} else if strings.HasSuffix(clusterRole.Role, DaemonSetFlag) {
		daemonSetName := strings.TrimSuffix(clusterRole.Role, DaemonSetFlag)
		daemonSet, err := kubeClient.AppsV1().DaemonSets(namespace).Get(daemonSetName, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}

		labelSelector := labels.Set(daemonSet.Spec.Selector.MatchLabels).AsSelector().String()
		pods, err = kubeClient.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: labelSelector})
		if err != nil {
			return nil, err
		}
	}
	return
}

func (p *KubeHandler) checkPodsCount(pods *corev1.PodList, count uint32) bool {
	return len(pods.Items) == int(count)
}

func (p *KubeHandler) checkPodsRunning(pods *corev1.PodList) bool {
	for _, pod := range pods.Items {
		if pod.Status.Phase != corev1.PodRunning {
			return false
		}
	}

	return true
}

func (p *KubeHandler) GetKubePodsAsClusterNodes(namespace, clusterId, owner string, clusterRoles map[string]*models.ClusterRole) ([]*pb.ClusterNode, error) {
	var pbClusterNodes []*pb.ClusterNode
	for _, clusterRole := range clusterRoles {
		pods, err := p.getPodsByClusterRole(namespace, clusterRole)
		if err != nil {
			return nil, err
		}

		if pods == nil {
			continue
		}

		p.appendKubePodsToClusterNodes(pbClusterNodes, pods, clusterId, owner)
	}

	return pbClusterNodes, nil
}

func (p *KubeHandler) appendKubePodsToClusterNodes(pbClusterNodes []*pb.ClusterNode, pods *corev1.PodList, clusterId, owner string) {
	for _, pod := range pods.Items {

		clusterNode := &models.ClusterNode{
			ClusterId:  clusterId,
			Name:       pod.GetName(),
			InstanceId: string(pod.GetUID()),
			PrivateIp:  pod.Status.PodIP,
			Status:     string(pod.Status.Phase),
			Owner:      owner,
			//GlobalServerId: pod.Spec.NodeName,
			CustomMetadata: GetLabelString(pod.GetObjectMeta().GetLabels()),
			CreateTime:     pod.GetObjectMeta().GetCreationTimestamp().Time,
			StatusTime:     pod.GetObjectMeta().GetCreationTimestamp().Time,
		}

		if len(pod.OwnerReferences) != 0 {
			clusterNode.Role = fmt.Sprintf("%s-%s", pod.OwnerReferences[0].Name, pod.OwnerReferences[0].Kind)
		}

		pbClusterNode := models.ClusterNodeToPb(clusterNode)
		pbClusterNodes = append(pbClusterNodes, pbClusterNode)
	}
}

func (p *KubeHandler) ValidateCredential(credential, zone string) error {
	if !NamespaceRegExp.MatchString(zone) {
		return fmt.Errorf(`namespace must match with regexp "[a-z0-9]([-a-z0-9]*[a-z0-9])?"`)
	}

	client, _, err := p.initKubeClientWithCredential(credential)
	if err != nil {
		return err
	}

	cli := client.CoreV1().Namespaces()
	_, err = cli.Get(KubeSystemNamespace, metav1.GetOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (p *KubeHandler) DescribeRuntimeProviderZones(credential string) ([]string, error) {
	client, _, err := p.initKubeClientWithCredential(credential)
	if err != nil {
		return nil, err
	}

	cli := client.CoreV1().Namespaces()
	out, err := cli.List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	namespaces := []string{}
	for _, ns := range out.Items {
		namespaces = append(namespaces, ns.Name)
	}

	return namespaces, nil
}
