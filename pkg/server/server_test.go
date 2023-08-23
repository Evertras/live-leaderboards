package server_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

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
	s := server.New(repo)

	round := genTestRound()

	buf, err := json.Marshal(round)

	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	body := bytes.NewReader(buf)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/round", body)

	req.Header.Add("Content-Type", "application/json")

	s.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Result().StatusCode)

	returnedRaw := w.Body.Bytes()

	var created api.CreatedRound

	err = json.Unmarshal(returnedRaw, &created)

	assert.NoError(t, err, "Failed to unmarshal created round")
	assert.NotEmpty(t, created.Id.String())
}
