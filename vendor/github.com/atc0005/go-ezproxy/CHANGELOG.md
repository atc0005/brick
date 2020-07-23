# Changelog

## Overview

All notable changes to this project will be documented in this file.

The format is based on [Keep a
Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Please [open an issue][repo-url-issues] for any
deviations that you spot; I'm still learning!.

## Types of changes

The following types of changes will be recorded in this file:

- `Added` for new features.
- `Changed` for changes in existing functionality.
- `Deprecated` for soon-to-be removed features.
- `Removed` for now removed features.
- `Fixed` for any bug fixes.
- `Security` in case of vulnerabilities.

## [Unreleased]

- placeholder

## [v0.1.4] - 2020-07-23

### Changed

- Dependencies
  - updated `actions/setup-go`
    - `v2.1.0` to `v2.1.1`
  - updated `actions/setup-node`
    - `v2.1.0` to `v2.1.1`

- Linting
  - `golangci-lint`: Disable default exclusions

- Logging
  - log original and sanitized filenames

### Fixed

- Linting
  - gosec: Wrap os.Open calls with filepath.Clean
  - golint: comment on exported const XYZ should be of the form XYZ ...
  - gosec (G204): Mute subprocess linting error for intentional `exec.Command`
    call which uses a a client-provided value
  - errcheck: Explicitly check file close return values

## [v0.1.3] - 2020-07-02

### Changed

- Document audit file and active user file fields

  - Extend GoDoc coverage
    - mention new subpackage doc coverage
    - auditlog package
      - add overview
      - add field types
      - explicit coverage of race condition
    - activefile package
      - add overview
      - add field types
        - known
        - unknown
      - line ordering
      - explicit coverage of race condition

  - README
    - Move project repo section up
    - Link to new Input File Formats doc

- Update dependencies
  - `actions/setup-go`
    - `v2.0.3` to `v2.1.0`
  - `actions/setup-node`
    - `v2.0.0` to `v2.1.0`

### Fixed

- CHANGELOG
  - fix v0.1.2 release URL

## [v0.1.2] - 2020-06-19

### Changed

- Embed `UserSession` within `TerminateSessionResult` instead of
  cherry-picking specific values. The intent is to allow deeper layers of
  client code to easily access the original `UserSession` field values (e.g.,
  IP Address).

- Update dependencies
  - `actions/checkout`
    - `v2.3.0` to `v2.3.1`

## [v0.1.1] - 2020-06-17

### Added

- New `TerminateUserSessionResults` type

- New `HasError()` method to report whether an error was recorded when
  terminating user sessions

### Changed

- Return type for multiple functions changed from
  `[]TerminateUserSessionResult` to `TerminateUserSessionResults`

- Enable Dependabot updates
  - GitHub Actions
  - Go Modules

- Update dependencies
  - `actions/setup-go`
    - `v1` to `v2.0.3`
  - `actions/checkout`
    - `v1` to `v2.3.0`
  - `actions/setup-node`
    - `v1` to `v2.0.0`

### Fixed

- Doc comment: Fix name of MatchingUserSessions func

## [v0.1.0] - 2020-06-09

Initial release!

This release provides an early release version of a library intended for use
with the processing of EZproxy related files and sessions. This library was
developed specifically to support the development of an in-progress
application, so the full context may not be entirely clear until that
application is released (currently pending review).

### Added

- generate a list of audit records for session-related events
  - for all usernames
  - for a specific username

- generate a list of active sessions using audit log
  - using entires without a corresponding logout event type

- generate a list of active sessions using active file
  - for all usernames
  - for a specific username

- terminate user sessions
  - single user session
  - bulk user sessions

- Go modules support (vs classic `GOPATH` setup)

### Missing

- Anything to do with traffic log entries
- Examples
  - the in-progress [atc0005/brick][related-brick-project] should serve well
    for this once it is released

<!-- Version header ref links here  -->

[Unreleased]: https://github.com/atc0005/go-ezproxy/compare/v0.1.4...HEAD
[v0.1.4]: https://github.com/atc0005/go-ezproxy/releases/tag/v0.1.4
[v0.1.3]: https://github.com/atc0005/go-ezproxy/releases/tag/v0.1.3
[v0.1.2]: https://github.com/atc0005/go-ezproxy/releases/tag/v0.1.2
[v0.1.1]: https://github.com/atc0005/go-ezproxy/releases/tag/v0.1.1
[v0.1.0]: https://github.com/atc0005/go-ezproxy/releases/tag/v0.1.0

<!-- General footnotes here  -->

[repo-url-home]: <https://github.com/atc0005/go-ezproxy>  "This project's GitHub repo"
[repo-url-issues]: <https://github.com/atc0005/go-ezproxy/issues>  "This project's issues list"
[repo-url-release-latest]: <https://github.com/atc0005/go-ezproxy/releases/latest>  "This project's latest release"

[docs-homepage]: <https://godoc.org/github.com/atc0005/go-ezproxy>  "GoDoc coverage"

[related-brick-project]: <https://github.com/atc0005/brick> "atc0005/brick project URL"
