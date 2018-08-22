// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package runtime

import (
	"context"
	"testing"
)

var ctx = context.TODO()

func TestValidateName(t *testing.T) {
	validNames := []string{"aa", "validatename!2313", "2131", "!!,."}
	invalidNames := []string{"", "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"}
	for _, validName := range validNames {
		err := ValidateName(ctx, validName)
		if err != nil {
			t.Fatalf("%+v should be validName", validName)
		}
	}
	for _, invalidName := range invalidNames {
		err := ValidateName(ctx, invalidName)
		if err == nil {
			t.Fatalf("%+v should be invalidName", invalidName)
		}
	}
}

func TestValidateURL(t *testing.T) {
	validURLs := []string{
		"http://foo.com/blah_blah",
		"http://userid:password@example.com:8080",
		"http://➡.ws/䨹",
		"http://例子.测试",
		"http://مثال.إختبار",
	}
	invalidURLs := []string{
		"http://??",
		"http://foo.bar?q=Spaces should be encoded",
		"//",
		"rdar://1234",
		"http://224.1.1.1",
	}
	for _, validURL := range validURLs {
		err := ValidateURL(ctx, validURL)
		if err != nil {
			t.Fatalf("%+v should be validURL", validURL)
		}
	}
	for _, invalidURL := range invalidURLs {
		err := ValidateURL(ctx, invalidURL)
		if err == nil {
			t.Fatalf("%+v should be validURL", invalidURL)
		}
	}
}

func TestValidateLabelKey(t *testing.T) {
	validLabelValues := []string{
		"a-b",
		"kubernetes_test",
		"test-test_test",
		"Aest---_test",
		"t-tBs_t",
	}
	invalidLabelValues := []string{
		"-a",
		"a-哈",
		"a-",
		"_test_",
		"!!@@-",
		"",
	}
	for _, validLabelValue := range validLabelValues {
		err := ValidateLabelKey(ctx, validLabelValue)
		if err != nil {
			t.Fatalf("%+v should be validLabelValue", validLabelValue)
		}
	}
	for _, invalidLabelValue := range invalidLabelValues {
		err := ValidateLabelKey(ctx, invalidLabelValue)
		if err == nil {
			t.Fatalf("%+v should be invalidLabelValue", invalidLabelValue)
		}
	}
}

func TestValidateLabelValue(t *testing.T) {
	validLabelValues := []string{"aa", "validatename!2313", "2131", "!!,."}
	invalidLabelValues := []string{"", "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"}
	for _, validLabelValue := range validLabelValues {
		err := ValidateLabelValue(ctx, validLabelValue)
		if err != nil {
			t.Fatalf("%+v should be validLabelValue", validLabelValue)
		}
	}
	for _, invalidLabelValue := range invalidLabelValues {
		err := ValidateLabelValue(ctx, invalidLabelValue)
		if err == nil {
			t.Fatalf("%+v should be invalidLabelValue", invalidLabelValue)
		}
	}
}

func TestValidateLabelMapFmt(t *testing.T) {
	validLabelMaps := []map[string][]string{
		{
			"a-b":             {"adsf"},
			"kubernetes_test": {"132432"},
		},
		{
			"test-test_test": {"validatename!2313"},
			"Aest---_test":   {"2131"},
			"t-tBs_t":        {"!!,."},
		},
	}
	invalidLabelMaps := []map[string][]string{
		{
			"_a-b":            {"adsf"},
			"kubernetes_test": {"132432"},
		},
		{
			"test-test_test": {"validatename!2313"},
			"Aest---_test":   {"2131"},
			"t-tBs_t-":       {"!!,."},
		},
		{
			"test-test_test": {"validatename!2313"},
			"Aest---_test":   {"2131"},
			"t-tBs_t":        {"!!,.", "12231"},
		},
		{
			"test-test_test": {"validatename!2313"},
			"Aest---_test":   {"2131"},
			"":               {"!!,."},
		},
		{
			"test-test_test": {"validatename!2313"},
			"Aest---_test":   {"2131"},
			"t-tBs_t-":       {""},
		},
	}
	for _, validLabelMap := range validLabelMaps {
		err := ValidateLabelMapFmt(ctx, validLabelMap)
		if err != nil {
			t.Fatalf("%+v should be validLabelMap", validLabelMap)
		}
	}
	for _, invalidLabelMap := range invalidLabelMaps {
		err := ValidateLabelMapFmt(ctx, invalidLabelMap)
		if err == nil {
			t.Fatalf("%+v should be invalidLabelMap", invalidLabelMap)
		}
	}
}

