package spruce_test

import (
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/dogmatiq/spruce"
	. "github.com/dogmatiq/spruce"
	"github.com/google/go-cmp/cmp"
)

func TestHandler_noAttributes(t *testing.T) {
	s := &testingTStub{T: t}
	l := NewTestLogger(s)

	l.Info("<message>")
	s.Expect(
		`[INFO <timestamp>] <message>`,
	)
}

func TestHandler_stringAttribute(t *testing.T) {
	s := &testingTStub{T: t}
	l := NewTestLogger(s)

	l.Info("<message>", "<key>", "<value>")
	s.Expect(
		`[INFO <timestamp>] <message>`,
		`╰── <key> ┈ <value>`,
	)
}

func TestHandler_stringerAttribute(t *testing.T) {
	s := &testingTStub{T: t}
	l := NewTestLogger(s)

	l.Info(
		"<message>",
		"duration",
		1*time.Second,
	)
	s.Expect(
		`[INFO <timestamp>] <message>`,
		`╰── duration ┈ 1s`,
	)
}

func TestHandler_attributeAlignment(t *testing.T) {
	s := &testingTStub{T: t}
	l := NewTestLogger(s)

	l.Info(
		"<message>",
		"<short>", "<value-1>",
		"<much longer>", "<value-2>",
	)
	s.Expect(
		`[INFO <timestamp>] <message>`,
		`├── <short> ┈┈┈┈┈┈┈ <value-1>`,
		`╰── <much longer> ┈ <value-2>`,
	)
}

func TestHandler_nestedAttributes(t *testing.T) {
	s := &testingTStub{T: t}
	l := NewTestLogger(s)

	l.Info(
		"<message>",
		slog.Group(
			"<group b>",
			"<key-1>", "<value-1>",
			"<key-2>", "<value-2>",
		),
		slog.Group(
			"<group a>",
			"<key-1>", "<value-1>",
		),
	)
	s.Expect(
		`[INFO <timestamp>] <message>`,
		`├─┬ <group b>`,
		`│ ├── <key-1> ┈ <value-1>`,
		`│ ╰── <key-2> ┈ <value-2>`,
		`╰─┬ <group a>`,
		`  ╰── <key-1> ┈ <value-1>`,
	)
}

func TestHandler_whitespaceEscaping(t *testing.T) {
	s := &testingTStub{T: t}
	l := NewTestLogger(s)

	l.Info(
		"<message>",
		"<key>", "value\nwith\nnewlines",
	)
	s.Expect(
		`[INFO <timestamp>] <message>`,
		`╰── <key> ┈ "value\nwith\nnewlines"`,
	)
}

func TestHandler_WithAttrs(t *testing.T) {
	s := &testingTStub{T: t}
	l := spruce.
		NewTestLogger(s).
		With("<key>", "<value>")

	l.Info("<message>")
	s.Expect(
		`[INFO <timestamp>] <message>`,
		`╰── <key> ┈ <value>`,
	)
}

func TestHandler_WithAttrs_sameKey(t *testing.T) {
	s := &testingTStub{T: t}
	l := spruce.
		NewTestLogger(s).
		With("<key>", "<value-1>")

	l.Info("<message>", "<key>", "<value-2>")
	s.Expect(
		`[INFO <timestamp>] <message>`,
		`├── <key> ┈ <value-1>`,
		`╰── <key> ┈ <value-2>`,
	)
}

func TestHandler_WithAttrs_inGroupKey(t *testing.T) {
	s := &testingTStub{T: t}
	l := spruce.
		NewTestLogger(s).
		WithGroup("<group>").
		With("<key>", "<value>")

	l.Info("<message>")
	s.Expect(
		`[INFO <timestamp>] <message>`,
		`╰─┬ <group>`,
		`  ╰── <key> ┈ <value>`,
	)
}

func TestHandler_WithGroup(t *testing.T) {
	s := &testingTStub{T: t}
	l := spruce.
		NewTestLogger(s).
		WithGroup("<group>")

	l.Info("<message>", "<key>", "<value>")
	s.Expect(
		`[INFO <timestamp>] <message>`,
		`╰─┬ <group>`,
		`  ╰── <key> ┈ <value>`,
	)
}

type testingTStub struct {
	T      *testing.T
	actual []string
}

func (s *testingTStub) Log(args ...any) {
	s.T.Log(args...)
	s.actual = append(s.actual, fmt.Sprint(args...))
}

func (s *testingTStub) Expect(lines ...string) {
	timestamp := regexp.MustCompile(`\[([A-Z]+) .+?\]`)

	expect := strings.Join(lines, "\n")
	for _, line := range s.actual {
		line = timestamp.ReplaceAllString(line, "[$1 <timestamp>]")
		if line == expect {
			return
		}
	}

	if diff := cmp.Diff([]string{expect}, s.actual); diff != "" {
		s.T.Errorf("unexpected logs (-want, +got): %s", diff)
	}
}
