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

Package activefile is intended for the processing of EZproxy active users and
hosts files.

OVERVIEW

There is only ever one Active Users and Hosts "state" file at a time. While
Host entries are (from what can be observed) consolidated on one line, session
entries are composed of multiple lines in a very specific order, each with
space-separated fields. These order-specific lines and fields are joined in
order to reconstruct a User Session that reflects active user sessions within
EZproxy.

KNOWN TYPES

Known entry types include (but may not be limited to):

Host (H)
Group (g)
Session (S)
Username or Login (L)

Currently, only the last two types (S, L) are relevant to our purposes.

UNKNOWN TYPES

These types have been observed, but not researched sufficiently to identify
their purpose (Pull Requests for this are welcome!):

P
M
s (lowercase letter)

LINE ORDERING

For our purposes, we match lines that start with a capital letter S and pair
it with the first line following it that begins with a capital letter L. We
skip over any line that begins with a lowercase letter s; we do not use the
value provided by this line.

When we match a line beginning with a capital S, these are the only supported
line orderings:

S
s
L

and:

S
L

FIELD NUMBERS

The line for for Logins (L) is composed of 2 fields:

01) Leading capital letter L
02) Username

The line for Sessions (S) is composed of 11 fields:

01) Leading capital letter S
02) Session ID
03) unknown, appears to be a UNIX timestamp
04) unknown, appears to be two UNIX timestamps separated by a literal dot
05) unknown integer; number 1 was common
06) EZproxy "MaxLifetime" or User Session timeout value
07) IP Address
08) unknown, 0 is recorded
09) unknown, 0 is recorded
10) unknown, 0 is recorded
11) unknown, asterisk is recorded

RACE CONDITION

NOTE: EZproxy does not immediately update the Active Users and Hosts "state"
file with state changes; when a user account logs in/out, there is a race
condition between when that information is updated and when a Reader created
from this package attempts to read the current state and reconstruct User
Sessions. In an effort to workaround this race condition, this package
attempts to retry session read attempts a limited number of times by default
before giving up. This retry behavior (including a delay between retry
attempts) can be modified by the caller as needed.

*/
package activefile
