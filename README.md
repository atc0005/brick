<!-- omit in toc -->
# brick

![brick project logo][brick-logo]

Automatically disable [EZproxy](docs/ezproxy.md) user accounts via incoming
webhook requests.

[![Latest Release](https://img.shields.io/github/release/atc0005/brick.svg?style=flat-square)](https://github.com/atc0005/brick/releases/latest)
[![GoDoc](https://godoc.org/github.com/atc0005/brick?status.svg)](https://godoc.org/github.com/atc0005/brick)
[![Validate Codebase](https://github.com/atc0005/brick/workflows/Validate%20Codebase/badge.svg)](https://github.com/atc0005/brick/actions?query=workflow%3A%22Validate+Codebase%22)
[![Validate Docs](https://github.com/atc0005/brick/workflows/Validate%20Docs/badge.svg)](https://github.com/atc0005/brick/actions?query=workflow%3A%22Validate+Docs%22)
[![Lint and Build using Makefile](https://github.com/atc0005/brick/workflows/Lint%20and%20Build%20using%20Makefile/badge.svg)](https://github.com/atc0005/brick/actions?query=workflow%3A%22Lint+and+Build+using+Makefile%22)
[![Quick Validation](https://github.com/atc0005/brick/workflows/Quick%20Validation/badge.svg)](https://github.com/atc0005/brick/actions?query=workflow%3A%22Quick+Validation%22)

<!-- omit in toc -->
## Table of contents

- [Project home](#project-home)
- [Status](#status)
- [Overview](#overview)
  - [EZproxy's purpose](#ezproxys-purpose)
  - [brick's purpose](#bricks-purpose)
  - [brick in action](#brick-in-action)
  - [The end result](#the-end-result)
- [Features](#features)
  - [Current](#current)
  - [Missing](#missing)
  - [Future](#future)
- [Changelog](#changelog)
- [Documentation](#documentation)
  - [TL;DR](#tldr)
  - [Hands-on / demo](#hands-on--demo)
  - [Index](#index)
- [License](#license)
- [References](#references)

## Project home

See [our GitHub repo](https://github.com/atc0005/brick) for the latest code,
to file an issue or submit improvements for review and potential inclusion
into the project.

## Status

Alpha quality; most
[MVP](https://en.wikipedia.org/wiki/Minimum_viable_product) functionality is
present, but a number of improvements to make this tool useful to a broad
audience are still missing.

## Overview

### EZproxy's purpose

Per the [EZproxy Wikipedia page](https://en.wikipedia.org/wiki/EZproxy):

> EZproxy is a web proxy server used by libraries to give access from outside
the library's computer network to restricted-access websites that authenticate
users by IP address. This allows library patrons at home or elsewhere to log
in through their library's EZproxy server and gain access to resources to
which their library subscribes, such as bibliographic databases.

### brick's purpose

`brick` is an application that is intended to run alongside an existing
[EZproxy](docs/ezproxy.md) instance and a Splunk monitoring system in order to
disable (suspected) compromised user accounts.

### brick in action

When Splunk, identifies suspect account behavior (based on sysadmin-specified
thresholds), it delivers a webhook payload to a HTTP endpoint on `brick` for
processing. `brick` then disables the suspect user account and (optionally)
generates one or more notifications listing the action(s) taken. All actions
are logged.

At this point, the associated user sessions can be optionally (and
automatically) terminated using two approaches:

1. using (not officially documented) `ezproxy` binary subcommand
1. using the provided fail2ban config files

If using native termination support, all active user sessions associated with
the reported username will be terminated using the kill subcommand provided by
the official `ezproxy` binary. The sysadmin will need to provide the path to
the `ezproxy` binary and the associated Active Users and Hosts "state" file
where sessions are tracked.

If installed and configured appropriately, fail2ban can be used to to monitor
the reported users log file and ban the associated IP address. This ban should
be for `MaxLifetime` minutes ([EZproxy](docs/ezproxy.md) setting) + a small
buffer to force active sessions associated with the disabled user account to
timeout and terminate.

### The end result

The net effect is that reported user accounts are immediately disabled and
compromised accounts can no longer be used with [EZproxy](docs/ezproxy.md)
until manually removed from the disabled users file.

**NOTE:** `brick` has not been designed to identify user accounts directly.
Instead, this application but rather relies on other systems (currently
limited to Splunk) to make the decision as to which accounts should be
disabled.

The hope is to extend payload format support to include Graylog and other
popular log monitoring systems in the future.

See also:

- [High-level overview](docs/start-here.md)
- our other [documentation](#documentation) for further instructions

## Features

### Current

- Highly configurable (with more configuration choices to be exposed in the future)

- Optional automatic (but not officially documented) termination of user
  sessions via official `ezproxy` binary

- `es` CLI application
  - small CLI app to list and optionally terminate user sessions for a
    specific username
  - intended for quick troubleshooting or as an optional replacement for
    logging into the admin UI to terminate user sessions for a specific
    username

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

- Optional notifications
  - Microsoft Teams
  - Email
  - generated for multiple events
    - alert received
    - disabled user
    - ignored user
    - ignored IP Address
    - error occurred
  - configurable retries and delay between retries
  - configurable notifications rate limit (used to respect remote API limits)

- Logging
  - Payload receipt from monitoring system
  - Action taken due to payload
    - username ignored
      - due to username inclusion in ignore file for usernames
      - due to IP Address inclusion in ignore file for IP Addresses
    - username disabled

- `contrib` files/content
  - intended solely for demo purposes
    - `postfix`
    - `docker`
      - `Maildev` container
    - sample JSON payloads for use with `curl` or other http/API clients
    - [demo environment](docs/demo.md) doc
    - slides from group presentation/demo
    - shell scripts to setup test/demo environment
  - intended for demo *and* as a template for production use
    - `fail2ban`
    - `brick`
    - `rsyslog`
    - `systemd`

The `contrib` content is provided both to allow for spinning up a [demo
environment](docs/demo.md) in order to provide a hands-on sense of what this
project can do and (at least some of the files) to use as a template for a
production installation (e.g., the `fail2ban` config files). At some point we
hope to provide one or more Ansible playbooks (GH-29) to replace the shell
scripts currently used by this project for setting up a test/demo environment.

### Missing

Known issues:

- Documentation
  - The docs are beginning to take overall shape, but still need a lot of work
- Payloads are accepted from any IP Address
  - the expectation is that host-level firewall rules will be used to protect
    against this until a feature can be added to filter access (see GH-18)

### Future

| Priority | Milestone                                                                                  | Description                                                                                                                                            |
| -------- | ------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------ |
| Low      | [Unplanned](https://github.com/atc0005/brick/milestone/5)                                  | Potential merit, but are either low demand or are more complex to implement than other issues under consideration.                                     |
| Medium   | [Future](https://github.com/atc0005/brick/milestone/3)                                     | Considered to have merit and sufficiently low complexity to fit within a near-future milestone.                                                        |
| High     | [vX.Y.Z](https://github.com/atc0005/brick/milestones?direction=desc&sort=title&state=open) | Milestones with a [semantic versioning](https://semver.org/) pattern reflect collections of issues that are in a planning or active development state. |

## Changelog

See the [`CHANGELOG.md`](CHANGELOG.md) file for the changes associated with
each release of this application. Changes that have been merged to `master`,
but not yet in an official release may also be noted in the file under the
`Unreleased` section. A helpful link to the Git commit history since the last
official release is also provided for further review.

## Documentation

### TL;DR

1. Install dependencies, including `fail2ban`
1. Setup environment
   1. create new service account/group
   1. create log directory
   1. create cache directory for disabled users file
   1. set ownership, permissions
   1. customize config files
      - `brick`
      - `fail2ban`
      - ...
   1. deploy config files (including systemd unit file, logrotate conf, ...)
1. Build `brick`
1. Deploy `brick`
1. Configure [EZproxy](docs/ezproxy.md) to use new disabled users file
1. Configure Splunk alerts
1. Test!

### Hands-on / demo

These resources were used to provide a demo to our team prior to gather
feedback, well before deployment plans were locked-in. The notes are detailed
(perhaps overly so) and were used as a reference before/during the team demo.

The PowerPoint is much like it sounds and the demo scripts are used to
(reasonably) quickly prepare a local "throwaway" VM for the demo. Reviewing
this material may help you decide if this project is a good fit for use at
your institution.

- [demo notes](docs/demo.md) (verbose)
- [demo scripts](contrib/demo/README.md)
- Slides used in live demo
  - [PowerPoint file](contrib/demo/ppt/brick-demo.pptx)
  - [PDF file](contrib/demo/ppt/brick-demo.pdf)

### Index

| Order | Name                             | Description                                                                                               |
| ----- | -------------------------------- | --------------------------------------------------------------------------------------------------------- |
| 1     | [Why?](docs/start-here.md)       | High-level overview of application design and purpose                                                     |
| *2*   | [Demo](docs/demo.md)             | Presentation material for a local demo that I will provide to showcase existing and future functionality. |
| 2     | [Build](docs/build.md)           | Building/compiling `brick`                                                                                |
| 3     | [Deploy](docs/deploy.md)         | Deploying `brick`                                                                                         |
| 4     | [Configure](docs/configure.md)   | Settings supported by `brick`                                                                             |
| 5     | [Fail2Ban](docs/fail2ban.md)     | Brief coverage on integrating with Fail2Ban (to monitor and take action on events recorded by `brick`)    |
| 6     | [EZproxy](docs/ezproxy.md)       | Brief coverage on integrating with EZproxy (suggested settings, using files generated by `brick`)         |
| 7     | [Rsyslog](docs/rsyslog.md)       | Brief coverage on adding a Rsyslog action to route messages from `brick`                                  |
| 8     | [Endpoints](docs/endpoints.md)   | Current endpoints offered by `brick`                                                                      |
| 9     | [Splunk](docs/splunk.md)         | Brief coverage on configuring an alert to send to `brick` (NOTE: *highly* environment specific)           |
| 10    | [References](docs/references.md) | Various reference material used while developing `brick`                                                  |

## License

Taken directly from the [`LICENSE`](LICENSE) and [`NOTICE.txt`](NOTICE.txt) files:

```License
Copyright 2020-Present Adam Chalkley

https://github.com/atc0005/brick/blob/master/LICENSE

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

Various references used when developing this project can be found in our
[references](docs/references.md) doc.

[brick-logo]: media/brick-logo-rounded.png "brick project logo"
