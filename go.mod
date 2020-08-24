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

module github.com/atc0005/brick

// replace github.com/atc0005/go-ezproxy => ../go-ezproxy

go 1.14

require (
	github.com/Showmax/go-fqdn v0.0.0-20180501083314-6f60894d629f
	github.com/alexflint/go-arg v1.3.0
	github.com/apex/log v1.9.0
	github.com/atc0005/go-ezproxy v0.1.6
	// temporarily use our fork; waiting on changes to be accepted upstream
	github.com/atc0005/go-teams-notify v1.3.1-0.20200419155834-55cca556e726
	github.com/atc0005/send2teams v0.4.6
	github.com/pelletier/go-toml v1.8.0
)
