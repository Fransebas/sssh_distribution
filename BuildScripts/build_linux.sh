#!/usr/bin/env bash

# $1 is the version
mkdir "builds/ubuntu$1"

cp "builds/ubuntu/install.sh" "builds/ubuntu$1/install.sh"
cp "builds/ubuntu/ssshserver" "builds/ubuntu$1/ssshserver"

go build -o "builds/ubuntu$1/sssh_server"
