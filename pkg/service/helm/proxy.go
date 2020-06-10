package helm

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/kube"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/util/funcutil"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
	"openpitrix.io/openpitrix/pkg/util/stringutil"

	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/release"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"openpitrix.io/openpitrix/pkg/models"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	runtimeclient "openpitrix.io/openpitrix/pkg/client/runtime"
)

const (
	Type       = "type"
	ExternalIp = "external_ip"
)

var (
	NamespaceReg    = `^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`
	NamespaceRegExp = regexp.MustCompile(NamespaceReg)

	ClusterNameReg    = `^[a-z]([-a-z0-9]*[a-z0-9])?$`
	ClusterNameRegExp = regexp.MustCompile(ClusterNameReg)
)

type CredentialGetter func(context.Context, string) (string, string, error)

var DefaultCredentialGetter CredentialGetter = func(ctx context.Context, runtimeId string) (string, string, error) {
	runtime, err := runtimeclient.NewRuntime(ctx, runtimeId)
	if err != nil {
		logger.Debug(nil, "get runtime error: [%s]", err.Error())
		return "", "", err
	}

	return runtime.Runtime.Zone, runtime.RuntimeCredentialContent, nil
}

//Helm kubernetes proxy
type Proxy struct {
	ctx          context.Context
	RuntimeId    string
	WorkloadInfo *Workload
}

type Workload struct {
	Deployments  []appsv1.Deployment  `json:"deployments,omitempty" description:"deployment list"`
	Statefulsets []appsv1.StatefulSet `json:"statefulsets,omitempty" description:"statefulset list"`
	Daemonsets   []appsv1.DaemonSet   `json:"daemonsets,omitempty" description:"daemonset list"`
	Services     []corev1.Service     `json:"services,omitempty" description:"application services"`
	Ingresses    []v1beta1.Ingress    `json:"ingresses,omitempty" description:"application ingresses"`
}

func NewProxy(ctx context.Context, runtimeId string) *Proxy {
	proxy := new(Proxy)
	proxy.RuntimeId = runtimeId
	proxy.ctx = ctx
	proxy.WorkloadInfo = new(Workload)
	return proxy
}

func (proxy *Proxy) GetKubeClient() (*kubernetes.Clientset, *rest.Config, error) {
	kubeconfigGetter := func() (*clientcmdapi.Config, error) {
		runtime, err := runtimeclient.NewRuntime(proxy.ctx, proxy.RuntimeId)
		if err != nil {
			return nil, err
		}

		return clientcmd.Load([]byte(runtime.RuntimeCredentialContent))
	}

	config, err := clientcmd.BuildConfigFromKubeconfigGetter("", kubeconfigGetter)
	if err != nil {
		return nil, nil, err
	}

	config.CAData = config.CAData[0:0]
	config.TLSClientConfig.Insecure = true

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}
	return clientset, config, err
}

func (proxy *Proxy) GetKubeClientWithCredential(credential string) (*kubernetes.Clientset, *rest.Config, error) {
	kubeconfigGetter := func() (*clientcmdapi.Config, error) {
		return clientcmd.Load([]byte(credential))
	}

	config, err := clientcmd.BuildConfigFromKubeconfigGetter("", kubeconfigGetter)
	if err != nil {
		return nil, nil, err
	}

	config.CAData = config.CAData[0:0]
	config.TLSClientConfig.Insecure = true

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}
	return clientset, config, err
}

