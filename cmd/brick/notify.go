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
	"strings"
	"text/template"
	"time"

	"github.com/apex/log"
	"github.com/atc0005/brick/config"
	"github.com/atc0005/brick/events"
)

// NotifyResult wraps the results of notification operations to make it easier
// to inspect the status of various tasks so that we can take action on either
// error or success conditions
type NotifyResult struct {

	// Val is the non-error condition message to return from a notification
	// operation
	Val string

	// Err is the error condition message to return from a notification
	// operation
	Err error

	// Success indicates whether the notification attempt succeeded or if it
	// failed for one reason or another (remote API, timeout, cancellation,
	// etc)
	Success bool
}

// NotifyQueue represents a channel used to queue input data and responses
// between the main application, the notifications manager and "notifiers".
type NotifyQueue struct {
	// The name of a queue. This is intended for display in log messages or
	// other output to identify queues with pending items.
	Name string

	// Channel is a channel used to transport input data and responses.
	Channel interface{}

	// Count is the number of items currently in the queue
	Count int

	// Capacity is the maximum number of items allowed in the queue
	Capacity int
}

// NotifyStats is a collection of stats for Teams and Email notifications
type NotifyStats struct {

	// These fields are collected directly
	IncomingMsgReceived int
	TeamsMsgSent        int
	TeamsMsgSuccess     int
	TeamsMsgFailure     int
	EmailMsgSent        int
	EmailMsgSuccess     int
	EmailMsgFailure     int

	// These fields are calculated from collected field values
	TeamsMsgPending int
	EmailMsgPending int

	TotalPendingMsg int
	TotalSuccessMsg int
	TotalFailureMsg int
}

// newNotifyScheduler takes a time.Duration value as a delay and returns a
// function that can be used to generate a new notification schedule. Each
// call to this function will produce a new schedule incremented by the
// time.Duration delay value. The intent is to provide an easy to use
// mechanism for delaying notifications to remote systems (e.g., in order to
// respect remote API limits).
func newNotifyScheduler(delay time.Duration) func() time.Time {

	log.Debugf("newNotifyScheduler: Initializing lastNotificationSchedule at %s",
		time.Now().Format("15:04:05"),
	)
	lastNotificationSchedule := time.Now()

	return func() time.Time {

		// if we haven't sent a message in a while we should make ensure
		// that we do not return a "next schedule" that has already passed
		if !lastNotificationSchedule.After(time.Now()) {

			expiredSchedule := lastNotificationSchedule.Add(delay)

			log.Debugf(
				"Expired next schedule: [Now: %v, Last: %v, Next: %v]",
				time.Now().Format("15:04:05.000"),
				lastNotificationSchedule.Format("15.04:05.000"),
				expiredSchedule.Format("15:04:05.000"),
			)

			replacementSchedule := time.Now().Add(delay)

			log.Debugf(
				"Replace expired schedule (%v) by resetting the schedule to now (%v) + delay (%v): %v",
				expiredSchedule.Format("15:04:05.000"),
				time.Now().Format("15:04:05.000"),
				delay,
				replacementSchedule.Format("15:04:05"),
			)

			lastNotificationSchedule = replacementSchedule

			return replacementSchedule
		}

		nextSchedule := lastNotificationSchedule.Add(delay)

		log.Debugf(
			"Next schedule not expired: [Last: %v, Now: %v, Next: %v]",
			lastNotificationSchedule.Format("15:04:05"),
			time.Now().Format("15:04:05"),
			nextSchedule.Format("15:04:05"),
		)

		lastNotificationSchedule = nextSchedule

		return nextSchedule
	}
}

