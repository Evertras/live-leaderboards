package server_test

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"

	"github.com/Evertras/live-leaderboards/pkg/api"
	"github.com/Evertras/live-leaderboards/pkg/repo"
	"github.com/Evertras/live-leaderboards/pkg/server"
)

type mockRoundRepo struct {
	sync.Mutex

	createdEvents map[string]*repo.EventRoundStart
}

var _ server.Repo = &mockRoundRepo{}

func newMockRoundRepo() *mockRoundRepo {
	return &mockRoundRepo{
		createdEvents: make(map[string]*repo.EventRoundStart),
	}
}

func (r *mockRoundRepo) CreateEventRoundStart(ctx context.Context, roundID uuid.UUID, req api.RoundRequest) error {
	r.Lock()
	defer r.Unlock()

	r.createdEvents[roundID.String()] = &repo.EventRoundStart{
		RoundID:      roundID.String(),
		SortKey:      "rnd_start",
		RoundRequest: req,
	}

	return nil
}

func (r *mockRoundRepo) GetEventRoundStart(ctx context.Context, roundID uuid.UUID) (*repo.EventRoundStart, error) {
	r.Lock()
	defer r.Unlock()

	if roundID.String() == "00000000-0000-0000-0000-000000000000" {
		return nil, fmt.Errorf("received empty uuid: %s", roundID.String())
	}

	round, exists := r.createdEvents[roundID.String()]

	if !exists {
		return nil, fmt.Errorf("id %q not found", roundID.String())
	}

	return round, nil
}

func ptr[K any](item K) *K {
	return &item
}
