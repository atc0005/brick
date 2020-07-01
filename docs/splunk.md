<!-- omit in toc -->
# brick: Integrating with Splunk

- [Project README](../README.md)

<!-- omit in toc -->
## Table of Contents

- [Overview](#overview)
- [Directions](#directions)
- [Payload schema / format](#payload-schema--format)
  - [Official example](#official-example)
  - [What we use](#what-we-use)

## Overview

If drawn on a whiteboard, Splunk is likely somewhere off to the side.

EZproxy is in the center, `fail2ban` and `brick` are nearby (or in the same
sphere) and other resources (such as maybe Microsoft Teams and email) are off
somewhere to another side. This quick sketch for illustration wouldn't
tell the whole picture however, and would greatly undervalue the role Splunk
plays in resource abuse "management" (it's hard to say with a straight face
that it can ever be completely prevented).

By taking the time to implement and refine alerts, Splunk will become an
invaluable tool to monitor and report abusive activity to your sysadmin team
responsible for managing your EZproxy server(s).

Splunk is also the primary dependency of this application. While this may
change in the future (e.g., Graylog testing), without Splunk, this current
iteration of the application serves little purpose.

## Directions

Once you have created, tested and refined one or more email-based alerts, you
are ready to begin using webhook payloads to report problematic user accounts
to `brick`.

1. Review the [official documentation](references.md) for setting up a
   "webhook alert action"
1. Follow those instructions and set the target webhook URL
   1. we'll assume that your EZproxy server has a FQDN of ezproxy.example.com
      and is normally accessible at <https://ezproxy.example.com/>
   2. using the [endpoints](endpoints.md) doc as our guide, set
      <https://ezproxy.example.com:8000/api/v1/users/disable> as the webhook
      URL
1. As noted in the [deploy](deploy.md) doc, make sure you have a firewall
   rule in place to limit payload delivery to the `disable` endpoint to only
   your Splunk server and any trusted SysAdmin / IT Support team members
1. If you haven't already done so, [build](build.md), [deploy](deploy.md) and
   [configure](configure.md) the `brick` application
1. The same goes for `fail2ban`, if you haven't yet, install and
   [configure](configure.md) [fail2ban](fail2ban.md) *or* exclusively use
   EZproxy's native (unofficial) support for terminating user sessions.
1. Test!

## Payload schema / format

### Official example

The example payload listed in the [official documentation](references.md) is
provided below.

```json
{
   "result": {
      "sourcetype" : "mongod",
      "count" : "8"
   },
   "sid" : "scheduler_admin_search_W2_at_14232356_132",
   "results_link" : "http://web.example.local:8000/app/search/@go?sid=scheduler_admin_search_W2_at_14232356_132",
   "search_name" : null,
   "owner" : "admin",
   "app" : "search"
}
```

Description:

> **Webhook data payload**
>
> The webhook POST request's JSON data payload includes the following details.
>
> - Search ID or SID for the saved search that triggered the alert
> - Link to search results
> - Search owner and app
> - First result row from the triggering search results

This example payload can also be found in the
[contrib/tests/splunk-test-submission.json](../contrib/tests/splunk-test-submission.json)
file.

### What we use

During testing we captured several iterations of the payload for comparison.
Here is a sanitized payload that was used when developing and testing this
application:

```json
{
    "results_link": "https://splunk.example.com:8000/app/search/@go?sid=scheduler__abc0001__search__RMD5267e440bddd8ef1f_at_1581522000_11805",
    "result": {
        "contextData": "",
        "date_month": "february",
        "forcecdn": "",
        "http_status_code": "200",
        "_bkt": "ezproxy-http~158~A31E767E-7887-4A15-86F1-2CB85DB0F805",
        "_indextime": "1581517935",
        "_kv": "1",
        "linecount": "1",
        "ezproxy_time": "12/Feb/2020:08:32:15 -0600",
        "_serial": "0",
        "_time": "1581517935",
        "eventtype": "lib_events",
        "_eventtype_color": "none",
        "sp": "",
        "bhskip": "",
        "_sourcetype": "ezproxy-http",
        "punct": "..._-__[//:::_-]_\"_://...:///_/.\"___\"/._(__.;_;_)_",
        "splunk_server": "splunk-index9000",
        "session": "",
        "host": "ezproxy",
        "url": "",
        "srcip": "192.168.2.3",
        "splunk_server_group": "",
        "_cd": "158:1576585",
        "_si": [
            "splunk-index9000",
            "ezproxy-http"
        ],
        "tag::eventtype": "library",
        "timestartpos": "25",
        "date_hour": "8",
        "date_second": "15",
        "timeendpos": "51",
        "username": "abc0001",
        "date_minute": "32",
        "date_mday": "12",
        "index": "ezproxy-http",
        "timeZoneId": "",
        "sourcetype": "ezproxy-http",
        "rs": "",
        "transitionType": "",
        "date_wday": "wednesday",
        "source": "/path/to/traffic/log/file.txt",
        "date_zone": "-360",
        "tag": "library",
        "date_year": "2020",
        "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.100 Safari/537.36",
        "_raw": "192.168.2.3 - abc0001 [12/Feb/2020:08:32:15 -0600] \"POST https://1.vendor.example.com:443/V1/Session/ExtendSessionActiveBrowser HTTP/1.1\" 200 0 \"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.100 Safari/537.36\"",
        "URL": "https://1.vendor.example.com:443/V1/Session/ExtendSessionActiveBrowser",
        "vr": ""
    },
    "sid": "scheduler__abc0001__search__RMD5267e440bddd8ef1f_at_1581522000_11805",
    "owner": "abc0001",
    "app": "search",
    "search_name": "TEST - Webhook - Echo"
}
```

`brick` parses and uses select fields from that payload for its work. As time
permits, we hope to further refine the delivered payload to exclude unwanted
fields and bring in new ones.

Files:

- [contrib/tests/splunk-sanitized-payload-unformatted.json](../contrib/tests/splunk-sanitized-payload-unformatted.json)
- [contrib/tests/splunk-sanitized-payload-formatted.json](../contrib/tests/splunk-sanitized-payload-formatted.json)
