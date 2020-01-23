# ------------------------------------------------------------------------
#
# Copyright 2019 WSO2, Inc. (http://wso2.com)
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License
#
# ------------------------------------------------------------------------

all: deps tests build

build: ## Builds the service broker
	go build -i github.com/wso2/service-broker-apim/cmd/servicebroker

tests: ## Runs the tests
	go test -v ./pkg/...

integration-test-start:
	./test/run-tests.sh

integration-test-stop:
	docker-compose -f ./test/integration-test-setup.yaml down

debug-setup-up:
	docker-compose -f ./test/debug-setup.yaml up -d

debug-setup-down:
	docker-compose -f ./test/debug-setup.yaml down

build-linux: ## Builds a Linux executable
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
	go build -o ./target/linux/servicebroker github.com/wso2/service-broker-apim/cmd/servicebroker

build-darwin: ## Builds a Darwin executable
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 \
	go build -o ./target/darwin/servicebroker github.com/wso2/service-broker-apim/cmd/servicebroker

clean: ## Cleans up build artifacts
	rm -f servicebroker
	rm -fr target

deps:
	dep ensure

setup-lint: ## Install golint
	go get -u golang.org/x/lint/golint

lint: ## Run golint on the code
	golint  ./pkg/* ./cmd/*

format: ## Run gofmt on the code
	gofmt -w ./pkg/* ./cmd/*

help: ## Shows the help
	@echo 'Usage: make <OPTIONS> ... <TARGETS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
        awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ''

