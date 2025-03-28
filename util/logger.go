package util

import (
	"github.com/rs/zerolog"
	"os"
	"time"
)

var Logger = zerolog.New(
	zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.DateTime},
).Level(zerolog.TraceLevel).With().Timestamp().Caller().Logger()

