#! /usr/bin/env bash
set -eu -o pipefail

wd=$(pwd)

go build psctl.go


####
./psctl load examples/ArchVizExplorer427.toml view

./psctl load examples/ArchVizExplorer427.toml new

./psctl load examples/ArchVizExplorer427.toml start

./psctl load examples/ArchVizExplorer427.toml kill

./psctl load examples/ArchVizExplorer427.toml sync

./psctl load examples/ArchVizExplorer427.toml syncLog

./psctl load examples/ArchVizExplorer427.toml syncStatus

./psctl load examples/ArchVizExplorer427.toml restart


####
./psctl load examples/ArchVizExplorer426.yaml view

./psctl load examples/ArchVizExplorer426.yaml start

./psctl load examples/ArchVizExplorer426.yaml sync

./psctl load examples/ArchVizExplorer426.yaml syncLog

./psctl load examples/ArchVizExplorer426.yaml syncStatus

./psctl load examples/ArchVizExplorer426.yaml kill
