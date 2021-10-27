package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gofor-little/xrand"
)

func CreateSubscription(ctx context.Context, emailAddress string) (*Subscription, error) {
	if err := checkPackage(); err != nil {
		return nil, err
	}

	id, err := xrand.UUIDV4()
	if err != nil {
		return nil, fmt.Errorf("failed to generate UUID: %w", err)
	}
	subscription := &Subscription{
		ID:           id,
		EmailAddress: emailAddress,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := subscription.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate subscription: %w", err)
	}

	attributeValues, err := attributevalue.MarshalMap(subscription)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal subscription into attribute values: %w", err)
	}
	attributeValues["pk"] = &types.AttributeValueMemberS{
		Value: "SUBSCRIPTION#" + subscription.EmailAddress,
	}
	attributeValues["sk"] = &types.AttributeValueMemberS{
		Value: "SUBSCRIPTION#" + subscription.ID,
	}

	input := &dynamodb.PutItemInput{
		Item:      attributeValues,
		TableName: aws.String(TableName),
	}

	if _, err := DynamoDBClient.PutItem(ctx, input); err != nil {
		return nil, err
	}

	return subscription, nil
}

func DeleteSubscription(ctx context.Context, id string, emailAddress string) error {
	if err := checkPackage(); err != nil {
		return err
	}

	input := &dynamodb.DeleteItemInput{
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{
				Value: "SUBSCRIPTION#" + emailAddress,
			},
			"sk": &types.AttributeValueMemberS{
				Value: "SUBSCRIPTION#" + id,
			},
		},
		TableName: aws.String(TableName),
	}

	if _, err := DynamoDBClient.DeleteItem(ctx, input); err != nil {
		return err
	}

	return nil
}

func GetSubscription(ctx context.Context, emailAddress string) (*Subscription, error) {
	if err := checkPackage(); err != nil {
		return nil, err
	}

	if emailAddress == "" {
		return nil, errors.New("emailAddress cannot be empty")
	}

	input := &dynamodb.QueryInput{
		KeyConditions: map[string]types.Condition{
			"pk": {
				ComparisonOperator: types.ComparisonOperatorEq,
				AttributeValueList: []types.AttributeValue{
					&types.AttributeValueMemberS{
						Value: "SUBSCRIPTION#" + emailAddress,
					},
				},
			},
		},
		Limit:     aws.Int32(1),
		TableName: aws.String(TableName),
	}

	output, err := DynamoDBClient.Query(ctx, input)
	if err != nil {
		return nil, err
	}

	if len(output.Items) == 0 {
		return nil, nil
	}

	var subscription *Subscription
	if err := attributevalue.UnmarshalMap(output.Items[0], &subscription); err != nil {
		return nil, fmt.Errorf("failed to unmarshal attribute values into Subscription: %w", err)
	}

	return subscription, nil
}

type Subscription struct {
	ID           string    `json:"id"`
	EmailAddress string    `json:"emailAddress"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func (s *Subscription) Validate() error {
	if len(s.EmailAddress) == 0 {
		return errors.New("email address cannot be empty")
	}

	if s.CreatedAt == (time.Time{}) {
		return errors.New("created at cannot be empty")
	}

	if s.UpdatedAt == (time.Time{}) {
		return errors.New("updated at cannot be empty")
	}

	return nil
}
