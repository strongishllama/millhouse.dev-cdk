package handler

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/strongishllama/xlambda"
)

func Handle(ctx context.Context, request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	return xlambda.NewProxyResponse(http.StatusOK, xlambda.ContentTypeApplicationJSON, nil, `{"message": "pong"}`)
}
