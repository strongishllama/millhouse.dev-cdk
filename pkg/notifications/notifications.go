package notifications

import (
	"embed"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

var (
	SQSClient sqsiface.SQSAPI

	//go:embed templates
	templates embed.FS
)

func Initialize(profile string, region string) error {
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

	return nil
}
