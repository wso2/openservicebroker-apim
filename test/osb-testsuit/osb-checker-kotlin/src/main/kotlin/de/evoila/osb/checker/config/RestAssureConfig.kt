package de.evoila.osb.checker.config

import io.restassured.RestAssured
import org.springframework.stereotype.Service
import java.util.*
import kotlin.test.assertTrue

@Service
class RestAssureConfig(
        configuration: Configuration
) {
    init {
        assertTrue(configuration.apiVersion in supportedApiVersions, noSupportedApiVersion)
        configuration.correctToken = encode(configuration.user, configuration.password)
        configuration.wrongUserToken = encode(UUID.randomUUID().toString(), configuration.password)
        configuration.wrongPasswordToken = encode(configuration.user, UUID.randomUUID().toString())

        RestAssured.baseURI = configuration.url
        RestAssured.port = configuration.port
        if (configuration.skipTLSVerification) {
            RestAssured.useRelaxedHTTPSValidation()
        }
    }

    private fun encode(user: String, password: String): String =
            "Basic ${Base64.getEncoder().encodeToString("$user:$password".toByteArray())}"

    companion object {
        private val supportedApiVersions = listOf(2.13, 2.14, 2.15)
        private val noSupportedApiVersion = "You entered a not supported Api Version. Please use one of the following:" +
                " $supportedApiVersions"
    }
}