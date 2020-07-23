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
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/apex/log"
	"github.com/atc0005/brick/internal/caller"
	"github.com/pelletier/go-toml"
)

// Version is updated via Makefile builds by referencing the fully-qualified
// path to this variable, including the package. We set a placeholder value so
// that something resembling a version string will be provided for
// non-Makefile builds.
var Version string = "x.y.z"

func (c *Config) String() string {
	return fmt.Sprintf(
		"UnifiedConfig: { "+
			"Network.LocalTCPPort: %v, "+
			"Network.LocalIPAddress: %v, "+
			"Logging.Level: %s, "+
			"Logging.Output: %s, "+
			"Logging.Format: %s, "+
			"DisabledUsers.File: %s, "+
			"DisabledUsers.EntrySuffix: %s, "+
			"DisabledUsers.FilePermissions: %v, "+
			"ReportedUsers.LogFile: %q, "+
			"ReportedUsers.LogFilePermissions: %v, "+
			"IgnoredUsers.File: %q, "+
			"IsSetIgnoredUsersFile: %t, "+
			"IgnoredIPAddresses.File: %q, "+
			"IsSetIgnoredIPAddressesFile: %t, "+
			"IgnoreLookupErrors: %t, "+
			"MSTeams.WebhookURL: %q, "+
			"MSTeams.RateLimit: %v, "+
			"MSTeams.Retries: %v, "+
			"MSTeams.RetryDelay: %v, "+
			"NotifyTeams: %t, "+
			"NotifyEmail: %t, "+
			"Email.Server: %q, "+
			"Email.Port: %v, "+
			"Email.ClientIdentity: %q, "+
			"Email.SenderAddress: %q, "+
			"Email.RecipientAddresses: %v, "+
			"Email.RateLimit: %v, "+
			"Email.Retries: %v, "+
			"Email.RetryDelay: %v, "+
			"EZproxy.ExecutablePath: %v, "+
			"EZproxy.ActiveFilePath: %v, "+
			"EZproxy.AuditFileDirPath: %v, "+
			"EZproxy.SearchRetries: %v, "+
			"EZproxy.SearchDelay: %v, "+
			"EZproxy.TerminateSessions: %t, "+
			"ConfigFile: %q}",
		c.LocalTCPPort(),
		c.LocalIPAddress(),
		c.LogLevel(),
		c.LogOutput(),
		c.LogFormat(),
		c.DisabledUsersFile(),
		c.DisabledUsersFileEntrySuffix(),
		c.DisabledUsersFilePermissions(),
		c.ReportedUsersLogFile(),
		c.ReportedUsersLogFilePermissions(),
		c.IgnoredUsersFile(),
		c.IsSetIgnoredUsersFile(),
		c.IgnoredIPAddressesFile(),
		c.IsSetIgnoredIPAddressesFile(),
		c.IgnoreLookupErrors(),
		c.TeamsWebhookURL(),
		c.TeamsNotificationRateLimit(),
		c.TeamsNotificationRetries(),
		c.TeamsNotificationRetryDelay(),
		c.NotifyTeams(),
		c.NotifyEmail(),
		c.EmailServer(),
		c.EmailServerPort(),
		c.EmailClientIdentity(),
		c.EmailSenderAddress(),
		c.EmailRecipientAddresses(),
		c.EmailNotificationRateLimit(),
		c.EmailNotificationRetries(),
		c.EmailNotificationRetryDelay(),
		c.EZproxyExecutablePath(),
		c.EZproxyActiveFilePath(),
		c.EZproxyAuditFileDirPath(),
		c.EZproxySearchRetries(),
		c.EZproxySearchDelay(),
		c.EZproxyTerminateSessions(),
		c.ConfigFile(),
	)
}

// Version emits version information and associated branding details whenever
// the user specifies the `--version` flag. The application exits after
// displaying this information.
func (c configTemplate) Version() string {
	return fmt.Sprintf("\n%s %s\n%s\n\n",
		MyAppName, Version, MyAppURL)
}

// Description emits branding information whenever the user specifies the `-h`
// flag. The application uses this as a header prior to displaying available
// CLI flag options.
func (c configTemplate) Description() string {
	return MyAppDescription
}

// GetNotificationTimeout calculates the timeout value for the entire message
// submission process, including the initial attempt and all retry attempts.
//
// This overall timeout value is computed using multiple values; (1) the base
// timeout value for a single message submission attempt, (2) the next
// scheduled notification (which was created using the configured delay we
// wish to force between message submission attempts), (3) the total number of
// retries allowed, (4) the delay between each retry attempt.
//
// This computed timeout value is intended to be used to cancel a notification
// attempt once it reaches this timeout threshold.
//
func GetNotificationTimeout(
	baseTimeout time.Duration,
	schedule time.Time,
	retries int,
	retriesDelay int,
) time.Duration {

	timeoutValue := (baseTimeout + time.Until(schedule)) +
		(time.Duration(retriesDelay) * time.Duration(retries))

	return timeoutValue
}

