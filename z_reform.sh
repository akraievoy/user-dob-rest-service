#!/bin/bash

set -e

export ACCOUNT_ID=FILLME
export VERSION=03301644

echo "
This script updates CloudFormation Stack without building/uploading Lambda binaries again.
ACCOUNT_ID=${ACCOUNT_ID}
VERSION=${VERSION}
Bail with Ctrl+C if you are not sure you have manually updated script vars before running this.
Sleeping for 3 seconds...
" 1>&2
sleep 3

if aws cloudformation list-stacks --output text | grep -v DELETE_COMPLETE | grep user-dob-rest-service | grep ' ROLLBACK_COMPLETE' > /dev/null ; then
    echo deleting rolled-back stack user-dob-rest-service
    aws cloudformation delete-stack \
        --stack-name user-dob-rest-service

    while aws cloudformation list-stacks --output text | grep -v DELETE_COMPLETE | grep user-dob-rest-service > /dev/null; do
        echo awaiting deletion of stack
        sleep 2
    done
fi

echo "Now deploying bigger stack..."
aws cloudformation deploy \
    --capabilities CAPABILITY_IAM \
    --template-file user-dob-rest-service.yaml \
    --stack-name user-dob-rest-service \
    --parameter-override CodeVersion=${VERSION} \
    `# that's the bucket we have just created via preceding deploy` \
    --s3-bucket ${ACCOUNT_ID}-cloudformation-temp \
    --force-upload \
    --no-fail-on-empty-changeset

aws cloudformation describe-stack-events \
    --stack-name user-dob-rest-service \
    --output table \
    --query 'StackEvents[*].{Status:ResourceStatus,Type:ResourceType,Reason:ResourceStatusReason}' | grep -v _IN_PROGRESS

aws cloudformation describe-stacks \
    --stack-name user-dob-rest-service \
    --output table \
    --query 'Stacks[0].{Outputs:Outputs}'

API_URL=$( \
    aws cloudformation describe-stacks \
        --stack-name user-dob-rest-service \
        --output text \
        --query 'Stacks[0].Outputs[?OutputKey==`ApiGatewayInternetUrl`].OutputValue'
)
USER_NAME="$( whoami )$( date +%m%d%H%M )"
CURL_VERBOSITY="" #set -i or -vvv if needed

echo "
--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--"
echo "getting resource which is not registered..."
echo "--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--"
curl --show-error --fail ${CURL_VERBOSITY} ${API_URL}/hello/${USER_NAME} || \
    echo "this command should fail with 404 and we ignore curl non-zero exit status"

echo "
--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--"
echo "setting a DOB from the past..."
echo "--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--"
curl --show-error --fail ${CURL_VERBOSITY} -X PUT \
    -H "Content-Type: application/json" \
    -d '{"dateOfBirth": "1982-09-18"}' \
    ${API_URL}/hello/${USER_NAME}

echo "
--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--"
echo "getting days to wait for a DOB as we've just set"
echo "--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--"
curl --show-error --fail ${CURL_VERBOSITY} ${API_URL}/hello/${USER_NAME}

echo "
--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--"
echo "setting a DOB from the near past..."
echo "--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--"
curl --show-error --fail ${CURL_VERBOSITY} -X PUT \
    -H "Content-Type: application/json" \
    -d "{\"dateOfBirth\": \"1982-$(date +%m-%d --date=yesterday)\"}" \
    ${API_URL}/hello/${USER_NAME}

echo "
--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--"
echo "getting 364 days to wait for a DOB as we've just set (could be 365 for leap years)"
echo "--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--"
curl --show-error --fail ${CURL_VERBOSITY} ${API_URL}/hello/${USER_NAME}

echo "
--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--"
echo "setting a DOB from the near future..."
echo "--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--"
curl --show-error --fail ${CURL_VERBOSITY} -X PUT \
    -H "Content-Type: application/json" \
    -d "{\"dateOfBirth\": \"1982-$(date +%m-%d --date=tomorrow)\"}" \
    ${API_URL}/hello/${USER_NAME}

echo "
--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--"
echo "getting one days to wait for a DOB as we've just set"
echo "--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--"
curl --show-error --fail ${CURL_VERBOSITY} ${API_URL}/hello/${USER_NAME}

echo "
--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--"
echo "setting a DOB from distant future..."
echo "--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--"
curl --show-error --fail ${CURL_VERBOSITY} -X PUT \
    -H "Content-Type: application/json" \
    -d '{"dateOfBirth": "2025-09-18"}' \
    ${API_URL}/hello/${USER_NAME}

echo "
--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--"
echo "getting much more days to wait for a DOB in future"
echo "--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--"
curl --show-error --fail ${CURL_VERBOSITY} ${API_URL}hello/${USER_NAME}

echo "

all non-ignored curl commands completed with zero exit statuses

ReST API URL (with username parameter): ${API_URL}hello/${USER_NAME}"
