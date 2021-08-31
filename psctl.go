package main

import (
	_ "embed"
	"os"
	"strings"

	"psctl/cmd"
	"psctl/pkg/misc"

	"github.com/spf13/cobra"
)

var (
	//go:embed .version
	version string
)

func init() {
	os.Setenv("PATH", os.Getenv("PATH")+":"+"~/.local/bin")
	misc.SetLogTimeFmt()
	version = strings.Fields(version)[0]
}

func main() {
	rootCmd := &cobra.Command{Use: "pixel streaming controller"}

	rootCmd.AddCommand(cmd.NewMd5Cmd("md5"))
	rootCmd.AddCommand(cmd.NewDemoCmd("demo"))
	rootCmd.AddCommand(cmd.NewLoadCmd("load"))
	rootCmd.AddCommand(cmd.NewVersionCmd("version", version))

	rootCmd.Execute()
}
