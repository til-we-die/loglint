package specials

import (
	"log/slog"
)

func main() {
	slog.Info("server started!") // want "multiple punctuation marks"
	slog.Info("failed!!")        // want "multiple punctuation marks"
	slog.Info("warning...")      // want "multiple punctuation marks"
	slog.Info("what?")
	slog.Info("what?!") // want "multiple punctuation marks"

	slog.Info("server started 😊") // want "special characters or emojis"
	slog.Info("🎉 success")        // want "special characters or emojis"
	slog.Info("error ‼")          // want "special characters or emojis"

	slog.Info("server started")
	slog.Info("connection failed")
	slog.Info("please wait")

	slog.Info("api/v1/users")
	slog.Info("key=value")
	slog.Info("user@example.com")
}
