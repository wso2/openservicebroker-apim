package de.evoila.osb.checker.request

import com.fasterxml.jackson.module.kotlin.jacksonObjectMapper
import de.evoila.osb.checker.config.Configuration
import de.evoila.osb.checker.config.Configuration.*
import io.restassured.response.ExtractableResponse
import io.restassured.response.Response
import java.util.*
import kotlin.test.assertTrue


abstract class RequestHandler(val configuration: Configuration) {

    val validRequestHeaders: MutableMap<String, Any> = mutableMapOf()
    val expectedResponseHeaders: MutableMap<String, Any> = mutableMapOf()

    init {
        validRequestHeaders["X-Broker-API-Version"] = configuration.apiVersion
        validRequestHeaders["Authorization"] = configuration.correctToken

        configuration.originatingIdentity?.let {
            validRequestHeaders[ORIGINATING_IDENTITY_KEY] = encodeOriginatingIdentity(it)
        }
    }

    fun useRequestIdentity(value: String) {
        validRequestHeaders[REQUEST_IDENTITY_KEY] = value
        expectedResponseHeaders[REQUEST_IDENTITY_KEY] = value
    }

    fun getGetResponseBodyAsString(response: ExtractableResponse<Response>): String = try {
        response.jsonPath().prettify()
    } catch (ex: Exception) {
        assertTrue(false, "Expected a response Body, but none was found.")
        ""
    }

    private fun encodeOriginatingIdentity(originatingIdentity: OriginatingIdentity): String {
        val mapper = jacksonObjectMapper()
        val jsonValue = mapper.writeValueAsString(originatingIdentity.value)

        return "${originatingIdentity.platform} ${Base64.getEncoder().encodeToString(jsonValue.toByteArray())}"
    }

    companion object {
        const val ORIGINATING_IDENTITY_KEY = "X-Broker-API-Originating-Identity"
        const val REQUEST_IDENTITY_KEY = "X-Broker-API-Request-Identity"
    }
}
