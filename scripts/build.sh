#!/bin/bash

#
# https://github.com/golang/go/issues/61229
#


dest=./dest
rm -rf $dest
mkdir -p $dest

go build -ldflags=-linkmode=internal -o $dest/hd-wallet cmd/hd-wallet/main.go
GOARCH=amd64 GOOS=windows go build -ldflags=-linkmode=internal -o $dest/hd-wallet.exe cmd/hd-wallet/main.go

