package logger

import (
	"errors"
	"log/slog"
	"os"
)

type logWithMeta struct {
	log  *slog.Logger
	meta []any
}

// New returns a logger that wraps metadata.
func New(args ...any) *logWithMeta {
	return &logWithMeta{
		instance,
		args,
	}
}

// SetMeta wraps metadata.
func (l *logWithMeta) SetMeta(args ...any) {
	l.meta = append(l.meta, args...)
}

// Debug outputs debug level log.
func (l *logWithMeta) Debug(msg string, args ...any) {
	slice := make([]any, 0)
	slice = append(slice, args...)
	slice = append(slice, l.meta...)
	l.log.Debug(msg, slice...)
}

// Info outputs information level log.
func (l *logWithMeta) Info(msg string, args ...any) {
	slice := make([]any, 0)
	slice = append(slice, args...)
	slice = append(slice, l.meta...)
	l.log.Info(msg, slice...)
}

// Warn outputs warning level log.
func (l *logWithMeta) Warn(msg string, args ...any) {
	slice := make([]any, 0)
	slice = append(slice, args...)
	slice = append(slice, l.meta...)
	l.log.Warn(msg, slice...)
}

// Warn outputs warning level log with a err field.
func (l *logWithMeta) WarnWithErr(msg string, err error, args ...any) {
	slice := make([]any, 0)
	slice = append(slice, "err", ErrMsg(err))
	slice = append(slice, args...)
	slice = append(slice, l.meta...)
	l.log.Warn(msg, slice...)
}

// Err outputs error level log.
//
// Returns new error with msg.
func (l *logWithMeta) Err(msg string, err error, args ...any) error {
	slice := make([]any, 0)
	slice = append(slice, "err", ErrMsg(err))
	slice = append(slice, args...)
	slice = append(slice, l.meta...)
	l.log.Error(msg, slice...)

	return errors.New(msg)
}

// Panic outputs error level log and causes panic.
func (l *logWithMeta) Panic(msg string, err error, args ...any) {
	slice := make([]any, 0)
	slice = append(slice, "severity", "Panic")
	slice = append(slice, args...)
	slice = append(slice, l.meta...)
	aErr := l.Err(msg, err, slice...)

	panic(aErr)
}

// Fatal outputs error level log and exit the program.
func (l *logWithMeta) Fatal(msg string, err error, args ...any) {
	slice := make([]any, 0)
	slice = append(slice, "severity", "Fatal")
	slice = append(slice, args...)
	slice = append(slice, l.meta...)
	l.Err(msg, err, slice...)

	os.Exit(1)
}

// NewErr outputs error level log and returns new error with err.
func (l *logWithMeta) NewErr(errMsg string, args ...any) error {
	slice := make([]any, 0)
	slice = append(slice, args...)
	slice = append(slice, l.meta...)
	l.log.Error(errMsg, slice...)

	return errors.New(errMsg)
}
