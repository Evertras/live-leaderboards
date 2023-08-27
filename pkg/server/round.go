package server

import (
	"context"
	"net/http"

	"github.com/Evertras/live-leaderboards/pkg/api"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type RoundRepo interface {
	CreateEventRoundStart(ctx context.Context, roundID uuid.UUID, req api.RoundRequest) error
	GetRound(ctx context.Context, roundID uuid.UUID) (*api.Round, error)
	SetScore(ctx context.Context, roundID uuid.UUID, scoreData api.PlayerScoreEvent) error
	GetLatestRoundID(ctx context.Context) (string, error)
}

func (s *Server) CreateRound(ctx echo.Context) error {
	ctx.Logger().Info("CreateRound")

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

func (s *Server) GetRound(ctx echo.Context, roundID string) error {
	logger := ctx.Logger()
	logger.Infof("GetRound: %s", roundID)

	id, err := uuid.Parse(roundID)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	logger.Debug(id)
	round, err := s.r.GetRound(ctx.Request().Context(), id)

	if err != nil {
		ctx.Logger().Errorf("Failed to get round: %v", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	err = ctx.JSON(http.StatusOK, round)

	if err != nil {
		ctx.Logger().Errorf("Failed to marshal JSON: %v", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return nil
}

func (s *Server) GetLatestRoundID(ctx echo.Context) error {
	logger := ctx.Logger()
	logger.Info("GetLatestRoundID")

	latest, err := s.r.GetLatestRoundID(ctx.Request().Context())

	if err != nil {
		logger.Errorf("Failed to get latest round ID: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	err = ctx.JSON(200, latest)

	if err != nil {
		logger.Errorf("Failed to marshal latest value: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return nil
}

func (s *Server) SendScore(ctx echo.Context, roundID string) error {
	logger := ctx.Logger()
	logger.Infof("SendScore: %s", roundID)

	id, err := uuid.Parse(roundID)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	logger.Debug(id)

	var scoreEvent api.PlayerScoreEvent

	err = ctx.Bind(&scoreEvent)

	if err != nil {
		ctx.Logger().Errorf("Failed to read body: %v", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = s.r.SetScore(ctx.Request().Context(), id, scoreEvent)

	if err != nil {
		logger.Errorf("Failed to set score: %v", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return echo.NewHTTPError(http.StatusNoContent)
}
