Feature: send scores
  In order to track the current round
  As a player
  I want to send scores from my round hole-by-hole

  Scenario: Create a new round and look at it
    When I create a new round with 2 players
    And I view the round
    Then the round is valid but empty

  Scenario: Create a new round and send one score
    When I create a new round with 3 players
    And player 2 scores a 3 on hole 5
    And I view the round
    Then the score for player 2 on hole 5 is 3

  Scenario: Create a new round and send multiple scores
    When I create a new round with 3 players
    And player 2 scores a 3 on hole 5
    And player 2 scores a 4 on hole 6
    And player 1 scores a 8 on hole 2
    And I view the round
    Then the score for player 2 on hole 5 is 3
    And the score for player 2 on hole 6 is 4
    And the score for player 1 on hole 2 is 8
