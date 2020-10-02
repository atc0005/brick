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

package config

import (
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/apex/log/handlers/discard"
	"github.com/apex/log/handlers/json"
	"github.com/apex/log/handlers/logfmt"
	"github.com/apex/log/handlers/text"
)

// configureLogging is a wrapper function to enable setting requested logging
// settings.
func (c Config) configureLogging() {

	switch c.LogLevel() {
	case LogLevelFatal:
		log.SetLevel(log.FatalLevel)
	case LogLevelError:
		log.SetLevel(log.ErrorLevel)
	case LogLevelWarn:
		log.SetLevel(log.WarnLevel)
	case LogLevelInfo:
		log.SetLevel(log.InfoLevel)
	case LogLevelDebug:
		log.SetLevel(log.DebugLevel)
	}

	// Apply user-specified logging output target
	var outputTarget *os.File
	switch c.LogOutput() {
	case LogOutputStdout:
		outputTarget = os.Stdout
	case LogOutputStderr:
		outputTarget = os.Stderr
	default:
		outputTarget = os.Stdout
	}

	switch c.LogFormat() {
	case LogFormatText:
		log.SetHandler(text.New(outputTarget))
	case LogFormatCLI:
		log.SetHandler(cli.New(outputTarget))
	case LogFormatLogFmt:
		log.SetHandler(logfmt.New(outputTarget))
	case LogFormatJSON:
		log.SetHandler(json.New(outputTarget))
	case LogFormatDiscard:
		log.SetHandler(discard.New())
	}

}
