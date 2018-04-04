// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package drone

import (
	"fmt"
	"net"
	"strings"

	"openpitrix.io/libconfd"
	pbdrone "openpitrix.io/openpitrix/pkg/pb/drone"
)

func MakeDroneId(suffix string) string {
	return fmt.Sprintf("drone@%s/%s", getLocalIP(), strings.TrimSpace(suffix))
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}

func To_pbdrone_ConfdConfig(cfg *libconfd.Config) *pbdrone.ConfdConfig {
	return &pbdrone.ConfdConfig{
		ConfDir:  cfg.ConfDir,
		Interval: int32(cfg.Interval),
		Prefix:   cfg.Prefix,
		SyncOnly: cfg.SyncOnly,
		LogLevel: cfg.LogLevel,
	}
}

func To_libconfd_Config(cfg *pbdrone.ConfdConfig) *libconfd.Config {
	return &libconfd.Config{
		ConfDir:  cfg.ConfDir,
		Interval: int(cfg.Interval),
		Prefix:   cfg.Prefix,
		SyncOnly: cfg.SyncOnly,
		LogLevel: cfg.LogLevel,
	}
}

func To_pbdrone_ConfdBackendConfig(bcfg *libconfd.BackendConfig) *pbdrone.ConfdBackendConfig {
	return &pbdrone.ConfdBackendConfig{
		Type:         bcfg.Type,
		Host:         append([]string{}, bcfg.Host...),
		Username:     bcfg.UserName,
		Password:     bcfg.Password,
		ClientCaKeys: bcfg.ClientCAKeys,
		ClientCert:   bcfg.ClientCert,
		ClientKey:    bcfg.ClientKey,
	}
}

func To_libconfd_BackendConfig(bcfg *pbdrone.ConfdBackendConfig) *libconfd.BackendConfig {
	return &libconfd.BackendConfig{
		Type:         bcfg.Type,
		Host:         append([]string{}, bcfg.Host...),
		UserName:     bcfg.Username,
		Password:     bcfg.Password,
		ClientCAKeys: bcfg.ClientCaKeys,
		ClientCert:   bcfg.ClientCert,
		ClientKey:    bcfg.ClientKey,
	}
}
