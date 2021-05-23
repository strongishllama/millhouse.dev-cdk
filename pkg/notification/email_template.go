package notification

import email "github.com/gofor-little/aws-email"

type EmailTemplate struct {
	FileName    string
	Subject     string
	ContentType email.ContentType
	Data        interface{}
}

type SubscriptionConfirmationTemplateData struct {
	WebsiteDomain  string
	APIDomain      string
	SubscriptionID string
	EmailAddress   string
}

type ReaderUnsubscribedTemplateData struct {
	EmailAddress string
}

type RecaptchaChallengeFailedTemplateData struct {
	EmailAddress string
	Score        float32
}
