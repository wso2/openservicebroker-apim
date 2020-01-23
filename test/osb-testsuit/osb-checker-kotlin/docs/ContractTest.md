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
    - [Provision](ProvisionTests.md#provision-tests)
        - [Test Procedure](ProvisionTests.md#test-procedure)
        - [Version specific Tests](ProvisionTests.md#version-specific-tests)
        - [Example Output](ProvisionTests.md#example-output)
    - [Binding](BindingTests.md#binding)
        - [Test Procedure](BindingTests.md#test-procedure)
        - [Version specific Tests](BindingTests.md#version-specific-tests)
        - [Example Output](BindingTests.md#example-output)
    - [Authentication](docs/AuthenticationTests.md)
    - [Contract](#contract)
- [Contribution](docs/Contribution.md)

# Contract

Tuns all standard requests and checks if they fail with 412 Precondition Failed, if the X-Broker-API-Version header is missing or does not match the given one.

## Example Output

```
╷
└─ JUnit Jupiter ✔
   └─ ContractJUnit5 ✔
      └─ testHeaderForAPIVersion() ✔
         └─ Requests should contain header X-Broker-API-Version ✔
            ├─ GET - v2/catalog should reject with 412 ✔
            ├─ PUT - v2/service_instance/instance_id should reject with 412 ✔
            ├─ DELETE - v2/service_instance/instance_id should reject with 412 ✔
            ├─ GET - v2/service_instance/instance_id/last_operation should reject with 412 ✔
            ├─ DELETE - v2/service_instance/instance_id?service_id=Invalid&plan_id=Invalid  should reject with 412) ✔
            ├─ PUT - v2/service_instance/instance_id/service_binding/binding_id  should reject with 412) ✔
            └─ DELETE - v2/service_instance/instance_id/service_binding/binding_id?service_id=Invalid&plan_id=Invalid should reject with 412 ✔

Test run finished after 4246 ms
[         4 containers found      ]
[         0 containers skipped    ]
[         4 containers started    ]
[         0 containers aborted    ]
[         4 containers successful ]
[         0 containers failed     ]
[         7 tests found           ]
[         0 tests skipped         ]
[         7 tests started         ]
[         0 tests aborted         ]
[         7 tests successful      ]
[         0 tests failed          ]
```