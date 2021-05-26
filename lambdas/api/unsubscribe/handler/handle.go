package handler

import (
	"context"
	"embed"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gofor-little/xerror"

	"github.com/strongishllama/millhouse.dev-cdk/pkg/db"
	"github.com/strongishllama/millhouse.dev-cdk/pkg/tmpl"
	"github.com/strongishllama/millhouse.dev-cdk/pkg/xlambda"
)

var (
	//go:embed templates
	templates embed.FS
)

func Handle(ctx context.Context, request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	if request.HTTPMethod != http.MethodGet {
		return xlambda.NewProxyResponse(http.StatusMethodNotAllowed, "", nil, nil)
	}

	// Fetch the unsubscibe template.
	data, err := tmpl.NewTemplateFromFile(templates, "templates/unsubscribe-successful.tmpl.html", nil)
	if err != nil {
		return xlambda.NewProxyResponse(http.StatusInternalServerError, xlambda.ContentTypeTextHTML, xerror.New("failed to create template from file", err), nil)
	}

	// Pull required query parameters.
	id, ok := request.QueryStringParameters["id"]
	if !ok {
		return xlambda.NewProxyResponse(http.StatusBadRequest, xlambda.ContentTypeTextHTML, xerror.Newf("required query parameter id is missing"), data)
	}
	emailAddress, ok := request.QueryStringParameters["emailAddress"]
	if !ok {
		return xlambda.NewProxyResponse(http.StatusBadRequest, xlambda.ContentTypeTextHTML, xerror.Newf("required query parameter emailAddress is missing"), data)
	}

	// Validate request data.
	requestData := &RequestData{
		ID:           id,
		EmailAddress: emailAddress,
	}
	if err := requestData.Validate(); err != nil {
		return xlambda.NewProxyResponse(http.StatusBadRequest, xlambda.ContentTypeTextHTML, xerror.New("failed to validate request data", err), data)
	}

	// Delete the subscription.
	if err := db.DeleteSubscription(ctx, requestData.ID, requestData.EmailAddress); err != nil {
		return xlambda.NewProxyResponse(http.StatusInternalServerError, xlambda.ContentTypeTextHTML, xerror.New("failed to delete subscription", err), data)
	}

	return xlambda.NewProxyResponse(http.StatusOK, xlambda.ContentTypeTextHTML, nil, data)
}
