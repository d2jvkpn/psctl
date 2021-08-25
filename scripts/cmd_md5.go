package main

import (
	"fmt"
	"os"

	"psctl/pkg/misc"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "not input args")
		os.Exit(1)
	}

	fmt.Println(misc.CmdMd5(os.Args[1:]))
}
