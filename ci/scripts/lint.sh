#!/bin/bash -eux

go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.5
npm install -g @redocly/cli

pushd dp-search-scrubber-api
  make lint
popd
