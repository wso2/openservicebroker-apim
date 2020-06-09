# Service Broker for API Manager

This repository contains Service Broker for WSO2 API Manager, implemented using the “Go”programming language. 

Let's take you through the steps of installation, and execution of WSO2 API-M Service Broker.

## Prerequisites

Before you begin, be sure you have met the installation prerequisites, and then follow the installation instructions.

* [Go tools](https://golang.org/doc/install)
* A running [WSO2 API Manager](https://docs.wso2.com/display/AM260/Installation+Guide)
* Database to preserve data on broker services and binds

## Quick Start Guide

1. Install openservicebroker_apim project to your ```$GOPATH```
```
$ go get github.com/wso2/openservicebroker_apim
```
2. Navigate to the ```openservicebroker_apim``` project directory in ```$GOPATH/src/github.com/wso2```.
```
$ cd $GOPATH/src/github.com/wso2/openservicebroker_apim
```
3. Update ```config/config.yaml``` in the project directory by replacing the following and save the yaml file. 


| Reference in the config.yaml  | Description                                               | Example           |
| ------------------------------|-----------------------------------------------------------|-------------------|
| ${wso2apim-endpoint}          | WSO2 API Manager endpoint.                                | localhost:9443    |
| ${keymanager-endpoint}        | Keymanager endpoint of WSO2 APIM                          | localhost:9443    |
| ${gateway-endpoint}           | Gateway endpoint of WSO2 APIM                             | localhost:8243    |
| ${broker-db-hostname}         | Database hostname                                         | localhost         |
| ${broker-db-port}             | Connecting port of your database                          | 3306              |
| ${broker-db-username}         | Username for authentication to databases                  | root              |
| ${broker-db-password}         | Password for authentication to databases                  | root              |
| ${broker-db-name}             | Name of the database that you want to store your data in  | broker            |

4. Set the environment variable for the ```config.yaml``` using the follwing command. 
```
$ export APIM_BROKER_CONF_FILE=$GOPATH/src/github.com/wso2/openservicebroker_apim/config/config.yaml
```
5. Create Service Broker binary executable, using the repective command below.
* Ubuntu OS users
```
$ make build-linux
```
* MacOs Users
```
$ make build-darwin
```

6. Follow the instructions given below and execute the binary to expose ```APIM Service Broker``` endpoints.
* Ubuntu OS users
```
$ cd target/linux
$ ./servicebroker
```
* MacOs Users
```
$ cd target/darwin
$ ./servicebroker
```