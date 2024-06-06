package strategy

import (
	"context"

	"github.com/tnosaj/gobench/internal"
)

// ExecutionStrategy defines what queries are run how
type ExecutionStrategy interface {
	Prepare()
	RunCommand()
	Cleanup()
	UpdateSettings(internal.Settings)
	Shutdown(context.Context)

	PopulateExistingValues([]string)
	ReturnExistingValues() []string
}
