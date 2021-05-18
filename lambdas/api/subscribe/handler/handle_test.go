package handler_test

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gofor-little/log"
	"github.com/stretchr/testify/require"

	"github.com/strongishllama/millhouse.dev-cdk/lambdas/api/subscribe/handler"
)

func TestHandle(t *testing.T) {
	log.Log = log.NewStandardLogger(os.Stdout, nil)

	response, err := handler.Handle(context.Background(), &events.APIGatewayProxyRequest{
		HTTPMethod: http.MethodPut,
		Body: `{"emailAddress":"test@example.com"}`,
	})

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode)
}
