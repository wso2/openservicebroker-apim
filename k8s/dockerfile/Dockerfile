# ------------------------------------------------------------------------
#
# Copyright 2020 WSO2, Inc. (http://wso2.com)
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

COPY Gopkg.toml Gopkg.lock $GOPATH/src/github.com/wso2/openservicebroker-apim/
COPY cmd/ $GOPATH/src/github.com/wso2/openservicebroker-apim/cmd/
COPY pkg/ $GOPATH/src/github.com/wso2/openservicebroker-apim/pkg/
COPY vendor/ $GOPATH/src/github.com/wso2/openservicebroker-apim/vendor/

WORKDIR $GOPATH/src/github.com/wso2/openservicebroker-apim
RUN dep ensure

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /servicebroker-linux github.com/wso2/openservicebroker-apim/cmd/servicebroker

FROM amd64/alpine:3.10.2
LABEL maintainer="WSO2 Docker Maintainers <dev@wso2.org>"

# build arguments for user/group configurations
ARG USER=wso2carbon
ARG USER_ID=802
ARG USER_GROUP=wso2
ARG USER_GROUP_ID=802
ARG USER_HOME=/home/${USER}

# create the non-root user and group
RUN \
    addgroup -S -g ${USER_GROUP_ID} ${USER_GROUP} \
    && adduser -S -u ${USER_ID} -h ${USER_HOME} -G ${USER_GROUP} ${USER}

# copy the binary from builder image
COPY --from=builder /servicebroker-linux ${USER_HOME}/

# set permission for the binary
RUN \
    chmod +x ${USER_HOME}/servicebroker-linux \
    && chown wso2carbon:wso2 ${USER_HOME}/servicebroker-linux

# set the user and work directory
USER ${USER_ID}
WORKDIR ${USER_HOME}

ENTRYPOINT ["/home/wso2carbon/servicebroker-linux"]
