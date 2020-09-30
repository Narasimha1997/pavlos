package core

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

//Network Manager for kontainer runtime

func networkMust(err error) {
	if err != nil {
		panic(err)
	}
}

//An interace to system
func system(command string, arguments string) {
	args := strings.Split(arguments, " ")
	executor := exec.Command(command, args...)
	executor.Stderr = os.Stderr
	executor.Stdout = os.Stdout

	networkMust(executor.Run())
}

//NetworkIfaceSetup : Creates network setup
func NetworkIfaceSetup(options *ContainerOpts) {
	id := "213"

	veth := fmt.Sprintf("link add veth%s type veth peer name veth1", id)
	system("ip", veth)

	vethUp := fmt.Sprintf("link set veth%s up", id)
	system("ip", vethUp)
}

//NetworkAddNamespace : Adds namespace
func NetworkAddNamespace(pid int) {
	vethAddNs := fmt.Sprintf("link set veth1 netns %d", pid)
	system("ip", vethAddNs)
}

//SetupContainerNetworking : Attaches isolated network to container
func SetupContainerNetworking(options *ContainerOpts, pid int) {
	system("ip", "link set veth3 up")
	vethIP := fmt.Sprintf("ip addr add %s/24 dev veth1", options.IP)
	system("ip", vethIP)
	system("route", "add default gw 172.16.0.100 veth1")
}
