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

package files

// NOTE: time.RFC3339 format should be used for flat-file log messages in
// order to increase fail2ban parsing reliability

const disabledUsersFileTemplateText string = `
# Username "{{ .Username }}" from source IP "{{ .UserIP }}" disabled at "{{ .ArrivalTime }}" per alert "{{ .AlertName }}" received by "{{ .PayloadSenderIP }}" (SearchID: "{{ .SearchID }}")
{{ .Username }}{{ .EntrySuffix }}
`

// This is a standard message and only indicates that a report was received,
// not that a user was disabled. This message should be followed by another
// message indicating whether the user was disabled or ignored
const reportedUserEventTemplateText string = `{{ .ArrivalTime }} [REPORTED] Username "{{ .Username }}" from source IP "{{ .UserIP }}" reported via alert "{{ .AlertName }}" received by "{{ .PayloadSenderIP }}" (SearchID: "{{ .SearchID }}")
`

const disabledUserFirstEventTemplateText string = `{{ .ArrivalTime }} [DISABLED] Username "{{ .Username }}" from source IP "{{ .UserIP }}" disabled due to alert "{{ .AlertName }}" received by "{{ .PayloadSenderIP }}" (SearchID: "{{ .SearchID }}")
`

const disabledUserRepeatEventTemplateText string = `{{ .ArrivalTime }} [DISABLED] Username "{{ .Username }}" from source IP "{{ .UserIP }}" already disabled, but would be again due to alert "{{ .AlertName }}" received by "{{ .PayloadSenderIP }}" (SearchID: "{{ .SearchID }}")
`

// NOTE: This template is used for ignored users and IP Addresses based on
// presence in the ignored users list and the ignored IP Addresses list.
const ignoredUserEventTemplateText string = `{{ .ArrivalTime }} [IGNORED] Username "{{ .Username }}" from source IP "{{ .UserIP }}" ignored per entry in "{{ .IgnoredEntriesFile }}" (SearchID: "{{ .SearchID }}")
`
