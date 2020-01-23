package de.evoila.osb.checker.request

enum class ResponseBodyType(val path: String) {
    VALID_PROVISION("provision-response-schema.json"),
    VALID_FETCH_INSTANCE("fetch-instance-response-schema.json"),
    VALID_BINDING("binding-response-schema.json"),
    VALID_FETCH_BINDING("binding-response-schema.json"),
    ERR_MAINTENANCE_INFO("service-broker-maintenance-info-error.json"),
    ERR_ASYNC_REQUIRED("service-broker-async-required.json"),
    ERR_CONCURRENCY("service-broker-concurrency-response.json"),
    ERR_REQUIRES_APP("service-broker-error-response.json"),
    ERR("service-broker-error-response.json"),
    NO_SCHEMA("no-schema.json")
}
