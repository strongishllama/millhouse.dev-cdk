package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var (
	DynamoDBClient *dynamodb.Client
	TableName      string
)

func Initialize(ctx context.Context, profile string, region string, tableName string) error {
	if len(tableName) == 0 {
		return errors.New("table name cannot be empty")
	}
	TableName = tableName

	var cfg aws.Config
	var err error

	if profile != "" && region != "" {
		cfg, err = config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(profile), config.WithRegion(region))
	} else {
		cfg, err = config.LoadDefaultConfig(ctx)
	}
	if err != nil {
		return fmt.Errorf("failed to load default config: %w", err)
	}

	DynamoDBClient = dynamodb.NewFromConfig(cfg)

	return nil
}

func checkPackage() error {
	if DynamoDBClient == nil {
		return errors.New("db.DynamoDBClient is nil, have you called db.Initialize()?")
	}

	if TableName == "" {
		return errors.New("db.TableName is empty, did you call db.Initialize()?")
	}

	return nil
}
