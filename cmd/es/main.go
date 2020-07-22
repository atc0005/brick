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
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"

	"github.com/atc0005/brick/config"
	"github.com/atc0005/go-ezproxy"
	"github.com/atc0005/go-ezproxy/activefile"
)

// Primarily used with branding
const myAppName string = "es"
const myAppURL string = "https://github.com/atc0005/brick"

// AppConfig represents the configuration used by this application
type AppConfig struct {

	// Username is the name of the user account that we are searching for.
	Username string

	// ActiveFilePath is the fully-qualified path to the Active Users and
	// Hosts "state" file used by EZproxy (and this application) to track
	// current sessions and hosts managed by EZproxy.
	ActiveFilePath string

	// EZproxyExecutable is the fully-qualified path to the EZproxy
	// executable/binary. This is the same executable that starts at boot.
	// This file is usually named 'ezproxy'.
	EZproxyExecutable string

	// TerminateSessions controls whether session termination support is
	// enabled.
	TerminateSessions bool

	// SearchDelay is the delay in seconds between searches of the audit log
	// or active file for a specified username. This is an attempt to work
	// around race conditions between EZproxy updating its state file (which
	// has been observed to have a delay of up to several seconds) and this
	// application *reading* the active file. This delay is applied to the
	// initial search and each subsequent retried search for the provided
	// username.
	SearchDelay int

	// SearchRetries is the number of retries allowed for the audit log and
	// active files before the application accepts that "cannot find matching
	// session IDs for specific user" is really the truth of it and not a race
	// condition between this application and the EZproxy application (e.g.,
	// EZproxy accepts a login, but delays writing the state information for
	// about 2 seconds to keep from hammering the storage device).
	SearchRetries int
}

// Branding is responsible for emitting application name, version and origin
func Branding() {
	fmt.Fprintf(flag.CommandLine.Output(), "\n%s %s\n%s\n\n", myAppName, config.Version, myAppURL)
}

// flagsUsage displays branding information and general usage details
func flagsUsage() func() {

	return func() {

		myBinaryName := filepath.Base(os.Args[0])

		Branding()

		fmt.Fprintf(flag.CommandLine.Output(), "Usage of \"%s\":\n",
			myBinaryName,
		)
		flag.PrintDefaults()

		fmt.Fprintf(flag.CommandLine.Output(), "\n")

	}
}

func main() {

	// logging controls for this application
	log.SetLevel(log.InfoLevel)
	log.SetHandler(cli.New(os.Stdout))

	// logging controls for imported ezproxy package and subpackages
	// ezproxy.EnableLogging()
	// ezproxy.DisableLogging()

	config := AppConfig{}

	flag.StringVar(&config.ActiveFilePath, "active-file-path", "", "The fully-qualified path to the EZproxy active users/state file")
	flag.StringVar(&config.ActiveFilePath, "auf", "", "The fully-qualified path to the EZproxy active users/state file")

	flag.StringVar(&config.EZproxyExecutable, "executable", "/opt/ezprozy/ezproxy", "The fully-qualified path to the EZproxy application/binary file")
	flag.StringVar(&config.EZproxyExecutable, "exe", "/opt/ezprozy/ezproxy", "The fully-qualified path to the EZproxy application/binary file")

	flag.StringVar(&config.Username, "username", "", "The name of the username to use when searching for active sessions")

	flag.BoolVar(&config.TerminateSessions, "terminate", false, "Whether active sessions for specified user should be terminated")
	flag.BoolVar(&config.TerminateSessions, "kill", false, "Whether active sessions for specified user should be terminated")

	flag.IntVar(&config.SearchDelay, "search-delay", 1, "The delay in seconds between search attempts.")
	flag.IntVar(&config.SearchDelay, "sd", 1, "The delay in seconds between search attempts.")
	flag.IntVar(&config.SearchDelay, "delay", 1, "The delay in seconds between search attempts.")

	flag.IntVar(&config.SearchRetries, "search-retries", 10, "The number of additional retries allowed when receiving zero search results")
	flag.IntVar(&config.SearchRetries, "sr", 10, "The number of additional retries allowed when receiving zero search results")
	flag.IntVar(&config.SearchRetries, "retries", 10, "The number of additional retries allowed when receiving zero search results")

	flag.Usage = flagsUsage()
	flag.Parse()

	handleError := func(err error) {
		if err != nil {
			log.Error(err.Error())

			// fine for this standalone tool, but *not* how we can implement this
			// within brick itself
			os.Exit(1)
		}
	}

	// flag validation
	if flagsErr := validate(config); flagsErr != nil {
		handleError(flagsErr)
	}

	reader, err := activefile.NewReader(config.Username, config.ActiveFilePath)
	handleError(err)

	// Adjust stubbornness of newly created reader (overridding library/package
	// default values with our own)
	handleError(reader.SetSearchDelay(config.SearchDelay))
	handleError(reader.SetSearchRetries(config.SearchRetries))

	log.Infof("Searching %q for %q", config.ActiveFilePath, config.Username)
	activeSessions, err := reader.MatchingUserSessions()
	handleError(err)

	if len(activeSessions) == 0 {
		log.Infof("No active sessions found for username %q", config.Username)
		return
	}

	fmt.Printf(
		"\nSessions (%d) active for username %q:\n\n",
		len(activeSessions),
		config.Username,
	)

	for _, session := range activeSessions {
		fmt.Printf(
			"Username: %q, SessionID: %q, IP Address: %q\n",
			session.Username,
			session.SessionID,
			session.IPAddress,
		)
	}

	if config.TerminateSessions {

		results := activeSessions.Terminate(config.EZproxyExecutable)

		writer := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', tabwriter.AlignRight)

		fmt.Fprintf(writer, "ID\tExitCode\tStdOut\tStdErr\tErrorMsg\n")

		// Separator row
		// TODO: I'm sure this can be handled better
		fmt.Fprintln(writer, "---\t---\t---\t---\t---\t")

		var termSuccess int

		for _, result := range results {
			if result.ExitCode == ezproxy.KillSubCmdExitCodeSessionTerminated {
				termSuccess++
			}
		}

		fmt.Printf("\nSessions (%d) terminated:\n\n", termSuccess)

		// check the results, report any issues
		for _, result := range results {

			// guard against (nil) lack of error in results slice entry
			errStr := ""
			if result.Error != nil {
				errStr = result.Error.Error()
			}

			fmt.Fprintf(writer, "%s\t%d\t%s\t%s\t%s\t\n",
				result.SessionID,
				result.ExitCode,
				result.StdOut,
				result.StdErr,
				errStr,
			)
		}
		fmt.Fprintln(writer)

		if err := writer.Flush(); err != nil {
			log.Errorf("error flushing tabwriter: %w", err.Error())
		}

	}
}
