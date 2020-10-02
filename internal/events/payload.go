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

import "net/http"

// SplunkSampleAlertPayload maps to the sample JSON payload provided by the
// Splunk webhook documentation. This payload is submitted via webhook request
// from a Splunk Alert action. The specific fields were constructed using
// https://mholt.github.io/json-to-go/
type SplunkSampleAlertPayload struct {
	Result struct {
		Sourcetype string `json:"sourcetype"`
		Count      string `json:"count"`
	} `json:"result"`
	Sid         string      `json:"sid"`
	ResultsLink string      `json:"results_link"`
	SearchName  interface{} `json:"search_name"`
	Owner       string      `json:"owner"`
	App         string      `json:"app"`
}

// SplunkAlertPayloadV1 maps to a captured JSON payload submitted by a webhook
// request from a test alert on 2020-02-12. The specific fields were
// constructed using the following web app as a starting point and then
// massaging the fieldnames to avoid conflicts:
// https://mholt.github.io/json-to-go/
type SplunkAlertPayloadV1 struct {
	ResultsLink string `json:"results_link"`
	Result      struct {
		ContextData       string   `json:"contextData"`
		DateMonth         string   `json:"date_month"`
		Forcecdn          string   `json:"forcecdn"`
		HTTPStatusCode    string   `json:"http_status_code"`
		Bkt               string   `json:"_bkt"`
		Indextime         string   `json:"_indextime"`
		Kv                string   `json:"_kv"`
		Linecount         string   `json:"linecount"`
		EzproxyTime       string   `json:"ezproxy_time"`
		Serial            string   `json:"_serial"`
		Time              string   `json:"_time"`
		Eventtype         string   `json:"eventtype"`
		EventtypeColor    string   `json:"_eventtype_color"`
		Sp                string   `json:"sp"`
		Bhskip            string   `json:"bhskip"`
		BhskipSourcetype  string   `json:"_sourcetype"`
		Punct             string   `json:"punct"`
		SplunkServer      string   `json:"splunk_server"`
		Session           string   `json:"session"`
		Host              string   `json:"host"`
		URL               string   `json:"url"`
		Srcip             string   `json:"srcip"`
		SplunkServerGroup string   `json:"splunk_server_group"`
		Cd                string   `json:"_cd"`
		Si                []string `json:"_si"`
		TagEventtype      string   `json:"tag::eventtype"`
		Timestartpos      string   `json:"timestartpos"`
		DateHour          string   `json:"date_hour"`
		DateSecond        string   `json:"date_second"`
		Timeendpos        string   `json:"timeendpos"`
		Username          string   `json:"username"`
		DateMinute        string   `json:"date_minute"`
		DateMday          string   `json:"date_mday"`
		Index             string   `json:"index"`
		TimeZoneID        string   `json:"timeZoneId"`
		Sourcetype        string   `json:"sourcetype"`
		Rs                string   `json:"rs"`
		TransitionType    string   `json:"transitionType"`
		DateWday          string   `json:"date_wday"`
		Source            string   `json:"source"`
		DateZone          string   `json:"date_zone"`
		Tag               string   `json:"tag"`
		DateYear          string   `json:"date_year"`
		UserAgent         string   `json:"user_agent"`
		Raw               string   `json:"_raw"`
		ResourceURL       string   `json:"URL"`
		Vr                string   `json:"vr"`
	} `json:"result"`
	Sid        string `json:"sid"`
	Owner      string `json:"owner"`
	App        string `json:"app"`
	SearchName string `json:"search_name"`
}

// SplunkAlertPayloadV2 maps (loosely) to a captured JSON payload submitted by
// a webhook request from a test alert on 2020-02-12. We've removed fields
// from this struct that we are choosing to ignore from the Splunk payload.
type SplunkAlertPayloadV2 struct {
	ResultsLink string `json:"results_link"`
	Result      struct {
		// What is this built from? The captured payload has an empty value.
		// Session           string   `json:"session"`

		SourceIP string `json:"srcip"`
		Username string `json:"username"`

		ResourceURL    string `json:"URL"`
		HTTPStatusCode string `json:"http_status_code"`
		UserAgent      string `json:"user_agent"`

		// TODO: Do we need this?
		TagEventtype string `json:"tag::eventtype"`

		// Splunk software stores timestamp values in the _time field, in
		// Coordinated Universal Time (UTC) format.
		// https://docs.splunk.com/Documentation/Splunk/latest/Data/HowSplunkextractstimestamps
		Time string `json:"_time"`

		EzproxyTime string `json:"ezproxy_time"` // original log date/time
		DateHour    string `json:"date_hour"`    // EZproxy parsed date/time fields
		DateSecond  string `json:"date_second"`
		DateMinute  string `json:"date_minute"`
		DateMday    string `json:"date_mday"`
		DateYear    string `json:"date_year"`
		DateWday    string `json:"date_wday"`
		DateMonth   string `json:"date_month"`
		DateZone    string `json:"date_zone"` // A time zone offset in minutes from UTC

		// TODO: These could potentially be useful to identify what data
		// source was consulted in order to generate the alert. This may be
		// (more) relevant once we ingest (or look at) the audit log data
		// and/or any other relevant source.
		Index      string `json:"index"`
		Sourcetype string `json:"sourcetype"`
		Source     string `json:"source"`

		// TODO: Do we need this?
		SplunkServer string `json:"splunk_server"`

		// TODO: What is Splunk considering this?
		Bkt string `json:"_bkt"`

		// TODO: Is there anything useful for this? Presumably everything
		// that comes from Splunk will show the same tag for our group?
		Tag string `json:"tag"`

		// TODO: Record this "archival" copy of the raw data?
		Raw string `json:"_raw"`
	} `json:"result"`

	// TODO: Are these three fields needed for anything?
	Sid   string `json:"sid"`
	Owner string `json:"owner"`
	App   string `json:"app"`

	// TODO: Use this to explain *why* a user account has been disabled?
	SearchName string `json:"search_name"`
}

// SplunkAlertEvent is a subset of the original alert payload received.
// TODO: Have ArrivalTime as time.Time type? Force formatting in template
// itself?
// TODO: Rename to `Alert` ?
type SplunkAlertEvent struct {

	// Username is the username reported by Splunk and represents a user logged
	// into EZproxy.
	Username string

	// UserIP is the IP Address of the user logged into EZproxy.
	UserIP string

	// PayloadSenderIP is the IP Address of the system submitting the payload.
	PayloadSenderIP string

	// ArrivalTime is the time when the Splunk alert was received.
	ArrivalTime string

	// LocalTime is the time when the Splunk alert was received recorded in
	// 24hr local time. This is a workaround for Teams choosing to ignore
	// time.RFC3339 designation that I encountered while developing
	// atc0005/bounce.
	LocalTime string

	// AlertName is the name of the Splunk alert.
	AlertName string

	// SearchID is the unique identifier for the Splunk search associated with
	// the alert.
	SearchID string

	// EndpointPath is the handler path where the payload was received.
	EndpointPath string

	// HTTPMethod is the HTTP verb or method used by the alert sender. POST is
	// the only supported HTTP method.
	HTTPMethod string

	// Headers is a set of HTTP headers sent with the alert payload.
	Headers http.Header
}
