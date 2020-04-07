#!/usr/bin/env bash

if [[ $# -eq 0 ]]
  then
    echo "No name supplied"
    exit
fi

export VERSION=${1}
export OUTPUT=${2}
export FOLDER=temp

mkdir -p ${FOLDER}/usr/local/bin
mkdir -p ${FOLDER}/Library/LaunchDaemons
mkdir -p ${FOLDER}/etc

go build -o ${FOLDER}/usr/local/bin/sssh_server
cp builds/macos/sssh_server.script.sh ${FOLDER}/usr/local/bin/sssh_server.script.sh
cp builds/macos/com.ssshserver.app.plist ${FOLDER}/Library/LaunchDaemons/com.ssshserver.app.plist
cp builds/macos/sssh.conf ${FOLDER}/etc/sssh.conf

chmod -R 755 ${FOLDER}

chmod u+x builds/macos/scripts/*

mkdir -p ${OUTPUT}

pkgbuild  --root ${FOLDER} --scripts builds/macos/scripts  --identifier com.ssshserver.app ${OUTPUT}/sssh_server${VERSION}.pkg

rm -r ${FOLDER}