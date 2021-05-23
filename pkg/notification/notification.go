package notification

import (
	"embed"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/gofor-little/xerror"
)

var (
	SQSClient sqsiface.SQSAPI
	QueueURL  string

	//go:embed templates
	templates embed.FS
)

func Initialize(profile string, region string, queueURL string) error {
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
		return fmt.Errorf("failed to create session.Session: %w", err)
	}

	SQSClient = sqs.New(sess)
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