// notifyStatsMonitor accepts a context, a delay and a channel for NotifyStats
// values in order to collect and emit summary information for notifications.
// This function is intended to be run as a goroutine.
func notifyStatsMonitor(ctx context.Context, delay time.Duration, statsQueue <-chan NotifyStats) {

	log.Debug("notifyStatsMonitor: Running")

	// this will be populated using values received on statsQueue
	stats := NotifyStats{}

	for {
		t := time.NewTimer(delay)

		// log.Debug("notifyQueueMonitor: Starting loop")

		// block until:
		//	- context cancellation
		//	- timer fires
		select {
		case <-ctx.Done():
			t.Stop()
			log.Debugf(
				"notifyStatsMonitor: Received Done signal: %v, shutting down ...",
				ctx.Err().Error(),
			)

			return

		// emit stats summary here
		case <-t.C:

			ctxLog := log.WithFields(log.Fields{
				"timestamp":  time.Now().Format("15:04:05"),
				"emit_stats": delay,
			})

			ctxLog.Infof(
				"notifyStatsMonitor: Total: "+
					"[%d received, %d pending, %d success, %d failure]",
				stats.IncomingMsgReceived,
				stats.TotalPendingMsg,
				stats.TotalSuccessMsg,
				stats.TotalFailureMsg,
			)

			ctxLog.Infof(
				"notifyStatsMonitor: Teams: "+
					"[%d total, %d pending, %d success, %d failure]",
				stats.TeamsMsgSent,
				stats.TeamsMsgPending,
				stats.TeamsMsgSuccess,
				stats.TeamsMsgFailure,
			)

			ctxLog.Infof(
				"notifyStatsMonitor: Email: "+
					"[%d total, %d pending, %d success, %d failure]",
				stats.EmailMsgSent,
				stats.EmailMsgPending,
				stats.EmailMsgSuccess,
				stats.EmailMsgFailure,
			)

		// received stats update; update our totals
		case statsUpdate := <-statsQueue:

			stats.IncomingMsgReceived += statsUpdate.IncomingMsgReceived

			stats.TeamsMsgSent += statsUpdate.TeamsMsgSent
			stats.TeamsMsgSuccess += statsUpdate.TeamsMsgSuccess
			stats.TeamsMsgFailure += statsUpdate.TeamsMsgFailure

			stats.EmailMsgSent += statsUpdate.EmailMsgSent
			stats.EmailMsgSuccess += statsUpdate.EmailMsgSuccess
			stats.EmailMsgFailure += statsUpdate.EmailMsgFailure

			// calculate non-collected stats here
			stats.TeamsMsgPending = stats.TeamsMsgSent -
				(stats.TeamsMsgSuccess + stats.TeamsMsgFailure)

			stats.EmailMsgPending = stats.EmailMsgSent -
				(stats.EmailMsgSuccess + stats.EmailMsgFailure)

			stats.TotalPendingMsg = stats.EmailMsgPending + stats.TeamsMsgPending
			stats.TotalFailureMsg = stats.EmailMsgFailure + stats.TeamsMsgFailure
			stats.TotalSuccessMsg = stats.EmailMsgSuccess + stats.TeamsMsgSuccess

		}
	}
}

// notifyQueueMonitor accepts a context, a delay and one or many NotifyQueue
// values to monitor for items yet to be processed. This function is intended
// to be run as a goroutine.
func notifyQueueMonitor(ctx context.Context, delay time.Duration, notifyQueues ...NotifyQueue) {

	if len(notifyQueues) == 0 {
		log.Error("received empty list of notifyQueues to monitor, exiting")
		return
	}

	log.Debug("notifyQueueMonitor: Running")

	for {

		t := time.NewTimer(delay)

		// log.Debug("notifyQueueMonitor: Starting loop")

		// block until:
		//	- context cancellation
		//	- timer fires
		select {
		case <-ctx.Done():
			t.Stop()
			log.Debugf(
				"notifyQueueMonitor: Received Done signal: %v, shutting down ...",
				ctx.Err().Error(),
			)
			return

		case <-t.C:

			// log.Debug("notifyQueueMonitor: Timer fired")

			// NOTE: Not needed since the channel is already drained as a
			// result of the case statement triggering and draining the
			// channel
			// t.Stop()

			// Attempt to receive message count updates, proceed without them
			// if they're not available

			var itemsFound bool
			//log.Debugf("Length of queues: %d", len(queues))
			for _, notifyQueue := range notifyQueues {

				switch queue := notifyQueue.Channel.(type) {

				// FIXME: Is there a generic way to match any channel type
				// here in order to calculate the length?
				case chan events.Record:
					notifyQueue.Count = len(queue)
					notifyQueue.Capacity = cap(queue)

				case <-chan events.Record:
					notifyQueue.Count = len(queue)
					notifyQueue.Capacity = cap(queue)

				case chan NotifyResult:
					notifyQueue.Count = len(queue)
					notifyQueue.Capacity = cap(queue)

				case chan NotifyStats:
					notifyQueue.Count = len(queue)
					notifyQueue.Capacity = cap(queue)

				default:
					log.Warn("Default case triggered (this should not happen")
					log.Warnf(
						"Unhandled channel: [Name: %s, Type: %T]",
						notifyQueue.Name, notifyQueue,
					)

				}

				// Show stats only for queues with content
				if notifyQueue.Count > 0 {
					itemsFound = true

					log.WithField("timestamp", time.Now().Format("15:04:05")).Debugf(
						"notifyQueueMonitor: %d/%d items in %s, %d goroutines running",
						notifyQueue.Count,
						notifyQueue.Capacity,
						notifyQueue.Name,
						runtime.NumGoroutine(),
					)
					continue
				}

			}

			if !itemsFound {
				log.WithField("timestamp", time.Now().Format("15:04:05")).Debugf(
					"notifyQueueMonitor: 0 items queued, %d goroutines running",
					runtime.NumGoroutine())
			}
		}
	}
}

