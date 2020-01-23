package de.evoila.osb.checker.tests

import de.evoila.osb.checker.request.ProvisionRequestRunner
import de.evoila.osb.checker.request.ResponseBodyType.*
import de.evoila.osb.checker.request.bodies.Person
import de.evoila.osb.checker.request.bodies.ProvisionBody
import de.evoila.osb.checker.request.bodies.ProvisionBody.ValidProvisioning
import de.evoila.osb.checker.response.catalog.MaintenanceInfo
import de.evoila.osb.checker.response.catalog.Plan
import de.evoila.osb.checker.response.catalog.Service
import org.junit.jupiter.api.DisplayName
import org.junit.jupiter.api.DynamicContainer
import org.junit.jupiter.api.DynamicContainer.dynamicContainer
import org.junit.jupiter.api.DynamicNode
import org.junit.jupiter.api.DynamicTest.dynamicTest
import org.junit.jupiter.api.TestFactory
import org.springframework.beans.factory.annotation.Autowired
import java.util.*
import kotlin.collections.LinkedHashMap

@DisplayName(value = "Provision test runs")
class ProvisionJUnit5 : TestBase() {

    @Autowired
    lateinit var provisionRequestRunner: ProvisionRequestRunner

    @TestFactory
    @DisplayName(value = "Run fetch, if fetchable, and synchronous operations")
    fun runGetInstanceAndSyncTests(): List<DynamicNode> {
        val catalog = configuration.initCustomCatalog(catalogRequestRunner.correctRequest())

        return catalog.services.flatMap { service ->
            service.plans.map { plan ->
                createSyncAndRetrievableTestForPlan(service, plan)
            }
        }
    }

    @TestFactory
    @DisplayName(value = "Run invalid asynchronous PUT requests")
    fun runInvalidAsyncPutTest(): List<DynamicNode> {
        val catalog = configuration.initCustomCatalog(catalogRequestRunner.correctRequest())

        return catalog.services.flatMap { service ->
            service.plans.map { plan ->
                createInvalidPutTestsForPlan(service, plan)
            }
        }
    }

    @TestFactory
    @DisplayName(value = "Run invalid asynchronous DELETE requests")
    fun runInvalidAsyncDeleteTest(): List<DynamicNode> {
        val catalog = configuration.initCustomCatalog(catalogRequestRunner.correctRequest())

        return catalog.services.flatMap { service ->
            service.plans.map { plan ->
                createInvalidAsyncDeleteTestsForPlan(service, plan)
            }
        }
    }

    private fun createInvalidPutTestsForPlan(service: Service, plan: Plan): DynamicContainer {
        val instanceId = UUID.randomUUID().toString()
        val dynamicNodes = listOf(
                TestCase(
                        requestBody = ValidProvisioning(
                                service_id = "",
                                plan_id = plan.id,
                                maintenance_info = if (configuration.apiVersion >= 2.15) {
                                    plan.maintenanceInfo
                                } else {
                                    null
                                }
                        ),
                        message = "should reject if missing service_id",
                        responseBodyType = ERR,
                        statusCode = 400
                ),
                TestCase(
                        requestBody = ValidProvisioning(
                                service_id = service.id,
                                plan_id = "",
                                maintenance_info = if (configuration.apiVersion >= 2.15) {
                                    plan.maintenanceInfo
                                } else {
                                    null
                                }
                        ),
                        message = "should reject if missing plan_id",
                        responseBodyType = ERR,
                        statusCode = 400
                ),
                TestCase(
                        requestBody = ProvisionBody.NoServiceFieldProvisioning(plan),
                        message = "should reject if missing service_id field",
                        responseBodyType = ERR,
                        statusCode = 400
                ),
                TestCase(
                        requestBody = ProvisionBody.NoPlanFieldProvisioning(service),
                        message = "should reject if missing plan_id field",
                        responseBodyType = ERR,
                        statusCode = 400
                ),
                TestCase(
                        requestBody = ProvisionBody.NoSpaceFieldProvisioning(service, plan),
                        message = "should reject if missing space_id field",
                        responseBodyType = ERR,
                        statusCode = 400
                ),
                TestCase(requestBody = ProvisionBody.NoOrgFieldProvisioning(service, plan),
                        message = "should reject if missing org_id field",
                        responseBodyType = ERR,
                        statusCode = 400
                ),
                TestCase(
                        requestBody = ValidProvisioning(
                                "Invalid", plan.id,
                                maintenance_info = if (configuration.apiVersion == 2.15) {
                                    plan.maintenanceInfo
                                } else {
                                    null
                                }
                        ),
                        message = "should reject if missing service_id is Invalid",
                        responseBodyType = ERR,
                        statusCode = 400
                ),
                TestCase(
                        requestBody = ValidProvisioning(
                                service.id, "Invalid",
                                maintenance_info = if (configuration.apiVersion >= 2.15) {
                                    plan.maintenanceInfo
                                } else {
                                    null
                                }
                        ),
                        message = "should reject if missing plan_id is Invalid",
                        responseBodyType = ERR,
                        statusCode = 400
                )
        ).map {
            dynamicTest("PUT ${it.message}") {
                provisionRequestRunner.runPutProvisionRequestAsync(instanceId,
                        requestBody = it.requestBody,
                        expectedFinalStatusCodes = *intArrayOf(it.statusCode),
                        expectedResponseBodyType = it.responseBodyType
                )
            }
        }

        return dynamicContainer(createDisplayName(service.name, plan.name, instanceId), if (configuration.apiVersion >= 2.15) {
            ValidProvisioning(service, plan, MaintenanceInfo("Invalid", "Should return 422"))
            dynamicNodes.plus(
                    dynamicContainer("Testing Maintenance Info ErrorCode and DELETE for clean up purposes.", listOf(
                            dynamicTest("PUT should reject if maintenance_info doesn't match") {
                                provisionRequestRunner.runPutProvisionRequestAsync(instanceId,
                                        requestBody = ValidProvisioning(
                                                service = service,
                                                plan = plan,
                                                maintenance_info = MaintenanceInfo("Invalid", "Should return 422")
                                        ),
                                        expectedFinalStatusCodes = *intArrayOf(422),
                                        expectedResponseBodyType = ERR_MAINTENANCE_INFO
                                )
                            },
                            dynamicTest("DELETE should return 410 when trying to delete a non existing service instance, as it should not have been created in the previous test.") {
                                provisionRequestRunner.runDeleteProvisionRequestAsync(
                                        serviceId = nullIfNotSet(service.id),
                                        planId = nullIfNotSet(plan.id),
                                        instanceId = instanceId,
                                        expectedFinalStatusCodes = intArrayOf(410)
                                )
                            }
                    ))
            )
        } else dynamicNodes
        )
    }

