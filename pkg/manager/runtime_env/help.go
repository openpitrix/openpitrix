// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package runtime_env

import (
	"fmt"
	"strings"

	"openpitrix.io/openpitrix/pkg/models"
)

func LabelStringToMap(labelString string) (map[string]string, error) {
	mapLabel := make(map[string]string)
	sliceArray := strings.Split(labelString, ",")
	for _, label := range sliceArray {
		labelArray := strings.Split(label, "=")
		if len(labelArray) != 2 {
			return nil, fmt.Errorf("bad label format %v", labelString)
		}
		labelArray[0] = strings.TrimSpace(labelArray[0])
		labelArray[1] = strings.TrimSpace(labelArray[1])
		if mapLabel[labelArray[0]] != "" {
			return nil, fmt.Errorf("bad label format %v", labelString)
		}
		mapLabel[labelArray[0]] = labelArray[1]
	}

	return mapLabel, nil
}

func LabelMapToString(labelMap map[string]string) (labelString string) {
	for labelKey, labelValue := range labelMap {
		labelString += fmt.Sprintf("%v=%v,", labelKey, labelValue)
	}
	labelString = strings.Trim(labelString, ",")
	return labelString
}

func SelectorStringToMap(selectorString string) (map[string][]string, error) {
	selectorMap := make(map[string][]string)
	sliceArray := strings.Split(selectorString, ",")
	for _, label := range sliceArray {
		labelArray := strings.Split(label, "=")
		if len(labelArray) != 2 {
			return nil, fmt.Errorf("bad selector format %v", selectorMap)
		}
		labelArray[0] = strings.TrimSpace(labelArray[0])
		labelArray[1] = strings.TrimSpace(labelArray[1])
		selectorMap[labelArray[0]] = append(selectorMap[labelArray[0]], labelArray[1])
	}
	return selectorMap, nil
}

func LabelMapDiff(oldLabelMap, newLabelMap map[string]string) (additions, deletions map[string]string) {
	additions = make(map[string]string)
	deletions = make(map[string]string)
	for i := 0; i < 2; i++ {
		for oldLabelKey, oldLabelValue := range oldLabelMap {
			found := false
			for newLabelKey, newLabelValue := range newLabelMap {
				if oldLabelKey == newLabelKey && oldLabelValue == newLabelValue {
					found = true
				}
			}
			if !found {
				if i == 0 {
					deletions[oldLabelKey] = oldLabelValue
				} else {
					additions[oldLabelKey] = oldLabelValue
				}
			}
		}
		if i == 0 {
			oldLabelMap, newLabelMap = newLabelMap, oldLabelMap
		}
	}
	return additions, deletions
}

func LabelStructToMap(labelStructs []*models.RuntimeEnvLabel) map[string]string {
	mapLabel := make(map[string]string)
	for _, labelStruct := range labelStructs {
		mapLabel[labelStruct.LabelKey] = labelStruct.LabelValue
	}
	return mapLabel
}
