package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/strongishllama/millhouse.dev-cdk/lambdas/api/ping/handler"
)

func main() {
	lambda.Start(handler.Handle)
}
