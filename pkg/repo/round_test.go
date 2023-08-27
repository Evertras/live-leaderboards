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

	latestID, err := repo.GetLatestRoundID(ctx)
	assert.NoError(t, err, "Failed to get latest round ID")
	assert.Equal(t, id.String(), latestID, "Latest ID is not the newly created round")
}

func TestSetScoreEvent(t *testing.T) {
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

	scoreEvent := api.PlayerScoreEvent{
		PlayerIndex: 1,
		Hole:        17,
		Score:       6,
	}

	err := repo.CreateEventRoundStart(ctx, id, req)
	if err != nil {
		t.Fatalf("Failed to create event round start: %v", err)
	}
	assert.NoError(t, err, "Failed to create round")

	err = repo.SetScore(ctx, id, scoreEvent)
	assert.NoError(t, err, "Failed to set score")

	round, err := repo.GetRound(ctx, id)

	assert.NoError(t, err, "Failed to get round")
	assert.Len(t, round.Players, len(req.Players), "Unexpected number of players returned")

	playerData := round.Players[scoreEvent.PlayerIndex]

	assert.Len(t, playerData.Scores, 1, "Unexpected number of scores")
	assert.Equal(t, playerData.Scores[0], api.HoleScore{
		Hole:  scoreEvent.Hole,
		Score: scoreEvent.Score,
	}, "Unexpected score values")
}
