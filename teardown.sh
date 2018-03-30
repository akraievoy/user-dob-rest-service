#!/bin/sh

set -e

echo "trying to delete user-dob-rest-service stack"
aws cloudformation delete-stack \
  --stack-name user-dob-rest-service
echo "...Succeeded"

echo "trying to delete cloudform-bootstrap stack"
aws cloudformation delete-stack \
  --stack-name cloudform-bootstrap
echo "...Succeeded"

echo "Trying to get hold of your AWS Account ID"
export ACCOUNT_ID=$(
    aws ec2 describe-security-groups \
        --group-names 'Default' \
        --query 'SecurityGroups[0].OwnerId' \
        --output text \
)
BUCKET_TO_REMOVE=s3://${ACCOUNT_ID}-cloudformation-temp

echo "
The bootstrap bucket still survives, and there's no way to force its deletion through CloudFormation
Let's try to forcefully remove ${BUCKET_TO_REMOVE}.

Bail with Ctrl+C if you'd like to retain it -- sleeping for 10 seconds...
" 1>&2
sleep 10

aws s3 rb --force "${BUCKET_TO_REMOVE}"
