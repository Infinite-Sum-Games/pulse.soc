package cmd

import (
	"errors"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

var Log *LoggerService

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string, err error)
	Fatal(msg string, err error)
}

type LoggerService struct {
	log zerolog.Logger
	env string
}

func NewLoggerService(environment string, file *os.File) *LoggerService {
	var output io.Writer

	if environment == "development" {
		// Logging to both file and std.out during development
		fileOut := zerolog.ConsoleWriter{Out: file, TimeFormat: time.RFC3339}
		consoleOut := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		output = zerolog.MultiLevelWriter(consoleOut, fileOut)

	} else if environment == "production" {
		// Logging only to file during production
		output = zerolog.ConsoleWriter{Out: file, TimeFormat: time.RFC3339}

	} else {
		panic(errors.New("Could not identify environment"))
	}

	logger := zerolog.New(output).With().Timestamp().Logger()
	return &LoggerService{
		log: logger,
		env: environment,
	}
}

func (l *LoggerService) Info(msg string) {
	l.log.WithLevel(zerolog.InfoLevel).Msgf("%s", msg)
}

func (l *LoggerService) Debug(msg string) {
	l.log.WithLevel(zerolog.DebugLevel).Msgf("%s", msg)
}

func (l *LoggerService) Warn(msg string) {
	l.log.WithLevel(zerolog.WarnLevel).Msgf("%s", msg)
}

func (l *LoggerService) Error(msg string, err error) {
	l.log.WithLevel(zerolog.ErrorLevel).Err(err).Msgf("%s", msg)
}

func (l *LoggerService) Fatal(msg string, err error) {
	l.log.WithLevel(zerolog.FatalLevel).Err(err).Msgf("%s", msg)
}
