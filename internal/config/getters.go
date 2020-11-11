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

	"github.com/Showmax/go-fqdn"
)

/******************************************************************
	Note: Validation of settings is performed as a later step.
******************************************************************/

// LogLevel returns the user-provided logging level or the default value if
// not provided. CLI flag values take precedence if provided.
func (c Config) LogLevel() string {

	switch {
	case c.cliConfig.Logging.Level != nil:
		return *c.cliConfig.Logging.Level
	case c.fileConfig.Logging.Level != nil:
		return *c.fileConfig.Logging.Level
	default:
		return defaultLogLevel
	}
}

// LogFormat returns the user-provided logging format or the default value if
// not provided. CLI flag values take precedence if provided.
func (c Config) LogFormat() string {

	switch {
	case c.cliConfig.Logging.Format != nil:
		return *c.cliConfig.Logging.Format
	case c.fileConfig.Logging.Format != nil:
		return *c.fileConfig.Logging.Format
	default:
		return defaultLogFormat
	}
}

// LogOutput returns the user-provided logging output or the default value if
// not provided. CLI flag values take precedence if provided.
func (c Config) LogOutput() string {

	switch {
	case c.cliConfig.Logging.Output != nil:
		return *c.cliConfig.Logging.Output
	case c.fileConfig.Logging.Output != nil:
		return *c.fileConfig.Logging.Output
	default:
		return defaultLogOutput
	}
}

// LocalTCPPort returns the user-provided logging format or the default value
// if not provided. CLI flag values take precedence if provided.
func (c Config) LocalTCPPort() int {

	switch {
	case c.cliConfig.Network.LocalTCPPort != nil:
		return *c.cliConfig.Network.LocalTCPPort
	case c.fileConfig.Network.LocalTCPPort != nil:
		return *c.fileConfig.Network.LocalTCPPort
	default:
		return defaultLocalTCPPort
	}
}

// LocalIPAddress returns the user-provided logging format or the default
// value if not provided. CLI flag values take precedence if provided.
func (c Config) LocalIPAddress() string {

	switch {
	case c.cliConfig.Network.LocalIPAddress != nil:
		return *c.cliConfig.Network.LocalIPAddress
	case c.fileConfig.Network.LocalIPAddress != nil:
		return *c.fileConfig.Network.LocalIPAddress
	default:
		return defaultLocalIP
	}
}

// RequireTrustedPayloadSender indicates whether the sysadmin specified a list
// of IP Addresses to trust for payload submission.
func (c Config) RequireTrustedPayloadSender() bool {
	return c.cliConfig.Network.TrustedIPAddresses != nil ||
		c.fileConfig.Network.TrustedIPAddresses != nil
}

// TrustedIPAddresses returns the user-provided list of IP Addresses that
// should be trusted to receive payloads or the the default value if not
// provided. CLI flag values take precedence if provided.
func (c Config) TrustedIPAddresses() []string {
	switch {
	case c.cliConfig.Network.TrustedIPAddresses != nil:
		return c.cliConfig.Network.TrustedIPAddresses
	case c.fileConfig.Network.TrustedIPAddresses != nil:
		return c.fileConfig.Network.TrustedIPAddresses
	default:
		return []string{}
	}
}

// DisabledUsersFile returns the user-provided path to the EZproxy include
// file where this application should write disabled user accounts or the
// default value if not provided. CLI flag values take precedence if provided.
func (c Config) DisabledUsersFile() string {

	switch {
	case c.cliConfig.DisabledUsers.File != nil:
		return *c.cliConfig.DisabledUsers.File
	case c.fileConfig.DisabledUsers.File != nil:
		return *c.fileConfig.DisabledUsers.File
	default:
		// FIXME: During development the default is set to a fixed/temporary
		// path. Before MVP deployment the defaults should be changed to empty
		// strings?
		return defaultDisabledUsersFile
	}
}

// DisabledUsersFilePermissions returns the user-provided permissions for the
// EZproxy include file where this application should write disabled user
// accounts or the default value if not provided. CLI flag values take
// precedence if provided.
func (c Config) DisabledUsersFilePermissions() os.FileMode {

	switch {
	case c.cliConfig.DisabledUsers.FilePermissions != nil:
		return *c.cliConfig.DisabledUsers.FilePermissions
	case c.fileConfig.DisabledUsers.FilePermissions != nil:
		return *c.fileConfig.DisabledUsers.FilePermissions
	default:
		return defaultDisabledUsersFilePerms
	}
}

