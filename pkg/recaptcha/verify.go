package recaptcha

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/strongishllama/millhouse.dev-cdk/pkg/xhttp"
)

var (
	HTTPClient xhttp.Client
)

func Verify(ctx context.Context, secret string, challengeResponseToken string) (float32, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://www.google.com/recaptcha/api/siteverify", nil)
	if err != nil {
		return 0, fmt.Errorf("failed to build HTTP request: %w", err)
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
		return 0, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code returned: %d", response.StatusCode)
	}

	responseData := &ResponseData{}
	if err := json.NewDecoder(response.Body).Decode(&responseData); err != nil {
		return 0, fmt.Errorf("failed to decode response body: %w", err)
	}

	if !responseData.Success {
		return 0, fmt.Errorf("verify challenge failed: %v", responseData.ErrorCodes)
	}

	return responseData.Score, nil
}