// teamsNotifier is a persistent goroutine used to receive incoming
// notification requests and spin off goroutines to create and send Microsoft
// Teams messages.
// TODO: Refactor per GH-37
func teamsNotifier(
	ctx context.Context,
	webhookURL string,
	sendTimeout time.Duration,
	sendRateLimit time.Duration,
	retries int,
	retriesDelay int,
	incoming <-chan events.Record,
	notifyMgrResultQueue chan<- NotifyResult,
	done chan<- struct{},
) {

	log.Debug("teamsNotifier: Running")

	// used by goroutines called by this function to return results
	ourResultQueue := make(chan NotifyResult)

	// Setup new scheduler that we can use to add an intentional delay between
	// Microsoft Teams notification attempts. This delay is added in order to
	// rate limit our outgoing messages to comply with remote API limits.
	// https://docs.microsoft.com/en-us/microsoftteams/platform/webhooks-and-connectors/how-to/connectors-using
	notifyScheduler := newNotifyScheduler(sendRateLimit)

	for {

		select {

		case <-ctx.Done():

			ctxErr := ctx.Err()
			result := NotifyResult{
				Val: fmt.Sprintf("teamsNotifier: Received Done signal: %v, shutting down", ctxErr.Error()),
			}
			log.Debug(result.Val)

			log.Debug("teamsNotifier: Sending back results")
			notifyMgrResultQueue <- result

			log.Debug("teamsNotifier: Closing notifyMgrResultQueue channel to signal shutdown")
			close(notifyMgrResultQueue)

			log.Debug("teamsNotifier: Closing done channel to signal shutdown")
			close(done)
			log.Debug("teamsNotifier: done channel closed, returning")
			return

		case record := <-incoming:

			log.Debugf("teamsNotifier: Request received at %v: %#v",
				time.Now(), record)

			log.Debug("Calculating next scheduled notification")
			nextScheduledNotification := notifyScheduler()

			log.Debugf("Now: %v, Next scheduled notification: %v",
				time.Now().Format("15:04:05"),
				nextScheduledNotification.Format("15:04:05"),
			)

			timeoutValue := config.GetNotificationTimeout(
				sendTimeout,
				nextScheduledNotification,
				retries,
				retriesDelay,
			)

			ctx, cancel := context.WithTimeout(ctx, timeoutValue)
			defer cancel()

			log.Debugf("teamsNotifier: child context created with timeout duration %v", timeoutValue)

			// if there is a message waiting *and* ctx.Done() case statements
			// are both valid, either path could be taken. If this one is
			// taken, then the message send timeout will be the only thing
			// forcing the attempt to loop back around and trigger the
			// ctx.Done() path, but only if this one isn't taken again by the
			// random case selection logic
			log.Debug("teamsNotifier: Checking context to determine whether we should proceed")

			if ctx.Err() != nil {
				result := NotifyResult{
					Success: false,
					Val:     "teamsNotifier: context has been cancelled, aborting notification attempt",
				}
				log.Debug(result.Val)
				notifyMgrResultQueue <- result

				continue
			}

			log.Debug("teamsNotifier: context not cancelled, proceeding with notification attempt")

			// launch task in separate goroutine, each with its own schedule
			log.Debug("teamsNotifier: Launching message creation/submission in separate goroutine")

			go func(
				ctx context.Context,
				webhookURL string,
				record events.Record,
				schedule time.Time,
				numRetries int,
				retryDelay int,
				resultQueue chan<- NotifyResult) {

				ourMessage := createTeamsMessage(record)
				resultQueue <- sendTeamsMessage(ctx, webhookURL, ourMessage, schedule, numRetries, retryDelay)

			}(ctx, webhookURL, record, nextScheduledNotification, retries, retriesDelay, ourResultQueue)

		case result := <-ourResultQueue:
			if result.Err != nil {
				log.Errorf("teamsNotifier: Error received from ourResultQueue: %v", result.Err)
			} else {
				log.Debugf("teamsNotifier: OK: non-error status received on ourResultQueue: %v", result.Val)
			}

			notifyMgrResultQueue <- result

		}

	}

}

