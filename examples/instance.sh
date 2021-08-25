#! /usr/bin/env bash
set -eu -o pipefail

wd=$(pwd)


####
go run psctl.go instance examples/ArchVizExplorer427.toml view

go run psctl.go instance examples/ArchVizExplorer427.toml run

go run psctl.go instance examples/ArchVizExplorer427.toml sync

go run psctl.go instance examples/ArchVizExplorer427.toml syncLog

go run psctl.go instance examples/ArchVizExplorer427.toml syncStatus

go run psctl.go instance examples/ArchVizExplorer427.toml kill


####
go run psctl.go instance examples/ArchVizExplorer426.yaml view

go run psctl.go instance examples/ArchVizExplorer426.yaml run

go run psctl.go instance examples/ArchVizExplorer426.yaml sync

go run psctl.go instance examples/ArchVizExplorer426.yaml syncLog

go run psctl.go instance examples/ArchVizExplorer426.yaml syncStatus

go run psctl.go instance examples/ArchVizExplorer426.yaml kill
