package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Setup initializes zerolog for the application
func Setup(isProd bool) {
	if isProd {
		// Production mode: JSON format for machine parsing
		log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	} else {
		// Development mode: Console format for human readability
		log.Logger = zerolog.New(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "15:04:05",
		}).With().Timestamp().Logger().Level(zerolog.DebugLevel)
	}
}
