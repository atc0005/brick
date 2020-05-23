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

# Purpose: Install core dependencies of this application.

if [[ "$UID" -eq 0 ]]; then
  echoerr "Run this script without sudo or as root, sudo will be called as needed."
  exit 1
fi

# Install dependencies for this app via APT repos
sudo apt-get update

# https://coderwall.com/p/lryimq/postfix-silent-install-on-ubuntu
echo "postfix postfix/mailname string example.com" | sudo debconf-set-selections
echo "postfix postfix/main_mailer_type string 'Internet Site'" | sudo debconf-set-selections

sudo apt-get install -y \
    fail2ban \
    geoip-bin \
    geoip-database \
    geoip-database-extra \
    mailutils \
    postfix \
    sqlite3

# Install Go toolchain
if [ ! -f /usr/local/go/bin/go ]
then
    # install Go toolchain
    echo "Downloading and installing Go toolchain"
    cd /tmp
    rm -f go1.14.3.linux-amd64.tar.gz
    curl -L -O https://dl.google.com/go/go1.14.3.linux-amd64.tar.gz
    sudo tar zxf go1.14.3.linux-amd64.tar.gz -C /usr/local/
    /usr/local/go/bin/go version
else
    echo "Go toolchain already present"
    /usr/local/go/bin/go version
fi

# Extend PATH to reference new Go installation
if ! grep -q "/usr/local/go/bin" $HOME/.profile
then
    echo "Go installation path is not in PATH, adding it."
    cat >> $HOME/.profile << HEREDOC

PATH="$PATH:/usr/local/go/bin:$HOME/go/bin"

HEREDOC

echo "Make sure to log out and back in to load new PATH settings"
fi
