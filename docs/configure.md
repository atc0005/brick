<!-- omit in toc -->
# brick: Configuring `brick`

- [Project README](../README.md)

<!-- omit in toc -->
## Table of contents

- [Precedence](#precedence)
- [Command-line Arguments](#command-line-arguments)
- [Environment Variables](#environment-variables)
- [Configuration File](#configuration-file)
- [Worth noting](#worth-noting)

## Precedence

The priority order is:

1. Command line flags (highest priority)
1. Environment variables
1. Configuration file
1. Default settings (lowest priority)

The intent is to support a *feathered* layering of configuration settings; if,
for example, a configuration file provides nearly all settings that you want,
specify just the settings that you wish to override via command-line flags (or
environment variables) and use the configuration file for the other settings.

## Command-line Arguments

- Flags marked as **`required`** must be set via CLI flag *or* within a
  TOML-formatted configuration file.
- Flags *not* marked as required are for settings where a useful default is
  already defined.

| Option                          | Required | Default                                        | Repeat | Possible                                   | Description                                                                                                                                                                                                                                                                                                                                                                                     |
| ------------------------------- | -------- | ---------------------------------------------- | ------ | ------------------------------------------ | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `h`, `help`                     | No       | `false`                                        | No     | `h`, `help`                                | Show Help text along with the list of supported flags.                                                                                                                                                                                                                                                                                                                                          |
| `config-file`                   | No       | *empty string*                                 | No     | *valid path to a file*                     | Fully-qualified path to a configuration file consulted for settings not provided via CLI flags or environment variables.                                                                                                                                                                                                                                                                        |
| `ignore-lookup-errors`          | No       | `false`                                        | No     | `true`, `false`                            | Whether application should continue if attempts to lookup existing disabled or ignored status for a username or IP Address fail. This is needed if you do not pre-create files used by this application ahead of time. WARNING: Because this can mask errors, you should probably only use it briefly when this application is first deployed, then later disabled once all files are in place. |
| `port`                          | No       | `8000`                                         | No     | *valid TCP port number*                    | TCP port that this application should listen on for incoming HTTP requests. Tip: Use an unreserved port between 1024:49151 (inclusive) for the best results.                                                                                                                                                                                                                                    |
| `ip-address`                    | No       | `localhost`                                    | No     | *valid fqdn, local name or IP Address*     | Local IP Address that this application should listen on for incoming HTTP requests.                                                                                                                                                                                                                                                                                                             |
| `log-level`                     | No       | `info`                                         | No     | `fatal`, `error`, `warn`, `info`, `debug`  | Log message priority filter. Log messages with a lower level are ignored.                                                                                                                                                                                                                                                                                                                       |
| `log-output`                    | No       | `stdout`                                       | No     | `stdout`, `stderr`                         | Log messages are written to this output target.                                                                                                                                                                                                                                                                                                                                                 |
| `log-format`                    | No       | `text`                                         | No     | `cli`, `json`, `logfmt`, `text`, `discard` | Use the specified `apex/log` package "handler" to output log messages in that handler's format.                                                                                                                                                                                                                                                                                                 |
| `disabled-users-file`           | No       | `/var/cache/brick/users.brick-disabled.txt`    | No     | *valid path to a file*                     | Fully-qualified path to the "disabled users" file                                                                                                                                                                                                                                                                                                                                               |
| `disabled-users-file-perms`     | No       | `0o644`                                        | No     | *valid permissions in octal format*        | Permissions (in octal) applied to newly created "disabled users" file. **NOTE:** `EZproxy` will need to be able to read this file.                                                                                                                                                                                                                                                              |
| `disabled-users-entry-suffix`   | No       | `::deny`                                       | No     | *valid EZproxy condition/action*           | String that is appended after every username added to the disabled users file in order to deny login access.                                                                                                                                                                                                                                                                                    |
| `reported-users-log-file`       | No       | `/var/log/brick/users.brick-reported.log`      | No     | *valid path to a file*                     | Fully-qualified path to the log file where this application should log user disable request events for fail2ban to ingest.                                                                                                                                                                                                                                                                      |
| `reported-users-log-file-perms` | No       | `0o644`                                        | No     | *valid permissions in octal format*        | Permissions (in octal) applied to newly created "reported users" log file. **NOTE:** `fail2ban` will need to be able to read this file.                                                                                                                                                                                                                                                         |
| `ignored-users-file`            | No       | `/usr/local/etc/brick/users.brick-ignored.txt` | No     | *valid path to a file*                     | Fully-qualified path to the file containing a list of user accounts which should not be disabled and whose IP Address reported in the same alert should not be banned by this application. Leading and trailing whitespace per line is ignored.                                                                                                                                                 |
| `ignored-ips-file`              | No       | `/usr/local/etc/brick/ips.brick-ignored.txt`   | No     | *valid path to a file*                     | Fully-qualified path to the file containing a list of individual IP Addresses which should not be disabled and which user account reported in the same alert should not be disabled by this application. Leading and trailing whitespace per line is ignored.                                                                                                                                   |
| `teams-webhook-url`             | No       | *empty string*                                 | No     | [*valid webhook url*](#worth-noting)       | The Webhook URL provided by a preconfigured Connector. If specified, this application will attempt to send client request details to the Microsoft Teams channel associated with the webhook URL.                                                                                                                                                                                               |
| `teams-notify-delay`            | No       | `5`                                            | No     | *valid whole number*                       | The number of seconds to wait between Microsoft Teams message delivery attempts.                                                                                                                                                                                                                                                                                                                |
| `teams-notify-retries`          | No       | `2`                                            | No     | *valid whole number*                       | The number of attempts that this application will make to deliver Microsoft Teams messages before giving up.                                                                                                                                                                                                                                                                                    |

## Environment Variables

If set, environment variables override settings provided by a configuration
file. If used, command-line arguments override the equivalent environment
variables listed below. See the [Command-line
Arguments](#command-line-arguments) table for more information.

| Flag Name                       | Environment Variable Name                   | Notes | Example (mostly using default values)                                                                                                                                                                                            |
| ------------------------------- | ------------------------------------------- | ----- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `config-file`                   | `BRICK_CONFIG_FILE`                         |       | `BRICK_CONFIG_FILE="/usr/local/etc/brick/config.toml"`                                                                                                                                                                           |
| `ignore-lookup-errors`          | `BRICK_IGNORE_LOOKUP_ERRORS`                |       | `BRICK_IGNORE_LOOKUP_ERRORS="false"`                                                                                                                                                                                             |
| `port`                          | `BRICK_LOCAL_TCP_PORT`                      |       | `BRICK_LOCAL_TCP_PORT="8000"`                                                                                                                                                                                                    |
| `ip-address`                    | `BRICK_LOCAL_IP_ADDRESS`                    |       | `BRICK_LOCAL_IP_ADDRESS="localhost"`                                                                                                                                                                                             |
| `log-level`                     | `BRICK_LOG_LEVEL`                           |       | `BRICK_LOG_LEVEL="info"`                                                                                                                                                                                                         |
| `log-output`                    | `BRICK_LOG_OUTPUT`                          |       | `BRICK_LOG_OUTPUT="stdout"`                                                                                                                                                                                                      |
| `log-format`                    | `BRICK_LOG_FORMAT`                          |       | `BRICK_LOG_FORMAT="text"`                                                                                                                                                                                                        |
| `disabled-users-file`           | `BRICK_DISABLED_USERS_FILE`                 |       | `BRICK_DISABLED_USERS_FILE="/var/cache/brick/users.brick-disabled.txt"`                                                                                                                                                          |
| `disabled-users-file-perms`     | `BRICK_DISABLED_USERS_FILE_PERMISSIONS`     |       | `BRICK_DISABLED_USERS_FILE_PERMISSIONS="0o644"`                                                                                                                                                                                  |
| `disabled-users-entry-suffix`   | `BRICK_DISABLED_USERS_ENTRY_SUFFIX`         |       | `BRICK_DISABLED_USERS_ENTRY_SUFFIX="::deny"`                                                                                                                                                                                     |
| `reported-users-log-file`       | `BRICK_REPORTED_USERS_LOG_FILE`             |       | `BRICK_REPORTED_USERS_LOG_FILE="/var/log/brick/users.brick-reported.log"`                                                                                                                                                        |
| `reported-users-log-file-perms` | `BRICK_REPORTED_USERS_LOG_FILE_PERMISSIONS` |       | `BRICK_REPORTED_USERS_LOG_FILE_PERMISSIONS="0o644"`                                                                                                                                                                              |
| `ignored-users-file`            | `BRICK_IGNORED_USERS_FILE`                  |       | `BRICK_IGNORED_USERS_FILE="/usr/local/etc/brick/users.brick-ignored.txt"`                                                                                                                                                        |
| `ignored-ips-file`              | `BRICK_IGNORED_IP_ADDRESSES_FILE`           |       | `BRICK_IGNORED_IP_ADDRESSES_FILE="/usr/local/etc/brick/ips.brick-ignored.txt"`                                                                                                                                                   |
| `teams-webhook-url`             | `BRICK_MSTEAMS_WEBHOOK_URL`                 |       | `BRICK_MSTEAMS_WEBHOOK_URL="https://outlook.office.com/webhook/a1269812-6d10-44b1-abc5-b84f93580ba0@9e7b80c7-d1eb-4b52-8582-76f921e416d9/IncomingWebhook/3fdd6767bae44ac58e5995547d66a4e4/f332c8d9-3397-4ac5-957b-b8e3fc465a8c"` |
| `teams-notify-delay`            | `BRICK_MSTEAMS_WEBHOOK_DELAY`               |       | `BRICK_MSTEAMS_WEBHOOK_DELAY="2"`                                                                                                                                                                                                |
| `teams-notify-retries`          | `BRICK_MSTEAMS_WEBHOOK_RETRIES`             |       | `BRICK_MSTEAMS_WEBHOOK_RETRIES="5"`                                                                                                                                                                                              |

## Configuration File

Configuration file settings have the lowest priority and are overridden by
settings specified in other configuration sources, except for default values.
See the [Command-line Arguments](#command-line-arguments) table for more
information, including the available values for the listed configuration
settings.

| Flag Name                       | Config file Setting Name | Section Name         | Notes |
| ------------------------------- | ------------------------ | -------------------- | ----- |
| `ignore-lookup-errors`          | `ignore_lookup_errors`   |                      |       |
| `port`                          | `local_tcp_port`         | `network`            |       |
| `ip-address`                    | `local_ip_address`       | `network`            |       |
| `log-level`                     | `level`                  | `logging`            |       |
| `log-format`                    | `format`                 | `logging`            |       |
| `log-out`                       | `output`                 | `logging`            |       |
| `disabled-users-file`           | `file_path`              | `disabledusers`      |       |
| `disabled-users-file-perms`     | `file_permissions`       | `disabledusers`      |       |
| `disabled-users-entry-suffix`   | `entry_suffix`           | `disabledusers`      |       |
| `reported-users-log-file`       | `file_path`              | `reportedusers`      |       |
| `reported-users-log-file-perms` | `file_permissions`       | `reportedusers`      |       |
| `ignored-users-file`            | `file_path`              | `ignoredusers`       |       |
| `ignored-ips-file`              | `file_path`              | `ignoredipaddresses` |       |
| `teams-webhook-url`             | `webhook_url`            | `msteams`            |       |
| `teams-notify-delay`            | `delay`                  | `msteams`            |       |
| `teams-notify-retries`          | `retries`                | `msteams`            |       |

The
[`contrib/brick/config.example.toml`](../contrib/brick/config.example.toml)
file is provided as a starting point for your own `config.toml` configuration
file. The default values provided by this configuration file should lineup
with the default application values if not specified.

Once reviewed and potentially adjusted, your copy of the `config.toml` file
can be placed in a location of your choosing and referenced using the
`--config-file` flag. See the [Command-line
arguments](#command-line-arguments) sections for usage details.

## Worth noting

- For best results, limit your choice of TCP port to an unprivileged user
  port between `1024` and `49151`

- Log format names map directly to the Handlers provided by the `apex/log`
  package. Their descriptions are copied from the [official
  README](https://github.com/apex/log/blob/master/Readme.md) and provided
  below for reference:

  | Log Format ("Handler") | Description                        |
  | ---------------------- | ---------------------------------- |
  | `cli`                  | human-friendly CLI output          |
  | `json`                 | provides log output in JSON format |
  | `logfmt`               | plain-text logfmt output           |
  | `text`                 | human-friendly colored output      |
  | `discard`              | discards all logs                  |

- Microsoft Teams webhook URLs
  - have one of two known FQDNs; both are valid as of this writing
     1. `outlook.office.com`
        - new webhook URLs use this one
     1. `outlook.office365.com`
        - older webhook URLs use this one
        - still referenced in official documentation
  - Example URL: <https://outlook.office.com/webhook/a1269812-6d10-44b1-abc5-b84f93580ba0@9e7b80c7-d1eb-4b52-8582-76f921e416d9/IncomingWebhook/3fdd6767bae44ac58e5995547d66a4e4/f332c8d9-3397-4ac5-957b-b8e3fc465a8c>
