#! /usr/bin/env bash
set -eu -o pipefail

wd=$(pwd)


branch=$(git rev-parse --abbrev-ref HEAD)

git add -A # include Untracked files for git diff
test -z "$(git diff HEAD)" || { echo "You have uncommitted changes!"; exit 1; }

test -z "$(git diff origin/$branch..HEAD --name-status)" ||
  { echo "You have unpushed commits!"; exit 1; }


go build -o psctl -ldflags="                 \
  -X psctl/cmd.BuildBranch="${branch}"       \
  -X psctl/cmd.BuildTime=$(date +'%FT%T%z')" \
  psctl.go

./psctl version
