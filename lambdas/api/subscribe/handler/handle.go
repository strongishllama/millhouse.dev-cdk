package handler

import (
	"context"
	"net/http"
	"net/mail"

	"github.com/aws/aws-lambda-go/events"
	email "github.com/gofor-little/aws-email"
	"github.com/gofor-little/log"
	"github.com/gofor-little/xerror"
	"github.com/strongishllama/xlambda"

	"github.com/strongishllama/millhouse.dev-cdk/pkg/db"
	"github.com/strongishllama/millhouse.dev-cdk/pkg/notification"
	"github.com/strongishllama/millhouse.dev-cdk/pkg/recaptcha"
)

var (
	AdminTo         string
	AdminFrom       string
	RecaptchaSecret string
)

func Handle(ctx context.Context, request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	if request.HTTPMethod != http.MethodPut {
		return xlambda.NewProxyResponse(http.StatusMethodNotAllowed, xlambda.ContentTypeApplicationJSON, nil, nil)
	}

	// Unmarshal and validate the request data.
	data := &RequestData{}
	if err := xlambda.UnmarshalAndValidate(request, data); err != nil {
		return xlambda.NewProxyResponse(http.StatusBadRequest, xlambda.ContentTypeApplicationJSON, xerror.Wrap("failed to unmarshal and validate request data", err), nil)
	}

	// Verify the recaptcha score.
	score, err := recaptcha.Verify(ctx, RecaptchaSecret, data.ReCaptchaChallengeToken)
	if err != nil {
		return xlambda.NewProxyResponse(http.StatusInternalServerError, xlambda.ContentTypeApplicationJSON, xerror.Wrap("recaptcha verification failed", err), nil)
	}

	// If the score is less than 0.5 enqueue an email notifying the admin, then exit.
	if score <= 0.5 {
		messageID, err := notification.EnqueueEmail(ctx, []string{AdminTo}, AdminFrom, notification.EmailTemplate{
			FileName:    "email/recaptcha-challenge-failed.tmpl.html",
			Subject:     "Recaptcha Challenge Failed",
			ContentType: email.ContentTypeTextHTML,
			Data: notification.RecaptchaChallengeFailedTemplateData{
				EmailAddress: data.EmailAddress,
				Score:        score,
			},
		})
		if err != nil {
			log.Error(log.Fields{
				"error":        xerror.Wrap("failed to enqueue recaptcha challenge failed email", err),
				"messageId":    messageID,
				"emailAddress": data.EmailAddress,
				"score":        score,
			})
		}

		return xlambda.NewProxyResponse(http.StatusOK, xlambda.ContentTypeApplicationJSON, nil, nil)
	}

	// Attempt to check if a subscription with that email already exists. If so, exit now.
	subscription, err := db.GetSubscription(ctx, data.EmailAddress)
	if err != nil {
		return xlambda.NewProxyResponse(http.StatusInternalServerError, xlambda.ContentTypeApplicationJSON, xerror.Wrap("failed to check if subscription already exists", err), nil)
	}
	if subscription != nil {
		return xlambda.NewProxyResponse(http.StatusOK, xlambda.ContentTypeApplicationJSON, nil, nil)
	}

	// Create the subscription.
	if _, err = db.CreateSubscription(ctx, data.EmailAddress); err != nil {
		return xlambda.NewProxyResponse(http.StatusInternalServerError, xlambda.ContentTypeApplicationJSON, xerror.Wrap("failed to create subscription", err), nil)
	}

	return xlambda.NewProxyResponse(http.StatusOK, xlambda.ContentTypeApplicationJSON, nil, nil)
}

type RequestData struct {
	EmailAddress            string `json:"emailAddress"`
	ReCaptchaChallengeToken string `json:"recaptchaChallengeToken"`
}

func (r *RequestData) Validate() error {
	if _, err := mail.ParseAddress(r.EmailAddress); err != nil {
		return xerror.Wrap("failed to validate EmailAddress", err)
	}

	if r.ReCaptchaChallengeToken == "" {
		return xerror.New("ReCaptchaChallengeToken cannot be empty")
	}

	return nil
}
