// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package helm

import (
	"context"
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
	"openpitrix.io/openpitrix/pkg/util/funcutil"
	"openpitrix.io/openpitrix/pkg/util/stringutil"
)

var (
	NamespaceRegExp = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)
)

type KubeHandler struct {
	ctx       context.Context
	RuntimeId string
}

func GetKubeHandler(ctx context.Context, runtimeId string) *KubeHandler {
	kubeHandler := new(KubeHandler)
	kubeHandler.ctx = ctx
	kubeHandler.RuntimeId = runtimeId
	return kubeHandler
}

func (p *KubeHandler) initKubeClient() (*kubernetes.Clientset, *rest.Config, error) {
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
	logger.Debug(p.ctx, "Get runtime [%s] supported versions [%+v]", p.RuntimeId, supportedVersions)
	logger.Debug(p.ctx, "Check api versions [%+v]", apiVersions)
	for _, apiVersion := range apiVersions {
		if !stringutil.StringIn(apiVersion, supportedVersions) {
			return gerr.New(p.ctx, gerr.PermissionDenied, gerr.ErrorUnsupportedApiVersion, apiVersion)
		}
	}
	return nil
}

func (p *KubeHandler) WaitWorkloadReady(runtimeId, namespace string, clusterRoles map[string]*models.ClusterRole, timeout time.Duration, waitInterval time.Duration) error {
	err := funcutil.WaitForSpecificOrError(func() (bool, error) {
		for _, clusterRole := range clusterRoles {
			if clusterRole.Role == "" {
				continue
			}

			pods, err := p.getPodsByClusterRole(namespace, clusterRole)
			if err != nil {
				return true, err
			}

			if pods == nil {
				continue
			}

			if clusterRole.ReadyReplicas != clusterRole.Replicas {
				return false, nil
			}
		}

		return true, nil
	}, timeout, waitInterval)
	return err
}

func (p *KubeHandler) DescribeClusterDetails(clusterWrapper *models.ClusterWrapper) error {
	runtime, err := runtimeclient.NewRuntime(p.ctx, p.RuntimeId)
	if err != nil {
		return err
	}
	namespace := runtime.Zone

	for k, clusterRole := range clusterWrapper.ClusterRoles {
		if clusterRole.Role == "" {
			continue
		}

		pods, err := p.getPodsByClusterRole(namespace, clusterRole)
		if err != nil {
			return err
		}

		if pods == nil {
			continue
		}

		(*clusterWrapper).ClusterRoles[k] = clusterRole

		p.addPodsToClusterNodes(&clusterWrapper.ClusterNodesWithKeyPairs, pods, clusterWrapper.Cluster.ClusterId, clusterWrapper.Cluster.Owner, clusterRole.Role)
	}

	return nil
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

	var namespaces []string
	for _, ns := range out.Items {
		namespaces = append(namespaces, ns.Name)
	}

	return namespaces, nil
}

