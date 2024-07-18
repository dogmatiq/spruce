package spruce_test

import (
	"strings"
	"testing"

	. "github.com/dogmatiq/spruce"
)

func TestNewStreamLogger(t *testing.T) {
	w := &strings.Builder{}
	l := NewStreamLogger(w)

	l.Info("<message>")

	got := w.String()
	want := "<message>\n"

	if !strings.HasSuffix(got, want) {
		t.Errorf("got %q, want suffix of %q", got, want)
	}
}
