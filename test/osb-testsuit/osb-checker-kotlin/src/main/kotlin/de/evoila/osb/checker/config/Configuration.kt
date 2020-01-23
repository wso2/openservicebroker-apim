package de.evoila.osb.checker.config

import de.evoila.osb.checker.response.catalog.Catalog
import java.util.*
import org.springframework.boot.context.properties.ConfigurationProperties
import org.springframework.stereotype.Component
import kotlin.collections.HashMap

@Component
@ConfigurationProperties(prefix = "config")
class Configuration {

    lateinit var url: String
    var port: Int = 80
    var apiVersion: Double = 0.0
    lateinit var user: String
    lateinit var password: String
    lateinit var correctToken: String
    lateinit var wrongUserToken: String
    lateinit var wrongPasswordToken: String
    var originatingIdentity: OriginatingIdentity? = null
    var useRequestIdentity: Boolean = false
    var skipTLSVerification: Boolean = false
    var testDashboard: Boolean = false
    var usingAppGuid: Boolean = true
    val provisionParameters: HashMap<String, HashMap<String, Any>> = hashMapOf()
    val bindingParameters: HashMap<String, HashMap<String, Any>> = hashMapOf()
    var services = mutableListOf<CustomService>()

    /*
     * This Method filters an provided catalog by the service and plan ids set in the application.yml
     * If no services are set the full catalog gets returned.
     * If only service id's are defined, all plans in the provided service are returned.
     */
    fun initCustomCatalog(fullCatalog: Catalog): Catalog {
        return if (services.isNotEmpty()) {
            fullCatalog.copy(
                    /*
                     * Creates a service list from the catalog only with the services which ids are defined in the
                     * customServices field from the application.yml.
                     */
                    services = fullCatalog.services.filter { service ->
                        services.firstOrNull { service.id == it.id }?.let { true } ?: false
                    }.map { filteredService ->
                        val customService = services.first { it.id == filteredService.id }
                        if (customService.plans.isNotEmpty()) {
                            filteredService.copy(
                                    /*
                                     * Filters by same predicate as before just for the plans on the filtered services.
                                     */
                                    plans = filteredService.plans.filter { plan ->
                                        customService.plans.firstOrNull { customPlan -> customPlan.id == plan.id }?.let { true }
                                                ?: false
                                    }
                            )
                        } else filteredService
                    }
            )
        } else fullCatalog
    }

    class OriginatingIdentity {
        var platform: String = ""
        var value: Map<String, Any> = hashMapOf()
    }

    class CustomService {
        lateinit var id: String
        var plans = mutableListOf<CustomPlan>()

        class CustomPlan {
            lateinit var id: String
        }
    }

    companion object {
        val notAnId = UUID.randomUUID().toString()
    }
}