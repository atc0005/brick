# Copyright 2023 Adam Chalkley
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
# Dependabot-maintained Dockerfile used by Makefile-driven release builds.
#

# https://github.com/atc0005/go-ci/releases
# https://github.com/atc0005/go-ci/pkgs/container/go-ci
FROM ghcr.io/atc0005/go-ci:go-ci-oldstable-build-v0.11.3

# Setup isolated build environment with a full copy of the Git repo contents
# MINUS any file or path listed in the .dockerignore file at the root of this
# repo.
RUN useradd --create-home --shell /bin/bash --user-group builduser

# Prevent Git from complaining when it encounters Git-tracked directories that
# are owned by someone other than the current user. We set this at the
# "system" level so that the setting is not specific to any one user account.
#
# https://stackoverflow.com/questions/71849415/i-cannot-add-the-parent-directory-to-safe-directory-in-git
# https://git-scm.com/docs/git-config/2.35.2#Documentation/git-config.txt-safedirectory
# https://github.com/git/git/commit/8959555cee7ec045958f9b6dd62e541affb7e7d9
RUN git config --system --add safe.directory '*'

# We skip setting a specific container user. This allows generating a
# container with a bind-mounted path using Podman and explicitly specifying
# `--user builduser:builduser` using Docker.
#
#USER builduser

WORKDIR /builds
COPY --chown=builduser:builduser . /builds
