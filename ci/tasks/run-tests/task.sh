#!/bin/bash

set -eux

mkdir -p go/src/code.cloudfoundry.org/uaa-cli
cp -R uaa-cli/* go/src/code.cloudfoundry.org/uaa-cli
export GOPATH="$(pwd)/go"
cd go/src/code.cloudfoundry.org/uaa-cli
make
