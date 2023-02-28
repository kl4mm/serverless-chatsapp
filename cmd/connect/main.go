package main

import (
	"context"
	"log"
	"os"

	"github.com/7junky/serverless_chatsapp/pkg/api"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, request api.Request) (api.Response, error) {
	s := session.Must(session.NewSession())
	cfg := aws.NewConfig().WithRegion("eu-west-2")

	db := dynamodb.New(s, cfg)
	table := os.Getenv("table")

	input := &dynamodb.PutItemInput{
		TableName: &table,
		Item: map[string]*dynamodb.AttributeValue{
			"connectionId": {
				S: &request.RequestContext.ConnectionID,
			},
		},
	}

	_, err := db.PutItem(input)
	if err != nil {
		log.Printf("Error: %v", err)
		return api.Response{
			StatusCode: 500,
			Body:       "Internal Server Error",
		}, nil
	}

	return api.Response{
		StatusCode: 200,
		Body:       "",
	}, nil
}