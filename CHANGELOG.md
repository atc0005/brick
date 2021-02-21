# Changelog

## Overview

All notable changes to this project will be documented in this file.

The format is based on [Keep a
Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Please [open an issue](https://github.com/atc0005/brick/issues) for any
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

## [v0.4.10] - 2021-02-21

### Changed

- dependencies
  - `go.mod` Go version
    - updated from `1.14` to `1.15`
  - built using Go 1.15.8
    - **Statically linked**
    - Windows (x86, x64)
    - Linux (x86, x64)
  - `atc0005/go-teams-notify`
    - `v2.3.0` to `v2.4.2`
  - `actions/setup-node`
    - `v2.1.2` to `v2.1.4`

### Fixed

- `files.appendToFile`: Fix invalid error var reference
- Fix deferred log call formatting
- Fix explicit exit code handling

## [v0.4.9] - 2020-11-16

### Changed

- Dependencies
  - built using Go 1.15.5
    - **Statically linked**
    - Windows
      - x86
      - x64
    - Linux
      - x86
      - x64
  - `atc0005/go-ezproxy`
    - `v0.1.6` to `v0.1.7`

### Fixed

- Fix CHANGELOG entries from v0.4.8 release
- Correct the version of Go noted for v0.4.8 release

## [v0.4.8] - 2020-11-11

### Added

- Add support for limiting payloads to specific IPs

### Changed

- Statically linked binary release
  - Built using Go 1.15.4
  - Windows
    - x86
    - x64
  - Linux
    - x86
    - x64

- Dependencies
  - `actions/checkout`
    - `v2.3.3` to `v2.3.4`

### Fixed

- Fix GitHub issue reference in v0.4.7 release entry

### Notes

- Windows builds
  - Windows builds are provided, but have not been tested. The current
    developer does not have access to a Windows + EZproxy test environment.
    Please [open an issue](https://github.com/atc0005/brick/issues) to share
    your experiences deploying tools from this project on a Windows EZproxy
    server.

## [v0.4.7] - 2020-10-11

### Added

- Binary release
  - Built using Go 1.15.2
  - **Statically linked** (GH-193)
  - Windows
    - x86
    - x64
  - Linux
    - x86
    - x64

Note: Windows builds are provided, but have not been tested. The current
developer does not have access to a Windows + EZproxy test environment. Please
[open an issue](https://github.com/atc0005/brick/issues) to share your
experiences deploying tools from this project on a Windows EZproxy server.

### Changed

- Add `-trimpath` build flag
- Restore explicit exit code handling (GH-191)

### Fixed

- Makefile build options do not generate static binaries (GH-189)

## [v0.4.6] - 2020-10-02

### Added

- Binary release
  - Built using Go 1.15.2
  - Windows
    - x86
    - x64
  - Linux
    - x86
    - x64

Note: Windows builds are provided, but have not been tested. The current
developer does not have access to a Windows + EZproxy test environment. Please
[open an issue](https://github.com/atc0005/brick/issues) to share your
experiences deploying tools from this project on a Windows EZproxy server.

### Changed

- Emit version number as part of startup message

- Move subpackages into `internal` directory

- Dependencies
  - upgrade `pelletier/go-toml`
    - `v1.8.0` to `v1.8.1`
  - upgrade `actions/checkout`
    - `v2.3.2` to `v2.3.3`
  - upgrade `actions/setup-node`
    - `v2.1.1` to `v2.1.2`

### Fixed

- Misc linting errors raised by latest `gocritic` release included with
  `golangci-lint` `v1.31.0`

- Flag for setting desired log output does not appear to work

- Documentation mistake: log-output CLI flag incorrectly listed as log-out

- Makefile generates checksums with qualified path

- Debug messages are emitted before logging settings are applied which would
  (potentially) allow them to be emitted

## [v0.4.5] - 2020-08-30

### Changed

- Dependencies
  - upgrade `go.mod` Go version
    - `1.13` to `1.14`
  - upgrade `atc0005/go-ezproxy`
    - `v0.1.5` to `v0.1.6`
  - upgrade `atc0005/go-teams-notify`
    - `v1.3.1-0.20200419155834-55cca556e726` to `v2.3.0`
      - NOTE: This is a significant change reflecting a merge of required
        functionality from the `atc0005/send2teams` project to the
        `atc0005/go-teams-notify` project
  - upgrade `Showmax/go-fqdn`
    - `v0.0.0-20180501083314-6f60894d629f` to `v1.0.0`
  - upgrade `apex/log`
    - `v1.7.0` to `v1.9.0`
  - upgrade `actions/checkout`
    - `v2.3.1` to `v2.3.2`
  - upgrade `atc0005/send2teams`
    - `v0.4.5` to `v0.4.6`
      - since removed

### Fixed

- Add missing filename reference in error message

## [v0.4.4] - 2020-08-05

### Added

- Docker-based GitHub Actions Workflows
  - Replace native GitHub Actions with containers created and managed through
    the `atc0005/go-ci` project.

  - New, primary workflow
    - with parallel linting, testing and building tasks
    - with three Go environments
      - "old stable" - currently `Go 1.13.14`
      - "stable" - currently `Go 1.14.6`
      - "unstable" - currently `Go 1.15rc1`
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

### Changed

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

- Dependencies
  - upgrade `apex/log`
    - `v1.6.0` to `v1.7.0`

### Fixed

- gosec linting errors
  - G404: Use of weak random number generator (`math/rand` instead of
    `crypto/rand`)
    - fixed this, though our use of `math/rand` wasn't for cryptographic
      purposes and was likely OK as-is
  - G304: Potential file inclusion via variable
    - marked this as ignored due to the variable being one we are
      intentionally allowing the sysadmin to set

- Lock MailDev container to specific, *proven* stable version used previously
  in demos
  - intent: reduce "gotchas" in future demo sessions if a drastically
    different/newer version were to get pulled in while resetting the demo
    environment

- Email notifications do not include `Session Termination Results` section
  - this was included with existing Microsoft Teams notifications, but not
    email notifications

## [v0.4.3] - 2020-07-24

### Changed

- Explicitly note notifications state

### Fixed

- Email templates: `MISSING VALUE - Please file a bug report!`; use generated
  email summary instead of `Record.Note`

## [v0.4.2] - 2020-7-23

### Changed

- Alert sender: Replace "received by" phrasing in file templates

- Dependencies
  - updated `atc0005/go-ezproxy`
    - `v0.1.4` to `v0.1.5`

### Fixed

- Documentation
  - Further work on EZproxy purpose
  - Add further information on integration with EZproxy, Splunk

- Reporting (monitoring) system referred to with "received by" phrasing
  instead of "received from"

- Deferred file close operations report "file already closed" error messages
  - note: the `atc0005/go-ezproxy` `v0.1.5` release includes the same type of
    changes

## [v0.4.1] - 2020-07-23

### Changed

- Dependencies
  - updated `atc0005/go-ezproxy`
    - `v0.1.3` to `v0.1.4`
  - updated `actions/setup-go`
    - `v2.1.0` to `v2.1.1`
  - updated `actions/setup-node`
    - `v2.1.0` to `v2.1.1`

- Linting
  - `golangci-lint`: Disable default exclusions

- Logging
  - Update `internal/fileutils.HasLine` function to emit name
  - Update `files.appendToFile` function to emit func name
  - Update `NewConfig` function to emit name

### Fixed

- Documentation
  - Add additional lead-in for `docs/ezproxy.md` to (hopefully) better explain
    what EZproxy is
  - Update main README to make majority of "EZproxy" references point to the
    updated `docs/ezproxy.md` doc

- Linting
  - Use `filepath.Clean` for all `os.Open` calls
    - even though this application is intended for use by sysadmins (who have
      no cause to try and exploit the system), it's better to go ahead and
      guard against potential exposure introduced by using externally-provided
      (e.g., config file or flags) filenames by sanitizing the paths
    - note: the `atc0005/go-ezproxy` `v0.1.4` release includes the same type
      of changes
  - errcheck: Explicitly check file close return values
  - errcheck: Explicitly check writer flush return value

## [v0.4.0] - 2020-07-19

### Added

- Email notifications
  - initial support

### Changed

- CI/Linting
  - re-enable separate `golint` step to work around what appears to be a bug
    in golangci-lint (golangci/golangci-lint#1249)

- Dependencies
  - upgrade `apex/log`
    - `v1.4.0` to `v1.6.0`
  - upgrade `atc0005/send2teams`
    - `v0.4.4` to `v0.4.5`

- Demo content
  - upgrade Go version from `v1.14.5` to `v1.14.6`
  - minor tweaks to output emitted by reset script

- Documentation
  - Cover new flags, environment variables and config file settings
  - Misc fixes for existing rate limit, number of retries and retry delay
  - Refresh existing setup/deploy steps to briefly cover email configuration

- Configuration
  - TOML config file
    - extended with new settings
    - rename some settings in an effort to better communicate intent

### Fixed

- golint reporting several "should have comment or be unexported" linting
  issues

- in-place modification of client/alert request headers for Teams message
  formatting leads to unintentional "spillover" to email notifications

## [v0.3.0] - 2020-07-11

### Added

- expose setting to configure Teams notifications rate limit
  - via CLI flag, config file and environment variable

### Changed

- retry delay setting renamed to emphasize intent
  - configuration file setting `delay` renamed to `retry_delay`
  - CLI flag setting `--teams-notify-delay` renamed to
    `--teams-notify-retry-delay`
  - environment variable `BRICK_MSTEAMS_WEBHOOK_DELAY` renamed to
    `BRICK_MSTEAMS_WEBHOOK_RETRY_DELAY`

### Fixed

- minor wording/grammatical tweaks in an effort to clarify intent
  - e.g., `config-file` flag
- invalid function call in validation function (oh the irony)
- Update documentation for rate limit and retry delay
  - prior code and documentation failed to properly communicate the difference
    between the two goals
- Add missing documentation comments to sample configuration file

## [v0.2.0] - 2020-07-03

This release brings two notable features:

- Optional automatic user sessions termination support
- `es` binary to list or optionally terminate user sessions

`fail2ban` should work as well with this release as it did before with the
v0.1.0 release.

### Added

- Add optional native EZproxy support to terminate user sessions
- New binaries
  - `es` binary
    - small CLI app to list and optionally terminate user sessions for a
      specific username
    - intended for quick troubleshooting or as an optional replacement for
      logging into the admin UI to terminate user sessions for a specific
      username
  - `ezproxy` (mock) binary
    - small CLI binary intended to be called by `brick` for development
      purposes
    - returns some expected response codes and text for valid input
    - returns some non-standard, "unexpected" results to help test error
      handling
- See also the new
  [`atc0005/go-ezproxy`](https://github.com/atc0005/go-ezproxy) project which
  is used by this one to perform most EZproxy-related session actions
  - `atc0005/go-ezproxy` `v0.1.3` is vendored with this release

### Changed

- Update dependencies
  - `actions/checkout`
    - `v2.3.0` to `v2.3.1`
  - `actions/setup-go`
    - `v2.0.3` to `v2.1.0`
  - `actions/setup-node`
    - `v2.0.0` to `v2.1.0`
- Teams notifications
  - explicit step X of Y labeling to notification titles
  - consistent use of Note (preferred) and Error (fallback) field values to
    generate primary "summary" text
  - rename "Request Annotations" to "Request Errors" to reflect dedicated
    single purpose vs blend of Note and Error field values as before
- Documentation
  - cover new v0.2.0 features
  - attempt to present `fail2ban` and the new v0.2.0 automatic user sessions
    termination as viable options

### Fixed

- TCP port range recommendation via config validation warning
- Clarify suggested port range in config settings doc
- Force writing disabled username entries as lowercase
  - other uses of the reported username are left as-is with the intent to aid
    in troubleshooting

## [v0.1.2] - 2020-06-18

### Added

- Dependabot
  - Enable Go Modules updates
  - Enable GitHub Actions updates

### Changed

- Update dependencies
  - `apex/log`
    - `v1.1.4` to `v1.4.0`
  - `actions/setup-go`
    - `v1` to `v2.0.3`
  - `actions/checkout`
    - `v1` to `v2.3.0`
  - `actions/setup-node`
    - `v1` to `v2.0.0`

### Fixed

- Remove duplicate steps in deploy doc
- Replace invalid config file parameters
- Fix debug comment so that it reflects current behavior

## [v0.1.1] - 2020-05-24

### Fixed

- (GH-33) Fix link to removed page content
- (GH-34) Missing doc coverage for deploying the `brick` binary
  - oh what shame ...

## [v0.1.0] - 2020-05-24

### Added

Features of the initial prototype release:

- Highly configurable (with more configuration choices to be exposed in the
  future)

- Supports configuration settings from multiple sources
  - command-line flags
  - environment variables
  - configuration file
  - reasonable default settings

- Ignore individual usernames (i.e., prevent disabling listed accounts)
- Ignore individual IP Addresses (i.e., prevent disabling associated account)

- User configurable logging settings
  - levels, format and output (see [configuration settings
    doc](docs/configure.md))

- Microsoft Teams notifications
  - generated for multiple events
    - alert received
    - disabled user
    - ignored user
    - ignored IP Address
    - error occurred
  - configurable retries
  - configurable notifications delay in order to respect remote API limits

- Logging
  - Payload receipt from monitoring system
  - Action taken due to payload
    - username ignored
      - due to username inclusion in ignore file for usernames
      - due to IP Address inclusion in ignore file for IP Addresses
    - username disabled

- `contrib` files/content provided to allow for spinning up a demo environment
  in order to provide a hands-on sense of what this project can do
  - `fail2ban`
  - `postfix`
  - `docker`
    - `Maildev` container
  - `brick`
  - `rsyslog`
  - `systemd`
  - sample JSON payloads for use with `curl` or other http/API clients
  - [demo environment](docs/demo.md) doc
  - slides from group presentation/demo

Worth noting:

- Go modules (vs classic `GOPATH` setup)
- GitHub Actions Workflows which apply linting and build checks
- Makefile for general use cases (including local linting)
  - Note: See [README](README.md) first if building on Windows

### Missing

Known issues:

- Email notifications are not currently supported (see GH-3)
- Payloads are accepted from any IP Address (GH-18)
  - the expectation is that host-level firewall rules will be used to protect
    against this until a feature can be added to filter access

[Unreleased]: https://github.com/atc0005/brick/compare/v0.4.10...HEAD
[v0.4.10]: https://github.com/atc0005/brick/releases/tag/v0.4.10
[v0.4.9]: https://github.com/atc0005/brick/releases/tag/v0.4.9
[v0.4.8]: https://github.com/atc0005/brick/releases/tag/v0.4.8
[v0.4.7]: https://github.com/atc0005/brick/releases/tag/v0.4.7
[v0.4.6]: https://github.com/atc0005/brick/releases/tag/v0.4.6
[v0.4.5]: https://github.com/atc0005/brick/releases/tag/v0.4.5
[v0.4.4]: https://github.com/atc0005/brick/releases/tag/v0.4.4
[v0.4.3]: https://github.com/atc0005/brick/releases/tag/v0.4.3
[v0.4.2]: https://github.com/atc0005/brick/releases/tag/v0.4.2
[v0.4.1]: https://github.com/atc0005/brick/releases/tag/v0.4.1
[v0.4.0]: https://github.com/atc0005/brick/releases/tag/v0.4.0
[v0.3.0]: https://github.com/atc0005/brick/releases/tag/v0.3.0
[v0.2.0]: https://github.com/atc0005/brick/releases/tag/v0.2.0
[v0.1.2]: https://github.com/atc0005/brick/releases/tag/v0.1.2
[v0.1.1]: https://github.com/atc0005/brick/releases/tag/v0.1.1
[v0.1.0]: https://github.com/atc0005/brick/releases/tag/v0.1.0
