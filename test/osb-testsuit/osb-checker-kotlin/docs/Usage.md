# Table of Contents
- [Description](../README.md#description)
- [Getting Started](../README.md#getting-started)
    - [Build Application](../README.md#build-application)
    - [Basic Run Configuration](../README.md#basic-run-configuration)
- [Changes](../README.md#changes)
- [Usage](#usage)
    - [Declaring Test Runs](#declaring-test-runs)
    - [Configuration](#configuration)
    - [Parameters](#parameters)
    - [Originating Identity](#originating-identity)
    - [Declaring Services](#declaring-services)
- [Test](#test)
    - [Catalog](#catalog)
    - [Provision](#provision)
        - [Test Procedure](ProvisionTests.md#test-procedure)
        - [Version specific Tests](ProvisionTests.md#version-specific-tests)
        - [Example Output](ProvisionTests.md#example-output)
    - [Binding](#binding)
    - [Authentication](#authentication)
    - [Contract](#contract)
    - [Example output](#example-output)

# Usage

osb-checker-kotlin provides a number of options to test service brokers with more detail, or shorten the run time, by leaving some tests out.

## Declaring Test Runs

If you want to run all tests just `java -jar osb-checker-kotlin-1.0.jar`

There are five different options to run tests. Possibles commands are:

* catalog: -cat/-catalog
* provision: -prov/-provision
* binding: -bind/-binding
* authentication: -auth/-authentication
* contract: -con/-contract

In case you want to run all tests call for example `java -jar osb-checker-kotlin-1.0.jar -cat -provision -bind -auth -con`
or just `java -jar osb-checker-kotlin-1.0.jar`

## Configuration

To run the application put a file with the name application.yml into the same location as the osb-checker-kotlin-1.0.jar file. For more information on how to configurate a spring boot application see [here](https://docs.spring.io/spring-boot/docs/current/reference/html/boot-features-external-config.html). 
 The .yml file needs the following schema.

```yaml

##Define the service broker connection here
config:
  url: http://localhost
  port : 80
  apiVersion: 2.15
  user: user
  password: password
##The following configuration are Optional
  skipTLSVerification: false
  usingAppGuid: true
  useRequestIdentity: true
  testDashboard: true
  
  originatingIdentity:
      platform: kubernetes
      value:
        username: duke,
        groups:
          - admin
          - dev
        extra:
          mydata:
            - data1
            - data3

  provisionParameters:
     plan-id-1-here:
        parameter1 : 1
        parameter2 : foo
      plan-id-2-here:
        parameter1 : 2
        parameter2 : bar

  bindingParameters:
      plan-id-here:
        key : value
        
  services:
    - id: service-id-here
      plans:
        - id: plan-id-here
```

**url**, **port**, **apiVersion**, **user** and **password** are mandatory and MUST be set.
Currently the application can test 2.13, 2.14 or 2.15 Service Brokers. Therefor **apiVersion** MUST be set to 2.13, 2.14 or 2.15.
**usingAppGuid**, **skipTLSVerification**, **useRequestIdentity**, **originatingIdentity**, **parameters**, **testDashboard** and **services** are optional.

Tests are created based upon the provided **apiVersion**. So the checker will not tests 2.15 functionality when this field is set to 2.13 or 2.14. More on version specific
testing here.

When **skipTLSVerification** is set request are 'http' is used instead of 'https' for every request. 

**usingAppGuid** sets the osb-checker to set a appGuid during provisioning. If no value it set it falls back to default true.

If **useRequestIdentity** is set to true, the osb-checker will set `X-Broker-API-Request-Identity` Header, for each request and verify if the header is present in the response.

**testDashboard** advises the checker to verify if an provided DashboardURL works after creating a service instance.

## Parameters

To set parameters for the provision, define them in parameters (Default is null).
specify the plan id as key for the parameters
example: a configuration with ...

```yaml
provisionParameters:
    plan-id-here:
      DB-name: db-name
      parameter1 : 1
      parameter2 : foo
      key : value
      schemaName: a_name
```

would run a provisions with the following request body.
```json
{
  "service_id": "service-id-here",
  "plan_id": "plan-id-here",
  "organization_guid": "org-guid-here",
  "space_guid": "space-guid-here",
  "parameters": {
    "DB-name": "db-name",
    "parameter1" : 1,
    "parameter2" : "foo",
    "key" : "value",
    "schemaName": "a_name"
    
    }
}
```

to declare parameters for a binding set them like this:

```yaml
bindingParameters:
    plan-id-here:
      key : value
      schemaName: a_name
```

### Originating Identity

If you wish to check the brokers behaviour with the header **X-Broker-API-Request-Identity** you define the value like in the following examples.
Set the unencoded value content in the value field.

```yaml
  originatingIdentity:
    platform: kubernetes
    value:
      username: duke,
      groups:
        - admin
        - dev
      extra:
        mydata:
          - data1
          - data3
```

```yaml
  originatingIdentity:
    platform: cloudfoundry
    value:
      user_id: myid
```

The declared content will be encoded and used in every request by the checker, according to [spec](https://github.com/openservicebrokerapi/servicebroker/blob/v2.15/spec.md#originating-identity)

## Declaring Services

The checker runs all it's tests on every service and plan defined in the service brokers catalog. If this is not desired, the developer can provide
a list of service and plan ids to direct the checker to the plans he would like to test. All details about the plans will be fetched from the actual catalog. It is only necessary to provide the ids.

```yaml
     services:
       - id: service-id-here
         plans:
           - id: plan-id-here
           - id: plan-id2-here
       - id: service-id2-here
         plans:
          - id: plan-id3-here
```
This config would run tests only for the defined plans.

```yaml
     services:
       - id: service-id-here
```
This configuration would run all plans that are defined within the service with the id 'service-id-here'.
