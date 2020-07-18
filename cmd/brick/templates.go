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

/*

	Structure of templates:

		Top-level Title

		Introduction block (e.g., msgCard.Text)

		Errors Section

		Disable User Request Details Section - Core of alert details

		Alert Request Summary Section - General client request details

		Alert Request Headers Section

		Branding / Trailer Section

*/

const defaultEmailTemplate string = `

{{ $missingValue := "MISSING VALUE - Please file a bug report!" }}

**Summary**

{{ if ne .Record.Note "" -}}
{{ .Record.Note }}
{{- else -}}
{{ $missingValue }}
{{- end }}


**Disable User Request Errors**

{{ if .Record.Error -}}
{{ .Record.Error }}
{{- else -}}
None
{{- end }}


**Disable User Request Details**

* Username: {{ if .Record.Alert.Username }}{{ .Record.Alert.Username }}{{ else }}{{ $missingValue }}{{ end }}
* User IP: {{ if .Record.Alert.UserIP }}{{ .Record.Alert.UserIP }}{{ else }}{{ $missingValue }}{{ end }}
* Alert/Search Name: {{ if .Record.Alert.AlertName }}{{ .Record.Alert.AlertName }}{{ else }}{{ $missingValue }}{{ end }}
* Alert/Search ID: {{ if .Record.Alert.SearchID }}{{ .Record.Alert.SearchID }}{{ else }}{{ $missingValue }}{{ end }}


**Alert Request Summary**

* Received at: {{ .Record.Alert.LocalTime }}
* Endpoint path: {{ .Record.Alert.EndpointPath }}
* HTTP Method: {{ .Record.Alert.HTTPMethod }}
* Alert Sender IP: {{ .Record.Alert.PayloadSenderIP }}


**Alert Request Headers**
{{ range $key, $slice := .Record.Alert.Headers }}
* {{ $key }}: {{ range $sliceValue := $slice }}{{ . }}{{ end }}
{{- else }}
* None
{{ end }}

{{ .Branding }}


`

const textileEmailTemplate string = `

{{ $missingValue := "MISSING VALUE - Please file a bug report!" }}

**Summary**

<pre>
{{ if ne .Record.Note "" -}}
{{ .Record.Note }}
{{- else -}}
{{ $missingValue }}
{{- end }}
</pre>


**Disable User Request Errors**

{{ if .Record.Error -}}
<pre>
{{ .Record.Error }}
</pre>
{{- else -}}
* None
{{- end }}


**Disable User Request Details**

| Username          | {{ if .Record.Alert.Username }}{{ .Record.Alert.Username }}{{ else }}{{ $missingValue }}{{ end }} |
| User IP           | {{ if .Record.Alert.UserIP }}{{ .Record.Alert.UserIP }}{{ else }}{{ $missingValue }}{{ end }} |
| Alert/Search Name | {{ if .Record.Alert.AlertName }}{{ .Record.Alert.AlertName }}{{ else }}{{ $missingValue }}{{ end }} |
| Alert/Search ID   | {{ if .Record.Alert.SearchID }}{{ .Record.Alert.SearchID }}{{ else }}{{ $missingValue }}{{ end }} |


**Alert Request Summary**

| Received at     | {{ .Record.Alert.LocalTime }} |
| Endpoint path   | {{ .Record.Alert.EndpointPath }} |
| HTTP Method     | {{ .Record.Alert.HTTPMethod }} |
| Alert Sender IP | {{ .Record.Alert.PayloadSenderIP }} |


**Alert Request Headers**
{{ range $key, $slice := .Record.Alert.Headers }}
| {{ $key }} | {{ range $sliceValue := $slice }}{{ . }}{{ end }} |
{{- else }}
| None | N/A |
{{ end }}

{{ .Branding }}

`
