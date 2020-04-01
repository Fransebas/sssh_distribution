#!/usr/bin/env bash

# needs to be root user
if [[ $(id -u) -ne 0 ]] ; then echo "Please run as root" ; exit 1 ; fi

cp ./sssh_server /usr/local/bin/sssh_server
cp ./sssh_server.script.sh /usr/local/bin/sssh_server.script.sh
cp ./com.ssshserver.app.plist /Library/LaunchDaemons/com.ssshserver.app.plist
cp ./sssh.conf /etc/sssh.conf

mkdir /etc/sssh

KEYFILE=/etc/sssh/rsa_host

if [[ ! -f ${KEYFILE} ]]; then
    ./sssh_server -mode=keygen -filename=${KEYFILE}
fi


echo "Installation finished"