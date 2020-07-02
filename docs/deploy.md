<!-- omit in toc -->
# brick: Deployment

- [Project README](../README.md)

<!-- omit in toc -->
## Table of contents

- [Order of deployment](#order-of-deployment)
- [Preparation](#preparation)
- [Building `brick`](#building-brick)
- [Deploying `brick`](#deploying-brick)
  - [Assumptions](#assumptions)
  - [Create service account](#create-service-account)
  - [Setup output directories](#setup-output-directories)
  - [Pre-create files](#pre-create-files)
  - [Log files](#log-files)
    - [Summary](#summary)
    - [Deploy logrotate snippet](#deploy-logrotate-snippet)
  - [Deploy `brick` binary](#deploy-brick-binary)
    - [Picking the right binary](#picking-the-right-binary)
    - [Deploying the binary](#deploying-the-binary)
  - [Configure `brick`](#configure-brick)
    - [Start with the template configuration file](#start-with-the-template-configuration-file)
  - [Firewall rules](#firewall-rules)
  - [Configure and deploy systemd unit](#configure-and-deploy-systemd-unit)
  - [Start `brick` now and enable it to start again at boot](#start-brick-now-and-enable-it-to-start-again-at-boot)

## Order of deployment

Technically `brick` can be deployed before or after `fail2ban`, before or
after Splunk alerts are enabled and before or after EZproxy is setup to
monitor a disabled users flat-file, though you'll have to `touch` or otherwise
pre-create the files monitored for both applications.

It is easier if you go ahead and build/deploy `brick` first before making
other changes covered by our documentation.

## Preparation

Backups. While the process of deploying `brick` and related applications is
normally non-destructive, it's a good idea to always assume the worst when
making system-wide changes; installing/using `fail2ban` is arguably a large
change.

If available to you, create a full VMware snapshot, Hyper-V checkpoint, etc.
Being able to quickly undo an "oops" moment is critical for these types of
changes.

## Building `brick`

See the [build](build.md) instructions doc.

## Deploying `brick`

### Assumptions

Mostly for the sake of complete examples, we will assume the following:

- Remote Splunk server is at 192.168.2.2
- EZproxy server is at 192.168.10.20
- `brick` is deployed to the EZproxy server on port `8000`
- `brick` generated files are saved to `/var/cache/brick/`
  - `/var/cache/brick/users.brick-disabled.txt`
- `brick` config files are deployed to `/usr/local/etc/brick/`
  - `/usr/local/etc/brick/config.toml`
    - originally named/found as `contrib/brick/config.example.toml`
  - `/usr/local/etc/brick/ips.brick-ignored.txt`
  - `/usr/local/etc/brick/users.brick-ignored.txt`
- `brick` user and `brick` group have been created
- `brick` runs as the `brick` user
- EZproxy server OS is Ubuntu Linux 16.04 or 18.04 LTS

### Create service account

1. `sudo useradd --system --shell /bin/false brick`

### Setup output directories

1. `sudo mkdir -vp /usr/local/etc/brick`
1. `sudo mkdir -vp /var/log/brick`
1. `sudo chown -Rv brick:syslog /var/log/brick`
1. `sudo chmod -Rv u=rwX,g+rX,o= /var/log/brick`
   - allow service account full access
   - allow rsyslog read-only access to directory and contents
   - deny access to all other non-root user accounts
1. `sudo mkdir -vp /var/cache/brick`
1. `sudo chown -Rv brick:brick /var/cache/brick`

### Pre-create files

1. `sudo touch /var/log/brick/users.brick-reported.log`
   - pre-creating this file satisfies `fail2ban` requirements that the file
     exist before the associated jail is activated
1. `sudo chown -v brick:brick /var/log/brick/users.brick-reported.log`
   - the `fail2ban` daemon runs with elevated privileges and will already be
     able to access this file
   - we explicitly ensure that the `brick` application won't be blocked from
     modifying it
1. `sudo touch /var/log/brick/syslog.log`
   - our rsyslog configuration snippet directs matching log messages here
1. `sudo chmod -v g+rw /var/log/brick/syslog.log`
    - allow rsyslog read-write access to log file

See also the rsyslog-specific instructions in the [rsyslog](rsyslog.md) doc.

### Log files

#### Summary

`brick` generates log output to `stdout`, `stderr` and directly to the
reported users log file. `systemd` handles routing `stdout`, `stderr`
messages, `rsyslog` handles fetching from `systemd` and everything ends up (if
using the provided [rsyslog](rsyslog.md) configuration) in a local log file
for further review (and processing by [fail2ban](fail2ban.md))

Here are those log files:

- `/var/log/brick.syslog.log`
  - written to by `rsyslog`
- `/var/log/brick/users.brick-reported.log`
  - written to by `brick`

These files will continue to grow until they're rotated out, which on most
Linux distros is handled by the `logrotate` utility. Configuring logs for
rotation often involves using an existing file as template and modifying to
match your specific requirements. The same holds true here. The
`contrib/logrotate/brick` logrotate configuration snippet/file handles
rotation for both log files currently generated.

See also the rsyslog-specific instructions in the [rsyslog](rsyslog.md) doc.

#### Deploy logrotate snippet

1. `sudo cp -vi contrib/logrotate/brick /etc/logrotate.d/`

### Deploy `brick` binary

#### Picking the right binary

In addition to all of the supporting files and settings, let's not forget
about deploying the actual brick executable (as an earlier version of these
docs did).

Depending on your target environment, you'll either need the 32-bit or 64-bit
version of the binary that you generated earlier by following the [build
instructions](build.md). Replace the `v0.1.0-0-g721e6d2` pattern below with
the latest available stable version.

| If you see this `uname -m` output | Use the filename with this pattern    | Your EZProxy server has this architecture |
| --------------------------------- | ------------------------------------- | ----------------------------------------- |
| `x86_64`                          | `brick-v0.1.0-0-g721e6d2-linux-amd64` | 64-bit                                    |
| `i686`                            | `brick-v0.1.0-0-g721e6d2-linux-386`   | 32-bit                                    |

For example, if you run `uname -m` on your EZproxy server and get `x86_64` as
the output, you will want to deploy the `brick-v0.1.0-0-g721e6d2-linux-amd64`
binary.

#### Deploying the binary

1. Copy the appropriate binary to `/usr/local/sbin/brick` on the EZproxy server
1. Set permissions on the brick birnary
   - `sudo chmod -v u=rwx,g=rx,o=rx /usr/local/sbin/brick`

If for example you learn that EZproxy is running on a 32-bit Linux
distribution, then you will want to deploy the
`brick-v0.1.0-0-g721e6d2-linux-386` binary to the EZproxy server as
`/usr/local/sbin/brick` and set `0755` permissions on the file.

Replace the `v0.1.0-0-g721e6d2` pattern with the latest available stable
version.

### Configure `brick`

#### Start with the template configuration file

1. Copy the starter/template configuration file from
   `contrib/brick/config.example.toml` and modify accordingly using the
   [configuration](configure.md) guide.
1. Decide whether you will enable automatic sessions termination or use
   `fail2ban`. See the [fail2ban](fail2ban.md) doc and the
   [configuration](configure.md) guide for more information.
1. Set a Microsoft Teams webhook URL to enable Teams channel notifications.
   - Skip this step if you don't use Microsoft Teams.
1. Copy the starter/template "ignore" files and modify accordingly
   - `/usr/local/etc/brick/ips.brick-ignored.txt`
   - `/usr/local/etc/brick/users.brick-ignored.txt`

### Firewall rules

You'll need to update the host firewall on the EZproxy server to permit
connections from Splunk to the chosen IP Address and TCP port where `brick`
will run. If possible, limit access to just the remote system submitting HTTP
requests (as our example shows).

1. `sudo ufw allow proto tcp from 192.168.2.2 to any tcp port 8000`
   - skip this step if you plan to run this application on your system for
     local testing
     - e.g., `localhost:8000`

Results:

1. allow tcp protocol connections
1. from any source port on remote Splunk search head at 192.168.2.2
1. to any IP Address on EZproxy server
1. to tcp port 8000 where `brick` will listen for connections

### Configure and deploy systemd unit

1. Copy the starter/template systemd unit file from
   `contrib/systemd/brick.service` and modify the `ExecStart` line to specify
   the path to a configuration file.
   - e.g., `ExecStart=/usr/local/sbin/brick --config-file
     /usr/local/etc/brick/config.toml`
   - optionally, specify environment variables using a supported method to
     override default settings as needed
   - see the [configuration](configure.md) guide for supported environment
      variables, command-line flags and configuration file settings
1. Copy to `/etc/systemd/system/brick.service`
1. Run `sudo systemctl daemon-reload`

### Start `brick` now and enable it to start again at boot

1. `sudo systemctl start brick`
1. `sudo systemctl enable brick`
