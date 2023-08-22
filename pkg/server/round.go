package server

import (
	"fmt"

	"github.com/Evertras/live-leaderboards/pkg/api"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (s *Server) PostRound(ctx echo.Context) error {
	ctx.Logger().Warn("PostRound")
	return nil
}

func (s *Server) GetRoundRoundID(ctx echo.Context, roundID string) error {
	ctx.Logger().Info("GetRoundRoundID")

	r := api.Round{}

	r.Id = uuid.New()

	r.Course = api.Course{
		Holes: []api.Hole{
			{
				Distance:    ptr(330),
				Hole:        1,
				Par:         4,
				StrokeIndex: ptr(16),
			},
		},
		Name: "Pebble Beach",
		Tees: ptr("White"),
	}

	r.Title = "Super Showdown"

	r.Players = []api.RoundPlayerData{
		{
			Name: "Evertras",
			Scores: &[]api.HoleScore{
				{
					Hole:  1,
					Score: 5,
				},
			},
		},
	}

	err := ctx.JSON(200, r)

	if err != nil {
		return fmt.Errorf("failed to send JSON: %w", err)
	}

	return nil
}
