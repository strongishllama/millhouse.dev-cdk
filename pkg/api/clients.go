package api

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"net/http"
	"text/template"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	email "github.com/gofor-little/aws-email"
	"github.com/gofor-little/xerror"
)

var (
	HTTPClient httpClient
	SQSClient  sqsiface.SQSAPI

	//go:embed templates
	templates embed.FS
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func Initialize(profile string, region string) error {
	HTTPClient = &http.Client{
		Timeout: time.Duration(3 * time.Second),
	}

	var sess *session.Session
	var err error

	if profile != "" && region != "" {
		sess, err = session.NewSessionWithOptions(session.Options{
			Config: aws.Config{
				Region: aws.String(region),
			},
			Profile: profile,
		})
	} else {
		sess, err = session.NewSession()
	}
	if err != nil {
		return xerror.New("failed to create session.Session", err)
	}

	SQSClient = sqs.New(sess)

	return nil
}

func sendEmail(ctx context.Context, queueURL string, to []string, from string, subject string, bodyTemplate string, contentType email.ContentType) (string, error) {
	tmpl := template.New("template")

	data, err := templates.ReadFile(bodyTemplate)
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
