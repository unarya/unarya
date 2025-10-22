package logging

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"
)

// Logger is the global structured logger
var Logger zerolog.Logger

func InitLogger() {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	Logger = zerolog.New(output).With().Timestamp().Logger()
}

// FromContext attaches logger to context
func FromContext(ctx context.Context) zerolog.Logger {
	if l, ok := ctx.Value("logger").(zerolog.Logger); ok {
		return l
	}
	return Logger
}

// WithContext injects logger into context
func WithContext(ctx context.Context, fields map[string]interface{}) context.Context {
	l := Logger.With().Fields(fields).Logger()
	return context.WithValue(ctx, "logger", l)
}
