#!/bin/bash

DIST=$(. /etc/os-release; echo $ID$VERSION_ID)
curl -s -L https://nvidia.github.io/libnvidia-container/gpgkey | \
  sudo apt-key add -
curl -s -L https://nvidia.github.io/libnvidia-container/$DIST/libnvidia-container.list | \
  sudo tee /etc/apt/sources.list.d/libnvidia-container.list
sudo apt-get update


curl -s -L https://nvidia.github.io/libnvidia-container/gpgkey | \
  sudo apt-key add -


sudo apt install libnvidia-container1 libnvidia-container-tools