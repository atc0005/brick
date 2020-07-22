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

package fileutils

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/apex/log"
	"github.com/atc0005/brick/internal/caller"
)

// HasLine accepts a search string, an optional pattern to ignore and a
// fully-qualified path to a file containing a list of such strings (e.g.,
// commonly usernames or single IP Addresses), one per line. Lines beginning
// with the optional ignore pattern (e.g., a `#` character) are ignored.
// Leading and trailing whitespace per line is ignored.
func HasLine(searchTerm string, ignorePrefix string, filename string) (bool, error) {

	myFuncName := caller.GetFuncName()

	log.Debugf("Attempting to open %q", filename)

	// TODO: How do we handle the situation where the file does not exist
	// ahead of time? Since this application will manage the file, it should
	// be able to create it with the desired permissions?
	f, err := os.Open(filename)
	if err != nil {
		return false, fmt.Errorf("error encountered opening file %q: %w", filename, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Errorf(
				"%s: failed to close file %q: %s",
				myFuncName,
				err.Error(),
			)
		}
	}()

	log.Debugf("Searching for: %q", searchTerm)

	s := bufio.NewScanner(f)
	var lineno int

	// TODO: Does Scan() perform any whitespace manipulation already?
	for s.Scan() {
		lineno++
		currentLine := s.Text()
		log.Debugf("Scanned line %d from %q: %q\n", lineno, filename, currentLine)

		currentLine = strings.TrimSpace(currentLine)
		log.Debugf("Line %d from %q after lowercasing and whitespace removal: %q\n",
			lineno, filename, currentLine)

		// explicitly ignore lines beginning with specified pattern, if
		// provided
		if ignorePrefix != "" {
			if strings.HasPrefix(currentLine, ignorePrefix) {
				log.Debugf("Ignoring line %d due to leading %q",
					lineno, ignorePrefix)
				continue
			}
		}

		log.Debugf("Checking whether line %d is a match for %q: %q",
			lineno, searchTerm, currentLine)
		if strings.EqualFold(currentLine, searchTerm) {
			log.Debugf("Match found on line %d, returning true to indicate this", lineno)
			return true, nil
		}

	}

	log.Debug("Exited s.Scan() loop")

	// report any errors encountered while scanning the input file
	if err := s.Err(); err != nil {
		return false, err
	}

	// otherwise, report that the requested searchTerm was not found
	return false, nil

}
