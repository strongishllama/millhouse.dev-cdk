package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	email "github.com/gofor-little/aws-email"
	"github.com/gofor-little/log"
	"github.com/gofor-little/xerror"

	"github.com/strongishllama/millhouse.dev-cdk/pkg/lambda"
	"github.com/strongishllama/millhouse.dev-cdk/pkg/notifications"
	"github.com/strongishllama/millhouse.dev-cdk/pkg/recaptcha"
)

func Handle(ctx context.Context, request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var requestData *RequestData
	if err := json.Unmarshal([]byte(request.Body), &requestData); err != nil {
		return lambda.NewProxyResponse(http.StatusBadRequest, xerror.New("failed to unmarshal request body into subscribeRequest", err), nil)
	}

	if err := requestData.Validate(); err != nil {
		return lambda.NewProxyResponse(http.StatusBadRequest, xerror.New("failed to validate request data", err), nil)
	}

	log.Info(log.Fields{
		"message":      "subscribe request initiated",
		"emailAddress": requestData.EmailAddress,
	})

	switch request.HTTPMethod {
	case http.MethodPut:
		score, err := recaptcha.Verify(ctx, Cfg.RecaptchaSecret, requestData.ReCaptchaChallengeToken)
		if err != nil {
			return lambda.NewProxyResponse(http.StatusInternalServerError, xerror.New("recaptcha verification failed", err), nil)
		}

		log.Info(log.Fields{
			"message": "recaptcha verification successful",
			"score":   score,
		})

		if score <= 0.5 {
			// Send email notifying of low sore, then exit.
		}

		messageID, err := notifications.SendEmail(ctx, Cfg.QueueURL, []string{requestData.EmailAddress}, Cfg.From, "Subscription Confirmation", "subscription-confirmation.tmpl", email.ContentTypeTextHTML)
		if err != nil {
			return lambda.NewProxyResponse(http.StatusInternalServerError, xerror.New("failed to send subscription confirmation email", err), nil)
		}

		log.Info(log.Fields{
			"message":   "subscription confirmation email sent",
			"messageID": messageID,
		})
	default:
		return lambda.NewProxyResponse(http.StatusMethodNotAllowed, nil, nil)
	}

	return lambda.NewProxyResponse(http.StatusOK, nil, nil)
}
