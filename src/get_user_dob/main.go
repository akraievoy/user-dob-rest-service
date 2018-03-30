package main

import (
	"fmt"
	"context"
	"time"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"strconv"
	"math"
	"encoding/json"
	"os"
)

// Response allows to define JSON format and marshaller for Lambda response
type Response struct {
	Message string `json:"message"`
}

func main() {
	awsSession, awsSessionError := session.NewSession()
	if awsSessionError != nil {
		panic(awsSessionError)
	}

	dynamoTable := os.Getenv("TABLE_NAME")
	consistentRead := os.Getenv("CONSISTENT_READ") == "true"
	dynamoClient := dynamodb.New(awsSession)
	projectionExpression := "dob"
	returnConsumedCapacity := "TOTAL"

	lambda.Start(
		func(
			ctx context.Context,
			request events.APIGatewayProxyRequest,
		) (events.APIGatewayProxyResponse, error) {
			usernameParam := request.PathParameters["username"]

			getItemInput :=
				&dynamodb.GetItemInput{
					TableName: aws.String(dynamoTable),
					Key: map[string]*dynamodb.AttributeValue{
						"username": {
							S: aws.String(usernameParam),
						},
					},
					ProjectionExpression:   &projectionExpression,
					ReturnConsumedCapacity: &returnConsumedCapacity,
					ConsistentRead:         &consistentRead,
				}

			getItemOutput, getItemErr := dynamoClient.GetItem(getItemInput)

			if getItemErr != nil {
				return events.APIGatewayProxyResponse{
					StatusCode: 503,
					Body:       getItemErr.Error(),
				}, nil
			}

			//	report to CloudWatch or Prometheus+Grafana
			respHeaders := map[string]string{
				"X-Consumed-Capacity-Total":
				fmt.Sprintf("%v", *(getItemOutput.ConsumedCapacity.CapacityUnits)),
			}

			if getItemOutput.Item == nil {
				return events.APIGatewayProxyResponse{
					StatusCode: 404,
					Body:       fmt.Sprintf("userDOBs record not found for user '%s'", usernameParam),
					Headers:    respHeaders,
				}, nil
			}

			for attrName, attrValue := range getItemOutput.Item {
				if attrName == projectionExpression && attrValue.N != nil {
					dobSeconds, parseError := strconv.ParseInt(*(attrValue.N), 10, 64)
					if parseError != nil {
						return events.APIGatewayProxyResponse{
							StatusCode: 503,
							Body:       parseError.Error(),
						}, nil
					}
					dobTime := time.Unix(dobSeconds, 0)
					nowTime := time.Now()
					nowSeconds := nowTime.Unix()
					message := "Hello " + usernameParam + "! "
					if nowSeconds < dobSeconds {
						message += fmt.Sprintf(
							"Your birthday is in %d days",
							int64(math.Ceil(float64(dobSeconds-nowSeconds)/86400.0)),
						)
					} else {
						if dobTime.Month() == nowTime.Month() && dobTime.Day() == nowTime.Day() {
							message += " Happy birthday!"
						} else {
							currentYearDOB :=
								time.Date(
									nowTime.Year(), dobTime.Month(), dobTime.Day(),
									0, 0, 0, 0,
									time.Local)

							var durationBeforeNextDOB time.Duration
							if currentYearDOB.Before(nowTime) {
								nextYearDOB := currentYearDOB.AddDate(1, 0 ,0)
								durationBeforeNextDOB = nextYearDOB.Sub(nowTime)
							} else {
								durationBeforeNextDOB = currentYearDOB.Sub(nowTime)
							}

							message +=
								fmt.Sprintf(
									"Your birthday is in %d days",
									int64(math.Ceil(durationBeforeNextDOB.Hours() / 24)),
								)
						}
					}

					resp := Response{
						Message: message,
					}
					jsonBytes, jsonErr := json.Marshal(resp)
					if jsonErr != nil {
						return events.APIGatewayProxyResponse{
							StatusCode: 500,
							Body:       "failed to generate JSON body",
							Headers:    respHeaders,
						}, nil
					}

					respHeaders["Content-Type"] = "application/json"
					return events.APIGatewayProxyResponse{
						StatusCode: 200,
						Body:       string(jsonBytes),
						Headers:    respHeaders,
					}, nil
				}
			}

			return events.APIGatewayProxyResponse{
				StatusCode: 503,
				Body:       "userDOBs record found but has required fields missing",
				Headers:    respHeaders,
			}, nil
		},
	)
}
