#!/bin/bash -eux

pushd dis-search-upstream-stub
  make build
  cp build/dis-search-upstream-stub Dockerfile.concourse ../build
popd
