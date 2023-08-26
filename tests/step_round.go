package tests

import (
	"fmt"
	"net/http"

	"github.com/Evertras/live-leaderboards/pkg/api"
)

func (t *testContext) iCreateANewRoundWithPlayers(numPlayers int) error {
	players := make([]api.PlayerData, numPlayers)

	for i := 0; i < numPlayers; i++ {
		players[i] = api.PlayerData{
			Name: fmt.Sprintf("Test Player %d", i),
		}
	}

	t.roundRequest = &api.RoundRequest{
		Course: api.Course{
			Holes: []api.Hole{
				{
					DistanceYards: ptr(370),
					Hole:          1,
					Par:           4,
					StrokeIndex:   ptr(17),
				},
			},
			Name: "Test Round",
			Tees: ptr("White"),
		},
		Players: players,
		Title:   ptr("Test Round"),
	}

	createdRound, err := t.client.CreateRound(t.execCtx, *t.roundRequest)

	if err != nil {
		return fmt.Errorf("failed to create round: %w", err)
	}

	if createdRound.Id.String() == "" {
		return fmt.Errorf("id was empty")
	}

	t.createdRoundID = createdRound.Id

	return nil
}

func (t *testContext) iViewTheRound() error {
	round, err := t.client.GetRound(t.execCtx, t.createdRoundID)

	if err != nil {
		return fmt.Errorf("failed to get round: %w", err)
	}

	t.returnedRound = round

	return nil
}

func (t *testContext) theRoundIsValidButEmpty() error {
	if t.returnedRound == nil {
		return fmt.Errorf("no round returned")
	}

	if t.createdRoundID.String() != t.returnedRound.Id.String() {
		return fmt.Errorf("mismatched id, created had %q but returned had %q", t.createdRoundID.String(), t.returnedRound.Id.String())
	}

	if t.returnedRound.Course.Name != t.roundRequest.Course.Name {
		return fmt.Errorf("course name mismatch, request had %q but returned had %q", t.roundRequest.Course.Name, t.returnedRound.Course.Name)
	}

	if len(t.returnedRound.Players) == 0 || len(t.returnedRound.Players) != len(t.roundRequest.Players) {
		return fmt.Errorf("incorrect number of players returned, request had %d but returned had %d", len(t.roundRequest.Players), len(t.returnedRound.Players))
	}

	for _, player := range t.returnedRound.Players {
		if len(player.Scores) != 0 {
			return fmt.Errorf("found %d scores but expected it to be empty", len(player.Scores))
		}
	}

	return nil
}

func (t *testContext) playerScoresOnHole(playerIndex, score, hole int) error {
	scoreEvent := api.PlayerScoreEvent{
		Hole:        hole,
		PlayerIndex: playerIndex - 1,
		Score:       score,
	}

	res, err := t.client.PutRoundRoundIDScore(t.execCtx, t.createdRoundID.String(), scoreEvent)

	if err != nil {
		return fmt.Errorf("failed to send score: %w", err)
	}

	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("expected status %d but got %d", http.StatusNoContent, res.StatusCode)
	}

	return nil
}

func (t *testContext) theScoreForPlayerOnHoleIs(playerIndex, hole, score int) error {
	if t.returnedRound == nil {
		return fmt.Errorf("no round to view")
	}

	if playerIndex <= 0 {
		return fmt.Errorf("player index is 1-indexed, but got %d", playerIndex)
	}

	if len(t.returnedRound.Players) < playerIndex {
		return fmt.Errorf("wanted to check player index %d but only %d players in round", playerIndex, len(t.returnedRound.Players))
	}

	player := t.returnedRound.Players[playerIndex-1]

	for _, entry := range player.Scores {
		if entry.Hole != hole {
			continue
		}

		if entry.Score != score {
			return fmt.Errorf("found score of %d but expected %d", entry.Score, score)
		}

		return nil
	}

	return fmt.Errorf("did not find player score for hole %d", hole)
}
