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

package events

import (
	"fmt"

	"github.com/apex/log"

	"github.com/atc0005/brick/internal/caller"
	"github.com/atc0005/go-ezproxy"
)

// This is a set of constants used with the Record.Action field to note the
// action taken in response to a received alert. The common use case is
// building a dynamic Title/Subject for various notifications.
const (
	ActionSuccessDisableRequestReceived string = "Disable user account request received"
	ActionSuccessDisabledUsername       string = "Username disabled"
	ActionSuccessDuplicatedUsername     string = "Username already disabled"
	ActionSuccessIgnoredUsername        string = "Username ignored due to ignore username entry"
	ActionSuccessIgnoredIPAddress       string = "Username ignored due to ignore IP entry"
	ActionSuccessTerminatedUserSession  string = "User sessions terminated"

	ActionSkippedTerminateUserSessions string = "User sessions termination not enabled; skipped"

	ActionFailureDisableRequestReceived   string = "Disable user account request log failure"
	ActionFailureDisabledUsername         string = "Username disable failure"
	ActionFailureDuplicatedUsername       string = "Username (duplicate) disable failure"
	ActionFailureIgnoredUsername          string = "Username ignore status check failure"
	ActionFailureIgnoredIPAddress         string = "IP Address ignore status check failure"
	ActionFailureUserSessionLookupFailure string = "Failed to lookup user sessions"
	ActionFailureTerminatedUserSession    string = "User session termination failure"
)

// Record is a collection of details that is saved to log files, sent by
// Microsoft Teams or email; this is a superset of types. This type contains
// the core details provided by the alert payload and select annotations
// associated with processing the alert payload.
type Record struct {

	// Alert is included since we will use the majority of the fields for
	// notifications and log entries
	Alert SplunkAlertEvent

	// Error optionally identifies the latest error with the associated event.
	// For Teams messages, this field is added as a "Fact" pair.
	Error error

	// Note is the additional message text used in notifications and log
	// entries. This field is not mutually exclusive; this field is displayed
	// alongside any error referenced by the `Error` field. For Teams
	// messages, this field is added as unformatted text via the associated
	// section's `Text` field.
	//
	// TODO: Rename this to `Summary?`
	Note string

	// Action briefly indicates what this application did in response to a
	// received alert. Values assigned to this field should only come from a
	// set of predefined constants in order to ensure consistency of log
	// messages, etc. This value is often used in single-event alert
	// notifications as part of the Subject or Title.
	Action string

	// SessionTerminationResults is a collection of results from attempts to
	// terminate sessions for the username specified in the alert payload.
	SessionTerminationResults []ezproxy.TerminateUserSessionResult
}

// NewRecord is a factory function that creates a Record from provided
// values. This function mostly exists as a way of having the compiler enforce
// that all required values for notifications are present.
func NewRecord(
	alert SplunkAlertEvent,
	err error,
	note string,
	action string,
	terminationResults []ezproxy.TerminateUserSessionResult,
) Record {

	record := Record{
		Alert:                     alert,
		Error:                     err,
		Note:                      note,
		Action:                    action,
		SessionTerminationResults: terminationResults,
	}

	if valid, err := record.Valid(); !valid {
		log.Errorf(
			"invalid Record '%#v' created at %s: %w",
			record,
			// Get the details for where NewRecord() was called, not the
			// details of where we are calling GetParentFuncFileLineInfo.
			caller.GetParentFuncFileLineInfo(),
			err,
		)
	}

	return record

}

// Valid is a helper method to perform light validation on Record fields in
// order to ensure that expected values are present. This is particularly
// important for the Action and Notes fields as they're heavily used in
// generated notifications.
func (rc Record) Valid() (bool, error) {

	// validate values

	// rely on compiler enforcing a valid SplunkAlertEvent is provided
	// alert

	// rely again on compiler
	// err

	// This is optional if an error value is already provided. The plan is to
	// have the msgCard.Text field use the Note field is if it is available,
	// otherwise call back to the Error field value.
	if rc.Error == nil {
		if rc.Note == "" {
			return false, fmt.Errorf(
				"empty or invalid Note field value provided: %s",
				rc.Note,
			)
		}
	}

	switch rc.Action {
	case "":
	case ActionSuccessDisableRequestReceived:
	case ActionSuccessDisabledUsername:
	case ActionSuccessDuplicatedUsername:
	case ActionSuccessIgnoredUsername:
	case ActionSuccessIgnoredIPAddress:
	case ActionSuccessTerminatedUserSession:
	case ActionSkippedTerminateUserSessions:
	case ActionFailureDisableRequestReceived:
	case ActionFailureDisabledUsername:
	case ActionFailureDuplicatedUsername:
	case ActionFailureIgnoredUsername:
	case ActionFailureIgnoredIPAddress:
	case ActionFailureUserSessionLookupFailure:
	case ActionFailureTerminatedUserSession:
	default:
		return false, fmt.Errorf(
			"empty or invalid Action field value provided: %s",
			rc.Action,
		)
	}

	return true, nil

}

// Records is a collection of Record values intended to allow easier bulk
// processing of event details.
type Records []Record

// Add is a helper method to collect event Records for later processing.
func (rcs *Records) Add(record Record) {

	// if pc, file, line, ok := runtime.Caller(1); ok {
	// 	log.Warnf(
	// 		"func %s called (from %q, line %d): ",
	// 		runtime.FuncForPC(pc).Name(),
	// 		file,
	// 		line,
	// 	)
	// }

	// TODO: Apply field validation to ensure we're collecting the absolute
	// minimum required details for later use. For example, we may need to
	// enforce that the Action field of a provided event Record only contains
	// known/valid Action constants for later use.

	// dereference pointer to Records to get access to slice
	// https://blog.golang.org/slices
	*rcs = append(*rcs, record)

}
