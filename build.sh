#!/bin/bash

#get the package for PCI probe
go get github.com/jaypipes/ghw

BASEDIR=${PWD}

cd ${BASEDIR}/pavlos
go build -o pavlosc
mv ${PWD}/pavlosc ${BASEDIR}/pavlosc
cd ${BASEDIR}

#attempt to install (if run as root, it will be automatically installed)
echo "Now the script will install pavlosc in /usr/bin, sudo to continue, Ctrl +C to skip"
sudo cp pavlosc /usr/bin/

echo "Done"