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

package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/apex/log"
	"github.com/atc0005/brick/config"
	"github.com/atc0005/brick/events"
	"github.com/atc0005/brick/internal/caller"
	"github.com/atc0005/go-ezproxy"

	// use our fork for now until recent work can be submitted for inclusion
	// in the upstream project
	goteamsnotify "github.com/atc0005/go-teams-notify"

	send2teams "github.com/atc0005/send2teams/teams"
)

// addFactPair accepts a MessageCard, MessageCardSection, a key and one or
// many values related to the provided key. An attempt is made to format the
// key and all values as code in an effort to keep Teams from parsing some
// special characters as Markdown code formatting characters. If an error
// occurs, this error is used as the MessageCard Text field to help make this
// failure more prominent, otherwise the top-level MessageCard fields remain
// untouched.
//
// FIXME: Rework and offer upstream?
func addFactPair(msg *goteamsnotify.MessageCard, section *goteamsnotify.MessageCardSection, key string, values ...string) {

	for idx := range values {
		values[idx] = send2teams.TryToFormatAsCodeSnippet(values[idx])
	}

	if err := section.AddFactFromKeyValue(
		key,
		values...,
	); err != nil {
		from := caller.GetFuncFileLineInfo()
		errMsg := fmt.Sprintf("%s error returned from attempt to add fact from key/value pair: %v", from, err)
		log.Errorf("%s %s", from, errMsg)
		msg.Text = msg.Text + "\n\n" + send2teams.TryToFormatAsCodeSnippet(errMsg)
	}
}

// FIXME: Move this elsewhere or remove; it does not look like a viable option
// func generateTable(sTermResults []ezproxy.TerminateUserSessionResult) string {
//
// 	// TODO: Can we generate a Markdown table here?
// 	sessionResultsTableHeader := `
// 		| ID | IP | ExitCode | StdOut | StdErr | ErrorMsg | <br>
// 		| --- | --- | --- | --- | --- | --- | <br>
// 		`
//
// 	// addFactPair(&msgCard, sessionTerminationResultsSection, "Error", record.Error.Error())
//
// 	var sessionResultsTableBody string
// 	for _, result := range sTermResults {
//
// 		// guard against (nil) lack of error in results slice entry
// 		errStr := "None"
// 		if result.Error != nil {
// 			errStr = result.Error.Error()
// 		}
//
// 		sessionResultsTableBody += fmt.Sprintf(
// 			"| %s | %s | %d | %s | %s | %s | <br>",
// 			result.SessionID,
// 			result.IPAddress,
// 			result.ExitCode,
// 			result.StdOut,
// 			result.StdErr,
// 			errStr,
// 		)
// 	}
//
// 	// TODO: No clue if this will work, figured it was worth a shot.
// 	sessionResultsTable := sessionResultsTableHeader + sessionResultsTableBody
//
// 	return sessionResultsTable
// }

// getTerminationResultsList generates a Markdown list summarizing session
// termination results. See also atc0005/go-ezproxy#18.
func getTerminationResultsList(sTermResults []ezproxy.TerminateUserSessionResult) string {

	var sessionResultsStringSets string
	for _, result := range sTermResults {

		// guard against (nil) lack of error in results slice entry
		errStr := "None"
		if result.Error != nil {
			errStr = result.Error.Error()
		}

		sessionResultsStringSets += fmt.Sprintf(
			"- { SessionID: %q, IPAddress: %q, ExitCode: %q, StdOut: %q, StdErr: %q, Error: %q }\n\n",
			result.SessionID,
			result.IPAddress,
			strconv.Itoa(result.ExitCode),
			result.StdOut,
			result.StdErr,
			errStr,
		)
	}

	return sessionResultsStringSets
}

