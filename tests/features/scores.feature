Feature: send scores
  In order to track the current round
  As a player
  I want to send scores from my round hole-by-hole

  Scenario: Create a new round and look at it
    When I create a new round
    And I view the round
    Then the round is valid but empty
