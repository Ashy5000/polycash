#!/bin/bash
# This script is used to build the project

# Build all combinations of os and arch
for os in linux darwin windows; do
  for arch in amd64 386 arm arm64; do
    if [ $os = "darwin" -a $arch = "386" ]; then
      continue
    fi
    if [ $os = "darwin" -a $arch = "arm" ]; then
      continue
    fi
    echo "Building for $os $arch"
    env GOOS=$os GOARCH=$arch go build -o builds/node/node_$os-$arch
  done
done