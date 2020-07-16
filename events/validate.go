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

package events

import (
	"errors"
	"fmt"
)

// ValidatePayload is used to perform very basic validation on all expected
// fields for the received payload.
func ValidatePayload(payloadV2 SplunkAlertPayloadV2) error {

	// The `SplunkAlertPayloadV2` type is currently composed of string fields
	// only, so unless we apply regex patterns we're mostly concerned with
	// just making sure that fields which we rely on have something in them.

	validationFailedErr := errors.New("payload validation failed; required field value missing")

	if payloadV2.ResultsLink == "" {
		return fmt.Errorf("%w: ResultsLink field empty", validationFailedErr)
	}

	if payloadV2.Result.SourceIP == "" {
		return fmt.Errorf("%w: Result.SourceIP field empty", validationFailedErr)
	}

	if payloadV2.Result.Username == "" {
		return fmt.Errorf("%w: Result.Username field empty", validationFailedErr)
	}

	if payloadV2.Result.ResourceURL == "" {
		return fmt.Errorf("%w: Result.ResourceURL field empty", validationFailedErr)
	}

	if payloadV2.Result.HTTPStatusCode == "" {
		return fmt.Errorf("%w: Result.HTTPStatusCode field empty", validationFailedErr)
	}

	if payloadV2.Result.UserAgent == "" {
		return fmt.Errorf("%w: Result.UserAgent field empty", validationFailedErr)
	}

	if payloadV2.Result.TagEventtype == "" {
		return fmt.Errorf("%w: Result.TagEventtype field empty", validationFailedErr)
	}

	if payloadV2.Result.Time == "" {
		return fmt.Errorf("%w: Result.Time field empty", validationFailedErr)
	}

	if payloadV2.Result.EzproxyTime == "" {
		return fmt.Errorf("%w: Result.EzproxyTime field empty", validationFailedErr)
	}

	if payloadV2.Result.DateHour == "" {
		return fmt.Errorf("%w: Result.DateHour field empty", validationFailedErr)
	}

	if payloadV2.Result.DateSecond == "" {
		return fmt.Errorf("%w: Result.DateSecond field empty", validationFailedErr)
	}

	if payloadV2.Result.DateMinute == "" {
		return fmt.Errorf("%w: Result.DateMinute field empty", validationFailedErr)
	}

	if payloadV2.Result.DateMday == "" {
		return fmt.Errorf("%w: Result.DateMday field empty", validationFailedErr)
	}

	if payloadV2.Result.DateYear == "" {
		return fmt.Errorf("%w: Result.DateYear field empty", validationFailedErr)
	}

	if payloadV2.Result.DateWday == "" {
		return fmt.Errorf("%w: Result.DateWday field empty", validationFailedErr)
	}

	if payloadV2.Result.DateMonth == "" {
		return fmt.Errorf("%w: Result.DateMonth field empty", validationFailedErr)
	}

	if payloadV2.Result.DateZone == "" {
		return fmt.Errorf("%w: Result.DateZone field empty", validationFailedErr)
	}

	if payloadV2.Result.Index == "" {
		return fmt.Errorf("%w: Result.Index field empty", validationFailedErr)
	}

	if payloadV2.Result.Sourcetype == "" {
		return fmt.Errorf("%w: Result.Sourcetype field empty", validationFailedErr)
	}

	if payloadV2.Result.Source == "" {
		return fmt.Errorf("%w: Result.Source field empty", validationFailedErr)
	}

	if payloadV2.Result.SplunkServer == "" {
		return fmt.Errorf("%w: Result.SplunkServer field empty", validationFailedErr)
	}

	if payloadV2.Result.Bkt == "" {
		return fmt.Errorf("%w: Result.Bkt field empty", validationFailedErr)
	}

	if payloadV2.Result.Tag == "" {
		return fmt.Errorf("%w: Result.Tag field empty", validationFailedErr)
	}

	if payloadV2.Result.Raw == "" {
		return fmt.Errorf("%w: Result.Raw field empty", validationFailedErr)
	}

	if payloadV2.Sid == "" {
		return fmt.Errorf("%w: Sid field empty", validationFailedErr)
	}

	if payloadV2.Owner == "" {
		return fmt.Errorf("%w: Owner field empty", validationFailedErr)
	}

	if payloadV2.App == "" {
		return fmt.Errorf("%w: App field empty", validationFailedErr)
	}

	if payloadV2.SearchName == "" {
		return fmt.Errorf("%w: SearchName field empty", validationFailedErr)
	}

	return nil

}
