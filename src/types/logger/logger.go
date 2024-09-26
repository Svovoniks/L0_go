package logger

import (
	"os"

	"github.com/rs/zerolog"
)

const LogFile = "log.file"

var (
	Logger zerolog.Logger
)

func SetupLogger() (*os.File, error) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	writer, err := os.OpenFile(LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		writer = os.Stdout
		Logger.Err(err).Str("filename", LogFile).Msg("Couldn't redirect logger to file")
	}

	Logger = zerolog.New(writer).With().Timestamp().Logger()
	Logger.Info().Msg("Logger setup")
	return writer, err
}
