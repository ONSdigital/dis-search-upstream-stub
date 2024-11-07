Feature: Resources

  Scenario: Posting and checking a response
    When I GET "/resources"
    Then I should receive a list of resources
    And the HTTP status code should be "200"

  Scenario: Valid offset and limit get response status 200
    When I GET "/resources?offset=3&limit=21"
    Then the HTTP status code should be "200"

  Scenario: Invalid limit gets bad request error
    When I GET "/resources?limit=badger"
    Then the HTTP status code should be "400"
    And I should receive the following response:
            """
            invalid limit query parameter
            """

  Scenario: Invalid offset gets bad request error
    When I GET "/resources?offset=turtle"
    Then the HTTP status code should be "400"
    And I should receive the following response:
            """
            invalid offset query parameter
            """

  Scenario: Invalid limit gets bad request error
    When I GET "/resources?offset=4&limit=2000"
    Then the HTTP status code should be "400"
    And I should receive the following response:
            """
            limit query parameter is larger than the maximum allowed
            """
    