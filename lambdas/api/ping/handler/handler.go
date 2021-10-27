package handler

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gofor-little/xlambda"
)

func Handler(ctx context.Context, request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	return xlambda.ProxyResponseJSON(http.StatusOK, nil, `{"message": "pong"}`)
}
