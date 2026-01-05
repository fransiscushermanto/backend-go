package utils

import (
	"os"

	"github.com/fransiscushermanto/backend/internal/constants"
	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

func init() {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: constants.TimeFormatDateTime}
	Logger = zerolog.New(output).With().Timestamp().Logger()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

func SetLogLevel(level string) {
	parsedLevel, err := zerolog.ParseLevel(level)

	if err != nil {
		Logger.Warn().Str("level", level).Msg("Invalid log level specified, defaulting to info")
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		return
	}

	zerolog.SetGlobalLevel(parsedLevel)
	Logger.Info().Str("log_level", parsedLevel.String()).Msg("Logger level set")
}

func Log() *zerolog.Logger {
	return &Logger
}
