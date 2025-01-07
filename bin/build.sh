#!/bin/sh

set -xe

BIN_DIR=$(dirname $0)
BIN_DIR=$(
    cd $BIN_DIR
    pwd
)
ROOT_DIR=$(
    cd $BIN_DIR/..
    pwd
)


mkdir -p $ROOT_DIR/dist/log
mkdir -p $ROOT_DIR/dist/config
mkdir -p $ROOT_DIR/dist/bin

go build $ROOT_DIR/main/go-counter.go

mv go-counter $ROOT_DIR/dist/bin
cp $ROOT_DIR/config/*.json $ROOT_DIR/dist/config
cp $ROOT_DIR/bin/* $ROOT_DIR/dist/bin

echo "build success"