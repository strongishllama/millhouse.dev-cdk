package db

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/gofor-little/xerror"
)

var (
	DynamoDBClient dynamodbiface.DynamoDBAPI
	TableName      string
)

func Initialize(profile string, region string, tableName string) error {
	var sess *session.Session
	var err error

	if profile != "" && region != "" {
		sess, err = session.NewSessionWithOptions(session.Options{
			Config: aws.Config{
				Region: aws.String(region),
			},
			Profile: profile,
		})
	} else {
		sess, err = session.NewSession()
	}
	if err != nil {
		return fmt.Errorf("failed to create session.Session: %w", err)
	}

	DynamoDBClient = dynamodb.New(sess)
	TableName = tableName

	return nil
}

func checkPackage() error {
	if DynamoDBClient == nil {
		return xerror.Newf("db.DynamoDBClient is nil, have you called db.Initialize()?")
	}

	if TableName == "" {
		return xerror.Newf("db.TableName is empty, did you call db.Initialize()?")
	}

	return nil
}
