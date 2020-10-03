package main

import (
	"os"

	"github.com/Narasimha1997/pavlospkg/core"
)

func main() {
	core.SetupDefaultDir()
	core.ParseCommandLineArgs(os.Args)
}
