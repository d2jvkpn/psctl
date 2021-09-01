package ueV1

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"psctl/pkg/misc"
)

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

func (inst *Instance) Ping() (err error) {
	cmd := exec.Command("ansible", inst.Host, "-m", "win_ping")
	return cmd.Run()
}

func (inst *Instance) Playbook(arg ...string) (err error) {
	cmds := make([]string, 0, 3+len(arg))
	cmds = append(cmds, "playbook.yaml", "--inventory", inst.Root+"/configs/hosts.ini")
	cmds = append(cmds, arg...)

	return inst.RunCmd("ansible-playbook", cmds...)
}

func (inst *Instance) View() (string, error) {
	var (
		ok     bool
		bts    []byte
		target string
		err    error
	)

	target = filepath.Join(inst.WorkPath(), "status.log")

	if ok, err = misc.FileExists(target); err != nil {
		return "", err
	}
	if !ok {
		return "", nil
	}

	if bts, err = ioutil.ReadFile(target); err != nil {
		return "", err
	}

	return string(bts), err
}

func (inst *Instance) Start() (err error) {
	if inst.NewPlaybook(false); err != nil {
		return nil
	}

	if err = inst.writeStatus("start"); err != nil {
		return err
	}

	if err = inst.Playbook("--tags", "start"); err != nil {
		return err
	}

	return nil
}

// sync log and status
func (inst *Instance) Sync() (err error) {
	return inst.Playbook("--tags", "get_log")
	// return inst.Playbook("--tags", "get_log,get_status")
}

func (inst *Instance) Status() (err error) {
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
