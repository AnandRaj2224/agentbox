package logger

import (
	"log/slog"
	"os"
)

// New returns a JSON logger.
func New() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}
