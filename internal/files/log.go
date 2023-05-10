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
	"fmt"
	"strings"

	"github.com/apex/log"

	"github.com/atc0005/brick/internal/caller"
	"github.com/atc0005/brick/internal/events"

	"github.com/atc0005/go-ezproxy"
)

// logEventDisableRequestReceived handles logging the event where a username
// has been reported by the remote monitoring system. This function emits the
// output to stdout for the init system to catch and also writes a templated
// message to the reported user events log for potential automation.
func logEventDisableRequestReceived(alert events.SplunkAlertEvent, reportedUserEventsLog *ReportedUserEventsLog) events.Record {

	requestReceivedMessage := fmt.Sprintf(
		"Disable request received from %q for username %q from IP %q",
		alert.PayloadSenderIP,
		alert.Username,
		alert.UserIP,
	)

	log.Debug(caller.GetFuncFileLineInfo())
	log.Infof(requestReceivedMessage)

	if err := appendToFile(
		fileEntry{
			Alert: alert,
		},
		reportedUserEventsLog.ReportTemplate,
		reportedUserEventsLog.FilePath,
		reportedUserEventsLog.FilePermissions,
	); err != nil {
		recordEventErr := fmt.Errorf(
			"func %s: error updating events log file %q: %w",
			caller.GetFuncName(),
			reportedUserEventsLog.FilePath,
			err,
		)

		return events.NewRecord(
			alert,
			recordEventErr,
			requestReceivedMessage,
			events.ActionFailureDisableRequestReceived,
			nil,
		)

	}

	return events.NewRecord(
		alert,
		nil,
		requestReceivedMessage,
		events.ActionSuccessDisableRequestReceived,
		nil,
	)

}

// logEventDisablingUsername handles logging the event where a username is
// being disabled. This function emits the output to stdout for the init
// system to catch. This function does NOT report the intent via
// notifications.
func logEventDisablingUsername(alert events.SplunkAlertEvent) {

	msgTemplate := "Disabling username %q from IP %q per report from %q"

	log.Debug(caller.GetFuncFileLineInfo())

	log.Infof(msgTemplate, alert.Username, alert.UserIP, alert.PayloadSenderIP)

}

// logEventDisabledUsername handles logging the event where a username
// has been successfully disabled. This function is responsible for emitting
// the success message to stdout for the init system to catch, write a
// templated message to the reported user events log for potential automation.
func logEventDisabledUsername(alert events.SplunkAlertEvent, reportedUserEventsLog *ReportedUserEventsLog) events.Record {

	disableSuccessMsg := fmt.Sprintf(
		"Disabled username %q from IP %q per report from %q",
		alert.Username,
		alert.UserIP,
		alert.PayloadSenderIP,
	)

	log.Debug(caller.GetFuncFileLineInfo())

	// emit to stdout right away in case we have problems recording this event
	// in the report users event log
	log.Info(disableSuccessMsg)

	if err := appendToFile(
		fileEntry{
			Alert: alert,
		},
		reportedUserEventsLog.DisableFirstEventTemplate,
		reportedUserEventsLog.FilePath,
		reportedUserEventsLog.FilePermissions,
	); err != nil {
		recordEventErr := fmt.Errorf(
			"func %s: error updating events log file %q: %w",
			caller.GetFuncName(),
			reportedUserEventsLog.FilePath,
			err,
		)

		return events.NewRecord(
			alert,
			recordEventErr,
			disableSuccessMsg,
			events.ActionFailureDisabledUsername,
			nil,
		)
	}

	return events.NewRecord(
		alert,
		nil,
		disableSuccessMsg,
		events.ActionSuccessDisabledUsername,
		nil,
	)

}

// logEventUsernameAlreadyDisabled handles logging the event where a username
// is already disabled, but another request has arrived to disable it, usually
// as a result of account compromise/sharing. This function emits the output
// to stdout for the init system to catch and also writes a templated message
// to the reported user events log for potential automation.
func logEventUsernameAlreadyDisabled(alert events.SplunkAlertEvent, reportedUserEventsLog *ReportedUserEventsLog) events.Record {

	// alreadyDisabledMsg := fmt.Sprintf(
	// 	"Received disable request from %q for user %q from IP %q;"+
	// 		" username is already disabled",
	// 	alert.PayloadSenderIP,
	// 	alert.Username,
	// 	alert.UserIP,
	// )

	alreadyDisabledMsg := fmt.Sprintf(
		"Username %q already disabled (current IP %q per report from %q)",
		alert.Username,
		alert.UserIP,
		alert.PayloadSenderIP,
	)

	log.Debug(caller.GetFuncFileLineInfo())
	log.Info(alreadyDisabledMsg)

	if err := appendToFile(
		fileEntry{
			Alert: alert,
		},
		reportedUserEventsLog.DisableRepeatEventTemplate,
		reportedUserEventsLog.FilePath,
		reportedUserEventsLog.FilePermissions,
	); err != nil {
		recordEventErr := fmt.Errorf(
			"func %s: error updating events log file %q: %w",
			caller.GetFuncName(),
			reportedUserEventsLog.FilePath,
			err,
		)

		return events.NewRecord(
			alert,
			recordEventErr,
			alreadyDisabledMsg,
			events.ActionFailureDuplicatedUsername,
			nil,
		)

	}

	return events.NewRecord(
		alert,
		nil,
		alreadyDisabledMsg,
		events.ActionSuccessDuplicatedUsername,
		nil,
	)

}

