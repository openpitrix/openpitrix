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

TARG:=_all.json

SWAGGER:=docker run --rm -it -v $(HOME):$(HOME) -w `pwd` quay.io/goswagger/swagger
MULTI_FILE_SWAGGER:=docker run --rm -it -v $(HOME):$(HOME) -w `pwd` vidsyhq/multi-file-swagger-docker

help:
	@echo "Please use \`make <target>\` where <target> is one of"
	@echo "  generate          to generate and check '$(TARG)'"
	@echo "  tools             to install depends tools"
	@echo "  clean             to clean the generated '$(TARG)' file"
	@echo "  help              to show this help informatrion"

generate:
	$(MULTI_FILE_SWAGGER) index.yaml > $(TARG)
	$(SWAGGER) validate $(TARG)
	@echo "ok"

tools:
	docker pull quay.io/goswagger/swagger
	docker pull vidsyhq/multi-file-swagger-docker
	@echo "ok"

clean:
	-rm $(TARG)
	@echo "ok"
