#!/bin/sh


set -xe

BIN_DIR=$(dirname $0)
BIN_DIR=$(
    cd $BIN_DIR
    pwd
)

ROOT_DIR=$(cd $BIN_DIR/..; pwd)

ls

${ROOT_DIR}/bin/go-counter -config=../config/config.json -variablesPath=../config/statisticVars.json