package ueV1

import (
	_ "embed"
	"fmt"
	"os"
	"text/template"
)

const (
	RootWorkingDir = "data"
)

var (
	//go:embed playbook.yaml
	playbook []byte

	//go:embed vars.tmpl
	varsYaml     string
	varsYamlTmpl *template.Template
)

func init() {
	var err error

	if varsYamlTmpl, err = template.New("varsYaml").Parse(varsYaml); err != nil {
		fmt.Fprintf(os.Stderr, "ue.varsYamlTmpl Parse: %v\n", err)
	}

}
