package de.evoila.osb.checker.response.catalog

import com.fasterxml.jackson.annotation.JsonIgnoreProperties
import com.fasterxml.jackson.annotation.JsonProperty

@JsonIgnoreProperties(ignoreUnknown = true)
data class Plan(
    val id: String,
    val name: String,
    val description: String,
    @JsonProperty("plan_updatable")
    val planUpdatable: Boolean?,
    val bindable: Boolean?,
    val metadata: PlanMetadata?,
    @JsonProperty("maximum_polling_duration")
    val maximumPollingDuration: Int = 86400,
    val maintenanceInfo: MaintenanceInfo?
) {

  constructor(id: String, name: String, description: String, bindable: Boolean) : this(
      id = id,
      name = name,
      description = description,
      planUpdatable = false,
      bindable = bindable,
      metadata = null,
      maintenanceInfo = null
  )

}