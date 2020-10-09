package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

//CGroupType Type definition to parse cgroup data
type CGroupType struct {
	memorylimit  int
	cpucorelimit float32
}

var cgroupPath string = "/sys/fs/cgroup"
var cfsPeriodus int = 10000

func cgroupMust(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

func handleErrors(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

func writeToFile(filePath string, file string, data []byte, perm os.FileMode) {
	_, err := os.Stat(filePath)

	if os.IsNotExist(err) {
		fmt.Println("File not found " + filePath)
		os.Exit(0)
	}

	//write to the file
	filePath = filepath.Join(filePath, file)
	cgroupMust(ioutil.WriteFile(filePath, data, perm))
}

func setupPathForCategory(containerName string, catrgory string) string {
	rootCgroupMemoryDir := filepath.Join(cgroupPath, "memory")
	rootCgroupMemoryDir = filepath.Join(rootCgroupMemoryDir, containerName)

	os.Mkdir(rootCgroupMemoryDir, 0755)
	return rootCgroupMemoryDir
}

func createMemoryGroup(containerName string, maxBytes int, pid int) {
	maxBytesString := strconv.Itoa(maxBytes)
	pidString := strconv.Itoa(pid)

	//set-up memory root
	rootCgroupMemoryDir := setupPathForCategory(containerName, "memory")

	//write maxmemory
	writeToFile(rootCgroupMemoryDir, "memory.limit_in_bytes", []byte(maxBytesString), 0700)
	writeToFile(rootCgroupMemoryDir, "notify_on_release", []byte("1"), 0700)

	//create a PID to apply common control to all child in this group
	writeToFile(rootCgroupMemoryDir, "cgroup.procs", []byte(pidString), 0700)
}

func createCPUGroup(containerName string, maxCores float32, pid int) {
	cfsQuotaus := maxCores * float32(cfsPeriodus)
	cfsQuotaString := fmt.Sprintf("%f", cfsQuotaus)
	pidString := strconv.Itoa(pid)

	cfsPeriodString := strconv.Itoa(cfsPeriodus)

	rootCgroupCPUDir := setupPathForCategory(containerName, "cpu")

	//write cpu quota periods
	writeToFile(rootCgroupCPUDir, "cpu.cfs_period_us", []byte(cfsPeriodString), 0700)
	writeToFile(rootCgroupCPUDir, "cpu.cfs_quota_us", []byte(cfsQuotaString), 0700)

	//create a PID to apply common control to all child in this group
	writeToFile(rootCgroupCPUDir, "cgroup.procs", []byte(pidString), 0700)
}

//InitCGroups initialize cgroups
func InitCGroups(cgroupData *CGroupType, containerName string, pid int) {
	fmt.Println("Settigng up cgroups..")
	fmt.Printf("memorylimit=%d cpucorelimit=%f\n", cgroupData.memorylimit, cgroupData.cpucorelimit)

	//setup maxmemory limit:
	if cgroupData.memorylimit != 0 {
		createMemoryGroup(containerName, cgroupData.memorylimit, pid)
	}

	if cgroupData.cpucorelimit != 0 {
		createCPUGroup(containerName, cgroupData.cpucorelimit, pid)
	}

	fmt.Println("Finished setting-up cgroups")
}
