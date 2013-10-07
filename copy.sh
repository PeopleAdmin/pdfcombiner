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

DEPLOY_BUCKET=${DEPLOY_BUCKET:-pdfcombiner-deploy}
DEPLOY_URI="s3://$DEPLOY_BUCKET"
DEPLOY_FILE="$DEPLOY_URI/pdfcombiner"

verify_aws() {
  if [[ -z "$AWS_ACCESS_KEY_ID" ]] || [[ -z "$AWS_SECRET_ACCESS_KEY" ]]; then
    echo "\$AWS_ACCESS_KEY_ID and \$AWS_SECRET_ACCESS_KEY must be set to the deploy keys"
    exit 1
  fi

  if ! which -s aws; then
    echo "aws must be installed and in the PATH, try \`pip install awscli\`"
    exit 1
  fi
}

check_bucket() {
  bucket_contents=$( aws --region us-east-1 s3 ls $DEPLOY_URI/ 2>&1 )
  if [[ "$?" != 0 ]]; then
    echo "Couldn't access $DEPLOY_URI: $bucket_contents"
    exit 1
  fi

  if ! echo "$bucket_contents" | grep -q 'cpdf$'; then
    echo "cpdf must be present in deploy bucket"
    exit 1
  fi
  export bucket_contents
}

correct_arch() {
  file ./pdfcombiner | grep -q "ELF 64-bit LSB executable"
}

previous_binary_exists() {
  echo "$bucket_contents" | grep -q 'pdfcombiner$'
}

backup() {
  now=`date +"%Y-%m-%d_%H.%M.%S"`
  backup_loc="$DEPLOY_URI/backups/pdfcombiner.$now"
  echo "previous copy found at $DEPLOY_FILE, making backup"
  aws --region us-east-1 s3 cp $DEPLOY_FILE $backup_loc
  if [[ $? != 0 ]]; then
    echo "problem backing up!" && exit 1
  fi
}

deploy() {
  echo "Deploying to $DEPLOY_FILE..."
  aws --region us-east-1 s3 cp ./pdfcombiner $DEPLOY_FILE
  if [[ $? != 0 ]]; then
    echo "problem deploying!" && exit 1
  fi
}

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
