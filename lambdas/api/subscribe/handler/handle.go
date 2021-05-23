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
	"github.com/strongishllama/millhouse.dev-cdk/pkg/lambda"
	"github.com/strongishllama/millhouse.dev-cdk/pkg/notification"
	"github.com/strongishllama/millhouse.dev-cdk/pkg/recaptcha"
)

func Handle(ctx context.Context, request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var requestData *RequestData
	if err := json.Unmarshal([]byte(request.Body), &requestData); err != nil {
		return lambda.NewProxyResponse(http.StatusBadRequest, "", xerror.New("failed to unmarshal request body into subscribeRequest", err), nil)
	}

	if err := requestData.Validate(); err != nil {
		return lambda.NewProxyResponse(http.StatusBadRequest, "", xerror.New("failed to validate request data", err), nil)
	}

	switch request.HTTPMethod {
	case http.MethodPut:
		score, err := recaptcha.Verify(ctx, Cfg.RecaptchaSecret, requestData.ReCaptchaChallengeToken)
		if err != nil {
			return lambda.NewProxyResponse(http.StatusInternalServerError, "", xerror.New("recaptcha verification failed", err), nil)
		}

		if score <= 0.5 {
			messageID, err := notification.EnqueueEmail(ctx, []string{Cfg.To}, Cfg.From, notification.EmailTemplate{
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
					"error":        xerror.New("failed to enqueue recaptcha challenge failed email", err),
					"messageId":    messageID,
					"emailAddress": requestData.EmailAddress,
					"score":        score,
				})
			}

			return lambda.NewProxyResponse(http.StatusOK, lambda.ContentTypeApplicationJSON, nil, nil)
		}

		subscription, err := db.GetSubscription(ctx, requestData.EmailAddress)
		if err != nil {
			return lambda.NewProxyResponse(http.StatusInternalServerError, "", xerror.New("failed to check if subscription already exists", err), nil)
		}

		// A subscription for that email already exists, exit now.
		if subscription != nil {
			return lambda.NewProxyResponse(http.StatusOK, "", nil, nil)
		}

		if subscription, err = db.CreateSubscription(ctx, requestData.EmailAddress); err != nil {
			return lambda.NewProxyResponse(http.StatusInternalServerError, "", xerror.New("failed to create subscription", err), nil)
		}

		messageID, err := notification.EnqueueEmail(ctx, []string{requestData.EmailAddress}, Cfg.From, notification.EmailTemplate{
			FileName:    "email/subscription-confirmation.tmpl.html",
			Subject:     "Subscription Confirmation",
			ContentType: email.ContentTypeTextHTML,
			Data: notification.SubscriptionConfirmationTemplateData{
				WebsiteDomain:  Cfg.WebsiteDomain,
				APIDomain:      Cfg.APIDomain,
				SubscriptionID: subscription.ID,
				EmailAddress:   subscription.EmailAddress,
			},
		})
		if err != nil {
			return lambda.NewProxyResponse(http.StatusInternalServerError, "", xerror.New("failed to enqueue subscription confirmation email", err), nil)
		}
		log.Info(log.Fields{
			"message":      "successfully enqueued subscription confirmation email",
			"messageID":    messageID,
			"emailAddress": requestData.EmailAddress,
			"score":        score,
		})
	default:
		return lambda.NewProxyResponse(http.StatusMethodNotAllowed, "", nil, nil)
	}

	return lambda.NewProxyResponse(http.StatusOK, "", nil, nil)
}
