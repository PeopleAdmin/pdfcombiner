# Common functions for deploying and building

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

