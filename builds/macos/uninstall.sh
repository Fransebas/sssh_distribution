#!/usr/bin/env bash

# needs to be root user
if [[ $(id -u) -ne 0 ]] ; then echo "Please run as root" ; exit 1 ; fi

sudo launchctl stop com.ssshserver.app

rm /usr/local/bin/sssh_server
rm /usr/local/bin/sssh_server.script.sh
rm /Library/LaunchDaemons/com.ssshserver.app.plist
rm /etc/sssh.conf

