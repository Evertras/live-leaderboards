package tests

import (
	"encoding/json"
	"fmt"
	"io"
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

	response, err := t.client.CreateRound(t.execCtx, *t.roundRequest)

	if err != nil {
		return fmt.Errorf("failed to create round: %w", err)
	}

	createdRound, err := api.ParseCreateRoundResponse(response)

	if err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if createdRound.JSON201.Id.String() == "" {
		return fmt.Errorf("id was empty")
	}

	t.createdRoundID = createdRound.JSON201.Id

	return nil
}

func (t *testContext) iViewTheRound() error {
	response, err := t.client.GetRound(t.execCtx, t.createdRoundID.String())

	if err != nil {
		return fmt.Errorf("failed to get round: %w", err)
	}

	round, err := api.ParseGetRoundResponse(response)

	if err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	t.returnedRound = round.JSON200

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

	res, err := t.client.SendScore(t.execCtx, t.createdRoundID.String(), scoreEvent)

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

func (t *testContext) iGetTheLatestRoundID() error {
	response, err := t.client.GetLatestRoundID(t.execCtx)

	if err != nil {
		return fmt.Errorf("failed to get latest: %w", err)
	}

	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected http status: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return fmt.Errorf("failed to read body: %w", err)
	}

	var latest api.RoundID

	err = json.Unmarshal(body, &latest)

	if err != nil {
		return fmt.Errorf("failed to unmarshal %q: %w", string(body), err)
	}

	t.returnedLatestRoundID = latest

	return nil
}

func (t *testContext) theLatestRoundIDMatches() error {
	if isZeroUUID(t.returnedLatestRoundID) {
		return fmt.Errorf("no latest round ID stored")
	}

	if t.returnedLatestRoundID.String() != t.createdRoundID.String() {
		return fmt.Errorf("last created round is %q but returned latest round ID is %q", t.createdRoundID.String(), t.returnedLatestRoundID.String())
	}

	return nil
}
