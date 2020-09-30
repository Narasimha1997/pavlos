#!/bin/bash

export GOPATH=${GOPATH}:$(pwd)

if [[ "$1" == "pavlos" ]]; then

    GOBIN=$(pwd)/bin

    go get github.com/Narasimha1997/pavlos
    go install github.com/Narasimha1997/pavlos

    echo "Installing pavlos under /usr/local/bin , enter sudo password"
    sudo cp $GOBIN/pavlos /usr/local/bin
fi

if [[ "$1" == "pavlospkg" ]]; then
    echo "To be implemented"
fi