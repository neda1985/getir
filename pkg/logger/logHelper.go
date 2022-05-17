package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var logger = zerolog.New(os.Stderr).With().Timestamp().Logger()

func Start() {
	log.Logger = logger
}

func LogError(err error) {
	log.Logger.Error().Err(err)
}

func LogInfo(msg string) {
	log.Logger.Info().Msg(msg)
}

func LogPanic(err error) {
	log.Logger.Panic().Err(err)
}
