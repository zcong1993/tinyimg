#!/bin/bash
set -ev
NAME=tinyimg
COMMIT=`git describe --always`

rm -rf bin/
go build \
    -ldflags="-X main.GitCommit=${COMMIT}" \
    ./...

mkdir bin
mv tinyimg ./bin/
