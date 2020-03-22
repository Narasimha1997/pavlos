#!/bin/bash

#wget http://cdimage.ubuntu.com/ubuntu-base/releases/18.04/release/ubuntu-base-18.04-base-amd64.tar.gz
mkdir -p $HOME/rootfs
tar -C $HOME/rootfs -xf ubuntu-base-18.04-base-amd64.tar.gz
rm ubuntu-base-18.04-base-amd64.tar.gz

echo "DONE!"