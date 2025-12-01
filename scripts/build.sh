#!/bin/bash

#
# https://github.com/golang/go/issues/61229
#


dest=./dest
rm -rf $dest
mkdir -p $dest

go build -ldflags=-linkmode=internal -o $dest/hd-wallets.darwin.macchip cmd/hd-wallets/main.go


build_single() {
  echo "build ${1}:${2} ${3}"
  GOOS=${1} GOARCH=${2} go build -ldflags "-s -w" -o "${3}" cmd/hd-wallets/main.go
  if [[ ${1} != "darwin" ]]; then
    upx --brute  "${3}"
  fi
}

oa_es=(linux:amd64 windows:amd64:.exe darwin:amd64 darwin:arm64)

for oa in "${oa_es[@]}"
do
  # shellcheck disable=SC2206
  oa_s=(${oa//:/ })
  build_single "${oa_s[0]}" "${oa_s[1]}" "${dest}/hd-wallets_${oa_s[0]}_${oa_s[1]}${oa_s[2]}"
done
