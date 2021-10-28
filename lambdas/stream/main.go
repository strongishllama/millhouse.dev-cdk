package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gofor-little/env"
	"github.com/gofor-little/log"

	"github.com/strongishllama/millhouse.dev-cdk/internal/db"
	"github.com/strongishllama/millhouse.dev-cdk/internal/notification"
	"github.com/strongishllama/millhouse.dev-cdk/lambdas/stream/handler"
)

func main() {
	log.Log = log.NewStandardLogger(os.Stdout, nil)

	if err := notification.Initialize(context.Background(), "", "", env.Get("EMAIL_QUEUE_URL", "")); err != nil {
		log.Error(log.Fields{"error": fmt.Errorf("failed to initialize the notification package: %w", err)})
		os.Exit(1)
	}

	if err := db.Initialize(context.Background(), "", "", env.Get("TABLE_NAME", "")); err != nil {
		log.Error(log.Fields{"error": fmt.Errorf("failed to initialize the db package: %w", err)})
		os.Exit(1)
	}

	var err error
	handler.FromAddress, err = env.MustGet("FROM_ADDRESS")
	if err != nil {
		log.Error(log.Fields{"error": err})
		os.Exit(1)
	}
	handler.APIDomain, err = env.MustGet("API_DOMAIN")
	if err != nil {
		log.Error(log.Fields{"error": err})
		os.Exit(1)
	}
	handler.WebsiteDomain, err = env.MustGet("WEBSITE_DOMAIN")
	if err != nil {
		log.Error(log.Fields{"error": err})
		os.Exit(1)
	}

	lambda.Start(handler.Handler)
}
