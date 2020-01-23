package de.evoila.osb.checker.tests

import org.junit.jupiter.api.DisplayName
import org.junit.jupiter.api.Test

@DisplayName(value = "Catalog test")
class CatalogJUnit5 : TestBase() {

  @Test
  @DisplayName(value = "Verify catalog schema")
  fun validateCatalog() {
    catalogRequestRunner.correctRequestAndValidateResponse()
  }
}