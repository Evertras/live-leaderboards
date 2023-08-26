package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Evertras/live-leaderboards/pkg/api"
	"github.com/cucumber/godog"
)

func (t *testContext) iCreateANewRound() error {
	//func (c *Client) PostRound(ctx context.Context, body PostRoundJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	roundRequest := api.RoundRequest{
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

	res, err := t.client.PostRound(t.execCtx, roundRequest)

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
	return godog.ErrPending
}

func (t *testContext) theRoundIsValidButEmpty() error {
	return godog.ErrPending
}
