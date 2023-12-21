#!/bin/bash

golangci-lint run --config ./.github/.golangci.yml && \
golangci-lint run --disable-all --config ./.github/.golangci.yml \
-E bodyclose \
-E exhaustive \
-E forbidigo \
-E forcetypeassert \
-E noctx \
-E whitespace \
-E gocognit \
-E unconvert \
-E gomnd \
-E errorlint \
-E gocyclo \
-E goimports \
-E gofmt \
-E errcheck \
-E gosec && \
pkgforbid -config=./.github/pkgforbid.yml ./...
