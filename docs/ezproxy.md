<!-- omit in toc -->
# brick: EZproxy integration

- [Project README](../README.md)

<!-- omit in toc -->
## Table of contents

- [Overview](#overview)
- [Settings](#settings)
  - [Summary](#summary)
  - [`MaxLifetime`](#maxlifetime)
  - [`::File`](#file)

## Overview

EZproxy supports authenticating user accounts using a variety (and mix) of
authentication methods. These methods include (among many others) flat-files,
LDAP and external scripts. Because EZproxy supports having multiple methods to
authenticate users, an EZproxy administrator can explicitly deny or `disable`
user accounts using an authentication method different from the one users are
normally authenticated with. This allows user accounts to authenticate using
LDAP, but be explicitly disabled using flat-files.

When asked, the OCLC Support team confirmed that EZproxy (April 2020, v7.0.16)
did not officially support automatic/programmatic termination of active
sessions. I later learned otherwise (GH-13, GH-31) and was told (paraphrasing)
that while it may work now, the feature I discovered was not officially
supported. I took this to mean that the feature could be pulled (or simply
stop working) in the future.

`brick` v0.2.0 added optional support for automatic user sessions termination.
See the [configure](configure.md) doc and the main [README](../README.md) file
for additional information for this feature. If enabled, `fail2ban` becomes an
optional layer in abuse control.

If automatic/native user sessions termination is not enabled:

- abusive user login sessions continue to be active until they are manually
  terminated or timeout
- an EZproxy administrator has to manually terminate sessions using the
  administrator web UI or use an alternative approach that deny access to
  EZproxy long enough for those sessions to expire.

Disabling user accounts using flat-files is the approach this application
currently supports. This is accomplished by adding user accounts to a specific
file, one per line, with an explicit `::deny` suffix per user account entry.
When those user accounts attempt to login again they will be denied access.

To handle timing out active sessions, this application relies upon a local
installation of `fail2ban` which runs alongside an EZproxy instance and this
application. See the [fail2ban](fail2ban.md) document for settings specific to
`fail2ban`.

This document covers specific settings used by EZproxy.

## Settings

### Summary

| Setting Name  | Filename     | Default value | Suggested value                                                    |
| ------------- | ------------ | ------------- | ------------------------------------------------------------------ |
| `MaxLifetime` | `config.txt` | 120 (minutes) | see notes                                                          |
| `::File`      | `user.txt`   | N/A           | `::File=/usr/local/etc/brick/users.brick-disabled.txt` (see notes) |

### `MaxLifetime`

`MaxLifetime` determines how long (in minutes) an EZproxy session should
remain valid after the last time it is accessed. The default is 120 minutes
and has proven to be far too generous a limit. Based on brief reading of the
official mailing list, many institutions have opted to use far shorter values
ranging from 30 minutes to as little as 15 minutes. The goal is to have this
just long enough to keep from inciting frustration from legitimate users
(e.g., whose sessions were timed out from briefly alt-tabbing elsewhere) and
yet short enough to allow user sessions associated with IP Addresses banned by
`fail2ban` to timeout (without greatly impacting service for legitimate users
that may be sharing the same IP Address).

Changing this setting in either direction *will* impact users on some level,
so a blanket "use value X" recommendation is hard to give. Instead, discuss
with your support team and faculty/student advocates to help balance impact
against the need to limit unauthorized activity against vendor resources.

### `::File`

As noted previously, EZproxy supports a mix of multiple authentication
methods. This allows authenticating the majority of user accounts via one
method, while explicitly disabling access by another source. We leverage that
support in this application by writing disabled user accounts to a flat-file
that EZproxy is configured to monitor. When a new entry appears in the file,
EZproxy honors it and denies login access to the user account when the next
login attempt occurs. As already noted, active sessions must be manually
terminated or timed out in order to interrupt access.

By default, this application is [configured](configure.md) to write entries to
`/var/cache/brick/users.brick-disabled.txt`. This can be overridden by
command-line flag (e.g., systemd unit file) or configuration file (usually
located at `/usr/local/etc/brick/config.toml`). To configure EZproxy to use
the default location, modify your `user.txt` file to include this line:

`::File=/var/cache/brick/users.brick-disabled.txt`

This line should go early in the configuration before other authentication
sources.
