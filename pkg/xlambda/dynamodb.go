package xlambda

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// UnmarshalDynamoDBImage unmarshalls image into v. For v to work with json.Unmarshal it must
// be a struct and a pointer.
//
// This helper function is required because the Lambda events package
// uses different types to the DynamoDB package. See this GitHub issue for more information
// https://github.com/aws/aws-lambda-go/issues/58
func UnmarshalDynamoDBImage(image map[string]events.DynamoDBAttributeValue, v interface{}) error {
	attributes := map[string]*dynamodb.AttributeValue{}

	for key, value := range image {
		data, err := value.MarshalJSON()
		if err != nil {
			return err
		}

		var attribute *dynamodb.AttributeValue
		if err := json.Unmarshal(data, &attribute); err != nil {
			return err
		}

		attributes[key] = attribute
	}

	return dynamodbattribute.UnmarshalMap(attributes, v)
}
