package cmd

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	//go:embed version.txt
	version string
)

func init() {
	version = strings.Fields(version)[0]
}

func NewVersionCmd(name string) (command *cobra.Command) {
	var (
		jsonFmt bool
		fSet    *pflag.FlagSet
	)

	command = &cobra.Command{
		Use: name,
		// Short: "program version",
		Long: `psctl version`,

		Run: func(cmd *cobra.Command, args []string) {
			info := BuildInfo{
				Version:   version,
				GoVersion: strings.Replace(runtime.Version(), "go", "", 1),
			}
			if jsonFmt {
				fmt.Println(info.JSON())
			} else {
				fmt.Println(info)
			}
		},
	}

	fSet = command.Flags()
	fSet.BoolVar(&jsonFmt, "json", false, "output in json format")

	return command
}

type BuildInfo struct {
	Version   string `json:"version"`
	GoVersion string `json:"goVersion"`
}

func (info BuildInfo) String() string {
	return fmt.Sprintf(
		"psctl:\n  version: %s\n  go version: %s",
		info.Version, info.GoVersion,
	)
}

func (info BuildInfo) JSON() string {
	bts, _ := json.Marshal(info)
	return string(bts)
}
