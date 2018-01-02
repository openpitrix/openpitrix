package k8s_runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const appConf = `
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: mydeploy
  labels:
    app: myapp
spec:
  replicas: 1
  selector:
    matchLabels:
      app: mypod
  template:
    metadata:
      labels:
        app: mypod
    spec:
      containers:
      - name: myapp-container
        image: busybox
        command: ['sh', '-c', 'echo Hello Kubernetes2! && sleep 3600']
        resources:
          limits:
            cpu: 1024m
            memory: 1000Mi
          requests:
            cpu: 1024m
            memory: 1000Mi
---
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: yourdeploy
  labels:
    app: yourapp
spec:
  replicas: 1
  selector:
    matchLabels:
      app: yourpod
  template:
    metadata:
      labels:
        app: yourpod
    spec:
      containers:
      - name: yourapp-container
        image: busybox
        command: ['sh', '-c', 'echo Hello Kubernetes2! && sleep 3600']
        resources:
          limits:
            cpu: 1024m
            memory: 1000Mi
          requests:
            cpu: 1024m
            memory: 1000Mi
`

func TestK8sRuntime(t *testing.T) {
	runtime := K8sRuntime{}

	clusterId, err := runtime.CreateCluster(appConf, true)
	assert.Empty(t, err)
	err = runtime.StopClusters(clusterId, true, appConf)
	assert.Empty(t, err)
	err = runtime.StartClusters(clusterId, true, appConf)
	assert.Empty(t, err)
	err = runtime.DeleteClusters(clusterId, true, appConf)
	assert.Empty(t, err)
	err = runtime.RecoverClusters(clusterId, true, appConf)
	assert.Empty(t, err)
	err = runtime.CeaseClusters(clusterId, true, appConf)
	assert.Empty(t, err)
}
