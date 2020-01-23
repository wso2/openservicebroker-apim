package de.evoila.osb.checker.response.catalog

import com.fasterxml.jackson.annotation.JsonIgnoreProperties

@JsonIgnoreProperties(ignoreUnknown = true)
data class ServiceMetadata(
    val displayName: String?,
    val listing: String?,
    val provider: Provider?
)