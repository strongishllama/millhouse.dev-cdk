package db

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gofor-little/xerror"
	"github.com/gofor-little/xrand"
)

func CreateSubscription(ctx context.Context, emailAddress string) (*Subscription, error) {
	if err := checkPackage(); err != nil {
		return nil, xerror.Wrap("checkPackage failed", err)
	}

	id, err := xrand.UUIDV4()
	if err != nil {
		return nil, xerror.Wrap("failed to generate UUID", err)
	}
	subscription := &Subscription{
		ID:           id,
		EmailAddress: emailAddress,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := subscription.Validate(); err != nil {
		return nil, xerror.Wrap("failed to validate subscription", err)
	}

	attributeValues, err := attributevalue.MarshalMap(subscription)
	if err != nil {
		return nil, xerror.Wrap("failed to marshal subscription into attribute values", err)
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
		return nil, xerror.Wrap("failed to put item", err)
	}

	return subscription, nil
}

func DeleteSubscription(ctx context.Context, id string, emailAddress string) error {
	if err := checkPackage(); err != nil {
		return xerror.Wrap("checkPackage failed", err)
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
		return xerror.Wrap("failed to delete subscription", err)
	}

	return nil
}

func GetSubscription(ctx context.Context, emailAddress string) (*Subscription, error) {
	if err := checkPackage(); err != nil {
		return nil, xerror.Wrap("checkPackage failed", err)
	}

	if emailAddress == "" {
		return nil, xerror.New("emailAddress cannot be empty")
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
		return nil, xerror.Wrap("failed to query subscription", err)
	}

	if len(output.Items) == 0 {
		return nil, nil
	}

	var subscription *Subscription
	if err := attributevalue.UnmarshalMap(output.Items[0], &subscription); err != nil {
		return nil, xerror.Wrap("failed to unmarshal attribute values into Subscription", err)
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
		return xerror.New("EmailAddress cannot be empty")
	}

	if s.CreatedAt == (time.Time{}) {
		return xerror.New("CreatedAt cannot be empty")
	}

	if s.UpdatedAt == (time.Time{}) {
		return xerror.New("UpdatedAt cannot be empty")
	}

	return nil
}
