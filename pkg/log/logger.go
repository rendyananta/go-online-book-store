package log

import (
	"io"
	"log/slog"
	"os"
)

func SetUp(cfg Config) {
	var output io.Writer

	if cfg.LogPath == "" {
		output = defaultLogOutput
	} else {
		file, err := os.OpenFile(cfg.LogPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			output = defaultLogOutput
		} else {
			output = file
		}
	}

	var handler slog.Handler

	if cfg.JSONFormatted {
		handler = slog.NewJSONHandler(output, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		})
	} else {
		handler = slog.NewTextHandler(output, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		})
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
}
