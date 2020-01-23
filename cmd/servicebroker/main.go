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

// Package main initialize and start the broker.
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/pivotal-cf/brokerapi"
	"github.com/wso2/service-broker-apim/pkg/apim"
	"github.com/wso2/service-broker-apim/pkg/broker"
	"github.com/wso2/service-broker-apim/pkg/client"
	"github.com/wso2/service-broker-apim/pkg/config"
	"github.com/wso2/service-broker-apim/pkg/db"
	"github.com/wso2/service-broker-apim/pkg/log"
	"github.com/wso2/service-broker-apim/pkg/model"
	"github.com/wso2/service-broker-apim/pkg/token"
)

const (
	ErrMsgUnableToStartServerTLS = "unable to start the server on Host: %s port: %s TLS key: %s TLS cert: %s"
	ErrMsgUnableToStartServer    = "unable to start the server on Host: %s port: %s"
	InfoMSGShutdownBroker        = "starting APIM Service Broker shutdown"
	InfoMSGServerStart           = "starting APIM broker"
	ErrMsgUnableToAddForeignKeys = "unable to add foreign keys"
)

func main() {

	// load configuration.
	conf, err := config.Load()
	if err != nil {
		log.HandleErrorAndExit("failed to load configuration", err)
	}
	// configure logging.
	logger, err := log.Configure(conf.Log.FilePath, conf.Log.Level)
	if err != nil {
		log.HandleErrorAndExit("failed to configure logger", err)
	}
	// configure HTTP client
	client.Configure(&conf.HTTP.Client)

	// Initialize DB.
	db.Init(&conf.DB)
	defer db.CloseDBCon()
	setupTables()

	// Initialize Token manager.
	tManager := &token.PasswordRefreshTokenGrantManager{
		TokenEndpoint:                    conf.APIM.TokenEndpoint,
		DynamicClientEndpoint:            conf.APIM.DynamicClientEndpoint,
		DynamicClientRegistrationContext: conf.APIM.DynamicClientRegistrationContext,
		UserName:                         conf.APIM.Username,
		Password:                         conf.APIM.Password,
	}
	tManager.Init([]string{token.ScopeSubscribe, token.ScopeAPIView})

	// Initialize API-M client.
	apim.Init(tManager, conf.APIM)

	brokerCreds := brokerapi.BrokerCredentials{
		Username: conf.HTTP.Server.Auth.Username,
		Password: conf.HTTP.Server.Auth.Password,
	}
	apimServiceBroker := &broker.APIM{}
	apimServiceBroker.Init()
	brokerAPI := brokerapi.New(apimServiceBroker, logger, brokerCreds)

	host := conf.HTTP.Server.Host
	port := conf.HTTP.Server.Port
	ld := log.NewData().
		Add("host", host).
		Add("port", port)

	server := http.Server{
		Handler: brokerAPI,
		Addr:    host + ":" + port,
	}

	// Handling terminating signal.
	idleConsClosed := make(chan struct{}, 1)
	go handleGracefulShutdown(idleConsClosed, &server)

	log.Info(InfoMSGServerStart, ld)
	if !conf.HTTP.Server.TLS.Enabled {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.HandleErrorAndExit(
				fmt.Sprintf(ErrMsgUnableToStartServer, host, port), err)
		}
	} else {
		if err := server.ListenAndServeTLS(conf.HTTP.Server.TLS.Cert, conf.HTTP.Server.TLS.Key); err != http.ErrServerClosed {
			log.HandleErrorAndExit(fmt.Sprintf(ErrMsgUnableToStartServerTLS,
				host,
				port,
				conf.HTTP.Server.TLS.Key,
				conf.HTTP.Server.TLS.Cert),
				err)
		}
	}
	log.Debug("waiting for idle connections to be closed", nil)
	<-idleConsClosed
}

// addForeignKeys configures foreign keys for Subscription table.
func addForeignKeys() {
	// With this foreign key mapping all the subscriptions are deleted respective once the service instance is deleted.
	err := db.AddForeignKey(&model.Subscription{}, model.ServiceInstanceIDFieldName, model.ForeignKeyDestSVCInstanceID, "CASCADE",
		"CASCADE")
	if err != nil {
		log.HandleErrorAndExit(ErrMsgUnableToAddForeignKeys, err)
	}
	// With this foreign key mapping  it is restricted to delete a bind of a existing service instance.
	err = db.AddForeignKey(&model.Bind{}, model.ServiceInstanceIDFieldName, model.ForeignKeyDestSVCInstanceID, "RESTRICT",
		"RESTRICT")
	if err != nil {
		log.HandleErrorAndExit(ErrMsgUnableToAddForeignKeys, err)
	}
}

// SetupTables creates the tables and add foreign keys.
func setupTables() {
	db.CreateTable(&model.ServiceInstance{})
	db.CreateTable(&model.Subscription{})
	db.CreateTable(&model.Bind{})
	addForeignKeys()
}

// handleGracefulShutdown shutdown the server gracefully.
func handleGracefulShutdown(idleConsClosed chan<- struct{}, server *http.Server) {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, os.Kill)
	log.Debug("graceful shutdown process is started. Waiting for interrupt or kill signal", nil)
	<-sigint
	log.Debug("interrupt or kill signal received", nil)
	log.Info(InfoMSGShutdownBroker, nil)
	if err := server.Shutdown(context.Background()); err != nil {
		// Error from closing listeners, or context timeout.
		log.Error("unable to shutdown server", err, nil)
	}
	close(idleConsClosed)
}
