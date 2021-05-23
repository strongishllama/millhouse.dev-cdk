package handler

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	email "github.com/gofor-little/aws-email"
	"github.com/gofor-little/log"
	"github.com/gofor-little/xerror"

	"github.com/strongishllama/millhouse.dev-cdk/pkg/db"
	"github.com/strongishllama/millhouse.dev-cdk/pkg/lambda"
	"github.com/strongishllama/millhouse.dev-cdk/pkg/notification"
	"github.com/strongishllama/millhouse.dev-cdk/pkg/tmpl"
)

func Handle(ctx context.Context, request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	data, err := tmpl.NewTemplateFromFile(templates, "templates/unsubscribe-successful.tmpl.html", nil)
	if err != nil {
		return lambda.NewProxyResponse(http.StatusInternalServerError, lambda.ContentTypeTextHTML, xerror.New("failed to create template from file", err), nil)
	}

	id, ok := request.QueryStringParameters["id"]
	if !ok {
		return lambda.NewProxyResponse(http.StatusBadRequest, lambda.ContentTypeTextHTML, xerror.Newf("required query parameter id is missing"), data)
	}
	emailAddress, ok := request.QueryStringParameters["emailAddress"]
	if !ok {
		return lambda.NewProxyResponse(http.StatusBadRequest, lambda.ContentTypeTextHTML, xerror.Newf("required query parameter emailAddress is missing"), data)
	}

	requestData := &RequestData{
		ID:           id,
		EmailAddress: emailAddress,
	}
	if err := requestData.Validate(); err != nil {
		return lambda.NewProxyResponse(http.StatusBadRequest, lambda.ContentTypeTextHTML, xerror.New("failed to validate request data", err), data)
	}

	if err := db.DeleteSubscription(ctx, requestData.ID, requestData.EmailAddress); err != nil {
		return lambda.NewProxyResponse(http.StatusInternalServerError, lambda.ContentTypeTextHTML, xerror.New("failed to delete subscription", err), data)
	}

	messageID, err := notification.EnqueueEmail(ctx, []string{Cfg.To}, Cfg.From, notification.EmailTemplate{
		FileName:    "email/reader-unsubscribed.tmpl.html",
		Subject:     "Reader Unsubscribed",
		ContentType: email.ContentTypeTextHTML,
		Data: notification.ReaderUnsubscribedTemplateData{
			EmailAddress: requestData.EmailAddress,
		},
	})
	if err != nil {
		log.Error(log.Fields{
			"error":        xerror.New("failed to enqueue reader unsubscribed email", err),
			"messageId":    messageID,
			"emailAddress": requestData.EmailAddress,
		})
	}

	return lambda.NewProxyResponse(http.StatusOK, lambda.ContentTypeTextHTML, nil, data)
}
