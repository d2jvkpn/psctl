package ueV1

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"psctl/pkg/misc"
)

func (inst *Instance) NewPlaybook0(override bool) (err error) {
	dir := inst.WorkPath()
	if override {
		err = os.MkdirAll(dir, 0755)
	} else {
		err = os.Mkdir(dir, 0755)
	}
	if err != nil {
		return fmt.Errorf("mkdir %s: %w", dir, err)
	}

	if err = misc.FileCopy("data/playbook.yaml", filepath.Join(dir, "playbook.yaml")); err != nil {
		return fmt.Errorf("copy playbook.yaml: %w", err)
	}

	data := struct {
		*Instance
		Md5Sum string
		Cmd    string
	}{
		Instance: inst,
		Md5Sum:   misc.CmdMd5(inst.commandline()),
		Cmd:      strings.Join(inst.Command, " "),
	}

	buf := bytes.NewBuffer([]byte{})
	if err = varsYamlTmpl.Execute(buf, data); err != nil {
		return fmt.Errorf("varsYamlTmpl.Execute: %w", err)
	}

	if err = ioutil.WriteFile(filepath.Join(dir, "vars.yaml"), buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("write vars.yaml: %w", err)
	}

	return nil
}

func (inst *Instance) NewPlaybook(override bool) (err error) {
	dir := inst.WorkPath()
	if override {
		err = os.MkdirAll(filepath.Join(dir, "logs"), 0755)
	} else {
		err = os.Mkdir(filepath.Join(dir, "logs"), 0755)
	}
	if err != nil {
		return fmt.Errorf("mkdir %s: %w", dir, err)
	}

	err = ioutil.WriteFile(filepath.Join(dir, "playbook.yaml"), playbook, 0644)
	if err != nil {
		return err
	}

	data := struct {
		*Instance
		Md5Sum string
		Cmd    string
	}{
		Instance: inst,
		Md5Sum:   misc.CmdMd5(inst.commandline()),
		Cmd:      strings.Join(inst.Command, " "),
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

func (inst *Instance) Ping() (err error) {
	return inst.RunCmd("ansible", inst.Host, "-m", "win_ping")
}

func (inst *Instance) Playbook(arg ...string) (err error) {
	cmds := make([]string, 0, 3+len(arg))
	cmds = append(cmds, "playbook.yaml", "--inventory", inst.Root+"/configs/hosts.ini")
	cmds = append(cmds, arg...)

	return inst.RunCmd("ansible-playbook", cmds...)
}

func (inst *Instance) View() (string, error) {
	var (
		bts    []byte
		target string
		err    error
	)

	target = filepath.Join(inst.WorkPath(), "status.log")

	if _, err = os.Stat(target); err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}

	if bts, err = ioutil.ReadFile(target); err != nil {
		return "", err
	}

	return string(bts), err
}

func (inst *Instance) Run() (err error) {
	if err = inst.writeStatus("starting"); err != nil {
		return err
	}

	if err = inst.Playbook("--tags", "run"); err != nil {
		return err
	}

	return nil
}

// sync log and status
func (inst *Instance) Sync() (err error) {
	return inst.Playbook("--tags", "get_log,get_status")
}

func (inst *Instance) SyncLog() (err error) {
	return inst.Playbook("--tags", "get_log")
}

func (inst *Instance) SyncStatus() (err error) {
	return inst.Playbook("--tags", "get_status")
}

func (inst *Instance) Kill() (err error) {
	if err = inst.Playbook("--tags", "execute", "--extra-vars", "call=kill"); err != nil {
		return err
	}

	if err = inst.writeStatus("killed"); err != nil {
		return err
	}

	if err = inst.end(); err != nil {
		return err
	}

	return nil
}
