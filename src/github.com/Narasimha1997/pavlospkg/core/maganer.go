package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//RootFSRoot : Root-filesystem root path
var RootFSRoot string = filepath.Join(os.Getenv("HOME"), ".rootfs/images")

//RootFSConfigPath : Base filesystem path to save config details
var RootFSConfigPath string = filepath.Join(os.Getenv("HOME"), ".rootfs/configs")

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

//utility functions:

func registryPainc(err error) {
	if err != nil {
		fmt.Printf("[Registry] Error %v\n", err)
		os.Exit(0)
	}
}

func checkNameExists(target string) {
	files, err := ioutil.ReadDir(RootFSRoot)
	registryPainc(err)

	for _, image := range files {
		if strings.Compare(target, image.Name()) == 0 {
			fmt.Printf("Name %s already exists in registry , try different name\n", target)
			os.Exit(0)
		}
	}
}

func runWget(name, url string) {

	if _, err := os.Stat(name + ".tar.gz"); err == nil {
		fmt.Printf("File exists,so not downloading.\n")
		return
	}

	command := exec.Command("wget", url, "-O", name+".tar.gz")
	command.Stderr = os.Stderr
	command.Stdout = os.Stdout

	registryPainc(command.Run())
}

//DeleteRootFs : Deletes root-fs image
func DeleteRootFs(name *string) {
	path := filepath.Join(RootFSRoot, *name)

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		fmt.Printf("Registry %s not exists, so skipped deletion\n", *name)
		os.Exit(0)
	}

	registryPainc(os.RemoveAll(path))
	fmt.Printf("deleted %s\n", path)
}

func runTarExtract(name string) {

	os.Mkdir(filepath.Join(RootFSRoot, name), 0777)

	fsPath := filepath.Join(RootFSRoot, name)
	command := exec.Command("tar", "-C", fsPath, "-xf", name+".tar.gz")
	command.Stderr = os.Stderr
	command.Stdout = os.Stdout

	registryPainc(command.Run())
}

func removeTarArchive(name string) {
	registryPainc(os.Remove(name + ".tar.gz"))
}

func setupRootFsImage(name string, url string) {
	checkNameExists(name)
	runWget(name, url)
	runTarExtract(name)
	removeTarArchive(name)
}

//ListRootFsImages : Lists all root-fs images
func ListRootFsImages() {
	files, err := ioutil.ReadDir(RootFSRoot)
	registryPainc(err)

	fmt.Println("* RootFS Images ")
	fmt.Println("=====================")

	for _, image := range files {
		fmt.Printf("%s\n", image.Name())
	}
}

//ShowHelp : Shows help information
func ShowHelp() {
	fmt.Println("pavlos-registry is a root-fs management system for pavlos-container-runtime.")
	fmt.Println("Available options :")
	fmt.Println("\t list : lists all the root-fs registered.")
	fmt.Println("\t create : an interactive cli that creates a container-root-fs image.")
	fmt.Println("\t remove <image-name> : removes the rootfs-image from registry.")
	fmt.Println("\t setup-default : creates a default ubuntu-18.04 rootfs image")
}

//CreateRootFs : creates a root-fs image, a series of prompts will be followed
func CreateRootFs(rootfsName *string, rootfsURI *string) {
	setupRootFsImage(*rootfsName, *rootfsURI)
	fmt.Println("Created rootfs " + *rootfsName)
}

//SetupDefault : Sets-up default root-fs image
func SetupDefault() {
	setupRootFsImage("default", "http://cdimage.ubuntu.com/ubuntu-base/releases/18.04/release/ubuntu-base-18.04-base-amd64.tar.gz")
}

//SetupDefaultDir Creates .rootfs at home if not exist
func SetupDefaultDir() {
	_, err := os.Stat(RootFSRoot)
	if os.IsNotExist(err) {
		err := os.MkdirAll(RootFSRoot, 0700)
		registryPainc(err)
	}

	_, err = os.Stat(RootFSConfigPath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(RootFSConfigPath, 0700)
		registryPainc(err)
	}
}

/*
	Configuration management
*/

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

//ListSavedConfigs : Lists all config files
func ListSavedConfigs() {
	if !pathExists(RootFSConfigPath) {
		fmt.Println("RootFS config path does not exist")
		os.Exit(0)
	}

	files, err := ioutil.ReadDir(RootFSConfigPath)
	registryPainc(err)

	fmt.Println("Rootfs Configs")
	fmt.Println("=====================")

	for _, config := range files {
		fmt.Printf("%s\n", config.Name())
	}
}

//CreateConfigFromFile : Creates config given a JSON file
func CreateConfigFromFile(file *string, rootfsName *string) {

	if !pathExists(*file) {
		fmt.Printf("Configuration source file not found")
		os.Exit(0)
	}

	jsonData, err := ioutil.ReadFile(*file)
	registryPainc(err)

	options := ContainerOpts{}

	json.Unmarshal([]byte(jsonData), &options)
	registryPainc(err)

	options.RootFs = *rootfsName

	//save file in registry location
	byteData, err := json.Marshal(options)
	registryPainc(err)

	fpOutput := filepath.Join(RootFSConfigPath, *rootfsName)
	err = ioutil.WriteFile(fpOutput, byteData, 0644)
	registryPainc(err)

	fmt.Printf("Successfully saved configuration %s.\n", *rootfsName)
}

//DeleteRootfsConfig : Delete RootFs configs
func DeleteRootfsConfig(name *string) {
	path := filepath.Join(RootFSConfigPath, *name)
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		fmt.Printf("Configuration %s does not exist", *name)
		os.Exit(0)
	}

	registryPainc(os.RemoveAll(path))
	fmt.Printf("deleted %s\n", path)
}
