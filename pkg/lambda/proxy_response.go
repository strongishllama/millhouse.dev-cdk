package lambda

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gofor-little/log"
	"github.com/gofor-little/xerror"
)

// NewProxyResponse builds an API gateway proxy response using the passed parameters. statusCode should
// be a valid HTTP status code. v should be a struct that is able to be marshaled into JSON. If err is
// nil no error will be returned. If v is nil nothing will be written to the response body.
func NewProxyResponse(statusCode int, err error, v interface{}) (*events.APIGatewayProxyResponse, error) {
	response := &events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type": "application/json",
			// "Access-Control-Expose-Headers":    "*",
			// "Access-Control-Allow-Headers":     "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token,X-Amz-User-Agent",
			// "Access-Control-Allow-Origin":      "*",
			// "Access-Control-Allow-Credentials": "true",
			// "Access-Control-Allow-Methods":     "OPTIONS,GET,PUT,POST,DELETE,PATCH,HEAD",
		},
		StatusCode: statusCode,
	}

	if v != nil {
		body, e := json.Marshal(v)
		if e != nil {
			log.Error(log.Fields{
				"error":      xerror.Newf("failed to marshal response body and API request failed", e, err),
				"statusCode": statusCode,
			})
			return response, xerror.Newf("failed to marshal response body and API request failed", e, err)
		}

		response.Body = string(body)
	}

	if err != nil {
		log.Error(log.Fields{
			"error":      xerror.New("api request failed", err),
			"statusCode": statusCode,
		})
		return response, xerror.New("api request failed", err)
	}

	log.Info(log.Fields{
		"message":    "api request succeeded",
		"statusCode": response.StatusCode,
	})

	return response, nil
}
