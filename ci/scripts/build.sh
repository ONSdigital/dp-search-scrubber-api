#!/bin/bash -eux

cwd=$(pwd)

export GOPATH=$cwd/go

pushd dp-search-scrubber-api
  make build-bin && mv build/$(go env GOOS)-$(go env GOARCH)/* $cwd/build
  cp Dockerfile.concourse $cwd/build
popd
