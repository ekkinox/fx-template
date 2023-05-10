package fxlogger

// see https://github.com/ipfans/fxlogger/tree/main

import (
	"go.uber.org/fx/fxevent"
	"strings"
)

func (l *Logger) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		l.log.Info().Str("callee", e.FunctionName).
			Str("caller", e.CallerName).
			Msg("OnStart hook executing")
	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			l.log.Warn().Err(e.Err).
				Str("callee", e.FunctionName).
				Str("caller", e.CallerName).
				Msg("OnStart hook failed")
		} else {
			l.log.Info().Str("callee", e.FunctionName).
				Str("caller", e.CallerName).
				Str("runtime", e.Runtime.String()).
				Msg("OnStart hook executed")
		}
	case *fxevent.OnStopExecuting:
		l.log.Info().Str("callee", e.FunctionName).
			Str("caller", e.CallerName).
			Msg("OnStop hook executing")
	case *fxevent.OnStopExecuted:
		if e.Err != nil {
			l.log.Warn().Err(e.Err).
				Str("callee", e.FunctionName).
				Str("callee", e.CallerName).
				Msg("OnStop hook failed")
		} else {
			l.log.Info().Str("callee", e.FunctionName).
				Str("caller", e.CallerName).
				Str("runtime", e.Runtime.String()).
				Msg("OnStop hook executed")
		}
	case *fxevent.Supplied:
		if e.Err != nil {
			l.log.Warn().Err(e.Err).Str("type", e.TypeName).Msg("supplied")
		} else {
			l.log.Info().Str("type", e.TypeName).Msg("supplied")
		}
	case *fxevent.Provided:
		for _, rtype := range e.OutputTypeNames {
			l.log.Info().Str("type", rtype).
				Str("constructor", e.ConstructorName).
				Msg("provided")
		}
		if e.Err != nil {
			l.log.Error().Err(e.Err).Msg("error encountered while applying options")
		}
	case *fxevent.Invoking:
		// Do nothing. Will log on Invoked.

	case *fxevent.Invoked:
		if e.Err != nil {
			l.log.Error().Err(e.Err).Str("stack", e.Trace).
				Str("function", e.FunctionName).Msg("invoke failed")
		} else {
			l.log.Info().Str("function", e.FunctionName).Msg("invoked")
		}
	case *fxevent.Stopping:
		l.log.Info().Str("signal", strings.ToUpper(e.Signal.String())).Msg("received signal")
	case *fxevent.Stopped:
		if e.Err != nil {
			l.log.Error().Err(e.Err).Msg("stop failed")
		}
	case *fxevent.RollingBack:
		l.log.Error().Err(e.StartErr).Msg("start failed, rolling back")
	case *fxevent.RolledBack:
		if e.Err != nil {
			l.log.Error().Err(e.Err).Msg("rollback failed")
		}
	case *fxevent.Started:
		if e.Err != nil {
			l.log.Error().Err(e.Err).Msg("start failed")
		} else {
			l.log.Info().Msg("started")
		}
	case *fxevent.LoggerInitialized:
		if e.Err != nil {
			l.log.Error().Err(e.Err).Msg("custom logger initialization failed")
		} else {
			l.log.Info().Str("function", e.ConstructorName).Msg("initialized custom fxevent.Logger")
		}
	}
}
