#!/usr/bin/env bash
export GOOS="linux"
export GOARCH="arm"
export CGO_ENABLED=1


# $1 is the version

mkdir builds/ubuntu$1

cp "builds/ubuntu/install.sh" builds/ubuntu$1/install.sh
cp "builds/ubuntu/ssshserver.sh" builds/ubuntu$1/ssshserver.sh

go build -o "builds/ubuntu$1/sssh_server"