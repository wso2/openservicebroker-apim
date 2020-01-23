package de.evoila.osb.checker.tests.contract

import de.evoila.osb.checker.request.BindingRequestRunner
import de.evoila.osb.checker.request.ProvisionRequestRunner
import de.evoila.osb.checker.tests.TestBase
import org.junit.jupiter.api.DynamicContainer
import org.junit.jupiter.api.DynamicNode
import org.junit.jupiter.api.DynamicTest.dynamicTest
import org.junit.jupiter.api.TestFactory
import org.springframework.beans.factory.annotation.Autowired

class ContractJUnit5 : TestBase() {

  @Autowired
  lateinit var bindingRequestRunner: BindingRequestRunner
  @Autowired
  lateinit var provisionRequestRunner: ProvisionRequestRunner

  @TestFactory
  fun testHeaderForAPIVersion(): List<DynamicNode> {
    return listOf(
        DynamicContainer.dynamicContainer("Requests should contain header X-Broker-API-Version",
            listOf(
                dynamicTest("GET - v2/catalog should reject with 412")
                { catalogRequestRunner.withoutHeader() },
                dynamicTest("PUT - v2/service_instance/instance_id should reject with 412")
                { provisionRequestRunner.putWithoutHeader() },
                dynamicTest("DELETE - v2/service_instance/instance_id should reject with 412")
                { provisionRequestRunner.deleteWithoutHeader() },
                dynamicTest("GET - v2/service_instance/instance_id/last_operation should reject with 412")
                { provisionRequestRunner.lastOperationWithoutHeader() },
                dynamicTest("DELETE - v2/service_instance/instance_id?service_id=Invalid&plan_id=Invalid  should reject with 412)")
                { provisionRequestRunner.deleteWithoutHeader() },
                dynamicTest("PUT - v2/service_instance/instance_id/service_binding/binding_id  should reject with 412)")
                { bindingRequestRunner.putWithoutHeader() },
                dynamicTest("DELETE - v2/service_instance/instance_id/service_binding/binding_id?service_id=Invalid&plan_id=Invalid should reject with 412")
                { bindingRequestRunner.deleteWithoutHeader() }
            )
        )
    )
  }
}