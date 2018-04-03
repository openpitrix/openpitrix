package test

import (
	"net/url"
	"testing"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/utils"
	"openpitrix.io/openpitrix/test/client/repo_manager"
	"openpitrix.io/openpitrix/test/models"
)

//var clientConfig = &ClientConfig{}
//
//func init() {
//	clientConfig = GetClientConfig()
//	log.Printf("Got Client Config: %+v", clientConfig)
//}

//func TestMain(m *testing.M) {
//	os.Exit(m.Run())
//}

func TestRepo(t *testing.T) {
	client := GetClient(clientConfig)

	// delete old repo
	testRepoName := "e2e_test_repo"
	describeParams := repo_manager.NewDescribeReposParams()
	describeParams.SetName([]string{testRepoName})
	describeParams.SetStatus([]string{constants.StatusActive})
	describeResp, err := client.RepoManager.DescribeRepos(describeParams)
	if err != nil {
		t.Fatal(err)
	}
	repos := describeResp.Payload.RepoSet
	for _, repo := range repos {
		deleteParams := repo_manager.NewDeleteRepoParams()
		deleteParams.SetBody(
			&models.OpenpitrixDeleteRepoRequest{
				RepoID: repo.RepoID,
			})
		_, err := client.RepoManager.DeleteRepo(deleteParams)
		if err != nil {
			t.Fatal(err)
		}
	}
	// create repo
	createParams := repo_manager.NewCreateRepoParams()
	createParams.SetBody(
		&models.OpenpitrixCreateRepoRequest{
			Name:        testRepoName,
			Description: "description",
			Type:        "https",
			URL:         "https://github.com/",
			Credential:  `{}`,
			Visibility:  "public",
		})
	createResp, err := client.RepoManager.CreateRepo(createParams)
	if err != nil {
		t.Fatal(err)
	}
	repoId := createResp.Payload.Repo.RepoID
	// modify repo
	modifyParams := repo_manager.NewModifyRepoParams()
	modifyParams.SetBody(
		&models.OpenpitrixModifyRepoRequest{
			RepoID:      repoId,
			Description: "cc",
			Type:        "https",
			URL:         "https://github.com/",
			Credential:  `{}`,
			Visibility:  "private",
		})
	modifyResp, err := client.RepoManager.ModifyRepo(modifyParams)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(modifyResp)
	// describe repo
	describeParams.WithRepoID([]string{repoId})
	describeResp, err = client.RepoManager.DescribeRepos(describeParams)
	if err != nil {
		t.Fatal(err)
	}
	repos = describeResp.Payload.RepoSet
	if len(repos) != 1 {
		t.Fatalf("failed to describe repos with params [%+v]", describeParams)
	}
	if repos[0].Name != testRepoName || repos[0].Description != "cc" || repos[0].URL != "https://github.com/" {
		t.Fatalf("failed to modify repo [%+v]", repos[0])
	}
	// delete repo
	deleteParams := repo_manager.NewDeleteRepoParams()
	deleteParams.WithBody(&models.OpenpitrixDeleteRepoRequest{
		RepoID: repoId,
	})
	deleteResp, err := client.RepoManager.DeleteRepo(deleteParams)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(deleteResp)
	// describe deleted repo
	describeParams.WithRepoID([]string{repoId})
	describeParams.WithStatus([]string{constants.StatusDeleted})
	describeParams.WithName(nil)
	describeResp, err = client.RepoManager.DescribeRepos(describeParams)
	if err != nil {
		t.Fatal(err)
	}
	repos = describeResp.Payload.RepoSet
	if len(repos) != 1 {
		t.Fatalf("failed to describe repos with params [%+v]", describeParams)
	}
	repo := repos[0]
	if repo.RepoID != repoId {
		t.Fatalf("failed to describe repo")
	}
	if repo.Status != constants.StatusDeleted {
		t.Fatalf("failed to delete repo, got repo status [%s]", repo.Status)
	}

	t.Log("test repo finish, all test is ok")
}

func generateLabels() string {
	v := url.Values{}
	v.Add("key1", utils.GetUuid(""))
	v.Add("key1", utils.GetUuid(""))
	v.Add("key1", utils.GetUuid(""))
	v.Add("key3", utils.GetUuid(""))
	v.Add("key4", utils.GetUuid(""))
	v.Add("key5", utils.GetUuid(""))
	v.Add("key6", utils.GetUuid(""))
	return v.Encode()
}

func TestRepoLabel(t *testing.T) {
	client := GetClient(clientConfig)
	// Create a test repo that can attach label on it
	testRepoName := "e2e_test_repo"
	labels := generateLabels()
	createParams := repo_manager.NewCreateRepoParams()
	createParams.SetBody(
		&models.OpenpitrixCreateRepoRequest{
			Name:        testRepoName,
			Description: "description",
			Type:        "https",
			URL:         "https://github.com/",
			Credential:  `{}`,
			Visibility:  "public",
			Labels:      labels,
		})
	createResp, err := client.RepoManager.CreateRepo(createParams)
	if err != nil {
		t.Fatal(err)
	}
	repoId := createResp.Payload.Repo.RepoID

	describeParams := repo_manager.NewDescribeReposParams()
	describeParams.Label = &labels
	describeParams.Status = []string{constants.StatusActive}
	describeResp, err := client.RepoManager.DescribeRepos(describeParams)
	if err != nil {
		t.Fatal(err)
	}
	if describeResp.Payload.RepoSet[0].RepoID != repoId {
		t.Fatalf("describe repo with filter failed")
	}

	// delete repo
	deleteParams := repo_manager.NewDeleteRepoParams()
	deleteParams.WithBody(&models.OpenpitrixDeleteRepoRequest{
		RepoID: repoId,
	})
	deleteResp, err := client.RepoManager.DeleteRepo(deleteParams)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(deleteResp)

	t.Log("test repo label finish, all test is ok")
}

func TestRepoSelector(t *testing.T) {
	client := GetClient(clientConfig)

	// Create a test repo that can attach selector on it
	testRepoName := "e2e_test_repo"
	labels := generateLabels()
	createParams := repo_manager.NewCreateRepoParams()
	createParams.SetBody(
		&models.OpenpitrixCreateRepoRequest{
			Name:        testRepoName,
			Description: "description",
			Type:        "https",
			URL:         "https://github.com/",
			Credential:  `{}`,
			Visibility:  "public",
			Selectors:   labels,
		})
	createResp, err := client.RepoManager.CreateRepo(createParams)
	if err != nil {
		t.Fatal(err)
	}
	repoId := createResp.Payload.Repo.RepoID

	describeParams := repo_manager.NewDescribeReposParams()
	describeParams.Selector = &labels
	describeParams.Status = []string{constants.StatusActive}
	describeResp, err := client.RepoManager.DescribeRepos(describeParams)
	if err != nil {
		t.Fatal(err)
	}
	if describeResp.Payload.RepoSet[0].RepoID != repoId {
		t.Fatalf("describe repo with filter failed")
	}

	// delete repo
	deleteParams := repo_manager.NewDeleteRepoParams()
	deleteParams.WithBody(&models.OpenpitrixDeleteRepoRequest{
		RepoID: repoId,
	})
	deleteResp, err := client.RepoManager.DeleteRepo(deleteParams)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(deleteResp)

	t.Log("test repo selector finish, all test is ok")
}
