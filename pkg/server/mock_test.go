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
	rounds        map[string]*api.Round
}

var _ server.Repo = &mockRoundRepo{}

func newMockRoundRepo() *mockRoundRepo {
	return &mockRoundRepo{
		createdEvents: make(map[string]*repo.EventRoundStart),
		rounds:        make(map[string]*api.Round),
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

func (r *mockRoundRepo) GetRound(ctx context.Context, roundID uuid.UUID) (*api.Round, error) {
	r.Lock()
	defer r.Unlock()

	if isZeroUUID(roundID) {
		return nil, fmt.Errorf("received empty uuid: %s", roundID.String())
	}

	round, exists := r.rounds[roundID.String()]

	if !exists {
		return nil, fmt.Errorf("id %q not found", roundID.String())
	}

	return round, nil
}

func ptr[K any](item K) *K {
	return &item
}

func (r *mockRoundRepo) storeRound(round *api.Round) {
	if round == nil {
		panic("can't store nil round")
	}

	if isZeroUUID(round.Id) {
		panic("zero UUID in round")
	}

	r.Lock()
	defer r.Unlock()

	r.rounds[round.Id.String()] = round
}

func isZeroUUID(id uuid.UUID) bool {
	// Maybe an easy way to do this in uuid but didn't see it off hand...
	return id.String() == "00000000-0000-0000-0000-000000000000"
}
