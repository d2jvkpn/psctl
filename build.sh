#! /usr/bin/env bash
set -eu -o pipefail

wd=$(pwd)

go build -o psctl -ldflags="                 \
  -X psctl/cmd.BuildTime=$(date +'%FT%T%z')" \
  psctl.go

./psctl version