// MessageTrailer generates a branded "footer" for use with notifications.
func MessageTrailer(format BrandingFormat) string {

	switch format {
	case BrandingMarkdownFormat:
	case BrandingTextileFormat:
	default:
		errMsg := fmt.Sprintf("Invalid branding format %q used!", format)
		log.Warn(errMsg)
		return errMsg
	}

	return fmt.Sprintf(
		string(format),
		MyAppName,
		MyAppURL,
		Version,
		time.Now().Format(time.RFC3339Nano),
	)
}

// NewConfig is a factory function that produces a new Config object based on
// user provided values. While the fields are exported (due to requirements of
// third-party config packages), the intent is that the "getter" methods be
// used to provided a unified view of the current configuration generated from
// one or more configuration sources.
func NewConfig() (*Config, error) {

	myFuncName := caller.GetFuncName()

	config := Config{}

	// Bundle the returned `*.arg.Parser` for later use One potential use:
	// from `main()` so that we can explicitly display usage or help details
	// should the user-provided settings fail validation.
	log.Debugf("%s: Parsing flags", myFuncName)
	config.flagParser = arg.MustParse(&config.cliConfig)

	// If user specified a config file, try to use it, fail if not found
	log.Debugf(
		"%s: Checking whether config file has been specified",
		myFuncName,
	)
	if config.ConfigFile() != "" {

		log.Debugf(
			"%s: Config file %q specified",
			myFuncName,
			config.ConfigFile(),
		)

		// Used to help reduce the number of filepath.Clean() in locations
		// where it is considered "safe" to do so. Using this variable with
		// os.Open (in particular) upsets the gosec linter.
		sanitizedFilePath := filepath.Clean(config.ConfigFile())

		log.Debugf(
			"%s: Confirming sanitized version of %q file exists",
			myFuncName,
			sanitizedFilePath,
		)

		// path not found
		if _, err := os.Stat(filepath.Clean(config.ConfigFile())); os.IsNotExist(err) {
			return nil, fmt.Errorf(
				"%s: sanitized version of requested config file not found: %v",
				myFuncName,
				err,
			)
		}

		log.Debugf(
			"%s: Config file %q exists, attempting to open it",
			myFuncName,
			sanitizedFilePath,
		)
		// use direct function call here instead of our variable to comply
		// with gosec linting rules
		fh, err := os.Open(filepath.Clean(config.ConfigFile()))
		if err != nil {
			return nil, fmt.Errorf(
				"%s: unable to open config file: %v",
				myFuncName,
				err,
			)
		}
		defer func() {
			if err := fh.Close(); err != nil {
				// Ignore "file already closed" errors
				if !errors.Is(err, os.ErrClosed) {
					log.Errorf(
						"%s: failed to close file %q: %s",
						myFuncName,
						err.Error(),
					)
				}
			}
		}()
		log.Debugf("%s: Config file %q opened", myFuncName, sanitizedFilePath)

		log.Debugf(
			"%s: Attempting to load config file %q",
			myFuncName,
			sanitizedFilePath,
		)
		if err := config.LoadConfigFile(fh); err != nil {
			return nil, fmt.Errorf(
				"%s: error loading config file %q: %v",
				myFuncName,
				sanitizedFilePath,
				err,
			)
		}
		log.Debugf(
			"%s: Config file %q successfully loaded",
			myFuncName,
			sanitizedFilePath,
		)

		// explicitly close file, bail if failure occurs
		if err := fh.Close(); err != nil {
			return nil, fmt.Errorf(
				"%s: failed to close file %q: %w",
				myFuncName,
				sanitizedFilePath,
				err,
			)
		}
	}

	// Apply initial logging settings based on user-supplied settings
	config.configureLogging()

	// If no errors were encountered during parsing, proceed to validation of
	// configuration settings (both user-specified and defaults)
	if err := validate(config); err != nil {
		config.flagParser.WriteHelp(os.Stderr)
		return nil, err
	}

	return &config, nil

}

// LoadConfigFile reads from an io.Reader and unmarshals a configuration file
// in TOML format into the associated Config struct.
func (c *Config) LoadConfigFile(fileHandle io.Reader) error {

	configFileEntries, err := ioutil.ReadAll(fileHandle)
	if err != nil {
		return err
	}

	if err := toml.Unmarshal(configFileEntries, &c.fileConfig); err != nil {
		return err
	}

	return nil
}
