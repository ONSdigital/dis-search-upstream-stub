Feature: Resources

  Scenario: Posting and checking a response
    When I GET "/resources"
    Then I should receive a list of resources

    