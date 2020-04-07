#!/usr/bin/env bash

if [[ $# -eq 1 ]]
  then
    echo "No name supplied"
    exit
fi

export VERSION=${1}
export FOLDER=${2}


rm -r ${FOLDER}

mkdir -p ${FOLDER}/usr/local/bin
mkdir -p ${FOLDER}/etc/init.d
mkdir -p ${FOLDER}/lib/systemd/system
mkdir -p ${FOLDER}/etc



go build -o ${FOLDER}/usr/local/bin/sssh_server
cp builds/ubuntu/ssshserver ${FOLDER}/etc/init.d/ssshserver
cp builds/ubuntu/sssh_server.service ${FOLDER}/lib/systemd/system/sssh_server.service
cp builds/ubuntu/sssh.conf ${FOLDER}/etc/sssh.conf

cp -r builds/ubuntu/DEBIAN ${FOLDER}/DEBIAN

chmod -R 0775 ${FOLDER}/DEBIAN

CWD=`pwd`

cd ${FOLDER}
find . -type f ! -regex '.*.hg.*' ! -regex '.*?debian-binary.*' ! -regex '.*?DEBIAN.*' -printf '%P ' | xargs md5sum > ${FOLDER}/DEBIAN/md5sums

cd ${CWD}

dpkg -b ${FOLDER} ~/sssh-server${VERSION}.deb
