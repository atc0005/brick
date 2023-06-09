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

## [v0.1.8] - 2023-06-09

### Changed

- Dependencies
  - drop individual GitHub Actions dependencies
  - import `atc0005/shared-project-resources` shared/importable CI workflows
- CI
  - (GH-73) Switch Docker image source from Docker Hub to GitHub Container
    Registry (GHCR)
  - (GH-74) Expand linting GitHub Actions Workflow to include `oldstable`,
    `unstable` container images
    (GH-78) Update project to Go 1.19
    (GH-81) Update Makefile and GitHub Actions Workflows
    (GH-85) Refactor GitHub Actions workflows to import logic
    (GH-86) Add Go Module Validation, Dependency Updates jobs
    (GH-87) Drop `Push Validation` workflow
    (GH-88) Rework workflow scheduling
    (GH-90) Remove `Push Validation` workflow status badge
    (GH-93) Update vuln analysis GHAW to remove on.push hook
    (GH-94) Add `quick` Makefile recipe (alias)

### Fixed

- (GH-56) Replace godoc.org badge with pkg.go.dev badge
- (GH-70) False positive `G307: Deferring unsafe method "Close" on type
  "*os.File" (gosec)` linting error
- (GH-79) Apply linting fixes for Go 1.19 release
- (GH-80) Remove broken ref links in CHANGELOG
- (GH-91) Update vuln analysis GHAW to use on.push hook
- (GH-92) Disable depguard linter

## [v0.1.7] - 2020-11-16

### Changed

- Dependencies
  - `actions/checkout`
    - `v2.3.2` to `v2.3.4`
  - `actions/setup-node`
    - `v2.1.1` to `v2.1.2`

### Fixed

- Date format for `v0.1.6` CHANGELOG entry
- gocritic `commentFormatting` linting errors

## [v0.1.6] - 2020-08-23

### Added

- Docker-based GitHub Actions Workflows
  - Replace native GitHub Actions with containers created and managed through
    the `atc0005/go-ci` project.

  - New, primary workflow
    - with parallel linting, testing and building tasks
    - with three Go environments
      - "old stable"
      - "stable"
      - "unstable"
    - Makefile is *not* used in this workflow
    - staticcheck linting using latest stable version provided by the
      `atc0005/go-ci` containers

  - Separate Makefile-based linting and building workflow
    - intended to help ensure that local Makefile-based builds that are
      referenced in project README files continue to work as advertised until
      a better local tool can be discovered/explored further
    - use `golang:latest` container to allow for Makefile-based linting
      tooling installation testing since the `atc0005/go-ci` project provides
      containers with those tools already pre-installed
      - linting tasks use container-provided `golangci-lint` config file
        *except* for the Makefile-driven linting task which continues to use
        the repo-provided copy of the `golangci-lint` configuration file

  - Add Quick Validation workflow
    - run on every push, everything else on pull request updates
    - linting via `golangci-lint` only
    - testing
    - no builds

- Add new README badges for additional CI workflows
  - each badge also links to the associated workflow results

### Changed

- dependencies
  - `go.mod` Go version
    - updated from `1.13` to `1.14`
  - `actions/setup-go`
    - updated from `v2.1.0` to `v2.1.2`
      - since replaced with Docker containers
  - `actions/checkout`
    - updated from `v2.3.1` to `v2.3.2`

- README
  - Link badges to applicable GitHub Actions workflows results

- Linting
  - Local
    - `Makefile`
      - install latest stable `golangci-lint` binary instead of using a fixed
          version
  - CI
    - remove repo-provided copy of `golangci-lint` config file at start of
      linting task in order to force use of Docker container-provided config
      file

## [v0.1.5] - 2020-07-23

### Fixed

- Deferred file close operations report "file already closed" error messages

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

[Unreleased]: https://github.com/atc0005/go-ezproxy/compare/v0.1.8...HEAD
[v0.1.8]: https://github.com/atc0005/go-ezproxy/releases/tag/v0.1.8
[v0.1.7]: https://github.com/atc0005/go-ezproxy/releases/tag/v0.1.7
[v0.1.6]: https://github.com/atc0005/go-ezproxy/releases/tag/v0.1.6
[v0.1.5]: https://github.com/atc0005/go-ezproxy/releases/tag/v0.1.5
[v0.1.4]: https://github.com/atc0005/go-ezproxy/releases/tag/v0.1.4
[v0.1.3]: https://github.com/atc0005/go-ezproxy/releases/tag/v0.1.3
[v0.1.2]: https://github.com/atc0005/go-ezproxy/releases/tag/v0.1.2
[v0.1.1]: https://github.com/atc0005/go-ezproxy/releases/tag/v0.1.1
[v0.1.0]: https://github.com/atc0005/go-ezproxy/releases/tag/v0.1.0

<!-- General footnotes here  -->

[repo-url-issues]: <https://github.com/atc0005/go-ezproxy/issues>  "This project's issues list"

[related-brick-project]: <https://github.com/atc0005/brick> "atc0005/brick project URL"
