package repo_test

import (
	"context"
	"testing"

	"github.com/Evertras/live-leaderboards/pkg/api"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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

	round, err := repo.GetRound(ctx, id)
	assert.NoError(t, err, "Failed to get round")
	assert.NotNil(t, round, "Returned round is nil")
	assert.Equal(t, id.String(), round.Id.String(), "Wrong ID")
	assert.Equal(t, req.Course, round.Course, "Course mismatch")
}
