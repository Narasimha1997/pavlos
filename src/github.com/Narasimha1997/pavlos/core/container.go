package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/jaypipes/ghw"
)

var home string = os.Getenv("HOME")

//DefaultRootFsAbsPath : Default rootfs for kontainer, change this according to your needs
var DefaultRootFsAbsPath string = filepath.Join(home, ".rootfs")

//IsolationOpts : Linux syscall isolation options
type IsolationOpts struct {
	EnableUTS   bool `json:"enableUTS"`   // Unix-Time-Sharing : Host names isolation
	EnablePID   bool `json:"enablePID"`   // Enables Process Isolation
	EnableRoot  bool `json:"enableRoot"`  // Enables chroot fs Isolation
	EnableNetNs bool `json:"enableNetNs"` // Enables Network Namespace isolation (Experimental still!)
}

//InternalFlags : Internal flags for container configuration, more flags will be added
type InternalFlags struct {
	HasNvidiaDevices    bool   `json:"hasNvidia"`
	NvidiaDeviceIndexes []int  `json:"deviceIdx"`
	Inet4Address        string `json:"inet4Address"`
	MappedPorts         []int  `json:"mappedPorts"`
	Devices             []*ghw.PCIDevice
}

//ContainerOpts : default container option
type ContainerOpts struct {
	Name                    string        `json:"name"`
	EnableIsolation         bool          `json:"enableIsolation"`
	EnableResourceIsolation bool          `json:"enableResourceIsolation"`
	IsolationOpts           IsolationOpts `json:"isolationOptions"`
	RootFs                  string        `json:"rootFs"`
	NvidiaGpus              []int         `json:"nvidiaGpus"`
	RuntimeArgs             []string      `json:"runtimeArgs"`
	Internals               InternalFlags `json:"internals"`
	InitScript              string        `json:"initScript"`
	IP                      string        `json:"ip"`
}

func catchError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

//ResolveRootFs : Finds the rootfs specified, handles default options
func ResolveRootFs(config *ContainerOpts) {
	if config.RootFs == "default" {
		fmt.Println("Setting ubuntu-base-18.04-base-amd64 as default rootfs image")
		config.RootFs = filepath.Join(DefaultRootFsAbsPath, "default")
	} else {
		fmt.Printf("Searching registry location for rootfs %s.", config.RootFs)
		config.RootFs = filepath.Join(DefaultRootFsAbsPath, config.RootFs)
	}

	//check if the rootfs exist :
	_, err := os.Stat(config.RootFs)
	catchError(err)
	if os.IsNotExist(err) {
		fmt.Printf("[Error] Rootfs %s not found, exiting runtime.\n", config.RootFs)
		os.Exit(0)
	}

}

//PciDeviceFilter : Filter out PCI devices
func PciDeviceFilter(config *ContainerOpts) {
	devices := PciDeviceProbe()

	if len(devices) == 0 {
		fmt.Printf("No PCI devices were found")
		config.Internals.HasNvidiaDevices = false
		return
	}

	config.Internals.Devices = devices

	if len(config.NvidiaGpus) == 0 {
		return
	}

	//query for list of supported devices
	devIndex := 0
	for _, dev := range devices {
		//NVIDA GPU probe
		if dev.Vendor.Name == "NVIDIA Corporation" && dev.Class.Name == "Display controller" {
			//fmt.Println("NVIDA GPU detected")
			fmt.Println(dev)
			config.Internals.NvidiaDeviceIndexes = append(config.Internals.NvidiaDeviceIndexes, devIndex)
			config.Internals.HasNvidiaDevices = true
		}
		devIndex++
	}
}

//PciDeviceProbe : Probe for active PCI devices
func PciDeviceProbe() []*ghw.PCIDevice {
	pci, err := ghw.PCI()
	if err != nil {
		fmt.Printf("PCI device probe failed")
		os.Exit(0)
	}

	devices := pci.ListDevices()
	return devices
}

//ProduceUnsupportedWarnings : Some features are not supported right now like network isolation
func ProduceUnsupportedWarnings(config *ContainerOpts) {

	if !config.EnableIsolation {
		fmt.Println("[Warning] No isolation is enabled. The container is executed as host process")
	}

	if config.EnableResourceIsolation {
		fmt.Println("[Warning] Resource isolation is not supported, disabling it")
		config.EnableResourceIsolation = false
	}

	if config.IsolationOpts.EnableNetNs {
		fmt.Println("[Warning] Network isolation is still beta")
		config.IsolationOpts.EnableNetNs = true
	}
}

//LoadConfigFromJSON : Configuration loader
func LoadConfigFromJSON(jsonFile string) ContainerOpts {
	jsonData, err := ioutil.ReadFile(jsonFile)
	catchError(err)
	options := ContainerOpts{}

	//fmt.Println(string(jsonData))
	json.Unmarshal([]byte(jsonData), &options)
	catchError(err)

	ProduceUnsupportedWarnings(&options)
	ResolveRootFs(&options)
	//fmt.Println(options)
	PciDeviceFilter(&options)
	return options
}

//PrettyPrintConfig : Disply config details
func PrettyPrintConfig(options *ContainerOpts) {
	fmt.Println("=========== Container Info ================")
	fmt.Printf("Container Name : %s\n", options.Name)
	fmt.Printf("Supports Isolation : %v\n", options.EnableIsolation)
	fmt.Printf("Root File system : %s\n", options.RootFs)
	fmt.Println("--- NVIDIA Runtime info ----")
	fmt.Printf("Requested devices : %v\n", options.NvidiaGpus)
	if options.Internals.HasNvidiaDevices && len(options.NvidiaGpus) != 0 {
		fmt.Printf("NVIDA support enabled. GPUs found : \n")
		for _, deviceIndex := range options.Internals.NvidiaDeviceIndexes {
			fmt.Print("*   ")
			fmt.Print(options.Internals.Devices[deviceIndex])
			fmt.Print("\n")
		}
	}
	fmt.Println("===========================================")
}