// ReportedUsersLogFile returns the fully-qualified path to the log file where
// this application should log user disable request events for fail2ban to
// ingest or the default value if not provided. CLI flag values take
// precedence if provided.
func (c Config) ReportedUsersLogFile() string {

	switch {
	case c.cliConfig.ReportedUsers.LogFile != nil:
		return *c.cliConfig.ReportedUsers.LogFile
	case c.fileConfig.ReportedUsers.LogFile != nil:
		return *c.fileConfig.ReportedUsers.LogFile
	default:
		// FIXME: During development the default is set to a
		// fixed/temporary path. Before MVP deployment the defaults
		// should be changed to empty strings?
		return defaultReportedUsersLogFile
	}
}

// ReportedUsersLogFilePermissions returns the user-provided permissions for
// the log file where this application should log user disable request events
// for fail2ban to ingest or the default value if not provided. CLI flag
// values take precedence if provided.
func (c Config) ReportedUsersLogFilePermissions() os.FileMode {

	switch {
	case c.cliConfig.ReportedUsers.LogFilePermissions != nil:
		return *c.cliConfig.ReportedUsers.LogFilePermissions
	case c.fileConfig.ReportedUsers.LogFilePermissions != nil:
		return *c.fileConfig.ReportedUsers.LogFilePermissions
	default:
		return defaultReportedUsersLogFilePerms
	}
}

// IgnoredUsersFile returns the user-provided path to the file containing a
// list of user accounts which should not be disabled and whose associated IP
// should not be banned by this application. If not specified, the default
// value is provided along. CLI flag values take precedence if provided.
func (c Config) IgnoredUsersFile() string {

	switch {
	case c.cliConfig.IgnoredUsers.File != nil:
		return *c.cliConfig.IgnoredUsers.File
	case c.fileConfig.IgnoredUsers.File != nil:
		return *c.fileConfig.IgnoredUsers.File
	default:
		return defaultIgnoredUsersFile
	}
}

// IsSetIgnoredUsersFile indicates whether a user-provided path to the file
// containing a list of user accounts which should not be disabled and whose
// associated IP should not be banned by this application was provided.
// Deprecated: See GH-46
func (c Config) IsSetIgnoredUsersFile() bool {
	switch {
	case c.cliConfig.IgnoredUsers.File != nil:
		return true
	case c.fileConfig.IgnoredUsers.File != nil:
		return true
	default:
		return false
	}
}

// IgnoredIPAddressesFile returns the user-provided path to the file
// containing a list of individual IP Addresses which should not be banned by
// this application. If not specified, the default value is provided.
func (c Config) IgnoredIPAddressesFile() string {
	switch {
	case c.cliConfig.IgnoredIPAddresses.File != nil:
		return *c.cliConfig.IgnoredIPAddresses.File
	case c.fileConfig.IgnoredIPAddresses.File != nil:
		return *c.fileConfig.IgnoredIPAddresses.File
	default:
		return defaultIgnoredIPAddressesFile
	}
}

// IsSetIgnoredIPAddressesFile indicates whether a user-provided path to the
// file containing a list of individual IP Addresses which should not be
// banned by this application was provided.
// Deprecated: See GH-46
func (c Config) IsSetIgnoredIPAddressesFile() bool {
	switch {
	case c.cliConfig.IgnoredIPAddresses.File != nil:
		return true
	case c.fileConfig.IgnoredIPAddresses.File != nil:
		return true
	default:
		return false
	}
}

// ConfigFile returns the user-provided path to the config file for this
// application or the default value if not provided. CLI flag or environment
// variables are the only way to specify a value for this setting.
func (c Config) ConfigFile() string {
	switch {
	case c.cliConfig.ConfigFile != nil:
		return *c.cliConfig.ConfigFile
	default:
		return defaultConfigFile
	}
}

