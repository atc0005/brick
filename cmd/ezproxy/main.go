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
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"regexp"
	"strings"

	"github.com/apex/log"

	"github.com/atc0005/go-ezproxy"
)

// Results is used to group together potential exit codes and output in a
// slice to be chosen pseudo-randomly.
type Results struct {
	ExitOutput string
	ExitCode   int
}

func main() {

	log.SetLevel(log.InfoLevel)

	log.Debug("Starting EZproxy mock binary")

	// called as: ezproxy kill SESSION_ID_HERE

	switch len(os.Args) {

	// ezproxy (without subcommand or arguments)
	case 1:
		fmt.Println("missing subcommand")
		os.Exit(1)

	// ezproxy kill (without specifying session ID)
	case 2:

		// subcommand is present, but not session ID
		if strings.EqualFold(os.Args[1], ezproxy.SubCmdNameSessionTerminate) {
			fmt.Println(ezproxy.KillSubCmdExitTextSessionNotSpecified)
			os.Exit(ezproxy.KillSubCmdExitCodeSessionNotSpecified)
		}

		// caller didn't even get the subcommand right, complain about that
		// and don't even mention that session id wasn't provided
		fmt.Println("invalid subcommand")
		os.Exit(1)

	// ezproxy kill SESSION_ID_HERE (valid number of arguments)
	case 3:

		// verify that the correct subcommand was used
		if !strings.EqualFold(os.Args[1], ezproxy.SubCmdNameSessionTerminate) {
			fmt.Println("invalid subcommand")
			os.Exit(1)
		}

		// ensure that the provided value is of the length previously observed
		// in the active users file
		//
		// TODO: Try to find an official OCLC reference for the session ID
		// restrictions (characters, length, etc)
		if len(os.Args[2]) < ezproxy.SessionIDLength {
			fmt.Println("invalid session ID length provided")
			os.Exit(1)
		}

		matchOK, matchErr := regexp.MatchString(
			ezproxy.SessionIDRegex,
			os.Args[2],
		)
		if !matchOK {
			if matchErr != nil {
				log.Error(matchErr.Error())
			}
			fmt.Println("invalid session ID pattern provided")
			os.Exit(1)
		}

		resultCodes := []Results{
			{
				ExitOutput: fmt.Sprintf(
					ezproxy.KillSubCmdExitTextTemplateSessionTerminated,
					os.Args[2],
				),
				ExitCode: ezproxy.KillSubCmdExitCodeSessionTerminated,
			},
			{
				ExitOutput: fmt.Sprintf(
					ezproxy.KillSubCmdExitTextTemplateSessionDoesNotExist,
					os.Args[2],
				),
				ExitCode: ezproxy.KillSubCmdExitCodeSessionDoesNotExist,
			},
			{
				ExitOutput: "Did you know that fish sticks pizza is really a thing?",
				ExitCode:   2,
			},
			{
				ExitOutput: "Did you know that macaroni pizza is really a thing?",
				ExitCode:   2,
			},
		}

		// randomly return one of the result codes from the list above
		maxRandomNumber := big.NewInt(int64(len(resultCodes)))
		nBig, rngErr := rand.Int(rand.Reader, maxRandomNumber)
		if rngErr != nil {
			fmt.Printf("unable to generate random number: %v", rngErr)
			os.Exit(1)
		}
		randomResultIdx := nBig.Int64()
		fmt.Println(resultCodes[randomResultIdx].ExitOutput)
		os.Exit(resultCodes[randomResultIdx].ExitCode)

	default:

		fmt.Println("too many options provided, I can't decide!")
		os.Exit(1)

	}

}
