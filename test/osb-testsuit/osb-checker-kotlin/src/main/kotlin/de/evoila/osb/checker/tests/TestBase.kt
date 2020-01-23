package de.evoila.osb.checker.tests

import de.evoila.osb.checker.Application
import de.evoila.osb.checker.config.Configuration
import de.evoila.osb.checker.request.CatalogRequestRunner
import de.evoila.osb.checker.request.ResponseBodyType
import de.evoila.osb.checker.request.bodies.RequestBody
import org.junit.jupiter.api.extension.ExtendWith
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.test.context.junit.jupiter.SpringExtension

@AutoConfigureMockMvc
@ExtendWith(SpringExtension::class)
@SpringBootTest(classes = [Application::class])
abstract class TestBase {

    @Autowired
    lateinit var catalogRequestRunner: CatalogRequestRunner

    @Autowired
    lateinit var configuration: Configuration
}

data class TestCase<out T : RequestBody>(
        val requestBody: T,
        val message: String,
        val responseBodyType: ResponseBodyType,
        val statusCode : Int
)