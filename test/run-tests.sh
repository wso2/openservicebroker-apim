#!/bin/sh
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
set -e

docker-compose -f $PWD/test/integration-test-setup.yaml up --build -d

until curl http://admin:admin@localhost:8444/v2/catalog -H "X-Broker-API-Version: 2.14" --silent --output /dev/null ; do
  >&2 echo "Broker is unavailable - sleeping"
  sleep 5
done

docker run -v $PWD/test/collections:/etc/newman --network="host" -t postman/newman:ubuntu \
    run "OSB-Integration-tests.postman_collection.json" \
    --environment="OSB-Integration-test.postman_environment.json" \
    --reporters="json,cli" --reporter-json-export="newman-results.json"

docker-compose -f $PWD/test/integration-test-setup.yaml down