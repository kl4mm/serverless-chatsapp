package main

import (
	"context"
	"fmt"
	"log"

	"github.com/7junky/serverless_chatsapp/pkg/api"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, request api.Request) (api.Response, error) {
	connectionId := request.RequestContext.ConnectionID

	s := session.Must(session.NewSession())
	cfg := aws.NewConfig().WithRegion("eu-west-2").WithEndpoint(
		request.RequestContext.DomainName + "/" + request.RequestContext.Stage)

	apigw := apigatewaymanagementapi.New(s, cfg)

	getInput := apigatewaymanagementapi.GetConnectionInput{
		ConnectionId: &connectionId,
	}

	info, err := apigw.GetConnection(&getInput)
	if err != nil {
		log.Printf("Error: %v", err)

		return api.Response{
			StatusCode: 500,
		}, nil
	}

	msg := fmt.Sprintf("Use the sendmessage route to send a message. Your info: %v", info)
	postInput := &apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: &connectionId,
		Data:         []byte(msg),
	}
	_, err = apigw.PostToConnection(postInput)
	if err != nil {
		log.Printf("Error: %v", err)

		return api.Response{
			StatusCode: 500,
		}, nil
	}

	return api.Response{
		StatusCode: 200,
	}, nil
}
