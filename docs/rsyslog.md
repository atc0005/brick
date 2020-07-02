<!-- omit in toc -->
# brick: Configure rsyslog action

- [Project README](../README.md)

<!-- omit in toc -->
## Table of contents

- [Overview](#overview)
- [Log files](#log-files)
- [Deployment](#deployment)
- [Test log message routing](#test-log-message-routing)

## Overview

The `contrib/rsyslog/brick.conf` file matches on syslog messages with `brick`
as the program name and writes those messages to the
`/var/log/brick/syslog.log` file. rsyslog is then asked to `stop` logging
those messages elsewhere. Depending on your environment, the `stop` directive
may unintentionally stop forwarding actions which centralize syslog messages
for further processing.

In our environment, we generally forward everything (with rare exception) and
then selectively log locally to application-specific log files. Said another
way, by the time the `stop` directive is processed by rsyslog in our
environment the log messages have already landed in the forward queue for
messages destined for our remote log server(s).

You may need to disable this directive if you encounter problems with messages
not forwarding as expected in your environment.

Tangent:

If you're not already forwarding/centralizing your server log messages, you
should consider doing so. Having log messages centralized in one searchable
location (e.g., Graylog, Splunk, ...) makes troubleshooting much more
effective. With lower cost alternatives such as Graylog, centralized log
management becomes approachable for even lightly staffed IT teams.

## Log files

As noted in the [deploy](deploy.md) doc, unhandled logs will continue to grow
until they become an issue. The `logrotate` utility is the most common tool
used to "rotate" log files on a set schedule so that they don't fill up local
storage. This project provides a `logrotate` configuration snippet/drop-in to
handle rotating the two log files associated with the `brick` application.

See the [deploy](deploy.md) doc for additional log file coverage.

## Deployment

1. `sudo cp -vi contrib/rsyslog/brick.conf /etc/rsyslog.d/`
1. `sudo service rsyslog restart`

## Test log message routing

1. `sudo systemctl restart brick`
   - optional step to force `brick` to emit output that systemd captures and
     makes available to rsyslog (e.g., `imuxsock`, `imjournal`)
1. `sudo tail /var/log/brick/syslog.log`

At this point you should see output from `brick` captured and forwarded on to
the `/var/log/brick/syslog.log` file.
