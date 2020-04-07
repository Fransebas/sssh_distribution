#!/usr/bin/env bash

# $1 is the version

mkdir builds/macos$1

cp -r builds/macos/* builds/macos$1/

go build -o "builds/macos$1/sssh_server"

##

#!/usr/bin/env bash

if [[ $# -eq 1 ]]
  then
    echo "No name supplied"
    exit
fi

export VERSION=${1}
export FOLDER=${2}


mkdir -p ${FOLDER}/usr/local/bin
mkdir -p ${FOLDER}/Library/LaunchDaemons
mkdir -p ${FOLDER}/etc

go build -o ${FOLDER}/usr/local/bin/sssh_server
cp builds/macos/sssh_server.script.sh ${FOLDER}/usr/local/bin/sssh_server.script.sh
cp builds/macos/com.ssshserver.app.plist ${FOLDER}/Library/LaunchDaemons/com.ssshserver.app.plist
cp builds/macos/sssh.conf ${FOLDER}/etc/sssh.conf


pkgbuild --root my_root --identifier my.fake.pkg my_package.pkg