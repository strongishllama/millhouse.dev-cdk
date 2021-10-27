package notification

import (
	"context"
	"embed"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

var (
	SQSClient *sqs.Client
	QueueURL  string

	//go:embed templates
	templates embed.FS
)

func Initialize(ctx context.Context, profile string, region string, queueURL string) error {
	if len(queueURL) == 0 {
		return errors.New("queue url cannot be empty")
	}
	QueueURL = queueURL

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

	return nil
}

func checkPackage() error {
	if SQSClient == nil {
		return errors.New("notification.SQSClient is nil, have you called notification.Initialize()?")
	}

	if QueueURL == "" {
		return errors.New("notification.QueueURL is empty, did you call notification.Initialize()?")
	}

	return nil
}
