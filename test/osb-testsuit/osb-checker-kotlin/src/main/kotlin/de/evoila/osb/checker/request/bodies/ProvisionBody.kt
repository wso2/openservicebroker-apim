package de.evoila.osb.checker.request.bodies

import com.fasterxml.jackson.annotation.JsonInclude
import de.evoila.osb.checker.response.catalog.MaintenanceInfo
import de.evoila.osb.checker.response.catalog.Plan
import de.evoila.osb.checker.response.catalog.Service
import java.util.*

abstract class ProvisionBody : RequestBody {

  data class ValidProvisioning(
      @JsonInclude(JsonInclude.Include.NON_EMPTY)
      var service_id: String,
      @JsonInclude(JsonInclude.Include.NON_EMPTY)
      var plan_id: String,
      @JsonInclude(JsonInclude.Include.NON_EMPTY)
      var organization_guid: String = UUID.randomUUID().toString(),
      @JsonInclude(JsonInclude.Include.NON_EMPTY)
      var space_guid: String = UUID.randomUUID().toString(),
      @JsonInclude(JsonInclude.Include.NON_NULL)
      var parameters: Map<String, Any>? = null,
      @JsonInclude(JsonInclude.Include.NON_NULL)
      var maintenance_info: MaintenanceInfo? = null
  ) : ProvisionBody() {

    constructor(service: Service, plan: Plan) : this(
        service_id = service.id,
        plan_id = plan.id
    )

    constructor(service: Service, plan: Plan, maintenance_info: MaintenanceInfo) : this(
        service_id = service.id,
        plan_id = plan.id,
        maintenance_info = maintenance_info
    )
  }

  data class NoPlanFieldProvisioning(
      var service_id: String?,
      var organization_guid: String? = UUID.randomUUID().toString(),
      var space_guid: String? = UUID.randomUUID().toString()
  ) : ProvisionBody() {

    constructor(service: Service) : this(
        service_id = service.id
    )
  }

  data class NoServiceFieldProvisioning(
      var plan_id: String?,
      var organization_guid: String? = UUID.randomUUID().toString(),
      var space_guid: String? = UUID.randomUUID().toString()
  ) : ProvisionBody() {

    constructor(plan: Plan) : this(
        plan_id = plan.id
    )
  }

  data class NoOrgFieldProvisioning(
      var service_id: String?,
      var plan_id: String?,
      var space_guid: String? = UUID.randomUUID().toString()
  ) : ProvisionBody() {

    constructor(service: Service, plan: Plan) : this(
        service_id = service.id,
        plan_id = plan.id
    )
  }

  data class NoSpaceFieldProvisioning(
      var service_id: String?,
      var plan_id: String?,
      var organization_guid: String? = UUID.randomUUID().toString()
  ) : ProvisionBody() {

    constructor(service: Service, plan: Plan) : this(
        service_id = service.id,
        plan_id = plan.id
    )
  }
}