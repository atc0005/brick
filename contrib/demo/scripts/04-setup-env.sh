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

# Purpose: Setup demo environment (create service account, directories,
# configure, ...)

# Create user account
sudo useradd --system --shell /bin/false brick

# Setup test ezproxy directory
sudo mkdir -vp /usr/local/ezproxy/audit
sudo chown -Rv brick:$USER /usr/local/ezproxy
sudo chmod -Rv u=rwX,g=rwX,o= brick:$USER /usr/local/ezproxy

# Setup output directories
sudo mkdir -vp /usr/local/etc/brick
sudo mkdir -vp /var/log/brick
sudo touch /var/log/brick/users.brick-reported.log
sudo chown -Rv brick:syslog /var/log/brick
sudo chmod -Rv g+rwX /var/log/brick
sudo mkdir -vp /var/cache/brick
sudo chown -Rv brick:brick /var/cache/brick

# Setup Docker container
sudo docker container stop maildev
sudo docker run --rm --detach --name maildev -p 1080:80 -p 1025:25 maildev/maildev

# Deploy Postfix configuration files
sudo cp -vf ../../postfix/main.cf /etc/postfix/
sudo cp -vf ../../postfix/header_checks.conf /etc/postfix/
sudo systemctl enable postfix
sudo systemctl restart postfix

# Deploy logrotate config
sudo cp -vf ../../logrotate/brick /etc/logrotate.d/
sudo chmod 644 /etc/logrotate.d/brick

# Deploy fail2ban config
sudo cp -vf ../../fail2ban/action.d/*.local /etc/fail2ban/action.d/
sudo cp -vf ../../fail2ban/filter.d/*.local /etc/fail2ban/filter.d/
sudo cp -vf ../../fail2ban/jail.local /etc/fail2ban/
sudo systemctl enable fail2ban
sudo systemctl restart fail2ban

# Deploy rsyslog config
sudo cp -vf ../../rsyslog/brick.conf /etc/rsyslog.d/00-brick.conf
sudo chmod -v 644 /etc/rsyslog.d/00-brick.conf
sudo systemctl enable rsyslog
sudo systemctl restart rsyslog

# Deploy brick unit file, modify to use our config file
sudo cp -vf ../../systemd/brick.service /etc/systemd/system/brick.service
sudo sed -i -r 's#^ExecStart=/usr/local/sbin/brick#ExecStart=/usr/local/sbin/brick --config-file /usr/local/etc/brick/config.toml#' /etc/systemd/system/brick.service
sudo systemctl daemon-reload

# Deploy app config files
[ ! -f /usr/local/etc/brick/config.toml ] && sudo cp -v ../../brick/config.example.toml /usr/local/etc/brick/config.toml
[ ! -f /usr/local/etc/brick/users.brick-ignored.txt ] && sudo cp -v ../../brick/users.brick-ignored.txt /usr/local/etc/brick/users.brick-ignored.txt
[ ! -f /usr/local/etc/brick/ips.brick-ignored.txt ] && sudo cp -v ../../brick/ips.brick-ignored.txt /usr/local/etc/brick/ips.brick-ignored.txt

sudo chown -v $USER:brick /usr/local/etc/brick/config.toml
sudo chmod -v 640 /usr/local/etc/brick/config.toml

sudo chown -v $USER:brick /usr/local/etc/brick/*.brick-ignored.txt
sudo chmod -v 640 /usr/local/etc/brick/*.brick-ignored.txt
