/*
 * Copyright (c) 2019 WSO2 Inc. (http:www.wso2.org) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http:www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
)

const (
	ConfigFilePath            = "../../config/config.yaml"
	ErrMsgTestIncorrectResult = "expected value: %v but then returned value: %v"
)

func setUpEnv(key, val string, t *testing.T) {
	err := os.Setenv(key, val)
	if err != nil {
		t.Error(err)
	}
}

func tearDownEnv(key string, t *testing.T) {
	err := os.Unsetenv(FilePathEnv)
	if err != nil {
		t.Error(err)
	}
}

func TestLoadConfigFile(t *testing.T) {
	setUpEnv(FilePathEnv, ConfigFilePath, t)
	err := loadConfigFile()
	if err != nil {
		t.Error(err)
	}
	if viper.ConfigFileUsed() != ConfigFilePath {
		t.Errorf(ErrMsgTestIncorrectResult, ConfigFilePath, viper.ConfigFileUsed())
	}
	viper.Reset()
	tearDownEnv(FilePathEnv, t)
}

func TestSetDefaultConf(t *testing.T) {
	setDefaultConf()
	testStringConf(t, "log.filePath", "server.log")
	testStringConf(t, "log.level", "info")
	testStringConf(t, "http.server.auth.username", "admin")
	testStringConf(t, "http.server.auth.password", "admin")
	testBooleanConf(t, "http.server.tls.enabled", false)
	testStringConf(t, "http.server.tls.key", "key.pem")
	testStringConf(t, "http.server.tls.cert", "cert.pem")
	testStringConf(t, "http.server.host", "0.0.0.0")
	testStringConf(t, "http.server.port", "8444")
	testIntegerConf(t, "http.client.timeout", 30)
	testIntegerConf(t, "http.client.minBackOff", 1)
	testIntegerConf(t, "http.client.maxBackOff", 60)
	testIntegerConf(t, "http.client.maxRetries", 3)
	testBooleanConf(t, "http.client.insecureCon", true)
	testStringConf(t, "apim.username", "admin")
	testStringConf(t, "apim.password", "admin")
	testStringConf(t, "apim.tokenEndpoint", "https://localhost:8243")
	testStringConf(t, "apim.dynamicClientEndpoint", "https://localhost:9443")
	testStringConf(t, "apim.dynamicClientRegistrationContext", "/client-registration/v0.14/register")
	testStringConf(t, "apim.publisherEndpoint", "https://localhost:9443")
	testStringConf(t, "apim.publisherAPIContext", "/api/am/publisher/v0.14/apis")
	testStringConf(t, "apim.publisherChangeAPILifeCycleContext", "/api/am/publisher/v0.14/apis/change-lifecycle")
	testStringConf(t, "apim.storeEndpoint", "https://localhost:9443")
	testStringConf(t, "apim.storeApplicationContext", "/api/am/store/v0.14/applications")
	testStringConf(t, "apim.storeSubscriptionContext", "/api/am/store/v0.14/subscriptions")
	testStringConf(t, "apim.generateApplicationKeyContext", "/api/am/store/v0.14/applications/generate-keys")
	testStringConf(t, "db.host", "localhost")
	testIntegerConf(t, "db.port", 3306)
	testStringConf(t, "db.username", "root")
	testStringConf(t, "db.password", "root123")
	testStringConf(t, "db.database", "broker")
	testBooleanConf(t, "db.logMode", false)
	testIntegerConf(t, "db.maxRetries", 3)
	viper.Reset()
}

func testBooleanConf(t *testing.T, key string, expected bool) {
	result := viper.GetBool(key)
	if expected != result {
		t.Errorf(ErrMsgTestIncorrectResult, expected, result)
	}
}

func testIntegerConf(t *testing.T, key string, expected int) {
	result := viper.GetInt(key)
	if expected != result {
		t.Errorf(ErrMsgTestIncorrectResult, expected, result)
	}
}

func testStringConf(t *testing.T, key, expected string) {
	result := viper.GetString(key)
	if expected != result {
		t.Errorf(ErrMsgTestIncorrectResult, expected, result)
	}
}

func TestEnvConf(t *testing.T) {
	setUpEnv(EnvPrefix+"_DB_DATABASE", "broker_test", t)
	c, err := Load()
	if err != nil {
		t.Error(err)
	}
	if c.DB.Database != "broker_test" {
		t.Errorf(ErrMsgTestIncorrectResult, "broker_test", c.DB.Database)
	}
	tearDownEnv(EnvPrefix+"_db_database", t)
	viper.Reset()
}
