package api

import "github.com/aws/aws-lambda-go/events"

type Request = events.APIGatewayWebsocketProxyRequest
type Response = events.APIGatewayProxyResponse
