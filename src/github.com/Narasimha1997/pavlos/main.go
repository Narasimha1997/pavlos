package main

import (
	"fmt"
	"os"

	"github.com/Narasimha1997/pavlos/core"

	"github.com/akamensky/argparse"
)

func handelArgsError(err error, args *argparse.Parser) {
	if err != nil {
		fmt.Println(args.Usage(nil))
		os.Exit(0)
	}
}

var pavlosDesc string = "Pavlos is a Linux rootfs emulator and container runtime that can emulate" +
	" any Linux Rootfs Image. Pavlos also supports GPU device assignment to container runtimes using libnvidiacontainer and PCI device scanning"

func printConfigExample() {
	exampleConfig := `
	{
		"name" : "TestContainer",
		"enableIsolation" : true,
		"enableResourceIsolation" : true,
		"isolationOptions" : {
			"enableUTS" : true,
			"enablePID" : true,
			"enableRoot" : true,
			"enableNetNs" : true
		},
		"rootFs" : "/path/to/root/fs : example $HOME/ubuntu_16.05_base",
		"nvidiaGpus" : [0],
		"runtimeArgs" : ["/bin/bash"]
	}
	`
	fmt.Println(exampleConfig)
}

func main() {
	//fmt.Printf("Hello, kontainer")

	if len(os.Args) < 2 {
		fmt.Println("Invalid command-line arguments")
		fmt.Println("Check sudo pavlos run --help")
		os.Exit(0)
	}

	arg := os.Args[1]

	if arg == "run" {

		parser := argparse.NewParser(
			"help",
			pavlosDesc,
		)

		run := parser.NewCommand("run", "Run a pavlos container")
		fromFile := run.Flag("f", "from-file", &argparse.Options{Help: "Take file input config", Required: false})
		config := run.String("c", "config", &argparse.Options{Help: "Config file/name", Required: true})

		err := parser.Parse(os.Args)
		handelArgsError(err, parser)

		if run.Happened() {
			if *config == "" {
				fmt.Println("Did not provide config")
				fmt.Println(run.Usage(nil))
				os.Exit(0)
			}

			core.CreateContainer(*config, *fromFile)
		}
	}

	if arg == "container" {
		core.ContainerRuntime()
	}
}