// getMsgCardMainSectionText evaluates the provided event Record and builds a
// primary message suitable for display as the main notification Text field.
// This message is generated first from the Note field if available, or from
// the Error field. This precedence allows for using a provided Note as a brief
// summary while still using the Error field in a dedicated section.
func getMsgCardMainSectionText(record events.Record) string {

	// This part of the message card is valuable "real estate" for eyeballs;
	// we should ensure we are communicating what just occurred instead
	// of using a mostly static block of text.

	var msgCardTextField string

	switch {
	case record.Note != "":
		msgCardTextField = "Summary: " + record.Note
	case record.Error != nil:
		msgCardTextField = "Error: " + record.Error.Error()

	// Attempting to use an empty string for the top-level message card Text
	// field results in a notification failure, so set *something* to meet
	// those requirements.
	default:
		msgCardTextField = "FIXME: Missing Note for this event record!"
	}

	return msgCardTextField

}

// getMsgCardTitle is a helper function used to generate the title for message
// cards. This function uses the provided prefix and event Record to generate
// stable titles reflecting the step in the disable user process at which the
// notification was generated; the intent is to quickly tell where the process
// halted for troubleshooting purposes.
func getMsgCardTitle(msgCardTitlePrefix string, record events.Record) string {

	var msgCardTitle string

	switch record.Action {

	// case record.Error != nil:
	// 	msgCardTitle = "[ERROR] " + record.Error.Error()

	// TODO: Calculate step labeling based off of enabled features (see GH-65).

	case events.ActionSuccessDisableRequestReceived, events.ActionFailureDisableRequestReceived:
		msgCardTitle = msgCardTitlePrefix + "[step 1 of 3] " + record.Action

	case events.ActionSuccessDisabledUsername, events.ActionFailureDisabledUsername:
		msgCardTitle = msgCardTitlePrefix + "[step 2 of 3] " + record.Action

	case events.ActionSuccessDuplicatedUsername, events.ActionFailureDuplicatedUsername:
		msgCardTitle = msgCardTitlePrefix + "[step 2 of 3] " + record.Action

	case events.ActionSuccessIgnoredUsername, events.ActionFailureIgnoredUsername:
		msgCardTitle = msgCardTitlePrefix + "[step 2 of 3] " + record.Action

	case events.ActionSuccessIgnoredIPAddress, events.ActionFailureIgnoredIPAddress:
		msgCardTitle = msgCardTitlePrefix + "[step 2 of 3] " + record.Action

	case events.ActionSuccessTerminatedUserSession,
		events.ActionFailureUserSessionLookupFailure,
		events.ActionFailureTerminatedUserSession,
		events.ActionSkippedTerminateUserSessions:
		msgCardTitle = msgCardTitlePrefix + "[step 3 of 3] " + record.Action

	default:
		msgCardTitle = msgCardTitlePrefix + " [UNKNOWN] " + record.Action
		log.Warnf("UNKNOWN record: %v+\n", record)
	}

	return msgCardTitle
}

