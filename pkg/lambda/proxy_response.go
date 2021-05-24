package lambda

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gofor-little/log"
	"github.com/gofor-little/xerror"
)

// NewProxyResponse builds an API gateway proxy response using the passed parameters. statusCode should
// be a valid HTTP status code. If contentType is ContentTypeApplicationJSON v should be a struct that is
// able to be marshaled into JSON. If contentType is ContentTypeTextHTML v should be a string or byte slice.
// If err is nil no error will be returned. If the content type is empty/unsupported or v is nil nothing will
// be written to the response body.
func NewProxyResponse(statusCode int, contentType ContentType, err error, v interface{}) (*events.APIGatewayProxyResponse, error) {
	response := &events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type": string(contentType),
			// "Access-Control-Expose-Headers":    "*",
			"Access-Control-Allow-Headers": "Date,Content-Type,Content-Length,Connection,x-amzn-RequestId,x-amz-apigw-id,X-Amzn-Trace-Id",
			"Access-Control-Allow-Origin":  "https://milhouse.dev,https://dev.millhouse.dev",
			"Access-Control-Allow-Methods": "OPTIONS,GET,PUT,POST,DELETE,PATCH,HEAD	",
		},
		StatusCode: statusCode,
	}

	if v != nil {
		switch contentType {
		case ContentTypeApplicationJSON:
			body, e := json.Marshal(v)
			if e != nil {
				log.Error(log.Fields{
					"error":      xerror.Newf("failed to marshal response body and API request failed", e, err),
					"statusCode": statusCode,
				})
				return response, xerror.Newf("failed to marshal response body and API request failed", e, err)
			}

			response.Body = string(body)
		case ContentTypeTextHTML:
			response.Body = fmt.Sprintf("%s", v)
		default:
			// Unsupported content type, do not write to the response body.
			log.Error(log.Fields{
				"error":      xerror.Newf("failed to write to response body because of unsupported content type: %s", contentType),
				"statusCode": statusCode,
			})
		}
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
