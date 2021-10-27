package handler_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gofor-little/env"
	"github.com/gofor-little/log"
	"github.com/gofor-little/xlambda"
	"github.com/gofor-little/xrand"
	"github.com/stretchr/testify/require"

	"github.com/strongishllama/millhouse.dev-cdk/internal/db"
	"github.com/strongishllama/millhouse.dev-cdk/lambdas/api/unsubscribe/handler"
)

var (
	setupSleep    time.Duration = 1
	teardownSleep time.Duration = 1
)

func TestHandle(t *testing.T) {
	subscription := setup(t)
	defer teardown(t)

	request, err := xlambda.ProxyRequest(http.MethodGet, map[string]string{
		"id":           subscription.ID,
		"emailAddress": subscription.EmailAddress,
	}, nil)
	require.NoError(t, err)

	response, err := handler.Handler(context.Background(), request)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode)
}

func setup(t *testing.T) *db.Subscription {
	if err := env.Load("../../../../.env"); err != nil {
		t.Logf("failed to load .env file, ignore if running in CI/CD: %v", err)
	}

	log.Log = log.NewStandardLogger(os.Stdout, nil)
	require.NoError(t, db.Initialize(context.Background(), env.Get("TEST_AWS_PROFILE", ""), env.Get("TEST_AWS_REGION", ""), fmt.Sprintf("millhouse-dev-handle-test_%d", time.Now().Unix())))

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("pk"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("sk"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		BillingMode: types.BillingModePayPerRequest,
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("pk"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("sk"),
				KeyType:       types.KeyTypeRange,
			},
		},
		TableName: aws.String(db.TableName),
	}

	_, err := db.DynamoDBClient.CreateTable(context.Background(), input)
	require.NoError(t, err)

	return createSubscription(t)
}

func teardown(t *testing.T) {
	input := &dynamodb.DeleteTableInput{
		TableName: aws.String(db.TableName),
	}

	_, err := db.DynamoDBClient.DeleteTable(context.Background(), input)
	if err != nil && teardownSleep <= 30 {
		aerr := &types.ResourceInUseException{}
		if !errors.As(err, &aerr) {
			require.NoError(t, err)
		}

		// If the table is still in use, such as being created. Wait, then try again.
		t.Logf("table still in use, trying again in %d seconds...", teardownSleep)
		time.Sleep(teardownSleep * time.Second)

		// Increment the sleep duration. If it reaches 30 seconds fail the test.
		teardownSleep++
		teardown(t)
	} else {
		require.NoError(t, err)
	}
}

func createSubscription(t *testing.T) *db.Subscription {
	id, err := xrand.UUIDV4()
	require.NoError(t, err)
	subscription := &db.Subscription{
		ID:           id,
		EmailAddress: "test@example.com",
		IsConfirmed:  false,
	}

	if err := subscription.Create(context.Background()); err != nil && setupSleep <= 30 {
		aerr := &types.ResourceNotFoundException{}
		if !errors.As(err, &aerr) {
			require.NoError(t, err)
		}

		// If the table is still being created. Wait, then try again.
		t.Logf("table still being created, trying again in %d seconds...", setupSleep)
		time.Sleep(setupSleep * time.Second)

		// Increment the sleep duration. If it reaches 30 seconds fail the test.
		setupSleep++
		return createSubscription(t)
	} else {
		require.NoError(t, err)
	}

	return subscription
}
