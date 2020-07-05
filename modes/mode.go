package modes

// Mode describes some operation on the optionally sampled input from the pedal.
type Mode interface {
	// Errors is channel to pass raising errors.
	Errors() chan error

	// Close this Mode and all its internal workers.
	Close() error
}
