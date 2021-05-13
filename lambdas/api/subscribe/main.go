package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gofor-little/log"

	"github.com/strongishllama/millhouse.dev-cdk/lambdas/api/subscribe/handler"
)

func main() {
	log.Log = log.NewStandardLogger(os.Stdout, nil)

	lambda.Start(handler.Handle)
}
