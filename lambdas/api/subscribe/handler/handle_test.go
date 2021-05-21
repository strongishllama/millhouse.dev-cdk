package handler_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/gofor-little/env"
	"github.com/gofor-little/log"
	"github.com/stretchr/testify/require"

	"github.com/strongishllama/millhouse.dev-cdk/lambdas/api/subscribe/handler"
	"github.com/strongishllama/millhouse.dev-cdk/pkg/lambda"
	"github.com/strongishllama/millhouse.dev-cdk/pkg/notifications"
	"github.com/strongishllama/millhouse.dev-cdk/pkg/recaptcha"
	"github.com/strongishllama/millhouse.dev-cdk/pkg/xhttp"
)

func setup(t *testing.T) string {
	if err := env.Load("../../../../.env"); err != nil {
		t.Logf("failed to load .env file, ignore if running in CI/CD: %v", err)
	}
	require.NoError(t, notifications.Initialize(env.Get("AWS_PROFILE", ""), env.Get("AWS_REGION", "")))

	input := &sqs.CreateQueueInput{
		QueueName: aws.String(fmt.Sprintf("HandleTest_%d", time.Now().Unix())),
	}
	require.NoError(t, input.Validate())

	output, err := notifications.SQSClient.CreateQueue(input)
	require.NoError(t, err)

	log.Log = log.NewStandardLogger(os.Stdout, nil)
	recaptcha.HTTPClient = &xhttp.MockClient{
		ResponseData: &recaptcha.ResponseData{
			Score:   0.9,
			Success: true,
		},
	}
	handler.Cfg = &handler.Config{
		RecaptchaSecret: "test-recaptcha-secret",
		QueueURL:        *output.QueueUrl,
		From:            "no-reply@millhouse.dev",
	}

	return *output.QueueUrl
}

func teardown(t *testing.T, queueURL string) {
	input := &sqs.DeleteQueueInput{
		QueueUrl: aws.String(queueURL),
	}
	require.NoError(t, input.Validate())

	_, err := notifications.SQSClient.DeleteQueue(input)
	require.NoError(t, err)
}

func TestHandle(t *testing.T) {
	queueURL := setup(t)
	defer teardown(t, queueURL)

	request, err := lambda.NewProxyRequest(http.MethodPut, &handler.RequestData{
		EmailAddress:            "test@example.com",
		ReCaptchaChallengeToken: "test-recaptcha-challenge-token",
	})
	require.NoError(t, err)

	response, err := handler.Handle(context.Background(), request)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode)
}
