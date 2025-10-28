package spruce

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"
)

type handler struct {
	log       func(string) error
	attrs     []slog.Attr
	groups    []string
	threshold slog.Level
	writeTime func(*strings.Builder, slog.Record)
}

func (h *handler) Enabled(context.Context, slog.Level) bool {
	return true
}

func (h *handler) Handle(_ context.Context, rec slog.Record) error {
	if rec.Level < h.threshold {
		return nil
	}

	buf := &strings.Builder{}
	var attrs []slog.Attr

	if len(h.groups) == 0 {
		attrs = slices.Clone(h.attrs)
		rec.Attrs(func(attr slog.Attr) bool {
			attrs = append(attrs, attr)
			return true
		})
	} else {
		var grouped []any

		for _, attr := range h.attrs {
			grouped = append(grouped, attr)
		}

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
	buf.WriteByte(' ')
	h.writeTime(buf, rec)

	for i, g := range h.groups {
		if i == 0 {
			buf.WriteByte(' ')
		} else {
			buf.WriteByte('.')
		}
		buf.WriteString(g)
	}

	buf.WriteString("] ")
	buf.WriteString(rec.Message)

	writeAttrs(buf, nil, attrs)

	return h.log(buf.String())
}

func writeAttrs(
	buf *strings.Builder,
	parents []bool,
	attrs []slog.Attr,
) {
	if len(attrs) == 0 {
		return
	}

	width := 0
	for _, attr := range attrs {
		if attr.Value.Kind() != slog.KindGroup {
			if len(attr.Key) > width {
				width = len(attr.Key)
			}
		}
	}

	for i, attr := range attrs {
		isLast := i == len(attrs)-1

		buf.WriteByte('\n')

		for _, parentIsLast := range parents {
			if parentIsLast {
				buf.WriteString("  ")
			} else {
				buf.WriteString("│ ")
			}
		}

		if isLast {
			buf.WriteString("╰─")
		} else {
			buf.WriteString("├─")
		}

		if attr.Value.Kind() == slog.KindGroup {
			buf.WriteString("┬ ")
			buf.WriteString(attr.Key)
			writeAttrs(
				buf,
				append(parents, isLast),
				attr.Value.Group(),
			)
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
		log:       h.log,
		attrs:     append(slices.Clone(h.attrs), attrs...),
		groups:    h.groups,
		threshold: h.threshold,
		writeTime: h.writeTime,
	}
}

func (h *handler) WithGroup(name string) slog.Handler {
	return &handler{
		log:       h.log,
		attrs:     h.attrs,
		groups:    append(slices.Clone(h.groups), name),
		threshold: h.threshold,
		writeTime: h.writeTime,
	}
}
