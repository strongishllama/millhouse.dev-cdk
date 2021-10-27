package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/mail"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gofor-little/xlambda"
	"github.com/gofor-little/xrand"

	"github.com/strongishllama/millhouse.dev-cdk/internal/db"
	"github.com/strongishllama/millhouse.dev-cdk/internal/recaptcha"
)

var (
	RecaptchaSecret string
)

func Handler(ctx context.Context, request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	data := &RequestData{}
	if err := xlambda.UnmarshalAndValidate(request, data); err != nil {
		return xlambda.ProxyResponseJSON(http.StatusBadRequest, err, nil)
	}

	score, err := recaptcha.Verify(ctx, RecaptchaSecret, data.ReCaptchaChallengeToken)
	if err != nil {
		return xlambda.ProxyResponseJSON(http.StatusInternalServerError, fmt.Errorf("recaptcha verification failed: %w", err), nil)
	}
	if score <= 0.5 {
		return xlambda.ProxyResponseJSON(http.StatusOK, nil, nil)
	}

	subscription, err := db.GetSubscription(ctx, data.EmailAddress)
	if err != nil {
		return xlambda.ProxyResponseJSON(http.StatusInternalServerError, fmt.Errorf("failed to check if subscription already exists: %w", err), nil)
	}
	if subscription != nil {
		return xlambda.ProxyResponseJSON(http.StatusOK, nil, nil)
	}

	id, err := xrand.UUIDV4()
	if err != nil {
		return nil, fmt.Errorf("failed to generate UUID: %w", err)
	}
	subscription = &db.Subscription{
		ID:           id,
		EmailAddress: data.EmailAddress,
		IsConfirmed:  false,
	}
	if err := subscription.Create(ctx); err != nil {
		return xlambda.ProxyResponseJSON(http.StatusInternalServerError, fmt.Errorf("failed to create subscription: %w", err), nil)
	}

	return xlambda.ProxyResponseJSON(http.StatusOK, nil, nil)
}

type RequestData struct {
	EmailAddress            string `json:"emailAddress"`
	ReCaptchaChallengeToken string `json:"recaptchaChallengeToken"`
}

func (r *RequestData) Validate() error {
	if _, err := mail.ParseAddress(r.EmailAddress); err != nil {
		return fmt.Errorf("failed to validate EmailAddress: %w", err)
	}

	if len(r.ReCaptchaChallengeToken) == 0 {
		return errors.New("ReCaptchaChallengeToken cannot be empty")
	}

	return nil
}
