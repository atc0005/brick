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
	"errors"
	"fmt"
	"os"
	"path/filepath"
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

	log.Debugf("%s: Request to open %q received", myFuncName, filename)
	log.Debugf("%s: Attempting to open sanitized version of file %q",
		myFuncName, filepath.Clean(filename))

	// TODO: How do we handle the situation where the file does not exist
	// ahead of time? Since this application will manage the file, it should
	// be able to create it with the desired permissions?
	f, err := os.Open(filepath.Clean(filename))
	if err != nil {
		return false, fmt.Errorf(
			"%s: error encountered opening file %q: %w",
			myFuncName,
			filename,
			err,
		)
	}
	defer func() {
		if err := f.Close(); err != nil {
			// Ignore "file already closed" errors
			if !errors.Is(err, os.ErrClosed) {
				log.Errorf(
					"%s: failed to close file %q: %s",
					myFuncName,
					filename,
					err.Error(),
				)
			}
		}
	}()

	log.Debugf("%s: Searching for: %q", myFuncName, searchTerm)

	s := bufio.NewScanner(f)
	var lineno int

	// TODO: Does Scan() perform any whitespace manipulation already?
	for s.Scan() {
		lineno++
		currentLine := s.Text()
		log.Debugf(
			"%s: Scanned line %d from %q: %q",
			myFuncName,
			lineno,
			filename,
			currentLine,
		)

		currentLine = strings.TrimSpace(currentLine)
		log.Debugf(
			"%s: Line %d from %q after lowercasing and whitespace removal: %q",
			myFuncName,
			lineno,
			filename,
			currentLine,
		)

		// explicitly ignore lines beginning with specified pattern, if
		// provided
		if ignorePrefix != "" {
			if strings.HasPrefix(currentLine, ignorePrefix) {
				log.Debugf(
					"%s: Ignoring line %d due to leading %q",
					myFuncName,
					lineno,
					ignorePrefix,
				)
				continue
			}
		}

		log.Debugf(
			"%s: Checking whether line %d is a match for %q: %q",
			myFuncName,
			lineno,
			searchTerm,
			currentLine,
		)
		if strings.EqualFold(currentLine, searchTerm) {
			log.Debugf(
				"%s: Match found on line %d, returning true to indicate this",
				myFuncName,
				lineno,
			)
			return true, nil
		}

	}

	log.Debugf("%s: Exited s.Scan() loop", myFuncName)

	// report any errors encountered while scanning the input file
	if err := s.Err(); err != nil {
		return false, err
	}

	// explicitly close file, bail if failure occurs
	if err := f.Close(); err != nil {
		return false, fmt.Errorf(
			"%s: failed to close file %q: %w",
			myFuncName,
			filepath.Clean(filename),
			err,
		)
	}

	// otherwise, report that the requested searchTerm was not found
	return false, nil

}
