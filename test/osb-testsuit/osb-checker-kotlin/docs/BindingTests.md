# Table of Contents
- [Description](../README.md#description)
- [Getting Started](../README.md#getting-started)
    - [Build Application](../README.md#build-application)
    - [Basic Run Configuration](../README.md#basic-run-configuration)
- [Changes](../README.md#changes)
- [Usage](Usage.md)
    - [Declaring Test Runs](Usage.md#declaring-test-runs)
    - [Configuration](Usage.md#configuration)
    - [Parameters](Usage.md#parameters)
    - [Originating Identity](Usage.md#originating-identity)
    - [Declaring Services](Usage.md#declaring-services)
- Test Classes
    - [Catalog](CatalogTest.md)
       - [Example Output](CatalogTest.md#example-output)
    - [Provision](ProvisionTests.md#provision-tests)
        - [Test Procedure](ProvisionTests.md#test-procedure)
        - [Version specific Tests](ProvisionTests.md#version-specific-tests)
        - [Example Output](ProvisionTests.md#example-output)
    - [Binding](#binding-tests)
        - [Test Procedure](#test-procedure)
        - [Version specific Tests](#version-specific-tests)
        - [Example Output](#example-output)
    - [Authentication](docs/AuthenticationTests.md)   
    - [Contract](docs/ContractTest.md)
- [Contribution](docs/Contribution.md)
   
# Binding Tests

This Test Class provisions every plan present in the service broker and if the service instance is bindable, it binds on it.
Afterwards the binding and provision is deleted. During this process the checker verifies that the service brokers reacts according to spec.
The checker also validates the behaviour of the service broker, when invalid binding requests are being send to it.

## Test Procedure

The tests created in this class depends upon the catalog. This means that a valid catalog is required for this test to generate useful debugging information. 
It is highly recommended to ensure the Catalog Test Class runs successfully, before using this test class. Naturally all binding tests will fail, if the underlying
service instances don't work, so it is recommended that the [provision tests](ProvisionTests.md#provision-tests) are succeeding too. 

At the beginning of a test the catalog is fetched. The checker then runs the following tests based upon the provided information:

- Valid Provision and Bindings
    - Create a valid provision.
        - If the Service broker creates service instances asynchronously, the checker will start polling and and verify the responses.
        - When configured to do so, the checker verifies if the dashboard URL works.
        are according to [spec](https://github.com/openservicebrokerapi/servicebroker/blob/v2.15/spec.md#polling-last-operation-for-service-instances) and finish successfully.
        - Test what happens when attempting to create a service instances with the same instance id and same parameters and with different parameters.
     Read [here](https://github.com/openservicebrokerapi/servicebroker/blob/v2.15/spec.md#polling-last-operation-for-service-instances) about the expected behaviour.
    - If the Service is bindable. Look [here](https://github.com/openservicebrokerapi/servicebroker/blob/v2.15/spec.md#binding) on how it should behave.
        - Binding Attempts with missing mandatory or malformed data.
        - Valid Bindings und un bindings.
    - Delete the provision
     
## Version specific Tests
**2.14**
- If the [instances](https://github.com/openservicebrokerapi/servicebroker/blob/v2.14/spec.md#fetching-a-service-instance)
 and [binding](https://github.com/openservicebrokerapi/servicebroker/blob/v2.14/spec.md#fetching-a-service-instance) are the checker tests they perform correctly.
- Tests how the broker reacts with synchronous and asynchronous binding attempts. More information [here](https://github.com/openservicebrokerapi/servicebroker/blob/v2.14/spec.md#fetching-a-service-instance)

**2.15**
- If the catalog contains [maintenance information](https://github.com/openservicebrokerapi/servicebroker/blob/master/spec.md#maintenance-info-object) the checker will uses it when requesting a provision.
- If the catalog contains a [maximum polling duration](https://github.com/openservicebrokerapi/servicebroker/blob/master/spec.md#polling-interval-and-duration)
 the checker will verify that the service broker acts accordingly.


### Example output
A Binding Test output with two services. 'base-sql-service-dev-managed' is bindable, instanceFetchable and bindingFetchable.
'base-sql-service-dev-unmanaged' is not. Note how tests run accordingly.

```
╷
└─ JUnit Jupiter ✔
   └─ Binding Tests ✔
      └─ Valid Provision and Binding Tests. ✔
         ├─ Running a valid provision and run binding tests. Delete both afterwards. In case of a asynchronous service broker polling after each operation. ✔
         │  ├─ Creating Service Instance, test dashboard URL, and try to fetch it. ✔
         │  │  ├─ Running valid PUT provision with instanceId 9eeea6a1-478a-4f47-abb4-eebd6d149c4d for service 'base-sql-service-dev-managed' and plan 's' ✔
         │  │  ├─ Running valid PUT provision with same attributes again. Expecting Status 200. ✔
         │  │  └─ Running valid PUT provision with different attributes again. Expecting Status 409. ✔
         │  ├─ Service base-sql-service-dev-managed Plan s is bindable. Testing binding operation with bindingId 27285fe4-c124-4a22-85be-f6db91de8373 ✔
         │  │  ├─ Run sync and invalid bindings attempts ✔
         │  │  │  ├─ should return status code 4XX when tying to fetch a non existing binding ✔
         │  │  │  ├─ should handle sync requests correctly ✔
         │  │  │  │  ├─ Sync PUT binding request ✔
         │  │  │  │  └─ Sync DELETE binding request ✔
         │  │  │  ├─ PUT should reject if missing service_id ✔
         │  │  │  ├─ DELETE should reject if missing service_id ✔
         │  │  │  ├─ PUT should reject if missing plan_id ✔
         │  │  │  └─ DELETE should reject if missing plan_id ✔
         │  │  └─ Running PUT binding and DELETE binding afterwards ✔
         │  │     ├─ Running valid PUT binding with bindingId 27285fe4-c124-4a22-85be-f6db91de8373 ✔
         │  │     ├─ Running PUT binding with same attribute again. Expecting StatusCode 200. ✔
         │  │     ├─ Running PUT binding with different attribute again. Expecting StatusCode 409. ✔
         │  │     ├─ Running GET for retrievable service binding and expecting StatusCode: 200 ✔
         │  │     └─ Deleting binding with bindingId 27285fe4-c124-4a22-85be-f6db91de8373 ✔
         │  └─ Deleting provision ✔
         │     ├─ DELETE provision and if the service broker is async polling afterwards ✔
         │     └─ Running valid DELETE provision with same parameters again. Expecting Status 410. ✔
         └─ Running a valid provision and if necessary polling. Deleting it afterwards. ✔
            ├─ Creating Service Instance, test dashboard URL. ✔
            │  ├─ Running valid PUT provision with instanceId a408b58d-aa4b-4c53-b851-f3d331cdfe43 for service 'base-sql-service-dev-unmanaged' and plan 's' ✔
            │  ├─ Running valid PUT provision with same attributes again. Expecting Status 200. ✔
            │  └─ Running valid PUT provision with different attributes again. Expecting Status 409. ✔
            └─ Deleting provision ✔
               ├─ DELETE provision and if the service broker is async polling afterwards ✔
               └─ Running valid DELETE provision with same parameters again. Expecting Status 410. ✔

Test run finished after 8786 ms
[        13 containers found      ]
[         0 containers skipped    ]
[        13 containers started    ]
[         0 containers aborted    ]
[        13 containers successful ]
[         0 containers failed     ]
[        22 tests found           ]
[         0 tests skipped         ]
[        22 tests started         ]
[         0 tests aborted         ]
[        22 tests successful      ]
[         0 tests failed          ]
```
