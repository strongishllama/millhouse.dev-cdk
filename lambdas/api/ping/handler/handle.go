package handler

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func Handle(ctx context.Context, request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "{\"message\":\"pong\"}",
	}, nil
}
