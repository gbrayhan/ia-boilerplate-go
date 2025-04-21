package integration

import (
	"github.com/cucumber/godog"
	"os"
	"testing"
	"time"
)

func TestIntegration(t *testing.T) {
	suite := godog.TestSuite{
		Name:                "integration",
		ScenarioInitializer: InitializeScenario,
		TestSuiteInitializer: func(tsc *godog.TestSuiteContext) {
			tsc.BeforeSuite(func() {
				// aquí podrías arrancar tu servidor en background
				go func() {
					// si tu main está en cmd/server/main.go:
					// cmd/server/main.Main()
				}()
				time.Sleep(3 * time.Second)
			})
		},
		Options: &godog.Options{
			Format: "pretty",
			Paths:  []string{"features"},
		},
	}

	if exitCode := suite.Run(); exitCode != 0 {
		os.Exit(exitCode)
	}
}
