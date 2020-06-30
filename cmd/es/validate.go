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
)

func validate(config AppConfig) error {

	if config.Username == "" {
		return fmt.Errorf(
			"error: missing username",
		)
	}

	if config.ActiveFilePath == "" {
		return fmt.Errorf(
			"error: missing filename",
		)
	}

	if config.EZproxyExecutable == "" {
		return fmt.Errorf(
			"error: missing EZproxy executable name",
		)
	}

	if config.SearchDelay < 0 {
		return fmt.Errorf(
			"%d is not a valid number of seconds for search delay",
			config.SearchDelay,
		)
	}

	if config.SearchRetries < 0 {
		return fmt.Errorf(
			"%d is not a valid number of search retries",
			config.SearchRetries,
		)
	}

	return nil

}
