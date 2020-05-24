<!-- omit in toc -->
# brick: References

- [Project README](../README.md)

<!-- omit in toc -->
## Table of contents

- [Overview](#overview)
- [References](#references)
  - [Dependencies](#dependencies)
  - [Instruction / Examples](#instruction--examples)
  - [Related projects](#related-projects)

## Overview

The links below are for resources that were found to be useful (if not
absolutely essential) while developing this application.

## References

### Dependencies

- `make` on Windows
  - <https://stackoverflow.com/questions/32127524/how-to-install-and-use-make-in-windows>
- `gcc` on Windows
  - <https://en.wikipedia.org/wiki/MinGW>
  - <http://mingw-w64.org/>
  - <https://www.msys2.org/>

- Libraries/packages
  - <https://github.com/alexflint/go-arg>
  - <https://github.com/apex/log>
  - <https://github.com/pelletier/go-toml>
  - <https://github.com/atc0005/send2teams>
  - `go-teams-notify`
    - (upstream) <https://github.com/dasrick/go-teams-notify>
    - (fork) <https://github.com/atc0005/go-teams-notify>

- External
  - [Splunk](https://www.splunk.com/​)
  - [EZproxy](https://www.oclc.org/en/ezproxy.html​)
  - [Fail2ban](https://www.fail2ban.org/​)

### Instruction / Examples

- EZproxy
  - <https://help.oclc.org/Library_Management/EZproxy/Authenticate_users/EZproxy_authentication_methods/Text_file_authentication>
  - <https://help.oclc.org/Library_Management/EZproxy/Get_started/Join_the_EZproxy_listserv_and_Community_Center>

- Logging
  - <https://github.com/apex/log>
  - <https://brandur.org/logfmt>

- Splunk / JSON payload
  - [Splunk Enterprise (v8.0.1) > Alerting Manual > Use a webhook alert action](https://docs.splunk.com/Documentation/Splunk/8.0.1/Alert/Webhooks)
  - [Splunk Enterprise > Getting Data In > How timestamp assignment works > How Splunk software assigns timestamps](https://docs.splunk.com/Documentation/Splunk/latest/Data/HowSplunkextractstimestamps)

- fail2ban
  - <https://www.the-art-of-web.com/system/fail2ban-filters/>

- systemd
  - <https://www.freedesktop.org/software/systemd/man/systemd.service.html>
  - <https://www.freedesktop.org/software/systemd/man/systemd.unit.html>
  - <https://vincent.bernat.ch/en/blog/2018-systemd-golang-socket-activation>
  - <https://vincent.bernat.ch/en/blog/2017-systemd-golang>
  - <https://fabianlee.org/2017/05/21/golang-running-a-go-binary-as-a-systemd-service-on-ubuntu-16-04/>

### Related projects

- <https://github.com/atc0005/bounce>
- <https://github.com/atc0005/send2teams>
