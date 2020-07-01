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

[Unreleased]: https://github.com/atc0005/brick/compare/v0.2.0...HEAD
[v0.2.0]: https://github.com/atc0005/brick/releases/tag/v0.2.0
[v0.1.2]: https://github.com/atc0005/brick/releases/tag/v0.1.2
[v0.1.1]: https://github.com/atc0005/brick/releases/tag/v0.1.1
[v0.1.0]: https://github.com/atc0005/brick/releases/tag/v0.1.0
