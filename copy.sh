#!/usr/bin/env bash
#
# Deploy the compiled binary to the S3 deploy bucket.
#
# Dependencies:
#   - DEPLOY_BUCKET environment variable (defaults to pdfcombiner-deploy)
#   - aws-cli
#     - pip install awscli
#     - export AWS_ACCESS_KEY_ID=DEPLOY_ID_HERE
#     - export AWS_SECRET_ACCESS_KEY=DEPLOY_KEY_HERE

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
source $DIR/deploy/common.sh

DEPLOY_BUCKET=${DEPLOY_BUCKET:-pdfcombiner-deploy}
DEPLOY_URI="s3://$DEPLOY_BUCKET"
DEPLOY_FILE="$DEPLOY_URI/pdfcombiner"

verify_aws
check_bucket

if [[ ! -f ./pdfcombiner ]]; then
  echo "Nothing to deploy" && exit 1
fi

if ! correct_arch; then
  echo "package to deploy is not built for linux.  Don't deploy this \
unless you are absolutely sure it's the right thing to do!" && exit 1
fi

if previous_binary_exists; then
  backup
fi

deploy
echo "deploy complete!"
