package de.evoila.osb.checker.response.operations

import com.fasterxml.jackson.annotation.JsonIgnoreProperties
import com.fasterxml.jackson.annotation.JsonProperty

@JsonIgnoreProperties
data class ProvisionResponse(
        @JsonProperty("dashboard_url")
        val dashboardUrl: String?
)