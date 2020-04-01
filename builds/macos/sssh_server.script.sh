#!/usr/bin/env bash


# The script is needed because in OSX there is a security thing that when running the program directly
# we get the following error fork/exec /bin/bash: operation not permitted inside of the library pty
# With this it doesn't happen

sudo /usr/local/bin/sssh_server 2> /tmp/com.ssshserver.app.err || exit 1