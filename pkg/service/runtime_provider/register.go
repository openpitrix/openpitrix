// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package runtime_provider

import (
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/plugins"
)

func RegisterRuntimeProvider(provider, config string) error {
	err := pi.Global().GlobalConfig().RegisterRuntimeProviderConfig(provider, config)
	if err != nil {
		return err
	}
	pi.Global().SetGlobalCfg(pi.Global().GlobalConfig())

	logger.Debug(nil, "Available plugins: %+v", plugins.GetAvailablePlugins())

	return nil
}
