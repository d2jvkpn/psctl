package cmd

import (
	"fmt"
	"log"

	"psctl/pkg/ueV1"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func NewLoadCmd(name string) (command *cobra.Command) {
	var (
		file string
		call string
		fSet *pflag.FlagSet
	)

	command = &cobra.Command{
		Use:   name,
		Short: `load a project and execute`,
		Long: `load a project(ue streamer instance): pstcl load <project.yaml>  <call>
  call: [new, start, kill, restart, ping, sync, status, view]`,

		Run: func(cmd *cobra.Command, args []string) {
			var (
				fArgs []string
				err   error
				inst  ueV1.Instance
			)

			fArgs = fSet.Args()
			if len(fArgs) < 2 {
				log.Fatalln("please provide instance file and call")
			}
			file, call = fArgs[0], fArgs[1]

			if inst, err = ueV1.InstanceFromFile(file); err != nil {
				log.Fatal(err)
			}

			inst.Debug = true
			fmt.Printf("### Instance %s: %s\n", inst.WorkPath(), inst.InstanceBase)

			if err = callFunc(&inst, call); err != nil {
				log.Fatal(err)
			}
		},
	}

	fSet = command.Flags()
	/*
		fSet.StringVar(&file, "file", "", "instance toml/yaml/json file")
		fSet.StringVar(&call, "call", "", "call a function")
	*/

	return command
}

func callFunc(inst *ueV1.Instance, call string) (err error) {
	var (
		output string
	)

	switch call {
	case "new":
		err = inst.NewPlaybook(true)
	case "start":
		err = inst.Start()
	case "kill":
		err = inst.Kill()
	case "restart":
		if err = inst.Kill(); err != nil {
			break
		}
		err = inst.Start()
	///
	//case "clear":
	//	err = inst.Clear()
	case "ping":
		err = inst.Ping()
	case "sync":
		err = inst.Sync()
	case "status":
		err = status(inst)
	case "view":
		if output, err = inst.View(); err == nil {
			fmt.Print(output)
		}
	default:
		err = fmt.Errorf("unknown call")
	}

	return err
}

func status(inst *ueV1.Instance) (err error) {
	var (
		yes    bool
		output string
	)

	if yes, err = inst.Exists(); err != nil {
		return err
	} else if !yes {
		return nil
	}

	if err = inst.Status(); err != nil {
		return err
	}
	if output, err = inst.View(); err != nil {
		return err
	}

	fmt.Print(output)
	return nil
}
