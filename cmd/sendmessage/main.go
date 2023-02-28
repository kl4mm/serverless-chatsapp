package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/7junky/serverless_chatsapp/pkg/api"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func main() {
	lambda.Start(handleRequest)
	fmt.Println("Hello world")
}

type message struct {
	Message string `json:"message"`
}

func handleRequest(ctx context.Context, request api.Request) (api.Response, error) {
	s := session.Must(session.NewSession())
	cfg := aws.NewConfig().WithRegion("eu-west-2")

	db := dynamodb.New(s, cfg)
	apigw := apigatewaymanagementapi.New(s, cfg)
	table := os.Getenv("table")

	scanInput := &dynamodb.ScanInput{
		TableName: &table,
	}

	connections, err := db.Scan(scanInput)
	if err != nil {
		log.Printf("Error: %v", err)

		return api.Response{
			StatusCode: 500,
		}, nil
	}

	msg := &message{}
	err = json.Unmarshal([]byte(request.Body), msg)
	if err != nil {
		log.Printf("Error: %v", err)

		return api.Response{
			StatusCode: 500,
		}, nil
	}

	log.Printf("Connections: %+v", connections.Items)
	for _, connection := range connections.Items {
		connectionId := connection["connectionId"].S
		if *connectionId != request.RequestContext.ConnectionID {
			_, err := apigw.PostToConnection(&apigatewaymanagementapi.PostToConnectionInput{
				ConnectionId: connectionId,
				Data:         []byte(msg.Message),
			})

			if err != nil {
				log.Printf("Error: %v", err)

				return api.Response{
					StatusCode: 500,
				}, nil
			}
		}
	}

	return api.Response{
		StatusCode: 200,
	}, nil
}