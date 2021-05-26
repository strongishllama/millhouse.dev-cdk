package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gofor-little/env"
	"github.com/gofor-little/log"
	"github.com/gofor-little/xerror"

	"github.com/strongishllama/millhouse.dev-cdk/lambdas/api/unsubscribe/handler"
	"github.com/strongishllama/millhouse.dev-cdk/pkg/db"
)

func main() {
	log.Log = log.NewStandardLogger(os.Stdout, nil)

	tableName, err := env.MustGet("TABLE_NAME")
	if err != nil {
		log.Error(log.Fields{"error": xerror.New("failed to get environment variable", err)})
		os.Exit(1)
	}
	if err := db.Initialize("", "", tableName); err != nil {
		log.Error(log.Fields{"error": xerror.New("failed to initialize the db package", err)})
		os.Exit(1)
	}

	lambda.Start(handler.Handle)
}
