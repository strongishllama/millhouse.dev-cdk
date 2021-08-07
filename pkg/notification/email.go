package notification

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	email "github.com/gofor-little/aws-email"
	"github.com/gofor-little/xerror"

	"github.com/strongishllama/millhouse.dev-cdk/pkg/tmpl"
)

func EnqueueEmail(ctx context.Context, to []string, from string, emailTemplate EmailTemplate) (string, error) {
	if err := checkPackage(); err != nil {
		return "", xerror.Wrap("checkPackage failed", err)
	}

	data, err := tmpl.NewTemplateFromFile(templates, "templates/"+emailTemplate.FileName, emailTemplate.Data)
	if err != nil {
		return "", xerror.Wrap("failed to create template from file", err)
	}

	body, err := json.Marshal(&email.Data{
		To:          to,
		From:        from,
		Subject:     emailTemplate.Subject,
		Body:        string(data),
		ContentType: emailTemplate.ContentType,
	})
	if err != nil {
		return "", xerror.Wrap("failed to marshal email.Data", err)
	}

	input := &sqs.SendMessageInput{
		MessageBody: aws.String(string(body)),
		QueueUrl:    aws.String(QueueURL),
	}

	output, err := SQSClient.SendMessage(ctx, input)
	if err != nil {
		return "", xerror.Wrap("failed to send message to SQS", err)
	}

	return *output.MessageId, nil
}
