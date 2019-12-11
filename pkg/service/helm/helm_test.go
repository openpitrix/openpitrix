package helm

import (
	"fmt"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/kube"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

const KEY string = `
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUN5RENDQWJDZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFWTVJNd0VRWURWUVFERXdwcmRXSmwKY201bGRHVnpNQjRYRFRJd01ESXhPREF5TXpZeU0xb1hEVE13TURJeE5UQXlNell5TTFvd0ZURVRNQkVHQTFVRQpBeE1LYTNWaVpYSnVaWFJsY3pDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBSmtjCllSL2d0Zm5jTis1UTN6NWNxdWhPd2pSRUxHbmpyd1RDOS82RXNiL0ZsbFkxMWhDNm5ma3RzVmE3Qy82UENUdjQKQStsNUVydjNLNXdsZ2UxaVM3NjN0L0wvanNiUmNwQTI0UFRISU9sUnk4eXgyZjJLTUdUSGYwYTU4T3VYT3RiTAp4UlVTdnFwUktDaktjK2lVUXZKdDlBSSs3M2JKOEpmN1RtK0UwWFpXZG0xRkY4UWo2MnRDQWRrdjFpL052c1ZvCkdxZFFENEd2dnpwQXpMaGZ0WVg2MDd6Wjl6U1NlWGZUeW9XVVVSY3BGcE1zZ0NTYmVVUklkcHgxRkdINU5qeTkKSTN6TzM4UkhIRVpnZ1VXN1k2Vk9wN0h6cUduYkgxcFZVWnlaYmhxTEJpL1hFblhxM2pSa2pLSEVMYmt1ejUxOQpWdUdLV0pxaytzcjM1MDAyaWdrQ0F3RUFBYU1qTUNFd0RnWURWUjBQQVFIL0JBUURBZ0trTUE4R0ExVWRFd0VCCi93UUZNQU1CQWY4d0RRWUpLb1pJaHZjTkFRRUxCUUFEZ2dFQkFJZnk3TGVUWnRyRHRpVlYzTlVQREpWV0U3ejMKTlAzM2xLZDh1UHRabUVCR05VMlIwRFpibHR5NXF3enBXMmg3ekIySzVYc3Zsa2gzZThMNDlIeGI1bTUvS3NFUgp0YUlKRVpGSEF2L29tbnhjdnJGOWlvYXVFb0JYQUxxbmp4WlZ0S3JsU3dybm52Um9Gc2pQWVJ5UXpkQW15Um9RCmpWaU1CZUVKZGJQQkc5aGhWWk5tc3JEaTJwL1lUSFl6TW1SWVVQM1k3Z1poc2RKSEpaVkxHalM0dllxdE50angKUjhkR2RCNDlPaTNTNkJhN0pUNi9RZXp1RXhiQ3lTa0EzSGVFM0dLenhlcXVXS0tQbTV0TjdrN1FoUnRRcW00aApsN1p3eEpmaFdlZFA2WEwvZWV5LytkcUVWelZpWkx5SWhVcGNNZnpUS3FqRWFTOTZkMUg5eThLMWZxbz0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=
    server: https://192.168.0.2:6443
  name: cluster.local
contexts:
- context:
    cluster: cluster.local
    user: kubernetes-admin
  name: kubernetes-admin@cluster.local
current-context: kubernetes-admin@cluster.local
kind: Config
preferences: {}
users:
- name: kubernetes-admin
  user:
    client-certificate-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUM4akNDQWRxZ0F3SUJBZ0lJY0dLaEZYSVk0Nnd3RFFZSktvWklodmNOQVFFTEJRQXdGVEVUTUJFR0ExVUUKQXhNS2EzVmlaWEp1WlhSbGN6QWVGdzB5TURBeU1UZ3dNak0yTWpOYUZ3MHlNVEF5TVRjd01qTTJNalphTURReApGekFWQmdOVkJBb1REbk41YzNSbGJUcHRZWE4wWlhKek1Sa3dGd1lEVlFRREV4QnJkV0psY201bGRHVnpMV0ZrCmJXbHVNSUlCSWpBTkJna3Foa2lHOXcwQkFRRUZBQU9DQVE4QU1JSUJDZ0tDQVFFQXl3Tlhjd051UHpEandURU0KdG43WkxTa3UrOVl0TVk2SnZPaVhwems2TVdDNWxVN01DNzE2clo0Y1NsUFRzQ3VIakZVVXJUYjIxaVhIZC9iNgpwazFXUFplckZ4NHN2ZTYwNnUzUUVqVG5SRWZ4OHdUV0swSmJLRGNnV0czSHExRjB4Y1RkQjY3blRLbVZ5SmtTCmwvRmd6cWRVTERNMXVTRFNpVmd1aElWVE4ySEJOYi96TCs1MldUMTc4TG5PNVY2cTBSR2kya21zamdheU01ZEEKQ2xsNUNWNlhwTG5HV3N0Z29pRElWbWp1cnZZTGVVQ1Q5cFFXL0E4VmNxRzJrTkFuWW9DWVFuZUZCRHJ5M2dKcQo0bkZNTkRWbjZJdVVlNHF0K0g3WE1yTnJBTlIvZS8zSkc2U25ZZHhkUzVRYXFaaVg1RmNiNkVyMU5wck0yK3c4CnRtSWhZd0lEQVFBQm95Y3dKVEFPQmdOVkhROEJBZjhFQkFNQ0JhQXdFd1lEVlIwbEJBd3dDZ1lJS3dZQkJRVUgKQXdJd0RRWUpLb1pJaHZjTkFRRUxCUUFEZ2dFQkFCbXJlQ1MvTGt5d2lGcVN4WW80TTJlZC8xaHI4UWZWOWg2cwpZZlA0Y251dzZJUjMvKzR1U2hoUFNiZDI5c0RlUktCbEFsS1dDekNLM0JvN211UFRkOVhtYWRTbGxWaEVUUzV6Ck1TdlE1T1JMdUNWOTdPSERDSHRaSHozZ2YzN1JUTGdTTzRkdEwvSkJrOVdQOUxwUlBmL3REaU9sYVZVd2t1aUEKN011akZ1OVVkWE1FZ0Y0bFVUekp0ZW1tVkJVK0FUazZvbmRhZVJHaUtWS0pya0VySzJaTVFQbzQxeFIyU3ZzdwpBMzBXeEx4dlZDV1IyeXVtOW9IUVdrdWFGemwyTCtzenFueFhxVjA3clVsem1WSEFLZ2Y1NTZBTU9PYlBEeWs4CmY3cTlCeGdEeDRvWjRNN1lTSXc4OGRnblQwNmdsdEhMUW9kd1dmKzVIdURscG9LSDlYbz0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=
    client-key-data: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcEFJQkFBS0NBUUVBeXdOWGN3TnVQekRqd1RFTXRuN1pMU2t1KzlZdE1ZNkp2T2lYcHprNk1XQzVsVTdNCkM3MTZyWjRjU2xQVHNDdUhqRlVVclRiMjFpWEhkL2I2cGsxV1BaZXJGeDRzdmU2MDZ1M1FFalRuUkVmeDh3VFcKSzBKYktEY2dXRzNIcTFGMHhjVGRCNjduVEttVnlKa1NsL0ZnenFkVUxETTF1U0RTaVZndWhJVlROMkhCTmIvegpMKzUyV1QxNzhMbk81VjZxMFJHaTJrbXNqZ2F5TTVkQUNsbDVDVjZYcExuR1dzdGdvaURJVm1qdXJ2WUxlVUNUCjlwUVcvQThWY3FHMmtOQW5Zb0NZUW5lRkJEcnkzZ0pxNG5GTU5EVm42SXVVZTRxdCtIN1hNck5yQU5SL2UvM0oKRzZTbllkeGRTNVFhcVppWDVGY2I2RXIxTnByTTIrdzh0bUloWXdJREFRQUJBb0lCQVFDQnlLUWFVZ0lrQWJSSQpxSHZSRzJ6cHN4OW5QamZzSzR5Z3FTMXlhV0pyZU1PTDBURWUvRVkyUWhNaDdVOHltOUZ2QkdGUWp3ZmtSWWlzCmg4Y2JrK3RqT3RmVTBxU1YwOG56T284L1pIVElzUm5iVzZjelJwdVNMUlBQbEhjR2JlK3lFeldlbU5FanNIS2kKS0VHN3cwTTVPYjNVOS9RTFl4RlZYbnQybXVsbFNFUEJHbTNnR25PSGtvL3YrTHp0TEFES1hiNWsxSGxyTG9uOApqeXo1SllRMXVXY2hBcEZ6OXB3aHdTVmZzOXgzK1RBNUlLWnQrMDlxWHhpRWdVa2lDV1Q3c055cWt0ZmFrVytyCjBrUmJOeXF2STU1UHNJS3V0ZVVsbHY0QU9mZGtBM1BkVm9GN3N2RGt4bFFURmN2MGxxQmNZS21HKzJid0dGZmIKYXFqcXVGVDVBb0dCQU5HUjVTc2V2K0ZHdENZQXNPRnh2RGhaOXNBZmNQUkRxOC85VDZaN2NUbGF2YTBKQmVlVgpoNUE2UFB4bXpPaGtyT0FIWnV1TGxFUmE4LzlKRXlobytUTnlPcEpGWnlTNEVMcFpKUVVyS0Y0eWI3WnUzNUZ1CkJETHQyUjFQUEdMN1czQmE2S2QzdlZQRG14OVVZNzNUbUFsQ1hVaW0yUG5FT0U2NWN3ZVlTZEZWQW9HQkFQZjkKa0ZRZHhGdTc0d2ZZZ2FLMzFsajd2aEttR21qTWJnQmlPc0pPS3BxRmQyUnJTWDBITjkvRzQ4N2RHelV6Q2R4SQo3THNJTGlYbXdNdm1zQkp6d3hhalZCU3hjL1U4NEJpZzBEajhkaTB3M0tQNFJVMk1IUVp4Zi9rN1hGZ2UwejJaClI5UTZoSnpEOFB3MkFMY0NIVXM1aHhCSDFrT0Fhdzg5MyttRkFnZlhBb0dCQUtBR01xTEJnUzFJNnVpVjRIclYKZVM3aWEwdDY5cXBlUGdTODNhUTNZRmEyVmwyWnBUdVg3NE1QSldCcU13OUZTTWhzZm9kZjlxQlhmemN4R29MaAozV0FPV25FMHM3VFRKRnJYRlRDa0t0ZjYxVmpOd3NOdTZaL21CTUtmclhHN2s2L3dpdlROdHZFT1RSWVlQMjFFCjlEUWx5OHRkTkJOTVpONmdOeGpXalk5ZEFvR0FYLzdtbExrbEhvRi9vN1Rha2J0cUhPM3VLTmZsbHpXelN6QzcKSUNZVDlkYStYYi91SlpqYXR5UU5ZVEZUNitjQzVTUFJoNkRtQkVQcjA4SkwzQWkxdHhpb1hvNUduZUxmdUlqZgpzWCtBMjROemxZRndpbEUzbHh2dWR2TFVqMFAzYjN2YlF6c1h4SHRRMk1DcXpDemtYQTg3eWtDVW4zS2hmcmZyCjZrQlRoZWNDZ1lCdkp0MWQwQVYxd3FTUUJnelBQeWdNV0RpUXRuWDZNZHFEbGdLYnhRMjI3MVJ5L2g1QWJzekwKWlVWb292aEVsU3dDWWdNSlY3eUc1UGVMd2lsVjZwcytOd3RxeXUrWEJLWUpoUDFORXNyeWZFT1RwVVlScjJxeQpaNzNDVElvaTRXQWRqTy9MUklrMlkxRkJzZ05qTkg0WVArVFFEUGc3d3dOeFRHU3B2SGlGMEE9PQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo=
`

var getter CredentialGetter = func(credential string) (string, error) {
	return string(KEY),nil
}

func TestHelmv3(t *testing.T) {
	chartPath := "/root/op/nginx-1.2.0.tgz"
	chart, err := loader.Load(chartPath)
	if err != nil {
		panic(err)
	}

	releaseName := fmt.Sprintf("test-helm-%d", time.Now().Unix())
	releaseNamespace := "helmv3"
	proxy := NewProxy(nil, "")
	cfg, err := proxy.GetHelmConfig(releaseNamespace, "configmap",getter)
	if err != nil {
		fmt.Println("config error:"+err.Error(), cfg)
	}
	err = proxy.InstallReleaseFromChart(cfg, chart, nil, releaseName)
	if err != nil {
		fmt.Printf("install release error: %s", err.Error())
	}
	rls,err := proxy.ListRelease(cfg,releaseName)
	if err != nil{
		fmt.Println("list release error: " + err.Error())
	}
	fmt.Println(rls.Name,rls.Version)
	stats,err := proxy.ReleaseStatus(cfg,releaseName)
	if err != nil{
		fmt.Println("release status error : " + err.Error())
	}
	fmt.Println("status " + stats.String())
	err = proxy.DeleteRelease(cfg,releaseName,true)
	if err != nil{
		fmt.Println("delete error: " + err.Error())
	}
	rls,err = proxy.ListRelease(cfg,releaseName)
	if err != nil{
		fmt.Println("list deleted release error: " + err.Error())
	}
	fmt.Println("release name is: " + rls.Name)

}

func VV(t *testing.T) {
	chartPath := "/root/op/nginx-1.2.0.tgz"
	chart, err := loader.Load(chartPath)
	if err != nil {
		panic(err)
	}

	kubeconfigPath := "/root/.kube/config"
	releaseName := "test-release"
	releaseNamespace := "helmv3"

	actionConfig := new(action.Configuration)
	// ues kubeconfig path
	if err := actionConfig.Init(kube.GetConfig(kubeconfigPath, "", releaseNamespace), releaseNamespace, os.Getenv("HELM_DRIVER"), func(format string, v ...interface{}) {
		fmt.Sprintf(format, v)
	}); err != nil {
		panic(err)
	}

	// use credential
	//actionConfig = NewActionConfig(false, []byte(KEY))
	actionConfig = NewActionConfigPath(releaseNamespace, []byte(KEY))

	//Delete
	uninstallClient := action.NewUninstall(actionConfig)
	uninstallClient.KeepHistory = false

	//Install
	iCli := action.NewInstall(actionConfig)

	iCli.Namespace = releaseNamespace
	iCli.ReleaseName = fmt.Sprintf("%s-%d", "test-helmv3", time.Now().Unix())
	//iCli.GenerateName = false
	_, _ = uninstallClient.Run(releaseName)
	rel, err := iCli.Run(chart, nil)
	if err != nil {
		fmt.Println("err: ", err.Error())
		panic(err)
	}
	fmt.Println("Successfully installed release: ", rel.Name)

	//List
	listClient := action.NewList(actionConfig)
	listClient.All = true
	lists, err := listClient.Run()
	if err != nil {
		panic(err)
	}
	for _, release := range lists {
		fmt.Println("list release:")
		fmt.Printf("release name: %s\n", release.Name)
	}

	//Status
	statusClient := action.NewStatus(actionConfig)
	for _, release := range lists {
		rel, err := statusClient.Run(release.Name)
		if err != nil {
			panic(err)
		}
		fmt.Printf("release %s,status: %s\n", rel.Name, rel.Info.Status)
	}

	//Uninstall
	resp, err := uninstallClient.Run(lists[0].Name)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp.Release.Name, resp.Info)

	lists, err = listClient.Run()
	if err != nil {
		panic(err)
	}

	for _, release := range lists {
		fmt.Printf("release name: %s", release.Name)
		rel, err := statusClient.Run(release.Name)
		if err != nil {
			panic(err)
		}
		fmt.Printf("release %s,status: %s", rel.Name, rel.Info.Status)
	}

	//clean release
	for _, release := range lists {
		resp, _ = uninstallClient.Run(release.Name)
		fmt.Printf("delete release %s\n", resp.Release.Name)
	}

}

func NewActionConfigPath(ns string, credentialContent []byte) *action.Configuration {
	file, err := ioutil.TempFile("", "config")
	if err != nil {
		fmt.Println("error1: " + err.Error())
	}

	_, err = file.Write(credentialContent)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	kubeConfigPath := file.Name()
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(kube.GetConfig(kubeConfigPath, "", ns), ns, os.Getenv("HELM_DRIVER"), func(format string, v ...interface{}) {
		fmt.Sprintf(format, v)
	}); err != nil {
		panic(err)
	}
	return actionConfig
}
