package spruce

import (
	"log/slog"
	"testing"
)

// TestingT is the subset of [testing.TB] that is used to write logs.
type TestingT interface {
	Log(...any)
}

var _ TestingT = (testing.TB)(nil)

// NewTestLogger returns a [slog.Logger] that writes to t.
func NewTestLogger(t TestingT, options ...Option) *slog.Logger {
	return slog.New(NewTestHandler(t, options...))
}

// NewTestHandler returns a new [slog.Handler] that writes to t.
func NewTestHandler(t TestingT, options ...Option) slog.Handler {
	h := &handler{
		log: func(s string) error {
			t.Log(s)
			return nil
		},
	}

	applyOptions(h, options)

	return h
}
