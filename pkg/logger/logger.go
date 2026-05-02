package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Init(logFile string) (*os.File, error) {
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	multi := io.MultiWriter(
		file,
		zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339},
	)

	log.Logger = zerolog.New(multi).
		With().
		Timestamp().
		Caller().
		Logger()

	zerolog.TimeFieldFormat = time.RFC3339

	return file, nil
}
