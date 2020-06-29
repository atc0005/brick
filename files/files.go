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

// Package files provides types and functions related to the various files
// created and/or used by this application
package files

import (
	"os"
	"strings"
	"text/template"

	"github.com/atc0005/brick/events"
)

// fileEntry represents the values that are used when generating entries via
// templates for flat-files: disable user accounts, log file of actions taken
type fileEntry struct {
	events.SplunkAlertEvent
	EntrySuffix        string
	IgnoredEntriesFile string
}

// FlatFile represents a text file that this application is responsible for
// populating. This includes the disable users file and the events log file
// parsed by fail2ban.
type FlatFile struct {
	// FileOwner represents the OS user account that owns this file
	FileOwner string

	// FileGroup represents the OS user group with defined permissions for this
	// file
	FileGroup string

	// FilePermissions represents the classic POSIX read, write, execute bits
	// granting (or denying) access to a file/directory. Because this file
	// *IS* read by EZproxy, the permissions on this file should permit
	// *read* access by that daemon's user/group.
	FilePermissions os.FileMode

	// Path is the fully-qualified path to the disables users file created and
	// managed by this application.
	FilePath string
}

// DisabledUsers represents the text file that EZproxy monitors for user
// accounts that should not be allowed to login. This application is
// responsible for recording user accounts in this file that it receives via
// alert payloads and are not otherwise excluded due to "ignored" user
// accounts or IP Addresses lists.
//
// TODO: Consider singular vs plural naming
// `DisabledUsers` vs `DisabledUser`
type DisabledUsers struct {
	FlatFile

	// Example future field (to help illustrate extensibility):
	// LDAPGroup

	// EntrySuffix is the string that is appended after every username added
	// to the disabled users file in order to deny login access.
	EntrySuffix string

	// Template is a parsed template representing the line written to this
	// file when a user account is disabled.
	Template *template.Template
}

// ReportedUserEventsLog represents a log file where this application
// records that a user account was reported and what action was taken. Actions
// include ignoring user accounts because they're in a external "safe" or
// "ignore" list (to prevent unintentional access disruption) and disabling
// user accounts (writing entries to `DisabledUsersFile`). This log file is
// intended to be human-readable, but also parsable by external tooling so
// that automatic actions can be performed (e.g, temporary banning of
// associated IP Addresses).
//
// TODO: Consider singular vs plural naming
// `ReportedUserEventsLog` vs `ReportedUserEvents` vs `ReportedUserEvent`
type ReportedUserEventsLog struct {
	FlatFile

	// Example future field (to help illustrate extensibility):
	// TeamsWebhookURL

	// ReportTemplate is a parsed template representing the log line written
	// when a user account is reported via alert payload. This entry occurs
	// regardless of whether an account is eventually ignored or disabled.
	ReportTemplate *template.Template

	// DisableTemplate is a parsed template representing the log line written
	// when a user account is reported the first time via alert payload and
	// the user account is disabled.
	DisableFirstEventTemplate *template.Template

	// DisableRepeatEventTemplate is a parsed template representing the log line written
	// when a user account is reported via alert payload again after the user
	// account is already disabled.
	DisableRepeatEventTemplate *template.Template

	// IgnoreTemplate is a parsed template representing the log line written
	// when a user account is reported via alert payload and the user account
	// or associated IP Address is ignored due to its presence in either the
	// specified "safe" or "ignored" user accounts file or IP Addresses file.
	IgnoreTemplate *template.Template
}

// IgnoredSources represents the various sources of "safe" or "ignore" entries
// for this application. This includes user account names and client IP
// Addresses.
type IgnoredSources struct {
	IgnoredUsersFile       string
	IgnoredIPAddressesFile string
	IgnoreLookupErrors     bool
}

// NewReportedUserEventsLog constructs a ReportedUserEventsLog type with
// parsed templates already set.
func NewReportedUserEventsLog(path string, permissions os.FileMode) *ReportedUserEventsLog {

	// parse templates
	reportedUserEventTemplate := template.Must(template.New(
		"reportedUserEventTemplate").Parse(reportedUserEventTemplateText))

	disabledUserFirstEventTemplate := template.Must(template.New(
		"disabledUserFirstEventTemplate").Parse(disabledUserFirstEventTemplateText))

	disabledUserRepeatEventTemplate := template.Must(template.New(
		"disabledUserRepeatEventTemplate").Parse(disabledUserRepeatEventTemplateText))

	ignoredUserEventTemplate := template.Must(template.New(
		"ignoredUserEventTemplate").Parse(ignoredUserEventTemplateText))

	ruel := ReportedUserEventsLog{
		FlatFile: FlatFile{
			FilePath:        path,
			FilePermissions: permissions,
		},
		ReportTemplate:             reportedUserEventTemplate,
		DisableFirstEventTemplate:  disabledUserFirstEventTemplate,
		DisableRepeatEventTemplate: disabledUserRepeatEventTemplate,
		IgnoreTemplate:             ignoredUserEventTemplate,
	}

	return &ruel

}

// NewDisabledUsers constructs a DisabledUsers type with parsed template
// already set.
func NewDisabledUsers(path string, entrySuffix string, permissions os.FileMode) *DisabledUsers {

	// parse template for disabled users file, provide a ToLower template
	// function that can be used to case-fold values written to the disabled
	// users file
	disabledUsersFileTemplate := template.Must(
		template.New(
			"disabledUsersFileTemplate",
		).Funcs(
			template.FuncMap{
				"ToLower": strings.ToLower,
			},
		).Parse(
			disabledUsersFileTemplateText,
		),
	)

	du := DisabledUsers{
		FlatFile: FlatFile{
			FilePath:        path,
			FilePermissions: permissions,
		},
		Template:    disabledUsersFileTemplate,
		EntrySuffix: entrySuffix,
	}

	return &du

}

// NewIgnoredSources constructs an IgnoredSources type
func NewIgnoredSources(
	ignoredUsersFile string,
	ignoredIPAddressesFile string,
	ignoreLookupErrors bool,
) IgnoredSources {

	ignoredSources := IgnoredSources{
		IgnoredUsersFile:       ignoredUsersFile,
		IgnoredIPAddressesFile: ignoredIPAddressesFile,
		IgnoreLookupErrors:     ignoreLookupErrors,
	}

	return ignoredSources
}
