<div align="center">

# Spruce

Spruce pretty-prints Go structured logs for humans.

[![Documentation](https://img.shields.io/badge/go.dev-documentation-007d9c?&style=for-the-badge)](https://pkg.go.dev/github.com/dogmatiq/spruce)
[![Latest Version](https://img.shields.io/github/tag/dogmatiq/spruce.svg?&style=for-the-badge&label=semver)](https://github.com/dogmatiq/spruce/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/dogmatiq/spruce/ci.yml?style=for-the-badge&branch=main)](https://github.com/dogmatiq/spruce/actions/workflows/ci.yml)
[![Code Coverage](https://img.shields.io/codecov/c/github/dogmatiq/spruce/main.svg?style=for-the-badge)](https://codecov.io/github/dogmatiq/spruce)

</div>

Spruce provides an [`slog.Handler`] implementation that acts as an adaptor
between an [`slog.Logger`] and the [`io.Writer`] and [`testing.TB`] interfaces.

It is intended for use in places where structured log messages will be read by
**humans**, such as within tests or during development. The output format is not
particular suitable for machine parsing; use the built-in JSON or text
formatters for that.

<!-- references -->

[`io.Writer`]: https://pkg.go.dev/io#Writer
[`slog.Handler`]: https://pkg.go.dev/log/slog#Handler
[`slog.Logger`]: https://pkg.go.dev/log/slog#Logger
[`testing.TB`]: https://pkg.go.dev/testing#TB

## Example Usage

```go
package pkg_test

import (
    "log/slog"
    "testing"

    "github.com/dogmatiq/spruce"
)

func TestSomething(t *testing.T) {
    // Create a new spruce logger that directs logs to t.
    logger := spruce.NewTestLogger(t)

    // Run the application code that is being tested, and direct its log output
    // to the test log.
    systemUnderTest(logger)
}

// systemUnderTest is your existing application code that uses an [slog.Logger].
func systemUnderTest(logger *slog.Logger) {
    // Log a message with some complex structured attributes.
    logger.Info(
		"hello, world!",
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
}
```

The above test produces the following output:

```
spruce.go:76: [INFO] hello, world!
    ├─┬ <group b>
    │ ├── <key-1> ┈ <value-1>
    │ ╰── <key-2> ┈ <value-2>
    ╰─┬ <group a>
      ╰── <key-1> ┈ <value-1>
```
