package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Evertras/live-leaderboards/pkg/api"
)

func (t *testContext) iCreateANewRound() error {
	//func (c *Client) PostRound(ctx context.Context, body PostRoundJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
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
		Players: []api.PlayerData{
			{
				Name: "Test Player",
			},
		},
		Title: ptr("Test Round"),
	}

	res, err := t.client.PostRound(t.execCtx, *t.roundRequest)

	if err != nil {
		return fmt.Errorf("t.client.PostRound: %w", err)
	}

	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d, expected %d", res.StatusCode, http.StatusCreated)
	}

	raw, err := io.ReadAll(res.Body)

	if err != nil {
		return fmt.Errorf("io.ReadAll: %w", err)
	}

	var createdRound api.CreatedRound

	err = json.Unmarshal(raw, &createdRound)

	if err != nil {
		return fmt.Errorf("json.Unmarshal result: %w", err)
	}

	if createdRound.Id.String() == "" {
		return fmt.Errorf("id was empty")
	}

	t.createdRoundID = createdRound.Id

	return nil
}

func (t *testContext) iViewTheRound() error {
	id := t.createdRoundID

	response, err := t.client.GetRoundRoundID(t.execCtx, id.String())

	if err != nil {
		return fmt.Errorf("failed to get round: %w", err)
	}

	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	raw, err := io.ReadAll(response.Body)

	if err != nil {
		return fmt.Errorf("failed to read body: %w", err)
	}

	var round api.Round

	err = json.Unmarshal(raw, &round)

	if err != nil {
		return fmt.Errorf("failed to unmarshal json: %w", err)
	}

	t.returnedRound = &round

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
