/*

brick automatically disables EZproxy accounts via incoming webhook requests.



PROJECT HOME

See our GitHub repo (https://github.com/atc0005/brick) for the latest code, to
file an issue or submit improvements for review and potential inclusion into
the project.

PURPOSE

This application is intended to be used as a HTTP endpoint that runs alongside
an EZproxy instance. This endpoint receives webhook requests from a monitoring
system (Splunk as of this writing), disables the user account identified by
the monitoring system rules and generates one or more notifications listing
the action taken. fail2ban is used to ban the offending IP for `MaxLifetime`
minutes (EZproxy setting) + a small buffer to force active sessions associated
with the disabled user account to timeout and terminate.

The combination of stopping new logins for the disabled user account and
timing out existing sessions works around the lack of native support for this
behavior in EZproxy itself.

The net effect is that reported user accounts are immediately disabled and
existing sessions forced to timeout, at which point compromised accounts can
no longer be used on EZproxy until manually removed from the disabled users
file.

NOTE: This application has not been designed to identify user accounts
directly, but rather relies on other systems (currently limited to Splunk) to
make the decision as to which accounts should be disabled.

FEATURES

• Highly configurable (with more configuration choices to be exposed in the future)

• Supports configuration settings from multiple sources (command-line flags, environment variables, configuration file, reasonable default settings)

• User configurable logging settings (levels, format, output)

• User configurable support for ignoring specific usernames (i.e., prevent disabling listed accounts)

• User configurable support for ignoring specific IP Addresses (i.e., prevent disabling associated account)

• Microsoft Teams notifications generated for multiple events with configurable retries and notification delay

• Logging of all events (e.g., payload receipt, action taken due to payload)



USAGE

Help output is below. See the README for examples.

$ ./brick -h

Automatically disable EZproxy users via webhook requests

brick x.y.z
https://github.com/atc0005/brick


Usage: brick [--port PORT] [--ip-address IP-ADDRESS] [--log-level LOG-LEVEL] [--log-output LOG-OUTPUT] [--log-format LOG-FORMAT] [--disabled-users-file DISABLED-USERS-FILE] [--disabled-users-entry-suffix DISABLED-USERS-ENTRY-SUFFIX] [--disabled-users-file-perms DISABLED-USERS-FILE-PERMS] [--reported-users-log-file REPORTED-USERS-LOG-FILE] [--reported-users-log-file-perms REPORTED-USERS-LOG-FILE-PERMS] [--ignored-users-file IGNORED-USERS-FILE] [--ignored-ips-file IGNORED-IPS-FILE] [--teams-webhook-url TEAMS-WEBHOOK-URL] [--teams-notify-delay TEAMS-NOTIFY-DELAY] [--teams-notify-retries TEAMS-NOTIFY-RETRIES] [--ignore-lookup-errors] [--config-file CONFIG-FILE]

Options:
  --port PORT            TCP port that this application should listen on for incoming HTTP requests.
  --ip-address IP-ADDRESS
                         Local IP Address that this application should listen on for incoming HTTP requests.
  --log-level LOG-LEVEL
                         Log message priority filter. Log messages with a lower level are ignored.
  --log-output LOG-OUTPUT
                         Log messages are written to this output target.
  --log-format LOG-FORMAT
                         Log messages are written in this format.
  --disabled-users-file DISABLED-USERS-FILE
                         fully-qualified path to the EZproxy include file where this application should write disabled user accounts.
  --disabled-users-entry-suffix DISABLED-USERS-ENTRY-SUFFIX
                         The string that is appended after every username added to the disabled users file in order to deny login access.
  --disabled-users-file-perms DISABLED-USERS-FILE-PERMS
                         Desired file permissions when this file is created. Note: The ezproxy daemon will need to be able to read this file.
  --reported-users-log-file REPORTED-USERS-LOG-FILE
                         Fully-qualified path to the log file where this application should log user disable request events for fail2ban to ingest.
  --reported-users-log-file-perms REPORTED-USERS-LOG-FILE-PERMS
                         Desired file permissions when this file is created. Note: fail2ban will need to be able to read this file.
  --ignored-users-file IGNORED-USERS-FILE
                         Fully-qualified path to the file containing a list of user accounts which should not be disabled and whose IP Address reported in the same alert should not be disabled by this application. Leading and trailing whitespace per line is ignored.
  --ignored-ips-file IGNORED-IPS-FILE
                         Fully-qualified path to the file containing a list of individual IP Addresses which should not be disabled and which user account reported in the same alert should not be disabled by this application. Leading and trailing whitespace per line is ignored.
  --teams-webhook-url TEAMS-WEBHOOK-URL
                         The Webhook URL provided by a preconfigured Connector. If specified, this application will attempt to send client request details to the Microsoft Teams channel associated with the webhook URL.
  --teams-notify-delay TEAMS-NOTIFY-DELAY
                         The number of seconds to wait between Microsoft Teams message delivery attempts.
  --teams-notify-retries TEAMS-NOTIFY-RETRIES
                         The number of attempts that this application will make to deliver Microsoft Teams messages before giving up.
  --ignore-lookup-errors
                         Whether application should continue if attempts to lookup existing disabled or ignored status for a username or IP Address fail.
  --config-file CONFIG-FILE
                         Full path to optional TOML-formatted configuration file. See contrib/brick/config.example.toml for a starter template.
  --help, -h             display this help and exit
  --version              display version and exit

*/
package main
