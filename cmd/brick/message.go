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
	"runtime"
	"time"

	"github.com/apex/log"
	"github.com/atc0005/brick/config"
	"github.com/atc0005/brick/events"

	// use our fork for now until recent work can be submitted for inclusion
	// in the upstream project
	goteamsnotify "github.com/atc0005/go-teams-notify"

	send2teams "github.com/atc0005/send2teams/teams"
)

func createMessage(record events.Record) goteamsnotify.MessageCard {

	log.Debugf("createMessage: alert received: %#v", record)

	// FIXME: Pull this out as a separate helper function?
	// FIXME: Rework and offer upstream?
	addFactPair := func(msg *goteamsnotify.MessageCard, section *goteamsnotify.MessageCardSection, key string, values ...string) {

		// attempt to format all values as code in an effort to keep Teams
		// from parsing some special characters as Markdown code formatting
		// characters
		for idx := range values {
			values[idx] = send2teams.TryToFormatAsCodeSnippet(values[idx])
		}

		if err := section.AddFactFromKeyValue(
			key,
			values...,
		); err != nil {

			// runtime.Caller(skip int) (pc uintptr, file string, line int, ok bool)
			_, file, line, ok := runtime.Caller(0)
			from := fmt.Sprintf("createMessage [file %s, line %d]:", file, line)
			if !ok {
				from = "createMessage:"
			}
			errMsg := fmt.Sprintf("%s error returned from attempt to add fact from key/value pair: %v", from, err)
			log.Errorf("%s %s", from, errMsg)
			msg.Text = msg.Text + "\n\n" + send2teams.TryToFormatAsCodeSnippet(errMsg)
		}
	}

	// build MessageCard for submission
	msgCard := goteamsnotify.NewMessageCard()

	msgCardTitlePrefix := config.MyAppName + ": "

	if record.Action == "" {
		msgCard.Title = msgCardTitlePrefix + "(1/2) Disable user account request received"
	} else {
		msgCard.Title = msgCardTitlePrefix + "(2/2) " + record.Action
	}

	msgCard.Text = fmt.Sprintf(
		"Disable request received for user %s at IP %s by alert %s",
		send2teams.TryToFormatAsCodeSnippet(record.Alert.Username),
		send2teams.TryToFormatAsCodeSnippet(record.Alert.UserIP),
		send2teams.TryToFormatAsCodeSnippet(record.Alert.AlertName),
	)

	/*
		Record/Annotations Section
	*/

	// TODO: Flesh this section out more by reviewing/updating app code to
	// provide more Error and Note details

	if record.Error != nil || record.Note != "" {
		disableUserRequestAnnotationsSection := goteamsnotify.NewMessageCardSection()
		disableUserRequestAnnotationsSection.Title = "## Disable User Request Annotations"
		disableUserRequestAnnotationsSection.StartGroup = true

		// TODO: Should we display these fields regardless?

		if record.Error != nil {
			addFactPair(&msgCard, disableUserRequestAnnotationsSection, "Error", record.Error.Error())
		}

		if record.Note != "" {
			disableUserRequestAnnotationsSection.Text = "Note: " + record.Note
		}

		if err := msgCard.AddSection(disableUserRequestAnnotationsSection); err != nil {
			errMsg := fmt.Sprintf("Error returned from attempt to add disableUserRequestAnnotationsSection: %v", err)
			log.Error("createMessage: " + errMsg)
			msgCard.Text = msgCard.Text + "\n\n" + send2teams.TryToFormatAsCodeSnippet(errMsg)
		}
	}

	/*
		Disable User Request Details Section - Core of SplunkAlertEvent details
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
