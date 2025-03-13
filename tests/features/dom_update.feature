Feature: Update Dorms As A Lessor
  As a lessor
  I want to update my property details
  So that I can keep the information accurate and up to date

  Scenario: Successfully update property as owner
    Given I am a registered lessor
    And I have a property listed
    And I am the owner of the property
    When I update property details
    And I submit the changes
    Then the updated details should be saved
    And I should see the updated information in my listings

  Scenario: Attempt to update property as non-owner
    Given I am a registered lessor
    And I have a property listed
    And I am not the owner of the property
    When I update property details
    And I submit the changes
    Then I should see an error message "You do not have permission to update this dorm"
    And my property should not be updated