func TestValidateLabelString(t *testing.T) {
	validLabelStrings := []string{
		"runtime=qingcloud&zone=test&test-test_test=314134",
		"runtime=kubernetes&Aest---_test=3242",
	}
	invalidLabelStrings := []string{
		"a!=b",
		"runtime=kubernetes&__=1",
	}
	for _, validLabelString := range validLabelStrings {
		err := ValidateLabelString(ctx, validLabelString)
		if err != nil {
			t.Fatalf("%+v should be validLabelString", validLabelString)
		}
	}
	for _, invalidLabelString := range invalidLabelStrings {
		err := ValidateLabelString(ctx, invalidLabelString)
		if err == nil {
			t.Fatalf("%+v should be invalidLabelString", invalidLabelString)
		}
	}

}

func TestValidateSelectorMapFmt(t *testing.T) {
	validSelectorMaps := []map[string][]string{
		{
			"a-b":             {"adsf"},
			"kubernetes_test": {"132432"},
		},
		{
			"test-test_test": {"validatename!2313"},
			"Aest---_test":   {"2131"},
			"t-tBs_t":        {"!!,."},
		},
		{
			"test-test_test": {"validatename!2313"},
			"Aest---_test":   {"2131"},
			"t-tBs_t":        {"!!,.", "12231"},
		},
	}
	invalidSelectorMaps := []map[string][]string{
		{
			"_a-b":            {"adsf"},
			"kubernetes_test": {"132432"},
		},
		{
			"test-test_test": {"validatename!2313"},
			"Aest---_test":   {"2131"},
			"t-tBs_t-":       {"!!,."},
		},
		{
			"test-test_test": {"validatename!2313"},
			"Aest---_test":   {"2131"},
			"":               {"!!,."},
		},
		{
			"test-test_test": {"validatename!2313"},
			"Aest---_test":   {"2131"},
			"t-tBs_t-":       {""},
		},
	}
	for _, validSelectorMap := range validSelectorMaps {
		err := ValidateSelectorMapFmt(ctx, validSelectorMap)
		if err != nil {
			t.Fatalf("%+v should be validSelectorMap", validSelectorMap)
		}
	}
	for _, invalidSelectorMap := range invalidSelectorMaps {
		err := ValidateSelectorMapFmt(ctx, invalidSelectorMap)
		if err == nil {
			t.Fatalf("%+v should be invalidSelectorMap", invalidSelectorMap)
		}
	}
}

func TestValidateSelectorString(t *testing.T) {
	validSelectorStrings := []string{
		"runtime=qingcloud&zone=test&test-test_test=314134",
		"runtime=kubernetes&Aest---_test=3242",
		"runtime=qingcloud&zone=pkea&zone=111,",
	}
	invalidSelectorStrings := []string{
		"zo!ne=peka",
		"acc%=b",
		"runtime=qingcloud&zone=pkea&zone_=111,",
		"runtime=kubernetes&__=1",
	}

	for _, validSelectorString := range validSelectorStrings {
		err := ValidateSelectorString(ctx, validSelectorString)
		if err != nil {
			t.Fatalf("%+v should be validSelectorString", validSelectorString)
		}
	}
	for _, invalidSelectorString := range invalidSelectorStrings {
		err := ValidateSelectorString(ctx, invalidSelectorString)
		if err == nil {
			t.Fatalf("%+v should be invalidSelectorString", invalidSelectorString)
		}
	}
}
