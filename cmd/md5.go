package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"psctl/pkg/misc"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func NewMd5Cmd(name string) (command *cobra.Command) {
	var (
		jsonFmt bool
		fSet    *pflag.FlagSet
	)

	command = &cobra.Command{
		Use:   name,
		Short: `command md5sum`,
		Long:  `calculate md5sum of command`,

		Run: func(cmd *cobra.Command, args []string) {
			fArgs := fSet.Args()
			if len(fArgs) < 1 {
				log.Fatalln("please provide args(command)")
			}

			if jsonFmt {
				bts, _ := json.Marshal(map[string][]string{"commandline": fArgs})
				fmt.Printf("%s\n", bts)
			}

			fmt.Println(misc.CmdMd5(fArgs))
		},
	}

	fSet = command.Flags()
	fSet.BoolVar(&jsonFmt, "json", false, "output in json format")

	return command
}
