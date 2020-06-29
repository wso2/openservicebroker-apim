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

## Working with APIM Service Broker for Pivotal Cloud Foundry (PCF)
WSO2 API Manager Service for PCF  furnish a set of subscribed APIs  to be consumed in a user application. 

#### Checking availability of APIM Service broker. 

To use the APIM service for PCF, the service has to be readily available in the marketplace. User can assess the availability by runing the following command.
```
$ cf marketplace
```
Note: If the output shows **wso2apim-service** under *service* column, it indicates that the WSO2 APIM Service is available in your space.

#### Check for any instances running in the current active space

To determine any active running instances on the Pivotal Space , Follow the steps below:

1. Log in to the org and space in PCF that contains your application. 

2. Run the following command in your terminal
```
$ cf services
```

3. Any **wso2apim-service** listing in the *service* column is a service instance of APIM Service. Now you can bind your app to an existing service instance or create a new instance for binding. 

#### Creating a Service Instance

To create an instance for APIM service, follow the below steps:

1. Run the following command

```
$ cf create-service [SERVICE] [PLAN] [SERVICE_INSTANCE] -c [PARAM_JSON]
```

2. Wait and monitor the status of your instance with the following command. 

```
$ watch cf services
```
Note: ensure all the  *last operation* column status indicates **create succeeded**.


#### Updating a service instance

Update the service instance by adding or deleting APIs to your service. Run the following command to update the service instance.
```
$ cf update-service [SERVICE_INSTANCE] -c [PARAM_JSON]
```

#### Using APIM Service in your application

1. Bind a service instance to your app

For your application to use APIM service, you have to bind the app with the created service instance. Run the following command
```
$ cf bind-service [APP] [SERVICE INSTANCE]
```

2. Unbind an app from the service instance 

To stop the app from using APIM Service, unbind from the service using the following command. 
```
$ cf unbind-service [APP] [SERVIECE_INSTANCE]
```
3. Deleting a Service Instance

To delete any previously created service instance run the following command:
```
$ cf delete-service [SERVICE_INSTANCE]
```

Note:  First confirm to unbind all the apps that were consuming the APIM service before deleting the service instance. Deletion of services is not feasible if there are application still consuming services.

```[APP]: Name of your application```

```[SERVICE_INSTANCE]: Service instance name```

```[SERVICE]: The name of the APIM service you want to create an instance. of (i.e:  wso2apim-service)```

```[PLAN]: The name of the APIM service plan that meets your need. (i.e: app )```

```[SERVICE_INSTANCE]: Provide any service instance name```

```[PARAM_JSON]: a valid JSON object containing service-specific configurations.``` 