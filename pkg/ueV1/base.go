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

type InstanceBase struct {
	Host    string `toml:"host" yaml:"host" json:"host,omitempty"`          // ansbile hostname or host group name
	Project string `toml:"project" yaml:"project" json:"project,omitempty"` // project directory in windows host(s)
	Program string `toml:"program" yaml:"program" json:"program,omitempty"` // exe program filename without extension
	SwsIp   string `toml:"swsIp" yaml:"swsIp" json:"swsIp,omitempty"`       // singaling and web server ip
	SwsPort string `toml:"swsPort" yaml:"swsPort" json:"swsPort,omitempty"` // sws streamer port
}

type Instance struct {
	Id         int64    `json:"id,omitempty"`
	Name       string   `json:"name,omitempty"`
	Root       string   `json:"root,omitempty"`
	Command    []string `json:"command,omitempty"`
	commandMd5 string
	Debug      bool `json:"-"`

	InstanceBase
}
