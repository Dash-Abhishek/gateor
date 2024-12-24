package pkg

import (
	"log/slog"
	"os"
	"sync"
)

var Log *slog.Logger
var loggerOnce sync.Once

func init() {

	loggerOnce.Do(func() {

		logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})

		Log = slog.New(logHandler)

	})

}

func GetLogger() *slog.Logger {

	return Log

}
