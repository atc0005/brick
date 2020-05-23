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

package files

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"text/template"

	"github.com/apex/log"

	"github.com/atc0005/brick/events"
)

func ProcessEvent(
	alert events.SplunkAlertEvent,
	disabledUsers *DisabledUsers,
	reportedUserEventsLog *ReportedUserEventsLog,
	ignoredSources IgnoredSources,
	notifyWorkQueue chan<- events.Record,
) error {

	if err := logEventReportingUsername(alert, reportedUserEventsLog); err != nil {
		return err
	}

	// check to see if username has already been disabled
	disabledUserEntry := alert.Username + disabledUsers.EntrySuffix
	disableEntryFound, disableEntryLookupErr := inList(disabledUserEntry, disabledUsers.FilePath)
	if disableEntryLookupErr != nil {
		errMsg := fmt.Errorf(
			"error while checking disabled status for user %q from IP %q: %w",
			alert.Username,
			alert.UserIP,
			disableEntryLookupErr,
		)

		if !ignoredSources.IgnoreLookupErrors {

			// only send notifications if lookup errors are not ignored
			go func() {
				notifyWorkQueue <- events.Record{
					Alert: alert,
					Error: errMsg,
				}
			}()

			return errMsg
		}

		// console, local logs, syslog via systemd, etc.
		log.Warn(errMsg.Error())
	}

	// if username has not been disabled yet, proceed with additional checks
	// before attempting to disable the account
	if !disableEntryFound {

		// check to see if username has been ignored
		userIgnoreEntryFound, userIgnoreLookupErr := inList(alert.Username, ignoredSources.IgnoredUsersFile)
		if userIgnoreLookupErr != nil {

			errMsg := fmt.Errorf(
				"error while checking ignored status for user %q from IP %q: %w",
				alert.Username,
				alert.UserIP,
				userIgnoreLookupErr,
			)

			if !ignoredSources.IgnoreLookupErrors {

				// only send notifications if lookup errors are not ignored
				go func() {
					notifyWorkQueue <- events.Record{
						Alert: alert,
						Error: errMsg,
					}
				}()

				return errMsg
			}

			// console, local logs, syslog via systemd, etc.
			log.Warn(errMsg.Error())
		}

		if userIgnoreEntryFound {

			logEventIgnoringUsernameErr := logEventIgnoringUsername(
				alert,
				reportedUserEventsLog,
				ignoredSources.IgnoredUsersFile,
			)

			ignoreUserMsg := fmt.Sprintf(
				"Ignoring disable request from %q for user %q from IP %q due to presence in %q file.",
				alert.PayloadSenderIP,
				alert.Username,
				alert.UserIP,
				ignoredSources.IgnoredUsersFile,
			)
			go func() {
				notifyWorkQueue <- events.Record{
					Alert:  alert,
					Error:  logEventIgnoringUsernameErr,
					Note:   ignoreUserMsg,
					Action: events.ActionSuccessIgnoredUsername,
				}
			}()

			return logEventIgnoringUsernameErr
		}

		// check to see if IP Address has been ignored
		ipAddressIgnoreEntryFound, ipAddressIgnoreLookupErr := inList(
			alert.UserIP, ignoredSources.IgnoredIPAddressesFile)
		if ipAddressIgnoreLookupErr != nil {

			errMsg := fmt.Errorf(
				"error while checking ignored status for IP %q associated with user %q: %w",
				alert.UserIP,
				alert.Username,
				ipAddressIgnoreLookupErr,
			)

			if !ignoredSources.IgnoreLookupErrors {

				// only send notifications if lookup errors are not ignored
				go func() {
					notifyWorkQueue <- events.Record{
						Alert: alert,
						Error: errMsg,
					}
				}()

				return errMsg
			}

			// console, local logs, syslog via systemd, etc.
			log.Warn(errMsg.Error())
		}

		if ipAddressIgnoreEntryFound {

			logEventIgnoringUsernameErr := logEventIgnoringUsername(
				alert,
				reportedUserEventsLog,
				ignoredSources.IgnoredIPAddressesFile,
			)

			ignoreIPAddressMsg := fmt.Sprintf(
				"Ignoring disable request from %q for user %q from IP %q due to presence in %q file.",
				alert.PayloadSenderIP,
				alert.Username,
				alert.UserIP,
				ignoredSources.IgnoredIPAddressesFile,
			)

			go func() {
				notifyWorkQueue <- events.Record{
					Alert:  alert,
					Error:  logEventIgnoringUsernameErr,
					Note:   ignoreIPAddressMsg,
					Action: events.ActionSuccessIgnoredIPAddress,
				}
			}()

			return logEventIgnoringUsernameErr

		}

		// disable user account
		if err := disableUser(alert, disabledUsers); err != nil {

			go func() {
				notifyWorkQueue <- events.Record{
					Alert: alert,
					Error: err,
				}
			}()

			return err
		}

		if err := logEventDisablingUsername(alert, reportedUserEventsLog); err != nil {

			go func() {
				notifyWorkQueue <- events.Record{
					Alert: alert,
					Error: err,
				}
			}()

			return err
		}

		disableSuccessMsg := fmt.Sprintf(
			"Disabled user %q from IP %q",
			alert.Username,
			alert.UserIP,
		)

		go func() {
			notifyWorkQueue <- events.Record{
				Alert:  alert,
				Note:   disableSuccessMsg,
				Action: events.ActionSuccessDisabledUsername,
			}
		}()

		log.Infof(disableSuccessMsg)

		// required to indicate that we successfully disabled user account
		return nil

	}

	// if username is already disabled, skip adding to disable file, but still
	// log the request so that all associated IPs can be banned by fail2ban
	alreadyDisabledMsg := fmt.Sprintf(
		"Received disable request from %q for user %q from IP %q;"+
			" username is already disabled",
		alert.PayloadSenderIP,
		alert.Username,
		alert.UserIP,
	)

	go func() {
		notifyWorkQueue <- events.Record{
			Alert: alert,
			// this isn't technically an "error", more of a "something worth
			// noting" event
			Note:   alreadyDisabledMsg,
			Action: events.ActionSuccessDuplicateUsername,
		}
	}()

	log.Debug(alreadyDisabledMsg)

	if err := logEventUserAlreadyDisabled(alert, reportedUserEventsLog); err != nil {

		go func() {
			notifyWorkQueue <- events.Record{
				Alert: alert,
				Error: err,
			}
		}()

		return err
	}

	// FIXME: Anything else needed at this point?
	return nil

}

