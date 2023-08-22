package repo_test

import (
	"context"
	"testing"

	"github.com/Evertras/live-leaderboards/pkg/api"
	"github.com/google/uuid"
)

func TestRoundEventCreateAndFetch(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	repo := newRepo(ctx, t)

	id := uuid.New()

	req := api.RoundRequest{
		Course: api.Course{
			Holes: []api.Hole{
				{
					DistanceYards: ptr(340),
					Hole:          1,
					Par:           4,
					StrokeIndex:   ptr(13),
				},
			},
			Name: "Test Course",
			Tees: ptr("White"),
		},
		Players: []api.PlayerData{
			{
				Name: "Test Player 1",
			},
			{
				Name: "Test Player 2",
			},
		},
		Title: ptr("Super Showdown"),
	}

	err := repo.CreateEventRoundStart(ctx, id, req)
	if err != nil {
		t.Fatalf("Failed to create event round start: %v", err)
	}

	returnedData, err := repo.GetEventRoundStart(ctx, id)
	if err != nil {
		t.Fatalf("Failed to fetch event: %v", err)
	}

	if returnedData == nil {
		t.Fatalf("Returned data was nil")
	}

	if returnedData.RoundID != id.String() {
		t.Fatalf("Round ID was %q but expected %q", returnedData.RoundID, id.String())
	}

	if returnedData.Course.Name != req.Course.Name {
		t.Fatalf("Course name was %q but expected %q", returnedData.Course.Name, req.Course.Name)
	}
}
