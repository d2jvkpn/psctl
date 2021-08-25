package cmd

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var (
	//go:embed version.txt
	version string
)

func init() {
	version = strings.Fields(version)[0]
}

func NewVersionCmd(name string) (command *cobra.Command) {
	return &cobra.Command{
		Use: name,
		// Short: "program version",
		Long: `ue_cloud version`,

		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	}
}
