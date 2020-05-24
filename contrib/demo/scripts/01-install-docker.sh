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

# Purpose: Install Docker for use in demo environment

if [[ "$UID" -eq 0 ]]; then
  echoerr "Run this script without sudo or as root, sudo will be called as needed."
  exit 1
fi

# Refresh package lists
sudo apt-get update

# Install Docker
# https://docs.docker.com/engine/install/ubuntu/
sudo apt-get install -y \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg-agent \
    software-properties-common
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"
sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io

# Explicitly enable and start Docker (likely already done by upstream packages)
sudo systemctl enable docker
sudo systemctl start docker

# Requires user interaction
#sudo systemctl status docker

# Verify that Docker works
sudo docker run hello-world

