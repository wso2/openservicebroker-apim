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
FROM golang:1.12.1 AS builder

# Download and install the latest release of dep
ADD https://github.com/golang/dep/releases/download/v0.5.1/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep

COPY Gopkg.toml Gopkg.lock $GOPATH/src/github.com/wso2/service-broker-apim/
COPY cmd/ $GOPATH/src/github.com/wso2/service-broker-apim/cmd/
COPY pkg/ $GOPATH/src/github.com/wso2/service-broker-apim/pkg/
COPY vendor/ $GOPATH/src/github.com/wso2/service-broker-apim/vendor/

WORKDIR $GOPATH/src/github.com/wso2/service-broker-apim
RUN dep ensure

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /servicebroker-linux github.com/wso2/service-broker-apim/cmd/servicebroker

FROM amd64/alpine:3.10.2
COPY --from=builder /servicebroker-linux ./
COPY ./test/init-scripts/wait-for-deps.sh ./wait-for-deps.sh
RUN chmod +x ./wait-for-deps.sh
ENTRYPOINT ["./wait-for-deps.sh"]