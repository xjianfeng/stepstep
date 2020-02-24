#!/bin/bash

curPath=`pwd`
GOPATH="${HOME}/.golang/path" 

echo $curPath
echo $GOPATH

export GOPATH="${GOPATH}:${curPath}"
echo $GOPATH


binSrcPath=${curPath}"/packages"
if [ ! -d $binSrcPath ] 
then
    ln -s ~/.golang/path/src packages
fi

echo "init env done ."