func (p *KubeHandler) getPodsByClusterRole(namespace string, clusterRole *models.ClusterRole) (*corev1.PodList, error) {
	kubeClient, _, err := p.initKubeClient()
	if err != nil {
		return nil, err
	}

	if strings.HasSuffix(clusterRole.Role, DeploymentFlag) {
		deploymentName := strings.TrimSuffix(clusterRole.Role, DeploymentFlag)
		switch clusterRole.ApiVersion {
		case "apps/v1":
			deployment, err := kubeClient.AppsV1().Deployments(namespace).Get(deploymentName, metav1.GetOptions{})
			if err != nil {
				return nil, err
			}

			(*clusterRole).ReadyReplicas = uint32(deployment.Status.ReadyReplicas)

			labelSelector := labels.Set(deployment.Spec.Selector.MatchLabels).AsSelector().String()
			pods, err := kubeClient.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: labelSelector})
			if err != nil {
				return nil, err
			}
			return pods, nil
		case "apps/v1beta2":
			deployment, err := kubeClient.AppsV1beta2().Deployments(namespace).Get(deploymentName, metav1.GetOptions{})
			if err != nil {
				return nil, err
			}

			(*clusterRole).ReadyReplicas = uint32(deployment.Status.ReadyReplicas)

			labelSelector := labels.Set(deployment.Spec.Selector.MatchLabels).AsSelector().String()
			pods, err := kubeClient.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: labelSelector})
			if err != nil {
				return nil, err
			}
			return pods, nil
		case "apps/v1beta1":
			deployment, err := kubeClient.AppsV1beta1().Deployments(namespace).Get(deploymentName, metav1.GetOptions{})
			if err != nil {
				return nil, err
			}

			(*clusterRole).ReadyReplicas = uint32(deployment.Status.ReadyReplicas)

			labelSelector := labels.Set(deployment.Spec.Selector.MatchLabels).AsSelector().String()
			pods, err := kubeClient.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: labelSelector})
			if err != nil {
				return nil, err
			}
			return pods, nil
		case "extensions/v1beta1":
			deployment, err := kubeClient.ExtensionsV1beta1().Deployments(namespace).Get(deploymentName, metav1.GetOptions{})
			if err != nil {
				return nil, err
			}

			(*clusterRole).ReadyReplicas = uint32(deployment.Status.ReadyReplicas)

			labelSelector := labels.Set(deployment.Spec.Selector.MatchLabels).AsSelector().String()
			pods, err := kubeClient.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: labelSelector})
			if err != nil {
				return nil, err
			}
			return pods, nil
		}
	} else if strings.HasSuffix(clusterRole.Role, StatefulSetFlag) {
		statefulSetName := strings.TrimSuffix(clusterRole.Role, StatefulSetFlag)

		switch clusterRole.ApiVersion {
		case "apps/v1":
			statefulSet, err := kubeClient.AppsV1().StatefulSets(namespace).Get(statefulSetName, metav1.GetOptions{})
			if err != nil {
				return nil, err
			}

			(*clusterRole).ReadyReplicas = uint32(statefulSet.Status.ReadyReplicas)

			labelSelector := labels.Set(statefulSet.Spec.Selector.MatchLabels).AsSelector().String()
			pods, err := kubeClient.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: labelSelector})
			if err != nil {
				return nil, err
			}
			return pods, nil
		case "apps/v1beta2":
			statefulSet, err := kubeClient.AppsV1beta2().StatefulSets(namespace).Get(statefulSetName, metav1.GetOptions{})
			if err != nil {
				return nil, err
			}

			(*clusterRole).ReadyReplicas = uint32(statefulSet.Status.ReadyReplicas)

			labelSelector := labels.Set(statefulSet.Spec.Selector.MatchLabels).AsSelector().String()
			pods, err := kubeClient.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: labelSelector})
			if err != nil {
				return nil, err
			}
			return pods, nil
		case "apps/v1beta1":
			statefulSet, err := kubeClient.AppsV1beta1().StatefulSets(namespace).Get(statefulSetName, metav1.GetOptions{})
			if err != nil {
				return nil, err
			}

			(*clusterRole).ReadyReplicas = uint32(statefulSet.Status.ReadyReplicas)

			labelSelector := labels.Set(statefulSet.Spec.Selector.MatchLabels).AsSelector().String()
			pods, err := kubeClient.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: labelSelector})
			if err != nil {
				return nil, err
			}
			return pods, nil
		}
	} else if strings.HasSuffix(clusterRole.Role, DaemonSetFlag) {
		daemonSetName := strings.TrimSuffix(clusterRole.Role, DaemonSetFlag)

		switch clusterRole.ApiVersion {
		case "apps/v1":
			daemonSet, err := kubeClient.AppsV1().DaemonSets(namespace).Get(daemonSetName, metav1.GetOptions{})
			if err != nil {
				return nil, err
			}

			(*clusterRole).Replicas = uint32(daemonSet.Status.DesiredNumberScheduled)
			(*clusterRole).ReadyReplicas = uint32(daemonSet.Status.NumberReady)

			labelSelector := labels.Set(daemonSet.Spec.Selector.MatchLabels).AsSelector().String()
			pods, err := kubeClient.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: labelSelector})
			if err != nil {
				return nil, err
			}
			return pods, nil
		case "apps/v1beta2":
			daemonSet, err := kubeClient.AppsV1beta2().DaemonSets(namespace).Get(daemonSetName, metav1.GetOptions{})
			if err != nil {
				return nil, err
			}

			(*clusterRole).Replicas = uint32(daemonSet.Status.DesiredNumberScheduled)
			(*clusterRole).ReadyReplicas = uint32(daemonSet.Status.NumberReady)

			labelSelector := labels.Set(daemonSet.Spec.Selector.MatchLabels).AsSelector().String()
			pods, err := kubeClient.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: labelSelector})
			if err != nil {
				return nil, err
			}
			return pods, nil
		case "extensions/v1beta1":
			daemonSet, err := kubeClient.ExtensionsV1beta1().DaemonSets(namespace).Get(daemonSetName, metav1.GetOptions{})
			if err != nil {
				return nil, err
			}

			(*clusterRole).Replicas = uint32(daemonSet.Status.DesiredNumberScheduled)
			(*clusterRole).ReadyReplicas = uint32(daemonSet.Status.NumberReady)

			labelSelector := labels.Set(daemonSet.Spec.Selector.MatchLabels).AsSelector().String()
			pods, err := kubeClient.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: labelSelector})
			if err != nil {
				return nil, err
			}
			return pods, nil
		}
	}

	return nil, nil
}

func (p *KubeHandler) addPodsToClusterNodes(clusterNodes *map[string]*models.ClusterNodeWithKeyPairs, pods *corev1.PodList, clusterId, owner, role string) {
	for _, pod := range pods.Items {

		clusterNode := &models.ClusterNodeWithKeyPairs{
			ClusterNode: &models.ClusterNode{
				NodeId:         models.NewClusterNodeId(),
				ClusterId:      clusterId,
				Name:           pod.GetName(),
				InstanceId:     string(pod.GetUID()),
				PrivateIp:      pod.Status.PodIP,
				Status:         string(pod.Status.Phase),
				Owner:          owner,
				Role:           role,
				CustomMetadata: GetLabelString(pod.GetObjectMeta().GetLabels()),
				CreateTime:     pod.GetObjectMeta().GetCreationTimestamp().Time,
				StatusTime:     pod.GetObjectMeta().GetCreationTimestamp().Time,
				HostId:         pod.Spec.NodeName,
				HostIp:         pod.Status.HostIP,
			},
		}

		//if len(pod.OwnerReferences) != 0 {
		//	clusterNode.Role = fmt.Sprintf("%s-%s", pod.OwnerReferences[0].Name, pod.OwnerReferences[0].Kind)
		//}

		(*clusterNodes)[clusterNode.NodeId] = clusterNode
	}
}
