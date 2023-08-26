package server_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"

	"github.com/Evertras/live-leaderboards/pkg/api"
	"github.com/Evertras/live-leaderboards/pkg/server"
)

func genTestRound() *api.RoundRequest {
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

func TestServerPostRound(t *testing.T) {
	repo := newMockRoundRepo()
	s := server.New(repo).WithLogLevel(log.DEBUG)

	round := genTestRound()

	buf, err := json.Marshal(round)

	assert.NoError(t, err, "Failed to marshal JSON for round")

	body := bytes.NewReader(buf)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/round", body)

	req.Header.Add("Content-Type", "application/json")

	s.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Result().StatusCode)

	returnedRaw := w.Body.Bytes()

	var created api.CreatedRound

	err = json.Unmarshal(returnedRaw, &created)

	assert.NoError(t, err, "Failed to unmarshal created round")
	assert.NotEmpty(t, created.Id.String())

	internalRound, err := repo.GetEventRoundStart(context.Background(), created.Id)

	assert.NoError(t, err, "Couldn't get round start event from internal repo")

	assert.NotNil(t, internalRound.Title, "Title wasn't saved properly")
	assert.Equal(t, "Test Round", *internalRound.Title)
}

func TestServerGetRoundByID(t *testing.T) {
	repo := newMockRoundRepo()
	s := server.New(repo).WithLogLevel(log.DEBUG)

	storedRound := genTestRound()
	roundID := uuid.New()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := repo.CreateEventRoundStart(ctx, roundID, *storedRound)

	assert.NoError(t, err, "Failed to create initial round, bad test setup")

	url := fmt.Sprintf("/round/%s", roundID)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", url, nil)

	req.Header.Add("Accept", "application/json")

	s.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	returnedRaw := w.Body.Bytes()

	var round api.Round

	err = json.Unmarshal(returnedRaw, &round)

	assert.NoError(t, err, "Failed to unmarshal returned round")
	assert.NotEmpty(t, round.Id.String())

	assert.NotNil(t, round.Title, "Title wasn't retrieved properly")
	assert.Equal(t, "Test Round", round.Title, "Title is wrong")
}
