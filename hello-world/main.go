package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
)

var (
	// DefaultHTTPGetAddress Default Address
	DefaultHTTPGetAddress = "https://checkip.amazonaws.com"

	// ErrNoIP No IP found in response
	ErrNoIP = errors.New("No IP in HTTP response")

	// ErrNon200Response non 200 status code in response
	ErrNon200Response = errors.New("Non 200 Response found")
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic(err)
	}
	dynamodbClient := dynamodb.New(cfg)

	req := dynamodbClient.GetItemRequest(&dynamodb.GetItemInput{
		TableName: aws.String("GeoIP-2020-05-29"),
		Key: map[string]dynamodb.AttributeValue{
			"networkHead": {
				N: aws.String("187"),
			},
			"firstIP": {
				N: aws.String("3137339392"),
			},
		},
	})

	res, err := req.Send(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println("res: ", res.Item)

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Hello"),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
