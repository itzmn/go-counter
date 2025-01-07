#!/bin/sh


set -xe

BIN_DIR=$(dirname $0)
BIN_DIR=$(
    cd $BIN_DIR
    pwd
)

ROOT_DIR=$(cd $BIN_DIR/..; pwd)

ls
cd $BIN_DIR

ls
exec ./go-counter -config=../config/docker-config.json -variablesPath=../config/statisticVars.json -log=../log