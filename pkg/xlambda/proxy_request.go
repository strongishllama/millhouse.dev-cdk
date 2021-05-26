package xlambda

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

func NewProxyRequest(method string, queryParameters map[string]string, body interface{}) (*events.APIGatewayProxyRequest, error) {
	request := &events.APIGatewayProxyRequest{
		HTTPMethod:            method,
		QueryStringParameters: queryParameters,
	}

	if body == nil {
		return request, nil
	}

	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	request.Body = string(data)

	return request, nil
}
