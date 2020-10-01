package core

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//RootFSRoot : Root-filesystem root path
var RootFSRoot string = filepath.Join(os.Getenv("HOME"), ".rootfs")

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
func DeleteRootFs(name string) {
	path := filepath.Join(RootFSRoot, name)

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		fmt.Printf("Registry %s not exists, so skipped deletion\n", name)
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
func CreateRootFs() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter the name of root-fs image")
	rootfsName, err := reader.ReadString('\n')
	registryPainc(err)

	rootfsName = strings.Replace(rootfsName, "\n", "", 1)

	fmt.Println("Enter the URL of root-fs image")
	rootfsURL, err := reader.ReadString('\n')
	registryPainc(err)

	rootfsURL = strings.Replace(rootfsURL, "\n", "", 1)

	setupRootFsImage(rootfsName, rootfsURL)
}

//SetupDefault : Sets-up default root-fs image
func SetupDefault() {
	setupRootFsImage("default", "http://cdimage.ubuntu.com/ubuntu-base/releases/18.04/release/ubuntu-base-18.04-base-amd64.tar.gz")
}

//SetupDefaultDir Creates .rootfs at home if not exist
func SetupDefaultDir() {
	_, err := os.Stat(RootFSRoot)
	if os.IsNotExist(err) {
		os.Mkdir(RootFSRoot, 0700)
	}
}
