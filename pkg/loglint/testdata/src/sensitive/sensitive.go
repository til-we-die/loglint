package sensitive

import (
	"log/slog"

	"go.uber.org/zap"
)

func main() {
	slog.Info("user password: secret123")    // want "sensitive data"
	slog.Info("api token: abc123")           // want "sensitive data"
	slog.Info("credit card: 1234-5678-9012") // want "sensitive data"

	slog.Info("password123") // want "sensitive data"
	slog.Info("my_password") // want "sensitive data"

	slog.Info("user authenticated successfully")
	slog.Info("api request completed")
	slog.Info("database connection established")

	logger := zap.NewExample()
	logger.Info("user login",
		zap.String("username", "admin"),
		zap.String("password", "secret"),
	)
}
