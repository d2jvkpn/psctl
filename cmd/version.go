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

	BuildTime   string
	BuildBranch string
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
				Version:     version,
				GoVersion:   strings.Replace(runtime.Version(), "go", "", 1),
				BuildBranch: BuildBranch,
				BuildTime:   BuildTime,
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
	Version     string `json:"version"`
	GoVersion   string `json:"goVersion"`
	BuildBranch string `json:"buildBranch"`
	BuildTime   string `json:"buildTime"`
}

func (info BuildInfo) String() string {
	return fmt.Sprintf(
		"psctl:\n  version: %s\n  go version: %s\n  build time: %s\n",
		info.Version, info.GoVersion, info.BuildTime,
	)
}

func (info BuildInfo) JSON() string {
	bts, _ := json.Marshal(info)
	return string(bts)
}
