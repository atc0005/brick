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
	"fmt"
	"net"

	"github.com/apex/log"
	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
)

// validateEmailAddress receives a string representing an email address and
// the intent or purpose for the email address (e.g., "sender", "recipient")
// and validates the email address. A error message indicating the reason for
// validation failure is returned or nil if no issues were found.
func validateEmailAddress(emailAddr string, purpose string) error {

	// https://golangcode.com/validate-an-email-address/
	// https://www.w3.org/TR/2016/REC-html51-20161101/sec-forms.html#email-state-typeemail
	// https://stackoverflow.com/a/574698
	// https://www.rfc-editor.org/errata_search.php?rfc=3696

	switch {
	case emailAddr == "":
		notSpecifiedErr := fmt.Errorf(
			"%s email address not specified",
			purpose,
		)
		return notSpecifiedErr

	case (len(emailAddr) < 3 || len(emailAddr) > 254):
		invalidLengthErr := fmt.Errorf(
			"%s email address %q has invalid length of %d",
			purpose,
			emailAddr,
			len(emailAddr),
		)
		return invalidLengthErr

	case !emailRegex.MatchString(emailAddr):
		invalidRegexExMatch := fmt.Errorf(
			"%s email address %q does not match expected W3C format",
			purpose,
			emailAddr,
		)
		return invalidRegexExMatch
	}

	return nil
}

// validateEmailAddresses receives a slice of email addresses and the
// associated intent (e.g., "sender", "recipient") and validates each of them
// returning the validation error if present or nil if no validation failures
// occur.
func validateEmailAddresses(emailAddresses []string, purpose string) error {

	if len(emailAddresses) == 0 {

		return fmt.Errorf(
			"%s email address not specified",
			purpose,
		)
	}

	for _, emailAddr := range emailAddresses {
		if err := validateEmailAddress(emailAddr, purpose); err != nil {
			log.Debug(err.Error())
			return err
		}
	}

	return nil
}

