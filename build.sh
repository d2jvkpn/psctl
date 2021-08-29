#! /usr/bin/env bash
set -eu -o pipefail

wd=$(pwd)


test -z "$(git diff HEAD)" || { echo "You have uncommitted changes!"; exit 1; }

branch=$(git rev-parse --abbrev-ref HEAD)

go build -o psctl -ldflags="                 \
  -X psctl/cmd.BuildBranch="${branch}"       \
  -X psctl/cmd.BuildTime=$(date +'%FT%T%z')" \
  psctl.go

./psctl version
