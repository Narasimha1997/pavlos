#!/bin/bash

export GOPATH=${GOPATH}:$(pwd)

if [[ "$1" == "pavlos" ]]; then

    GOBIN=$(pwd)/bin

    go get github.com/Narasimha1997/pavlos
    go install github.com/Narasimha1997/pavlos

    echo "Installing pavlos under /usr/local/bin , enter sudo password if asked"
    sudo cp $GOBIN/pavlos /usr/local/bin

    echo "Installed pavlos at /usr/local/bin"
fi

if [[ "$1" == "pavlospkg" ]]; then
    GOBIN=$(pwd)/bin

    go get github.com/Narasimha1997/pavlospkg
    go install github.com/Narasimha1997/pavlos

    echo "Installing pavlospkg under /usr/local/bin, enter sudo password if asked"
    sudo cp $GOBIN/pavlospkg /usr/local/bin

    echo "Installed pavlospkg at /usr/local/bin"
fi