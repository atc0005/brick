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
	"fmt"
	"net/http"

	"github.com/apex/log"
)

func viewDisabledUsersHandler(w http.ResponseWriter, r *http.Request) {

	log.Info("viewDisabledUsersHandler endpoint hit")
	fmt.Fprintf(w, "viewDisabledUsersHandler endpoint hit")
}

func viewDisabledUserStatusHandler(w http.ResponseWriter, r *http.Request) {

	// TODO: Pull historical disable status details from database
	// TODO: Pull current disable status from flat-files
	//  * Check file that this web app maintains
	//  * Check file that this web app knows about, but does not maintain
	// Examples:
	// users.disabled.txt (manual)
	// users.brick-disabled.txt (automatic)

	log.Info("viewDisabledUserStatusHandler endpoint hit")
	fmt.Fprintf(w, "viewDisabledUserStatusHandler endpoint hit")
}
