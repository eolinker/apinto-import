#!/usr/bin/env bash

echo $0
. $(dirname $0)/common.sh

#echo ${BasePath}
#echo ${CMD}
#echo ${Hour}

VERSION=$(genVersion $1)
OUTPATH="${BasePath}/out/apinto-import-${VERSION}"
buildApp apinto-import $VERSION

cp -a ${BasePath}/build/resources/*  ${OUTPATH}/

cd ${ORGPATH}
