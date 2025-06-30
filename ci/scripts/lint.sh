#!/bin/bash -eux

pushd dis-search-upstream-stub
  go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.6
  make lint
  npm install -g @redocly/cli
  make validate-specification
popd
