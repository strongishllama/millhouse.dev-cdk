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
	"github.com/strongishllama/millhouse.dev-cdk/pkg/db"
	"github.com/strongishllama/millhouse.dev-cdk/pkg/notification"
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

	emailQueueURL, err := env.MustGet("EMAIL_QUEUE_URL")
	if err != nil {
		log.Error(log.Fields{"error": xerror.New("failed to get environment variable", err)})
		os.Exit(1)
	}
	if err := notification.Initialize("", "", emailQueueURL); err != nil {
		log.Error(log.Fields{"error": xerror.New("failed to initialize the notifications package", err)})
		os.Exit(1)
	}

	handler.AdminTo, err = env.MustGet("ADMIN_TO")
	if err != nil {
		log.Error(log.Fields{"error": xerror.New("failed to get environment variable", err)})
		os.Exit(1)
	}
	handler.AdminFrom, err = env.MustGet("ADMIN_FROM")
	if err != nil {
		log.Error(log.Fields{"error": xerror.New("failed to get environment variable", err)})
		os.Exit(1)
	}

	configSecretARN, err := env.MustGet("CONFIG_SECRET_ARN")
	if err != nil {
		log.Error(log.Fields{"error": xerror.New("failed to environment variable", err)})
		os.Exit(1)
	}
	if err := cfg.Initialize("", ""); err != nil {
		log.Error(log.Fields{"error": xerror.New("failed to initialize the cfg package", err)})
		os.Exit(1)
	}

	config := &struct {
		RecaptchaSecret string `json:"recaptchaSecret"`
	}{}
	if err := cfg.Load(context.Background(), configSecretARN, config); err != nil {
		log.Error(log.Fields{"error": xerror.New("failed to load config", err)})
		os.Exit(1)
	}
	handler.RecaptchaSecret = config.RecaptchaSecret

	lambda.Start(handler.Handle)
}
