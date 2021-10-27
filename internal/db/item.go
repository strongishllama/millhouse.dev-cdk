package db

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type itemType string

const (
	itemTypeSubscription itemType = "SUBSCRIPTION"
	itemTypeCount        itemType = "COUNT"
)

// item represents an item in the DynamoDB table. If implementing this interface,
// be sure to add the 'dynamodbav' tags to the struct's properties.
type item interface {
	// Returns the primary key of the item.
	pk() string
	// Returns the sort key of the item.
	sk() string
	// Returns the primary key of the secondary item that keeps a count of the main item.
	countPK() string
	// Returns the sort key of the secondary item that keeps a count of the main item.
	countSK() string
	// Returns the type of the item.
	itemType() itemType
	// Returns the expression to update the item.
	updateExpression() (expression.Expression, error)
	// Validates the items properties.
	validate() error
}

// deleteItem deletes an item based on its primary key and sort key from the
// DynamoDB table.
func deleteItem(ctx context.Context, it itemType, pk string, sk string) error {
	if err := checkPackage(); err != nil {
		return err
	}

	// Delete the item in a transaction so we can update a secondary item that tracks the
	// number of this type of item in the DynamoDB table.
	if _, err := DynamoDBClient.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{
				Delete: &types.Delete{
					Key: map[string]types.AttributeValue{
						"pk": &types.AttributeValueMemberS{Value: pk},
						"sk": &types.AttributeValueMemberS{Value: sk},
					},
					TableName: aws.String(TableName),
				},
			},
			{
				Update: &types.Update{
					Key: map[string]types.AttributeValue{
						"pk": &types.AttributeValueMemberS{Value: string(itemTypeCount)},
						"sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("%s#%s", itemTypeCount, it)},
					},
					TableName:        aws.String(TableName),
					UpdateExpression: aws.String("ADD #count :count"),
					ExpressionAttributeNames: map[string]string{
						"#count": "count",
					},
					ExpressionAttributeValues: map[string]types.AttributeValue{
						":count": &types.AttributeValueMemberN{
							Value: "-1",
						},
					},
				},
			},
		},
	}); err != nil {
		return err
	}

	return nil
}

// getItem fetches a single item based on its primary key and sort key from the
// DynamoDB table. The 'item' parameter must be a non-nil pointer to a slice of
// elements that implement the item interface.
func getItem(ctx context.Context, pk string, sk string, item interface{}) error {
	if err := checkPackage(); err != nil {
		return err
	}

	key := map[string]types.AttributeValue{
		"pk": &types.AttributeValueMemberS{Value: pk},
	}
	if len(sk) != 0 {
		key["sk"] = &types.AttributeValueMemberS{Value: sk}
	}

	output, err := DynamoDBClient.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(TableName),
		Key:       key,
	})
	if err != nil {
		return err
	}

	if output.Item == nil {
		return nil
	}

	if err := attributevalue.UnmarshalMap(output.Item, &item); err != nil {
		return fmt.Errorf("failed to unmarshal item into interface: %w", err)
	}

	return nil
}

// getItems fetches a slice of items based on their itemType from the DynamoDB
// table. The 'items' parameter must be a non-nil pointer to a slice of elements
// that implement the item interface and have their itemType equal the itemType
// of the 'it' parameter.
func getItems(ctx context.Context, it itemType, items interface{}) error {
	if err := checkPackage(); err != nil {
		return err
	}

	expr, err := expression.NewBuilder().WithKeyCondition(expression.Key("gsiPk1").Equal(expression.Value(it))).Build()
	if err != nil {
		return err
	}

	output, err := DynamoDBClient.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(TableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		IndexName:                 aws.String("Gsi1"),
		KeyConditionExpression:    expr.KeyCondition(),
	})
	if err != nil {
		return err
	}

	if err := attributevalue.UnmarshalListOfMaps(output.Items, &items); err != nil {
		return fmt.Errorf("failed to unmarshal items into slice: %w", err)
	}

	return nil
}

// putItem inserts a new item into the DynamoDB table.
func putItem(ctx context.Context, i item) error {
	if err := checkPackage(); err != nil {
		return err
	}

	if err := i.validate(); err != nil {
		return fmt.Errorf("failed to validate %s: %w", i.itemType(), err)
	}

	attributeValues, err := attributevalue.MarshalMap(i)
	if err != nil {
		return fmt.Errorf("failed to marshal account into attribute values: %w", err)
	}
	attributeValues["pk"] = &types.AttributeValueMemberS{Value: i.pk()}
	attributeValues["sk"] = &types.AttributeValueMemberS{Value: i.sk()}
	attributeValues["itemType"] = &types.AttributeValueMemberS{Value: string(i.itemType())}

	// Create the item in a transaction so we can update a secondary item that tracks the
	// number of this type of item in the DynamoDB table.
	if _, err := DynamoDBClient.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{
				Put: &types.Put{
					Item:      attributeValues,
					TableName: aws.String(TableName),
				},
			},
			{
				Update: &types.Update{
					Key: map[string]types.AttributeValue{
						"pk": &types.AttributeValueMemberS{Value: i.countPK()},
						"sk": &types.AttributeValueMemberS{Value: i.countSK()},
					},
					TableName:        aws.String(TableName),
					UpdateExpression: aws.String("ADD #count :count"),
					ExpressionAttributeNames: map[string]string{
						"#count": "count",
					},
					ExpressionAttributeValues: map[string]types.AttributeValue{
						":count": &types.AttributeValueMemberN{
							Value: "1",
						},
					},
				},
			},
		},
	}); err != nil {
		return err
	}

	return nil
}

// updateItem updates an existing item in the DynamoDB table.
func updateItem(ctx context.Context, i item) error {
	if err := checkPackage(); err != nil {
		return err
	}

	if err := i.validate(); err != nil {
		return fmt.Errorf("failed to validate %s: %w", i.itemType(), err)
	}

	expr, err := i.updateExpression()
	if err != nil {
		return err
	}

	if _, err := DynamoDBClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: i.pk()},
			"sk": &types.AttributeValueMemberS{Value: i.sk()},
		},
		TableName:                 aws.String(TableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
	}); err != nil {
		return err
	}

	return nil
}
