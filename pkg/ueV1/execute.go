package ueV1

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
)

func (inst *Instance) RunCmd(name string, arg ...string) (err error) {
	var buf bytes.Buffer

	cmd := exec.Command(name, arg...)
	cmd.Dir = inst.WorkPath()
	cmd.Env = append(cmd.Env, "ANSIBLE_LOG_PATH="+filepath.Join("logs", "ansible.log"))
	cmd.Stdout, cmd.Stderr = &buf, &buf

	if inst.Debug {
		log.Printf(">>> $ %s %s", name, strings.Join(arg, " "))
	}

	if err = cmd.Run(); err != nil {
		if out := buf.String(); len(out) > 0 {
			return fmt.Errorf("%s\nRunCmd error: %w", out, err)
		}
		return fmt.Errorf("RunCmd error: %w", err)
	}
	return nil
}

func (inst *Instance) prepare() (err error) {
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
	cmds := make([]string, 0, 4+len(arg))
	// cmds = append(cmds, "playbook.yaml", "--inventory", inst.Root+"/configs/hosts.ini")
	cmds = append(cmds, "playbook.yaml")
	if inst.Debug {
		cmds = append(cmds, "-v")
	}
	cmds = append(cmds, "--inventory", "../../../../configs/hosts.ini")
	cmds = append(cmds, arg...)

	return inst.RunCmd("ansible-playbook", cmds...)
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
	if err = inst.writeStatus("kill"); err != nil {
		return err
	}

	if err = inst.Playbook("--tags", "execute", "--extra-vars", "call=kill"); err != nil {
		return err
	}

	if err = inst.end(); err != nil {
		return err
	}

	return nil
}
