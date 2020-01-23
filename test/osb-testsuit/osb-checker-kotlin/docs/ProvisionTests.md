# Table of Contents
- [Description](../README.md#description)
- [Getting Started](../README.md#getting-started)
    - [Build Application](../README.md#build-application)
    - [Basic Run Configuration](../README.md#basic-run-configuration)
- [Usage](Usage.md)
    - [Declaring Test Runs](Usage.md#declaring-test-runs)
    - [Configuration](Usage.md#configuration)
    - [Parameters](Usage.md#parameters)
    - [Originating Identity](Usage.md#originating-identity)
    - [Declaring Services](Usage.md#declaring-services)
- Test Classes
    - [Catalog](CatalogTest.md)
       - [Example Output](CatalogTest.md#example-output)
    - [Provision](#provision-tests)
        - [Test Procedure](#test-procedure)
        - [Version specific Tests](#version-specific-tests)
        - [Example Output](#example-output)
    - [Binding](BindingTests.md#binding)
        - [Test Procedure](BindingTests.md#test-procedure)
        - [Version specific Tests](BindingTests.md#version-specific-tests)
        - [Example Output](BindingTests.md#example-output)
    - [Authentication](docs/AuthenticationTests.md)
    - [Contract](docs/ContractTest.md)
- [Contribution](docs/Contribution.md)
   
# Provision Tests

This Test Class checks that the service brokers behaviour when somebody tries to create or delete a service instance
with malformed request bodies, is spec compliant. Additionally it verifies that the service broker reacts correctly to synchronous
provision and deprovision attempts. 

## Test Procedure

The tests created in this class depends upon the catalog. This means that a valid catalog is required for this test to generate useful debugging information. 
It is highly recommended to ensure the Catalog Test Class runs successfully, before using this test class.

At the beginning of a test the catalog is fetched. The checker then runs the following tests based upon the provided information:

- Synchronous Test run
    - Runs a valid provision without query parameter `accepts_incomplete=true` and verifies if the response is correct according to [spec](https://github.com/openservicebrokerapi/servicebroker/blob/v2.15/spec.md#synchronous-and-asynchronous-operations)
    - Runs a valid deprovision without query parameter `accepts_incomplete=true` and verifies if if the response is correct according to [spec](https://github.com/openservicebrokerapi/servicebroker/blob/v2.15/spec.md#synchronous-and-asynchronous-operations)
- Runs invalid provisions where either the service_id or plan_id is empty, missing or not defined in the catalog and verifies if the response is correct according to [spec](https://github.com/openservicebrokerapi/servicebroker/blob/v2.15/spec.md#provisioning)
- Runs invalid de provisions where either the service_id or plan_id is empty, missing or not defined in the catalog and verifies if the response is correct according to [spec](https://github.com/openservicebrokerapi/servicebroker/blob/v2.15/spec.md#provisioning)
 
If you wish to avoid testing a entire catalog every time see [here](Usage.md#declaring-services) on how to cherry pick plans.
Depending on the API version and content of the catalog, additional tests are added.

## Version specific Tests

**2.14**
- If the service is declared to be fetchable, this class checks the service broker returns a 4XX error code when the requested service instance does not exist.
Look [here](https://github.com/openservicebrokerapi/servicebroker/blob/v2.15/spec.md#fetching-a-service-instance) for more information about fetching a service instance.

**2.15**
- If maintenance information is declared for the service plan, the checker verifies the correct error code is returned when a provision request contains a invalid maintenance_info field.
Afterwards a deprovision gets called with the same instanceId, to ensure no instances remain after the test.
Look [here](https://github.com/openservicebrokerapi/servicebroker/blob/v2.15/spec.md#error-codes) for more information about maintenance_info error code.

# Example Output

A successful Provision test run of a v2.15 service broker. The tested plan is fetchable and contains a maintenance information.

```
╷
└─ JUnit Jupiter ✔
   └─ Provision test runs ✔
      ├─ Run fetch, if fetchable, and synchronous operations ✔
      │  └─ Testing service 'base-sql-service-dev-managed' plan 's'. Using instanceId: 3183a819-6435-49d4-8a33-e907e105e7ad ✔
      │     ├─ should return 4XX when trying to retrieve a non existing service instance. ✔
      │     └─ should handle sync requests correctly ✔
      │        ├─ Sync PUT provision request ✔
      │        └─ Sync DELETE provision request ✔
      ├─ Run invalid asynchronous PUT requests ✔
      │  └─ Testing service 'base-sql-service-dev-managed' plan 's'. Using instanceId: a558bbd9-45fe-4038-8928-d974091d48d0 ✔
      │     ├─ PUT should reject if missing service_id ✔
      │     ├─ PUT should reject if missing plan_id ✔
      │     ├─ PUT should reject if missing service_id field ✔
      │     ├─ PUT should reject if missing plan_id field ✔
      │     ├─ PUT should reject if missing service_id field ✔
      │     ├─ PUT should reject if missing service_id field ✔
      │     ├─ PUT should reject if missing service_id is Invalid ✔
      │     ├─ PUT should reject if missing plan_id is Invalid ✔
      │     └─ Testing Maintenance Info ErrorCode and DELETE for clean up purposes. ✔
      │        ├─ PUT should reject if maintenance_info doesn't match ✔
      │        └─ DELETE should return 410 when trying to delete a non existing service instance, as it should not have been created in the previous test. ✔
      └─ Run invalid asynchronous DELETE requests ✔
         └─ Testing service 'base-sql-service-dev-managed' plan 's'. Using instanceId: 07e7b043-d17d-4d69-b278-92506c9dfe99 ✔
            ├─ DELETE should reject if service_id is missing ✔
            ├─ DELETE should reject if plan_id is missing ✔
            └─ DELETE should return 410 when trying to delete a non existing service instance ✔

Test run finished after 13487 ms
[        10 containers found      ]
[         0 containers skipped    ]
[        10 containers started    ]
[         0 containers aborted    ]
[        10 containers successful ]
[         0 containers failed     ]
[        16 tests found           ]
[         0 tests skipped         ]
[        16 tests started         ]
[         0 tests aborted         ]
[        16 tests successful      ]
[         0 tests failed          ]
```