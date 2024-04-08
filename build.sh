#!/bin/bash
# Copyright 2024, Asher Wrobel

# This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

# This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

# You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

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

# Build the GUI wallet
cd gui_wallet
for triple in x86_64-unknown-linux-gnu; do
  echo "Building for $triple"
  cargo build --release --target $triple
  mv target/$triple/release/gui_wallet ../builds/gui_wallet/gui_wallet_$triple
done
