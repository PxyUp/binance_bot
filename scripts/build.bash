#!/usr/bin/env bash

rm -rf bin

version_inject=$1

platforms=("windows/amd64" "windows/386" "darwin/amd64" "linux/arm64" "linux/amd64")

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name=bin/$version_inject-binance_bot'-'$GOOS'-'$GOARCH
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi
    if [ $GOOS = "darwin" ]; then
      continue
    fi
    env GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name -ldflags "-s -w" main.go
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
done