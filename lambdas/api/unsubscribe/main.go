package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gofor-little/env"
	"github.com/gofor-little/log"
	"github.com/gofor-little/xlambda"

	"github.com/strongishllama/millhouse.dev-cdk/internal/db"
	"github.com/strongishllama/millhouse.dev-cdk/lambdas/api/unsubscribe/handler"
)

func main() {
	log.Log = log.NewStandardLogger(os.Stdout, nil)

	if err := db.Initialize(context.Background(), "", "", env.Get("TABLE_NAME", "")); err != nil {
		log.Error(log.Fields{"error": fmt.Errorf("failed to initialize the db package: %w", err)})
		os.Exit(1)
	}

	if err := xlambda.Initialize(env.Get("ACCESS_CONTROL_ALLOW_ORIGIN", "*")); err != nil {
		log.Error(log.Fields{"error": fmt.Errorf("failed to initialize the xlambda package; %w", err)})
		os.Exit(1)
	}

	lambda.Start(handler.Handler)
}
