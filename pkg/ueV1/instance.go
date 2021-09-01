package ueV1

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"psctl/pkg/misc"

	"github.com/BurntSushi/toml"
	yaml "gopkg.in/yaml.v2"
)

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

func (inst *Instance) NewPlaybook(override bool) (err error) {
	dir := inst.WorkPath()
	target := filepath.Join(dir, "logs")
	yes := false

	if yes, err = misc.DirExists(target); err != nil {
		return err
	} else if yes && !override {
		return nil
	}

	if err = os.MkdirAll(target, 0755); err != nil {
		return fmt.Errorf("mkdir %s: %w", dir, err)
	}

	err = ioutil.WriteFile(filepath.Join(dir, "playbook.yaml"), playbook, 0644)
	if err != nil {
		return err
	}

	data := struct {
		*Instance
		CommandMd5 string
		Cmd        string
	}{
		Instance:   inst,
		CommandMd5: inst.commandMd5,
		Cmd:        strings.Join(inst.Command, " "),
	}

	buf := bytes.NewBuffer([]byte{})
	if err = varsYamlTmpl.Execute(buf, data); err != nil {
		return fmt.Errorf("varsYamlTmpl.Execute: %w", err)
	}

	if err = ioutil.WriteFile(filepath.Join(dir, "vars.yaml"), buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("write vars.yaml: %w", err)
	}

	if err = inst.Playbook("--tags", "prepare"); err != nil {
		return fmt.Errorf("ansile-playbook prepare: %w", err)
	}

	return nil
}

func (inst *Instance) Exists() (yes bool, err error) {
	var target string

	target = filepath.Join(inst.WorkPath(), "playbook.yaml")
	if yes, err = misc.FileExists(target); err != nil {
		return false, err
	}

	target = filepath.Join(inst.WorkPath(), "vars.yaml")
	if yes, err = misc.FileExists(target); err != nil {
		return false, err
	}

	return yes, nil
}

func (inst *Instance) View() (string, error) {
	var (
		yes    bool
		bts    []byte
		target string
		err    error
	)

	target = filepath.Join(inst.WorkPath(), "status.log")
	if yes, err = misc.FileExists(target); err != nil {
		return "", err
	}
	if !yes {
		return "", nil
	}

	if bts, err = ioutil.ReadFile(target); err != nil {
		return "", err
	}

	return string(bts), err
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
