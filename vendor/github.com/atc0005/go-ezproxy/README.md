<!-- omit in toc -->
# go-ezproxy

Go library and tooling for working with EZproxy.

[![Latest Release](https://img.shields.io/github/release/atc0005/go-ezproxy.svg?style=flat-square)][repo-url-release-latest]
[![Go Reference](https://pkg.go.dev/badge/github.com/atc0005/go-ezproxy.svg)][docs-homepage]
[![go.mod Go version](https://img.shields.io/github/go-mod/go-version/atc0005/go-ezproxy)](https://github.com/atc0005/go-ezproxy)
[![Lint and Build](https://github.com/atc0005/go-ezproxy/actions/workflows/lint-and-build.yml/badge.svg)](https://github.com/atc0005/go-ezproxy/actions/workflows/lint-and-build.yml)
[![Project Analysis](https://github.com/atc0005/go-ezproxy/actions/workflows/project-analysis.yml/badge.svg)](https://github.com/atc0005/go-ezproxy/actions/workflows/project-analysis.yml)

<!-- omit in toc -->
## Table of contents

- [Project home](#project-home)
- [Status](#status)
- [Overview](#overview)
- [Features](#features)
  - [Current](#current)
  - [Missing](#missing)
- [Changelog](#changelog)
- [Documentation](#documentation)
- [Examples](#examples)
- [License](#license)
- [References](#references)
  - [Related projects](#related-projects)
  - [Official EZproxy docs](#official-ezproxy-docs)

## Project home

See [our GitHub repo][repo-url-home] for the latest code, to file an issue or
submit improvements for review and potential inclusion into the project.

## Status

Alpha; very much getting a feel for how the project will be structured
long-term and what functionality will be offered.

The existing functionality was added specifically to support the
in-development [atc0005/brick][related-brick-project]. This library is subject
to change in order to better support that project.

## Overview

This library is intended to provide common EZproxy-related functionality such
as reporting or terminating active login sessions (either for all usernames or
specific usernames), filtering (or not) audit file entries or traffic patterns
(not implemented yet) for specific usernames or domains.

See the [Input file formats](docs/input-files.md) doc for additional details
regarding known and supported input file formats.

**Just to be perfectly clear**:

- this library is intended to supplement the provided functionality of the
  official OCLC-developed/supported `EZproxy` application, not in any way
  replace it.
- this library is not in any way associated with OCLC, `EZproxy` or other
  services offered by OCLC.

## Features

### Current

- generate a list of audit records for session-related events
  - for all usernames
  - for a specific username

- generate a list of active sessions using the audit log
  - using entires without a corresponding logout event type

- generate a list of active sessions using the active file
  - for all usernames
  - for a specific username

- terminate user sessions
  - single user session
  - bulk user sessions

### Missing

- Anything to do with traffic log entries
- [Examples](examples/README.md)

## Changelog

See the [`CHANGELOG.md`](CHANGELOG.md) file for the changes associated with
each release of this application. Changes that have been merged to `master`,
but not yet an official release may also be noted in the file under the
`Unreleased` section. A helpful link to the Git commit history since the last
official release is also provided for further review.

## Documentation

Please see:

- our [GoDoc][docs-homepage] coverage
- our [topic-specific](docs/README.md) coverage

If something doesn't make sense, please [file an issue][repo-url-issues] and
note what is (or was) unclear.

## Examples

Please see our [GoDoc][docs-homepage] coverage for general usage and the
[examples](examples/README.md) doc for a list of applications developed using
this module.

## License

Taken directly from the [`LICENSE`](LICENSE) and [`NOTICE.txt`](NOTICE.txt) files:

```License
Copyright 2020-Present Adam Chalkley

https://github.com/atc0005/go-ezproxy/blob/master/LICENSE

Licensed under the Apache License, Version 2.0 (the "License"); you may not use
this file except in compliance with the License. You may obtain a copy of the
License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied. See the License for the
specific language governing permissions and limitations under the License.
```

## References

### Related projects

- [atc0005/brick][related-brick-project] project
  - this project uses this library to provides tools (two as of this writing)
    intended to help manage login sessions.

- <https://github.com/calvinm/ezproxy-abuse-checker>
  - <https://github.com/calvinm/ezproxy-abuse-checker/blob/d7202e617305745cf272df9918b1e95ff030f63f/block_user.pl#L32>
  - this is the project that proved to me that EZproxy sessions *can* be
    terminated programatically.

### Official EZproxy docs

- <https://help.oclc.org/Library_Management/EZproxy/EZproxy_configuration/EZproxy_system_elements>
- <https://help.oclc.org/Library_Management/EZproxy/Configure_resources/Audit>
- <https://help.oclc.org/Library_Management/EZproxy/Configure_resources/LogFormat>
- <https://help.oclc.org/Library_Management/EZproxy/Configure_resources/Option_LogSession>
- <https://help.oclc.org/Library_Management/EZproxy/Configure_resources/Option_LogUser>
- <https://help.oclc.org/Library_Management/EZproxy/Get_started/Join_the_EZproxy_listserv_and_Community_Center>

<!-- Footnotes here  -->

[repo-url-home]: <https://github.com/atc0005/go-ezproxy>  "This project's GitHub repo"
[repo-url-issues]: <https://github.com/atc0005/go-ezproxy/issues>  "This project's issues list"
[repo-url-release-latest]: <https://github.com/atc0005/go-ezproxy/releases/latest>  "This project's latest release"

[docs-homepage]: <https://pkg.go.dev/github.com/atc0005/go-ezproxy>  "GoDoc coverage"

[related-brick-project]: <https://github.com/atc0005/brick> "atc0005/brick project URL"

<!-- []: PLACEHOLDER "DESCRIPTION_HERE" -->
