package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gofor-little/cfg"
	"github.com/gofor-little/env"
	"github.com/gofor-little/log"
	"github.com/gofor-little/xerror"

	"github.com/strongishllama/millhouse.dev-cdk/lambdas/api/subscribe/handler"
	"github.com/strongishllama/millhouse.dev-cdk/pkg/notifications"
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

	if err := cfg.Initialize("", ""); err != nil {
		log.Error(log.Fields{
			"error": xerror.New("failed to initialize cfg package", err),
		})
		os.Exit(1)
	}

	handler.Cfg = &handler.Config{}
	if err := cfg.Load(context.Background(), configSecretARN, handler.Cfg); err != nil {
		log.Error(log.Fields{
			"error": xerror.New("failed to load config", err),
		})
		os.Exit(1)
	}

	if err := notifications.Initialize("", ""); err != nil {
		log.Error(log.Fields{
			"error": xerror.New("failed to initialize the notifications package", err),
		})
		os.Exit(1)
	}

	lambda.Start(handler.Handle)
}
