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

echo "This script resets output files."
script_names="$(ls *setup-env* *deploy*)"
echo "Run scripts ${script_names} to reset other files also."

# FIXME: Is this even needed?
#
# sudo touch /var/log/brick/users.brick-disabled.txt
# sudo touch /var/log/brick/users.brick-reported.log
# sudo touch /var/log/brick/syslog.log
# sudo touch /var/log/fail2ban.log
# sudo touch /var/log/mail.log

# Stop services
sudo systemctl stop fail2ban
sudo systemctl stop postfix
sudo systemctl stop rsyslog
sudo systemctl stop brick

# Truncate log files
sudo truncate -s 0 /var/cache/brick/users.brick-disabled.txt
sudo truncate -s 0 /var/log/brick/users.brick-reported.log
sudo truncate -s 0 /var/log/brick/syslog.log
sudo truncate -s 0 /var/log/fail2ban.log
sudo truncate -s 0 /var/log/mail.log

# Fix permissions
sudo chown -Rv brick:syslog /var/log/brick
sudo chmod -v 664 /var/log/brick/syslog.log
sudo chmod -v 644 /var/log/brick/users.brick-reported.log
sudo chown -Rv brick:brick /var/cache/brick
sudo chmod -v 644 /var/cache/brick/users.brick-disabled.txt

# Clear any stuck emails
sudo postsuper -d ALL

# Reset Docker container
sudo docker container stop maildev
sudo docker run --rm --detach --name maildev -p 1080:80 -p 1025:25 maildev/maildev

# Reset fail2ban state
# https://unix.stackexchange.com/questions/286119
sudo systemctl stop fail2ban
for lin in {200..1}; do
   sudo iptables -D f2b-brick $lin 2>&1 | grep -v 'iptables: No'
done
sudo rm -vf /var/lib/fail2ban/fail2ban.sqlite3

# Restart affected services so that they'll start over fresh
sudo systemctl restart rsyslog
sudo systemctl restart postfix
sudo systemctl restart fail2ban
sudo systemctl restart brick
