package modes

import (
	"os/exec"
)

// Action defines some action to be executed on a mode's behalf.
type Action interface {
	// Execute this Action.
	Execute() error

	// Close this Action. Afterwards Execute is no longer allowed to be called.
	Close() error
}

// CommandAction is an Action which executes some (shell) command.
type CommandAction struct {
	commandStr string
}

// NewCommandAction creates a CommandAction to execute a command.
func NewCommandAction(commandStr string) CommandAction {
	return CommandAction{commandStr}
}

// Execute this CommandAction.
func (commandAction CommandAction) Execute() error {
	cmd := exec.Command("sh", "-c", commandAction.commandStr)
	return cmd.Start()
}

// Close this CommandAction. Afterwards Execute is no longer allowed to be called.
func (commandAction CommandAction) Close() error {
	return nil
}
