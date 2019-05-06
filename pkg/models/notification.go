// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

type EmailNotification struct {
	Title       string
	Content     string
	Owner       string
	ContentType string
	Addresses   []string
}

func UniqueAddresses(addresses []string) []string {
	addressMap := make(map[string]string)
	var uniqueAddresses []string

	for _, address := range addresses {
		addressMap[address] = address
	}

	for _, address := range addressMap {
		uniqueAddresses = append(uniqueAddresses, address)
	}
	return uniqueAddresses
}

func UniqueEmailNotifications(emailNotifications []*EmailNotification) []*EmailNotification {
	emailNotificationMap := make(map[string]*EmailNotification)
	var uniqueEmailNotifications []*EmailNotification

	for _, emailNotification := range emailNotifications {
		key := emailNotification.Title + emailNotification.Content
		_, isExist := emailNotificationMap[key]
		if isExist {
			emailNotificationMap[key].Addresses = append(emailNotificationMap[key].Addresses, emailNotification.Addresses...)
		} else {
			emailNotificationMap[key] = emailNotification
		}
	}

	for _, emailNotification := range emailNotificationMap {
		emailNotification.Addresses = UniqueAddresses(emailNotification.Addresses)
		uniqueEmailNotifications = append(uniqueEmailNotifications, emailNotification)
	}
	return uniqueEmailNotifications
}
