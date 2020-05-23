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

// Credit:

// https://golangcode.com/get-the-request-ip-addr/
// https://github.com/eddturtle/golangcode-site

package events

import (
	"net/http"

	"github.com/apex/log"
)

// GetIP gets a request's IP address by reading off the forwarded-for
// header (for proxies) and falls back to using the remote address.
func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	log.WithFields(log.Fields{
		"forwarded_header": forwarded,
	}).Debug("logging X-FORWARDED-FOR header")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}
