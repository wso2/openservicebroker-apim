package de.evoila.osb.checker.tests.contract

import de.evoila.osb.checker.request.BindingRequestRunner
import de.evoila.osb.checker.request.ProvisionRequestRunner
import de.evoila.osb.checker.tests.TestBase
import org.junit.jupiter.api.DynamicContainer
import org.junit.jupiter.api.DynamicNode
import org.junit.jupiter.api.DynamicTest
import org.junit.jupiter.api.TestFactory
import org.springframework.beans.factory.annotation.Autowired


class AuthenticationJUnit5 : TestBase() {
  @Autowired
  lateinit var bindingRequestRunner: BindingRequestRunner
  @Autowired
  lateinit var provisionRequestRunner: ProvisionRequestRunner

  @TestFactory
  fun testAuthentication(): List<DynamicNode> = listOf(
      DynamicContainer.dynamicContainer("Requests without authentication should be rejected",
          listOf(
              DynamicContainer.dynamicContainer("GET - v2/catalog should reject with 401",
                  listOf(
                      DynamicTest.dynamicTest(NO_AUTH) { catalogRequestRunner.noAuth() },
                      DynamicTest.dynamicTest(WRONG_USER) { catalogRequestRunner.wrongUser() },
                      DynamicTest.dynamicTest(WRONG_PW) { catalogRequestRunner.wrongPassword() }
                  )),
              DynamicContainer.dynamicContainer("PUT - v2/service_instance/instance_id should reject with 401",
                  listOf(
                      DynamicTest.dynamicTest(NO_AUTH) { provisionRequestRunner.putNoAuth() },
                      DynamicTest.dynamicTest(WRONG_USER) { provisionRequestRunner.putWrongUser() },
                      DynamicTest.dynamicTest(WRONG_PW) { provisionRequestRunner.putWrongPassword() }
                  )),
              DynamicContainer.dynamicContainer("DELETE - v2/service_instance/instance_id should reject with 401",
                  listOf(
                      DynamicTest.dynamicTest(NO_AUTH) { provisionRequestRunner.deleteNoAuth() },
                      DynamicTest.dynamicTest(WRONG_USER) { provisionRequestRunner.deleteWrongUser() },
                      DynamicTest.dynamicTest(WRONG_PW) { provisionRequestRunner.deleteWrongPassword() }
                  )),
              DynamicContainer.dynamicContainer("GET - v2/service_instance/instance_id/last_operation should reject with 401",
                  listOf(
                      DynamicTest.dynamicTest(NO_AUTH) { provisionRequestRunner.lastOpNoAuth() },
                      DynamicTest.dynamicTest(WRONG_USER) { provisionRequestRunner.lastOpWrongUser() },
                      DynamicTest.dynamicTest(WRONG_PW) { provisionRequestRunner.lastOpWrongPassword() }
                  )),
              DynamicContainer.dynamicContainer("PUT - v2/service_instance/instance_id/service_binding/binding_id  should reject with 401",
                  listOf(
                      DynamicTest.dynamicTest(NO_AUTH) { bindingRequestRunner.putNoAuth() },
                      DynamicTest.dynamicTest(WRONG_USER) { bindingRequestRunner.putWrongUser() },
                      DynamicTest.dynamicTest(WRONG_PW) { bindingRequestRunner.putWrongPassword() }
                  )),
              DynamicContainer.dynamicContainer("DELETE - v2/service_instance/instance_id/service_binding/binding_id  should reject with 401",
                  listOf(
                      DynamicTest.dynamicTest(NO_AUTH) { bindingRequestRunner.deleteNoAuth() },
                      DynamicTest.dynamicTest(WRONG_USER) { bindingRequestRunner.deleteWrongUser() },
                      DynamicTest.dynamicTest(WRONG_PW) { bindingRequestRunner.deleteWrongPassword() }
                  ))
          )
      )
  )

  companion object {
    const val NO_AUTH = "Without authentication"
    const val WRONG_USER = "With wrong Username"
    const val WRONG_PW = "With wrong Password"
  }
}