#!/usr/bin/env bash
#
# Build pdfcombiner as a linux(amd64) package.
#
# Dependency:
#   - go source tree:
#     - `hg clone https://code.google.com/p/go`
#     - `export GOROOT=$PWD/go`

REQUIRED_HEADER="include/plan9/amd64/u.h"

verify_goroot()
{
  if [[ ! -f $GOROOT/$REQUIRED_HEADER ]]; then
    echo "Go source needs to be set to \$GOROOT and \$GOROOT/$REQUIRED_HEADER should exist"
    exit 1
  fi
}

build_go()
{
  OLDPWD=$PWD
  cd $(go env GOROOT)/src
  GOOS=linux GOARCH=amd64 ./make.bash --no-clean 2>&1
  cd $OLDPWD
}

build_pdfcombiner()
{
  GOOS=linux GOARCH=amd64 go build
}

verify_goroot

if ! build_go; then
  echo "failed to build go" && exit 1
fi

if ! build_pdfcombiner; then
  echo "crosscompiling pdfcombiner failed" && exit 1
fi

if ! file ./pdfcombiner | grep -q "ELF 64-bit LSB executable"; then
  echo "Something went wrong - the package built for the wrong arch, bailing!" && exit 1
fi

echo "Done crosscompiling pdfcombiner for linux-amd64 at ./pdfcombiner"