    private fun createInvalidAsyncDeleteTestsForPlan(service: Service, plan: Plan): DynamicContainer {
        val instanceId = UUID.randomUUID().toString()

        val dynamicNodes = listOf(
                TestCase(
                        message = "should reject if service_id is missing",
                        requestBody = ValidProvisioning("", plan.id),
                        responseBodyType = ERR,
                        statusCode = 400
                ),
                TestCase(
                        message = "should reject if plan_id is missing",
                        requestBody = ValidProvisioning(service.id, ""),
                        responseBodyType = ERR,
                        statusCode = 400
                ),
                TestCase(
                        message = "should return 410 when trying to delete a non existing service instance",
                        requestBody = ValidProvisioning(service, plan),
                        responseBodyType = NO_SCHEMA,
                        statusCode = 410
                )
        ).map {
            dynamicTest("DELETE ${it.message}") {
                val provisionBody = it.requestBody
                provisionRequestRunner.runDeleteProvisionRequestAsync(
                        serviceId = nullIfNotSet(provisionBody.service_id),
                        planId = nullIfNotSet(provisionBody.plan_id),
                        instanceId = instanceId,
                        expectedFinalStatusCodes = intArrayOf(it.statusCode)
                )
            }
        }

        return dynamicContainer(createDisplayName(service.name, plan.name, instanceId), dynamicNodes)
    }

    private fun createSyncAndRetrievableTestForPlan(service: Service, plan: Plan): DynamicContainer {
        val instanceId = UUID.randomUUID().toString()
        val provisionRequestBody = if (configuration.apiVersion >= 2.15 && plan.maintenanceInfo != null) {
            ValidProvisioning(service, plan, plan.maintenanceInfo)
        } else {
            ValidProvisioning(service, plan)
        }
        // set configuration parameters for API plan provisioning request.
        configuration.provisionParameters.let {
            if (it.containsKey(plan.id)) {
                val appPlanID = "00e851cd-ce8b-43eb-bc27-ac4d4fbb3204"

                if (plan.id == appPlanID) {
//                     configuration.provisionParameters[appPlanID]?.get("apis") = arrayOf(Person("PizzaShackAPI", "1.0.0"))
//                    apiParam.let { itParam ->
//                        (itParam) as LinkedHashMap<String, Any>
//                        itParam["transport"] = arrayOf<String>("http", "https")
//                        itParam["tiers"] = arrayOf<String>("Unlimited")
//                        itParam["tags"] = arrayOf<String>("Pizza")
//                    }
                    val apiParamW = mapOf("apis" to arrayOf(Person("PizzaShackAPI", "1.0.0")))
                    provisionRequestBody.parameters = apiParamW as Map<String, Any>
                } else {
                    provisionRequestBody.parameters = it[plan.id]
                }
            }
        }

        val dynamicNodes = mutableListOf<DynamicNode>()
        service.instancesRetrievable?.let { instancesRetrievable ->
            if (configuration.apiVersion >= 2.14 && instancesRetrievable) {
                dynamicNodes.add(dynamicTest("should return 4XX when trying to retrieve a non existing service instance.") {
                    provisionRequestRunner.getProvision(instanceId, *IntArray(100) { 400 + it })
                })
            }
        }

        dynamicNodes.add(
                dynamicContainer("should handle sync requests correctly", listOf(
                        dynamicTest("Sync PUT provision request") {
                            provisionRequestRunner.runPutProvisionRequestSync(instanceId, provisionRequestBody)
                        },
                        dynamicTest("Sync DELETE provision request") {
                            provisionRequestRunner.runDeleteProvisionRequestSync(
                                    instanceId = instanceId,
                                    serviceId = provisionRequestBody.service_id,
                                    planId = provisionRequestBody.plan_id)
                        }
                ))
        )

        return dynamicContainer(createDisplayName(service.name, plan.name, instanceId), dynamicNodes)
    }

    private fun createDisplayName(serviceName: String, planName: String, instanceId: String): String =
            "Testing service '${serviceName}'" +
                    " plan '${planName}'. Using instanceId: $instanceId"

    private fun nullIfNotSet(value: String): String? {
        return if (value.isNotEmpty()) {
            value
        } else {
            null
        }
    }
}