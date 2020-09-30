# Pavlos
A light-weight container runtime for Linux with NVIDIA gpu support, allows developers to quicky setup development environments for dev and test. Pavlos can emulate any Linux rootfs image as a container.
Pavlos is a greek word which means "small"

### Building Pavlos from srouce :
Follow the steps to build Pavlos from source :
Tools required : golang, git , a working NVIDIA driver and libnvidia-container (for GPU support)

1. #### Download the source
```
git clone git@github.com:Narasimha1997/pavlos.git
```

2. #### Build and install in GOPATH

a. Set `$GOPATH`
``` 
export GOPATH=${GOPATH}:$(pwd)
```
b. Install pavlos:
```
go install pavlos
```

Or you can simply run
```
./build.sh pavlos
```
This script will build pavlos and install it under `/usr/local/bin`

This will build and install pavlosc ( Container cli ) in /usr/bin

3. #### (Optional) Add NVIDIA support :
If you have a working NVIDIA graphics card and want to use it with Pavlos container, you must install libnvida-container. Run `add-nvidia.sh` to install nvidia support. The installation script is for debian based distributions. For RHEL , follow libnvidia-container install instructions.

4. #### (Optional) Download a test rootfs image of Ubuntu 18.04 base :
You can run `get-rootfs.sh` to download Ubuntu 18.04 base filesystem and install it at `$HOME/rootfs`, you can also manually set-up a file-system image as `get-rootfs` is just a collection of simple download and extract commands.

### Create and run the container : 
Pavlos accepts a configuration file of the container image to create. This config file has the format shown below : 

```json
{
    "name" : "TestContainer",
    "enableIsolation" : true,
    "enableResourceIsolation" : true,
    "isolationOptions" : {
        "enableUTS" : true,
        "enablePID" : true,
        "enableRoot" : true,
        "enableNetNs" : false
    },
    "rootFs" : "/path/to/root/fs/absloute/path",
    "nvidiaGpus" : [0],
    "ip" : "172.16.0.12",
    "runtimeArgs" : ["/bin/bash"]
}
```
Keys are self-explanatory, specifically `nvidiaGpus` is an array that takes device ids. `runtimeArgs` is the entrypoint. `isolationOptions` are resources to isolate (Maintain default as true).

`rootfs` points to rootfs absloute path, for example if you had downloaded a rootfs Linux image and placed it in $HOME/rootfs, then `"rootfs" : "/home/username/rootfs`

##### Running a pavlos container
To run use the following command : (A root user or a sudo capable user is required)

```
sudo pavlosc run test-container-config-file.json
```

For example to run test ubuntu container image presented in the repo :
```
sudo pavlosc run ubuntu-container.json
```

For help information : 
`
sudo pavlosc help show
`

If any nvidia devices are specified, the pavlos container runtime will automatically find the PCI devices and drivers, later it uses libnvidia-container to provide isolated gpu to container.

Pavlos automatically copies your host's DNS information (`/etc/resolv.conf`) inside the container, so internet is accessable inside the container.

### Roadmap : 
1. SW defined container overaly networking (like CNI) (This functionality is broken as of now)
2. Resource isolation using Linux - cgroups (Partial support is already implemented)
3. Union File system like the one in Docker (Reduces storage space by reusing existing fs components)
4. More advanced hardware discovery by traversing PCI Device tree (Partial functionality is already implemented to discover NVIDIA GPUs using jaypipe's [GHW](https://github.com/jaypipes/ghw) library)
5. Encrypted file-system support (see dm_crypt and cryptmount) using Linux Kernel device-mapper driver.
6. Volume mounting support like Docker - can mount encrypted volumes as well (See LVM encrypted volume mounting)


