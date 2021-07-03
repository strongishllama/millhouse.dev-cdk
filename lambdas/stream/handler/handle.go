package handler

import (
	"context"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	email "github.com/gofor-little/aws-email"
	"github.com/gofor-little/log"
	"github.com/gofor-little/xerror"

	"github.com/strongishllama/millhouse.dev-cdk/pkg/db"
	"github.com/strongishllama/millhouse.dev-cdk/pkg/notification"
	"github.com/strongishllama/millhouse.dev-cdk/pkg/xlambda"
)

const (
	EventInsert = "INSERT"
	EventModify = "MODIFY"
	EventRemove = "REMOVE"
)

var (
	AdminTo       string
	AdminFrom     string
	APIDomain     string
	WebsiteDomain string
)

func Handle(ctx context.Context, event *events.DynamoDBEvent) error {
	for _, r := range event.Records {
		// Only handle items that are subscriptions from INSERT and REMOVE events.
		if !strings.HasPrefix(r.Change.Keys["pk"].String(), "SUBSCRIPTION#") || r.EventName == EventModify {
			continue
		}

		// If it's a REMOVE event, enqueue an email notifying the admin of the reader unsubscribing.
		if r.EventName == EventRemove {
			emailAddress := strings.TrimPrefix(r.Change.Keys["pk"].String(), "SUBSCRIPTION#")
			messageID, err := notification.EnqueueEmail(ctx, []string{AdminTo}, AdminFrom, notification.EmailTemplate{
				FileName:    "email/reader-unsubscribed.tmpl.html",
				Subject:     "Reader Unsubscribed",
				ContentType: email.ContentTypeTextHTML,
				Data: notification.ReaderUnsubscribedTemplateData{
					EmailAddress: emailAddress,
				},
			})
			if err != nil {
				log.Error(log.Fields{
					"error":        xerror.Wrap("failed to enqueue reader unsubscribed email", err),
					"messageId":    messageID,
					"emailAddress": emailAddress,
				})
			}
		}

		// If it's a INSERT event, unmarshal the subscription and enqueue a subscription confirmation email.
		var subscription *db.Subscription
		if err := xlambda.UnmarshalDynamoDBImage(r.Change.NewImage, &subscription); err != nil {
			return xerror.Wrap("failed to unmarshal DynamoDB record into db.Subscription", err)
		}

		messageID, err := notification.EnqueueEmail(ctx, []string{subscription.EmailAddress}, AdminFrom, notification.EmailTemplate{
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
			return xerror.Wrap("failed to enqueue subscription confirmation email", err)
		}
		log.Info(log.Fields{
			"message":      "successfully enqueued subscription confirmation email",
			"messageID":    messageID,
			"emailAddress": subscription,
		})
	}

	return nil
}
