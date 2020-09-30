package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}
func containerMust(err error) {
	if err != nil {
		fmt.Println("Container exited")
	}
}

func isolateGPUDevice(gpuDeviceID int, rootFsPath string) {
	//nvidia-container-cli --load-kmods configure --ldconfig=@/sbin/ldconfig.real --no-cgroups --utility --device 0 $(pwd)
	executor := "nvidia-container-cli"
	options := []string{"--load-kmods"}

	deviceOpt := "--device=" + strconv.Itoa(gpuDeviceID)

	config := []string{"configure", "--ldconfig=@/sbin/ldconfig.real", "--no-cgroups", "--utility", deviceOpt, rootFsPath}

	args := append(options, config...)
	//fmt.Println(args)
	command := exec.Command(executor, args...)
	command.Stderr = os.Stderr
	command.Stdout = os.Stdout

	must(command.Run())
}

func makeNvidiaContainerRunime(options *ContainerOpts) {

	if len(options.NvidiaGpus) == 0 {
		return
	}

	if !options.Internals.HasNvidiaDevices {
		panic("NVIDIA container runtime not supported")
	}

	for index := range options.Internals.NvidiaDeviceIndexes {
		//fmt.Printf("Mapping GPU device %d\n", index)
		isolateGPUDevice(index, options.RootFs)
	}
}

func customBind(args []string) {
	command := exec.Command("mount", args...)
	must(command.Run())
}

func addDNS() {
	//Adds host DNS configuration to container
	input, err := ioutil.ReadFile("/etc/resolv.conf")
	must(err)

	must(ioutil.WriteFile("etc/resolv.conf", input, 0664))
}

func prepareContainerLinux() {

	//file systems
	must(syscall.Mount("proc", "proc", "proc", 0, ""))
	must(syscall.Mount("tmp", "tmp", "tmpfs", 0, ""))
	must(syscall.Mount("/sys", "sys", "sysfs", 0, ""))
	must(syscall.Mount("/dev", "dev", "devtmpfs", syscall.MS_BIND, ""))

	//DNS resolution
	addDNS()
}

func unmountFileSystems() {
	must(syscall.Unmount("proc", 0))
	must(syscall.Unmount("tmp", 0))
	must(syscall.Unmount("sys", 0))
}

func makeConfigsPassthrough(config *ContainerOpts) string {
	byteData, err := json.Marshal(config)
	if err != nil {
		panic("Argument passthrough error")
	}
	return string(byteData)
}

func jsonToConfigs(jsonData string, containerConfig *ContainerOpts) {
	byteData := []byte(jsonData)
	must(json.Unmarshal(byteData, &containerConfig))
}

//CreateContainer : Main function that creates a container
func CreateContainer(configFile string) {
	_, err := os.Stat(configFile)
	catchError(err)

	if os.IsNotExist(err) {
		fmt.Printf("[Error ] Container config %s not found, exiting runtime.\n", configFile)
		os.Exit(0)
	}

	options := LoadConfigFromJSON(configFile)

	konArgs := append([]string{makeConfigsPassthrough(&options)}, options.RuntimeArgs...)
	executorHandle := MakeContainerRuntime(&options, konArgs)
	executorHandle.Stdout = os.Stdout
	executorHandle.Stderr = os.Stderr
	executorHandle.Stdin = os.Stdin

	clearSrc := exec.Command("clear")
	clearSrc.Stdout = os.Stdout

	must(clearSrc.Run())
	must(executorHandle.Run())
}

//ResolveRuntimeHooks : The function that resolves runtime hooks for container
func ResolveRuntimeHooks(options *ContainerOpts) {
	//More elegant hooks will be provided later
	makeNvidiaContainerRunime(options)
}

//ContainerRuntime : The function child that runs the container
func ContainerRuntime() {
	configJSON := os.Args[2]
	//fmt.Println("Inside container")

	options := ContainerOpts{}
	jsonToConfigs(configJSON, &options)

	if options.EnableResourceIsolation {
		fmt.Printf("CGroup isolation will be added soon\n")
	}

	if options.IsolationOpts.EnableNetNs {
		pidContainer := os.Getpid()
		NetworkAddNamespace(pidContainer)
	}

	command := exec.Command(os.Args[3], os.Args[4:]...)

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	//set up hostname:
	must(os.Chdir(options.RootFs))

	prepareContainerLinux()
	ResolveRuntimeHooks(&options)
	must(syscall.Chroot(options.RootFs))
	must(os.Chdir("/"))

	if options.IsolationOpts.EnableNetNs {
		SetupContainerNetworking(&options, os.Getpid())
	}

	must(syscall.Sethostname([]byte(options.Name)))
	PrettyPrintConfig(&options)

	containerMust(command.Run())

	if options.IsolationOpts.EnableNetNs {
		system("ip", "link delete veth3")
	}
	fmt.Println("Container exited")

	//unmountFileSystems()
}
