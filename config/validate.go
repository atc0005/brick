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

	"github.com/apex/log"

	// use our fork for now until recent work can be submitted for inclusion
	// in the upstream project
	goteamsnotify "github.com/atc0005/go-teams-notify"
	send2teams "github.com/atc0005/send2teams/teams"
)

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

		// TODO: Do we really need both of these?
		if ok, err := goteamsnotify.IsValidWebhookURL(c.TeamsWebhookURL()); !ok {
			return err
		}

		if err := send2teams.ValidateWebhook(c.TeamsWebhookURL()); err != nil {
			return err
		}

	}

	if c.TeamsNotificationDelay() < 0 {
		log.Debugf("unsupported delay specified for MS Teams notifications: %d ", c.TeamsNotificationDelay())
		return fmt.Errorf(
			"invalid delay specified for MS Teams notifications: %d",
			c.TeamsNotificationDelay(),
		)
	}

	if c.TeamsNotificationRetries() < 0 {
		log.Debugf("unsupported retry limit specified for MS Teams notifications: %d ", c.TeamsNotificationRetries())
		return fmt.Errorf(
			"invalid retries limit specified for MS Teams notifications: %d",
			c.TeamsNotificationRetries(),
		)
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
