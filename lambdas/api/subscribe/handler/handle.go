package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/mail"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gofor-little/xerror"

	"github.com/strongishllama/millhouse.dev-cdk/pkg/api"
	"github.com/strongishllama/millhouse.dev-cdk/pkg/lambda"
)

func Handle(ctx context.Context, request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var requestData *SubscribeRequest
	if err := json.Unmarshal([]byte(request.Body), &requestData); err != nil {
		return lambda.NewProxyResponse(http.StatusBadRequest, xerror.New("failed to unmarshal request body into subscribeRequest", err), nil)
	}

	if err := requestData.Validate(); err != nil {
		return lambda.NewProxyResponse(http.StatusBadRequest, xerror.New("failed to validate subscribeRequest", err), nil)
	}

	switch request.HTTPMethod {
	case http.MethodPut:
		if err := api.Subscribe(ctx, config.RecaptchaSecret, requestData.ReCaptchaChallengeToken, requestData.EmailAddress); err != nil {
			return lambda.NewProxyResponse(http.StatusInternalServerError, xerror.New("failed to add new email subscription", err), nil)
		}
	default:
		return lambda.NewProxyResponse(http.StatusMethodNotAllowed, nil, nil)
	}

	return lambda.NewProxyResponse(http.StatusOK, nil, nil)
}

type SubscribeRequest struct {
	EmailAddress            string `json:"emailAddress"`
	ReCaptchaChallengeToken string `json:"recaptchaChallengeToken"`
}

func (s *SubscribeRequest) Validate() error {
	if _, err := mail.ParseAddress(s.EmailAddress); err != nil {
		return xerror.New("failed to validate EmailAddress", err)
	}

	if s.ReCaptchaChallengeToken == "" {
		return xerror.Newf("ReCaptchaChallengeToken cannot be empty")
	}

	return nil
}
