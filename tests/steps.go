package tests

import "github.com/cucumber/godog"

func InitializeScenario(sc *godog.ScenarioContext) {
	t := newTestContext()

	sc.Step(`^I create a new round$`, t.iCreateANewRound)
	sc.Step(`^I view the round$`, t.iViewTheRound)
	sc.Step(`^the round is valid but empty$`, t.theRoundIsValidButEmpty)
}