func (proxy *Proxy) GetHelmConfig(driver, namespace string, getter CredentialGetter) (*action.Configuration, error) {
	file, err := ioutil.TempFile("", "config")
	defer os.Remove(file.Name())
	if err != nil {
		logger.Debug(proxy.ctx, "get helm config error: [%s]", err.Error())
		return nil, err
	}

	ns, credentialContent, err := getter(proxy.ctx, proxy.RuntimeId)
	if namespace != "" {
		ns = namespace
	}
	if len(credentialContent) == 0 {
		return nil, err
	}

	_, err = file.Write([]byte(credentialContent))
	if err != nil {
		logger.Debug(proxy.ctx, "write crendential content error: [%s]", err.Error())
		return nil, err
	}
	kubeConfigPath := file.Name()
	actionConfig := new(action.Configuration)

	// todo
	var FMT = func(format string, v ...interface{}) {
		fmt.Sprintf(format, v)
	}
	//todo context
	if err := actionConfig.Init(kube.GetConfig(kubeConfigPath, "", ns), ns, driver, FMT); err != nil {
		logger.Debug(proxy.ctx, "Init ActionConfig Error: [%s]", err.Error())
		return nil, err
	}
	return actionConfig, nil

}

func (proxy *Proxy) ListRelease(cfg *action.Configuration, name string) (*release.Release, error) {
	listCli := action.NewList(cfg)
	rls, err := listCli.Run()
	if err != nil {
		logger.Debug(proxy.ctx, "list release error: [%s]", err.Error())
		return nil, err
	}
	if len(rls) <= 0 {
		logger.Debug(proxy.ctx, "release: [%s] not found [%s]", name, err.Error())
		return nil, err
	}
	for _, r := range rls {
		if r.Name == name {
			return r, nil
		}
	}
	return nil, errors.New("release not found")

}

func (proxy *Proxy) InstallReleaseFromChart(cfg *action.Configuration, c *chart.Chart, rawVals map[string]interface{}, releaseName string) error {
	installCli := action.NewInstall(cfg)
	installCli.ReleaseName = releaseName
	rls, err := installCli.Run(c, rawVals)
	if err != nil {
		errDelete := proxy.DeleteRelease(cfg, rls.Name, true)
		if errDelete != nil && !strings.Contains(errDelete.Error(), "release: not found") {
			logger.Debug(proxy.ctx, "release: [%s] not found [%s]", releaseName, err.Error())
			return fmt.Errorf("Release %q failed: %v. Unable to delete failed release: %v", rls.Name, err, errDelete)
		}
		return err
	}
	return err
}

func (proxy *Proxy) UpdateReleaseFromChart(cfg *action.Configuration, releaseName string, c *chart.Chart, rawVals map[string]interface{}) error {
	_, err := proxy.ListRelease(cfg, releaseName)
	if err != nil {
		logger.Debug(proxy.ctx, "release: [%s] not found [%s]!", releaseName, err.Error())
		return err
	}
	updateCli := action.NewUpgrade(cfg)
	_, err = updateCli.Run(releaseName, c, rawVals)
	if err != nil {
		logger.Debug(proxy.ctx, "update release [%s] error [%s]", releaseName, err.Error())
		return err
	}
	return nil
}

func (proxy *Proxy) RollbackRelease(cfg *action.Configuration, releaseName string) error {
	rollbackCli := action.NewRollback(cfg)
	err := rollbackCli.Run(releaseName)
	if err != nil {
		logger.Debug(proxy.ctx, "rollback release [%s] error [%s]", releaseName, err.Error())
		return err
	}
	return nil
}

func (proxy *Proxy) DeleteRelease(cfg *action.Configuration, releaseName string, purge bool) error {
	deleteCli := action.NewUninstall(cfg)
	deleteCli.KeepHistory = true
	if purge {
		deleteCli.KeepHistory = false
	}
	_, err := deleteCli.Run(releaseName)
	if err != nil {
		logger.Debug(proxy.ctx, "delete release [%s] error [%s]", releaseName, err.Error())
		return err
	}
	return err
}

func (proxy *Proxy) ReleaseStatus(cfg *action.Configuration, releaseName string) (release.Status, error) {
	statusCli := action.NewStatus(cfg)
	rls, err := statusCli.Run(releaseName)
	if err != nil {
		logger.Debug(proxy.ctx, "get release [%s] status error [%s]", releaseName, err.Error())
		return "", err
	}
	return rls.Info.Status, nil
}

