#!/usr/bin/env bash

# $1 is the version
mkdir "builds/linux$1"

cp -r builds/ubuntu/* "builds/linux$1/"

go build -o "builds/linux$1/sssh_server"
