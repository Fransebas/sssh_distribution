#!/usr/bin/env bash

# needs to be root user
if [[ $(id -u) -ne 0 ]] ; then echo "Please run as root" ; exit 1 ; fi

sudo launchctl stop com.ssshserver.app

rm /usr/local/bin/sssh_server
rm /etc/init.d/ssshserver
rm /lib/systemd/system/sssh_server.service
rm /etc/sssh.conf
