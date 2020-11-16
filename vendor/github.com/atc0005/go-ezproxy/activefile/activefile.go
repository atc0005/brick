// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/go-ezproxy
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

package activefile

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/atc0005/go-ezproxy"
)

const (
	// SessionLinePrefix is a single letter prefix found at the start of all
	// lines containing a session ID (among other details).
	SessionLinePrefix string = "S"

	// SessionLineEvenNumbered indicates that this line should be found on
	// even numbered lines.
	SessionLineEvenNumbered bool = true

	// UsernameLinePrefix is a single letter prefix found at the start of all
	// lines containing a username.
	UsernameLinePrefix string = "L"

	// UsernameLineEvenNumbered indicates that this line should be found on
	// odd numbered lines.
	UsernameLineEvenNumbered bool = false
)

// SessionEntry reflects a line in the ezproxy.hst file that contains session
// information. We have to tie this information back to a specific username
// based on line ordering. The session line comes first in the set followed by
// one or more additional lines, one of which contains the username.
type SessionEntry struct {

	// Type is the first field in the file. Observed entries thus far are "P,
	// M, H, S, L, g, s". Those lines relevant to our purposes of matching
	// session IDs to usernames are "S" for Session and "L" for Login.
	Type string

	// SessionID is the second field for a line in the  ActiveFile that starts
	// with capital letter 'S'. We need to tie this back to a specific
	// username in order to reliably terminate active sessions.
	SessionID string

	// IPAddress is the seventh field for a line in the ActiveFile that starts
	// with capital letter 'S'. We *could* use this value to determine which
	// session ID to terminate, though using this value from a remote
	// payload/report by itself has a greater chance of terminating the wrong
	// user session.
	IPAddress string
}

// activeFileReader represents a file reader specific to the EZProxy active
// users and hosts file.
type activeFileReader struct {
	// SearchDelay is the intentional delay between each attempt to open and
	// search the specified filename for the specified username.
	SearchDelay time.Duration

	// SearchRetries is the number of additional search attempts that will be
	// made whenever the initial search attempt returns zero results. Each
	// attempt to read the active file is subject to a race condition; EZproxy
	// does not immediately write session information to disk when creating or
	// terminating sessions, so some amount of delay and a number of retry
	// attempts are used in an effort to work around that write delay.
	SearchRetries int

	// Username is the name of the user account to search for within the
	// specified file.
	Username string

	// Filename is the name of the file which will be parsed/searched for the
	// specified username.
	Filename string
}

// NewReader creates a new instance of a SessionReader that provides access to
// a collection of user sessions for the specified username.
func NewReader(username string, filename string) (ezproxy.SessionsReader, error) {

	if username == "" {
		return nil, errors.New(
			"func NewReader: missing username",
		)
	}

	if filename == "" {
		return nil, errors.New(
			"func NewReader: missing filename",
		)
	}

	reader := activeFileReader{
		SearchDelay:   ezproxy.DefaultSearchDelay,
		SearchRetries: ezproxy.DefaultSearchRetries,
		Username:      username,
		Filename:      filename,
	}

	return &reader, nil
}

