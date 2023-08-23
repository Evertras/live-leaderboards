package tests

import (
	"context"

	"github.com/google/uuid"

	"github.com/Evertras/live-leaderboards/pkg/api"
)

type testContext struct {
	client *api.Client

	createdRoundID uuid.UUID

	execCtx context.Context
}

func newTestContext() *testContext {
	client, err := api.NewClient("http://localhost:8037")

	if err != nil {
		panic(err)
	}

	return &testContext{
		client: client,
	}
}
