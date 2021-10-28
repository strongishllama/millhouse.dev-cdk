package handler

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	email "github.com/gofor-little/aws-email"
	"github.com/gofor-little/log"
	"github.com/gofor-little/xlambda"

	"github.com/strongishllama/millhouse.dev-cdk/internal/db"
	"github.com/strongishllama/millhouse.dev-cdk/internal/notification"
)

const (
	EventInsert = "INSERT"
	EventModify = "MODIFY"
	EventRemove = "REMOVE"
)

var (
	FromAddress   string
	APIDomain     string
	WebsiteDomain string
)

func Handler(ctx context.Context, event *events.DynamoDBEvent) error {
	for _, r := range event.Records {
		if !strings.HasPrefix(r.Change.Keys["pk"].String(), "SUBSCRIPTION#") || r.EventName != EventInsert {
			continue
		}

		var subscription *db.Subscription
		if err := xlambda.UnmarshalDynamoDBEventAttributeValues(r.Change.NewImage, &subscription); err != nil {
			return fmt.Errorf("failed to unmarshal DynamoDB record into db.Subscription: %w", err)
		}

		if subscription.IsConfirmed {
			continue
		}

		_, err := notification.EnqueueEmail(ctx, []string{subscription.EmailAddress}, FromAddress, notification.EmailTemplate{
			FileName:    "email/subscription-confirmation.tmpl.html",
			Subject:     "Subscription Confirmation",
			ContentType: email.ContentTypeTextHTML,
			Data: notification.SubscriptionConfirmationTemplateData{
				WebsiteDomain:  WebsiteDomain,
				APIDomain:      APIDomain,
				SubscriptionID: subscription.ID,
				EmailAddress:   subscription.EmailAddress,
			},
		})
		if err != nil {
			log.Error(log.Fields{"error": err})
			return err
		}

		subscription.IsConfirmed = true
		if err := subscription.Update(ctx); err != nil {
			log.Error(log.Fields{"error": err})
			return err
		}
	}

	return nil
}
