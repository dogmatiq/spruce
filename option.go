package spruce

import (
	"log/slog"
	"strings"
	"time"
)

// Option is a function that configures a handler.
type Option func(h *handler)

// WithRelativeTimestamps configures a handler to write timestamps as relative
// durations from the time at which the handler was created.
func WithRelativeTimestamps() Option {
	return func(h *handler) {
		epoch := time.Now()

		h.writeTime = func(w *strings.Builder, rec slog.Record) {
			elapsed := rec.Time.Sub(epoch)
			if elapsed > 0 {
				w.WriteByte('+')
			}
			w.WriteString(elapsed.String())
		}
	}
}

// WithAbsoluteTimestamps configures a handler to write timestamps as absolute
// RFC3339 formatted strings.
func WithAbsoluteTimestamps() Option {
	return func(h *handler) {
		h.writeTime = func(w *strings.Builder, rec slog.Record) {
			w.WriteString(
				rec.Time.Format(time.RFC3339),
			)
		}
	}
}

// WithThreshold configures a handler to only log messages at or above the
// specified level.
func WithThreshold(level slog.Level) Option {
	return func(h *handler) {
		h.threshold = level
	}
}

func applyOptions(h *handler, options []Option) {
	options = append(
		[]Option{
			WithRelativeTimestamps(),
			WithThreshold(slog.LevelDebug),
		},

		options...,
	)

	for _, opt := range options {
		opt(h)
	}
}
