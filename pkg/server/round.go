package server

import (
	"context"
	"net/http"

	"github.com/Evertras/live-leaderboards/pkg/api"
	"github.com/Evertras/live-leaderboards/pkg/repo"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type RoundRepo interface {
	CreateEventRoundStart(ctx context.Context, roundID uuid.UUID, req api.RoundRequest) error
	GetEventRoundStart(ctx context.Context, roundID uuid.UUID) (*repo.EventRoundStart, error)
}

func (s *Server) PostRound(ctx echo.Context) error {
	ctx.Logger().Info("PostRound")

	r := api.RoundRequest{}

	// TODO: add validator
	// https://echo.labstack.com/docs/request
	err := ctx.Bind(&r)

	if err != nil {
		ctx.Logger().Errorf("Failed to read body: %v", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	id := uuid.New()

	err = s.r.CreateEventRoundStart(ctx.Request().Context(), id, r)

	if err != nil {
		ctx.Logger().Errorf("Failed to create round: %v", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	err = ctx.JSON(http.StatusCreated, api.CreatedRound{
		Id: id,
	})

	if err != nil {
		ctx.Logger().Errorf("Failed to marshal JSON: %v", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return nil
}

func (s *Server) GetRoundRoundID(ctx echo.Context, roundID string) error {
	ctx.Logger().Info("GetRoundRoundID")

	id, err := uuid.Parse(roundID)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	round, err := s.r.GetEventRoundStart(ctx.Request().Context(), id)

	if err != nil {
		ctx.Logger().Errorf("Failed to get round: %v", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	err = ctx.JSON(200, round)

	if err != nil {
		ctx.Logger().Errorf("Failed to marshal JSON: %v", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return nil
}