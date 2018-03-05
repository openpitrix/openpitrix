package test

import (
	"testing"

	"openpitrix.io/openpitrix/pkg/constants"
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
			URL:         "http://q.io",
			Credential:  `{"accesskey": "ds", "secretkey": "ds"}`,
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
			URL:         "http://q.q",
			Credential:  `{"accesskey": "cc", "secretkey": "cc"}`,
			Visibility:  "labeled",
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
	if repos[0].Name != testRepoName || repos[0].Description != "cc" || repos[0].URL != "http://q.q" {
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
