package strategy

import "gitlab.otters.xyz/jason.tevnan/gobench/internal"

// ExecutionStrategy defines what queries are run how
type ExecutionStrategy interface {
	Prepare()
	RunCommand()
	Cleanup()
	UpdateSettings(internal.Settings)
}
