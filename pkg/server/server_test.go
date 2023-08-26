package server_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

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
