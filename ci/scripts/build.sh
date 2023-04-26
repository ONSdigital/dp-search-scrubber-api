#!/bin/bash -eux


cwd=$(pwd)

export GOPATH=$cwd/go

pushd dp-nlp-search-scrubber
  make build-bin && mv build/$(go env GOOS)-$(go env GOARCH)/*
  cp build/dp-nlp-search-scrubber Dockerfile.concourse $cwd/build
popd

