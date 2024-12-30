#!/bin/sh


set -xe

BIN_DIR=$(dirname $0)
BIN_DIR=$(
    cd $BIN_DIR
    pwd
)

cd $BIN_DIR
mkdir -p log
ls

./go-counter