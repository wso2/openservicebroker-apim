package de.evoila.osb.checker.response.catalog

import com.fasterxml.jackson.annotation.JsonIgnoreProperties
import com.fasterxml.jackson.annotation.JsonProperty

@JsonIgnoreProperties(ignoreUnknown = true)
data class DashboardClient(
    val id: String?,
    @JsonProperty("redirect_uri")
    val redirectUri: String?,
    val secret: String?
)