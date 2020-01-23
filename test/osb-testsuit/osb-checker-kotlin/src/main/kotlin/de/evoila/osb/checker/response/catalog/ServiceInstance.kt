package de.evoila.osb.checker.response.catalog

import com.fasterxml.jackson.annotation.JsonIgnoreProperties
import com.fasterxml.jackson.annotation.JsonProperty

@JsonIgnoreProperties(ignoreUnknown = true)
data class ServiceInstance(
    @JsonProperty("service_id")
    val serviceId: String?,
    @JsonProperty("plan_id")
    val planId: String?,
    @JsonProperty("dashboard_url")
    var dashboardUrl: String?,
    val parameters: HashMap<String, Any> = hashMapOf()
)