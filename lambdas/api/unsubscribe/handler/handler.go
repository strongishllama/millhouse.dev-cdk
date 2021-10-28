package handler

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"net/http"
	"net/mail"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gofor-little/xlambda"

	"github.com/strongishllama/millhouse.dev-cdk/internal/db"
	"github.com/strongishllama/millhouse.dev-cdk/internal/tmpl"
)

var (
	//go:embed templates
	templates embed.FS
)

func Handler(ctx context.Context, request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	template, err := tmpl.NewTemplateFromFile(templates, "templates/unsubscribe-successful.tmpl.html", nil)
	if err != nil {
		return xlambda.ProxyResponseHTML(http.StatusInternalServerError, fmt.Errorf("failed to create template from file: %w", err), nil)
	}

	data := &RequestData{}
	if err := xlambda.ParseAndValidate(request, data); err != nil {
		return xlambda.ProxyResponseHTML(http.StatusBadRequest, err, nil)
	}

	if err := db.DeleteSubscription(ctx, data.EmailAddress, data.ID); err != nil {
		return xlambda.ProxyResponseHTML(http.StatusInternalServerError, fmt.Errorf("failed to delete subscription: %w", err), template)
	}

	return xlambda.ProxyResponseHTML(http.StatusOK, nil, template)
}

type RequestData struct {
	ID           string `mapstructure:"id"`
	EmailAddress string `mapstructure:"emailAddress"`
}

func (r *RequestData) Validate() error {
	if len(r.ID) == 0 {
		return errors.New("id cannot be empty")
	}
	if _, err := mail.ParseAddress(r.EmailAddress); err != nil {
		return fmt.Errorf("failed to validate EmailAddress: %w", err)
	}
	return nil
}
