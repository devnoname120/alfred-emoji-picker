#!/usr/bin/env bash

echo -e "Note: you need to run ./gen-emojis.sh if you want to regenerate the emoji images\n"

set -exuo pipefail


GOOS=darwin GOARCH=amd64 go build -o main_amd64 main.go
GOOS=darwin GOARCH=arm64 go build -o main_arm64 main.go

lipo -create -output alfred-emoji-picker main_amd64 main_arm64
