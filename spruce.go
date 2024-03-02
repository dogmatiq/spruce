package spruce

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"testing"
)

// TestingT is the subset of the [testing.TB] interface that is used to write
// logs.
type TestingT interface {
	Log(...any)
}

var _ TestingT = (testing.TB)(nil)

// NewLogger returns a [slog.Logger] that writes to t.
func NewLogger(t TestingT) *slog.Logger {
	return slog.New(NewHandler(t))
}

// NewHandler returns a new [slog.Handler] that writes to t.
func NewHandler(t TestingT) slog.Handler {
	return &handler{T: t}
}

type handler struct {
	T      TestingT
	attrs  []slog.Attr
	groups []string
}

func (h *handler) Enabled(context.Context, slog.Level) bool {
	return true
}

func (h *handler) Handle(_ context.Context, rec slog.Record) error {
	buf := &strings.Builder{}
	attrs := slices.Clone(h.attrs)

	if len(h.groups) == 0 {
		rec.Attrs(func(attr slog.Attr) bool {
			attrs = append(attrs, attr)
			return true
		})
	} else {
		var grouped []any
		rec.Attrs(func(attr slog.Attr) bool {
			grouped = append(grouped, attr)
			return true
		})

		var group slog.Attr
		for i := len(h.groups) - 1; i >= 0; i-- {
			group = slog.Group(h.groups[i], grouped...)
			grouped = []any{group}
		}

		attrs = append(attrs, group)
	}

	level := rec.Level.String()
	buf.WriteByte('[')
	buf.WriteString(level)
	buf.WriteString("] ")
	buf.WriteString(rec.Message)

	writeAttrs(buf, 0, attrs, true)

	h.T.Log(buf.String())

	return nil
}

func writeAttrs(
	buf *strings.Builder,
	depth int,
	attrs []slog.Attr,
	lastInParent bool,
) {
	if len(attrs) == 0 {
		return
	}

	width := 0
	for _, attr := range attrs {
		if len(attr.Key) > width {
			width = len(attr.Key)
		}
	}

	for i, attr := range attrs {
		last := i == len(attrs)-1

		buf.WriteByte('\n')

		for i := 0; i < depth; i++ {
			if lastInParent {
				buf.WriteString("  ")
			} else {
				buf.WriteString("│ ")
			}
		}

		if last {
			buf.WriteString("╰─")
		} else {
			buf.WriteString("├─")
		}

		if attr.Value.Kind() == slog.KindGroup {
			buf.WriteString("┬ ")
			buf.WriteString(attr.Key)
			writeAttrs(buf, depth+1, attr.Value.Group(), last)
		} else {
			buf.WriteString("─ ")
			buf.WriteString(attr.Key)
			buf.WriteString(" ")
			for i := len(attr.Key); i < width+1; i++ {
				buf.WriteString("┈")
			}
			buf.WriteString(" ")

			v := attr.Value.String()
			if strings.ContainsAny(v, " \t\n\r") {
				fmt.Fprintf(buf, "%q", v)
			} else {
				buf.WriteString(v)
			}
		}
	}
}

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &handler{
		T:      h.T,
		attrs:  append(slices.Clone(h.attrs), attrs...),
		groups: h.groups,
	}
}

func (h *handler) WithGroup(name string) slog.Handler {
	return &handler{
		T:      h.T,
		attrs:  h.attrs,
		groups: append(slices.Clone(h.groups), name),
	}
}
