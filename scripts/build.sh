#!/bin/bash

#
# https://github.com/golang/go/issues/61229
#


dest=./dest
rm -rf $dest
mkdir -p $dest

go build -ldflags=-linkmode=internal -o $dest/hd-wallets cmd/hd-wallets/main.go
GOARCH=amd64 GOOS=windows go build -o $dest/hd-wallets.exe cmd/hd-wallet/smain.go

