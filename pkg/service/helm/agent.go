package helm

import (
	"fmt"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/release"
)

//todo
func DescribeRelease(cfg *action.Configuration, releaseName string) (map[string]string, error) {
	rls, err := GetRelease(cfg, releaseName)
	if err != nil {
		return nil, fmt.Errorf("release [%s] not found", releaseName)
	}
	roles, _, err := ExactResources(rls)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func GetRelease(cfg *action.Configuration, releaseName string) (*release.Release, error) {
	cmd := action.NewGet(cfg)
	releaseInfo, err := cmd.Run(releaseName)
	if err != nil {
		return nil, err
	}
	return releaseInfo, nil
}

func ListReleases(cfg *action.Configuration, releaseName, ns, status string, offset, limit int) ([]*release.Release, error) {
	results := make([]*release.Release, 0)

	if releaseName != "" {
		res, err := GetRelease(cfg, releaseName)
		if err != nil {
			return nil, err
		}
		if status == "" || (status == res.Info.Status.String()) {
			results = append(results, res)
			return results, nil
		} else {
			return nil, fmt.Errorf("release [%s] not found", releaseName)
		}
	}
	cmd := action.NewList(cfg)
	cmd.Offset = offset
	cmd.Limit = limit
	allNamespaces := false
	if ns == "" {
		allNamespaces = true
	}
	if status == "deployed" {
		cmd.StateMask = action.ListDeployed
	} else if status == "failed" {
		cmd.StateMask = action.ListFailed
	} else {
		cmd.StateMask = action.ListAll
	}
	releases, err := cmd.Run()
	if err != nil {
		return nil, err
	}

	for _, r := range releases {
		if allNamespaces || r.Namespace == ns {
			results = append(results, r)
		}
	}
	return releases, nil
}

func CreateRelease(cfg *action.Configuration, releaseName, ns string, ch *chart.Chart, vals map[string]interface{}) (*release.Release, error) {
	cmd := action.NewInstall(cfg)
	cmd.Namespace = ns
	if releaseName == "" {
		cmd.GenerateName = true
	} else {
		cmd.ReleaseName = releaseName
	}

	rls, err := cmd.Run(ch, vals)
	if err != nil {

		return nil, err
	}
	return rls, nil
}

func UpgradeRelease(cfg *action.Configuration, name string, ch *chart.Chart, vals map[string]interface{}) (*release.Release, error) {
	_, err := GetRelease(cfg, name)
	if err != nil {

	}
	cmd := action.NewUpgrade(cfg)
	releaseInfo, err := cmd.Run(name, ch, vals)
	if err != nil {
		return nil, err
	}

	return releaseInfo, nil

}

func RollbackRelease(cfg *action.Configuration, name string, revision int32) (*release.Release, error) {
	cmd := action.NewRollback(cfg)
	cmd.Version = int(revision)
	err := cmd.Run(name)

	if err != nil {
		return nil, err
	}

	releaseInfo, err := GetRelease(cfg, name)
	return releaseInfo, err
}

func DeleteRelease(cfg *action.Configuration, name string, keepHistory bool) error {
	cmd := action.NewUninstall(cfg)
	cmd.KeepHistory = keepHistory
	_, err := cmd.Run(name)
	return err
}
