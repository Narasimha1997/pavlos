package core

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

//ChildNamespaceName : Child process name that creates container
const ChildNamespaceName = "container"

//MakeCloneFlags : Creates a set clone flag from UNIX SYSCALL options
func MakeCloneFlags(config *ContainerOpts) uintptr {
	if config.EnableIsolation {
		var flag uintptr = 0
		if config.IsolationOpts.EnableUTS {
			fmt.Println("Unix time sharing enabled")
			flag = syscall.CLONE_NEWUTS
		}

		if config.IsolationOpts.EnablePID {
			fmt.Println("Process clones enabled")
			flag = flag | syscall.CLONE_NEWPID
		}

		if config.IsolationOpts.EnableNetNs {
			fmt.Println("Net network namespace enabled")
			flag = flag | syscall.CLONE_NEWNET
		}

		flag = flag | syscall.CLONE_NEWNS
		return flag
	}
	return 0
}

//MakeContainerRuntime : Creates a exec interface with options specified
func MakeContainerRuntime(config *ContainerOpts, konArgs []string) *exec.Cmd {
	command := exec.Command("/proc/self/exe", append([]string{ChildNamespaceName}, konArgs...)...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	sysFlags := MakeCloneFlags(config)
	command.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:   sysFlags,
		Unshareflags: syscall.CLONE_NEWNS, //Unshare the namespace
	}

	fmt.Println(command.String())

	return command
}
