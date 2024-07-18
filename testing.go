package spruce

import (
	"log/slog"
	"strings"
	"testing"
	"time"
)

// TestingT is the subset of [testing.TB] that is used to write logs.
type TestingT interface {
	Log(...any)
}

var _ TestingT = (testing.TB)(nil)

// NewTestLogger returns a [slog.Logger] that writes to t.
func NewTestLogger(t TestingT) *slog.Logger {
	return slog.New(NewTestHandler(t))
}

// NewTestHandler returns a new [slog.Handler] that writes to t.
func NewTestHandler(t TestingT) slog.Handler {
	epoch := time.Now()

	return &handler{
		log: func(s string) error {
			t.Log(s)
			return nil
		},
		writeTime: func(w *strings.Builder, rec slog.Record) {
			elapsed := rec.Time.Sub(epoch)
			if elapsed > 0 {
				w.WriteByte('+')
			}
			w.WriteString(elapsed.String())
		},
	}
}
