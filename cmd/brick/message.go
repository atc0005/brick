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
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/smtp"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/apex/log"
	"github.com/atc0005/brick/internal/caller"
	"github.com/atc0005/brick/internal/config"
	"github.com/atc0005/brick/internal/events"
	"github.com/atc0005/go-ezproxy"

	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
	"github.com/atc0005/go-teams-notify/v2/messagecard"
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
func addFactPair(msg *messagecard.MessageCard, section *messagecard.Section, key string, values ...string) {

	for idx := range values {
		values[idx] = messagecard.TryToFormatAsCodeSnippet(values[idx])
	}

	if err := section.AddFactFromKeyValue(
		key,
		values...,
	); err != nil {
		from := caller.GetFuncFileLineInfo()
		errMsg := fmt.Sprintf("%s error returned from attempt to add fact from key/value pair: %v", from, err)
		log.Errorf("%s %s", from, errMsg)
		msg.Text = msg.Text + "\n\n" + messagecard.TryToFormatAsCodeSnippet(errMsg)
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

// getMsgSummaryText evaluates the provided event Record and builds a message
// suitable for display as the main or summary notification text. This message
// is generated first from the Note field if available, second from the Error
// field and finally if neither are set a fallback message is used. This
// precedence allows for using a provided Note as a brief summary while still
// using the Error field in a dedicated section of the notification.
func getMsgSummaryText(record events.Record) string {

	// This part of the message is valuable "real estate" for eyeballs; we
	// should ensure we are communicating what just occurred instead of using
	// a mostly static block of text.

	var msgSummaryText string

	switch {
	case record.Note != "":
		msgSummaryText = "Summary: " + record.Note
	case record.Error != nil:
		msgSummaryText = "Error: " + record.Error.Error()

	// Attempting to use an empty string for the top-level message card Text
	// field results in a notification failure, so set *something* to meet
	// those requirements. This "guard rail" is also useful for ensuring that
	// email notifications are similarly provided a fallback message.
	default:
		msgSummaryText = "FIXME: Missing Note and Error for this event record!"
	}

	return msgSummaryText

}

// getMsgTitle is a helper function used to generate the title for outgoing
// notifications. This function uses the provided prefix and event Record to
// generate stable titles reflecting the step in the disable user process at
// which the notification was generated; the intent is to quickly tell where
// the process halted for troubleshooting purposes.
func getMsgTitle(msgTitlePrefix string, record events.Record) string {

	var msgCardTitle string

	switch record.Action {

	// case record.Error != nil:
	// 	msgCardTitle = "[ERROR] " + record.Error.Error()

	// TODO: Calculate step labeling based off of enabled features (see GH-65).

	case events.ActionSuccessDisableRequestReceived, events.ActionFailureDisableRequestReceived:
		msgCardTitle = msgTitlePrefix + recordActionStep2of3 + " " + record.Action

	case events.ActionSuccessDisabledUsername, events.ActionFailureDisabledUsername:
		msgCardTitle = msgTitlePrefix + recordActionStep2of3 + " " + record.Action

	case events.ActionSuccessDuplicatedUsername, events.ActionFailureDuplicatedUsername:
		msgCardTitle = msgTitlePrefix + recordActionStep2of3 + " " + record.Action

	case events.ActionSuccessIgnoredUsername, events.ActionFailureIgnoredUsername:
		msgCardTitle = msgTitlePrefix + recordActionStep2of3 + " " + record.Action

	case events.ActionSuccessIgnoredIPAddress, events.ActionFailureIgnoredIPAddress:
		msgCardTitle = msgTitlePrefix + recordActionStep2of3 + " " + record.Action

	case events.ActionSuccessTerminatedUserSession,
		events.ActionFailureUserSessionLookupFailure,
		events.ActionFailureTerminatedUserSession,
		events.ActionSkippedTerminateUserSessions:
		msgCardTitle = msgTitlePrefix + recordActionStep3of3 + " " + record.Action

	default:
		msgCardTitle = msgTitlePrefix + " " + recordActionUnknownRecord + " " + record.Action
		log.Warnf("UNKNOWN record: %v+\n", record)
	}

	return msgCardTitle
}

// createTeamsMessage receives an event Record and generates a MessageCard
// which is used to generate a Microsoft Teams message.
func createTeamsMessage(record events.Record) *messagecard.MessageCard {

	myFuncName := caller.GetFuncName()

	log.Debugf("%s: alert received: %#v", myFuncName, record)

	// build MessageCard for submission
	msgCard := messagecard.NewMessageCard()

	msgCardTitlePrefix := config.MyAppName + ": "

	msgCard.Title = getMsgTitle(msgCardTitlePrefix, record)

	// msgCard.Text = record.Note
	msgCard.Text = getMsgSummaryText(record)

	/*
		Errors Section
	*/

	disableUserRequestErrors := messagecard.NewSection()
	disableUserRequestErrors.Title = "## Disable User Request Errors"
	disableUserRequestErrors.StartGroup = true

	switch {
	case record.Error != nil:
		addFactPair(msgCard, disableUserRequestErrors, "Error", record.Error.Error())
	case record.Error == nil:
		disableUserRequestErrors.Text = "None"
	}

	if err := msgCard.AddSection(disableUserRequestErrors); err != nil {
		errMsg := fmt.Sprintf("Error returned from attempt to add disableUserRequestErrors: %v", err)
		log.Errorf("%s: %v", myFuncName, errMsg)
		msgCard.Text = msgCard.Text + "\n\n" + messagecard.TryToFormatAsCodeSnippet(errMsg)
	}

	// If Session Termination is enabled, create Termination Results section
	if record.SessionTerminationResults != nil {

		sessionTerminationResultsSection := messagecard.NewSection()
		sessionTerminationResultsSection.Title = "## Session Termination Results"
		sessionTerminationResultsSection.StartGroup = true

		sessionTerminationResultsSection.Text = getTerminationResultsList(record.SessionTerminationResults)

		if err := msgCard.AddSection(sessionTerminationResultsSection); err != nil {
			errMsg := fmt.Sprintf("Error returned from attempt to add sessionTerminationResultsSection: %v", err)
			log.Errorf("%s: %v", myFuncName, errMsg)
			msgCard.Text = msgCard.Text + "\n\n" + messagecard.TryToFormatAsCodeSnippet(errMsg)
		}

	}

	/*
		Disable User Request Details Section - Core of alert details
	*/

	disableUserRequestDetailsSection := messagecard.NewSection()
	disableUserRequestDetailsSection.Title = "## Disable User Request Details"
	disableUserRequestDetailsSection.StartGroup = true

	addFactPair(msgCard, disableUserRequestDetailsSection, "Username", record.Alert.Username)
	addFactPair(msgCard, disableUserRequestDetailsSection, "User IP", record.Alert.UserIP)
	addFactPair(msgCard, disableUserRequestDetailsSection, "Alert/Search Name", record.Alert.AlertName)
	addFactPair(msgCard, disableUserRequestDetailsSection, "Alert/Search ID", record.Alert.SearchID)

	if err := msgCard.AddSection(disableUserRequestDetailsSection); err != nil {
		errMsg := fmt.Sprintf("Error returned from attempt to add disableUserRequestDetailsSection: %v", err)
		log.Errorf("%s: %v", myFuncName, errMsg)
		msgCard.Text = msgCard.Text + "\n\n" + messagecard.TryToFormatAsCodeSnippet(errMsg)
	}

	/*
		Alert Request Summary Section - General client request details
	*/

	alertRequestSummarySection := messagecard.NewSection()
	alertRequestSummarySection.Title = "## Alert Request Summary"
	alertRequestSummarySection.StartGroup = true

	addFactPair(msgCard, alertRequestSummarySection, "Received at", record.Alert.LocalTime)
	addFactPair(msgCard, alertRequestSummarySection, "Endpoint path", record.Alert.EndpointPath)
	addFactPair(msgCard, alertRequestSummarySection, "HTTP Method", record.Alert.HTTPMethod)
	addFactPair(msgCard, alertRequestSummarySection, "Alert Sender IP", record.Alert.PayloadSenderIP)

	if err := msgCard.AddSection(alertRequestSummarySection); err != nil {
		errMsg := fmt.Sprintf("Error returned from attempt to add alertRequestSummarySection: %v", err)
		log.Errorf("%s: %v", myFuncName, errMsg)
		msgCard.Text = msgCard.Text + "\n\n" + messagecard.TryToFormatAsCodeSnippet(errMsg)
	}

	/*
		Alert Request Headers Section
	*/

	alertRequestHeadersSection := messagecard.NewSection()
	alertRequestHeadersSection.StartGroup = true
	alertRequestHeadersSection.Title = "## Alert Request Headers"

	alertRequestHeadersSection.Text = fmt.Sprintf(
		"%d alert request headers provided",
		len(record.Alert.Headers),
	)

	// process alert request headers

	// Create a copy of the original so that we don't modify the original
	// alert headers; other notifications (e.g., email) will need a fresh copy
	// of those values so that any formatting applied here doesn't "spill
	// over" to those notifications.
	requestHeadersCopy := make(http.Header)
	for key, value := range record.Alert.Headers {
		requestHeadersCopy[key] = value
	}

	for header, values := range requestHeadersCopy {

		// As with the enclosing map, we create a copy here so that we don't
		// modify the original (which is used also by email notifications).
		headerValuesCopy := make([]string, len(values))
		copy(headerValuesCopy, values)

		for index, value := range headerValuesCopy {
			// update value with code snippet formatting, assign back using
			// the available index value
			headerValuesCopy[index] = messagecard.TryToFormatAsCodeSnippet(value)
		}
		addFactPair(msgCard, alertRequestHeadersSection, header, headerValuesCopy...)
	}

	if err := msgCard.AddSection(alertRequestHeadersSection); err != nil {
		errMsg := fmt.Sprintf("Error returned from attempt to add alertRequestHeadersSection: %v", err)
		log.Errorf("%s: %v", myFuncName, errMsg)
		msgCard.Text = msgCard.Text + "\n\n" + messagecard.TryToFormatAsCodeSnippet(errMsg)
	}

	/*
		Message Card Branding/Trailer Section
	*/

	trailerSection := messagecard.NewSection()
	trailerSection.StartGroup = true
	trailerSection.Text = messagecard.ConvertEOLToBreak(config.MessageTrailer(config.BrandingMarkdownFormat))
	if err := msgCard.AddSection(trailerSection); err != nil {
		errMsg := fmt.Sprintf("Error returned from attempt to add trailerSection: %v", err)
		log.Errorf("%s: %v", myFuncName, errMsg)
		msgCard.Text = msgCard.Text + "\n\n" + messagecard.TryToFormatAsCodeSnippet(errMsg)
	}

	return msgCard
}

// sendTeamsMessage wraps and orchestrates an external library function call
// to send Microsoft Teams messages. This includes honoring the provided
// schedule in order to comply with remote API rate limits.
func sendTeamsMessage(
	ctx context.Context,
	webhookURL string,
	msgCard *messagecard.MessageCard,
	schedule time.Time,
	retries int,
	retriesDelay int,
) NotifyResult {

	myFuncName := caller.GetFuncName()

	// Note: We already do validation elsewhere, and the library call does
	// even more validation, but we can handle this obvious empty argument
	// problem directly
	if webhookURL == "" {
		return NotifyResult{
			Err: fmt.Errorf(
				"%s: webhookURL not defined, skipping message submission to Microsoft Teams channel",
				myFuncName,
			),
			Success: false,
		}
	}

	log.Debugf("%s: Time now is %v", myFuncName, time.Now().Format("15:04:05"))
	log.Debugf("%s: Notification scheduled for: %v", myFuncName, schedule.Format("15:04:05"))

	// Set delay timer to meet received notification schedule. This helps
	// ensure that we delay the appropriate amount of time before we make our
	// first attempt at sending a message to Microsoft Teams.
	notificationDelay := time.Until(schedule)

	notificationDelayTimer := time.NewTimer(notificationDelay)
	defer notificationDelayTimer.Stop()
	log.Debugf("%s: notificationDelayTimer created at %v with duration %v",
		myFuncName,
		time.Now().Format("15:04:05"),
		notificationDelay,
	)

	log.Debugf(
		"%s: Waiting for either context or notificationDelayTimer to expire before sending notification",
		myFuncName,
	)

	select {
	case <-ctx.Done():
		ctxErr := ctx.Err()
		msg := NotifyResult{
			Val: fmt.Sprintf("%s: Received Done signal at %v: %v, shutting down",
				myFuncName,
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

		log.Debugf("%s: Waited %v before notification attempt at %v",
			myFuncName,
			notificationDelay,
			time.Now().Format("15:04:05"),
		)

		ctxExpires, ctxExpired := ctx.Deadline()
		if ctxExpired {
			log.Debugf(
				"%s: WaitTimeout context expires at: %v",
				myFuncName,
				ctxExpires.Format("15:04:05"),
			)
		}

		// check to see if context has expired during our delay
		if ctx.Err() != nil {
			msg := NotifyResult{
				Val: fmt.Sprintf(
					"%s: context expired or cancelled at %v: %v, attempting to abort message submission",
					myFuncName,
					time.Now().Format("15:04:05"),
					ctx.Err().Error(),
				),
				Success: false,
			}

			log.Debug(msg.Val)

			return msg
		}

		// Create Microsoft Teams client
		mstClient := goteamsnotify.NewTeamsClient()

		// Submit message card using Microsoft Teams client, retry submission
		// if needed up to specified number of retry attempts.
		if err := mstClient.SendWithRetry(ctx, webhookURL, msgCard, retries, retriesDelay); err != nil {
			errMsg := NotifyResult{
				Err: fmt.Errorf(
					"%s: ERROR: Failed to submit message to Microsoft Teams at %v: %v",
					myFuncName,
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
				"%s: Message successfully sent to Microsoft Teams at %v",
				myFuncName,
				time.Now().Format("15:04:05"),
			),
			Success: true,
		}

		// Note success for potential troubleshooting
		log.Debug(successMsg.Val)

		return successMsg

	}

}

// emailConfig represents user-provided settings specific to sending email
// notifications. This is mostly a helper for passing email settings "down"
// until sufficient work on GH-22 can be applied.
type emailConfig struct {
	server                 string
	serverPort             int
	senderAddress          string
	recipientAddresses     []string
	clientIdentity         string
	timeout                time.Duration
	notificationRateLimit  time.Duration
	notificationRetries    int
	notificationRetryDelay int
	template               *template.Template
}

// createEmailMessage receives an event record and a collection of settings
// used to generate a formatted email message.
func createEmailMessage(record events.Record, emailCfg emailConfig) string {

	myFuncName := caller.GetFuncName()

	emailSubjectPrefix := config.MyAppName + ": "

	// getMsgTitle() func is used for both email and Teams messages
	emailSubject := getMsgTitle(emailSubjectPrefix, record)

	// Use this function to generate the summary instead of directly
	// referencing the Record.Note field; the Record.Note field may not always
	// have a value (GH-134).
	emailSummary := getMsgSummaryText(record)

	data := struct {
		Record       events.Record
		EmailSubject string
		EmailSummary string
		Branding     string
	}{
		Record:       record,
		EmailSubject: emailSubject,
		EmailSummary: emailSummary,
		Branding:     config.MessageTrailer(config.BrandingTextileFormat),
	}

	var renderedTmpl bytes.Buffer
	var emailBody string

	tmplErr := emailCfg.template.Execute(&renderedTmpl, data)
	switch {
	case tmplErr != nil:
		errMsg := fmt.Sprintf(
			"Error returned from attempt to parse email template: %v",
			tmplErr,
		)
		log.Errorf("%s: %v", myFuncName, errMsg)

		emailBody = errMsg
	default:
		emailBody = renderedTmpl.String()
	}

	email := fmt.Sprintf(
		"To: %s\r\n"+
			"From: %s\r\n"+
			"Subject: %s\r\n"+
			"\r\n"+
			"%s\r\n",
		strings.Join(emailCfg.recipientAddresses, ", "),
		emailCfg.senderAddress,
		emailSubject,
		emailBody,
	)

	return email

}

// sendEmail is an analogue of the abstraction/functionality provided by
// goteamsnotify.SendMessage(...). The plan is to refactor this function as part
// of the work for GH-22.
func sendEmail(
	ctx context.Context,
	emailCfg emailConfig,
	emailMsg string,
) error {

	myFuncName := caller.GetFuncName()

	smtpServer := fmt.Sprintf("%s:%d", emailCfg.server, emailCfg.serverPort)

	// FIXME: This function both logs *and* returns the error, which is
	// duplication that will require fixing at some point. Leaving both in for
	// the time being until this code proves stable.
	send := func(ctx context.Context, emailCfg emailConfig, emailMsg string, myFuncName string) error {

		// Connect to the remote SMTP server.
		c, dialErr := smtp.Dial(smtpServer)
		if dialErr != nil {
			errMsg := fmt.Errorf(
				"%s: failed to connect to SMTP server %q on port %v: %w",
				myFuncName,
				emailCfg.server,
				emailCfg.serverPort,
				dialErr,
			)
			log.Error(errMsg.Error())

			return errMsg
		}

		// At this point we have a Client, so we need to ensure that the QUIT
		// command is sent to the SMTP server to clean up; close connection and
		// send the QUIT command.
		defer func() {
			if err := c.Quit(); err != nil {

				fmt.Printf("Error type: %+v", err)

				errMsg := fmt.Errorf(
					"%s: failure occurred sending QUIT command to %q on port %v: %w",
					myFuncName,
					emailCfg.server,
					emailCfg.serverPort,
					err,
				)
				log.Error(errMsg.Error())

				return
			}

			log.Debugf(
				"%s: Successfully sent QUIT command to %q on port %v",
				myFuncName,
				emailCfg.server,
				emailCfg.serverPort,
			)
		}()

		// Use user-specified (or default) client identity in our greeting to the
		// SMTP server.
		log.Debugf(
			"%s: Sending greeting to SMTP server %q on port %v with identity of %q",
			myFuncName,
			emailCfg.server,
			emailCfg.serverPort,
			emailCfg.clientIdentity,
		)
		if err := c.Hello(emailCfg.clientIdentity); err != nil {

			errMsg := fmt.Errorf(
				"%s: failure occurred sending greeting to SMTP server %q on port %v: %w",
				myFuncName,
				emailCfg.server,
				emailCfg.serverPort,
				err,
			)
			log.Error(errMsg.Error())

			return errMsg
		}

		// Set the sender
		if err := c.Mail(emailCfg.senderAddress); err != nil {
			errMsg := fmt.Errorf(
				"%s: failed to set sender address %q for email: %w",
				myFuncName,
				emailCfg.senderAddress,
				err,
			)
			log.Error(errMsg.Error())

			return errMsg
		}

		// Process one or more user-provided destination email addresses
		for _, emailAddr := range emailCfg.recipientAddresses {
			if err := c.Rcpt(emailAddr); err != nil {
				errMsg := fmt.Errorf(
					"%s: failed to set recipient address %q: %w",
					myFuncName,
					emailAddr,
					err,
				)
				log.Error(errMsg.Error())

				return errMsg
			}
		}

		// Send the email body.
		//
		// Data issues a DATA command to the server and returns a writer that can
		// be used to write the mail headers and body. The caller should close the
		// writer before calling any more methods on c. A call to Data must be
		// preceded by one or more calls to Rcpt.
		wc, dataErr := c.Data()
		if dataErr != nil {
			errMsg := fmt.Errorf(
				"%s: failure occurred when sending DATA command to SMTP server %q on port %v: %w",
				myFuncName,
				emailCfg.server,
				emailCfg.serverPort,
				dataErr,
			)
			log.Error(errMsg.Error())

			return dataErr
		}

		defer func() {

			if err := wc.Close(); err != nil {

				fmt.Printf("Error type: %+v", err)

				errMsg := fmt.Errorf(
					"%s: failure occurred closing mail headers and body writer: %w",
					myFuncName,
					err,
				)
				log.Error(errMsg.Error())

				return
			}

			log.Debugf(
				"%s: Successfully closed mail headers and body writer",
				myFuncName,
			)
		}()

		if _, err := fmt.Fprint(wc, emailMsg); err != nil {
			errMsg := fmt.Errorf(
				"%s: failure occurred when writing message to connection for SMTP server %q on port %v: %w",
				myFuncName,
				emailCfg.server,
				emailCfg.serverPort,
				dataErr,
			)
			log.Error(errMsg.Error())

			return dataErr

		}

		return nil
	}

	var result error

	// initial attempt + number of specified retries
	attemptsAllowed := 1 + emailCfg.notificationRetries

	ourRetryDelay := time.Duration(emailCfg.notificationRetryDelay) * time.Second

	// attempt to send message to Microsoft Teams, retry specified number of
	// times before giving up
	for attempt := 1; attempt <= attemptsAllowed; attempt++ {

		// Check here at the start of the loop iteration (either first or
		// subsequent) in order to return early in an effort to prevent
		// undesired message attempts after the context has been cancelled.
		if ctx.Err() != nil {
			errMsg := fmt.Errorf(
				"%s: context cancelled or expired: %v; aborting message submission after %d of %d attempts",
				myFuncName,
				ctx.Err().Error(),
				attempt,
				attemptsAllowed,
			)

			// If this is set, we're looking at the second (incomplete)
			// iteration. Let's combine our error above with the last result
			// in an effort to provide a more meaningful error.
			if result != nil {
				// TODO: How to properly combine these two errors?
				errMsg = fmt.Errorf("%w: %v", result, errMsg)
			}

			log.Error(errMsg.Error())
			return errMsg
		}

		result = send(ctx, emailCfg, emailMsg, myFuncName)
		if result != nil {

			errMsg := fmt.Errorf(
				"%s: Attempt %d of %d to send message failed: %v",
				myFuncName,
				attempt,
				attemptsAllowed,
				result,
			)

			log.Error(errMsg.Error())

			// apply retry delay if our context hasn't been cancelled yet,
			// otherwise continue with the loop to allow context cancellation
			// handling logic to be applied at the top of the loop
			if ctx.Err() == nil {
				log.Debugf(
					"%s: Context not cancelled yet, applying retry delay of %v",
					myFuncName,
					ourRetryDelay,
				)
				time.Sleep(ourRetryDelay)
			}

			// retry send attempt (if attempts remain)
			continue

		}

		log.Debugf(
			"%s: successfully sent message after %d of %d attempts\n",
			myFuncName,
			attempt,
			attemptsAllowed,
		)

		// break out of retry loop: we're done!
		break

	}

	return result

}

// sendEmailMessage wraps and orchestrates another function call to send email
// messages. This includes honoring the provided schedule in order to comply
// with remote API rate limits. The plan is to refactor this function as part
// of the work for GH-22.
func sendEmailMessage(
	ctx context.Context,
	emailCfg emailConfig,
	emailMsg string,
	schedule time.Time,

) NotifyResult {

	myFuncName := caller.GetFuncName()

	log.Debugf("%s: Time now is %v", myFuncName, time.Now().Format("15:04:05"))
	log.Debugf("%s: Notification scheduled for: %v", myFuncName, schedule.Format("15:04:05"))

	// Set delay timer to meet received notification schedule. This helps
	// ensure that we delay the appropriate amount of time before we make our
	// first attempt at sending a message to the specified SMTP server.
	notificationDelay := time.Until(schedule)

	notificationDelayTimer := time.NewTimer(notificationDelay)
	defer notificationDelayTimer.Stop()
	log.Debugf("%s: notificationDelayTimer created at %v with duration %v",
		myFuncName,
		time.Now().Format("15:04:05"),
		notificationDelay,
	)

	log.Debugf(
		"%s: Waiting for either context or notificationDelayTimer to expire before sending notification",
		myFuncName,
	)

	select {
	case <-ctx.Done():
		ctxErr := ctx.Err()
		msg := NotifyResult{
			Val: fmt.Sprintf("%s: Received Done signal at %v: %v, shutting down",
				myFuncName,
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

		log.Debugf("%s: Waited %v before notification attempt at %v",
			myFuncName,
			notificationDelay,
			time.Now().Format("15:04:05"),
		)

		ctxExpires, ctxExpired := ctx.Deadline()
		if ctxExpired {
			log.Debugf(
				"%s: WaitTimeout context expires at: %v",
				myFuncName,
				ctxExpires.Format("15:04:05"),
			)
		}

		// check to see if context has expired during our delay
		if ctx.Err() != nil {
			msg := NotifyResult{
				Val: fmt.Sprintf(
					"%s: context expired or cancelled at %v: %v, attempting to abort message submission",
					myFuncName,
					time.Now().Format("15:04:05"),
					ctx.Err().Error(),
				),
				Success: false,
			}

			log.Debug(msg.Val)

			return msg
		}

		if err := sendEmail(ctx, emailCfg, emailMsg); err != nil {

			errMsg := NotifyResult{
				Err: fmt.Errorf(
					"%s: ERROR: Failed to submit message to %s on port %v at %v: %v",
					myFuncName,
					emailCfg.server,
					emailCfg.serverPort,
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
				"%s: Message successfully sent to SMTP server %q on port %v at %v",
				myFuncName,
				emailCfg.server,
				emailCfg.serverPort,
				time.Now().Format("15:04:05"),
			),
			Success: true,
		}

		// Note success for potential troubleshooting
		log.Debug(successMsg.Val)

		return successMsg

	}

}
