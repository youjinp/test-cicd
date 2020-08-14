package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sts"
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

	fmt.Println("Loading default config..")
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic(err)
	}

	// https://github.com/aws/aws-sdk-go-v2/blob/master/example/service/sts/assumeRole/assumeRole.go
	fmt.Println("Assuming role..")
	roleARN := "arn:aws:iam::556352520741:role/access-database-role"
	svc := sts.New(cfg)
	roleReq := svc.AssumeRoleRequest(&sts.AssumeRoleInput{RoleArn: aws.String(roleARN), RoleSessionName: aws.String("assumeTestRole")})
	roleRes, err := roleReq.Send(ctx)
	if err != nil {
		panic(err)
	}

	awsConfig := svc.Config.Copy()
	awsConfig.Credentials = CredentialsProvider{Credentials: roleRes.Credentials}

	//
	fmt.Println("Accessing database..")
	dynamodbClient := dynamodb.New(awsConfig)

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

type CredentialsProvider struct {
	*sts.Credentials
}

func (s CredentialsProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	if s.Credentials == nil {
		return aws.Credentials{}, errors.New("sts credentials are nil")
	}

	return aws.Credentials{
		AccessKeyID:     aws.StringValue(s.AccessKeyId),
		SecretAccessKey: aws.StringValue(s.SecretAccessKey),
		SessionToken:    aws.StringValue(s.SessionToken),
		Expires:         aws.TimeValue(s.Expiration),
	}, nil
}
