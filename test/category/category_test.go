// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// +build integration

package category

import (
	"io/ioutil"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/require"

	"openpitrix.io/openpitrix/test/client/category_manager"
	"openpitrix.io/openpitrix/test/models"
	"openpitrix.io/openpitrix/test/testutil"
)

var clientConfig = testutil.GetClientConfig()

func getTestIcon(t *testing.T) strfmt.Base64 {
	b, err := ioutil.ReadFile("testdata/logo.png")
	require.NoError(t, err)
	return strfmt.Base64(b)
}

func getTestIcon1(t *testing.T) strfmt.Base64 {
	b, err := ioutil.ReadFile("testdata/logo1.png")
	require.NoError(t, err)
	return strfmt.Base64(b)
}

func TestCategory(t *testing.T) {
	client := testutil.GetClient(clientConfig)

	// delete old category
	testCategoryName := "test_category_name"
	testCategoryName2 := "test_category_name2"
	testCategoryLocale := "{}"
	describeParams := category_manager.NewDescribeCategoriesParams()
	describeParams.SetName([]string{testCategoryName})
	describeResp, err := client.CategoryManager.DescribeCategories(describeParams, nil)
	require.NoError(t, err)
	categories := describeResp.Payload.CategorySet
	for _, category := range categories {
		deleteParams := category_manager.NewDeleteCategoriesParams()
		deleteParams.SetBody(
			&models.OpenpitrixDeleteCategoriesRequest{
				CategoryID: []string{category.CategoryID},
			})
		_, err := client.CategoryManager.DeleteCategories(deleteParams, nil)
		require.NoError(t, err)
	}
	// create category
	createParams := category_manager.NewCreateCategoryParams()
	createParams.SetBody(
		&models.OpenpitrixCreateCategoryRequest{
			Name:   testCategoryName,
			Locale: testCategoryLocale,
			Icon:   getTestIcon(t),
		})
	createResp, err := client.CategoryManager.CreateCategory(createParams, nil)
	require.NoError(t, err)

	categoryId := createResp.Payload.CategoryID
	// modify category
	modifyParams := category_manager.NewModifyCategoryParams()
	modifyParams.SetBody(
		&models.OpenpitrixModifyCategoryRequest{
			CategoryID: categoryId,
			Name:       testCategoryName2,
			Icon:       getTestIcon1(t),
		})
	modifyResp, err := client.CategoryManager.ModifyCategory(modifyParams, nil)
	require.NoError(t, err)

	t.Log(modifyResp)
	// describe category
	describeParams.WithCategoryID([]string{categoryId})
	describeParams.WithName([]string{testCategoryName2})
	describeResp, err = client.CategoryManager.DescribeCategories(describeParams, nil)
	require.NoError(t, err)

	categories = describeResp.Payload.CategorySet

	require.Equal(t, 1, len(categories))
	require.Equal(t, categoryId, categories[0].CategoryID)
	require.Equal(t, testCategoryName2, categories[0].Name)

	// delete category
	deleteParams := category_manager.NewDeleteCategoriesParams()
	deleteParams.WithBody(&models.OpenpitrixDeleteCategoriesRequest{
		CategoryID: []string{categoryId},
	})
	deleteResp, err := client.CategoryManager.DeleteCategories(deleteParams, nil)
	require.NoError(t, err)

	t.Log(deleteResp)
	// describe deleted category
	describeParams.WithCategoryID([]string{categoryId})
	describeParams.WithName(nil)
	describeResp, err = client.CategoryManager.DescribeCategories(describeParams, nil)
	require.NoError(t, err)

	categories = describeResp.Payload.CategorySet
	require.Equal(t, 0, len(categories))

	t.Log("test category finish, all test is ok")
}