// IgnoreLookupErrors returns the user-provided choice regarding ignoring
// lookup errors or the default value if not provided. CLI flag values take
// precedence if provided.
//
// TODO: See GH-62.
func (c Config) IgnoreLookupErrors() bool {
	switch {
	case c.cliConfig.IgnoreLookupErrors != nil:
		return *c.cliConfig.IgnoreLookupErrors
	case c.fileConfig.IgnoreLookupErrors != nil:
		return *c.fileConfig.IgnoreLookupErrors
	default:
		return defaultIgnoreLookupErrors
	}
}

// DisabledUsersFileEntrySuffix returns the user-provided disabled users entry
// suffix or the default value if not provided. CLI flag values take
// precedence if provided.
func (c Config) DisabledUsersFileEntrySuffix() string {
	// TODO: Set this as a method on the DisabledUsers type instead/also?
	switch {
	case c.cliConfig.DisabledUsers.EntrySuffix != nil:
		return *c.cliConfig.DisabledUsers.EntrySuffix
	case c.fileConfig.DisabledUsers.EntrySuffix != nil:
		return *c.fileConfig.DisabledUsers.EntrySuffix
	default:
		return defaultDisabledUsersFileEntrySuffix
	}
}

// TeamsWebhookURL returns the user-provided webhook URL used for Teams
// notifications or the default value if not provided. CLI flag values take
// precedence if provided.
func (c Config) TeamsWebhookURL() string {

	switch {
	case c.cliConfig.MSTeams.WebhookURL != nil:
		return *c.cliConfig.MSTeams.WebhookURL
	case c.fileConfig.MSTeams.WebhookURL != nil:
		return *c.fileConfig.MSTeams.WebhookURL
	default:
		// FIXME: During development the default is set to a
		// fixed/temporary path. Before MVP deployment the defaults
		// should be changed to empty strings?
		return defaultMSTeamsWebhookURL
	}
}

// TeamsNotificationRateLimit returns a time.Duration value based on the
// user-provided rate limit in seconds between Microsoft Teams notifications
// or the default value if not provided. CLI flag values take precedence if
// provided.
func (c Config) TeamsNotificationRateLimit() time.Duration {
	var rateLimitSeconds int
	switch {
	case c.cliConfig.MSTeams.RateLimit != nil:
		rateLimitSeconds = *c.cliConfig.MSTeams.RateLimit
	case c.fileConfig.MSTeams.RateLimit != nil:
		rateLimitSeconds = *c.fileConfig.MSTeams.RateLimit
	default:
		rateLimitSeconds = defaultMSTeamsRateLimit
	}

	return time.Duration(rateLimitSeconds) * time.Second
}

// TeamsNotificationRetries returns the user-provided retry limit before
// giving up on message delivery or the default value if not provided. CLI
// flag values take precedence if provided.
func (c Config) TeamsNotificationRetries() int {

	switch {
	case c.cliConfig.MSTeams.Retries != nil:
		return *c.cliConfig.MSTeams.Retries
	case c.fileConfig.MSTeams.Retries != nil:
		return *c.fileConfig.MSTeams.Retries
	default:
		return defaultMSTeamsRetries
	}
}

// TeamsNotificationRetryDelay returns the user-provided delay between retry
// attempts for Microsoft Teams notifications or the default value if not
// provided. CLI flag values take precedence if provided.
func (c Config) TeamsNotificationRetryDelay() int {

	switch {
	case c.cliConfig.MSTeams.RetryDelay != nil:
		return *c.cliConfig.MSTeams.RetryDelay
	case c.fileConfig.MSTeams.RetryDelay != nil:
		return *c.fileConfig.MSTeams.RetryDelay
	default:
		return defaultMSTeamsRetryDelay
	}
}

// NotifyTeams indicates whether or not notifications should be sent to a
// Microsoft Teams channel.
func (c Config) NotifyTeams() bool {

	// Assumption: config.validate() has already been called for the existing
	// instance of the Config type and this method is now being called by
	// later stages of the codebase to determine only whether an attempt
	// should be made to send a message to Teams.

	// For now, use the same logic that validate() uses to determine whether
	// validation checks should be run.
	return c.TeamsWebhookURL() != ""

}

// NotifyEmail indicates whether or not notifications should be generated and
// sent via email to specified recipients.
func (c Config) NotifyEmail() bool {

	// Assumption: config.validate() has already been called for the existing
	// instance of the Config type and this method is now being called by
	// later stages of the codebase to determine only whether an attempt
	// should be made to send a message to a SMTP server.

	// For now, use the same logic that validate() uses to determine whether
	// validation checks should be run.
	return c.EmailServer() != ""

}

