// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package k8s

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const test_appConf = `
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
	clientConf := "~/.kube/config"
	_, err := os.Stat(strings.Replace(clientConf, "~/", os.Getenv("HOME")+"/", 1))
	if err != nil {
		fmt.Printf("K8s runtime test skipped because no [%s] %s", clientConf, err)
		t.Skip()
	}

	runtime := K8sRuntime{}

	clusterId, err := runtime.CreateCluster(test_appConf, true)
	assert.Empty(t, err)
	err = runtime.StopClusters(clusterId, true, test_appConf)
	assert.Empty(t, err)
	err = runtime.StartClusters(clusterId, true, test_appConf)
	assert.Empty(t, err)
	err = runtime.DeleteClusters(clusterId, true, test_appConf)
	assert.Empty(t, err)
	err = runtime.RecoverClusters(clusterId, true, test_appConf)
	assert.Empty(t, err)
	err = runtime.CeaseClusters(clusterId, true, test_appConf)
	assert.Empty(t, err)
}
