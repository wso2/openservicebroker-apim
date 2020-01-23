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

until nc -z mysql 3306; do
  >&2 echo "Mysql is unavailable - sleeping"
  sleep 1
done

until nc -z wso2apim 9443; do
  >&2 echo "WSO2 APIM is unavailable - sleeping"
  sleep 1
done


exec ./servicebroker-linux