// EmailNotificationRateLimit returns a time.Duration value based on the
// user-provided rate limit in seconds between email notifications or the
// default value if not provided. CLI flag values take precedence if provided.
func (c Config) EmailNotificationRateLimit() time.Duration {
	var rateLimitSeconds int
	switch {
	case c.cliConfig.Email.RateLimit != nil:
		rateLimitSeconds = *c.cliConfig.Email.RateLimit
	case c.fileConfig.Email.RateLimit != nil:
		rateLimitSeconds = *c.fileConfig.Email.RateLimit
	default:
		rateLimitSeconds = defaultSMTPRateLimit
	}

	return time.Duration(rateLimitSeconds) * time.Second
}

// EmailServer returns the user-provided SMTP server to be used for email
// notifications or the default value if not provided. CLI flag values take
// precedence if provided.
func (c Config) EmailServer() string {
	switch {
	case c.cliConfig.Email.Server != nil:
		return *c.cliConfig.Email.Server
	case c.fileConfig.Email.Server != nil:
		return *c.fileConfig.Email.Server
	default:
		return defaultSMTPServerFQDN
	}
}

// EmailServerPort returns the user-provided TCP port for email notifications
// or the default value if not provided. CLI flag values take precedence if
// provided.
func (c Config) EmailServerPort() int {
	switch {
	case c.cliConfig.Email.Port != nil:
		return *c.cliConfig.Email.Port
	case c.fileConfig.Email.Port != nil:
		return *c.fileConfig.Email.Port
	default:
		return defaultSMTPServerPort
	}
}

// EmailClientIdentity returns the user-provided identity for the server that
// this application sends email notifications on behalf of. If not provided,
// attempt to get the fully-qualified domain name for the system where this
// application is running. If there are issues resolving the fqdn use our
// fallback value. or the default CLI flag values take precedence if provided.
func (c Config) EmailClientIdentity() string {
	switch {
	case c.cliConfig.Email.ClientIdentity != nil:
		return *c.cliConfig.Email.ClientIdentity
	case c.fileConfig.Email.ClientIdentity != nil:
		return *c.fileConfig.Email.ClientIdentity
	default:
		// Since sysadmin did not specify a value, attempt to get
		// fully-qualified domain name for the system where this application
		// is running. If there are issues resolving the fqdn use our fallback
		// value.
		hostname, err := fqdn.FqdnHostname()
		if err != nil {
			hostname = defaultSMTPClientIdentity
		}

		return hostname
	}
}

// EmailSenderAddress returns the user-provided email address used as the
// sender for all outgoing email notifications from this application or the
// default value if not provided. CLI flag values take precedence if provided.
func (c Config) EmailSenderAddress() string {
	switch {
	case c.cliConfig.Email.SenderAddress != nil:
		return *c.cliConfig.Email.SenderAddress
	case c.fileConfig.Email.SenderAddress != nil:
		return *c.fileConfig.Email.SenderAddress
	default:
		return defaultSMTPSenderAddress
	}
}

// EmailRecipientAddresses returns the user-provided list of email addresess
// to receive all outgoing email notifications from this application or the
// default value if not provided. CLI flag values take precedence if provided.
func (c Config) EmailRecipientAddresses() []string {
	switch {
	case c.cliConfig.Email.RecipientAddresses != nil:
		return c.cliConfig.Email.RecipientAddresses
	case c.fileConfig.Email.RecipientAddresses != nil:
		return c.fileConfig.Email.RecipientAddresses
	default:
		// Validation should catch this and fail the start-up attempt IF the
		// sysadmin opted to specify a SMTP mail server, but *not* provide one
		// or more destination addresses. If the SMTP server is not specified,
		// this method (and thus the value) should remain unused.
		return make([]string, 0)
	}
}

// EmailNotificationRetries returns the user-provided retry limit before
// giving up on email message delivery or the default value if not provided.
// CLI flag values take precedence if provided.
func (c Config) EmailNotificationRetries() int {
	switch {
	case c.cliConfig.Email.Retries != nil:
		return *c.cliConfig.Email.Retries
	case c.fileConfig.Email.Retries != nil:
		return *c.fileConfig.Email.Retries
	default:
		return defaultSMTPRetries
	}
}

