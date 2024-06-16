package logger

import (
	"log/slog"
	"time"
)

// Meta returns slog.Attr.
func Meta(field string, value interface{}) slog.Attr {
	switch v := value.(type) {
	case bool:
		return slog.Bool(field, v)
	case time.Duration:
		return slog.Duration(field, v)
	case time.Time:
		return slog.Time(field, v)
	case uint64:
		return slog.Uint64(field, v)
	case int64:
		return slog.Int64(field, v)
	case int:
		return slog.Int(field, v)
	case float64:
		return slog.Float64(field, v)
	case string:
		return slog.String(field, v)
	}
	return slog.Any(field, value)
}

// Grout returns slog.Attr group.
func Group(field string, args ...any) slog.Attr {
	return slog.Group(field, args...)
}

// ErrMsg returns error content string.
//
// Mainly used to prevent memory panic caused by nil error.
func ErrMsg(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}
