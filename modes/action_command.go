package modes

import "os/exec"

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
