package testdata

import (
	"log/slog"
)

func main() {
	slog.Info("Starting server on port 8080") // want "log message should start with a lowercase letter"
	slog.Debug("Debug message")               // want "log message should start with a lowercase letter"
	slog.Warn("Warning message")              // want "log message should start with a lowercase letter"
	slog.Error("Error message")               // want "log message should start with a lowercase letter"

	println("Not a log call")
}
