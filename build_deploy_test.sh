#!/bin/sh

set -e

export VERSION=$(
  date +%m%d%H%M
)

export GOPATH=`pwd`
export GOBIN=${GOPATH}/bin
export PATH=${PATH}:${GOPATH}/bin
export DONT_GLOB_JUST_YET='*'

echo "getting golint..."
go get -v github.com/golang/lint/golint
echo "golint should be available now"

for PACKAGE in $(ls -1 src | grep _user_dob ) ; do

    echo "Building ${PACKAGE}..."
    bash -c "\
        set -e; \
        rm -f bin/${PACKAGE}.zip; \
        cd src/${PACKAGE}; \
        go get -v -t ./...; \
        golint -set_exit_status; \
        go vet . ; \
        go test . ; \
        go build -o main ./${DONT_GLOB_JUST_YET}.go ;
        zip -9 ../../bin/${PACKAGE}.zip main \
    "
    echo "Successfully built ${PACKAGE}"

done

echo "Trying to get hold of your AWS Account ID"
echo "...if this explodes your aws cli install is missing or misconfigured"
export ACCOUNT_ID=$(
    aws ec2 describe-security-groups \
        --group-names 'Default' \
        --query 'SecurityGroups[0].OwnerId' \
        --output text \
)
echo "Got ACCOUNT_ID = ${ACCOUNT_ID}"

echo "To deploy bigger stack we have to bootstrap smaller stack..."
echo "...this will create a temp bucket for cloudformation templates and lambda code"
echo "...located at s3://${ACCOUNT_ID}-cloudformation-temp"
aws cloudformation deploy \
    --template-file cloudform-bootstrap.yaml \
    --stack-name cloudform-bootstrap \
    --no-fail-on-empty-changeset

aws s3 sync \
    --exclude '*' \
    --include '*_user_dob.zip' \
    bin/ \
    s3://${ACCOUNT_ID}-cloudformation-temp/lambda_handlers/$VERSION/

echo "Now deploying stack user-dob-rest-service..."

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
curl --show-error --fail ${CURL_VERBOSITY} ${API_URL}/hello/${USER_NAME}

echo "

all non-ignored curl commands completed with zero exit statuses"
