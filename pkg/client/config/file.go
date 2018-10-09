// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package config

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type ClientConfig struct {
	ClientSecret string `json:"client_secret"`
	ClientID     string `json:"client_id"`
	EndpointURL  string `json:"endpoint_url"`
}

func (f *ClientConfig) tokenSource(ctx context.Context, scopes []string) (oauth2.TokenSource, error) {
	cfg := &clientcredentials.Config{
		ClientID:     f.ClientID,
		ClientSecret: f.ClientSecret,
		Scopes:       scopes,
		TokenURL:     f.EndpointURL + "/v1/oauth2/token",
	}

	oauth2.RegisterBrokenAuthHeaderProvider(f.EndpointURL)

	return cfg.TokenSource(ctx), nil
}

func GetClient(ctx context.Context, filename string) (*http.Client, error) {
	ts, err := GetTokenSource(ctx, filename)
	if err != nil {
		return nil, err
	}
	return oauth2.NewClient(ctx, ts), nil
}

func GetTokenSource(ctx context.Context, filename string) (oauth2.TokenSource, error) {
	f, err := ReadConfigFile(ctx, filename)
	if err != nil {
		return nil, err
	}
	return f.tokenSource(ctx, []string{""})
}

func ReadConfigFile(_ context.Context, filename string) (*ClientConfig, error) {
	jsonData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var f ClientConfig
	if err := json.Unmarshal(jsonData, &f); err != nil {
		return nil, err
	}
	return &f, nil
}
