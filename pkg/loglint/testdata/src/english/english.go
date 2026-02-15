package english

import (
	"log/slog"
)

func main() {
	slog.Info("запуск сервера")     // want "non-Latin character detected"
	slog.Info("ошибка подключения") // want "non-Latin character detected"
	slog.Info("starting сервер")    // want "non-Latin character detected"

	slog.Info("starting server")
	slog.Info("failed to connect")
	slog.Info("user 123 connected")
	slog.Info("api/v1/users")

	slog.Info("server запуск") // want "non-Latin character detected"
}
