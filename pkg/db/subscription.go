package db

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gofor-little/xerror"
	"github.com/gofor-little/xrand"
)

func CreateSubscription(ctx context.Context, emailAddress string) (*Subscription, error) {
	if err := checkPackage(); err != nil {
		return nil, xerror.New("checkPackage failed", err)
	}

	id, err := xrand.UUIDV4()
	if err != nil {
		return nil, xerror.New("failed to generate UUID", err)
	}
	subscription := &Subscription{
		ID:           id,
		EmailAddress: emailAddress,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := subscription.Validate(); err != nil {
		return nil, xerror.New("failed to validate subscription", err)
	}

	attributeValues, err := dynamodbattribute.MarshalMap(subscription)
	if err != nil {
		return nil, xerror.New("failed to marshal subscription into attribute values", err)
	}
	attributeValues["pk"] = &dynamodb.AttributeValue{
		S: aws.String("SUBSCRIPTION#" + subscription.EmailAddress),
	}
	attributeValues["sk"] = &dynamodb.AttributeValue{
		S: aws.String("SUBSCRIPTION#" + subscription.ID),
	}

	input := &dynamodb.PutItemInput{
		Item:      attributeValues,
		TableName: aws.String(TableName),
	}

	if err := input.Validate(); err != nil {
		return nil, xerror.New("failed to validate dynamodb.PutItemInput", err)
	}

	if _, err := DynamoDBClient.PutItemWithContext(ctx, input); err != nil {
		return nil, xerror.New("failed to put item", err)
	}

	return subscription, nil
}

func DeleteSubscription(ctx context.Context, id string, emailAddress string) error {
	if err := checkPackage(); err != nil {
		return xerror.New("checkPackage failed", err)
	}

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"pk": {
				S: aws.String("SUBSCRIPTION#" + emailAddress),
			},
			"sk": {
				S: aws.String("SUBSCRIPTION#" + id),
			},
		},
	}

	if err := input.Validate(); err != nil {
		return xerror.New("failed to validate dynamodb.DeleteItemInput", err)
	}

	if _, err := DynamoDBClient.DeleteItemWithContext(ctx, input); err != nil {
		return xerror.New("failed to delete subscription", err)
	}

	return nil
}

func GetSubscription(ctx context.Context, emailAddress string) (*Subscription, error) {
	if err := checkPackage(); err != nil {
		return nil, xerror.New("checkPackage failed", err)
	}

	if emailAddress == "" {
		return nil, xerror.Newf("emailAddress cannot be empty")
	}

	input := &dynamodb.QueryInput{
		KeyConditions: map[string]*dynamodb.Condition{
			"pk": {
				ComparisonOperator: aws.String(dynamodb.ComparisonOperatorEq),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String("SUBSCRIPTION#" + emailAddress),
					},
				},
			},
		},
		Limit:     aws.Int64(1),
		TableName: aws.String(TableName),
	}

	if err := input.Validate(); err != nil {
		return nil, xerror.New("failed to validate dynamodb.QueryInput", err)
	}

	output, err := DynamoDBClient.QueryWithContext(ctx, input)
	if err != nil {
		return nil, xerror.New("failed to query subscription", err)
	}

	if len(output.Items) == 0 {
		return nil, nil
	}

	var subscription *Subscription
	if err := dynamodbattribute.UnmarshalMap(output.Items[0], &subscription); err != nil {
		return nil, xerror.New("failed to unmarshal attribute values into Subscription", err)
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
	if s.EmailAddress == "" {
		return xerror.Newf("EmailAddress cannot be empty")
	}

	if s.CreatedAt == (time.Time{}) {
		return xerror.Newf("CreatedAt cannot be empty")
	}

	if s.UpdatedAt == (time.Time{}) {
		return xerror.Newf("UpdatedAt cannot be empty")
	}

	return nil
}
