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

// Package config is responsible for loading, parsing the configuration.
package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const (
	// EnvPrefix is used as the prefix for Environment variable configuration.
	EnvPrefix = "APIM_BROKER"
	// FilePathEnv is used as a key to get the configuration file location.
	FilePathEnv = EnvPrefix + "_CONF_FILE"
	// FileType constant is used to specify the Configuration file type(YAML).
	FileType = "yaml"

	InfoMsgSettingUp        = "loading the configuration file: %s "
	ErrMsgUnableToReadConf  = "unable to read configuration: %s"
	ErrMsgUnableToParseConf = "unable to parse configuration"
)

// DB represent the Database configuration.
type DB struct {
	Host       string `mapstructure:"host"`
	Port       int    `mapstructure:"port"`
	Username   string `mapstructure:"username"`
	Password   string `mapstructure:"password"`
	Database   string `mapstructure:"database"`
	LogMode    bool   `mapstructure:"logMode"`
	MaxRetries int    `mapstructure:"maxRetries"`
}

// APIM represents the information required to interact with the APIM.
type APIM struct {
	Username                         string `mapstructure:"username"`
	Password                         string `mapstructure:"password"`
	TokenEndpoint                    string `mapstructure:"tokenEndpoint"`
	DynamicClientEndpoint            string `mapstructure:"dynamicClientEndpoint"`
	DynamicClientRegistrationContext string `mapstructure:"dynamicClientRegistrationContext"`
	PublisherEndpoint                string `mapstructure:"publisherEndpoint"`
	PublisherAPIContext              string `mapstructure:"publisherAPIContext"`
	StoreApplicationContext          string `mapstructure:"storeApplicationContext"`
	StoreSubscriptionContext         string `mapstructure:"storeSubscriptionContext"`
	StoreMultipleSubscriptionContext string `mapstructure:"storeMultipleSubscriptionContext"`
	StoreEndpoint                    string `mapstructure:"storeEndpoint"`
	GenerateApplicationKeyContext    string `mapstructure:"generateApplicationKeyContext"`
}

// Auth represents the username and the password for basic auth.
type Auth struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// TLS represents configuration needed for HTTPS.
type TLS struct {
	Enabled bool   `mapstructure:"enabled"`
	Key     string `mapstructure:"key"`
	Cert    string `mapstructure:"cert"`
}

// Log represents the configuration related to logging.
type Log struct {
	FilePath string `mapstructure:"filePath"`
	Level    string `mapstructure:"level"`
}

// Server represents configuration needed for the HTTP server.
type Server struct {
	Auth Auth   `mapstructure:"auth"`
	TLS  TLS    `mapstructure:"tls"`
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

// Client represents configuration needed for the HTTP client.
type Client struct {
	InsecureCon bool `mapstructure:"insecureCon"`
	Timeout     int  `mapstructure:"timeout"`
	MinBackOff  int  `mapstructure:"minBackOff"`
	MaxBackOff  int  `mapstructure:"maxBackOff"`
	MaxRetries  int  `mapstructure:"maxRetries"`
}

// HTTP represents configuration needed for the HTTP server and client.
type HTTP struct {
	Server Server `mapstructure:"server"`
	Client Client `mapstructure:"client"`
}

// Broker main struct which holds  sub configurations.
type Broker struct {
	Log  Log  `mapstructure:"log"`
	HTTP HTTP `mapstructure:"http"`
	APIM APIM `mapstructure:"apim"`
	DB   DB   `mapstructure:"db"`
}

// Load loads configuration into Broker object.
// Returns a pointer to the created Broker object or any error encountered.
func Load() (*Broker, error) {
	viper.SetConfigType(FileType)
	viper.SetEnvPrefix(EnvPrefix)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	setDefaultConf()
	if err := loadConfigFile(); err != nil {
		return nil, err
	}

	var brokerConfig Broker
	err := viper.Unmarshal(&brokerConfig)
	if err != nil {
		return nil, errors.Wrapf(err, ErrMsgUnableToParseConf)
	}
	return &brokerConfig, nil
}

// loadConfigFile loads the configuration into the Viper file only if the configuration file is pointed with "APIM_BROKER_CONF_FILE" environment variable.
// Returns an error if it is unable to read the config into Viper.
func loadConfigFile() error {
	confFile, exists := os.LookupEnv(FilePathEnv)
	if exists {
		fmt.Println(fmt.Sprintf(InfoMsgSettingUp, confFile))
		viper.SetConfigFile(confFile)
		if err := viper.ReadInConfig(); err != nil {
			return errors.Wrapf(err, ErrMsgUnableToReadConf, confFile)
		}
	}
	return nil
}

// setDefaultConf sets the default configurations for Viper.
func setDefaultConf() {
	viper.SetDefault("log.filePath", "server.log")
	viper.SetDefault("log.level", "info")

	viper.SetDefault("http.server.auth.username", "admin")
	viper.SetDefault("http.server.auth.password", "admin")
	viper.SetDefault("http.server.tls.enabled", false)
	viper.SetDefault("http.server.tls.key", "key.pem")
	viper.SetDefault("http.server.tls.cert", "cert.pem")
	viper.SetDefault("http.server.host", "0.0.0.0")
	viper.SetDefault("http.server.port", "8444")

	viper.SetDefault("http.client.insecureCon", true)
	viper.SetDefault("http.client.minBackOff", 1)
	viper.SetDefault("http.client.maxBackOff", 60)
	viper.SetDefault("http.client.timeout", 30)
	viper.SetDefault("http.client.maxRetries", 3)

	viper.SetDefault("apim.username", "admin")
	viper.SetDefault("apim.password", "admin")
	viper.SetDefault("apim.tokenEndpoint", "https://localhost:8243")
	viper.SetDefault("apim.dynamicClientEndpoint", "https://localhost:9443")
	viper.SetDefault("apim.dynamicClientRegistrationContext", "/client-registration/v0.14/register")
	viper.SetDefault("apim.publisherEndpoint", "https://localhost:9443")
	viper.SetDefault("apim.publisherAPIContext", "/api/am/publisher/v0.14/apis")
	viper.SetDefault("apim.storeEndpoint", "https://localhost:9443")
	viper.SetDefault("apim.storeApplicationContext", "/api/am/store/v0.14/applications")
	viper.SetDefault("apim.storeSubscriptionContext", "/api/am/store/v0.14/subscriptions")
	viper.SetDefault("apim.publisherChangeAPILifeCycleContext", "/api/am/publisher/v0.14/apis/change-lifecycle")
	viper.SetDefault("apim.storeMultipleSubscriptionContext", "/api/am/store/v0.14/subscriptions/multiple")
	viper.SetDefault("apim.generateApplicationKeyContext", "/api/am/store/v0.14/applications/generate-keys")

	viper.SetDefault("db.host", "localhost")
	viper.SetDefault("db.port", "3306")
	viper.SetDefault("db.username", "root")
	viper.SetDefault("db.password", "root123")
	viper.SetDefault("db.database", "broker")
	viper.SetDefault("db.logMode", false)
	viper.SetDefault("db.maxRetries", 3)
}
