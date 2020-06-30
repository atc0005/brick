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

/*

Package ezproxy is intended for the processing of EZproxy related files and
sessions.

PROJECT HOME

See our GitHub repo (https://github.com/atc0005/go-ezproxy) for the latest
code, to file an issue or submit improvements for review and potential
inclusion into the project.

PURPOSE

Process EZproxy related files and sessions.

FEATURES

• generate a list of audit records for session-related events for all usernames or just for a specific username

• generate a list of active sessions using the audit log using entires without a corresponding logout event type

• generate a list of active sessions using the active file for all usernames or just for a specific username

• terminate single user session or bulk user sessions

OVERVIEW

Ultimately, this package was written in order to support retrieving session
information for a specific username so that the session can be terminated.
Because of this the majority of the functionality is specific to user
sessions.

General workflow:

1. Import this package
2. Import one or more of the subpackages
3. Create a new reader for the file type you need to work with
4. Using the new reader, generate a UserSessions collection
5. Use the Terminate method to terminate user sessions

If using the ezproxy/auditlog package, you can also generate a SessionEntries
collection representing all SessionEntry values from a specified audit log
file or just the values applicable to a specifc user.

FUTURE

This package currently provides functionality for working with an active user
or audit log file, but not for EZproxy traffic log files. Having minimal
support for traffic log files could provide a way to die activity for a
specific user account to specifc resources accessed by that user account. This
could prove invaluable where automation is used to automatically terminate
user sessions; after account termination, a report could be generated for the
incident listing at a high-level the providers accessed and general statistics
associated with the access (e.g., PDF downloads, total bandwidth, etc.).

*/
package ezproxy
