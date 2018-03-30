package main

import (
	"fmt"
	"context"
	"encoding/json"
	"time"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"os"
)

// Request allows to define JSON format and marshaller for Lambda request
type Request struct {
	DateOfBirth string `json:"dateOfBirth"`
}

func main() {
	awsSession, awsSessionError := session.NewSession()
	if awsSessionError != nil {
		panic(awsSessionError)
	}

	dynamoTable := os.Getenv("TABLE_NAME")
	dynamoClient := dynamodb.New(awsSession)
	updateExpr := "SET dob = :dob"
	returnCollectionMetrics := "SIZE"
	returnConsumedCapacity := "TOTAL"

	lambda.Start(
		func(
			ctx context.Context,
			request events.APIGatewayProxyRequest,
		) (events.APIGatewayProxyResponse, error) {
			usernameParam := request.PathParameters["username"]

			if request.Headers["Content-Type"] != "application/json" {
				return events.APIGatewayProxyResponse{
					StatusCode: 400,
					Body:       "application/json request content type expected",
				}, nil
			}

			// TODO: defend from lengthy requests
			// TODO: test request gzip compression
			req := Request{}
			jsonUnmarshalErr := json.Unmarshal([]byte(request.Body), &req)
			if jsonUnmarshalErr != nil {
				return events.APIGatewayProxyResponse{
					StatusCode: 400,
					Body:       jsonUnmarshalErr.Error(),
				}, nil
			}

			dob, dobParseError := time.Parse("2006-01-02", req.DateOfBirth)
			if dobParseError != nil {
				return events.APIGatewayProxyResponse{
					StatusCode: 400,
					Body:       dobParseError.Error(),
				}, nil
			}

			updateItemInput :=
				&dynamodb.UpdateItemInput{
					TableName: aws.String(dynamoTable),
					Key: map[string]*dynamodb.AttributeValue{
						"username": {
							S: aws.String(usernameParam),
						},
					},
					UpdateExpression: &updateExpr,
					ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
						":dob": {
							N: aws.String(fmt.Sprintf("%d", int64(dob.Unix()))),
						},
					},
					ReturnConsumedCapacity:      &returnConsumedCapacity,
					ReturnItemCollectionMetrics: &returnCollectionMetrics,
				}

			updateItemOutput, updateErr := dynamoClient.UpdateItem(updateItemInput)

			if updateErr != nil {
				return events.APIGatewayProxyResponse{
					StatusCode: 503,
					Body:       updateErr.Error(),
				}, nil
			}

			sizeEstimateMin := float64(-1)
			if updateItemOutput.ItemCollectionMetrics != nil &&
				updateItemOutput.ItemCollectionMetrics.SizeEstimateRangeGB[0] != nil {
				sizeEstimateMin = *(updateItemOutput.ItemCollectionMetrics.SizeEstimateRangeGB[0])
			}
			sizeEstimateMax := float64(-1)
			if updateItemOutput.ItemCollectionMetrics != nil &&
				updateItemOutput.ItemCollectionMetrics.SizeEstimateRangeGB[1] != nil {
				sizeEstimateMax = *(updateItemOutput.ItemCollectionMetrics.SizeEstimateRangeGB[1])
			}

			//	report to CloudWatch or Prometheus+Grafana
			respHeaders := map[string]string{
				"X-Consumed-Capacity-Total":
				fmt.Sprintf("%v", *(updateItemOutput.ConsumedCapacity.CapacityUnits)),
				"X-Size-Estimate-Range-GByte-Min":
				fmt.Sprintf("%v", sizeEstimateMin),
				"X-Size-Estimate-Range-GByte-Max":
				fmt.Sprintf("%v", sizeEstimateMax),
			}

			return events.APIGatewayProxyResponse{
				StatusCode: 201,
				Body:       "",
				Headers:    respHeaders,
			}, nil
		},
	)
}
