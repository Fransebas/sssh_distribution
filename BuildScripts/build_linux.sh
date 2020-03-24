#!/usr/bin/env bash

# $1 is the version
mkdir "builds/ubuntu$1"

cp -R "builds/ubuntu/" "builds/ubuntu$1/"

go build -o "builds/ubuntu$1/sssh_server"
