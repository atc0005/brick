#!/bin/bash

# Copyright 2020 Adam Chalkley
#
# https://github.com/atc0005/brick
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Purpose:
#
# Part of a set of install/deploy scripts to handle setting up a demo
# environment. The plan is to replace this script with an Ansible playbook
# before taking this project public in an effort to make the provided files
# easier to deploy.

# Assumption: This script is run with the contrib/demo directory as the
# current working directory.

if [[ "$UID" -eq 0 ]]; then
  echo "Run this script without sudo or as root, sudo will be called as needed."
  exit 1
fi

# Build applications
go build -mod=vendor ../../../cmd/brick/
go build -mod=vendor ../../../cmd/es/
go build -mod=vendor ../../../cmd/ezproxy/

# Deploy applications
sudo service brick stop
sudo mv -vf brick /usr/local/sbin/brick
sudo mv -vf ezproxy /usr/local/ezproxy/ezproxy
sudo mv -vf es /usr/local/sbin/es

# Set executable bit
sudo chmod -v +x /usr/local/sbin/brick
sudo chmod -v +x /usr/local/ezproxy/ezproxy
sudo chmod -v +x /usr/local/sbin/es

# Start app
sudo systemctl enable brick
sudo systemctl start brick
sudo systemctl status brick
