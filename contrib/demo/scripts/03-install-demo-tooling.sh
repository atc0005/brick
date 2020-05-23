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

# Purpose: Install applications specific to demo environment

if [[ "$UID" -eq 0 ]]; then
  echoerr "Run this script without sudo or as root, sudo will be called as needed."
  exit 1
fi

# Refresh package lists
sudo apt-get update

# Install demo tooling
sudo apt-get install -y \
    ccze \
    guake \
    guake-indicator

# Install Visual Studio Code
cd /tmp
curl -L https://go.microsoft.com/fwlink/?LinkID=760868 > vscode-installer.deb
sudo apt-get install -y ./vscode-installer.deb
