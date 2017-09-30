# +-------------------------------------------------------------------------
# | Copyright (C) 2017 Yunify, Inc.
# +-------------------------------------------------------------------------
# | Licensed under the Apache License, Version 2.0 (the "License");
# | you may not use this work except in compliance with the License.
# | You may obtain a copy of the License in the LICENSE file, or at:
# |
# | http://www.apache.org/licenses/LICENSE-2.0
# |
# | Unless required by applicable law or agreed to in writing, software
# | distributed under the License is distributed on an "AS IS" BASIS,
# | WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# | See the License for the specific language governing permissions and
# | limitations under the License.
# +-------------------------------------------------------------------------

# https://github.com/go-swagger/go-swagger
# https://github.com/mohsen1/multi-file-swagger-example
# http://apihandyman.io/writing-openapi-swagger-specification-tutorial-part-1-introduction/

default:
	cat copyright.yaml > _all.yaml
	multi-file-swagger -o yaml index.yaml >> _all.yaml
	swagger validate _all.yaml

install-tools:
	go get -u github.com/go-swagger/go-swagger/cmd/swagger
	npm install -g multi-file-swagger

clean:
