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

package caller

import (
	"fmt"
	"runtime"
)

// GetFuncFileLineInfo is a wrapper around the runtime.Caller() function. This
// function returns the calling function name, filename and line number to
// help with debugging efforts.
func GetFuncFileLineInfo() string {

	if pc, file, line, ok := runtime.Caller(1); ok {
		return fmt.Sprintf(
			"func %s called (from %q, line %d): ",
			runtime.FuncForPC(pc).Name(),
			file,
			line,
		)
	}

	return "error: unable to recover caller origin via runtime.Caller()"
}

// GetFuncName is a wrapper around the runtime.Caller() function. This
// function returns the calling function name and discards other return
// values.
func GetFuncName() string {

	if pc, _, _, ok := runtime.Caller(1); ok {
		return runtime.FuncForPC(pc).Name()
	}

	return "error: unable to recover caller origin via runtime.Caller()"
}

// GetParentFuncFileLineInfo is a wrapper around the runtime.Caller()
// function. This function returns the parent calling function name, filename
// and line number to help with debugging efforts.
func GetParentFuncFileLineInfo() string {

	if pc, file, line, ok := runtime.Caller(2); ok {
		return fmt.Sprintf(
			"func %s called (from %q, line %d): ",
			runtime.FuncForPC(pc).Name(),
			file,
			line,
		)
	}

	return "error: unable to recover caller parent origin via runtime.Caller()"
}

// GetParentFuncName is a wrapper around the runtime.Caller() function. This
// function returns the parent calling function name and discards other return
// values.
func GetParentFuncName() string {

	if pc, _, _, ok := runtime.Caller(2); ok {
		return runtime.FuncForPC(pc).Name()
	}

	return "error: unable to recover caller parent origin via runtime.Caller()"
}
