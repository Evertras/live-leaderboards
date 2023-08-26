package server_test

import (
	"github.com/Evertras/live-leaderboards/pkg/api"
	"github.com/google/uuid"
)

func genTestRoundRequest() *api.RoundRequest {
	return &api.RoundRequest{
		Course: api.Course{
			Holes: []api.Hole{
				{
					DistanceYards: ptr(370),
					Hole:          1,
					Par:           4,
					StrokeIndex:   ptr(17),
				},
			},
			Name: "Test Course",
			Tees: ptr("White"),
		},
		Players: []api.PlayerData{
			{
				Name: "Test Player",
			},
		},
		Title: ptr("Test Round"),
	}
}

func genTestRound() *api.Round {
	return &api.Round{
		Course: api.Course{
			Holes: []api.Hole{
				{
					DistanceYards: ptr(370),
					Hole:          1,
					Par:           4,
					StrokeIndex:   ptr(17),
				},
			},
			Name: "Test Course",
			Tees: ptr("White"),
		},
		Id: uuid.New(),
		Players: []api.RoundPlayerData{{
			Name: "Test Player",
			Scores: []api.HoleScore{
				{
					Hole:  1,
					Score: 5,
				},
			},
		}},
		Title: "Test Round",
	}
}