func (proxy *Proxy) CheckClusterNameIsUnique(clusterName, namespace string) error {
	if clusterName == "" {
		return fmt.Errorf("cluster name must be provided")
	}

	if !ClusterNameRegExp.MatchString(clusterName) {
		return fmt.Errorf(`cluster name must match with regexp "%s"`, ClusterNameReg)
	}

	// Related to https://github.com/helm/helm/pull/1080
	if len(clusterName) > 14 {
		return fmt.Errorf("the length of config [Name] must be less than 15")
	}

	err := funcutil.WaitForSpecificOrError(func() (bool, error) {
		//todo
		cfg, err := proxy.GetHelmConfig("", namespace, DefaultCredentialGetter)
		_, err = proxy.ReleaseStatus(cfg, clusterName)
		if err != nil {
			if isConnectionError(err) {
				return false, nil
			}
			return true, nil
		}

		return true, gerr.New(proxy.ctx, gerr.PermissionDenied, gerr.ErrorHelmReleaseExists, clusterName)
	}, constants.DefaultServiceTimeout, constants.WaitTaskInterval)
	return err
}

func (proxy *Proxy) DescribeVersionInfo() (*version.Info, error) {
	kubeClient, _, err := proxy.GetKubeClient()
	if err != nil {
		return nil, err
	}

	return kubeClient.ServerVersion()
}

func (proxy *Proxy) CheckApiVersionsSupported(apiVersions []string) error {
	if len(apiVersions) == 0 {
		return nil
	}

	_, config, err := proxy.GetKubeClient()
	if err != nil {
		return err
	}

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return err
	}

	apiGroupList, err := discoveryClient.ServerGroups()
	if err != nil {
		return err
	}
	var supportedVersions []string
	for _, group := range apiGroupList.Groups {
		for _, version := range group.Versions {
			supportedVersions = append(supportedVersions, version.GroupVersion)
		}
	}
	logger.Debug(proxy.ctx, "Get runtime [%s] supported versions [%+v]", proxy.RuntimeId, supportedVersions)
	logger.Debug(proxy.ctx, "Check api versions [%+v]", apiVersions)
	for _, apiVersion := range apiVersions {
		if !stringutil.StringIn(apiVersion, supportedVersions) {
			return gerr.New(proxy.ctx, gerr.PermissionDenied, gerr.ErrorUnsupportedApiVersion, apiVersion)
		}
	}
	return nil
}

