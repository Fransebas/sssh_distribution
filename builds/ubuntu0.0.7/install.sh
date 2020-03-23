#!/usr/bin/env bash

# needs to be root user
if [[ $(id -u) -ne 0 ]] ; then echo "Please run as root" ; exit 1 ; fi



cp ./sssh_server /usr/local/bin/sssh_server
cp ./ssshserver //etc/init.d/ssshserver

echo "Installation finished"