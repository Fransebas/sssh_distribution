#!/usr/bin/env bash
export GOOS="linux"
export GOARCH="arm"
export CGO_ENABLED=1

go build -o builds/ubuntu/sssh_server