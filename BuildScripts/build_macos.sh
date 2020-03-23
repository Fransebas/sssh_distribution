#!/usr/bin/env bash

# $1 is the version

mkdir macos$1

go build -o "builds/macos$1/sssh_server"