// validate confirms that all config struct fields have reasonable values
func validate(c Config) error {

	switch {

	// WARNING: User opted to use a privileged system port
	case (c.LocalTCPPort() >= TCPSystemPortStart) && (c.LocalTCPPort() <= TCPSystemPortEnd):

		log.Debugf(
			"unprivileged system port %d chosen. ports between %d and %d require elevated privileges",
			c.LocalTCPPort(),
			TCPSystemPortStart,
			TCPSystemPortEnd,
		)

		// log at WARNING level
		log.Warnf(
			"Binding to a port < %d requires elevated permissions. If you encounter errors with this application, please re-run this application and specify a port number between %d and %d",
			TCPUserPortStart,
			TCPUserPortStart,
			TCPUserPortEnd,
		)

	// OK: User opted to use a valid and non-privileged port number
	case (c.LocalTCPPort() >= TCPUserPortStart) && (c.LocalTCPPort() <= TCPUserPortEnd):
		log.Debugf(
			"Valid, non-privileged user port between %d and %d configured: %d",
			TCPUserPortStart,
			TCPUserPortEnd,
			c.LocalTCPPort(),
		)

	// WARNING: User opted to use a dynamic or private TCP port
	case (c.LocalTCPPort() >= TCPDynamicPrivatePortStart) && (c.LocalTCPPort() <= TCPDynamicPrivatePortEnd):
		log.Warnf(
			"WARNING: Valid, non-privileged, but dynamic/private port between %d and %d configured. This range is reserved for dynamic (usually outgoing) connections. If you encounter errors with this application, please re-run this application and specify a port number between %d and %d",
			TCPDynamicPrivatePortStart,
			TCPDynamicPrivatePortEnd,
			TCPUserPortStart,
			TCPUserPortEnd,
		)

	default:
		log.Debugf("invalid port %d specified", c.LocalTCPPort())
		return fmt.Errorf(
			"port %d is not a valid TCP port for this application",
			c.LocalTCPPort(),
		)
	}

	if c.LocalIPAddress() == "" {
		return fmt.Errorf("local IP Address not provided")
	}

	// true if sysadmin specified a value via CLI or config file
	if c.RequireTrustedPayloadSender() {
		switch {
		case len(c.TrustedIPAddresses()) < 1:
			return fmt.Errorf("empty list of trusted IP Addresses provided")
		default:
			for _, ipAddr := range c.TrustedIPAddresses() {
				if net.ParseIP(ipAddr) == nil {
					return fmt.Errorf(
						"invalid IP Address %q provided for trusted IPs list",
						ipAddr,
					)
				}
			}
		}
	}

	switch c.LogLevel() {
	case LogLevelFatal:
	case LogLevelError:
	case LogLevelWarn:
	case LogLevelInfo:
	case LogLevelDebug:
	default:
		return fmt.Errorf("invalid option %q provided for log level",
			c.LogLevel())
	}

	switch c.LogOutput() {
	case LogOutputStderr:
	case LogOutputStdout:
	default:
		return fmt.Errorf("invalid option %q provided for log output",
			c.LogOutput())
	}

	switch c.LogFormat() {
	case LogFormatCLI:
	case LogFormatJSON:
	case LogFormatLogFmt:
	case LogFormatText:
	case LogFormatDiscard:
	default:
		return fmt.Errorf("invalid option %q provided for log format",
			c.LogFormat())
	}

	// TODO: Decide on how to best validate. If we offer default values, can
	// we also logically enforce that something be specified?
	//
	if c.DisabledUsersFile() == "" {
		return fmt.Errorf("path to disabled users file not provided")
	}

	if c.ReportedUsersLogFile() == "" {
		return fmt.Errorf("path to reported users log file not provided")
	}

	// Verify that the user did not opt to set an empty string as the value,
	// otherwise we fail the config validation by returning an error.
	// DEPRECATED: See GH-46
	if c.IsSetIgnoredUsersFile() && c.IgnoredUsersFile() == "" {
		return fmt.Errorf("empty path to ignored users file provided")
	}

	// Verify that the user did not opt to set an empty string as the value,
	// otherwise we fail the config validation by returning an error.
	// DEPRECATED: See GH-46
	if c.IsSetIgnoredIPAddressesFile() && c.IgnoredIPAddressesFile() == "" {
		return fmt.Errorf("empty path to ignored ip addresses file provided")
	}

	// Not having a webhook URL is a valid choice. Perform validation if value
	// is provided.
	if c.TeamsWebhookURL() != "" {
		log.Debugf("Microsoft Teams WebhookURL provided: %v", c.TeamsWebhookURL())

		// Create Microsoft Teams client
		mstClient := goteamsnotify.NewTeamsClient()

		if err := mstClient.ValidateWebhook(c.TeamsWebhookURL()); err != nil {
			return fmt.Errorf("webhook URL validation failed: %w", err)
		}
	}

	if c.TeamsNotificationRateLimit() < 0 {
		log.Debugf("unsupported rate limit specified for MS Teams notifications: %d ", c.TeamsNotificationRateLimit())
		return fmt.Errorf(
			"invalid rate limit specified for MS Teams notifications: %d",
			c.TeamsNotificationRateLimit(),
		)
	}

	if c.TeamsNotificationRetryDelay() < 0 {
		log.Debugf("unsupported retry delay specified for MS Teams notifications: %d ", c.TeamsNotificationRetryDelay())
		return fmt.Errorf(
			"invalid retry delay specified for MS Teams notifications: %d",
			c.TeamsNotificationRetryDelay(),
		)
	}

	if c.TeamsNotificationRetries() < 0 {
		log.Debugf("unsupported retry limit specified for MS Teams notifications: %d ", c.TeamsNotificationRetries())
		return fmt.Errorf(
			"invalid retries limit specified for MS Teams notifications: %d",
			c.TeamsNotificationRetries(),
		)
	}

	// Not specifying an email server is a valid choice. Perform validation of
	// this and other related values if the server name is provided.
	if c.EmailServer() != "" {

		log.Debugf("Email server provided: %q", c.EmailServer())

		// TODO: This needs more work to be truly useful
		if len(c.EmailServer()) <= 1 {
			log.Debugf(
				"unsupported email server name of %d characters specified for email notifications: %q",
				len(c.EmailServer()),
				c.EmailServer(),
			)
			return fmt.Errorf(
				"invalid server name specified for email notifications: %q",
				c.EmailServer(),
			)
		}

		if !(c.EmailServerPort() >= TCPSystemPortStart) && (c.EmailServerPort() <= TCPDynamicPrivatePortEnd) {
			log.Debugf("invalid port %d specified", c.EmailServerPort())
			return fmt.Errorf(
				"port %d is not a valid TCP port for the destination SMTP server",
				c.EmailServerPort(),
			)
		}

		if err := validateEmailAddress(c.EmailSenderAddress(), "sender"); err != nil {
			log.Debug(err.Error())
			return err
		}

		if err := validateEmailAddresses(c.EmailRecipientAddresses(), "recipient"); err != nil {
			log.Debug(err.Error())
			return err
		}

		if c.EmailClientIdentity() != "" {
			if len(c.EmailClientIdentity()) <= 1 {
				log.Debugf(
					"unsupported email client identity of %d characters specified for email server connection: %q",
					len(c.EmailClientIdentity()),
					c.EmailClientIdentity(),
				)
				return fmt.Errorf(
					"unsupported email client identity specified for email server connection: %q",
					c.EmailClientIdentity(),
				)
			}
		}

		if c.EmailNotificationRateLimit() < 0 {
			log.Debugf(
				"unsupported rate limit specified for email notifications: %d ",
				c.EmailNotificationRateLimit(),
			)
			return fmt.Errorf(
				"invalid rate limit specified for email notifications: %d",
				c.EmailNotificationRateLimit(),
			)
		}

		if c.EmailNotificationRetryDelay() < 0 {
			log.Debugf(
				"unsupported delay specified for email notifications: %d ",
				c.EmailNotificationRetryDelay(),
			)
			return fmt.Errorf(
				"invalid delay specified for email notifications: %d",
				c.EmailNotificationRetryDelay(),
			)
		}

		if c.EmailNotificationRetries() < 0 {
			log.Debugf(
				"unsupported retry limit specified for email notifications: %d",
				c.EmailNotificationRetries(),
			)
			return fmt.Errorf(
				"invalid retries limit specified for email notifications: %d",
				c.EmailNotificationRetries(),
			)
		}

	}

	if c.EZproxyExecutablePath() == "" {
		return fmt.Errorf("path to EZproxy executable file not provided")
	}

	if c.EZproxyActiveFilePath() == "" {
		return fmt.Errorf("path to EZproxy active users state file not provided")
	}

	if c.EZproxyAuditFileDirPath() == "" {
		return fmt.Errorf("path to EZproxy audit file directory not provided")
	}

	if c.EZproxySearchDelay() < 0 {
		log.Debugf("unsupported delay specified for EZproxy session lookup attempts: %d ", c.EZproxySearchDelay())
		return fmt.Errorf(
			"invalid delay specified for EZproxy session lookup attempts: %d",
			c.EZproxySearchDelay(),
		)
	}

	if c.EZproxySearchRetries() < 0 {
		log.Debugf("unsupported retry limit specified for EZproxy session lookup attempts: %d ", c.EZproxySearchRetries())
		return fmt.Errorf(
			"invalid retries limit specified for EZproxy session lookup attempts: %d",
			c.EZproxySearchRetries(),
		)
	}

	// if we made it this far then we signal all is well
	return nil

}