// emailNotifier is a persistent goroutine used to receive incoming
// notification requests and spin off goroutines to create and send email
// messages.
// TODO: Refactor per GH-37
func emailNotifier(
	ctx context.Context,
	emailCfg emailConfig,
	incoming <-chan events.Record,
	notifyMgrResultQueue chan<- NotifyResult,
	done chan<- struct{},
) {

	log.Debug("emailNotifier: Running")

	// used by goroutines called by this function to return results
	ourResultQueue := make(chan NotifyResult)

	// Setup new scheduler that we can use to add an intentional delay between
	// email notification attempts. This delay is added in order to rate limit
	// our outgoing messages to comply with any destination email server
	// limits.
	notifyScheduler := newNotifyScheduler(emailCfg.notificationRateLimit)

	for {

		select {

		case <-ctx.Done():

			ctxErr := ctx.Err()
			result := NotifyResult{
				Val: fmt.Sprintf("emailNotifier: Received Done signal: %v, shutting down", ctxErr.Error()),
			}
			log.Debug(result.Val)

			log.Debug("emailNotifier: Sending back results")
			notifyMgrResultQueue <- result

			log.Debug("emailNotifier: Closing notifyMgrResultQueue channel to signal shutdown")
			close(notifyMgrResultQueue)

			log.Debug("emailNotifier: Closing done channel to signal shutdown")
			close(done)
			log.Debug("emailNotifier: done channel closed, returning")
			return

		case record := <-incoming:

			log.Debugf("emailNotifier: Request received at %v: %#v",
				time.Now(), record)

			log.Debug("Calculating next scheduled notification")

			nextScheduledNotification := notifyScheduler()

			log.Debugf("Now: %v, Next scheduled notification: %v",
				time.Now().Format("15:04:05"),
				nextScheduledNotification.Format("15:04:05"),
			)

			timeoutValue := config.GetNotificationTimeout(
				emailCfg.timeout,
				nextScheduledNotification,
				emailCfg.notificationRetries,
				emailCfg.notificationRetryDelay,
			)

			ctx, cancel := context.WithTimeout(ctx, timeoutValue)
			defer cancel()

			log.Debugf("emailNotifier: child context created with timeout duration %v", timeoutValue)

			// if there is a message waiting *and* ctx.Done() case statements
			// are both valid, either path could be taken. If this one is
			// taken, then the message send timeout will be the only thing
			// forcing the attempt to loop back around and trigger the
			// ctx.Done() path, but only if this one isn't taken again by the
			// random case selection logic
			log.Debug("emailNotifier: Checking context to determine whether we should proceed")

			if ctx.Err() != nil {
				result := NotifyResult{
					Success: false,
					Val:     "emailNotifier: context has been cancelled, aborting notification attempt",
				}
				log.Debug(result.Val)
				notifyMgrResultQueue <- result

				continue
			}

			log.Debug("emailNotifier: context not cancelled, proceeding with notification attempt")

			// launch task in separate goroutine, each with its own schedule
			log.Debug("emailNotifier: Launching message creation/submission in separate goroutine")

			go func(
				ctx context.Context,
				record events.Record,
				schedule time.Time,
				emailCfg emailConfig,
				resultQueue chan<- NotifyResult,
			) {
				ourMessage := createEmailMessage(record, emailCfg)
				resultQueue <- sendEmailMessage(ctx, emailCfg, ourMessage, schedule)
			}(ctx, record, nextScheduledNotification, emailCfg, ourResultQueue)

		case result := <-ourResultQueue:

			if result.Err != nil {
				log.Errorf("emailNotifier: Error received from ourResultQueue: %v", result.Err)
			} else {
				log.Debugf("emailNotifier: OK: non-error status received on ourResultQueue: %v", result.Val)
			}

			notifyMgrResultQueue <- result

		}
	}

}

