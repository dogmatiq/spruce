# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog], and this project adheres to
[Semantic Versioning].

<!-- references -->

[Keep a Changelog]: https://keepachangelog.com/en/1.0.0/
[Semantic Versioning]: https://semver.org/spec/v2.0.0.html

## [Unreleased]

### Changed

- The group name is now displayed in the "meta-data" section at the beginning of
  the log line.

### Fixed

- `With()` now places attributes in the current group instead of at the root.

## [0.2.0] - 2024-07-19

### Added

- Added `NewStreamLogger()` and `NewStreamHandler()` functions that create
  loggers and handlers that write to an `io.Writer`.

### Changed

- **[BC]** Renamed `NewLogger` to `NewTestLogger`
- **[BC]** Renamed `NewHandler` to `NewTestHandler`

## [0.1.1] - 2024-04-09

- Include the elapsed duration since the logger was created in each log message.

## [0.1.0] - 2024-03-02

- Initial release.

<!-- references -->

[Unreleased]: https://github.com/dogmatiq/spruce
[0.1.0]: https://github.com/dogmatiq/spruce/releases/tag/v0.1.0
[0.1.1]: https://github.com/dogmatiq/spruce/releases/tag/v0.1.1
[0.2.0]: https://github.com/dogmatiq/spruce/releases/tag/v0.2.0

<!-- version template
## [0.0.1] - YYYY-MM-DD

### Added
### Changed
### Deprecated
### Removed
### Fixed
### Security
-->