// inList accepts a string and a fully-qualified path to a file containing a
// list of such strings (commonly usernames or single IP Addresses), one per
// line. Lines beginning with a `#` character are ignored. Leading and
// trailing whitespace per line is ignored.
func inList(needle string, haystack string) (bool, error) {

	log.Debugf("Attempting to open %q", haystack)

	// TODO: How do we handle the situation where the file does not exist
	// ahead of time? Since this application will manage the file, it should
	// be able to create it with the desired permissions?
	f, err := os.Open(haystack)
	if err != nil {
		return false, fmt.Errorf("error encountered opening file %q: %w", haystack, err)
	}
	defer f.Close()

	log.Debugf("Searching for: %q", needle)

	s := bufio.NewScanner(f)
	var lineno int

	// TODO: Does Scan() perform any whitespace manipulation already?
	for s.Scan() {
		lineno++
		currentLine := s.Text()
		log.Debugf("Scanned line %d from %q: %q\n", lineno, haystack, currentLine)

		currentLine = strings.TrimSpace(currentLine)
		log.Debugf("Line %d from %q after lowercasing and whitespace removal: %q\n",
			lineno, haystack, currentLine)

		// explicitly ignore comments
		if strings.HasPrefix(currentLine, "#") {
			log.Debugf("Ignoring comment line %d", lineno)
			continue
		}

		log.Debugf("Checking whether line %d is a match: %q %q", lineno, currentLine, needle)
		if strings.EqualFold(currentLine, needle) {
			log.Debugf("Match found on line %d, returning true to indicate this", lineno)
			return true, nil
		}

	}

	log.Debug("Exited s.Scan() loop")

	// report any errors encountered while scanning the input file
	if err := s.Err(); err != nil {
		return false, err
	}

	// otherwise, report that the requested needle was not found
	return false, nil

}

// disableUser adds the specified username to the disabled users file. This
// function is intended to be called from within another function that first
// confirms that the specified user account has not already been disabled.
func disableUser(alert events.SplunkAlertEvent, disabledUsers *DisabledUsers) error {

	// NOTE: Notifications are handled by the caller

	log.Debug("DisableUser: disabling user per alert")
	if err := appendToFile(
		fileEntry{
			SplunkAlertEvent: alert,
			EntrySuffix:      disabledUsers.EntrySuffix,
		},
		disabledUsers.Template,
		disabledUsers.FilePath,
		disabledUsers.FilePermissions,
	); err != nil {
		return fmt.Errorf(
			"error updating disabled user file %q: %w",
			disabledUsers.FilePath,
			err,
		)
	}

	return nil

}

// appendToFile is a helper function that accepts a new message, a destination
// filename and intended permissions for the filename if it does not already
// exist. All leading and trailing whitespace is removed from the new message
// and one trailing newline appended.
func appendToFile(entry fileEntry, tmpl *template.Template, filename string, perms os.FileMode) error {

	var mutex = &sync.Mutex{}

	log.Debugf("Attempting to open %q", filename)

	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, perms)
	if err != nil {
		return fmt.Errorf("error encountered opening file %q: %w", filename, err)
	}
	defer f.Close()
	log.Debugf("Successfully opened %q", filename)

	log.Debug("Locking mutex")
	mutex.Lock()
	defer func() {
		log.Debug("Unlocking mutex from deferred anonymous func")
		mutex.Unlock()
	}()

	log.Debugf("Executing template to update %q", filename)
	if err := tmpl.Execute(f, entry); err != nil {
		f.Close() // ignore error; Write error takes precedence
		return fmt.Errorf("error writing to file %q: %w", filename, err)
	}
	log.Debugf("Successfully executed template to update %q", filename)

	log.Debug("Syncing file modifications")
	if err := f.Sync(); err != nil {
		return fmt.Errorf(
			"failed to explicitly sync file %q after writing: %s",
			filename,
			err,
		)
	}
	log.Debugf("Successfully synced modifications to %q", filename)

	log.Debugf("Closing %q", filename)
	if err := f.Close(); err != nil {
		return fmt.Errorf("error closing file %q: %w", filename, err)
	}
	log.Debugf("Successfully closed %q", filename)

	return nil
}
