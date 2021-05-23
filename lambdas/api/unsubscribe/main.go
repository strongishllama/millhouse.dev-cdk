package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gofor-little/env"
	"github.com/gofor-little/log"
	"github.com/gofor-little/xerror"

	"github.com/strongishllama/millhouse.dev-cdk/lambdas/api/unsubscribe/handler"
)

func main() {
	log.Log = log.NewStandardLogger(os.Stdout, nil)

	configSecretARN, err := env.MustGet("CONFIG_SECRET_ARN")
	if err != nil {
		log.Error(log.Fields{
			"error": xerror.New("failed to environment variable", err),
		})
		os.Exit(1)
	}

	if err := handler.Initialize(configSecretARN); err != nil {
		log.Error(log.Fields{
			"error": xerror.New("failed to initialize handler package", err),
		})
		os.Exit(1)
	}

	lambda.Start(handler.Handle)
}
