package core

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
)

func handleInvalidErrors(err error, args *argparse.Parser) {
	if err != nil {
		fmt.Println(args.Usage(nil))
		os.Exit(0)
	}
}

//ParseCommandLineArgs : Execute package manager based on supplied command line args
func ParseCommandLineArgs(args []string) {
	/*
		Route to appropriate functions based on the sub-commands
	*/
	parserMain := argparse.NewParser(
		"help",
		"A small package manager written for Pavlos container runtime",
	)

	//implement the sub-commands:
	rootfs := parserMain.NewCommand(
		"rootfs",
		"Perform rootfs functions, like create, delete and delete",
	)

	//rootfs commands
	rootfsCreate := rootfs.NewCommand("create", "Create a rootfs image")
	rootfsName := rootfsCreate.String("n", "name", &argparse.Options{Help: "Name of rootfs image", Required: true})
	rootfsURI := rootfsCreate.String("u", "uri", &argparse.Options{Help: "Rootfs Uri", Required: true})

	rootfsList := rootfs.NewCommand("list", "Get a list of available rootfs images")
	rootfsDelete := rootfs.NewCommand("delete", "Delete a rootfs image")
	rootfsDelName := rootfsDelete.String("n", "name", &argparse.Options{Help: "Rootfs image name to delete", Required: true})

	configs := parserMain.NewCommand(
		"config",
		"Perform list, delete and create of predefined configs",
	)

	configsCreate := configs.NewCommand("create", "Creates config from json file")
	configsName := configsCreate.String("n", "name", &argparse.Options{Help: "Name of the rootfs image", Required: true})
	configsFP := configsCreate.String("f", "file", &argparse.Options{Help: "file path to JSON file", Required: true})

	configsList := configs.NewCommand("list", "Lists saved configs")
	configsDelete := configs.NewCommand("delete", "Deletes a config from saved configs")
	configsDelName := configsDelete.String("n", "name", &argparse.Options{Help: "Name of config file", Required: true})

	configsEdit := configs.NewCommand("edit", "Edits a configuration file with vi, export EDITOR env variable targetting the custom editor you would like to use.")
	configsEditName := configsEdit.String("n", "name", &argparse.Options{Help: "Name of the config file to edit"})

	err := parserMain.Parse(args)
	handleInvalidErrors(err, parserMain)

	if rootfs.Happened() && rootfsCreate.Happened() {
		if *rootfsName == "" || *rootfsURI == "" {
			fmt.Println("Did not provide name or uri of rootfs image")
			fmt.Println(rootfsCreate.Help(nil))
			os.Exit(0)
		}

		CreateRootFs(rootfsName, rootfsURI)
		os.Exit(0)
	}
	if rootfs.Happened() && rootfsList.Happened() {
		ListRootFsImages()
		os.Exit(0)
	}
	if rootfs.Happened() && rootfsDelete.Happened() {
		if *rootfsDelName == "" {
			fmt.Println("Did not provide rootfs name")
			fmt.Println(rootfsDelete.Help(nil))
			os.Exit(0)
		}
		DeleteRootFs(rootfsDelName)
		os.Exit(0)
	}

	if configs.Happened() && configsCreate.Happened() {
		if *configsName == "" || *configsFP == "" {
			fmt.Println("Did not provide config name or configs file path")
			os.Exit(0)
		}

		CreateConfigFromFile(configsFP, configsName)
		os.Exit(0)
	}

	if configs.Happened() && configsList.Happened() {
		ListSavedConfigs()
		os.Exit(0)
	}

	if configs.Happened() && configsDelete.Happened() {
		if *configsDelName == "" {
			fmt.Println("Did not provide config name")
			os.Exit(0)
		}

		DeleteRootfsConfig(configsDelName)
		os.Exit(0)
	}

	if configs.Happened() && configsEdit.Happened() {
		if *configsEditName == "" {
			fmt.Printf("Did not provide config name.")
			os.Exit(0)
		}

		EditRootfsConfig(configsEditName)
		os.Exit(0)
	}
}
