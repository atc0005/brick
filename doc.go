/*
brick automatically disables EZproxy accounts via incoming webhook requests.

# Project Home

See our GitHub repo (https://github.com/atc0005/brick) for the latest code, to
file an issue or submit improvements for review and potential inclusion into
the project.

# Purpose

This application is intended to be used as a HTTP endpoint that runs alongside
an EZproxy instance. This endpoint receives webhook requests from a monitoring
system (Splunk as of this writing), disables the user account identified by
the rules enabled on the monitoring system and generates one or more
notifications listing the action taken. At this point, the associated user
sessions can be optionally (and automatically) terminated using two
approaches:

 1. using (not officially documented) EZproxy binary subcommand
 2. using the provided fail2ban config files

If using native termination support, all active user sessions associated with
the reported username will be terminated using the kill subcommand provided by
the official ezproxy binary. The sysadmin will need to provide the path to the
ezproxy and the associated Active Users and Hosts "state" file where sessions
are tracked.

If installed and configured appropriately, fail2ban can be used to to monitor
the reported users log file and ban the associated IP address for
`MaxLifetime` minutes (EZproxy setting) + a small buffer to force active
sessions associated with the disabled user account to timeout and terminate.

The net effect is that reported user accounts are immediately disabled and
compromised accounts can no longer be used with EZproxy until manually removed
from the disabled users file.

NOTE: This application has not been designed to identify user accounts
directly, but rather relies on other systems (currently limited to Splunk) to
make the decision as to which accounts should be disabled.

# Features

  - Highly configurable (with more configuration choices to be exposed in the
    future)
  - Supports configuration settings from multiple sources (command-line flags,
    environment variables, configuration file, reasonable default settings)
  - User configurable logging settings (levels, format, output)
  - User configurable support for ignoring specific usernames (i.e., prevent
    disabling listed accounts)
  - User configurable support for ignoring specific IP Addresses (i.e.,
    prevent disabling associated account)
  - Logging of all events (e.g., payload receipt, action taken due to payload)
  - Optional Microsoft Teams notifications generated for multiple events with
    configurable rate limit, notification retries, and retry delay
  - Optional email notifications generated for multiple events with
    configurable rate limit, notification retries, and retry delay
  - Optional automatic (but not officially documented) termination of user
    sessions via official EZproxy binary
  - Optional filtering of JSON payload sender IP Addresses

# Usage

See the README for examples.
*/
package main
