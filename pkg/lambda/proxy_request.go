package lambda

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

func NewProxyRequest(method string, v interface{}) (*events.APIGatewayProxyRequest, error) {
	request := &events.APIGatewayProxyRequest{
		HTTPMethod: method,
	}

	if v == nil {
		return request, nil
	}

	body, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	request.Body = string(body)

	return request, nil
}