func (proxy *Proxy) WaitWorkloadReady(runtimeId, namespace string, clusterRoles map[string]*models.ClusterRole, timeout time.Duration, waitInterval time.Duration) error {
	err := funcutil.WaitForSpecificOrError(func() (bool, error) {
		for _, clusterRole := range clusterRoles {

			pods, err := proxy.getPodsByClusterRole(namespace, clusterRole)
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

func (proxy *Proxy) describeAdditionalInfo(namespace string, cluster *models.Cluster) error {
	kubeClient, _, err := proxy.GetKubeClient()
	if err != nil {
		return err
	}

	var additionalInfo map[string][]map[string]interface{}
	err = jsonutil.Decode([]byte(cluster.AdditionalInfo), &additionalInfo)
	if err != nil {
		return err
	}

	for t, v := range additionalInfo {
		switch t {
		case "service":
			for i, svc := range v {
				service, err := kubeClient.CoreV1().Services(namespace).Get(svc["name"].(string), metav1.GetOptions{})
				if err != nil {
					return err
				}
				proxy.WorkloadInfo.Services = append(proxy.WorkloadInfo.Services, *service)

				additionalInfo[t][i][Type] = string(service.Spec.Type)
				additionalInfo[t][i]["cluster_ip"] = service.Spec.ClusterIP
				if service.Status.LoadBalancer.Ingress != nil && len(service.Status.LoadBalancer.Ingress) != 0 {
					additionalInfo[t][i][ExternalIp] = service.Status.LoadBalancer.Ingress[0].IP
				} else {
					if additionalInfo[t][i][Type] == "LoadBalancer" {
						additionalInfo[t][i][ExternalIp] = "pending"
					} else {
						additionalInfo[t][i][ExternalIp] = "none"
					}
				}

				ports := []string{}
				for _, port := range service.Spec.Ports {
					if port.NodePort == 0 {
						ports = append(ports, fmt.Sprintf("%d/%s", port.Port, port.Protocol))
					} else {
						ports = append(ports, fmt.Sprintf("%d:%d/%s", port.Port, port.NodePort, port.Protocol))
					}
				}
				additionalInfo[t][i]["ports"] = strings.Join(ports, ",")
			}
		case "configmap":
			for i, cm := range v {
				configMap, err := kubeClient.CoreV1().ConfigMaps(namespace).Get(cm["name"].(string), metav1.GetOptions{})
				if err != nil {
					return err
				}

				additionalInfo[t][i]["data_count"] = uint32(len(configMap.Data))
			}
		case "secret":
			for i, sec := range v {
				secret, err := kubeClient.CoreV1().Secrets(namespace).Get(sec["name"].(string), metav1.GetOptions{})
				if err != nil {
					return err
				}

				additionalInfo[t][i]["data_count"] = uint32(len(secret.Data))
			}
		case "pvc":
			for i, p := range v {
				pvc, err := kubeClient.CoreV1().PersistentVolumeClaims(namespace).Get(p["name"].(string), metav1.GetOptions{})
				if err != nil {
					return err
				}

				additionalInfo[t][i]["status"] = string(pvc.Status.Phase)
				additionalInfo[t][i]["volume"] = pvc.Spec.VolumeName
				additionalInfo[t][i]["capacity"] = pvc.Status.Capacity.StorageEphemeral().String()

				if len(pvc.Status.AccessModes) != 0 {
					additionalInfo[t][i]["access_mode"] = string(pvc.Status.AccessModes[0])
				}
			}
		case "ingress":
			for i, ing := range v {
				ingress, err := kubeClient.ExtensionsV1beta1().Ingresses(namespace).Get(ing["name"].(string), metav1.GetOptions{})
				if err != nil {
					return err
				}
				proxy.WorkloadInfo.Ingresses = append(proxy.WorkloadInfo.Ingresses, *ingress)

				hosts := []string{}
				for _, rule := range ingress.Spec.Rules {
					hosts = append(hosts, rule.Host)
				}

				additionalInfo[t][i]["hosts"] = strings.Join(hosts, ",")

				if ingress.Status.LoadBalancer.Ingress != nil && len(ingress.Status.LoadBalancer.Ingress) != 0 {
					additionalInfo[t][i]["address"] = ingress.Status.LoadBalancer.Ingress[0].IP
				}
			}
		}
	}

	(*cluster).AdditionalInfo = jsonutil.ToString(additionalInfo)

	return nil
}

func (proxy *Proxy) DescribeClusterDetails(clusterWrapper *models.ClusterWrapper) error {
	runtime, err := runtimeclient.NewRuntime(proxy.ctx, proxy.RuntimeId)
	if err != nil {
		return err
	}
	namespace := runtime.Zone
	if clusterWrapper.Cluster.Zone != "" {
		namespace = clusterWrapper.Cluster.Zone
	}

	for k, clusterRole := range clusterWrapper.ClusterRoles {

		pods, err := proxy.getPodsByClusterRole(namespace, clusterRole)
		if err != nil {
			return err
		}

		if pods == nil {
			continue
		}

		(*clusterWrapper).ClusterRoles[k] = clusterRole

		proxy.addPodsToClusterNodes(&clusterWrapper.ClusterNodesWithKeyPairs, pods, clusterWrapper.Cluster.ClusterId, clusterWrapper.Cluster.Owner, clusterRole.Role)
	}

	err = proxy.describeAdditionalInfo(namespace, clusterWrapper.Cluster)
	if err != nil {
		return err
	}

	return nil
}

func (proxy *Proxy) ValidateRuntime(zone string, runtimeCredential *models.RuntimeCredential, needCreate bool) error {
	if len(zone) == 0 {
		zone = "default"
	}
	if !NamespaceRegExp.MatchString(zone) {
		err := fmt.Errorf(`namespace must match with regexp "%s"`, NamespaceReg)
		return gerr.NewWithDetail(nil, gerr.PermissionDenied, err, gerr.ErrorNamespaceNotMatchWithRegex, zone, NamespaceReg)
	}
	client, _, err := proxy.GetKubeClientWithCredential(runtimeCredential.RuntimeCredentialContent)
	if err != nil {
		return gerr.NewWithDetail(nil, gerr.InvalidArgument, err, gerr.ErrorCredentialIllegal, "kubeconfig")
	}

	cli := client.CoreV1().Namespaces()
	if !needCreate {
		_, err = cli.Get(zone, metav1.GetOptions{})
		if err != nil {
			return gerr.NewWithDetail(nil, gerr.PermissionDenied, err, gerr.ErrorNamespaceUnavailable, zone)
		}
	} else {
		// create runtime
		// if not exist namespace, create new namespace with annotations
		// if exist namespace, should not with annotations
		namespace, err := cli.Get(zone, metav1.GetOptions{})
		if err != nil {
			logger.Info(proxy.ctx, "namespace [%s] not exist, need create", fmt.Sprintf("namespace: %s", zone))
			_, err = cli.Create(&corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: zone,
					Annotations: map[string]string{
						RuntimeAnnotationKey: proxy.RuntimeId,
					},
				},
			})
			if err != nil {
				return gerr.NewWithDetail(nil, gerr.Internal, err, gerr.ErrorCreateResourceFailed, zone)
			}
		} else {
			runtimeAnnotation, isExist := namespace.Annotations[RuntimeAnnotationKey]
			if isExist {
				err = fmt.Errorf("namespace [%s] annotations %s:%s already exist", zone, RuntimeAnnotationKey, runtimeAnnotation)
				return gerr.NewWithDetail(nil, gerr.AlreadyExists, err, gerr.ErrorNamespaceExists, zone)
			} else {
				logger.Info(proxy.ctx, "namespace [%s] exist, need update", zone)
				_, err = cli.Patch(zone, types.StrategicMergePatchType,
					[]byte(fmt.Sprintf(`{"metadata": {"annotations": {"%s": "%s"}}}`, RuntimeAnnotationKey, proxy.RuntimeId)))
				if err != nil {
					return gerr.NewWithDetail(nil, gerr.Internal, err, gerr.ErrorUpdateResourceFailed, fmt.Sprintf("namespace: %s", zone))
				}
			}
		}
	}

	return nil
}

func (proxy *Proxy) DescribeRuntimeProviderZones(runtimeCredential *models.RuntimeCredential) ([]string, error) {
	client, _, err := proxy.GetKubeClientWithCredential(runtimeCredential.RuntimeCredentialContent)
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

func (proxy *Proxy) getPodsByClusterRole(namespace string, clusterRole *models.ClusterRole) (*corev1.PodList, error) {
	kubeClient, _, err := proxy.GetKubeClient()
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
			proxy.WorkloadInfo.Deployments = append(proxy.WorkloadInfo.Deployments, *deployment)

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
			proxy.WorkloadInfo.Statefulsets = append(proxy.WorkloadInfo.Statefulsets, *statefulSet)

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
			proxy.WorkloadInfo.Daemonsets = append(proxy.WorkloadInfo.Daemonsets, *daemonSet)

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

func (proxy *Proxy) addPodsToClusterNodes(clusterNodes *map[string]*models.ClusterNodeWithKeyPairs, pods *corev1.PodList, clusterId, owner, role string) {
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
