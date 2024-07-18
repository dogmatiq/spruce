package spruce

import (
	"io"
	"log/slog"
	"strings"
	"time"
)

var newLine = []byte{'\n'}

// NewStreamLogger returns a [slog.Logger] that writes to w.
func NewStreamLogger(w io.Writer) *slog.Logger {
	return slog.New(NewStreamHandler(w))
}

// NewStreamHandler returns a new [slog.Handler] that writes to w.
func NewStreamHandler(w io.Writer) slog.Handler {
	return &handler{
		log: func(s string) error {
			if _, err := w.Write([]byte(s)); err != nil {
				return err
			}
			if _, err := w.Write(newLine); err != nil {
				return err
			}
			return nil
		},
		writeTime: func(w *strings.Builder, rec slog.Record) {
			w.WriteString(
				rec.Time.Format(time.RFC3339),
			)
		},
	}
}
