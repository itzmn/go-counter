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

go build $ROOT_DIR/main/go-counter.go

mv go-counter $ROOT_DIR/dist
cp $ROOT_DIR/config/config.json $ROOT_DIR/dist
cp $ROOT_DIR/bin/* $ROOT_DIR/dist

echo "build success"