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

const (
	ActionSuccessDisabledUsername  string = "Username disabled"
	ActionSuccessDuplicateUsername string = "Username already disabled"
	ActionSuccessIgnoredUsername   string = "Username ignored due to ignore username entry"
	ActionSuccessIgnoredIPAddress  string = "Username ignored due to ignore IP entry"

	// ActionFailureDisabledUsername  string = ""
	// ActionFailureDuplicateUsername string = ""
	// ActionFailureIgnoredUsername   string = ""
	// ActionFailureIgnoredIPAddress  string = ""
)

// AlertResponse is meant to help indicate what specific action we took based
// on the received alert.
// type AlertResponse string

// Record is a collection of details that is saved to log files, sent by
// Microsoft Teams or email; this is a superset of types. This type contains
// the core details received by the alert payload and select annotations
// associated with processing the alert payload.
type Record struct {

	// Alert is included since we will use the majority of the fields for
	// notifications and log entries
	Alert SplunkAlertEvent

	// Error optionally identifies the latest error with the associated event
	Error error

	// Note is the additional message text used in notifications and log
	// entries
	// FIXME: Not a fan of the field name
	Note string

	// Action notes what this application did in response to a received alert
	// Action AlertResponse
	Action string
}
