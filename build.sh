#!/usr/bin/env bash
# Build pdfcombiner as a linux(amd64) package and copy it to an S3 bucket.
# dependencies:
#   - golang-crosscompile:
#     - `git clone git://github.com/davecheney/golang-crosscompile.git`
#     - `export -f go-crosscompile-build
#   - go source tree:
#     - `hg clone https://code.google.com/p/go`
#     - `export GOROOT=$PWD/go`
#

REQUIRED_HEADER="include/plan9/amd64/u.h"

verify_crosscompile()
{
  if [[ `type -t go-crosscompile-build` != "function" ]];then
    echo "Need to \`source golang-crosscompile/crosscompile.bash\` and \`export -f go-crosscompile-build\` "
    exit 1
  fi
}

verify_goroot()
{
  if [[ ! -f $GOROOT/$REQUIRED_HEADER ]]; then
    echo "Go source needs to be set to \$GOROOT and \$GOROOT/$REQUIRED_HEADER should exist"
    exit 1
  fi
}

verify_aws()
{
  if [[ -z "$AWS_ACCESS_KEY_ID" ]] || [[ -z "$AWS_SECRET_ACCESS_KEY" ]]; then
    echo "\$AWS_ACCESS_KEY_ID and \$AWS_SECRET_ACCESS_KEY should be set and have access to deploy bucket"
    exit 1
  fi
}

verify_goroot
verify_aws
verify_crosscompile
OLDPWD=$PWD

if ! go-crosscompile-build linux/amd64; then
  "building go failed" && exit 1
fi

cd $OLDPWD
if ! GOOS=linux GOARCH=amd64 go build; then
  echo "building pdfcombiner failed" && exit 1
fi

if ! file ./pdfcombiner | grep -q "ELF 64-bit LSB executable"; then
  echo "Something went wrong, bailing!" && exit 1
fi

echo "Done building pdfcombiner for linux-amd64 at ./pdfcombiner"
