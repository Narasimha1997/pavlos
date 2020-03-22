package main

import (
	"fmt"
	"os"
)

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
	if len(os.Args) == 1 || len(os.Args) == 2 {
		fmt.Println("Usage : sudo pavlosc run CONFIG_FILE.json")
		fmt.Println(" For more info : pavlosc help show")
		os.Exit(0)
	}

	argument := os.Args[1]

	switch argument {
	case "run":
		CreateContainer(os.Args[2])
	case "container":
		ContainerRuntime()
	default:
		fmt.Println(
			"pavalosc is a CLI for pavalos container runtime -- a light weight container emulator, it can run any" +
				" arbitrary Linux Rootfs image as a container with NVIDA gpu support" +
				" you modify pavalos source to add support for many other PCI devices.")
		fmt.Println("To run : sudo pavalosc run CONTAINER_CONFIG_FILE.json")
		fmt.Println("Example config file : ")
		fmt.Println("==================================")
		printConfigExample()
	}
}
