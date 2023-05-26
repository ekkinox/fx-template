package fxhttpserver

import (
	"fmt"
	"github.com/ekkinox/fx-template/modules/fxlogger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/rs/zerolog"
	"io"
	"strings"
)

var headerRequestID = strings.ToLower(echo.HeaderXRequestID)
var headerTraceParent = "traceparent"

type EchoLogger struct {
	logger *fxlogger.Logger
	prefix string
}

func NewEchoLogger(logger *fxlogger.Logger) *EchoLogger {
	return &EchoLogger{
		logger: logger,
	}
}

func (e *EchoLogger) SetLogger(logger *fxlogger.Logger) *EchoLogger {
	e.logger = logger

	return e
}

func (e *EchoLogger) Output() io.Writer {
	return e.logger
}

func (e *EchoLogger) SetOutput(w io.Writer) {
	e.logger.Output(w)
}

func (e *EchoLogger) Prefix() string {
	return e.prefix
}

func (e *EchoLogger) SetPrefix(p string) {
	e.prefix = p
}

func (e *EchoLogger) Level() log.Lvl {
	return convertZeroLevel(e.logger.GetLevel())
}

func (e *EchoLogger) SetLevel(v log.Lvl) {
	e.logger.Level(convertEchoLevel(v))
}

func (e *EchoLogger) SetHeader(h string) {
	e.logger.With().Str("header", h)
}

func (e *EchoLogger) Debug(i ...interface{}) {
	e.logger.Debug().Msg(fmt.Sprint(i...))
}

func (e *EchoLogger) Debugf(format string, args ...interface{}) {
	e.logger.Debug().Msgf(format, args...)
}

func (e *EchoLogger) Debugj(j log.JSON) {
	e.logJSON(e.logger.Debug(), j)
}

func (e *EchoLogger) Info(i ...interface{}) {
	e.logger.Info().Msg(fmt.Sprint(i...))
}

func (e *EchoLogger) Infof(format string, args ...interface{}) {
	e.logger.Info().Msgf(format, args...)
}

func (e *EchoLogger) Infoj(j log.JSON) {
	e.logJSON(e.logger.Info(), j)
}

func (e *EchoLogger) Warn(i ...interface{}) {
	e.logger.Warn().Msg(fmt.Sprint(i...))
}

func (e *EchoLogger) Warnf(format string, args ...interface{}) {
	e.logger.Warn().Msgf(format, args...)
}

func (e *EchoLogger) Warnj(j log.JSON) {
	e.logJSON(e.logger.Warn(), j)
}

func (e *EchoLogger) Error(i ...interface{}) {
	e.logger.Error().Msg(fmt.Sprint(i...))
}

func (e *EchoLogger) Errorf(format string, args ...interface{}) {
	e.logger.Error().Msgf(format, args...)
}

func (e *EchoLogger) Errorj(j log.JSON) {
	e.logJSON(e.logger.Error(), j)
}

func (e *EchoLogger) Fatal(i ...interface{}) {
	e.logger.Fatal().Msg(fmt.Sprint(i...))
}

func (e *EchoLogger) Fatalj(j log.JSON) {
	e.logJSON(e.logger.Fatal(), j)
}

func (e *EchoLogger) Fatalf(format string, args ...interface{}) {
	e.logger.Fatal().Msgf(format, args...)
}

func (e *EchoLogger) Panic(i ...interface{}) {
	e.logger.Panic().Msg(fmt.Sprint(i...))
}

func (e *EchoLogger) Panicj(j log.JSON) {
	e.logJSON(e.logger.Panic(), j)
}

func (e *EchoLogger) Panicf(format string, args ...interface{}) {
	e.logger.Panic().Msgf(format, args...)
}

func (e *EchoLogger) Print(i ...interface{}) {
	e.logger.WithLevel(zerolog.NoLevel).Str("level", "-").Msg(fmt.Sprint(i...))
}

func (e *EchoLogger) Printf(format string, i ...interface{}) {
	e.logger.WithLevel(zerolog.NoLevel).Str("level", "-").Msgf(format, i...)
}

func (e *EchoLogger) Printj(j log.JSON) {
	e.logJSON(e.logger.WithLevel(zerolog.NoLevel).Str("level", "-"), j)
}

func (e *EchoLogger) logJSON(event *zerolog.Event, j log.JSON) {
	for k, v := range j {
		event = event.Interface(k, v)
	}

	event.Msg("")
}

func convertZeroLevel(level zerolog.Level) log.Lvl {
	switch level {
	case zerolog.TraceLevel:
		return log.DEBUG
	case zerolog.DebugLevel:
		return log.DEBUG
	case zerolog.InfoLevel:
		return log.INFO
	case zerolog.WarnLevel:
		return log.WARN
	case zerolog.ErrorLevel:
		return log.ERROR
	case zerolog.NoLevel:
		return log.OFF
	default:
		return log.INFO
	}
}

func convertEchoLevel(level log.Lvl) zerolog.Level {
	switch level {
	case log.DEBUG:
		return zerolog.DebugLevel
	case log.INFO:
		return zerolog.InfoLevel
	case log.WARN:
		return zerolog.WarnLevel
	case log.ERROR:
		return zerolog.ErrorLevel
	case log.OFF:
		return zerolog.NoLevel
	default:
		return zerolog.InfoLevel
	}
}
