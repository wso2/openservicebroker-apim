package de.evoila.osb.checker.request

import de.evoila.osb.checker.config.Configuration
import de.evoila.osb.checker.request.ResponseBodyType.*
import de.evoila.osb.checker.request.bodies.RequestBody
import de.evoila.osb.checker.response.operations.LastOperationResponse
import io.restassured.RestAssured
import io.restassured.http.ContentType
import io.restassured.http.Header
import io.restassured.module.jsv.JsonSchemaValidator
import io.restassured.response.ExtractableResponse
import io.restassured.response.Response
import org.hamcrest.collection.IsIn
import org.springframework.stereotype.Service
import java.time.Instant
import java.util.*
import kotlin.test.assertTrue

@Service
class BindingRequestRunner(configuration: Configuration) : PollingRequestHandler(configuration) {

    fun runGetBindingRequest(instanceId: String, bindingId: String, vararg expectedStatusCodes: Int) {
        if (configuration.apiVersion >= 2.15 && configuration.useRequestIdentity) {
            useRequestIdentity("OSB-Checker-GET-binding-${UUID.randomUUID()}")
        }

        val response = RestAssured.with()
                .log().ifValidationFails()
                .headers(validRequestHeaders)
                .contentType(ContentType.JSON)
                .get(SERVICE_INSTANCE_PATH + instanceId + SERVICE_BINDING_PATH + bindingId)
                .then()
                .log().ifValidationFails()
                .assertThat()
                .headers(expectedResponseHeaders)
                .statusCode(IsIn(expectedStatusCodes.asList()))
                .extract()
                .response()

        if (expectedStatusCodes.contentEquals(intArrayOf(200))) {
            assertTrue(JsonSchemaValidator.matchesJsonSchemaInClasspath(VALID_FETCH_BINDING.path)
                    .matches(response.jsonPath().prettify()),
                    "\nExpected a valid GET binding response but was:\n\"${response.jsonPath().prettify()}\"")
        }
    }

    fun runPutBindingRequestSync(
            requestBody: RequestBody,
            instanceId: String,
            bindingId: String
    ) {
        if (configuration.apiVersion >= 2.15 && configuration.useRequestIdentity) {
            useRequestIdentity("OSB-Checker-PUT-binding-${UUID.randomUUID()}")
        }

        val response = RestAssured.with()
                .log().ifValidationFails()
                .headers(validRequestHeaders)
                .contentType(ContentType.JSON)
                .body(requestBody)
                .put(SERVICE_INSTANCE_PATH + instanceId + SERVICE_BINDING_PATH + bindingId)
                .then()
                .assertThat()
                .log().ifValidationFails()
                .statusCode(IsIn(listOf(201, 422)))
                .extract()

        if (response.statusCode() == 201) {
            val responseBodyString = getGetResponseBodyAsString(response)
            assertTrue(JsonSchemaValidator.matchesJsonSchemaInClasspath(VALID_BINDING.path)
                    .matches(responseBodyString), "\nExpected a valid binding response but was:\n$responseBodyString")
        } else {
            val responseBodyString = getGetResponseBodyAsString(response)
            assertTrue(JsonSchemaValidator.matchesJsonSchemaInClasspath(ERR_ASYNC_REQUIRED.path)
                    .matches(responseBodyString), "\nExpected OSB error code async required but was:\n$responseBodyString")
        }
    }

    fun runPutBindingRequestAsync(
            requestBody: RequestBody,
            instanceId: String,
            bindingId: String,
            vararg expectedStatusCodes: Int,
            expectedResponseBody: ResponseBodyType
    ): ExtractableResponse<Response> {
        if (configuration.apiVersion >= 2.15 && configuration.useRequestIdentity) {
            useRequestIdentity("OSB-Checker-PUT-binding-${UUID.randomUUID()}")
        }

        val response = RestAssured.with()
                .log().ifValidationFails()
                .headers(validRequestHeaders)
                .contentType(ContentType.JSON)
                .body(requestBody)
                .param("accepts_incomplete", configuration.apiVersion > 2.13)
                .put(SERVICE_INSTANCE_PATH + instanceId + SERVICE_BINDING_PATH + bindingId)
                .then()
                .log().ifValidationFails()
                .assertThat()
                .headers(expectedResponseHeaders)
                .statusCode(IsIn(expectedStatusCodes.asList()))
                .extract()

        if (response.statusCode() == 200) {
            val responseBodyString = response.body().jsonPath().prettify()
            assertTrue(JsonSchemaValidator.matchesJsonSchemaInClasspath(PATH_TO_BINDING).matches(responseBodyString),
                    "\nExpected a valid binding ResponseBody, but was:\n$responseBodyString")
        }

        return response
    }

    fun polling(
            instanceId: String,
            bindingId: String,
            expectedFinalStatusCode: Int,
            operationData: String?,
            maxPollingDuration: Int
    ): LastOperationResponse.State {
        val latestAcceptablePollingInstant = Instant.now().plusSeconds(maxPollingDuration.toLong())

        return waitForFinish(
                path = SERVICE_INSTANCE_PATH + instanceId + SERVICE_BINDING_PATH + bindingId + LAST_OPERATION,
                expectedFinalStatusCode = expectedFinalStatusCode,
                operationData = operationData,
                latestAcceptablePollingInstant = latestAcceptablePollingInstant
        )
    }