// EmailNotificationRetryDelay returns the user-provided delay for email
// notifications or the default value if not provided. CLI flag values take
// precedence if provided. This delay is added regardless of whether a
// previous notification delivery attempt has been made.
func (c Config) EmailNotificationRetryDelay() int {
	switch {
	case c.cliConfig.Email.RetryDelay != nil:
		return *c.cliConfig.Email.RetryDelay
	case c.fileConfig.Email.RetryDelay != nil:
		return *c.fileConfig.Email.RetryDelay
	default:
		return defaultSMTPRetryDelay
	}
}

// EZproxyExecutablePath returns the user-provided, fully-qualified path to
// the EZproxy executable or the default value if not provided. CLI flag
// values take precedence if provided.
func (c Config) EZproxyExecutablePath() string {
	switch {
	case c.cliConfig.EZproxy.ExecutablePath != nil:
		return *c.cliConfig.EZproxy.ExecutablePath
	case c.fileConfig.EZproxy.ExecutablePath != nil:
		return *c.fileConfig.EZproxy.ExecutablePath
	default:
		return defaultEZproxyExecutablePath
	}
}

// EZproxyActiveFilePath returns the user-provided, fully-qualified path to
// the EZproxy Active Users and Hosts "state" file or the default value if not
// provided. CLI flag values take precedence if provided.
func (c Config) EZproxyActiveFilePath() string {
	switch {
	case c.cliConfig.EZproxy.ActiveFilePath != nil:
		return *c.cliConfig.EZproxy.ActiveFilePath
	case c.fileConfig.EZproxy.ActiveFilePath != nil:
		return *c.fileConfig.EZproxy.ActiveFilePath
	default:
		return defaultEZproxyActiveFilePath
	}
}

// EZproxyAuditFileDirPath returns the user-provided, fully-qualified path to
// the EZproxy audit files directory or the default value if not provided. CLI
// flag values take precedence if provided.
func (c Config) EZproxyAuditFileDirPath() string {
	switch {
	case c.cliConfig.EZproxy.AuditFileDirPath != nil:
		return *c.cliConfig.EZproxy.AuditFileDirPath
	case c.fileConfig.EZproxy.AuditFileDirPath != nil:
		return *c.fileConfig.EZproxy.AuditFileDirPath
	default:
		return defaultEZproxyAuditFileDirPath
	}
}

// EZproxySearchRetries returns the user-provided number of retry attempts to
// make for session lookup attempts that return zero results or the default
// value if not provided. CLI flag values take precedence if provided.
func (c Config) EZproxySearchRetries() int {
	switch {
	case c.cliConfig.EZproxy.SearchRetries != nil:
		return *c.cliConfig.EZproxy.SearchRetries
	case c.fileConfig.EZproxy.SearchRetries != nil:
		return *c.fileConfig.EZproxy.SearchRetries
	default:
		return defaultEZproxySearchRetries
	}
}

// EZproxySearchDelay returns the user-provided number of seconds between
// session lookup attempts or the default value if not provided. CLI flag
// values take precedence if provided.
func (c Config) EZproxySearchDelay() int {
	switch {
	case c.cliConfig.EZproxy.SearchDelay != nil:
		return *c.cliConfig.EZproxy.SearchDelay
	case c.fileConfig.EZproxy.SearchDelay != nil:
		return *c.fileConfig.EZproxy.SearchDelay
	default:
		return defaultEZproxySearchDelay
	}
}

// EZproxyTerminateSessions indicates whether attempts should be made to
// terminate sessions for reported user accounts. The user-provided value is
// returned or the default value if not provided. CLI flag values take
// precedence if provided.
func (c Config) EZproxyTerminateSessions() bool {
	switch {
	case c.cliConfig.EZproxy.TerminateSessions != nil:
		return *c.cliConfig.EZproxy.TerminateSessions
	case c.fileConfig.EZproxy.TerminateSessions != nil:
		return *c.fileConfig.EZproxy.TerminateSessions
	default:
		return defaultEZproxyTerminateSessions
	}
}
