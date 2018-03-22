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

func TestRepoLabel(t *testing.T) {
	client := GetClient(clientConfig)

	// Create a test repo that can attach label on it
	testRepoName := "e2e_test_repo"
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

	// create repo label on the repo just created
	createRepoLabelParams := repo_manager.NewCreateRepoLabelParams()
	createRepoLabelParams.SetBody(
		&models.OpenpitrixCreateRepoLabelRequest{
			RepoID:     repoId,
			LabelKey:   "department",
			LabelValue: "marketing",
		})
	createRepoLabelResp, err := client.RepoManager.CreateRepoLabel(createRepoLabelParams)
	if err != nil {
		t.Fatal(err)
	}
	repoLabelId := createRepoLabelResp.Payload.RepoLabel.RepoLabelID

	// modify repo label
	modifyRepoLabelParams := repo_manager.NewModifyRepoLabelParams()
	modifyRepoLabelParams.SetBody(
		&models.OpenpitrixModifyRepoLabelRequest{
			RepoLabelID: repoLabelId,
			LabelKey:    "department",
			LabelValue:  "develop",
		})
	modifyRepoLabelResp, err := client.RepoManager.ModifyRepoLabel(modifyRepoLabelParams)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(modifyRepoLabelResp)

	// check whether modification is successful
	describeRepoLabelParams := repo_manager.NewDescribeRepoLabelsParams()
	describeRepoLabelParams.WithRepoLabelID([]string{repoLabelId})
	describeRepoLabelResp, err := client.RepoManager.DescribeRepoLabels(describeRepoLabelParams)
	if err != nil {
		t.Fatal(err)
	}
	repoLabels := describeRepoLabelResp.Payload.RepoLabelSet
	if len(repoLabels) != 1 {
		t.Fatalf("failed to describe repo labels with params [%+v]", describeRepoLabelParams)
	}
	if repoLabels[0].LabelKey != "department" || repoLabels[0].LabelValue != "develop" {
		t.Fatalf("failed to modify repo label [%+v]", repoLabels[0])
	}

	// delete repo label
	deleteRepoLabelParams := repo_manager.NewDeleteRepoLabelParams()
	deleteRepoLabelParams.WithBody(&models.OpenpitrixDeleteRepoLabelRequest{
		RepoLabelID: repoLabelId,
	})
	deleteRepoLabelResp, err := client.RepoManager.DeleteRepoLabel(deleteRepoLabelParams)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(deleteRepoLabelResp)

	// describe deleted repo label
	describeRepoLabelParams.WithRepoLabelID([]string{repoLabelId})
	describeRepoLabelParams.WithStatus([]string{constants.StatusDeleted})
	describeRepoLabelResp, err = client.RepoManager.DescribeRepoLabels(describeRepoLabelParams)
	if err != nil {
		t.Fatal(err)
	}
	repoLabels = describeRepoLabelResp.Payload.RepoLabelSet
	if len(repoLabels) != 1 {
		t.Fatalf("failed to describe repo labels with params [%+v]", describeRepoLabelParams)
	}
	repoLabel := repoLabels[0]
	if repoLabel.RepoLabelID != repoLabelId {
		t.Fatalf("failed to describe repo label")
	}
	if repoLabel.Status != constants.StatusDeleted {
		t.Fatalf("failed to delete repo label, got repo status [%s]", repoLabel.Status)
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

	// create repo selector on the repo just created
	createRepoSelectorParams := repo_manager.NewCreateRepoSelectorParams()
	createRepoSelectorParams.SetBody(
		&models.OpenpitrixCreateRepoSelectorRequest{
			RepoID:        repoId,
			SelectorKey:   "runtime",
			SelectorValue: "aws",
		})
	createRepoSelectorResp, err := client.RepoManager.CreateRepoSelector(createRepoSelectorParams)
	if err != nil {
		t.Fatal(err)
	}
	repoSelectorId := createRepoSelectorResp.Payload.RepoSelector.RepoSelectorID

	// modify repo selector
	modifyRepoSelectorParams := repo_manager.NewModifyRepoSelectorParams()
	modifyRepoSelectorParams.SetBody(
		&models.OpenpitrixModifyRepoSelectorRequest{
			RepoSelectorID: repoSelectorId,
			SelectorKey:    "runtime",
			SelectorValue:  "qingcloud",
		})
	modifyRepoSelectorResp, err := client.RepoManager.ModifyRepoSelector(modifyRepoSelectorParams)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(modifyRepoSelectorResp)

	// check whether modification is successful
	describeRepoSelectorParams := repo_manager.NewDescribeRepoSelectorsParams()
	describeRepoSelectorParams.WithRepoSelectorID([]string{repoSelectorId})
	describeRepoSelectorResp, err := client.RepoManager.DescribeRepoSelectors(describeRepoSelectorParams)
	if err != nil {
		t.Fatal(err)
	}
	repoSelectors := describeRepoSelectorResp.Payload.RepoSelectorSet
	if len(repoSelectors) != 1 {
		t.Fatalf("failed to describe repo selectors with params [%+v]", describeRepoSelectorParams)
	}
	if repoSelectors[0].SelectorKey != "runtime" || repoSelectors[0].SelectorValue != "qingcloud" {
		t.Fatalf("failed to modify repo selector [%+v]", repoSelectors[0])
	}

	// delete repo selector
	deleteRepoSelectorParams := repo_manager.NewDeleteRepoSelectorParams()
	deleteRepoSelectorParams.WithBody(&models.OpenpitrixDeleteRepoSelectorRequest{
		RepoSelectorID: repoSelectorId,
	})
	deleteRepoSelectorResp, err := client.RepoManager.DeleteRepoSelector(deleteRepoSelectorParams)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(deleteRepoSelectorResp)

	// describe deleted repo selector
	describeRepoSelectorParams.WithRepoSelectorID([]string{repoSelectorId})
	describeRepoSelectorParams.WithStatus([]string{constants.StatusDeleted})
	describeRepoSelectorResp, err = client.RepoManager.DescribeRepoSelectors(describeRepoSelectorParams)
	if err != nil {
		t.Fatal(err)
	}
	repoSelectors = describeRepoSelectorResp.Payload.RepoSelectorSet
	if len(repoSelectors) != 1 {
		t.Fatalf("failed to describe repo selectors with params [%+v]", describeRepoSelectorParams)
	}
	repoSelector := repoSelectors[0]
	if repoSelector.RepoSelectorID != repoSelectorId {
		t.Fatalf("failed to describe repo selector")
	}
	if repoSelector.Status != constants.StatusDeleted {
		t.Fatalf("failed to delete repo selector, got repo status [%s]", repoSelector.Status)
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
