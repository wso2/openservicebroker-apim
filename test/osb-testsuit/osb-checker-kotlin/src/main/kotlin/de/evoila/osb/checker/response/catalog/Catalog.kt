package de.evoila.osb.checker.response.catalog

import com.fasterxml.jackson.annotation.JsonIgnoreProperties

@JsonIgnoreProperties(ignoreUnknown = true)
data class Catalog(
    val services: List<Service>
)
