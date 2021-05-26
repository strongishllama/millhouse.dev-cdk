package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gofor-little/env"
	"github.com/gofor-little/log"
	"github.com/gofor-little/xerror"

	"github.com/strongishllama/millhouse.dev-cdk/lambdas/stream/handler"
	"github.com/strongishllama/millhouse.dev-cdk/pkg/notification"
)

func main() {
	log.Log = log.NewStandardLogger(os.Stdout, nil)

	emailQueueURL, err := env.MustGet("EMAIL_QUEUE_URL")
	if err != nil {
		log.Error(log.Fields{"error": xerror.New("failed to get environment variable", err)})
		os.Exit(1)
	}
	if err := notification.Initialize("", "", emailQueueURL); err != nil {
		log.Error(log.Fields{"error": xerror.New("failed to initialize the notification package", err)})
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
	handler.APIDomain, err = env.MustGet("API_DOMAIN")
	if err != nil {
		log.Error(log.Fields{"error": xerror.New("failed to get environment variable", err)})
		os.Exit(1)
	}
	handler.WebsiteDomain, err = env.MustGet("WEBSITE_DOMAIN")
	if err != nil {
		log.Error(log.Fields{"error": xerror.New("failed to get environment variable", err)})
		os.Exit(1)
	}

	lambda.Start(handler.Handle)
}
