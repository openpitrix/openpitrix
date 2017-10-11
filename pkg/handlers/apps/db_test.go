// +-------------------------------------------------------------------------
// | Copyright (C) 2017 Yunify, Inc.
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

// start mysql in docker for test

package apps

import (
	"flag"
	"os"
	"testing"
)

var (
	fFalgDbName = flag.String("test-db-name", "root:password@tcp(mysql:3306)/openpitrix", "set mysql database name")
)

func TestMain(m *testing.M) {
	flag.Parse()

	rv := m.Run()
	os.Exit(rv)
}

func TestAppDatabase(t *testing.T) {
	// TODO
}
