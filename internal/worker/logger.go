package worker

// Logger defines the interface for logging in the worker package
type Logger interface {
	// Debug logs a debug message with key-value pairs
	Debug(msg string, args ...interface{})

	// Info logs an info message with key-value pairs
	Info(msg string, args ...interface{})

	// Error logs an error message with key-value pairs
	Error(msg string, args ...interface{})
}
