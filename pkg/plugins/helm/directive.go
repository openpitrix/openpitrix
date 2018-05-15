package helm

import (
	"encoding/json"
	"fmt"

	runtimeclient "openpitrix.io/openpitrix/pkg/client/runtime"
	"openpitrix.io/openpitrix/pkg/models"
)

type JobDirective struct {
	Namespace   string
	RuntimeId   string
	Values      string
	ClusterName string
}

func getJobDirective(data string) (*JobDirective, error) {
	clusterWrapper, err := models.NewClusterWrapper(data)
	if err != nil {
		return nil, err
	}

	runtimeId := clusterWrapper.Cluster.RuntimeId

	runtime, err := runtimeclient.NewRuntime(runtimeId)
	if err != nil {
		return nil, err
	}

	namespace := runtime.Zone
	clusterRole, ok := clusterWrapper.ClusterRoles[""]
	if !ok {
		return nil, fmt.Errorf("env is missing")
	}

	j := &JobDirective{
		Namespace:   namespace,
		RuntimeId:   runtimeId,
		Values:      clusterRole.Env,
		ClusterName: clusterWrapper.Cluster.Name,
	}

	return j, nil
}

type TaskDirective struct {
	VersionId   string
	Namespace   string
	RuntimeId   string
	Values      string
	ClusterName string
}

func getTaskDirectiveJson(v interface{}) string {
	r, _ := json.Marshal(v)
	return string(r)
}

func getTaskDirective(data string) (*TaskDirective, error) {
	var v TaskDirective
	err := json.Unmarshal([]byte(data), &v)
	if err != nil {
		return nil, err
	}
	return &v, nil
}
