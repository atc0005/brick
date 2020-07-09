<!-- omit in toc -->
# brick: Demo

![brick project logo][brick-logo]

- [Project README](../README.md)

<!-- omit in toc -->
## Table of contents

- [Why brick was created](#why-brick-was-created)
- [Preparation](#preparation)
  - [Environment](#environment)
  - [Web browser / MailDev](#web-browser--maildev)
  - [Visual Studio Code](#visual-studio-code)
  - [brick](#brick)
  - [Terminal](#terminal)
- [Presentation](#presentation)
  - [Slideshow](#slideshow)
  - [First payload with as-is demo settings](#first-payload-with-as-is-demo-settings)
  - [Simulate a repeat Splunk alert](#simulate-a-repeat-splunk-alert)
  - [Enable Teams and email notifications](#enable-teams-and-email-notifications)
  - [Same test username, different IP Address](#same-test-username-different-ip-address)
  - [Ignored username, different IP Address](#ignored-username-different-ip-address)
  - [Ignored IP Address, different username](#ignored-ip-address-different-username)
- [High-level overview](#high-level-overview)
- [Improvements](#improvements)
  - [Recent](#recent)
  - [Planned](#planned)
  - [Potential](#potential)
- [Upstream feature requests](#upstream-feature-requests)
- [References](#references)

## Why brick was created

This demo is intended to answer that question by covering:

1. brief history of dealing with account abuse
1. current functionality provided by this application
1. future improvements, pending approval of time and resources

## Preparation

### Environment

1. Boot throwaway Ubuntu Linux 18.04 LTS desktop VM
   older or newer versions *may* work, but have not received as much testing
1. Install `git`
   - `sudo apt-get update && sudo apt-get install -y git`
1. Clone remote repo
   1. `cd $HOME/Desktop`
   1. `git clone https://github.com/atc0005/brick`
      - Note: Depending on project development timing, this might actually be
        a different repo path
1. Run demo scripts 1-4
   1. `contrib/demo/scripts/01-install-docker.sh`
   1. `contrib/demo/scripts/02-install-dependencies.sh`
   1. `contrib/demo/scripts/03-install-demo-tooling.sh`
   1. `contrib/demo/scripts/04-setup-env.sh`
1. Shutdown VM
1. Take snapshot
1. Boot VM
1. Run `go version` to ensure that `go` is properly registered in `PATH`
1. Run demo script 5
   - `contrib/demo/scripts/05-deploy.sh`
1. Reset demo environment
   - run `contrib/demo/scripts/06-reset-env.sh`
   - this is even for the first part of the demo to ensure we are working from
     a consistent state throughout the demo
   - this does not modify existing copies of the config files in the
     `/usr/local/etc/brick/` directory

### Web browser / MailDev

1. Open a web browser directly within the demo VM or on the VM host
1. Navigate to the IP Address of the VM + `:1080`
   - e.g., `http://192.168.92.136:1080/` or `http://localhost:1080/`
1. Enable automatic refresh of mailbox contents

### Visual Studio Code

1. Add `brick` repo path to workspace
1. Add `/var/log/brick` path to workspace
1. Add `/usr/local/etc/brick` path to workspace
1. Add `/var/cache/brick` path to workspace
1. Add `/usr/local/ezproxy` path to workspace
1. Install extensions
   - `Even Better TOML` (`tamasfe.even-better-toml`)
1. Configure settings

   ```json
   {
      "workbench.startupEditor": "newUntitledFile",
      "telemetry.enableTelemetry": false,
      "telemetry.enableCrashReporter": false,
      "workbench.activityBar.visible": true,
      "editor.minimap.enabled": false
   }
   ```

1. Open `/usr/local/etc/brick/config.toml`
1. Open `/usr/local/etc/brick/users.brick-ignored.txt`
1. Open `/usr/local/etc/brick/ips.brick-ignored.txt`

We're going to make changes to these latter files in a bit and having them
open already during the demo should help speed things up.

### brick

Configure `brick`:

1. Open `/usr/local/etc/brick/config.toml`
1. Set `ignore_lookup_errors` to `true`
1. Retrieve webhook URL from Microsoft Teams Connector
   - assumption: We will use an existing "testing" channel
1. Set `msteams.webhook_url` to the webhook URL retrieved from the test
   Microsoft Teams channel
1. Comment out `msteams.webhook_url` line
   - this disables Microsoft Teams notifications for now
1. Set `ezproxy.terminate_sessions` to `true`
   - TODO: Decide if this will be enabled initially, or as a follow-up item
1. Configure email settings
   1. `email.server` to `localhost`
   1. `email.port` to `25`
   1. `email.client_identity` to `brick`
      - a production system should set this to the fully-qualified hostname of
        the sending system
      - a production system should also have a [Forward-confirmed reverse DNS
        (FCrDNS](https://en.wikipedia.org/wiki/Forward-confirmed_reverse_DNS)
        configuration to enable the most reliable delivery possible
   1. `email.sender_address` to `brick@example.org`
   1. `email.recipient_addresses` to `["help@example.org"]`
   1. `email.rate_limit` to `3`
   1. `email.retries` to `2`
   1. `email.retry_delay` to `2`
1. Comment out `email.server`
   - this disables email notifications for now

### Terminal

Configure `guake` for demo:

1. create labeled tabs
   1. `Deploy`
   1. `Splunk (Payload)`
   1. `syslog entries`
   1. `disabled user entries`
   1. `reported user entries`
    - aka, "action" entries
   1. `fail2ban log`
1. Configure 100% height/width
1. Disable transparency
1. Enable 9000 lines scroll-back (hopefully not needed)
1. Clear screen in each `guake` tab, tail log files, pipe
   through `ccze -A`
   1. `Deploy`
      - none
   1. `Splunk (Payload)`
      - No log file; stage `curl` command
        1. `cd $HOME/Desktop/brick/contrib/tests`
        1. `curl -X POST -H "Content-Type: application/json" -d
           @splunk-sanitized-payload-formatted.json
           http://localhost:8000/api/v1/users/disable`
   1. `syslog entries`
      - `clear && tail -f /var/log/brick/syslog.log | ccze -A`
   1. `disabled user entries`
      - `clear && tail -f /var/cache/brick/users.brick-disabled.txt | ccze -A`
   1. `reported user entries`
      - `clear && tail -f /var/log/brick/users.brick-reported.log | ccze -A`
   1. `fail2ban log`
      - `clear && tail -f /var/log/fail2ban.log | ccze -A`
   1. `email log`
      1. `clear && tail -f /var/log/mail.log | ccze -A`

## Presentation

### Slideshow

1. Start slideshow
1. Step through each slide, checking for notes in the presenter's Notes box
1. Make sure to pause momentarily for questions
1. On the Demonstration slide, switch back to this document (offscreen)

### First payload with as-is demo settings

1. Start with payload tab
   - Mention that this emulates receiving a Splunk alert
   - We have permission from NetSec to run the search every 5 minutes
   - Splunk alerts will continue to trigger for every Username/IP pair
     associated with the thresholds we have defined for the Splunk alerts.
   - We may have to adjust those thresholds to get the desired behavior over
     time.
1. Switch to syslog entries tab
   1. syslog entries are automatically forwarded to sawmill1 and on to Graylog
      where they are searchable
   1. Mention the periodic stats output showing how many notification messages
      were generated vs what was actually set
      - At this point Teams should have all zeros; we haven't enabled Teams
        notifications yet
1. Switch to disabled user entries tab
   - Emphasize format and fields recorded
      - comment
          1. `Username`
          1. `Source IP` (aka, `User IP`)
          1. `Alert name` (1:1 with name shown in `Splunk Alerts` panel)
          1. `Alert sender IP` (aka, `Distributed IT Splunk search head`)
              - Useful if we have someone intentionally (or otherwise) attempt
                to spoof a Splunk alert payload
              - Note: Splunk does not appear to support providing any sort of
                authentication token/key to prove its identity, so we'll need
                to use host-level firewall rules to help control abuse
          1. SearchID (should be 1:1 with value recorded in Splunk's logs if
             we ever need to dig that deep)
      - `username::deny` format
        - matches existing `users.disabled.txt` file format which we maintain
          by hand
        - plan for disabled user files
          - keep one file for manual entries
          - add another to be automatically maintained by `brick`
1. Switch to reported user entries
   1. Mention that the entries come in pairs
       1. *an alert was received*
       1. *this is what we did about it*
   1. Mention that there is no guarantee that they'll be in order (though they
      often will be)
       - use `SearchID` to match up report/action pairs
   1. Mention that we'll come back to this tab later to show additional entry
      types
1. Switch to fail2ban log tab
   1. Note `bantime` value
   1. Note `found` line
   1. Note `Ban` line
   1. Mention that we'll come back in a moment
1. Open web browser
   1. Navigate to `http://localhost:1080/`
   1. Show the fail2ban alert email
      1. mention that GeoIP functionality would be available for non-private
         IPs
      1. mention that we can customize the alerts to use Redmine "include"
         pages to allow a standardized "what do I do with this alert?" guide

### Simulate a repeat Splunk alert

1. Switch back to `Splunk (Payload)` tab
   1. Submit another payload for the same username
   1. Mention that this simulates a repeat alert in case the Splunk
      Agent/Forwarder on the EZproxy server gets "stuck" after the initial
      alert and finally unsticks sending in queued entries after the
      "cooldown" timer on the Splunk "search head" expires, which from
      `brick`'s perspective is a duplicate disable request
1. Switch to the `syslog entries` tab
   - Note "already disabled" wording
1. Switch to the `disabled user entries` tab
   - Note that there is still only one entry
     - this is as expected
1. Switch to the `reported user entries` tab
   - Note the new pair of entries
     1. Report: same format as before
     1. Action: same "DISABLE" prefix as before, but …
        1. The log entries explicitly notes that the user account has already
           been disabled, but had it not, it would be again due to XYZ alert
           (should be a different SearchID than the original disable action)
        1. The IP Address could be different and often would be different for
           cases where a compromised account is shared. We want to disable ALL
           IPs associated with a compromised account (note that we'll return
           to that shortly; (we're going to demo the "ignored users/ips"
           support, but don't mention that just yet)
1. Switch to the `fail2ban log` tab
   1. Note that it recognizes the "duplicate" report
      - The current configuration doesn't do this, but fail2ban offers the
        ability to reset unban timers for each additional match
      - We don't do that on purpose; we want to ban the IP of the original
        disable request only long enough to terminate the matching user
        session as other legitimate users may be coming from that IP.
   1. Note that once the fail2ban timer expires for the original IP, other
      users from that IP Address will be allowed to connect to EZproxy again.
      The disabled account will stay disabled.

### Enable Teams and email notifications

Use the same test username as before.

---

1. Enable Teams webhook in `/usr/local/etc/brick/config.toml`
   - Uncomment webhook URL entry added earlier
1. Reset demo environment
1. Submit two payloads
1. Switch to Teams channel
   - First message indicates payload receipt
   - Second message indicates errors (if any), actions taken
1. Switch to the `syslog entries` tab
   - point out Teams notifications details
1. Switch to the disabled `user entries` tab
   - show that there is still only one entry
1. Switch to the `reported user entries` tab
   - Show familiar entries
1. Switch to the `fail2ban log` tab
   - Show familiar entries
1. Switch to the MailDev container
   - Accessible within demo VM or externally on NAT network
   - Step through all notifications

### Same test username, different IP Address

Do not reset demo environment; we need to have the prior entries/settings
as-is.

---

1. Modify the `contrib/tests/splunk-sanitized-payload-formatted.json` file
   - Replace test account IP Address with `123.123.123.123`
1. Switch to `Splunk (Payload)` tab
   - Submit payload
1. Switch to `syslog entries` tab
1. Switch to `disabled user entries` tab
   - Note that the user account won't be disabled twice
1. Switch to the `reported user entries` tab
   - Note that the "already disabled" message is present as before, but now we
     see a different IP Address than when the user was first banned
1. Switch to the `fail2ban log` tab
   - Mention that the new IP associated with the disabled user account has
     been banned
     - This will cause the associated session to timeout

### Ignored username, different IP Address

1. Modify the `contrib/tests/splunk-sanitized-payload-formatted.json` file
   1. Replace earlier test account with mine
   1. Replace test IP Address with `4.4.4.4`
1. Add my user account to the ignored user file
   - `/usr/local/etc/brick/users.brick-ignored.txt`
   - Should already be open in VS Code from earlier prep step
1. Submit the updated payload
1. Switch to the syslog entries tab
   - Note the Ignored entry there
1. Switch to the disabled user entries tab
   - show that no new entries were added
1. Switch to the reported user entries tab
   - show that an IGNORED entry was added along with sufficient details
     explaining why
   - "ignored" user accounts are "protected" from disable actions (see next)
1. Switch to the fail2ban log tab
   - show that no new bans are present
   - Note: Other requests from Splunk tied to the same IP for other usernames
     won't disable this account, but they will be temporarily banned due to
     the shared IP Address. Once the IP is unbanned, non-disabled accounts can
     use EZproxy as before
1. Switch to Microsoft Teams
   a. Show the notification pair

### Ignored IP Address, different username

1. Modify the `contrib/tests/splunk-sanitized-payload-formatted.json` file
   - Replace my account with another team member's
1. Add IP Address to the ignored IP Address file
   - IP: `10.10.10.10`
     - arbitrary, just writing it out here for reference
   - File: `/usr/local/etc/brick/ips.brick-ignored.txt`
     - Note: should already be open in VS Code from earlier prep step
1. Switch to the payload tab
   - Submit payload
1. Switch to the syslog entries tab
   - Note ignored entry there
1. Switch to the disabled user entries tab
   - show that no new entries were added
1. Switch to the reported user entries tab
   - show that an IGNORED entry was added along with sufficient details
     explaining why
1. Switch to the fail2ban log tab
   - show that no new bans are present
   - Note: all non-disabled user accounts associated with the ignored IP entry
     are "protected" or ignored
1. Switch to Microsoft Teams
   - Show the notification pair

## High-level overview

See [Overview](start-here.md) doc for additional details.

## Improvements

### Recent

- Email notifications directly from `brick`
  - analogue to Teams notifications

- Support for automatic sessions termination
  - using official `ezproxy` binary
  - using unsupported `kill` subcommand of official `ezproxy` binary

### Planned

- Additional endpoints
  - list disabled user accounts
    - user accounts disabled by this application (`users.brick-disabled.txt`)
    - user accounts disabled manually (`users.disabled.txt`)
  - list log messages associated with disabled user accounts

### Potential

Many of these improvements are likely on hold pending investment of
time/support from leadership.

- "Warning-only" behavior based on specific fields or alert name prefixes
  - Disable for high-confidence alerts
  - Warning for anything else
- LDAP
  - Use LDAP package to add disabled users to a Library-specific AD group
  - Update EZproxy (test instance first) to deny login access to users in this
    AD group
- Refactoring to allow `brick` to take general actions from Splunk, Graylog or
  another alert
  - restart a service
  - force a user account to logout & ban them
  - submit API request to Service Now to indicate problems with a system
  - submit API requests to Redmine to add, create or update tickets
- Refactoring to allow `brick` to receive payloads from other services (and
  take specific actions, perhaps branching paths based on the notification
  source or endpoint used (most likely)
  - GitHub
  - Office 365
  - Our own services/tooling
    - e.g., RSS or Atom feed updates monitoring by another could result in a
      payload submission to `brick` which kicks off a pipeline

## Upstream feature requests

The following enhancement requests would help regardless of whether we use
`brick`:

- Add support for dynamically denying access to specified IPs
  - without restarting EZproxy
- Add (official) support to terminate active user sessions via API or other
  external control
  - without restarting EZproxy
  - Note: This support is available as of right now (learned this after the
    May 2020 demo), but OCLC Support has indicated it isn't official and could
    go away without notice.

## References

These are references intended for the audience to review just after a demo
wraps up.

- "Brick" image​
  - <https://www.flickr.com/photos/sfantti/239849911/​>

- Splunk​
  - <https://www.splunk.com/​>

- EZproxy​
  - <https://www.oclc.org/en/ezproxy.html​>

- Fail2ban​
  - <https://www.fail2ban.org/​>

- Brick​
  - <https://github.com/atc0005/brick>

[brick-logo]: ../media/brick-logo-rounded.png "brick project logo"
