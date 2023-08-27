Feature: latest rounds
  In order to easily see the latest round
  As a viewer
  I want to get the latest round ID

  Scenario: Create a new round and check the latest
    When I create a new round with 2 players
    And I get the latest round ID
    Then the latest round ID matches

  Scenario: Create multiple new rounds and check the latest
    When I create a new round with 2 players
    And I create a new round with 2 players
    And I get the latest round ID
    Then the latest round ID matches
