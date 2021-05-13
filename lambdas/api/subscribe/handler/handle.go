package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/mail"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gofor-little/xerror"

	"github.com/strongishllama/millhouse.dev-cdk/pkg/api"
)

func Handle(ctx context.Context, request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var requestData *SubscribeRequest
	if err := json.Unmarshal([]byte(request.Body), &requestData); err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
		}, xerror.New("failed to unmarshal request body into subscribeRequest", err)
	}

	if err := requestData.Validate(); err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
		}, xerror.New("failed to validate subscribeRequest", err)
	}

	switch request.HTTPMethod {
	case http.MethodPut:
		if err := api.Subscribe(requestData.EmailAddress); err != nil {
			return &events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
			}, xerror.New("failed to add new email subscription", err)
		}
	default:
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusMethodNotAllowed,
		}, nil
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
	}, nil
}

type SubscribeRequest struct {
	EmailAddress string `json:"emailAddress"`
}

func (s *SubscribeRequest) Validate() error {
	if _, err := mail.ParseAddress(s.EmailAddress); err != nil {
		return xerror.New("failed to validate EmailAddress", err)
	}

	return nil
}
