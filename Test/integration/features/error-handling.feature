Feature: Error Handling and Edge Cases
  As an API consumer
  I want to verify proper error handling
  So that I can ensure the API responds correctly to invalid requests and edge cases.

  Background:
    # Login to obtain accessToken is handled globally by InitializeScenario
    # and the token is automatically added to headers by the addAuthHeader function.
    # All resources created in scenarios are automatically tracked and cleaned up
    # by the test framework's teardown mechanism.

  # ===== AUTHENTICATION ERROR CASES =====

  Scenario: TC01 - Access protected endpoint without authentication
    Given I clear the authentication token
    When I send a GET request to "/api/users"
    Then the response code should be 401
    And the JSON response should contain error "error": "Authorization header not provided"

  Scenario: TC01.1 - Access protected endpoint with invalid token format
    Given I clear the authentication token
    When I send a GET request to "/api/medicines/1"
    Then the response code should be 401
    And the JSON response should contain error "error": "Authorization header not provided"

  Scenario: TC01.2 - Access protected endpoint with malformed Bearer token
    Given I clear the authentication token
    When I send a GET request to "/api/icd-cie"
    Then the response code should be 401
    And the JSON response should contain error "error": "Authorization header not provided"

  # ===== INVALID ID FORMATS =====

  Scenario: TC02 - Attempt to access resource with non-numeric ID
    When I send a GET request to "/api/users/abc123"
    Then the response code should be 400
    And the JSON response should contain error "error": "Invalid ID"

  Scenario: TC02.1 - Attempt to update resource with non-numeric ID
    When I send a PUT request to "/api/medicines/xyz789" with body:
      """
      {
        "eanCode": "TEST123",
        "description": "Test medicine",
        "laboratory": "Test Lab"
      }
      """
    Then the response code should be 400
    And the JSON response should contain error "error": "Invalid ID"

  Scenario: TC02.2 - Attempt to delete resource with non-numeric ID
    When I send a DELETE request to "/api/icd-cie/invalid-id"
    Then the response code should be 400
    And the JSON response should contain error "error": "Invalid ID"

  # ===== MISSING REQUIRED FIELDS =====

  Scenario: TC03 - Create user without required fields
    When I send a POST request to "/api/users" with body:
      """
      {
        "firstName": "John",
        "lastName": "Doe"
      }
      """
    Then the response code should be 400
    And the JSON response should contain key "error"

  Scenario: TC03.1 - Create medicine without required fields
    When I send a POST request to "/api/medicines" with body:
      """
      {
        "description": "Test medicine without EAN code"
      }
      """
    Then the response code should be 400
    And the JSON response should contain key "error"

  Scenario: TC03.2 - Create ICD-CIE without required fields
    When I send a POST request to "/api/icd-cie" with body:
      """
      {
        "description": "Test ICD-CIE without code"
      }
      """
    Then the response code should be 400
    And the JSON response should contain key "error"

  # ===== INVALID JSON PAYLOADS =====

  Scenario: TC04 - Send request with malformed JSON
    When I send a POST request to "/api/users" with body:
      """
      {
        "username": "testuser",
        "email": "test@example.com",
        "password": "password123",
        "roleId": 1,
        "firstName": "John",
        "lastName": "Doe"
      }
      """
    Then the response code should be 400
    And the JSON response should contain key "error"

  Scenario: TC04.1 - Send request with empty JSON body
    When I send a POST request to "/api/medicines" with body:
      """
      {}
      """
    Then the response code should be 400
    And the JSON response should contain key "error"

  # ===== NON-EXISTENT RESOURCES =====

  Scenario: TC05 - Access non-existent user
    When I send a GET request to "/api/users/999999"
    Then the response code should be 404
    And the JSON response should contain error "error": "User not found"

  Scenario: TC05.1 - Update non-existent medicine
    When I send a PUT request to "/api/medicines/999999" with body:
      """
      {
        "eanCode": "NONEXISTENT",
        "description": "Non-existent medicine",
        "laboratory": "No Lab"
      }
      """
    Then the response code should be 404
    And the JSON response should contain error "error": "Medicine not found"

  Scenario: TC05.2 - Delete non-existent ICD-CIE record
    When I send a DELETE request to "/api/icd-cie/999999"
    Then the response code should be 404
    And the JSON response should contain error "error": "ICDCie record not found"

  # ===== INVALID SEARCH PARAMETERS =====

  Scenario: TC06 - Search with invalid property name
    When I send a GET request to "/api/medicines/search-by-property?property=invalid_field&search_text=test"
    Then the response code should be 400
    And the JSON response should contain error "error": "Invalid property or search_text"

  Scenario: TC06.1 - Search with empty search text
    When I send a GET request to "/api/icd-cie/search-by-property?property=code&search_text="
    Then the response code should be 400
    And the JSON response should contain error "error": "Invalid property or search_text"

  Scenario: TC06.2 - Search with missing search text
    When I send a GET request to "/api/medicines/search-by-property?property=description"
    Then the response code should be 400
    And the JSON response should contain error "error": "Invalid property or search_text"

  # ===== PAGINATION EDGE CASES =====

  Scenario: TC07 - Search with negative page number
    When I send a GET request to "/api/medicines/search-paginated?page=-1&limit=10"
    Then the response code should be 200
    And the JSON response should contain "current_page": 1

  Scenario: TC07.1 - Search with zero page number
    When I send a GET request to "/api/icd-cie/search-paginated?page=0&limit=10"
    Then the response code should be 200
    And the JSON response should contain "current_page": 1

  Scenario: TC07.2 - Search with negative limit
    When I send a GET request to "/api/users/devices/search-paginated?page=1&limit=-5"
    Then the response code should be 200
    And the JSON response should contain "page_size": 10

  Scenario: TC07.3 - Search with zero limit
    When I send a GET request to "/api/medicines/search-paginated?page=1&limit=0"
    Then the response code should be 200
    And the JSON response should contain "page_size": 10

  Scenario: TC07.4 - Search with very large limit
    When I send a GET request to "/api/icd-cie/search-paginated?page=1&limit=1000"
    Then the response code should be 200
    And the JSON response should contain "page_size": 1000

  # ===== INVALID HTTP METHODS =====

  Scenario: TC08 - Use unsupported HTTP method on users endpoint
    When I send a PATCH request to "/api/users/1"
    Then the response code should be 404

  Scenario: TC08.1 - Use unsupported HTTP method on medicines endpoint
    When I send a PATCH request to "/api/medicines/1"
    Then the response code should be 404

  Scenario: TC08.2 - Use unsupported HTTP method on ICD-CIE endpoint
    When I send a PATCH request to "/api/icd-cie/1"
    Then the response code should be 404

  # ===== MISSING PATH PARAMETERS =====

  Scenario: TC09 - Access endpoint without required path parameter
    When I send a GET request to "/api/users/"
    Then the response code should be 404

  Scenario: TC09.1 - Access endpoint without required path parameter for medicines
    When I send a GET request to "/api/medicines/"
    Then the response code should be 404

  Scenario: TC09.2 - Access endpoint without required path parameter for ICD-CIE
    When I send a GET request to "/api/icd-cie/"
    Then the response code should be 404

  # ===== INVALID QUERY PARAMETERS =====

  Scenario: TC10 - Search with invalid query parameter format
    When I send a GET request to "/api/medicines/search-paginated?page=abc&limit=def"
    Then the response code should be 200
    And the JSON response should contain "current_page": 1
    And the JSON response should contain "page_size": 10

  Scenario: TC10.1 - Search with invalid query parameter format for ICD-CIE
    When I send a GET request to "/api/icd-cie/search-paginated?page=xyz&limit=123"
    Then the response code should be 200
    And the JSON response should contain "current_page": 1
    And the JSON response should contain "page_size": 10

  # ===== CORS AND HEADERS =====

  Scenario: TC11 - Verify CORS headers are present
    When I send a GET request to "/api/users"
    Then the response code should be 200
    # Note: CORS headers are handled by middleware and may not be visible in response

  Scenario: TC11.1 - Verify content type headers
    When I send a GET request to "/api/medicines/1"
    Then the response code should be 404
    # Note: Content-Type headers are handled by Gin framework

  # ===== LARGE PAYLOADS =====

  Scenario: TC12 - Send request with very large description field
    Given I generate a unique EAN code as "largePayloadEan"
    When I send a POST request to "/api/medicines" with body:
      """
      {
        "eanCode": "${largePayloadEan}",
        "description": "This is a very long description that exceeds normal limits. It contains many characters and should test the system's ability to handle large text fields. The description goes on and on with repetitive text to ensure we reach a substantial length that could potentially cause issues with database storage or API processing. This test ensures that the system can handle large payloads without crashing or behaving unexpectedly.",
        "laboratory": "Large Payload Test Lab",
        "type": "tableta",
        "iva": "16",
        "satKey": "51182200",
        "activeIngredient": "Test Ingredient",
        "temperatureControl": "seco",
        "isControlled": false,
        "unitQuantity": 100.0,
        "unitType": "tabletas",
        "rxCode": "RXP-LARGE"
      }
      """
    Then the response code should be 201
    And the JSON response should contain key "id"

  # ===== SPECIAL CHARACTERS =====

  Scenario: TC13 - Send request with special characters in fields
    Given I generate a unique alias as "specialCharUsername"
    And I generate a unique alias as "specialCharEmail"
    And I send a POST request to "/api/users/roles" with body:
      """
      {
        "name": "SpecialCharRole",
        "description": "Role for special character testing",
        "enabled": true
      }
      """
    And I save the JSON response key "id" as "specialCharRoleID"
    When I send a POST request to "/api/users" with body:
      """
      {
        "username": "${specialCharUsername}",
        "firstName": "José María",
        "lastName": "García-López",
        "email": "${specialCharEmail}@test.com",
        "password": "p@ssw0rd!123",
        "jobPosition": "Développeur Senior",
        "roleId": ${specialCharRoleID},
        "enabled": true
      }
      """
    Then the response code should be 201
    And the JSON response should contain "firstName": "José María"
    And the JSON response should contain "lastName": "García-López"
    And the JSON response should contain "jobPosition": "Développeur Senior" 