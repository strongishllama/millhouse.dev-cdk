import * as cdk from '@aws-cdk/core';
import * as apigateway from '@aws-cdk/aws-apigateway';
import * as certificatemanager from '@aws-cdk/aws-certificatemanager';
import * as dynamodb from '@aws-cdk/aws-dynamodb';
import * as iam from '@aws-cdk/aws-iam';
import * as go_lambda from '@aws-cdk/aws-lambda-go';
import * as route53 from '@aws-cdk/aws-route53';
import * as targets from '@aws-cdk/aws-route53-targets';
import * as ssm from '@aws-cdk/aws-ssm';
import * as sqs from '@aws-cdk/aws-sqs';
import { DynamoDB, SecretsManager, SQS } from '@strongishllama/aws-iam-constants';
import { bundling } from './lambda';
import { Method } from './method';

export interface ApiStackProps extends cdk.StackProps {
  readonly accessControlAllowOrigin: string;
  readonly recaptchaSecretArn: string;
  readonly baseDomainName: string;
  readonly fullDomainName: string;
}

export class ApiStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props: ApiStackProps) {
    super(scope, id, props);

    const table = dynamodb.Table.fromTableArn(this, 'subscription-table', ssm.StringParameter.fromStringParameterName(this, 'table-arn', 'table-arn').stringValue);
    const emailQueue = sqs.Queue.fromQueueArn(this, 'email-queue', ssm.StringParameter.fromStringParameterName(this, 'email-queue-arn', 'email-queue-arn').stringValue);

    const api = new apigateway.RestApi(this, 'rest-api', {
      defaultCorsPreflightOptions: {
        allowOrigins: [props.accessControlAllowOrigin]
      }
    });

    // Add ping method - /
    api.root.addMethod(Method.GET, new apigateway.LambdaIntegration(new go_lambda.GoFunction(this, 'ping-function', {
      entry: 'lambdas/api/ping',
      bundling: bundling
    })));

    // Add subscribe method - /subscribe
    api.root.addResource('subscribe').addMethod(Method.PUT, new apigateway.LambdaIntegration(new go_lambda.GoFunction(this, 'subscribe-function', {
      entry: 'lambdas/api/subscribe',
      bundling: bundling,
      environment: {
        'ACCESS_CONTROL_ALLOW_ORIGIN': props.accessControlAllowOrigin,
        'RECAPTCHA_SECRET_ARN': props.recaptchaSecretArn,
        'EMAIL_QUEUE_URL': emailQueue.queueUrl,
        'TABLE_NAME': table.tableName
      },
      initialPolicy: [
        new iam.PolicyStatement({
          actions: [
            SecretsManager.GET_SECRET_VALUE
          ],
          resources: [
            props.recaptchaSecretArn
          ]
        }),
        new iam.PolicyStatement({
          actions: [
            DynamoDB.PUT_ITEM,
            DynamoDB.QUERY,
          ],
          resources: [
            table.tableArn,
            `${table.tableArn}/index/*`
          ]
        }),
        new iam.PolicyStatement({
          actions: [
            SQS.SEND_MESSAGE
          ],
          resources: [
            emailQueue.queueArn
          ]
        }),
      ]
    })));

    // Add unsubscribe method - /unsubscribe
    api.root.addResource('unsubscribe').addMethod(Method.GET, new apigateway.LambdaIntegration(new go_lambda.GoFunction(this, 'unsubscribe-function', {
      entry: 'lambdas/api/unsubscribe',
      bundling: bundling,
      environment: {
        'ACCESS_CONTROL_ALLOW_ORIGIN': props.accessControlAllowOrigin,
        'TABLE_NAME': table.tableName
      },
      initialPolicy: [
        new iam.PolicyStatement({
          actions: [
            DynamoDB.DELETE_ITEM
          ],
          resources: [
            table.tableArn
          ]
        })
      ]
    })));

    const hostedZone = route53.HostedZone.fromLookup(this, 'hosted-zone', {
      domainName: props.baseDomainName
    });

    const certificate = new certificatemanager.DnsValidatedCertificate(this, 'api-certificate', {
      domainName: props.fullDomainName,
      hostedZone: hostedZone
    });

    const domain = new apigateway.DomainName(this, 'api-domain-name', {
      domainName: props.fullDomainName,
      certificate: certificate,
    });
    domain.addBasePathMapping(api);

    new route53.ARecord(this, 'a-record', {
      zone: hostedZone,
      recordName: props.fullDomainName,
      ttl: cdk.Duration.seconds(60),
      target: route53.RecordTarget.fromAlias(new targets.ApiGatewayDomain(domain))
    });
  }

}