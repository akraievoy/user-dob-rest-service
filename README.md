# References for used GoLang tooling

* [AWS SDK -- Docs](https://docs.aws.amazon.com/sdk-for-go/api/aws/)
* [AWS SDK -- Sources](https://github.com/aws/aws-sdk-go/tree/master/aws)
* [Proper AWS CodeBuild template](https://github.com/aws-samples/lambda-go-samples/blob/master/buildspec.yml)
* [GoLang Programming Model for AWS Lambda](https://docs.aws.amazon.com/lambda/latest/dg/go-programming-model.html)

## Cloudforming API Gateway for Lambda

* [Use proper Integration.Type](https://docs.aws.amazon.com/apigateway/api-reference/resource/integration/#type)
* [More Docs on Integration](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apitgateway-method-integration.html#cfn-apigateway-method-integration-requesttemplates)
* [Jayway's example - simplest yet still working, *Stage also cloudfromed*](https://blog.jayway.com/2016/08/17/introduction-to-cloudformation-for-api-gateway/)
* [Lambda Versions](https://docs.aws.amazon.com/lambda/latest/dg/versioning-intro.html)
* [tmaslen's example](https://gist.github.com/tmaslen/ddb53c6393fd80160dd2a8c75900cfb3)
* [matsev's example](https://github.com/matsev/cloudformation-api-gateway/blob/master/cloudformation.template)
* [zombie apocalypse example](https://github.com/aws-samples/aws-lambda-zombie-workshop/blob/master/CloudFormation/CreateZombieWorkshop.json)
* [SO Example](https://stackoverflow.com/a/39772260)
* [This is what happens if you shoot GET API requests as GET requests to Lambda Functions](https://stackoverflow.com/questions/41371970/accessdeniedexception-unable-to-determine-service-operation-name-to-be-authoriz)

## Cloudforming DynamoDB

* [Lambda-backed Rest-API sample](https://github.com/jthomerson/cloudformation-template-for-lambda-backed-api-gateway-with-dynamodb/blob/master/cloudformation/complete-api.template) <-- used that as starting point

## CloudFormation

* [Main Landing](https://aws.amazon.com/cloudformation/?p=tile)
* [Intrinsic functions](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference.html?shortFooter=true)
* [cfn_flip](https://github.com/awslabs/aws-cfn-template-flip)
* [stelligent template collection](https://github.com/stelligent/cloudformation_templates)
* [Json to Yaml with CF Designer](https://aws.amazon.com/blogs/mt/the-virtues-of-yaml-cloudformation-and-using-cloudformation-designer-to-convert-json-to-yaml/), but cfn_flip *is* much better

# CI/CD maybes

* [vamp](https://vamp.io/)
* GoCD

## API Gateway <-> Amazon Lambda

* [Dev Guide](https://docs.aws.amazon.com/lambda/latest/dg/lambda-introduction-function.html)
* [AWS Dev Guide, blueprints](https://docs.aws.amazon.com/lambda/latest/dg/get-started-step3-optional.html)
* [AWS Dev Guide, create event source](https://docs.aws.amazon.com/lambda/latest/dg/with-on-demand-https-example-configure-event-source_2.html)
* [Good AWS Sample with API dicsussion](https://aws.amazon.com/blogs/compute/announcing-go-support-for-aws-lambda/)
* [Cloud Guru sample relying on http client](https://read.acloud.guru/serverless-golang-api-with-aws-lambda-34e442385a6a)
* [aws-lambda-go GH samples](https://github.com/aws/aws-lambda-go/blob/master/events/README_ApiGatewayEvent.md)

### Old, not needed for lambda, rewrite/transform magick

* [Integration Pass-Through Behaviours](https://docs.aws.amazon.com/apigateway/latest/developerguide/integration-passthrough-behaviors.html)
* [Integration Request Parameters](https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-swagger-extensions-integration-requestParameters.html)
* [API Gateway Transforms Pictorial](https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-create-api-step-by-step.html) 
* [Same without pictures as a Developer Guide Book](https://www.scribd.com/document/361524947/Apigateway-Developer-guide)
* [There's better integration for Lambda](https://stackoverflow.com/a/46114185/148926)
* [Cloudonaut example -- Swagger, Models, STAAHP](https://cloudonaut.io/create-a-serverless-restful-api-with-api-gateway-cloudformation-lambda-and-dynamodb/)
* [More Swagger Madness](https://stackoverflow.com/questions/47638957/use-timestamp-in-aws-api-gateway-in-the-path-override-of-the-integration-request)
* [Swagger at AwsLabs](https://github.com/awslabs/aws-apigateway-importer/issues/121)
* [That is what happens if you use old style integration](https://stackoverflow.com/questions/36122250/requestparameters-returning-invalid-mapping-expression-specified-true)
* [...and more of the same](https://stackoverflow.com/questions/46250206/aws-api-gateway-http-passthrough-path-parameters/46259271)
* [That is what happens if you use old style integration and/or interpreted language](https://stackoverflow.com/questions/43708017/aws-lambda-api-gateway-error-malformed-lambda-proxy-response)

## Talking to DynamoDB:

* [Update of dynamo from golang lambda](https://hackernoon.com/aws-lambda-go-vs-node-js-performance-benchmark-1c8898341982)
* [updates are upserts](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/WorkingWithItems.html)
* [expressions and value maps](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/GettingStarted.Java.03.html#GettingStarted.Java.03.03)

### Testing against local DynamoDB mock

* [AWS News](https://aws.amazon.com/blogs/aws/dynamodb-local-for-desktop-development/)
* [SO, easier testing](https://stackoverflow.com/questions/26901613/easier-dynamodb-local-testing)
* [GH, Dockerized version](https://github.com/cnadiminti/docker-dynamodb-local)
* [Survey from Serebrov](http://serebrov.github.io/html/2015-02-01-dynamodb-local.html)

### Local Testing for AWS services + SAM 

* [LocalStack](https://blog.utar.co/blog/localstack)
* [SAM Local](https://docs.aws.amazon.com/lambda/latest/dg/test-sam-local.html)
* [SAM + Lambda @ Edge example](https://github.com/awslabs/serverless-application-model/tree/master/examples/2016-10-31/lambda_edge)

## AWS Security Model

* [Service Assume Role with Perms for Resources](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_create_for-service.html?shortFooter=true)
* [IAM Policies](https://docs.aws.amazon.com/IAM/latest/UserGuide/access_policies.html)
* [Identity vs Resource policies](https://docs.aws.amazon.com/IAM/latest/UserGuide/access_policies_identity-vs-resource.html?shortFooter=true)
* [Demystifying Resource-level perms (last section on auth tracing)](https://aws.amazon.com/blogs/security/demystifying-ec2-resource-level-permissions/)
* [The Diagram](https://docs.aws.amazon.com/IAM/latest/UserGuide/images/access-diagram_800.png)
* [Access for Dynamo Table](https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_examples_dynamodb_specific-table.html?shortFooter=true)

### Custom CloudWatch Metrics

* [SO](https://stackoverflow.com/questions/17019069/what-can-i-use-custom-cloudwatch-metrics-for)
* [AWS user Guide](https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/publishingMetrics.html)
* [SO, Python](https://serverfault.com/questions/824510/access-to-custom-cloudwatch-metrics-using-lambdaboto3)

#### Alternatives to Lambda

* [Fargate](https://aws.amazon.com/fargate/?p=tile)
* [Kinesis DataStreams](https://aws.amazon.com/kinesis/data-streams/)
* [Reference: Kinesis -> Lambda -> DynamoDB](https://s3.amazonaws.com/awslambda-reference-architectures/stream-processing/lambda-refarch-streamprocessing.pdf)
  [Source](https://github.com/aws-samples/lambda-refarch-streamprocessing)

#### Sparta

* [Sparta Code Pipeline + go dep](https://github.com/mweagle/SpartaCodePipeline/blob/master/buildspec.yml)
* [Sparta - Behind the Shield, Apr 2016](https://medium.com/@mweagle/sparta-behind-the-shield-7a6e178f1b72)
* [Go Framework for AWS Lambda](https://medium.com/@mweagle/a-go-framework-for-aws-lambda-ab14f0c42cb)
* [Sparta - CodePipelines, Sept 2017](https://medium.com/statuscode/serverless-serverfull-and-weaving-pipelines-c9f83eec9227)
* [GoFormation](https://github.com/awslabs/goformation) <-- not really what's needed, forces evaluation of intrinsics (completely mad and unuseable due to this) 

### AWS CLI Goodness

* [Controlling Output](https://docs.aws.amazon.com/cli/latest/userguide/controlling-output.html)