// logEventIgnoredIPAddress handles logging the event where an IP Address has
// been ignored due to inclusion of that IP Address in an "ignore file" for IP
// Addresses. This function emits the output to stdout for the init system to
// catch and also writes a templated message to the reported user events log
// for potential automation.
func logEventIgnoredIPAddress(alert events.SplunkAlertEvent, reportedUserEventsLog *ReportedUserEventsLog, ignoredEntriesFile string) events.Record {

	ignoreIPAddressMsg := fmt.Sprintf(
		"Ignored disable request from %q for user %q from IP %q due to presence in %q file.",
		alert.PayloadSenderIP,
		alert.Username,
		alert.UserIP,
		ignoredEntriesFile,
	)

	log.Debug(caller.GetFuncFileLineInfo())

	log.Info(ignoreIPAddressMsg)

	if err := appendToFile(
		fileEntry{
			Alert:              alert,
			IgnoredEntriesFile: ignoredEntriesFile,
		},
		reportedUserEventsLog.IgnoreTemplate,
		reportedUserEventsLog.FilePath,
		reportedUserEventsLog.FilePermissions,
	); err != nil {
		recordEventErr := fmt.Errorf(
			"func %s: error updating events log file %q: %w",
			caller.GetFuncName(),
			reportedUserEventsLog.FilePath,
			err,
		)

		return events.NewRecord(
			alert,
			recordEventErr,
			ignoreIPAddressMsg,
			events.ActionFailureIgnoredIPAddress,
			nil,
		)
	}

	return events.NewRecord(
		alert,
		nil,
		ignoreIPAddressMsg,
		events.ActionSuccessIgnoredIPAddress,
		nil,
	)

}

// logEventIgnoredUsername handles logging the event where a username has been
// ignored due to inclusion of that username in an "ignore file" for
// usernames. This function emits the output to stdout for the init system to
// catch and also writes a templated message to the reported user events log
// for potential automation.
func logEventIgnoredUsername(alert events.SplunkAlertEvent, reportedUserEventsLog *ReportedUserEventsLog, ignoredEntriesFile string) events.Record {

	ignoreUsernameMsg := fmt.Sprintf(
		"Ignored disable request from %q for user %q from IP %q due to presence in %q file.",
		alert.PayloadSenderIP,
		alert.Username,
		alert.UserIP,
		ignoredEntriesFile,
	)

	log.Debug(caller.GetFuncFileLineInfo())

	log.Info(ignoreUsernameMsg)

	if err := appendToFile(
		fileEntry{
			Alert:              alert,
			IgnoredEntriesFile: ignoredEntriesFile,
		},
		reportedUserEventsLog.IgnoreTemplate,
		reportedUserEventsLog.FilePath,
		reportedUserEventsLog.FilePermissions,
	); err != nil {
		recordEventErr := fmt.Errorf(
			"func %s: error updating events log file %q: %w",
			caller.GetFuncName(),
			reportedUserEventsLog.FilePath,
			err,
		)

		return events.NewRecord(
			alert,
			recordEventErr,
			ignoreUsernameMsg,
			events.ActionFailureIgnoredUsername,
			nil,
		)
	}

	return events.NewRecord(
		alert,
		nil,
		ignoreUsernameMsg,
		events.ActionSuccessIgnoredUsername,
		nil,
	)

}

