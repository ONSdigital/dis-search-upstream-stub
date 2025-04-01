#!/bin/bash -eux

pushd dis-search-upstream-stub
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.5
  make lint
  npm install -g @redocly/cli
  make validate-specification
popd
