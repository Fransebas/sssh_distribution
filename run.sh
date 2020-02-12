#!/usr/bin/env bash
set START_INSTALL=`date +%s`
go install
set END_INSTALL=`date +%s`
DIFF=`${START_INSTALL} - ${END_INSTALL}`
echo "installation time ${DIFF}"
sudo sssh_server &> err.txt