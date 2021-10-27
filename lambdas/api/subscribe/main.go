package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gofor-little/cfg"
	"github.com/gofor-little/env"
	"github.com/gofor-little/log"
	"github.com/gofor-little/xlambda"

	"github.com/strongishllama/millhouse.dev-cdk/lambdas/api/subscribe/handler"
	"github.com/strongishllama/millhouse.dev-cdk/pkg/db"
	"github.com/strongishllama/millhouse.dev-cdk/pkg/notification"
)

func main() {
	log.Log = log.NewStandardLogger(os.Stdout, nil)

	if err := db.Initialize(context.Background(), "", "", env.Get("TABLE_NAME", "")); err != nil {
		log.Error(log.Fields{"error": fmt.Errorf("failed to initialize the db package: %w", err)})
		os.Exit(1)
	}

	if err := notification.Initialize(context.Background(), "", "", env.Get("EMAIL_QUEUE_URL", "")); err != nil {
		log.Error(log.Fields{"error": fmt.Errorf("failed to initialize the notifications package: %w", err)})
		os.Exit(1)
	}

	if err := xlambda.Initialize(env.Get("ACCESS_CONTROL_ALLOW_ORIGIN", "*")); err != nil {
		log.Error(log.Fields{"error": fmt.Errorf("failed to initialize the xlambda package: %w", err)})
		os.Exit(1)
	}

	if err := cfg.Initialize(context.Background(), "", ""); err != nil {
		log.Error(log.Fields{"error": fmt.Errorf("failed to initialize the cfg package: %w", err)})
		os.Exit(1)
	}

	var err error
	handler.RecaptchaSecret, err = cfg.LoadString(context.Background(), env.Get("RECAPTCHA_SECRET_ARN", ""))
	if err != nil {
		log.Error(log.Fields{"error": fmt.Errorf("failed to load recaptcha secret: %w", err)})
		os.Exit(1)
	}

	lambda.Start(handler.Handler)
}
