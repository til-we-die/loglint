package lowercase

import (
	"log/slog"

	"go.uber.org/zap"
)

func main() {
	slog.Info("Server started")     // want "log message should start with a lowercase letter"
	slog.Error("Failed to connect") // want "log message should start with a lowercase letter"

	slog.Info("server started")
	slog.Error("failed to connect")

	slog.Info("123 started")
	slog.Info("!start")

	logger := zap.NewExample()
	logger.Info("Server running") // want "log message should start with a lowercase letter"
	logger.Info("server running") // OK
}