// NotifyMgr receives event details from elsewhere in the application and
// sends notifications to any enabled service (e.g., Microsoft Teams).
func NotifyMgr(ctx context.Context, cfg *config.Config, notifyWorkQueue <-chan events.Record, done chan<- struct{}) {

	log.Debug("NotifyMgr: Running")

	// TODO: Refactor as part of GH-22
	//
	// Create separate, buffered channels to hand-off event details for
	// processing for each service, e.g., one channel for Microsoft Teams
	// outgoing notifications, another for email and so on. Buffered channels
	// are used both to enable async tasks and to provide a means of
	// monitoring the number of items queued for each channel; unbuffered
	// channels have a queue depth (and thus length) of 0.
	teamsNotifyWorkQueue := make(chan events.Record, config.NotifyMgrQueueDepth)
	teamsNotifyResultQueue := make(chan NotifyResult, config.NotifyMgrQueueDepth)
	teamsNotifyDone := make(chan struct{})

	emailNotifyWorkQueue := make(chan events.Record, config.NotifyMgrQueueDepth)
	emailNotifyResultQueue := make(chan NotifyResult, config.NotifyMgrQueueDepth)
	emailNotifyDone := make(chan struct{})

	notifyStatsQueue := make(chan NotifyStats, 1)

	if !cfg.NotifyTeams() && !cfg.NotifyEmail() {
		log.Warn("Teams and email notifications from this application are not enabled.")
		log.Debug("NotifyMgr: Teams and email notifications not requested, not starting notifier goroutines")
		// NOTE: Do not return/exit here.
		//
		// We cannot return/exit the function here because NotifyMgr HAS
		// to run in order to keep the notifyWorkQueue from filling up and
		// blocking other parts of this application that send messages to this
		// channel.
	}

	// If enabled, start persistent goroutine to process request details and
	// submit messages to Microsoft Teams.
	if cfg.NotifyTeams() {
		log.Info("NotifyMgr: Teams notifications enabled")
		log.Debug("NotifyMgr: Starting up teamsNotifier")
		go teamsNotifier(
			ctx,
			cfg.TeamsWebhookURL(),
			config.NotifyMgrTeamsNotificationTimeout,
			cfg.TeamsNotificationRateLimit(),
			cfg.TeamsNotificationRetries(),
			cfg.TeamsNotificationRetryDelay(),
			teamsNotifyWorkQueue,
			teamsNotifyResultQueue,
			teamsNotifyDone,
		)
	}

	// If enabled, start persistent goroutine to process request details and
	// submit messages by email.
	if cfg.NotifyEmail() {
		log.Info("NotifyMgr: Email notifications enabled")
		log.Debug("NotifyMgr: Starting up emailNotifier")

		// TODO: Replace with a more dynamic process that allows for use
		// of user-specified, file-based templates. For now, this is the
		// minimum necessary to complete a first pass at GH-3.

		// TODO: Move these to external files
		// activeTemplate := defaultEmailTemplate
		// activeTemplate := textileEmailTemplate

		// FIXME: Keep linter from complaining about this being unused for
		// now.
		_ = defaultEmailTemplate
		activeTemplate := textileEmailTemplate

		emailTemplate := template.Must(
			template.New(
				"emailTemplate",
			).Funcs(template.FuncMap{
				"trim": strings.TrimSpace,
			}).Parse(activeTemplate))

		// TODO: Refactor as fields for new email notifier (not sure of name
		// yet) type as part of GH-22.
		emailCfg := emailConfig{
			server:                 cfg.EmailServer(),
			serverPort:             cfg.EmailServerPort(),
			senderAddress:          cfg.EmailSenderAddress(),
			recipientAddresses:     cfg.EmailRecipientAddresses(),
			clientIdentity:         cfg.EmailClientIdentity(),
			timeout:                config.NotifyMgrEmailNotificationTimeout,
			notificationRateLimit:  cfg.EmailNotificationRateLimit(),
			notificationRetries:    cfg.EmailNotificationRetries(),
			notificationRetryDelay: cfg.EmailNotificationRetryDelay(),
			template:               emailTemplate,
		}

		go emailNotifier(
			ctx,
			emailCfg,
			emailNotifyWorkQueue,
			emailNotifyResultQueue,
			emailNotifyDone,
		)
	}

	// Monitor queues and report stats for each, even if the user has not
	// opted to use notifications. This is done since we are tracking at least
	// one queue (notifyStatsQueue) which is active even with notifiers
	// disabled.
	queuesToMonitor := []NotifyQueue{
		{
			Name:    "notifyWorkQueue",
			Channel: notifyWorkQueue,
		},
		{
			Name:    "emailNotifyWorkQueue",
			Channel: emailNotifyWorkQueue,
		},
		{
			Name:    "emailNotifyResultQueue",
			Channel: emailNotifyResultQueue,
		},
		{
			Name:    "teamsNotifyWorkQueue",
			Channel: teamsNotifyWorkQueue,
		},
		{
			Name:    "teamsNotifyResultQueue",
			Channel: teamsNotifyResultQueue,
		},
		{
			Name:    "notifyStatsQueue",
			Channel: notifyStatsQueue,
		},
	}

	// periodically print current queue items
	go notifyQueueMonitor(
		ctx,
		config.NotifyQueueMonitorDelay,
		queuesToMonitor...,
	)

	// collect and periodically emit summary of notification details
	go notifyStatsMonitor(
		ctx,
		config.NotifyStatsMonitorDelay,
		notifyStatsQueue,
	)

	for {

		select {

		// NOTE: This should ONLY ever be done when shutting down the entire
		// application, as otherwise goroutines associated with client
		// requests will likely hang, likely until client/server timeout
		// settings are reached
		case <-ctx.Done():
			ctxErr := ctx.Err()
			log.Debugf("NotifyMgr: Received Done signal: %v, shutting down ...", ctxErr.Error())

			evalResults := func(queueName string, result NotifyResult) {
				if result.Err != nil {
					log.Errorf("NotifyMgr: Error received from %s: %v", queueName, result.Err)
					return
				}
				log.Debugf("NotifyMgr: OK: non-error status received on %s: %v", queueName, result.Val)
			}

			// Process any waiting results before blocking and waiting
			// on final completion response from notifier goroutines
			if cfg.NotifyTeams() {
				log.Debug("NotifyMgr: Teams notifications are enabled, shutting down teamsNotifier")

				log.Debug("NotifyMgr: Ranging over teamsNotifyResultQueue")
				for result := range teamsNotifyResultQueue {
					evalResults("teamsNotifyResultQueue", result)
				}

				log.Debug("NotifyMgr: Waiting on teamsNotifyDone")
				select {
				case <-teamsNotifyDone:
					log.Debug("NotifyMgr: Received from teamsNotifyDone")
				case <-time.After(config.NotifyMgrServicesShutdownTimeout):
					log.Debug("NotifyMgr: Timeout occurred while waiting for teamsNotifyDone")
					log.Debug("NotifyMgr: Proceeding with shutdown")
				}

			}

			if cfg.NotifyEmail() {
				log.Debug("NotifyMgr: Email notifications are enabled, shutting down emailNotifier")

				log.Debug("NotifyMgr: Ranging over emailNotifyResultQueue")
				for result := range emailNotifyResultQueue {
					evalResults("emailNotifyResultQueue", result)
				}

				log.Debug("NotifyMgr: Waiting on emailNotifyDone")
				select {
				case <-emailNotifyDone:
					log.Debug("NotifyMgr: Received from emailNotifyDone")
				case <-time.After(config.NotifyMgrServicesShutdownTimeout):
					log.Debug("NotifyMgr: Timeout occurred while waiting for emailNotifyDone")
					log.Debug("NotifyMgr: Proceeding with shutdown")
				}

			}

			log.Debug("NotifyMgr: Closing done channel")
			close(done)

			log.Debug("NotifyMgr: About to return")
			return

		case record := <-notifyWorkQueue:

			log.Debug("NotifyMgr: Input received from notifyWorkQueue")

			go func() {
				notifyStatsQueue <- NotifyStats{
					IncomingMsgReceived: 1,
				}
			}()

			// If we don't have *any* notifications enabled we will just
			// discard the item we have pulled from the channel
			if !cfg.NotifyEmail() && !cfg.NotifyTeams() {
				log.Debug("NotifyMgr: Notifications are not currently enabled; ignoring notification request")
				continue
			}

			if cfg.NotifyTeams() {
				log.Debug("NotifyMgr: Creating new goroutine to place record into teamsNotifyWorkQueue")

				// TODO: Perhaps record this *after* sending the record
				// down the teamsNotifyWorkQueue channel? See other cases
				// where we're using the same "record stat, then do it"
				// approach.

				go func() {
					notifyStatsQueue <- NotifyStats{
						TeamsMsgSent: 1,
					}
				}()

				go func() {
					log.Debugf("NotifyMgr: Existing items in teamsNotifyWorkQueue: %d", len(teamsNotifyWorkQueue))
					log.Debug("NotifyMgr: Pending; placing record into teamsNotifyWorkQueue")
					teamsNotifyWorkQueue <- record
					log.Debug("NotifyMgr: Done; placed record into teamsNotifyWorkQueue")
					log.Debugf("NotifyMgr: Items now in teamsNotifyWorkQueue: %d", len(teamsNotifyWorkQueue))
				}()
			}

			if cfg.NotifyEmail() {
				log.Debug("NotifyMgr: Creating new goroutine to place record in emailNotifyWorkQueue")

				go func() {
					notifyStatsQueue <- NotifyStats{
						EmailMsgSent: 1,
					}
				}()

				go func() {
					log.Debugf("NotifyMgr: Existing items in emailNotifyWorkQueue: %d", len(emailNotifyWorkQueue))
					log.Debug("NotifyMgr: Pending; placing record into emailNotifyWorkQueue")
					emailNotifyWorkQueue <- record
					log.Debug("NotifyMgr: Done; placed record into emailNotifyWorkQueue")
					log.Debugf("NotifyMgr: Items now in emailNotifyWorkQueue: %d", len(emailNotifyWorkQueue))
				}()
			}

		case result := <-teamsNotifyResultQueue:

			statsUpdate := NotifyStats{}

			// NOTE: Only consider explicit success, not a non-error condition
			// because cancellations and timeouts are (currently) treated as
			// non-error, but they're not successful notifications.

			if !result.Success {
				if result.Err != nil {
					log.Errorf("NotifyMgr: Error received from teamsNotifyResultQueue: %v", result.Err)
				}
				statsUpdate.TeamsMsgFailure = 1
			}

			if result.Success {
				log.Debugf("NotifyMgr: OK: non-error status received on teamsNotifyResultQueue: %v", result.Val)
				log.Infof("NotifyMgr: %v", result.Val)
				statsUpdate.TeamsMsgSuccess = 1
			}

			//log.Debugf("statsUpdate: %#v", statsUpdate)

			go func() {
				notifyStatsQueue <- statsUpdate
			}()

		case result := <-emailNotifyResultQueue:

			statsUpdate := NotifyStats{}

			// NOTE: Only consider explicit success, not a non-error condition
			// because cancellations and timeouts are (currently) treated as
			// non-error, but they're not successful notifications.

			if !result.Success {
				if result.Err != nil {
					log.Errorf("NotifyMgr: Error received from emailNotifyResultQueue: %v", result.Err)
				}
				statsUpdate.EmailMsgFailure = 1
			}

			if result.Success {
				log.Debugf("NotifyMgr: non-error status received on emailNotifyResultQueue: %v", result.Val)
				log.Infof("NotifyMgr: %v", result.Val)
				statsUpdate.EmailMsgSuccess = 1
			}

			go func() {
				notifyStatsQueue <- statsUpdate
			}()

		}

	}
}
