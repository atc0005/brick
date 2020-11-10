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
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/atc0005/brick/internal/config"
	"github.com/atc0005/brick/internal/events"
	"github.com/atc0005/brick/internal/files"
	goteamsnotify "github.com/atc0005/go-teams-notify/v2"

	"github.com/apex/log"
)

// See splunk-test-submissions.json and splunk-test-.http for sample test data

func main() {

	// Toggle debug logging from library packages as needed to troubleshoot
	// implementation work
	goteamsnotify.DisableLogging()

	// Emulate returning exit code from main function by "queuing up" a
	// default exit code that matches expectations, but allow explicitly
	// setting the exit code in such a way that is compatible with using
	// deferred function calls throughout the application.
	var appExitCode *int
	defer func(code *int) {
		var exitCode int
		if code != nil {
			exitCode = *code
		}
		os.Exit(exitCode)
	}(appExitCode)

	appConfig, err := config.NewConfig()
	if err != nil {
		log.Errorf("Failed to initialize application: %s", err)
		*appExitCode = 1
		return
	}
	log.Debug("Initializing application")

	log.Debugf("AppConfig: %+v", appConfig)

	mux := http.NewServeMux()

	// Apply "default" timeout settings provided by Simon Frey; override the
	// default "wait forever" configuration.
	// FIXME: Refine these settings to apply values more appropriate for a
	// small-to-medium on-premise API (e.g., not over a public Internet link
	// where clients are expected to be slow)
	httpServer := &http.Server{
		ReadHeaderTimeout: config.HTTPServerReadHeaderTimeout,
		ReadTimeout:       config.HTTPServerReadTimeout,
		WriteTimeout:      config.HTTPServerWriteTimeout,
		Handler:           mux,
		Addr:              fmt.Sprintf("%s:%d", appConfig.LocalIPAddress(), appConfig.LocalTCPPort()),
	}

	// Create context that can be used to cancel background jobs.
	ctx, cancel := context.WithCancel(context.Background())

	// Defer cancel() to cover edge cases where it might not otherwise be
	// called
	defer cancel()

	// Use signal.Notify() to send a message on dedicated channel when when
	// interrupt is received (e.g., Ctrl+C) so that we can cleanly shutdown
	// the application.
	//
	// Q: Why are these channels buffered?
	// A: In order to make them asynchronous.
	// Per Bakul Shah (golang-nuts/QEORIGKZO24): In general, synchronize only
	// when you have to. Here the main thread wants to know when the worker
	// thread terminates but the worker thread doesn't care when the main
	// thread gets around to reading from "done". Using a 1 deep buffer
	// channel exactly captures this usage pattern. An unbuffered channel
	// would make the worker thread "rendezvous" with the main thread, which
	// is unnecessary.
	// done := make(chan struct{}, 1)
	//
	// NOTE: Setting up a separate done channel for notify mgr and another
	// for when the http server has been shutdown.
	httpDone := make(chan struct{}, 1)
	notifyDone := make(chan struct{}, 1)
	quit := make(chan os.Signal, 1)

	// override default Go handling of specified signals in order to customize
	// the shutdown process
	signal.Notify(quit,
		syscall.SIGINT,  // stop
		syscall.SIGTERM, // full restart
	)

	// Where events will be sent for processing. We use a buffered channel in
	// an effort to reduce the delay for client requests.
	notifyWorkQueue := make(chan events.Record, config.NotifyMgrQueueDepth)

	// Create "notifications manager" function as persistent goroutine to
	// process incoming notification requests.
	go NotifyMgr(ctx, appConfig, notifyWorkQueue, notifyDone)

	// Setup "listener" to cancel the parent context when Signal.Notify()
	// indicates that SIGINT has been received
	go shutdownListener(ctx, quit, cancel)

	// Setup "listener" to shutdown the running http server when
	// the parent context has been cancelled
	go gracefulShutdown(ctx, httpServer, config.HTTPServerShutdownTimeout, httpDone)

	// build objects representing output files, the templates used to generate
	// those files and any files containing "ignored" users/IPs using our
	// newly constructed config object.
	reportedUserEventsLog := files.NewReportedUserEventsLog(
		appConfig.ReportedUsersLogFile(),
		appConfig.ReportedUsersLogFilePermissions(),
	)

	disabledUsers := files.NewDisabledUsers(
		appConfig.DisabledUsersFile(),
		appConfig.DisabledUsersFileEntrySuffix(),
		appConfig.DisabledUsersFilePermissions(),
	)

	ignoredSources := files.NewIgnoredSources(
		appConfig.IgnoredUsersFile(),
		appConfig.IgnoredIPAddressesFile(),
		appConfig.IgnoreLookupErrors(),
	)

	// log this to help troubleshoot why payloads are (or are not) filtered
	switch {
	case appConfig.RequireTrustedPayloadSender():
		log.Info("OK: Restricting payload sender IP Addresses enabled")
	default:
		log.Warn("CAUTION: Restricting payload sender IP Addresses disabled")
	}

	// GET requests
	mux.HandleFunc(frontpageEndpointPattern, frontPageHandler)
	mux.HandleFunc(apiV1ViewDisabledUsersEndpointPattern, viewDisabledUsersHandler)
	mux.HandleFunc(apiV1ViewDisabledUsersStatusEndpointPattern, viewDisabledUserStatusHandler)

	// POST request
	mux.HandleFunc(
		apiV1DisableUserEndpointPattern,
		disableUserHandler(
			appConfig.RequireTrustedPayloadSender(),
			appConfig.TrustedIPAddresses(),
			reportedUserEventsLog,
			disabledUsers,
			ignoredSources,
			notifyWorkQueue,
			appConfig.EZproxyTerminateSessions(),
			appConfig.EZproxyActiveFilePath(),
			appConfig.EZproxySearchDelay(),
			appConfig.EZproxySearchRetries(),
			appConfig.EZproxyExecutablePath(),
		),
	)

	// listen on specified port and IP Address, block until app is terminated
	log.Infof("%s %s is listening on %s port %d",
		config.MyAppName,
		config.Version,
		appConfig.LocalIPAddress(),
		appConfig.LocalTCPPort(),
	)

	// TODO: This can be handled in a cleaner fashion?
	if err := httpServer.ListenAndServe(); err != nil {

		// Calling Shutdown() will immediately return ErrServerClosed, but
		// based on reading the docs it sounds like any errors from closing
		// connections will instead overwrite this default error message with
		// a real one, so receiving ErrServerClosed can be treated as a
		// "successful shutdown" message of sorts, so ignore it and look for
		// any other error message.
		if !errors.Is(err, http.ErrServerClosed) {
			log.Errorf("error occurred while running httpServer: %v", err)
			*appExitCode = 1
			return
		}
	}

	log.Debug("Waiting on gracefulShutdown completion signal")
	<-httpDone
	log.Debug("Received gracefulShutdown completion signal")

	log.Debug("Waiting on NotifyMgr completion signal")
	<-notifyDone
	log.Debug("Received NotifyMgr completion signal")

	log.Infof("%s successfully shutdown", config.MyAppName)
}
