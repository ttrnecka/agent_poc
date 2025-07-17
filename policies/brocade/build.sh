#!/usr/bin/env bash
set -x
# STEP 1: Determinate the required values

VERSION=$1
echo $VERSION

# STEP 2: Build the ldflags

LDFLAGS=(
  "-X 'main.Version=${VERSION}'"
)

# STEP 3: Actual Go build process

go build -ldflags="${LDFLAGS[*]}" -o "brocade_cli_${VERSION}.exe"