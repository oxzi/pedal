package modes

// Action defines some action to be executed on a mode's behalf.
type Action interface {
	// Execute this Action.
	Execute() error
}
