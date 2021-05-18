package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gofor-little/log"
	"github.com/gofor-little/xerror"
)

func Subscribe(ctx context.Context, recaptchaSecret string, challengeResponseToken string, emailAddress string) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://www.google.com/recaptcha/api/siteverify", nil)
	if err != nil {
		return xerror.New("failed to build HTTP request", err)
	}

	query := request.URL.Query()
	query.Add("secret", recaptchaSecret)
	query.Add("response", challengeResponseToken)
	request.URL.RawQuery = query.Encode()

	response, err := HTTPClient.Do(request)
	if err != nil {
		return xerror.New("failed to send HTTP request", err)
	}
	defer response.Body.Close()

	responseData := &verifyResponse{}
	if err := json.NewDecoder(response.Body).Decode(&responseData); err != nil {
		xerror.New("failed to decode response body", err)
	}

	if responseData.Score <= 0.5 {
		// Send email notifying of low sore, then exit.
	}

	// Send confirmation email of subscription.
	// sendEmail(ctx, "TODO", []string{emailAddress}, "TODO", "Subscription Confirmation", "subscription-confirmation.tmpl", email.ContentTypeTextHTML)

	log.Info(log.Fields{
		"emailAddress":   emailAddress,
		"verifyResponse": responseData,
	})
	return nil
}

type verifyResponse struct {
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	Score       float32   `json:"score"`
	Success     bool      `json:"success"`
	ErrorCodes  []string  `json:"error-codes"`
}
