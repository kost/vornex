#!/bin/bash

echo 'Building functional-test binary'
go build

echo 'Building VORNEX binary from current branch'
go build -o vornex_dev ../vornex

echo 'Installing latest release of VORNEX'
GO111MODULE=on go build -v github.com/kost/vornex/v2/cmd/vornex

echo 'Starting VORNEX functional test'
./functional-test -main ./vornex -dev ./vornex_dev -testcases testcases.txt
