#!/bin/bash -eux


cwd=$(pwd)

export GOPATH=$cwd/go

pushd dp-nlp-search-scrubber
  make build-bin
  cp build/dp-nlp-search-scrubber Dockerfile.concourse $cwd/build
popd

