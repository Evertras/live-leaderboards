package server_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"

	"github.com/Evertras/live-leaderboards/pkg/api"
	"github.com/Evertras/live-leaderboards/pkg/server"
)

func TestServerPostRound(t *testing.T) {
	repo := newMockRoundRepo()
	s := server.New(repo).WithLogLevel(log.DEBUG)

	roundRequest := genTestRoundRequest()

	buf, err := json.Marshal(roundRequest)

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

	internalRound, exists := repo.createdEvents[created.Id.String()]

	assert.True(t, exists, "Round creation not found")
	assert.NotNil(t, internalRound.Title, "Title was nil")
	assert.Equal(t, "Test Round", *internalRound.Title, "Mismatched title")
}

func TestServerGetRoundByID(t *testing.T) {
	repo := newMockRoundRepo()
	s := server.New(repo).WithLogLevel(log.DEBUG)

	storedRound := genTestRound()

	repo.storeRound(storedRound)

	url := fmt.Sprintf("/round/%s", storedRound.Id.String())
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", url, nil)

	req.Header.Add("Accept", "application/json")

	s.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	returnedRaw := w.Body.Bytes()

	var round api.Round

	err := json.Unmarshal(returnedRaw, &round)

	assert.NoError(t, err, "Failed to unmarshal returned round")
	assert.NotEmpty(t, round.Id.String())

	assert.Equal(t, storedRound.Id, round.Id, "Mismatched ID")
	assert.NotNil(t, round.Title, "Title wasn't retrieved properly")
	assert.Equal(t, "Test Round", round.Title, "Title is wrong")
}

func TestServerSendScore(t *testing.T) {
	repo := newMockRoundRepo()
	s := server.New(repo).WithLogLevel(log.DEBUG)

	storedRound := genTestRound()

	repo.storeRound(storedRound)

	// Sanity check
	storedScores := repo.getScoreEvents(storedRound.Id)
	assert.Len(t, storedScores, 0, "Unexpected number of stored scores, bad test setup")

	url := fmt.Sprintf("/round/%s/score", storedRound.Id.String())

	playerScoreEvent := api.PlayerScoreEvent{
		Hole:        17,
		PlayerIndex: 2,
		Score:       5,
	}

	buf, err := json.Marshal(playerScoreEvent)

	assert.NoError(t, err, "Failed to marshal request json")

	body := bytes.NewReader(buf)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", url, body)
	req.Header.Add("Content-Type", "application/json")

	s.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Result().StatusCode)

	storedScores = repo.getScoreEvents(storedRound.Id)
	assert.Len(t, storedScores, 1, "Unexpected number of stored scores")
	storedScore := storedScores[0]
	assert.Equal(t, playerScoreEvent, storedScore)
}

func TestServerGetLatestRoundID(t *testing.T) {
	repo := newMockRoundRepo()
	s := server.New(repo).WithLogLevel(log.DEBUG)

	latestID := uuid.New()

	repo.setLatestRoundID(latestID)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/latest/roundID", nil)
	req.Header.Add("Accept", "application/json")

	s.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	var returnedID uuid.UUID

	err := json.Unmarshal(w.Body.Bytes(), &returnedID)

	assert.NoError(t, err, "Failed to unmarshal ID response")
	assert.Equal(t, latestID, returnedID)
}