// filterEntries is a helper function that returns all entries from the
// provided active file that have the required line prefix. Other methods
// handle converting these entries to UserSession values.
func (afr activeFileReader) filterEntries(validPrefixes []string) ([]ezproxy.FileEntry, error) {

	ezproxy.Logger.Printf(
		"filterEntries: Request to open %q received\n",
		afr.Filename,
	)
	ezproxy.Logger.Printf(
		"filterEntries: Attempting to open sanitized version of file %q\n",
		filepath.Clean(afr.Filename),
	)

	f, err := os.Open(filepath.Clean(afr.Filename))
	if err != nil {
		return nil, fmt.Errorf("func filterEntries: error encountered opening file %q: %w", afr.Filename, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			// Ignore "file already closed" errors
			if !errors.Is(err, os.ErrClosed) {
				ezproxy.Logger.Printf(
					"filterEntries: failed to close file %q: %s",
					afr.Filename,
					err.Error(),
				)
			}
		}
	}()

	s := bufio.NewScanner(f)

	var lineno int

	var validLines []ezproxy.FileEntry

	// TODO: Does Scan() perform any whitespace manipulation already?
	for s.Scan() {
		lineno++
		currentLine := s.Text()
		// ezproxy.Logger.Printf("Scanned line %d from %q: %q\n", lineno, filename, currentLine)

		currentLine = strings.TrimSpace(currentLine)
		// ezproxy.Logger.Printf("Line %d from %q after whitespace removal: %q\n",
		// 	lineno, filename, currentLine)

		if currentLine != "" {
			for _, validPrefix := range validPrefixes {
				if strings.HasPrefix(currentLine, validPrefix) {
					validLines = append(validLines, ezproxy.FileEntry{
						Text:   currentLine,
						Number: lineno,
					})
				}
			}
		}
	}

	ezproxy.Logger.Printf("Exited s.Scan() loop")

	// report any errors encountered while scanning the input file
	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("func filterEntries: errors encountered while scanning the input file: %w", err)
	}

	// explicitly close file, bail if failure occurs
	if err := f.Close(); err != nil {
		return nil, fmt.Errorf(
			"func filterEntries: failed to close file %q: %w",
			afr.Filename,
			err,
		)
	}

	return validLines, nil
}

// SetSearchRetries is a helper method for setting the number of additional
// retries allowed when receiving zero search results.
func (afr *activeFileReader) SetSearchRetries(retries int) error {
	if retries < 0 {
		return fmt.Errorf("func SetSearchRetries: %d is not a valid number of search retries", retries)
	}

	afr.SearchRetries = retries

	return nil
}

// SetSearchDelay is a helper method for setting the delay in seconds between
// search attempts.
func (afr *activeFileReader) SetSearchDelay(delay int) error {
	if delay < 0 {
		return fmt.Errorf("func SetSearchDelay: %d is not a valid number of seconds for search delay", delay)
	}

	afr.SearchDelay = time.Duration(delay) * time.Second

	return nil
}

