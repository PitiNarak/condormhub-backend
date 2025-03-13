Feature: Add Dorms as a Lessor
  As a lessor
  I want to add multiple properties to the system
  So that I can manage them efficiently

  Scenario: Successfully add a new dorm property
    Given the user is logged in as a lessor
    And the user prepares valid dorm data
    When the user submits the dorm data
    Then the dorm should be created successfully

  Scenario: Reject invalid dorm data
    Given the user is logged in as a lessor
    And the user prepares invalid dorm data
    When the user submits the dorm data
    Then the dorm validation should fail