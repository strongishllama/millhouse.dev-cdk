package handler

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	email "github.com/gofor-little/aws-email"
	"github.com/gofor-little/log"
	"github.com/gofor-little/xlambda"

	"github.com/strongishllama/millhouse.dev-cdk/pkg/db"
	"github.com/strongishllama/millhouse.dev-cdk/pkg/notification"
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
		// Only handle items that are subscriptions from INSERT.
		if !strings.HasPrefix(r.Change.Keys["pk"].String(), "SUBSCRIPTION#") || r.EventName != EventInsert {
			continue
		}

		var subscription *db.Subscription
		if err := xlambda.UnmarshalDynamoDBEventAttributeValues(r.Change.NewImage, &subscription); err != nil {
			return fmt.Errorf("failed to unmarshal DynamoDB record into db.Subscription: %w", err)
		}

		messageID, err := notification.EnqueueEmail(ctx, []string{subscription.EmailAddress}, FromAddress, notification.EmailTemplate{
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
			return fmt.Errorf("failed to enqueue subscription confirmation email: %w", err)
		}
		log.Info(log.Fields{
			"message":      "successfully enqueued subscription confirmation email",
			"messageID":    messageID,
			"emailAddress": subscription,
		})
	}

	return nil
}
