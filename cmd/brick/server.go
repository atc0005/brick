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
	"net/http"
	"os"
	"time"

	"github.com/apex/log"
	"github.com/atc0005/brick/internal/config"
)

// shutdownListener listens for an os.Signal on the provided quit channel.
// When this signal is received, the provided parent context cancel() function
// is used to cancel all child contexts. This is intended to be run as a
// goroutine.
func shutdownListener(ctx context.Context, quit <-chan os.Signal, parentContextCancel context.CancelFunc) {

	// FIXME: If we're passing in the parent context's CancelFunc, do we need
	// the `ctx` that we're passing in?

	// monitor for shutdown signal
	osSignal := <-quit

	log.Debugf("shutdownListener: Received shutdown signal: %v", osSignal)

	// Attempt to trigger a cancellation of the parent context
	log.Debug("shutdownListener: Cancelling context ...")
	parentContextCancel()
	log.Debug("shutdownListener: context canceled")

}

// gracefullShutdown listens for a context cancellation and then shuts down
// the running http server. Once the http server is shutdown, this function
// signals back that work is complete by closing the provided done channel.
// This function is intended to be run as a goroutine.
func gracefulShutdown(ctx context.Context, server *http.Server, timeout time.Duration, done chan<- struct{}) {

	log.Debug("gracefulShutdown: started; now waiting on <-ctx.Done()")

	// monitor for cancellation context
	<-ctx.Done()

	log.Debugf("gracefulShutdown: context is done: %v", ctx.Err())
	log.Warnf("%s is shutting down, please wait ...", config.MyAppName)

	// Disable HTTP keep-alives to prevent connections from persisting
	server.SetKeepAlivesEnabled(false)

	ctxShutdown, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// From the docs:
	// Shutdown returns the context's error, otherwise it returns any error
	// returned from closing the Server's underlying Listener(s).
	if err := server.Shutdown(ctxShutdown); err != nil {
		log.Errorf("gracefulShutdown: could not gracefully shutdown the server: %v", err)
	}
	close(done)
}
