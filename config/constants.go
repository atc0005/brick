// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/brick
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"os"
	"time"
)

const (

	// MyAppName is the public name of this application
	MyAppName string = "brick"

	// MyAppURL is the location of the repo for this application
	MyAppURL string = "https://github.com/atc0005/brick"

	// MyAppDescription is the description for this application shown in
	// HelpText output.
	MyAppDescription string = "Automatically disable EZproxy users via webhook requests"
)

// Default (flag, config file, etc) settings if not overridden by user input
const (
	defaultLocalTCPPort int    = 8000
	defaultLocalIP      string = "localhost"
	defaultLogLevel     string = "info"
	defaultLogOutput    string = "stdout"
	defaultLogFormat    string = "text"

	// This is appended to each username as it is written to the file in order
	// for EZproxy to treat the user account as ineligible to login
	defaultDisabledUsersFileEntrySuffix string      = "::deny"
	defaultDisabledUsersFile            string      = "/var/cache/brick/users.brick-disabled.txt"
	defaultDisabledUsersFilePerms       os.FileMode = 0o644

	defaultReportedUsersLogFile      string      = "/var/log/brick/users.brick-reported.log"
	defaultReportedUsersLogFilePerms os.FileMode = 0o644
	defaultIgnoredUsersFile          string      = "/usr/local/etc/brick/users.brick-ignored.txt"
	defaultIgnoredIPAddressesFile    string      = "/usr/local/etc/brick/ips.brick-ignored.txt"

	defaultIgnoreLookupErrors bool = true

	// No assumptions can be safely made here; user has to supply this
	defaultMSTeamsWebhookURL string = ""

	// the number of attempts to deliver messages before giving up; applies to
	// Microsoft Teams notifications only
	defaultMSTeamsRetries int = 2

	// the number of seconds to wait between retry attempts; applies to
	// Microsoft Teams notifications only
	defaultMSTeamsDelay int = 5
)

// TODO: Expose these settings via flags, config file
//
// Timeout settings applied to our instance of http.Server
const (
	HTTPServerReadHeaderTimeout time.Duration = 20 * time.Second
	HTTPServerReadTimeout       time.Duration = 1 * time.Minute
	HTTPServerWriteTimeout      time.Duration = 2 * time.Minute
)

// ReadHeaderTimeout:

// HTTPServerShutdownTimeout is used by the graceful shutdown process to
// control how long the shutdown process should wait before forcefully
// terminating.
const HTTPServerShutdownTimeout time.Duration = 30 * time.Second

// NotifyMgrServicesShutdownTimeout is used by the NotifyMgr to determine how
// long it should wait for results from each notifier or notifier "service"
// before continuing on with the shutdown process.
const NotifyMgrServicesShutdownTimeout time.Duration = 2 * time.Second

// Timing-related settings (delays, timeouts) used by our notification manager
// and child goroutines to concurrently process notification requests.
const (

	// NotifyMgrTeamsTimeout is the timeout setting applied to each Microsoft
	// Teams notification attempt. This value does NOT take into account the
	// number of configured retries and retry delays. The final value timeout
	// applied to each notification attempt should be based on those
	// calculations. The GetTimeout method does just that.
	NotifyMgrTeamsTimeout time.Duration = 10 * time.Second

	// NotifyMgrTeamsSendAttemptTimeout

	// NotifyMgrEmailTimeout is the timeout setting applied to each email
	// notification attempt. This value does NOT take into account the number
	// of configured retries and retry delays. The final value timeout applied
	// to each notification attempt should be based on those calculations. The
	// GetTimeout method does just that.
	NotifyMgrEmailTimeout time.Duration = 30 * time.Second

	// NotifyStatsMonitorDelay limits notification stats logging to no more
	// often than this duration. This limiter is to keep from logging the
	// details so often that the information simply becomes noise.
	NotifyStatsMonitorDelay time.Duration = 5 * time.Minute

	// NotifyQueueMonitorDelay limits notification queue stats logging to no
	// more often than this duration. This limiter is to keep from logging the
	// details so often that the information simply becomes noise.
	NotifyQueueMonitorDelay time.Duration = 15 * time.Second

	// NotifyMgrTeamsNotificationDelay is the delay between Microsoft Teams
	// notification attempts. This delay is intended to help prevent
	// unintentional abuse of remote services.
	NotifyMgrTeamsNotificationDelay time.Duration = 5 * time.Second

	// NotifyMgrEmailNotificationDelay is the delay between email notification
	// attempts. This delay is intended to help prevent unintentional abuse of
	// remote services.
	NotifyMgrEmailNotificationDelay time.Duration = 5 * time.Second
)

// NotifyMgrQueueDepth is the number of items allowed into the queue/channel
// at one time. Senders with items for the notification "pipeline" that do not
// fit within the allocated space will block until space in the queue opens.
// Best practice for channels advocates that a smaller number is better than a
// larger one, so YMMV if this is set either too high or too low.
//
// Brief testing (as of this writing) shows that a depth as low as 1 works for
// our purposes, but results in a greater number of stalled goroutines waiting
// to place items into the queue.
const NotifyMgrQueueDepth int = 5

// TCP port ranges
// http://www.iana.org/assignments/port-numbers
// Port numbers are assigned in various ways, based on three ranges: System
// Ports (0-1023), User Ports (1024-49151), and the Dynamic and/or Private
// Ports (49152-65535)
const (
	TCPReservedPort            int = 0
	TCPSystemPortStart         int = 1
	TCPSystemPortEnd           int = 1023
	TCPUserPortStart           int = 1024
	TCPUserPortEnd             int = 49151
	TCPDynamicPrivatePortStart int = 49152
	TCPDynamicPrivatePortEnd   int = 65535
)

// Log levels
const (
	// https://godoc.org/github.com/apex/log#Level

	// LogLevelFatal is used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	LogLevelFatal string = "fatal"

	// LogLevelError is for errors that should definitely be noted.
	LogLevelError string = "error"

	// LogLevelWarn is for non-critical entries that deserve eyes.
	LogLevelWarn string = "warn"

	// LogLevelInfo is for general application operational entries.
	LogLevelInfo string = "info"

	// LogLevelDebug is for debug-level messages and is usually enabled
	// when debugging. Very verbose logging.
	LogLevelDebug string = "debug"
)

// 	apex/log Handlers
// ---------------------------------------------------------
// cli - human-friendly CLI output
// discard - discards all logs
// es - Elasticsearch handler
// graylog - Graylog handler
// json - JSON output handler
// kinesis - AWS Kinesis handler
// level - level filter handler
// logfmt - logfmt plain-text formatter
// memory - in-memory handler for tests
// multi - fan-out to multiple handlers
// papertrail - Papertrail handler
// text - human-friendly colored output
// delta - outputs the delta between log calls and spinner
const (
	// LogFormatCLI provides human-friendly CLI output
	LogFormatCLI string = "cli"

	// LogFormatJSON provides JSON output
	LogFormatJSON string = "json"

	// LogFormatLogFmt provides logfmt plain-text output
	LogFormatLogFmt string = "logfmt"

	// LogFormatText provides human-friendly colored output
	LogFormatText string = "text"

	// LogFormatDiscard discards all logs
	LogFormatDiscard string = "discard"
)

const (

	// LogOutputStdout represents os.Stdout
	LogOutputStdout string = "stdout"

	// LogOutputStderr represents os.Stderr
	LogOutputStderr string = "stderr"
)
