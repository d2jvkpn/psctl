package main

import (
	"os"

	"psctl/cmd"
	"psctl/pkg/misc"

	"github.com/spf13/cobra"
)

func init() {
	os.Setenv("PATH", os.Getenv("PATH")+":"+"~/.local/bin")
	misc.SetLogTimeFmt()
}

func main() {
	rootCmd := &cobra.Command{Use: "pixel streaming controller"}

	rootCmd.AddCommand(cmd.NewVersionCmd("version"))
	rootCmd.AddCommand(cmd.NewDemoCmd("demo"))
	rootCmd.AddCommand(cmd.NewInstanceCmd("instance"))

	rootCmd.Execute()
}
