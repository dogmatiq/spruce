<div align="center">

# Spruce

Spruce adapts Go [structured logging] to pretty test output.

[![Documentation](https://img.shields.io/badge/go.dev-documentation-007d9c?&style=for-the-badge)](https://pkg.go.dev/github.com/dogmatiq/spruce)
[![Latest Version](https://img.shields.io/github/tag/dogmatiq/spruce.svg?&style=for-the-badge&label=semver)](https://github.com/dogmatiq/spruce/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/dogmatiq/spruce/ci.yml?style=for-the-badge&branch=main)](https://github.com/dogmatiq/spruce/actions/workflows/ci.yml)
[![Code Coverage](https://img.shields.io/codecov/c/github/dogmatiq/spruce/main.svg?style=for-the-badge)](https://codecov.io/github/dogmatiq/spruce)

</div>

## Example

```go
package pkg_test

import (
    "log/slog"
    "testing"

    "github.com/dogmatiq/spruce"
)

func TestSomething(t *testing.T) {
    // Create a new spruce logger that directs logs to t.
    logger := spruce.NewLogger(t)

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

[structured logging]: https://pkg.go.dev/slog
