package fxgorm

import (
	"context"
	"fmt"
	"time"

	"github.com/ekkinox/fx-template/modules/fxlogger"
	"github.com/rs/zerolog"
	"gorm.io/gorm/logger"
)

type GormLogger struct {
	logger     *fxlogger.Logger
	withValues bool
}

func NewGormLogger(logger *fxlogger.Logger, withValues bool) *GormLogger {
	return &GormLogger{
		logger:     logger,
		withValues: withValues,
	}
}

func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	l.logger.Level(convertGormLevel(level))

	return l
}

func (l *GormLogger) Error(ctx context.Context, msg string, opts ...interface{}) {
	zerolog.Ctx(ctx).Error().Msg(fmt.Sprintf(msg, opts...))
}

func (l *GormLogger) Warn(ctx context.Context, msg string, opts ...interface{}) {
	zerolog.Ctx(ctx).Info().Msg(fmt.Sprintf(msg, opts...))
}

func (l *GormLogger) Info(ctx context.Context, msg string, opts ...interface{}) {
	zerolog.Ctx(ctx).Info().Msg(fmt.Sprintf(msg, opts...))
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, f func() (string, int64), err error) {
	zl := zerolog.Ctx(ctx)
	var event *zerolog.Event

	if err != nil {
		event = zl.Error()
	} else {
		event = zl.Debug()
	}

	event.Str("latency", time.Since(begin).String())

	sql, rows := f()
	if sql != "" {
		event.Str("sql_query", sql)
	}
	if rows > -1 {
		event.Int64("sql_rows", rows)
	}

	event.Send()

	return
}

func (l *GormLogger) ParamsFilter(ctx context.Context, sql string, params ...interface{}) (string, []interface{}) {
	if !l.withValues {
		return sql, nil
	}
	return sql, params
}

func convertGormLevel(level logger.LogLevel) zerolog.Level {
	switch level {
	case logger.Silent:
		return zerolog.NoLevel
	case logger.Info:
		return zerolog.DebugLevel
	case logger.Warn:
		return zerolog.WarnLevel
	case logger.Error:
		return zerolog.ErrorLevel
	default:
		return zerolog.WarnLevel
	}
}
