package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	email "github.com/gofor-little/aws-email"
	"github.com/gofor-little/log"
	"github.com/gofor-little/xerror"

	"github.com/strongishllama/millhouse.dev-cdk/pkg/db"
	"github.com/strongishllama/millhouse.dev-cdk/pkg/notification"
	"github.com/strongishllama/millhouse.dev-cdk/pkg/recaptcha"
	"github.com/strongishllama/millhouse.dev-cdk/pkg/xlambda"
)

var (
	AdminTo         string
	AdminFrom       string
	RecaptchaSecret string
)

func Handle(ctx context.Context, request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	if request.HTTPMethod != http.MethodPut {
		return xlambda.NewProxyResponse(http.StatusMethodNotAllowed, "", nil, nil)
	}

	// Unmarshal and validate the request data.
	var requestData *RequestData
	if err := json.Unmarshal([]byte(request.Body), &requestData); err != nil {
		return xlambda.NewProxyResponse(http.StatusBadRequest, "", xerror.Wrap("failed to unmarshal request body into subscribeRequest", err), nil)
	}
	if err := requestData.Validate(); err != nil {
		return xlambda.NewProxyResponse(http.StatusBadRequest, "", xerror.Wrap("failed to validate request data", err), nil)
	}

	// Verify the recaptcha score.
	score, err := recaptcha.Verify(ctx, RecaptchaSecret, requestData.ReCaptchaChallengeToken)
	if err != nil {
		return xlambda.NewProxyResponse(http.StatusInternalServerError, "", xerror.Wrap("recaptcha verification failed", err), nil)
	}

	// If the score is less than 0.5 enqueue an email notifying the admin, then exit.
	if score <= 0.5 {
		messageID, err := notification.EnqueueEmail(ctx, []string{AdminTo}, AdminFrom, notification.EmailTemplate{
			FileName:    "email/recaptcha-challenge-failed.tmpl.html",
			Subject:     "Recaptcha Challenge Failed",
			ContentType: email.ContentTypeTextHTML,
			Data: notification.RecaptchaChallengeFailedTemplateData{
				EmailAddress: requestData.EmailAddress,
				Score:        score,
			},
		})
		if err != nil {
			log.Error(log.Fields{
				"error":        xerror.Wrap("failed to enqueue recaptcha challenge failed email", err),
				"messageId":    messageID,
				"emailAddress": requestData.EmailAddress,
				"score":        score,
			})
		}

		return xlambda.NewProxyResponse(http.StatusOK, xlambda.ContentTypeApplicationJSON, nil, nil)
	}

	// Attempt to check if a subscription with that email already exists. If so, exit now.
	subscription, err := db.GetSubscription(ctx, requestData.EmailAddress)
	if err != nil {
		return xlambda.NewProxyResponse(http.StatusInternalServerError, "", xerror.Wrap("failed to check if subscription already exists", err), nil)
	}
	if subscription != nil {
		return xlambda.NewProxyResponse(http.StatusOK, "", nil, nil)
	}

	// Create the subscription.
	if _, err = db.CreateSubscription(ctx, requestData.EmailAddress); err != nil {
		return xlambda.NewProxyResponse(http.StatusInternalServerError, "", xerror.Wrap("failed to create subscription", err), nil)
	}

	return xlambda.NewProxyResponse(http.StatusOK, "", nil, nil)
}