// AllUserSessions returns a list of all session IDs along with their associated
// IP Address in the form of a slice of UserSession values. This list of
// session IDs is intended for further processing such as filtering to a
// specific username or aggregating to check thresholds.
func (afr activeFileReader) AllUserSessions() (ezproxy.UserSessions, error) {

	// Lines containing the session entries
	validPrefixes := []string{
		SessionLinePrefix,
		UsernameLinePrefix,
	}

	var allUserSessions ezproxy.UserSessions

	fileEntryDelimiter := " "

	validLines, filterErr := afr.filterEntries(validPrefixes)
	if filterErr != nil {
		return nil, fmt.Errorf(
			"failed to filter active file entries while generating list of user sessions: %w",
			filterErr,
		)
	}

	// Ensure that the gathered lines consist of pairs, otherwise we are
	// likely dealing with an invalid active users file. At this point we
	// should bail as continuing would likely mean identifying the wrong user
	// session for termination.
	if !(len(validLines)%2 == 0) {
		errMsg := fmt.Sprintf(
			"error: Incomplete data pairs (%d lines) found in file %q while searching for %q user sessions",
			len(validLines),
			afr.Filename,
			afr.Username,
		)
		ezproxy.Logger.Println(errMsg)
		return nil, errors.New(errMsg)
	}

	for idx, currentLine := range validLines {

		activeFileEntry := strings.Split(currentLine.Text, fileEntryDelimiter)
		lineno := currentLine.Number
		switch activeFileEntry[0] {
		case SessionLinePrefix:
			// line 1 of 2 (even numbered idx)
			// session ID as field 2, IP Address as field 7

			// if not even numbered line, but the username only occurs on even
			// numbered lines
			// if !(idx%2 == 0) {
			if (idx%2 == 1) && (SessionLineEvenNumbered) {

				// We have found a "S" prefixed line out of expected order.
				// This suggests an invalid active users file or a bug in the
				// earlier application logic used when generating the list of
				// valid file entries.

				errMsg := fmt.Sprintf(
					"error: Unexpected data pair ordering encountered at line %d in the active users file %q while searching for %q; "+
						"session line is odd numbered",
					lineno,
					afr.Filename,
					afr.Username,
				)
				ezproxy.Logger.Println(errMsg)
				return nil, errors.New(errMsg)
			}

			allUserSessions = append(allUserSessions, ezproxy.UserSession{
				SessionID: activeFileEntry[1],
				IPAddress: activeFileEntry[6],
			})
		case UsernameLinePrefix:
			// line 2 of 2 (odd numbered idx)
			// username as field 2

			if (idx%2 == 1) && (UsernameLineEvenNumbered) {

				// We have found a "L" prefixed line out of expected order.
				// This suggests an invalid active users file or a bug in the
				// earlier application logic used when generating the list of
				// valid file entries.

				errMsg := fmt.Sprintf(
					"error: Unexpected data pair ordering encountered at line %d in the active users file %q while searching for %q; "+
						"session line is odd numbered",
					lineno,
					afr.Filename,
					afr.Username,
				)
				ezproxy.Logger.Println(errMsg)
				return nil, errors.New(errMsg)
			}

			// Use the length of the collected user sessions minus 1 as the
			// index into the allUserSessions slice. The intent is to get
			// access to the partial ActiveUserSession that was just
			// constructed from the previous 'S' line in order to include the
			// username alongside the existing Session ID and IP Address
			// fields.
			prevSessionIdx := len(allUserSessions) - 1
			if prevSessionIdx < 0 {

				ezproxy.Logger.Printf(
					"Current text from username line %d: %v",
					lineno,
					activeFileEntry,
				)

				errMsg := fmt.Sprintf(
					"error: unable to update partial ActiveUserSession from line %d; "+
						"unable to reliably determine session ID for %q",
					lineno-1,
					afr.Username,
				)
				ezproxy.Logger.Println(errMsg)
				return nil, errors.New(errMsg)
			}
			allUserSessions[prevSessionIdx].Username = activeFileEntry[1]
		default:
			continue
		}

	}

	ezproxy.Logger.Printf(
		"Found %d active sessions\n",
		len(allUserSessions),
	)

	return allUserSessions, nil

}

// MatchingUserSessions uses the previously provided username to return a list
// of all matching session IDs along with their associated IP Address in the
// form of a slice of UserSession values.
func (afr activeFileReader) MatchingUserSessions() (ezproxy.UserSessions, error) {

	// What we will return to the the caller
	requestedUserSessions := make([]ezproxy.UserSession, 0, ezproxy.SessionsLimit)

	searchAttemptsAllowed := afr.SearchRetries + 1

	// Perform the search up to X times
	for searchAttempts := 1; searchAttempts <= searchAttemptsAllowed; searchAttempts++ {

		ezproxy.Logger.Printf(
			"Beginning search attempt %d of %d for %q\n",
			searchAttempts,
			searchAttemptsAllowed,
			afr.Username,
		)

		// Intentional delay in an effort to better avoid stale data due to
		// potential race condition with EZproxy write delays.
		ezproxy.Logger.Printf(
			"Intentionally delaying for %v to help avoid race condition due to delayed EZproxy writes\n",
			afr.SearchDelay,
		)
		time.Sleep(afr.SearchDelay)

		allUserSessions, err := afr.AllUserSessions()
		if err != nil {
			return nil, fmt.Errorf(
				"func UserSessions: failed to retrieve all user sessions in order to filter to specific username: %w",
				err,
			)
		}

		// filter all user sessions found earlier just to the requested user
		for _, session := range allUserSessions {
			if strings.EqualFold(afr.Username, session.Username) {
				requestedUserSessions = append(requestedUserSessions, session)
			}
		}

		// skip further attempts to find sessions if we already found some
		if len(requestedUserSessions) > 0 {
			break
		}

		// try again (unless we hit our limit)
		continue

	}

	ezproxy.Logger.Printf(
		"Found %d active sessions for %q\n",
		len(requestedUserSessions),
		afr.Username,
	)

	return requestedUserSessions, nil

}
