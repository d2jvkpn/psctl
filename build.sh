#! /usr/bin/env bash
set -eu -o pipefail

_wd=$(pwd)
_path=$(dirname $0)


test -z "$(git status --short)" || { echo "You have uncommitted changes!"; exit 1; }

b1=$(git rev-parse --abbrev-ref HEAD)
test -z "$(git diff origin/$b1..HEAD --name-status)" ||
  { echo "You have unpushed commits!"; exit 1; }


branch=$b1

go build -o psctl -ldflags="                 \
  -X psctl/cmd.BuildBranch="${branch}"       \
  -X psctl/cmd.BuildTime=$(date +'%FT%T%z')" \
  psctl.go

./psctl version
