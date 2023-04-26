#!/bin/bash -eux

pushd dp-nlp-search-scrubber
  make build-bin
  cp build/dp-nlp-search-scrubber Dockerfile.concourse build
popd
