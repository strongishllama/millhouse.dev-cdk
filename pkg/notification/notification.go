package notification

import (
	"context"
	"embed"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/gofor-little/xerror"
)

var (
	SQSClient *sqs.Client
	QueueURL  string

	//go:embed templates
	templates embed.FS
)

func Initialize(ctx context.Context, profile string, region string, queueURL string) error {
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

	SQSClient = sqs.NewFromConfig(cfg)
	QueueURL = queueURL

	return nil
}

func checkPackage() error {
	if SQSClient == nil {
		return xerror.Newf("notification.SQSClient is nil, have you called notification.Initialize()?")
	}

	if QueueURL == "" {
		return xerror.Newf("notification.QueueURL is empty, did you call notification.Initialize()?")
	}

	return nil
}
