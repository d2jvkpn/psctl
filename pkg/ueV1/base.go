package ueV1

import (
	_ "embed"
)

const (
	RootWorkingDir = "data"
)

var (
	//go:embed playbook.yaml
	playbook []byte
)
