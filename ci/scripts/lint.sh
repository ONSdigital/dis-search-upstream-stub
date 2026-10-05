#!/bin/bash -eux

pushd dis-search-upstream-stub
  make lint
  make validate-specification
popd
