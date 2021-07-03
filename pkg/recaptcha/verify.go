package recaptcha

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gofor-little/xerror"

	"github.com/strongishllama/millhouse.dev-cdk/pkg/xhttp"
)

var (
	HTTPClient xhttp.Client
)

func Verify(ctx context.Context, secret string, challengeResponseToken string) (float32, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://www.google.com/recaptcha/api/siteverify", nil)
	if err != nil {
		return 0, xerror.Wrap("failed to build HTTP request", err)
	}

	query := request.URL.Query()
	query.Add("secret", secret)
	query.Add("response", challengeResponseToken)
	request.URL.RawQuery = query.Encode()

	if HTTPClient == nil {
		HTTPClient = &http.Client{
			Timeout: time.Duration(3 * time.Second),
		}
	}

	response, err := HTTPClient.Do(request)
	if err != nil {
		return 0, xerror.Wrap("failed to send HTTP request", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return 0, xerror.Newf("unexpected status code returned", response.StatusCode)
	}

	responseData := &ResponseData{}
	if err := json.NewDecoder(response.Body).Decode(&responseData); err != nil {
		return 0, xerror.Wrap("failed to decode response body", err)
	}

	if !responseData.Success {
		return 0, xerror.Newf("verify challenge failed", responseData.ErrorCodes)
	}

	return responseData.Score, nil
}
