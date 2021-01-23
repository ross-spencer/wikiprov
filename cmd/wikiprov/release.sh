#!/usr/bin/env bash
set -eux

WPROV="wikiprov"
DIR="release"
mkdir -p "$DIR"
export GOOS=windows
export GOARCH=386
go build
mv "$WPROV".exe "${DIR}/${WPROV}"-win386.exe
export GOOS=windows
export GOARCH=amd64
go build
mv "$WPROV".exe "${DIR}/${WPROV}"-win64.exe
export GOOS=linux
export GOARCH=amd64
go build
mv "$WPROV" "${DIR}/${WPROV}"-linux64
export GOOS=darwin
export GOARCH=386
go build
mv "$WPROV" "${DIR}/${WPROV}"-darwin386
export GOOS=darwin
export GOARCH=amd64
go build
mv "$WPROV" "${DIR}/${WPROV}"-darwinAmd64
export GOOS=
export GOARCH=
