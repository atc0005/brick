<!-- omit in toc -->
# brick: fail2ban integration

- [Project README](../README.md)

<!-- omit in toc -->
## Table of contents

- [Overview](#overview)
- [Directions](#directions)
  - [Assumptions](#assumptions)
  - [Summary](#summary)
  - [Details](#details)
    - [Install `fail2ban`](#install-fail2ban)
    - [Modify `fail2ban` config files](#modify-fail2ban-config-files)
    - [Deploy config files](#deploy-config-files)
    - [Restart `fail2ban`](#restart-fail2ban)
    - [Review `fail2ban` log file for potential issues](#review-fail2ban-log-file-for-potential-issues)

## Overview

As noted elsewhere, `fail2ban` is responsible for banning IP Addresses
associated with disable user actions. This causes active sessions associated
with the disabled username/IP Address pair to timeout and eventually
terminate. New Splunk alerts should trigger causing further IP Addresses to
be logged (and banned) until all (or nearly all) sessions for IP Addresses
associated with the disabled user account are terminated.

`fail2ban` is in a lot of ways the "missing piece" for this application.
Without `fail2ban`, `brick` serves the purpose by helping control new abuse,
but is otherwise (at this time) unable to reliably stop *active* abuse.

## Directions

### Assumptions

- your EZproxy server is a Debian-based Linux system.
  - substitute with the appropriate package manager for your Linux
    distribution (e.g., replace `apt-get` with `yum`).
- your EZproxy server is *not* a Windows system.
  - Note: Due to `brick` relying on `fail2ban`, Windows is not a supported
    platform at this time.
- your EZproxy server does not already have `fail2ban` installed
- you are going to test these (and all other setup instructions) in a test
  environment *first* before deploying to production
  - sorry, had to say it!
- your EZproxy server has a local Postfix (or other SMTP server) installation
  configured to accept mail on localhost and forward to a central relay or
  remote SMTP server
  - if this isn't the case, you will need to modify the `fail2ban`
    configuration files provided by this project to use another mail server

### Summary

1. Install `fail2ban` and dependencies
1. Modify copies of the config files provided by this project
1. Deploy your copies of the config files (modified and not)
1. Restart `fail2ban` to confirm there are no configuration errors
   - potential immediate error output to screen
   - potential error output caught by systemd
     - `sudo systemctl status fail2ban`
1. Review `/var/log/fail2ban.log` for potential issues

### Details

#### Install `fail2ban`

1. Install `fail2ban` and dependencies by:

1. `sudo apt-get update`
1. `sudo apt-get install fail2ban geoip-bin geoip-database
   geoip-database-extra sqlite3`
   - `sqlite3` is not strictly required, but may be used when troubleshooting
     in order to interact with the SQLite database used by `fail2ban`

#### Modify `fail2ban` config files

Modify `contrib/fail2ban/jail.local` file provided by this project:

- update `ignoreip` with a space-separated list of IP Addresses that should
  not be banned by `fail2ban`
  - this is a different setting, entirely separate from `brick`'s built-in
    support to ignore specific IP Addresses
  - this setting should be seen as a safety net; while `brick` can be
    configured to ignore individual IP Addresses, it does not currently
    support ignoring an entire range
- update `logpath` if you have opted to change the reported users log file
  path from the default of `/var/log/brick/users.brick-reported.log`
- update `bantime` to match the EZproxy `MaxLifetime` setting + some padding
  - e.g., `MaxLifetime` setting is `30`, so we use `2100` (seconds) to
    indicate 35 minutes (for 30 minutes `MaxLifetime` +5 minutes padding)
- update `destemail` to reflect the email address that should receive email
  alerts
  - suggestion: set this to your ticketing system's intake email address
- update `sender` to reflect the fully-qualified email address/alias that
  bounce notifications should be sent to (if there is a problem delivering
  mail to `destemail`)
- update `sendername` to whatever name you wish to have `fail2ban`
  notifications use
  - e.g., `fail2ban on EZproxy server`

Review the `contrib/fail2ban/action.d/sendmail-geoip-lines.local` file and
make any desired changes to the email template (e.g, add Redmine ticket
routing keywords).

#### Deploy config files

Deploy your copies of the config files (modified and not) provided by this
project:

- `contrib/fail2ban/jail.local`
- `contrib/fail2ban/action.d/sendmail-common.local`
- `contrib/fail2ban/action.d/sendmail-geoip-lines.local`
- `contrib/fail2ban/filter.d/brick.local`

**WARNING**: If you already have a `fail2ban` instance on your EZproxy server
you will need to merge the new settings with the previous; dropping in the new
`contrib/fail2ban/jail.local` file is intended as a safe action, but at
present at least one setting (`action_mgl`) will conflict with local settings.
See GH-28 for additional information.

#### Restart `fail2ban`

Restart `fail2ban` to confirm there are no configuration errors.

**Caution**:

- Keep a SSH session going in case you forget to whitelist your sysadmin
  workstation's IP and somehow trip a rule somewhere that would match your IP
- Make sure you inserted critical/related IP Addresses and ranges in the
  `ignoreip` config setting mentioned earlier

#### Review `fail2ban` log file for potential issues

Review `/var/log/fail2ban.log` for any errors, warnings or other details of
note.
