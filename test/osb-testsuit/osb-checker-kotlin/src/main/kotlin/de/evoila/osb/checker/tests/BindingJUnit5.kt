package de.evoila.osb.checker.tests

import de.evoila.osb.checker.request.BindingRequestRunner
import de.evoila.osb.checker.request.ResponseBodyType.ERR
import de.evoila.osb.checker.request.ResponseBodyType.VALID_BINDING
import de.evoila.osb.checker.request.bodies.BindingBody
import de.evoila.osb.checker.request.bodies.Person
import de.evoila.osb.checker.request.bodies.ProvisionBody
import de.evoila.osb.checker.response.catalog.Plan
import de.evoila.osb.checker.response.catalog.Service
import de.evoila.osb.checker.tests.containers.BindingContainers
import org.junit.jupiter.api.DisplayName
import org.junit.jupiter.api.DynamicContainer
import org.junit.jupiter.api.DynamicContainer.dynamicContainer
import org.junit.jupiter.api.DynamicNode
import org.junit.jupiter.api.DynamicTest.dynamicTest
import org.junit.jupiter.api.TestFactory
import org.springframework.beans.factory.annotation.Autowired
import java.util.*

@DisplayName(value = "Binding Tests")
class BindingJUnit5 : TestBase() {

    @Autowired
    lateinit var bindingRequestRunner: BindingRequestRunner
    @Autowired
    lateinit var bindingContainerFactory: BindingContainers

    @TestFactory
    @DisplayName(value = "Valid Provision and Binding Tests.")
    fun validProvisionAndBindingTestRuns(): List<DynamicNode> {

        return configuration.initCustomCatalog(catalogRequestRunner.correctRequest()).services.flatMap { service ->
            service.plans.map { plan ->
                val instanceId = UUID.randomUUID().toString()
                val bindingId = UUID.randomUUID().toString()
                val needsAppGuid = needAppGUID(plan)
                val provision = if (configuration.apiVersion >= 2.15 && plan.maintenanceInfo != null)
                    ProvisionBody.ValidProvisioning(service, plan, plan.maintenanceInfo)
                else ProvisionBody.ValidProvisioning(service, plan)

                val binding = if (needsAppGuid) BindingBody(
                        serviceId = service.id,
                        planId = plan.id,
                        appGuid = UUID.randomUUID().toString()
                )
                else BindingBody(serviceId = service.id, planId = plan.id)

                configuration.provisionParameters.let {
                    if (it.containsKey(plan.id)) {
                        provision.parameters = mapOf("apis" to arrayOf(Person("PizzaShackAPI", "1.0.0"))) as Map<String, Any>
                    }
                }
                configuration.bindingParameters.let {
                    if (it.containsKey(plan.id)) {
                        binding.parameters = it[plan.id]
                    }
                }
                val dynamicContainers = mutableListOf(
                        bindingContainerFactory.validProvisionContainer(
                                instanceId = instanceId,
                                plan = plan,
                                provision = provision,
                                isRetrievable = configuration.apiVersion > 2.13 && (service.instancesRetrievable
                                        ?: false),
                                serviceName = service.name,
                                planName = plan.name
                        )
                )

                val bindable = planIsBindable(service, plan)
                if (bindable) {
                    dynamicContainers.add(dynamicContainer("Service ${service.name} Plan ${plan.name} is bindable. Testing binding operation with bindingId $bindingId", mutableListOf(
                            createSyncAndInvalidBindingTests(service, plan, needsAppGuid, instanceId, bindingId),
                            bindingContainerFactory.validBindingContainer(
                                    binding = binding,
                                    instanceId = instanceId,
                                    bindingId = bindingId,
                                    isRetrievable = configuration.apiVersion > 2.13 && service.bindingsRetrievable ?: false,
                                    plan = plan
                            ))
                    ))
                }

                dynamicContainers.add(bindingContainerFactory.validDeleteProvisionContainer(instanceId, service, plan))
                dynamicContainer(if (bindable) {
                    BINDABLE_MESSAGE
                } else {
                    PROVISION_DISPLAY_NAME
                }, dynamicContainers)
            }
        }
    }

    private fun createSyncAndInvalidBindingTests(
            service: Service,
            plan: Plan,
            needsAppGuid: Boolean,
            instanceId: String,
            bindingId: String
    ): DynamicContainer {
        val bindingTests = mutableListOf<DynamicNode>()

        service.bindingsRetrievable?.let { bindingRetrievable ->
            if (configuration.apiVersion >= 2.14 && bindingRetrievable) {
                bindingTests.add(dynamicTest("should return status code 4XX when tying to fetch a non existing binding") {
                    bindingRequestRunner.runGetBindingRequest(instanceId, bindingId, *IntArray(100) { 400 + it })
                })
            }
        }
        if (configuration.apiVersion >= 2.14) {
            bindingTests.add(dynamicContainer("should handle sync requests correctly",
                    bindingContainerFactory.createSyncBindingTest(
                            binding = if (needsAppGuid) BindingBody(
                                    serviceId = service.id,
                                    planId = plan.id,
                                    appGuid = UUID.randomUUID().toString())
                            else BindingBody(service.id, plan.id),
                            bindingId = bindingId,
                            instanceId = instanceId
                    )
            ))
        }
        bindingTests.addAll(listOf(
                TestCase(
                        requestBody =
                        if (needsAppGuid) BindingBody(null, plan.id, UUID.randomUUID().toString())
                        else BindingBody(null, plan.id),
                        message = "should reject if missing service_id",
                        responseBodyType = VALID_BINDING,
                        statusCode = 400
                ),
                TestCase(
                        requestBody = if (needsAppGuid) BindingBody(service.id, null, UUID.randomUUID().toString())
                        else BindingBody(service.id, null),
                        message = "should reject if missing plan_id",
                        responseBodyType = VALID_BINDING,
                        statusCode = 400
                )
        ).flatMap { testCase ->
            listOf(
                    dynamicTest("PUT ${testCase.message}") {
                        bindingRequestRunner.runPutBindingRequestAsync(
                                requestBody = testCase.requestBody,
                                instanceId = instanceId,
                                bindingId = bindingId,
                                expectedStatusCodes = *intArrayOf(testCase.statusCode),
                                expectedResponseBody = ERR
                        )
                    },
                    dynamicTest("DELETE ${testCase.message}") {
                        val bindingRequestBody = testCase.requestBody
                        bindingRequestRunner.runDeleteBindingRequestAsync(
                                serviceId = bindingRequestBody.serviceId,
                                planId = bindingRequestBody.planId,
                                instanceId = instanceId,
                                bindingId = bindingId,
                                expectedStatusCodes = *intArrayOf(testCase.statusCode)
                        )
                    }
            )
        })
        return dynamicContainer("Run sync and invalid bindings attempts", bindingTests)
    }

    private fun planIsBindable(service: Service, plan: Plan): Boolean = plan.bindable ?: service.bindable

    fun needAppGUID(plan: Plan): Boolean = plan.metadata?.customParameters?.usesServicesKeys
            ?: configuration.usingAppGuid

    companion object {
        private const val BINDABLE_MESSAGE = "Running a valid provision and run binding tests." +
                " Delete both afterwards. In case of a asynchronous service broker polling after each operation."
        private const val PROVISION_DISPLAY_NAME = "Running a valid provision and if necessary polling. Deleting it afterwards."
    }
}
