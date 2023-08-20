Feature: send scores
  In order to track the current round
  As a player
  I want to send scores from my round hole-by-hole

  Background:
    Given a completely fresh environment

  Scenario: Create a round
    When I create a round
    Then I see the round with no entered scores