// logEventTerminatingUserSession handles logging the event where a session
// for a username is being terminated. This function may be called multiple
// times, once per session associated with a username. This function emits the
// output to stdout for the init system to catch. This function does NOT
// record the event within the reported users event log nor does it generate a
// notification.
func logEventTerminatingUserSession(
	alert events.SplunkAlertEvent,
	userSession ezproxy.UserSession,
) {

	msgTemplate := "Terminating session %q (associated with IP %q) for username %q (from IP %q) per report from %q"

	log.Debug(caller.GetFuncFileLineInfo())

	log.Infof(
		msgTemplate,
		userSession.SessionID,
		userSession.IPAddress,
		alert.Username,
		alert.UserIP,
		alert.PayloadSenderIP,
	)

}

// logEventTerminatedUserSessions handles logging the event where sessions for
// a username have been terminated. This function is called once for a
// collection of termination results associated with a username. This function
// emits the output to stdout for the init system to catch and also sends a
// summary of the termination results as a notification.
func logEventTerminatedUserSessions(
	alert events.SplunkAlertEvent,
	reportedUserEventsLog *ReportedUserEventsLog,
	terminationResults ezproxy.TerminateUserSessionResults,
) events.Record {

	// Record origin *before* we start processing via loop
	log.Debug(caller.GetFuncFileLineInfo())

	// emit each termination result to stdout, write to the reported user
	// events log file
	successfulTerminationTmplPrefix := "Successfully terminated"
	failedTerminationTmplPrefix := "Failed to terminate"
	terminatedMsgTmpl := "%s session %q (associated with IP %q) for username %q (from IP %q) per report from %q [ExitCode: %d, StdOut: %q, StdErr: %q, Error: %q]"

	var failedTerminationsNum int
	var failedTerminationsSessionIDs []string
	for _, result := range terminationResults {

		// be optimistic!
		terminatedMsgPrefix := successfulTerminationTmplPrefix

		if result.Error != nil {
			terminatedMsgPrefix = failedTerminationTmplPrefix
			failedTerminationsSessionIDs = append(failedTerminationsSessionIDs, result.SessionID)
			failedTerminationsNum++
		}

		// guard against (nil) lack of error in results slice entry
		errStr := "None"
		if result.Error != nil {
			errStr = result.Error.Error()
		}

		terminatedMsg := fmt.Sprintf(
			terminatedMsgTmpl,
			terminatedMsgPrefix,
			result.SessionID,
			result.IPAddress,
			alert.Username,
			alert.UserIP,
			alert.PayloadSenderIP,
			result.ExitCode,
			result.StdOut,
			result.StdErr,
			errStr,
		)

		log.Info(terminatedMsg)

		// only record successful terminations in the reported user events log
		if result.Error == nil {

			var recordEventErr error
			if err := appendToFile(
				fileEntry{
					Alert:       alert,
					UserSession: result.UserSession,
				},
				reportedUserEventsLog.TerminateUserSessionEventTemplate,
				reportedUserEventsLog.FilePath,
				reportedUserEventsLog.FilePermissions,
			); err != nil {

				recordEventErr = fmt.Errorf(
					"func %s: error updating events log file %q: %w",
					caller.GetFuncName(),
					reportedUserEventsLog.FilePath,
					recordEventErr,
				)

				return events.NewRecord(
					alert,
					recordEventErr,
					terminatedMsg,
					events.ActionFailureTerminatedUserSession,
					terminationResults,
				)

			}

		}

	}

	successfulTerminations := len(terminationResults) - failedTerminationsNum

	// emit via stdout (for systemd/syslog)
	log.Infof(
		"Session termination summary for %q: [success: %d, failure: %d]",
		alert.Username,
		successfulTerminations,
		failedTerminationsNum,
	)

	if terminationResults.HasError() {

		terminationResultsFailureMsg := fmt.Sprintf(
			"%d errors occurred while terminating %d sessions for username %s",
			failedTerminationsNum,
			len(terminationResults),
			alert.Username,
		)

		terminationResultsError := fmt.Errorf(
			"failed to terminate sessions: %s",
			strings.Join(failedTerminationsSessionIDs, ", "),
		)

		return events.NewRecord(
			alert,
			terminationResultsError,
			terminationResultsFailureMsg,
			events.ActionFailureTerminatedUserSession,
			terminationResults,
		)
	}

	// TODO: We need to determine our % of success and convey that at a
	// glance. The current "User sessions terminated" suffix conveys at
	// *something* was terminated, but this may not be true (e.g., the recent
	// test case was a complete failure and still that suffix is used)

	sessionTerminationResultsSuccessMsg := fmt.Sprintf(
		"Successfully terminated all %d user sessions for %q",
		len(terminationResults),
		alert.Username,
	)

	return events.NewRecord(
		alert,
		nil,
		sessionTerminationResultsSuccessMsg,
		events.ActionSuccessTerminatedUserSession,
		terminationResults,
	)

}
