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

	"github.com/apex/log"

	"github.com/atc0005/brick/events"
)

func logEventReportingUsername(alert events.SplunkAlertEvent, reportedUserEventsLog *ReportedUserEventsLog) error {

	msgTemplate := "Disable request received from %q for username %q from IP %q"

	log.Debugf(
		"LogEventReportedUser: "+msgTemplate,
		alert.Username,
		alert.PayloadSenderIP,
		alert.UserIP,
	)

	if err := appendToFile(
		fileEntry{
			SplunkAlertEvent: alert,
		},
		reportedUserEventsLog.ReportTemplate,
		reportedUserEventsLog.FilePath,
		reportedUserEventsLog.FilePermissions,
	); err != nil {
		return fmt.Errorf(
			"func LogEventReportedUser: error updating events log file %q: %w",
			reportedUserEventsLog.FilePath,
			err,
		)
	}

	log.Infof(msgTemplate, alert.PayloadSenderIP, alert.Username, alert.UserIP)

	return nil

}

func logEventDisablingUsername(alert events.SplunkAlertEvent, reportedUserEventsLog *ReportedUserEventsLog) error {

	msgTemplate := "Disabling username %q from IP %q per report from %q"

	log.Debugf(
		"LogEventDisablingUsername: "+msgTemplate,
		alert.Username,
		alert.UserIP,
		alert.PayloadSenderIP,
	)

	if err := appendToFile(
		fileEntry{
			SplunkAlertEvent: alert,
		},
		reportedUserEventsLog.DisableFirstEventTemplate,
		reportedUserEventsLog.FilePath,
		reportedUserEventsLog.FilePermissions,
	); err != nil {
		return fmt.Errorf(
			"func LogEventDisablingUsername: error updating events log file %q: %w",
			reportedUserEventsLog.FilePath,
			err,
		)
	}

	log.Infof(msgTemplate, alert.Username, alert.UserIP, alert.PayloadSenderIP)

	return nil

}

func logEventUserAlreadyDisabled(alert events.SplunkAlertEvent, reportedUserEventsLog *ReportedUserEventsLog) error {

	msgTemplate := "Username %q already disabled (current IP %q per report from %q)"

	log.Debugf(
		"logEventUserAlreadyDisabled: "+msgTemplate,
		alert.Username,
		alert.UserIP,
		alert.PayloadSenderIP,
	)

	if err := appendToFile(
		fileEntry{
			SplunkAlertEvent: alert,
		},
		reportedUserEventsLog.DisableRepeatEventTemplate,
		reportedUserEventsLog.FilePath,
		reportedUserEventsLog.FilePermissions,
	); err != nil {
		return fmt.Errorf(
			"func LogEventDisabledUser: error updating events log file %q: %w",
			reportedUserEventsLog.FilePath,
			err,
		)
	}

	log.Infof(msgTemplate, alert.Username, alert.UserIP, alert.PayloadSenderIP)

	return nil

}

func logEventIgnoringUsername(alert events.SplunkAlertEvent, reportedUserEventsLog *ReportedUserEventsLog, ignoredEntriesFile string) error {

	msgTemplate := "Ignoring disable request from %q for user %q from IP %q due to presence in %q file."

	log.Debugf(
		"LogEventIgnoredUserOrIP: "+msgTemplate,
		alert.PayloadSenderIP,
		alert.Username,
		alert.UserIP,
		ignoredEntriesFile,
	)

	if err := appendToFile(
		fileEntry{
			SplunkAlertEvent:   alert,
			IgnoredEntriesFile: ignoredEntriesFile,
		},
		reportedUserEventsLog.IgnoreTemplate,
		reportedUserEventsLog.FilePath,
		reportedUserEventsLog.FilePermissions,
	); err != nil {
		return fmt.Errorf(
			"func LogEventIgnoredUserOrIP: error updating events log file %q: %w",
			reportedUserEventsLog.FilePath,
			err,
		)
	}

	log.Infof(
		msgTemplate,
		alert.PayloadSenderIP,
		alert.Username,
		alert.UserIP,
		ignoredEntriesFile,
	)

	return nil

}
