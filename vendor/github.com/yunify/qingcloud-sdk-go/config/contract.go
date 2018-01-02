// +-------------------------------------------------------------------------
// | Copyright (C) 2016 Yunify, Inc.
// +-------------------------------------------------------------------------
// | Licensed under the Apache License, Version 2.0 (the "License");
// | you may not use this work except in compliance with the License.
// | You may obtain a copy of the License in the LICENSE file, or at:
// |
// | http://www.apache.org/licenses/LICENSE-2.0
// |
// | Unless required by applicable law or agreed to in writing, software
// | distributed under the License is distributed on an "AS IS" BASIS,
// | WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// | See the License for the specific language governing permissions and
// | limitations under the License.
// +-------------------------------------------------------------------------

package config

import (
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strings"
)

// DefaultConfigFileContent is the default config file content.
const DefaultConfigFileContent = `# QingCloud services configuration

#qy_access_key_id: 'ACCESS_KEY_ID'
#qy_secret_access_key: 'SECRET_ACCESS_KEY'

host: 'api.qingcloud.com'
port: 443
protocol: 'https'
uri: '/iaas'
connection_retries: 3
connection_timeout: 30

# Valid log levels are "debug", "info", "warn", "error", and "fatal".
log_level: 'warn'

`

// DefaultConfigFile is the default config file.
const DefaultConfigFile = "~/.qingcloud/config.yaml"

// GetUserConfigFilePath returns the user config file path.
func GetUserConfigFilePath() string {
	return strings.Replace(DefaultConfigFile, "~/", getHome()+"/", 1)
}

// InstallDefaultUserConfig install the default user config file.
func InstallDefaultUserConfig() error {
	err := os.MkdirAll(path.Dir(GetUserConfigFilePath()), 0755)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(GetUserConfigFilePath(), []byte(DefaultConfigFileContent), 0644)
}

func getHome() string {
	home := os.Getenv("HOME")
	if runtime.GOOS == "windows" {
		home = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
	}

	return home
}
