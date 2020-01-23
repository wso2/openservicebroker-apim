# osb-checker-kotlin

## Table of Contents
- [Description](#description)
- [Getting Started](#getting-started)
    - [Build Application](#build-application)
    - [Basic Run Configuration](#basic-run-configuration)
- [Changes](#changes)
- [Usage](docs/Usage.md)
    - [Declaring Test Runs](docs/Usage.md##declaring-test-runs)
    - [Parameters](docs/Usage.md#parameters)
    - [Originating Identity](docs/Usage.md#originating-identity)
    - [Declaring Services](docs/Usage.md#declaring-services)
- Test Classes
    - [Catalog](docs/CatalogTest.md)
    - [Provision](docs/ProvisionTests.md)
        - [Test Procedure](docs/ProvisionTests.md#test-procedure)
        - [Version specific Tests](docs/ProvisionTests.md#version-specific-tests)
        - [Example Output](docs/ProvisionTests.md#example-output)
    - [Binding](docs/BindingTests.md#binding)
        - [Test Procedure](docs/BindingTests.md#test-procedure)
        - [Version specific Tests](docs/BindingTests.md#version-specific-tests)
        - [Example Output](docs/BindingTests.md#example-output)
    - [Authentication](docs/AuthenticationTests.md)   
    - [Contract](docs/ContractTest.md)
- [Contribution](docs/Contribution.md)
   
## Description
This application is a generalized test program for service brokers. It runs rest calls against the defined service broker and checks if it
behaves as expected to the [service broker API specification](link=https://github.com/openservicebrokerapi/servicebroker)
Tests are created dynamically based upon the service broker catalog or custom input by the operator.

## Getting Started

### Build Application

to build the application run `{path}/osb-checker-kotlin/gradlew build` on linux and MacOS or `{path}/osb-checker-kotlin/gradlew.bat build` on windows.
Afterwards you can find `osb-checker-kotlin-1.0.jar` in `osb-checker-kotlin/build/libs`.

### Basic Run Configuration

The start the basic run configuration create a application.yml file at the same location of the .jar file similar to the example below.

```yaml
##Define the service broker connection here
config:
  url: http://localhost
  port : 80
  apiVersion: 2.15
  user: user
  password: password
```

Then call `java -jar osb-checker-kotlin-1.0.jar` on the commandline to start checker. In this configuration the checker will run all tests for every service-plan listed 
in the catalog. See the chapter [Usage](docs/Usage.md) for more details about configuring this test-application.

## Changes

Changes since v1.0:
- Added 2.15 feature testing:
    - Using [maintenance_info](docs/ProvisionTests.md#version-specific-tests) if defined in catalog.
    - Option to set [X-Broker-API-Request-Identity](docs/Usage.md#Configuration) header.
    - When testing polling maximum polling duration is tested.
- Added Option to set [X-Broker-API-Originating-Identity](docs/Usage.md#originating-identity).
- Added checks for osb-error-codes in response bodies.
- Added tests for fetching non existing bindings / provisions.
- Tests for what happens, when trying to create already existing provisions / bindings.
- Various improvements of logging such as using names of services and plans instead of id.
- setting up a custom catalog requires only the id's of the plan and no more additional information. This is gathered from the catalog now.
- Add Optional test for Dashboard URLs.
- Restructuring of binding test:
 - All plans in the catalog are now tested for invalid binding attempts.
 - Too reduce runtime invalid and valid binding tests use the same provision for testing instead.
