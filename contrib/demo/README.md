<!-- omit in toc -->
# Brick: Demo deployment scripts

<!-- omit in toc -->
## Table of contents

- [Purpose](#purpose)
- [Target environment](#target-environment)
- [List](#list)

## Purpose

A set of installation and deployment scripts to handle setting up a **demo**
environment. The plan is to replace these scripts with an Ansible playbook in
an effort to make the provided files easier to deploy.

Assumptions for these scripts:

- run **within a test environment**
- run with the contrib/demo directory as the current working directory

## Target environment

The scripts assume that they're running in one of these two tested
environments:

- Ubuntu 16.04 LTS
- Ubuntu 18.04 LTS

## List

| Step | Script name                  | Purpose                                                               |
| ---- | ---------------------------- | --------------------------------------------------------------------- |
| 1    | `01-install-docker.sh`       | Install Docker                                                        |
| 2    | `02-install-dependencies.sh` | Install core application dependencies                                 |
| 3    | `03-install-demo-tooling.sh` | Install apps used in demo                                             |
| 4    | `04-setup-env.sh`            | Configure dependencies, test environment                              |
| 5    | `05-deploy.sh`               | Build and deploy `brick` application                                  |
| 6    | `06-reset-env.sh`            | Reset environment for demo (truncate log files, restart daemons, etc) |