    fun runDeleteBindingRequestAsync(
            serviceId: String?,
            planId: String?,
            instanceId: String,
            bindingId: String,
            vararg expectedStatusCodes: Int
    ): ExtractableResponse<Response> {
        var path = SERVICE_INSTANCE_PATH + instanceId + SERVICE_BINDING_PATH + bindingId
        path = serviceId?.let { "$path?service_id=$serviceId" } ?: path

        if (configuration.apiVersion >= 2.15 && configuration.useRequestIdentity) {
            useRequestIdentity("OSB-Checker-DELETE-binding-${UUID.randomUUID()}")
        }

        return RestAssured.with()
                .log().ifValidationFails()
                .headers(validRequestHeaders)
                .contentType(ContentType.JSON)
                .param("service_id", serviceId)
                .param("plan_id", planId)
                .param("accepts_incomplete", configuration.apiVersion > 2.13)
                .delete(path)
                .then()
                .log().ifValidationFails()
                .assertThat()
                .headers(expectedResponseHeaders)
                .statusCode(IsIn(expectedStatusCodes.asList()))
                .extract()
    }

    fun runDeleteBindingRequestSync(instanceId: String, bindingId: String, serviceId: String?, planId: String?) {
        var path = SERVICE_INSTANCE_PATH + instanceId + SERVICE_BINDING_PATH + bindingId
        path = serviceId?.let { "$path?service_id=$serviceId" } ?: path
        path = planId?.let { "$path&plan_id=$planId" } ?: path

        if (configuration.apiVersion >= 2.15 && configuration.useRequestIdentity) {
            useRequestIdentity("OSB-Checker-DELETE-binding-${UUID.randomUUID()}")
        }

        val response = RestAssured.with()
                .log().ifValidationFails()
                .headers(validRequestHeaders)
                .contentType(ContentType.JSON)
                .delete(path)
                .then()
                .log().ifValidationFails()
                .headers(expectedResponseHeaders)
                .statusCode(IsIn(listOf(200, 422)))
                .extract()

        if (response.statusCode() != 200) {
            val responseBodyString = getGetResponseBodyAsString(response)
            assertTrue(JsonSchemaValidator.matchesJsonSchemaInClasspath(ERR_ASYNC_REQUIRED.path)
                    .matches(responseBodyString), "Expected OSB error code async required but was:\n$responseBodyString")
        }
    }

    fun putWithoutHeader() {
        RestAssured.with()
                .log().ifValidationFails()
                .header(Header("Authorization", configuration.correctToken))
                .put(SERVICE_INSTANCE_PATH + Configuration.notAnId + SERVICE_BINDING_PATH + Configuration.notAnId)
                .then()
                .log().ifValidationFails()
                .assertThat()
                .statusCode(412)
    }

    fun deleteWithoutHeader() {
        RestAssured.with()
                .log().ifValidationFails()
                .header(Header("Authorization", configuration.correctToken))
                .put(SERVICE_INSTANCE_PATH + Configuration.notAnId + SERVICE_BINDING_PATH + Configuration.notAnId)
                .then()
                .log().ifValidationFails()
                .assertThat()
                .statusCode(412)
    }

    fun putNoAuth() {
        RestAssured.with()
                .log().ifValidationFails()
                .header(Header("X-Broker-API-Version", "${configuration.apiVersion}"))
                .put(SERVICE_INSTANCE_PATH + Configuration.notAnId + SERVICE_BINDING_PATH + Configuration.notAnId)
                .then()
                .log().ifValidationFails()
                .assertThat()
                .statusCode(401)
    }

    fun putWrongUser() {
        RestAssured.with()
                .header(Header("Authorization", configuration.wrongUserToken))
                .log().ifValidationFails()
                .header(Header("X-Broker-API-Version", "$configuration.apiVersion"))
                .put(SERVICE_INSTANCE_PATH + Configuration.notAnId + SERVICE_BINDING_PATH + Configuration.notAnId)
                .then()
                .log().ifValidationFails()
                .assertThat()
                .statusCode(401)
    }

    fun putWrongPassword() {
        RestAssured.with()
                .log().ifValidationFails()
                .header(Header("Authorization", configuration.wrongPasswordToken))
                .header(Header("X-Broker-API-Version", "${configuration.apiVersion}"))
                .put(SERVICE_INSTANCE_PATH + Configuration.notAnId + SERVICE_BINDING_PATH + Configuration.notAnId)
                .then()
                .log().ifValidationFails()
                .assertThat()
                .statusCode(401)
    }

    fun deleteNoAuth() {
        RestAssured.with()
                .log().ifValidationFails()
                .header(Header("X-Broker-API-Version", "${configuration.apiVersion}"))
                .delete(SERVICE_INSTANCE_PATH + Configuration.notAnId + SERVICE_BINDING_PATH + Configuration.notAnId)
                .then()
                .log().ifValidationFails()
                .assertThat()
                .statusCode(401)
    }

    fun deleteWrongUser() {
        RestAssured.with()
                .log().ifValidationFails()
                .header(Header("Authorization", configuration.wrongUserToken))
                .header(Header("X-Broker-API-Version", "${configuration.apiVersion}"))
                .delete(SERVICE_INSTANCE_PATH + Configuration.notAnId + SERVICE_BINDING_PATH + Configuration.notAnId)
                .then()
                .log().ifValidationFails()
                .assertThat()
                .statusCode(401)
    }

    fun deleteWrongPassword() {
        RestAssured.with()
                .log().ifValidationFails()
                .header(Header("Authorization", configuration.wrongPasswordToken))
                .header(Header("X-Broker-API-Version", "${configuration.apiVersion}"))
                .delete(SERVICE_INSTANCE_PATH + Configuration.notAnId + SERVICE_BINDING_PATH + Configuration.notAnId)
                .then()
                .log().ifValidationFails()
                .assertThat()
                .statusCode(401)
    }

    companion object {
        const val SERVICE_BINDING_PATH = "/service_bindings/"
        private const val PATH_TO_BINDING = "binding-response-schema.json"
    }
}