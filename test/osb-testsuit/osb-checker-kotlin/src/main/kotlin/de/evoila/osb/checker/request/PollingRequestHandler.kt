package de.evoila.osb.checker.request

import de.evoila.osb.checker.config.Configuration
import de.evoila.osb.checker.response.operations.LastOperationResponse
import de.evoila.osb.checker.response.operations.LastOperationResponse.State.FAILED
import io.restassured.RestAssured
import io.restassured.http.ContentType
import io.restassured.module.jsv.JsonSchemaValidator
import org.hamcrest.collection.IsIn
import java.time.Instant
import java.util.*
import kotlin.test.assertTrue

abstract class PollingRequestHandler(
        configuration: Configuration
) : RequestHandler(configuration) {

    fun waitForFinish(path: String,
                      expectedFinalStatusCode: Int,
                      operationData: String?,
                      latestAcceptablePollingInstant: Instant
    ): LastOperationResponse.State {
        if (configuration.apiVersion >= 2.15 && configuration.useRequestIdentity) {
            useRequestIdentity("OSB-Checker-GET-last-operation-${UUID.randomUUID()}")
        }

        val request = RestAssured.with()
                .log().ifValidationFails()
                .headers(validRequestHeaders)
                .contentType(ContentType.JSON)

        operationData?.let { request.queryParam("operation", it) }

        val response = request.get(path)
                .then()
                .log().ifValidationFails()
                .assertThat()
                .headers(expectedResponseHeaders)
                .statusCode(IsIn(listOf(expectedFinalStatusCode, 200)))
                .extract()
                .response()

        if (response.statusCode == 410) {
            return LastOperationResponse.State.GONE
        }

        val responseBody = response.jsonPath()
                .getObject("", LastOperationResponse::class.java)

        if (!Instant.now().isBefore(latestAcceptablePollingInstant)) {
            assertTrue(responseBody.state != FAILED,
                    "\nWhen it takes more time than defined in maximumPollingDuration the Operation needs be defined as FAILED but was ${responseBody.state}")
            assertTrue(false, "\nInstance creation took longer than it should!!")
        }

        val responseBodyString = response.jsonPath().prettify()
        assert(JsonSchemaValidator.matchesJsonSchemaInClasspath("polling-response-schema.json")
                .matches(responseBodyString)) { "Expected a valid polling result body but was $responseBodyString" }

        return if (responseBody.state == LastOperationResponse.State.IN_PROGRESS) {
            Thread.sleep(10000)
            waitForFinish(path, expectedFinalStatusCode, operationData, latestAcceptablePollingInstant)
        } else {
            responseBody.state
        }
    }

    companion object {
        const val ACCEPTS_INCOMPLETE = "?accepts_incomplete=true"
        const val LAST_OPERATION = "/last_operation"
        const val SERVICE_INSTANCE_PATH = "/v2/service_instances/"
    }
}