// createMessage receives an event Record and generates a MessageCard which is
// used to generate a Microsoft Teams message.
func createMessage(record events.Record) goteamsnotify.MessageCard {

	log.Debugf("createMessage: alert received: %#v", record)

	// build MessageCard for submission
	msgCard := goteamsnotify.NewMessageCard()

	msgCardTitlePrefix := config.MyAppName + ": "

	msgCard.Title = getMsgCardTitle(msgCardTitlePrefix, record)

	// msgCard.Text = record.Note
	msgCard.Text = getMsgCardMainSectionText(record)

	/*
		Errors Section
	*/

	disableUserRequestErrors := goteamsnotify.NewMessageCardSection()
	disableUserRequestErrors.Title = "## Disable User Request Errors"
	disableUserRequestErrors.StartGroup = true

	switch {
	case record.Error != nil:
		addFactPair(&msgCard, disableUserRequestErrors, "Error", record.Error.Error())
	case record.Error == nil:
		disableUserRequestErrors.Text = "None"
	}

	if err := msgCard.AddSection(disableUserRequestErrors); err != nil {
		errMsg := fmt.Sprintf("Error returned from attempt to add disableUserRequestErrors: %v", err)
		log.Error("createMessage: " + errMsg)
		msgCard.Text = msgCard.Text + "\n\n" + send2teams.TryToFormatAsCodeSnippet(errMsg)
	}

	// If Session Termination is enabled, create Termination Results section
	if record.SessionTerminationResults != nil {

		sessionTerminationResultsSection := goteamsnotify.NewMessageCardSection()
		sessionTerminationResultsSection.Title = "## Session Termination Results"
		sessionTerminationResultsSection.StartGroup = true

		sessionTerminationResultsSection.Text = getTerminationResultsList(record.SessionTerminationResults)

		if err := msgCard.AddSection(sessionTerminationResultsSection); err != nil {
			errMsg := fmt.Sprintf("Error returned from attempt to add sessionTerminationResultsSection: %v", err)
			log.Error("createMessage: " + errMsg)
			msgCard.Text = msgCard.Text + "\n\n" + send2teams.TryToFormatAsCodeSnippet(errMsg)
		}

	}

	/*
		Disable User Request Details Section - Core of alert details
	*/

	disableUserRequestDetailsSection := goteamsnotify.NewMessageCardSection()
	disableUserRequestDetailsSection.Title = "## Disable User Request Details"
	disableUserRequestDetailsSection.StartGroup = true

	addFactPair(&msgCard, disableUserRequestDetailsSection, "Username", record.Alert.Username)
	addFactPair(&msgCard, disableUserRequestDetailsSection, "User IP", record.Alert.UserIP)
	addFactPair(&msgCard, disableUserRequestDetailsSection, "Alert/Search Name", record.Alert.AlertName)
	addFactPair(&msgCard, disableUserRequestDetailsSection, "Alert/Search ID", record.Alert.SearchID)

	if err := msgCard.AddSection(disableUserRequestDetailsSection); err != nil {
		errMsg := fmt.Sprintf("Error returned from attempt to add disableUserRequestDetailsSection: %v", err)
		log.Error("createMessage: " + errMsg)
		msgCard.Text = msgCard.Text + "\n\n" + send2teams.TryToFormatAsCodeSnippet(errMsg)
	}

	/*
		Alert Request Summary Section - General client request details
	*/

	alertRequestSummarySection := goteamsnotify.NewMessageCardSection()
	alertRequestSummarySection.Title = "## Alert Request Summary"
	alertRequestSummarySection.StartGroup = true

	addFactPair(&msgCard, alertRequestSummarySection, "Received at", record.Alert.LocalTime)
	addFactPair(&msgCard, alertRequestSummarySection, "Endpoint path", record.Alert.EndpointPath)
	addFactPair(&msgCard, alertRequestSummarySection, "HTTP Method", record.Alert.HTTPMethod)
	addFactPair(&msgCard, alertRequestSummarySection, "Alert Sender IP", record.Alert.PayloadSenderIP)

	if err := msgCard.AddSection(alertRequestSummarySection); err != nil {
		errMsg := fmt.Sprintf("Error returned from attempt to add alertRequestSummarySection: %v", err)
		log.Error("createMessage: " + errMsg)
		msgCard.Text = msgCard.Text + "\n\n" + send2teams.TryToFormatAsCodeSnippet(errMsg)
	}

	/*
		Alert Request Headers Section
	*/

	alertRequestHeadersSection := goteamsnotify.NewMessageCardSection()
	alertRequestHeadersSection.StartGroup = true
	alertRequestHeadersSection.Title = "## Alert Request Headers"

	alertRequestHeadersSection.Text = fmt.Sprintf(
		"%d alert request headers provided",
		len(record.Alert.Headers),
	)

	// process alert request headers

	for header, values := range record.Alert.Headers {
		for index, value := range values {
			// update value with code snippet formatting, assign back using
			// the available index value
			values[index] = send2teams.TryToFormatAsCodeSnippet(value)
		}
		addFactPair(&msgCard, alertRequestHeadersSection, header, values...)
	}

	if err := msgCard.AddSection(alertRequestHeadersSection); err != nil {
		errMsg := fmt.Sprintf("Error returned from attempt to add alertRequestHeadersSection: %v", err)
		log.Error("createMessage: " + errMsg)
		msgCard.Text = msgCard.Text + "\n\n" + send2teams.TryToFormatAsCodeSnippet(errMsg)
	}

	/*
		Message Card Branding/Trailer Section
	*/

	trailerSection := goteamsnotify.NewMessageCardSection()
	trailerSection.StartGroup = true
	trailerSection.Text = send2teams.ConvertEOLToBreak(config.MessageTrailer())
	if err := msgCard.AddSection(trailerSection); err != nil {
		errMsg := fmt.Sprintf("Error returned from attempt to add trailerSection: %v", err)
		log.Error("createMessage: " + errMsg)
		msgCard.Text = msgCard.Text + "\n\n" + send2teams.TryToFormatAsCodeSnippet(errMsg)
	}

	return msgCard
}

