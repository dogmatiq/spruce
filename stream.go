package spruce

import (
	"io"
	"log/slog"
)

var newLine = []byte{'\n'}

// NewStreamLogger returns a [slog.Logger] that writes to w.
func NewStreamLogger(w io.Writer, options ...Option) *slog.Logger {
	return slog.New(NewStreamHandler(w, options...))
}

// NewStreamHandler returns a new [slog.Handler] that writes to w.
func NewStreamHandler(w io.Writer, options ...Option) slog.Handler {
	h := &handler{
		log: func(s string) error {
			if _, err := w.Write([]byte(s)); err != nil {
				return err
			}
			if _, err := w.Write(newLine); err != nil {
				return err
			}
			return nil
		},
		initialDepth: 1,
	}

	WithAbsoluteTimestamps()(h)
	applyOptions(h, options)

	return h
}
