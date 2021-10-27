package notification

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	email "github.com/gofor-little/aws-email"

	"github.com/strongishllama/millhouse.dev-cdk/pkg/tmpl"
)

func EnqueueEmail(ctx context.Context, to []string, from string, emailTemplate EmailTemplate) (string, error) {
	if err := checkPackage(); err != nil {
		return "", err
	}

	data, err := tmpl.NewTemplateFromFile(templates, "templates/"+emailTemplate.FileName, emailTemplate.Data)
	if err != nil {
		return "", fmt.Errorf("failed to create template from file: %w", err)
	}

	body, err := json.Marshal(&email.Data{
		To:          to,
		From:        from,
		Subject:     emailTemplate.Subject,
		Body:        string(data),
		ContentType: emailTemplate.ContentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal email.Data: %w", err)
	}

	input := &sqs.SendMessageInput{
		MessageBody: aws.String(string(body)),
		QueueUrl:    aws.String(QueueURL),
	}

	output, err := SQSClient.SendMessage(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to send message to SQS: %w", err)
	}

	return *output.MessageId, nil
}
