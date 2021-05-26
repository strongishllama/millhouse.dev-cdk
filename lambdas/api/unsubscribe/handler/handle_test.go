package handler_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gofor-little/env"
	"github.com/gofor-little/log"
	"github.com/stretchr/testify/require"

	"github.com/strongishllama/millhouse.dev-cdk/lambdas/api/unsubscribe/handler"
	"github.com/strongishllama/millhouse.dev-cdk/pkg/db"
	"github.com/strongishllama/millhouse.dev-cdk/pkg/xlambda"
)

var (
	setupSleep    time.Duration = 1
	teardownSleep time.Duration = 1
)

func TestHandle(t *testing.T) {
	subscription := setup(t)
	defer teardown(t)

	request, err := xlambda.NewProxyRequest(http.MethodGet, map[string]string{
		"id":           subscription.ID,
		"emailAddress": subscription.EmailAddress,
	}, nil)
	require.NoError(t, err)

	response, err := handler.Handle(context.Background(), request)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode)
}

func setup(t *testing.T) *db.Subscription {
	if err := env.Load("../../../../.env"); err != nil {
		t.Logf("failed to load .env file, ignore if running in CI/CD: %v", err)
	}

	log.Log = log.NewStandardLogger(os.Stdout, nil)
	require.NoError(t, db.Initialize(env.Get("TEST_AWS_PROFILE", ""), env.Get("TEST_AWS_REGION", ""), fmt.Sprintf("millhouse-dev-handle-test_%d", time.Now().Unix())))

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("pk"),
				AttributeType: aws.String(dynamodb.ScalarAttributeTypeS),
			},
			{
				AttributeName: aws.String("sk"),
				AttributeType: aws.String(dynamodb.ScalarAttributeTypeS),
			},
		},
		BillingMode: aws.String(dynamodb.BillingModePayPerRequest),
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("pk"),
				KeyType:       aws.String(dynamodb.KeyTypeHash),
			},
			{
				AttributeName: aws.String("sk"),
				KeyType:       aws.String(dynamodb.KeyTypeRange),
			},
		},
		TableName: aws.String(db.TableName),
	}
	require.NoError(t, input.Validate())

	_, err := db.DynamoDBClient.CreateTable(input)
	require.NoError(t, err)

	return createSubscription(t)
}

func teardown(t *testing.T) {
	input := &dynamodb.DeleteTableInput{
		TableName: aws.String(db.TableName),
	}
	require.NoError(t, input.Validate())

	_, err := db.DynamoDBClient.DeleteTable(input)
	if err != nil && teardownSleep <= 30 {
		aerr, ok := err.(awserr.Error)
		if !ok || aerr.Code() != dynamodb.ErrCodeResourceInUseException {
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
	subscription, err := db.CreateSubscription(context.Background(), "test@example.com")
	if err != nil && setupSleep <= 30 {
		aerr, ok := errors.Unwrap(err).(awserr.Error)
		if !ok || aerr.Code() != dynamodb.ErrCodeResourceNotFoundException {
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
