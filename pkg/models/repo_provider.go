// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import "openpitrix.io/openpitrix/pkg/db"

var RepoProviderColumns = db.GetColumnsFromStruct(&RepoProvider{})

type RepoProvider struct {
	RepoId   string
	Provider string
}

func RepoProvidersMap(repoProviders []*RepoProvider) map[string][]*RepoProvider {
	providersMap := make(map[string][]*RepoProvider)
	for _, l := range repoProviders {
		repoId := l.RepoId
		providersMap[repoId] = append(providersMap[repoId], l)
	}
	return providersMap
}
