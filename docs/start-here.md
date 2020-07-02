<!-- omit in toc -->
# brick: Why brick was created / a high-level overview

- [Project README](../README.md)

<!-- omit in toc -->
## Table of contents

- [Problem Scenario](#problem-scenario)
- [Requirements / Dependencies](#requirements--dependencies)
- [Configuration](#configuration)
- [Behavior](#behavior)
  - [If using EZproxy's native (unofficial) session termination support](#if-using-ezproxys-native-unofficial-session-termination-support)
  - [If using fail2ban to timeout sessions](#if-using-fail2ban-to-timeout-sessions)
- [Notifications](#notifications)

## Problem Scenario

1. A compromised account logs in to EZproxy on Friday at 6 pm and abuses
   vendor resources throughout the weekend.
1. We learn about this from a Splunk alert (e.g., from an email that we find
   during our next maintenance window) or after the vendor notifies us.

This last scenario was usually when we learned that some number of our campus
IP Addresses were blocked by the vendor and would not be unblocked until we
identified the specific accounts responsible (and dealt with them). While a
reasonable request, this was not the best way to start the week.

We needed a way to take automatic action on Splunk alerts when we were not
around to manually do so. `brick` is one answer to that problem.

## Requirements / Dependencies

| Application         | brick's role                                                                        | Application's role                                                                                            | (Potential) Alternatives                    |
| ------------------- | ----------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------- | ------------------------------------------- |
| Splunk              | brick accepts Splunk payloads, processes them (including generating notifications)  | Splunk provides payloads that match specified alert criteria                                                  | Graylog                                     |
| fail2ban (optional) | brick logs alert receipt, logs actions taken (including disable and ignore actions) | fail2ban uses brick's disable action log messages to ban/unban IP Addresses associated with a disable request | Built-in/native session termination support |
| EZproxy             | brick maintains disabled users flat-file                                            | EZproxy monitors brick's generated disabled users flat-file to deny login access to specified user accounts   |                                             |

## Configuration

1. The `brick` application is installed on our EZproxy server and runs as a
   separate service
1. The `fail2ban` application is (optionally) installed on our EZproxy server
   and runs as a separate service
1. `EZproxy` is configured to look at a new, automatically-maintained disabled
   users file
1. Splunk alerts are configured & enabled
   - using existing, high-confidence, time-tested, email-based alerts for
     specific abuse patterns

## Behavior

As of v0.2.0, `brick` supports using one or both of EZproxy's native
(unofficial) session termination support and `fail2ban` to force sessions to
timeout. The sysadmin deploying `brick` can choose one or both of the options.
The behavior for these choices is noted below.

### If using EZproxy's native (unofficial) session termination support

In this scenario the native session termination support is enabled and
`fail2ban` is either not installed or not configured to monitor the report
users log file.

1. Splunk sends alert based on time-tested thresholds
1. `brick` receives the alert, logs it, evaluates it (based on ignore
   lists, etc), takes an action, sends notifications
   - If not asked to ignore a username or IP Address, `brick`
       1. adds the username to a local "disabled users" flat-file that EZproxy
          is configured to monitor
       1. logs that action
   - If ignoring a username or IP Address
       - logs that the username or IP Address was ignored
1. EZproxy sees that the user account is disabled and denies future logins,
   but refrains from terminating existing sessions
1. `brick` calls `/fully/qualified/path/to/ezproxy kill SESSION_ID` for each
   active user session for the reported username
1. EZproxy sessions for the disabled user account are terminated

### If using fail2ban to timeout sessions

In this scenario the native session termination support is *not* enabled and
`fail2ban` is used exclusively to halt abusive user accounts.

1. Splunk sends alert based on time-tested thresholds
1. `brick` receives the alert, logs it, evaluates it (based on ignore
   lists, etc), takes an action, sends notifications
   - If not asked to ignore a username or IP Address, `brick`
       1. adds the username to a local "disabled users" flat-file that EZproxy
          is configured to monitor
       1. logs that action
   - If ignoring a username or IP Address
       - logs that the username or IP Address was ignored
1. EZproxy sees that the user account is disabled and denies future logins,
   but refrains from terminating existing sessions
1. Fail2ban (if applicable) finds an entry that indicates a user account
   was disabled
   1. Bans the IP for `MaxLifetime` (EZproxy setting) minutes + some padding
      - to force session timeout for the user account associated with the
        banned IP Address
   1. Sends an email notification indicating the IP was banned, the GeoIP
      details (if available) and the matching log records from `brick`
      (includes username, IP Address and action taken)
1. EZproxy sessions for the disabled user account timeout and terminate
1. Fail2ban automatically unbans the IP after the ban timer (`MaxLifetime`
   EZproxy setting + some padding) expires

## Notifications

Notifications are generated by `brick` as these actions occur:

- alert received
- disabled user
- ignored user
- ignored IP Address
- error occurred
