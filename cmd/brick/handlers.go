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
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/apex/log"

	"github.com/atc0005/brick/internal/events"
	"github.com/atc0005/brick/internal/files"
)

// API endpoint patterns supported by this application
//
// TODO: Find a better location for these values
const (
	frontpageEndpointPattern                    string = "/"
	apiV1DisableUserEndpointPattern             string = "/api/v1/users/disable"
	apiV1ViewDisabledUsersEndpointPattern       string = "/api/v1/users/list"
	apiV1ViewDisabledUsersStatusEndpointPattern string = "/api/v1/users/status"
)

// frontPageHandler is our catch-all handler. By default it tells clients to
// get off its lawn.
func frontPageHandler(w http.ResponseWriter, r *http.Request) {

	ctxLog := log.WithFields(log.Fields{
		"url_path": r.URL.Path,
	})

	ctxLog.Debug("frontPageHandler endpoint hit")

	if r.Method != http.MethodGet {

		ctxLog.WithFields(log.Fields{
			"http_method": r.Method,
		}).Debug("non-GET request received on GET-only endpoint")
		errorMsg := fmt.Sprintf(
			"Sorry, this endpoint only accepts %s requests. "+
				"You submitted a %s request to %q. "+
				"Please see the README for a list of available endpoints "+
				"(and their supported HTTP methods) and then try again.",
			http.MethodGet, r.Method, r.URL.Path,
		)
		// TODO: Can apex/log hook into this and handle output?
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		fmt.Fprint(w, errorMsg)
		return
	}

	// https://github.com/golang/go/issues/4799
	// https://github.com/golang/go/commit/1a819be59053fa1d6b76cb9549c9a117758090ee
	if r.URL.Path != "/" {
		ctxLog.Debug("Rejecting request not explicitly handled by a route")
		http.NotFound(w, r)
		return
	}

	ctxLog.Debug("Sending Forbidden status code; we are not currently auto-indexing supported endpoints")
	http.Error(
		w,
		"Please see the README for available endpoints and supported HTTP methods and then try again.",
		http.StatusForbidden,
	)

}

func disableUserHandler(
	reportedUserEventsLog *files.ReportedUserEventsLog,
	disabledUsers *files.DisabledUsers,
	ignoredSources files.IgnoredSources,
	notifyWorkQueue chan<- events.Record,
	terminateSessions bool,
	ezproxyActiveFilePath string,
	ezproxySessionsSearchDelay int,
	ezproxySessionSearchRetries int,
	ezproxyExecutable string,
) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// fmt.Fprintf(mw, "disableUserHandler endpoint hit\n")
		log.Debug("disableUserHandler handler hit")

		if r.Method != http.MethodPost {

			log.WithFields(log.Fields{
				"url_path":    r.URL.Path,
				"http_method": r.Method,
			}).Debug("non-POST request received on POST-only endpoint")
			errorMsg := fmt.Sprintf(
				"Sorry, this endpoint only accepts %s requests. "+
					"Please see the README for examples and then try again.",
				http.MethodPost,
			)
			// TODO: Can apex/log hook into this and handle output?
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			fmt.Fprint(w, errorMsg)
			return
		}

		// Limit request body to 1 MB
		r.Body = http.MaxBytesReader(w, r.Body, 1*MB)

		// read everything from the (size-limited) request body so that we
		// have the option of displaying it in a raw format (e.g.,
		// troubleshooting), replace the Body with a new io.ReadCloser to
		// allow later access to r.Body for JSON-decoding purposes
		requestBody, requestBodyReadErr := ioutil.ReadAll(r.Body)
		if requestBodyReadErr != nil {
			http.Error(w, requestBodyReadErr.Error(), http.StatusBadRequest)
			return
		}

		log.Debugf("raw requestBody: %s", requestBody)

		r.Body = ioutil.NopCloser(bytes.NewReader(requestBody))

		// Try to decode the request body into the struct. If there is an
		// error, respond to the client with the error message and appropriate
		// status code.
		var payloadV2 events.SplunkAlertPayloadV2
		if err := json.NewDecoder(r.Body).Decode(&payloadV2); err != nil {
			log.Errorf("Error decoding r.Body into payloadV2:\n%v\n\n", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Debugf("disableUserHandler: Splunk Alert payload decoded into v2 format:\n%+v\n\n", payloadV2)

		// Validate payload fields by ensuring that *something* is present for
		// all fields that we've included in SplunkAlertPayloadV2
		if err := events.ValidatePayload(payloadV2); err != nil {
			log.Error(err.Error())

			// Inform Splunk that we received an invalid payload
			http.Error(w, err.Error(), http.StatusBadRequest)
			return

		}

		// Explicitly confirm that the payload was received so that the sender
		// can go ahead and disconnect. This prevents holding up the sender
		// while this application performs further (unrelated from the
		// sender's perspective) processing.
		//
		// FIXME: Is having a newline here best practice, or no?
		if _, err := io.WriteString(w, "OK: Payload received\n"); err != nil {
			log.Error("disableUserHandler: Failed to send OK status response to payload sender")
		}

		// Manually flush http.ResponseWriter in an additional effort to
		// prevent undue wait time for payload sender
		if f, ok := w.(http.Flusher); ok {
			log.Debug("disableUserHandler: Manually flushing http.ResponseWriter")
			f.Flush()
		} else {
			log.Warn("disableUserHandler: http.Flusher interface not available, cannot flush http.ResponseWriter")
			log.Warn("disableUserHandler: Not flushing http.ResponseWriter may cause a noticeable delay between requests")
		}

		// if we made it this far, the payload checks out and we should be
		// able to safely retrieve values that we need. We will also append
		// payload sender metadata values such as headers, endpoint path, etc
		// so that we can report those later.
		alert := events.SplunkAlertEvent{
			Username:        payloadV2.Result.Username,
			UserIP:          payloadV2.Result.SourceIP,
			PayloadSenderIP: events.GetIP(r),
			ArrivalTime:     time.Now().Format(time.RFC3339),
			LocalTime:       time.Now().Format("2006-01-02 15:04:05"),
			AlertName:       payloadV2.SearchName,
			SearchID:        payloadV2.Sid,
			EndpointPath:    r.URL.Path,
			HTTPMethod:      r.Method,
			Headers:         r.Header,
		}

		// All return values from subfunction calls are dropped into the
		// notifyWorkQueue channel; nothing is returned here for further
		// processing.
		//
		// NOTE: Because this is executed in a goroutine, the client (e.g.,
		// monitoring system) gets a near-immediate response back and the
		// connection is closed. There are probably other/better ways to
		// achieve that specific result without using a goroutine, but the
		// effect is worth noting for further exploration later.
		go files.ProcessDisableEvent(
			alert,
			disabledUsers,
			reportedUserEventsLog,
			ignoredSources,
			notifyWorkQueue,
			terminateSessions,
			ezproxyActiveFilePath,
			ezproxySessionsSearchDelay,
			ezproxySessionSearchRetries,
			ezproxyExecutable,
		)

	}
}
