package cmd

import (
	"fmt"
	"log"

	"psctl/pkg/ueV1"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func NewInstanceCmd(name string) (command *cobra.Command) {
	var (
		file string
		call string
		fSet *pflag.FlagSet
	)

	command = &cobra.Command{
		Use:   name,
		Short: `run a instance`,
		Long:  `run a ue streamer instance: view, run, sync, syncLog, syncStatus, kill`,

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

			fmt.Println("### Instance:", inst.InstanceBase)

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
	switch call {
	case "view":
		var output string
		if output, err = inst.View(); err != nil {
			return err
		}
		fmt.Print(output)
	case "run":
		if err = inst.NewPlaybook(true); err != nil {
			return err
		}
		err = inst.Run()
	case "sync":
		err = inst.Sync()
	case "syncLog":
		err = inst.SyncLog()
	case "syncStatus":
		err = inst.SyncStatus()
	case "kill":
		err = inst.Kill()
	//case "clear":
	//	err = inst.Clear()
	default:
		err = fmt.Errorf("unknow call")
	}

	return err
}
