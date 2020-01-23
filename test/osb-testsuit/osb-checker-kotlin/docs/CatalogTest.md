# Table of Contents
- [Description](../README.md#description)
- [Getting Started](../README.md#getting-started)
    - [Build Application](../README.md#build-application)
    - [Basic Run Configuration](../README.md#basic-run-configuration)
- [Changes](../README.md#changes)
- [Usage](#usage)
    - [Declaring Test Runs](Usage.md#declaring-test-runs)
    - [Configuration](Usage.md#configuration)
    - [Parameters](Usage.md#parameters)
    - [Originating Identity](Usage.md#originating-identity)
    - [Declaring Services](Usage.md#declaring-services)
- Test Classes
    - [Catalog](#catalog)
    - [Provision](#provision)
        - [Test Procedure](ProvisionTests.md#test-procedure)
        - [Version specific Tests](ProvisionTests.md#version-specific-tests)
        - [Example Output](ProvisionTests.md#example-output)
    - [Binding](BindingTests.md#binding)
    - [Authentication](AuthenticationTests.md#authentication)
    - [Contract](ContractTest.md#contract)
- [Contribution](docs/Contribution.md)

# Catalog Test

The Catalog Tests verifies that catalog endpoint returns status-code 200 OK and a valid service broker catalog according to spec.

when starting the application with the parameter -cat/-catalog, it will:
- call `curl http://username:password@broker-url/v2/catalog -X GET -H "X-Broker-API-Version: api-version-here" -H "Content-Type: application/json"`
- check if the service broker returns 200 and validate if the catalog from the response follows schema.

Look [here](https://github.com/openservicebrokerapi/servicebroker/blob/v2.15/spec.md#catalog-management), for more information about service broker catalogs.
A valid catalog is crucial for all following tests, since the checker uses it's content to figure out what tests should run. It's highly recommended to make sure
the catalog is implemented correctly before continuing implementing the other endpoints.

## Example Output

```
╷
└─ JUnit Jupiter ✔
   └─ Catalog test ✔
      └─ Verify catalog schema ✔

```