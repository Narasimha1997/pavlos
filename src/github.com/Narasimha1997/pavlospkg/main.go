package main

import (
	"os"

	"github.com/Narasimha1997/pavlospkg/core"
)

func main() {
	var Args []string = os.Args

	if len(Args) < 2 {
		core.ShowHelp()
		os.Exit(0)
	}

	operation := Args[1]

	core.SetupDefaultDir()

	switch operation {
	case "list":
		core.ListRootFsImages()
		break
	case "create":
		core.CreateRootFs()
		break
	case "remove":
		core.DeleteRootFs(Args[2])
	case "setup-default":
		core.SetupDefault()
	default:
		core.ShowHelp()
	}
}
