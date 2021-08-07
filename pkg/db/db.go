package db

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gofor-little/xerror"
)

var (
	DynamoDBClient *dynamodb.Client
	TableName      string
)

func Initialize(ctx context.Context, profile string, region string, tableName string) error {
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
	TableName = tableName

	return nil
}

func checkPackage() error {
	if DynamoDBClient == nil {
		return xerror.New("db.DynamoDBClient is nil, have you called db.Initialize()?")
	}

	if TableName == "" {
		return xerror.New("db.TableName is empty, did you call db.Initialize()?")
	}

	return nil
}
