package notifications

import (
	"bytes"
	"context"
	"encoding/json"
	"text/template"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	email "github.com/gofor-little/aws-email"
	"github.com/gofor-little/xerror"
)

func SendEmail(ctx context.Context, queueURL string, to []string, from string, subject string, bodyTemplate string, contentType email.ContentType) (string, error) {
	tmpl := template.New("template")

	data, err := templates.ReadFile("templates/" + bodyTemplate)
	if err != nil {
		return "", xerror.New("failed to read file data", err)
	}

	tmpl, err = tmpl.Parse(string(data))
	if err != nil {
		return "", xerror.New("failed to parse template", err)
	}

	buff := &bytes.Buffer{}
	if err := tmpl.Execute(buff, nil); err != nil {
		return "", xerror.New("failed to execute template", err)
	}

	body, err := json.Marshal(&email.Data{
		To:          to,
		From:        from,
		Subject:     subject,
		Body:        buff.String(),
		ContentType: contentType,
	})
	if err != nil {
		return "", xerror.New("failed to marshal email.Data", err)
	}

	input := &sqs.SendMessageInput{
		MessageBody: aws.String(string(body)),
		QueueUrl:    aws.String(queueURL),
	}

	if err := input.Validate(); err != nil {
		return "", xerror.New("failed to validate sqs.SendMessageInput", err)
	}

	output, err := SQSClient.SendMessageWithContext(ctx, input)
	if err != nil {
		return "", xerror.New("failed to send message to SQS", err)
	}

	return *output.MessageId, nil
}
