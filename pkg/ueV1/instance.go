package ueV1

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"psctl/pkg/misc"

	"github.com/BurntSushi/toml"
	yaml "gopkg.in/yaml.v2"
)

var (
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
	Host    string `toml:"host" yaml:"host" json:"host"`          // ansbile hostname or host group name
	Project string `toml:"project" yaml:"project" json:"project"` // project directory in windows host(s)
	Program string `toml:"program" yaml:"program" json:"program"` // exe program filename without extension
	SwsIp   string `toml:"swsIp" yaml:"swsIp" json:"swsIp"`       // singaling and web server ip
	SwsPort string `toml:"swsPort" yaml:"swsPort" json:"swsPort"` // sws streamer port
}

type Instance struct {
	Id         int64    `json:"id"`
	Name       string   `json:"name"`
	Root       string   `json:"root"`
	Command    []string `json:"command"`
	commandMd5 string

	InstanceBase
}

func (inst InstanceBase) String() string {
	return fmt.Sprintf(
		"host=%q, project=%q, program=%q, swsIp=%q, swsPort=%q",
		inst.Host, inst.Project, inst.Program, inst.SwsIp, inst.SwsPort,
	)
}

func InstanceFromFile(fp string, ids ...int64) (inst Instance, err error) {
	var (
		bts  []byte
		base InstanceBase
	)

	if bts, err = ioutil.ReadFile(fp); err != nil {
		return inst, err
	}

	switch {
	case strings.HasSuffix(fp, ".json"):
		err = json.Unmarshal(bts, &base)
	case strings.HasSuffix(fp, ".toml"):
		_, err = toml.DecodeFile(fp, &base)
	case strings.HasSuffix(fp, ".yaml") || strings.HasSuffix(fp, ".yml"):
		err = yaml.Unmarshal(bts, &base)
	}

	if err != nil {
		return
	}

	if len(ids) > 0 {
		return base.Inst(ids[0]), nil
	} else {
		return base.Inst(0), nil
	}
}

func (base InstanceBase) commandline() []string {
	return []string{base.Program + ".exe", "-AudioMixer", "-PixelStreamingIP=" + base.SwsIp,
		"-PixelStreamingPort=" + base.SwsPort, "-RenderOffScreen",
	}
}

func (base InstanceBase) Inst(id int64) (inst Instance) {
	cmd := base.commandline()
	inst = Instance{
		Id:           id,
		Name:         inst.Project + " on " + inst.Host,
		Command:      cmd,
		commandMd5:   misc.CmdMd5(cmd),
		InstanceBase: base,
	}
	inst.Root, _ = os.Getwd()
	return inst
}

func (inst *Instance) WorkPath() string {
	strs := strings.Split(inst.Project, "\\")

	return filepath.Join(RootWorkingDir, inst.Host, strs[len(strs)-1], inst.commandMd5)
}

func (inst *Instance) RunCmd(name string, arg ...string) (err error) {
	var buf bytes.Buffer

	cmd := exec.Command(name, arg...)
	cmd.Dir = inst.WorkPath()
	cmd.Env = append(cmd.Env, "ANSIBLE_LOG_PATH="+filepath.Join("logs", "ansible.log"))
	cmd.Stdout, cmd.Stderr = &buf, &buf

	if err = cmd.Run(); err != nil {
		if out := buf.String(); len(out) > 0 {
			return fmt.Errorf("%s\nRunCmd error: %w", out, err)
		}
		return fmt.Errorf("RunCmd error: %w", err)
	}
	return nil
}

func (inst *Instance) writeStatus(status string) (err error) {
	var file *os.File

	file, err = os.OpenFile(
		filepath.Join(inst.WorkPath(), "status.log"),
		os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644,
	)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%s\t%s\n", time.Now().Format(time.RFC3339), status))

	return err
}

func (inst *Instance) end() (err error) {
	var (
		bn, wd, suffix string
		now            time.Time
	)

	if err = inst.Sync(); err != nil {
		return err
	}

	now = time.Now()
	wd, suffix = inst.WorkPath(), now.Format(".backup_2006-01-03T15-04-05.log")

	bn = filepath.Join(wd, "logs", inst.commandMd5)
	os.Rename(bn+".log", bn+suffix)

	bn = filepath.Join(wd, "logs", inst.Program)
	os.Rename(bn+".log", bn+suffix)

	bn = filepath.Join(wd, "logs", "ansible")
	os.Rename(bn+".log", bn+suffix)

	return err
}

func (inst *Instance) clear() (err error) {
	return os.RemoveAll(inst.WorkPath())
}
