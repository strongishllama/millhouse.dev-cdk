import * as cdk from '@aws-cdk/core';
import * as apigateway from '@aws-cdk/aws-apigateway';
import * as lambda from '@aws-cdk/aws-lambda-go';
import { Stage } from './stage';
import { Method } from './method';

export interface ApiStackProps extends cdk.StackProps {
  prefix: string;
  stage: Stage;
}

export class ApiStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props: ApiStackProps) {
    super(scope, id, props);

    // Create a REST API for the website to interact with.
    const api = new apigateway.RestApi(this, `${props.prefix}-rest-api-${props.stage}`, {
      deployOptions: {
        stageName: props.stage
      }
    });

    // Add ping method - /
    api.root.addMethod(Method.GET, new apigateway.LambdaIntegration(new lambda.GoFunction(this, `${props.prefix}-ping-function-${props.stage}`, {
      entry: 'lambdas/api/ping'
    })));

    // Add subscribe method - /subscribe
    api.root.addResource('subscribe').addMethod(Method.PUT, new apigateway.LambdaIntegration(new lambda.GoFunction(this, `${props.prefix}-subscribe-function-${props.stage}`, {
      entry: 'lambdas/api/subscribe'
    })));
  }

}