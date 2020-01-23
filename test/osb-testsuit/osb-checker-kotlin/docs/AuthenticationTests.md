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
    - [Binding](BindingTests.md#binding)
        - [Test Procedure](BindingTests.md#test-procedure)
        - [Version specific Tests](BindingTests.md#version-specific-tests)
        - [Example Output](BindingTests.md#example-output)
    - [Authentication](#authentication)
    - [Contract](docs/ContractTest.md)
- [Contribution](docs/Contribution.md)
   
# Authentication

Runs a all requests without a user and password, a wrong username and a wrong password. It checks if service broker replies with HttpStatus 401 unauthorized.

## Example Output

```
╷
└─ JUnit Jupiter ✔
   └─ AuthenticationJUnit5 ✔
      └─ testAuthentication() ✔
         └─ Requests without authentication should be rejected ✔
            ├─ GET - v2/catalog should reject with 401 ✔
            │  ├─ Without authentication ✔
            │  ├─ With wrong Username ✔
            │  └─ With wrong Password ✔
            ├─ PUT - v2/service_instance/instance_id should reject with 401 ✔
            │  ├─ Without authentication ✔
            │  ├─ With wrong Username ✔
            │  └─ With wrong Password ✔
            ├─ DELETE - v2/service_instance/instance_id should reject with 401 ✔
            │  ├─ Without authentication ✔
            │  ├─ With wrong Username ✔
            │  └─ With wrong Password ✔
            ├─ GET - v2/service_instance/instance_id/last_operation should reject with 401 ✔
            │  ├─ Without authentication ✔
            │  ├─ With wrong Username ✔
            │  └─ With wrong Password ✔
            ├─ PUT - v2/service_instance/instance_id/service_binding/binding_id  should reject with 401 ✔
            │  ├─ Without authentication ✔
            │  ├─ With wrong Username ✔
            │  └─ With wrong Password ✔
            └─ DELETE - v2/service_instance/instance_id/service_binding/binding_id  should reject with 401 ✔
               ├─ Without authentication ✔
               ├─ With wrong Username ✔
               └─ With wrong Password ✔

Test run finished after 4720 ms
[        10 containers found      ]
[         0 containers skipped    ]
[        10 containers started    ]
[         0 containers aborted    ]
[        10 containers successful ]
[         0 containers failed     ]
[        18 tests found           ]
[         0 tests skipped         ]
[        18 tests started         ]
[         0 tests aborted         ]
[        18 tests successful      ]
[         0 tests failed          ]
```