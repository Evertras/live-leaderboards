package tests

import (
	"context"

	"github.com/cucumber/godog"
)

func InitializeScenario(sc *godog.ScenarioContext) {
	t := newTestContext()

	sc.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		t.execCtx = ctx

		return ctx, nil
	})

	sc.Step(`^I create a new round$`, t.iCreateANewRound)
	sc.Step(`^I view the round$`, t.iViewTheRound)
	sc.Step(`^the round is valid but empty$`, t.theRoundIsValidButEmpty)
}
