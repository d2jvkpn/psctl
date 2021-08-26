#! /usr/bin/env bash
set -eu -o pipefail

wd=$(pwd)

go build psctl.go


####
./psctl instance examples/ArchVizExplorer427.toml ping

./psctl instance examples/ArchVizExplorer427.toml view

./psctl instance examples/ArchVizExplorer427.toml run

./psctl instance examples/ArchVizExplorer427.toml sync

./psctl instance examples/ArchVizExplorer427.toml syncLog

./psctl instance examples/ArchVizExplorer427.toml syncStatus

./psctl instance examples/ArchVizExplorer427.toml kill


####


./psctl instance examples/ArchVizExplorer426.yaml view

./psctl instance examples/ArchVizExplorer426.yaml run

./psctl instance examples/ArchVizExplorer426.yaml sync

./psctl instance examples/ArchVizExplorer426.yaml syncLog

./psctl instance examples/ArchVizExplorer426.yaml syncStatus

./psctl instance examples/ArchVizExplorer426.yaml kill
