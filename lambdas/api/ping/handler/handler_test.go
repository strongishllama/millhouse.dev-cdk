package handler_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/strongishllama/millhouse.dev-cdk/lambdas/api/ping/handler"
)

func TestHandler(t *testing.T) {
	response, err := handler.Handler(context.Background(), nil)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode)
}
