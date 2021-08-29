package cmd

import (
	// "flag"
	"fmt"
	"log"
	"time"

	"psctl/pkg/ueV1"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func NewDemoCmd(name string) (command *cobra.Command) {
	var (
		fSet *pflag.FlagSet
		base ueV1.InstanceBase
	)

	command = &cobra.Command{
		Use:   name,
		Short: `ue demo`,
		Long:  `run a ue instance as demo`,

		Run: func(cmd *cobra.Command, args []string) {
			var err error

			// fmt.Println("~~~", args)
			inst := base.Inst(0)
			fmt.Printf("### Instance %s: %s\n", inst.WorkPath(), inst)

			if err = demo(inst); err != nil {
				log.Fatal(err)
			}
		},
	}

	fSet = command.Flags()

	fSet.StringVar(&base.Host, "host", "win01", "hostname in the ansible inventory file")
	fSet.StringVar(&base.Project, "project", "D:\\projects\\Project_001", "project directory in windows")
	fSet.StringVar(&base.Program, "program", "Project", "program name without .exe extension")
	fSet.StringVar(&base.SwsIp, "swsIp", "192.168.0.171", "signaling and web server(SWS) ip address")
	fSet.StringVar(&base.SwsPort, "swsPort", "8201", "SWS port")

	/*
		base := ueV1.InstanceBase{
			Host:    "win01",
			Project: "Project_001",
			Program: "Project",
			SwsIp:   "192.168.0.171",
			SwsPort: "8204",
		}
	*/

	return command
}

func demo(inst ueV1.Instance) (err error) {
	if err = inst.NewPlaybook(true); err != nil {
		return err
	}

	/*
		log.Println(">>> instance: ping")
		if err = inst.Ping(); err != nil {
			return err
		}
	*/

	log.Println(">>> instance: start")
	if err = inst.Start(); err != nil {
		return err
	}
	log.Println(">>> instance: run")

	time.Sleep(10 * time.Second)
	log.Println(">>> instance: sync")
	if err = inst.Sync(); err != nil {
		return err
	}

	time.Sleep(90 * time.Second)
	log.Println(">>> instance: sync")
	if err = inst.Sync(); err != nil {
		return err
	}

	time.Sleep(20 * time.Second)
	log.Println(">>> instance: kill")
	if err = inst.Kill(); err != nil {
		return err
	}

	return nil
}