// define function/wrapper for sending details to Microsoft Teams
func sendMessage(
	ctx context.Context,
	webhookURL string,
	msgCard goteamsnotify.MessageCard,
	schedule time.Time,
	retries int,
	retriesDelay int,
) NotifyResult {

	// Note: We already do validation elsewhere, and the library call does
	// even more validation, but we can handle this obvious empty argument
	// problem directly
	if webhookURL == "" {
		return NotifyResult{
			Err:     fmt.Errorf("sendMessage: webhookURL not defined, skipping message submission to Microsoft Teams channel"),
			Success: false,
		}
	}

	log.Debugf("sendMessage: Time now is %v", time.Now().Format("15:04:05"))
	log.Debugf("sendMessage: Notification scheduled for: %v", schedule.Format("15:04:05"))

	// Set delay timer to meet received notification schedule. This helps
	// ensure that we delay the appropriate amount of time before we make our
	// first attempt at sending a message to Microsoft Teams.
	notificationDelay := time.Until(schedule)

	notificationDelayTimer := time.NewTimer(notificationDelay)
	defer notificationDelayTimer.Stop()
	log.Debugf("sendMessage: notificationDelayTimer created at %v with duration %v",
		time.Now().Format("15:04:05"),
		notificationDelay,
	)

	log.Debug("sendMessage: Waiting for either context or notificationDelayTimer to expire before sending notification")

	select {
	case <-ctx.Done():
		ctxErr := ctx.Err()
		msg := NotifyResult{
			Val: fmt.Sprintf("sendMessage: Received Done signal at %v: %v, shutting down",
				time.Now().Format("15:04:05"),
				ctxErr.Error(),
			),
			Success: false,
		}
		log.Debug(msg.Val)
		return msg

	// Delay between message submission attempts; this will *always*
	// delay, regardless of whether the attempt is the first one or not
	case <-notificationDelayTimer.C:

		log.Debugf("sendMessage: Waited %v before notification attempt at %v",
			notificationDelay,
			time.Now().Format("15:04:05"),
		)

		ctxExpires, ctxExpired := ctx.Deadline()
		if ctxExpired {
			log.Debugf("sendMessage: WaitTimeout context expires at: %v", ctxExpires.Format("15:04:05"))
		}

		// check to see if context has expired during our delay
		if ctx.Err() != nil {
			msg := NotifyResult{
				Val: fmt.Sprintf(
					"sendMessage: context expired or cancelled at %v: %v, attempting to abort message submission",
					time.Now().Format("15:04:05"),
					ctx.Err().Error(),
				),
				Success: false,
			}

			log.Debug(msg.Val)

			return msg
		}

		// Submit message card, retry submission if needed up to specified number
		// of retry attempts.
		if err := send2teams.SendMessage(ctx, webhookURL, msgCard, retries, retriesDelay); err != nil {
			errMsg := NotifyResult{
				Err: fmt.Errorf(
					"sendMessage: ERROR: Failed to submit message to Microsoft Teams at %v: %v",
					time.Now().Format("15:04:05"),
					err,
				),
				Success: false,
			}
			log.Error(errMsg.Err.Error())
			return errMsg
		}

		successMsg := NotifyResult{
			Val: fmt.Sprintf(
				"sendMessage: Message successfully sent to Microsoft Teams at %v",
				time.Now().Format("15:04:05"),
			),
			Success: true,
		}

		// Note success for potential troubleshooting
		log.Debug(successMsg.Val)

		return successMsg

	}

}
