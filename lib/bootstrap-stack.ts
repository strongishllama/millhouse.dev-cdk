import * as cdk from '@aws-cdk/core';
import * as backup from '@aws-cdk/aws-backup';
import * as dynamodb from '@aws-cdk/aws-dynamodb';
import * as iam from '@aws-cdk/aws-iam';
import * as lambda from '@aws-cdk/aws-lambda';
import * as go_lambda from '@aws-cdk/aws-lambda-go';
import * as lambda_events from '@aws-cdk/aws-lambda-event-sources';
import * as ssm from '@aws-cdk/aws-ssm';
import * as sqs from '@aws-cdk/aws-sqs';
import { EmailService } from '@strongishllama/email-service-cdk';
import { SQS } from '@strongishllama/iam-constants-cdk';
import { bundling } from './lambda';
import { Stage } from './stage';

export interface BootstrapStackProps extends cdk.StackProps {
  namespace: string;
  stage: Stage;
  adminTo: string;
  adminFrom: string;
}

export class BootstrapStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props: BootstrapStackProps) {
    super(scope, id, props);

    if (props.env === undefined || props.env.account === undefined || props.env.region === undefined) {
      throw Error('Error: env property is undefined and is required to deploy this stack');
    }

    // Create the email service.
    const emailService = new EmailService(this, `${props.namespace}-email-service-${props.stage}`, {
      namespace: props.namespace,
      stage: props.stage,
      receiveMessageWaitTime: cdk.Duration.seconds(20)
    });

    // Create a string parameter for the queue ARN so other stacks can reference it.
    new ssm.StringParameter(this, `${props.namespace}-queue-arn-${props.stage}`, {
      parameterName: `${props.namespace}-email-queue-${props.stage}`,
      tier: ssm.ParameterTier.STANDARD,
      stringValue: emailService.queue.queueArn
    });

    // Create a table to store subscription emails.
    const table = new dynamodb.Table(this, `${props.namespace}-table-${props.stage}`, {
      partitionKey: {
        name: 'pk',
        type: dynamodb.AttributeType.STRING
      },
      sortKey: {
        name: 'sk',
        type: dynamodb.AttributeType.STRING
      },
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
      stream: dynamodb.StreamViewType.NEW_IMAGE,
      removalPolicy: props.stage === Stage.PROD ? cdk.RemovalPolicy.RETAIN : cdk.RemovalPolicy.DESTROY
    });

    // Create a function to receive stream events from the table and add it as an event source.
    const streamFunction = new go_lambda.GoFunction(this, `${props.namespace}-stream-functions-${props.stage}`, {
      entry: 'lambdas/stream',
      bundling: bundling,
      environment: {
        'ADMIN_TO': props.adminTo,
        'ADMIN_FROM': props.adminFrom,
        'EMAIL_QUEUE_URL': emailService.queue.queueUrl,
        'API_DOMAIN': props.stage === Stage.PROD ? 'api.millhouse.dev' : 'dev.api.millhouse.dev',
        'WEBSITE_DOMAIN': props.stage === Stage.PROD ? 'millhouse.dev' : 'dev.millhouse.dev'
      },
      initialPolicy: [
        new iam.PolicyStatement({
          actions: [
            SQS.SEND_MESSAGE
          ],
          resources: [
            emailService.queue.queueArn
          ]
        }),
      ]
    });
    streamFunction.addEventSource(new lambda_events.DynamoEventSource(table, {
      bisectBatchOnError: true,
      onFailure: new lambda_events.SqsDlq(new sqs.Queue(this, `${props.namespace}-stream-dead-letter-queue-${props.stage}`, {
        receiveMessageWaitTime: cdk.Duration.seconds(20)
      })),
      retryAttempts: 3,
      startingPosition: lambda.StartingPosition.TRIM_HORIZON
    }));

    // Create a string parameter for the table ARN so other stacks can reference it.
    new ssm.StringParameter(this, `${props.namespace}-table-arn-${props.stage}`, {
      parameterName: `${props.namespace}-table-arn-${props.stage}`,
      tier: ssm.ParameterTier.STANDARD,
      stringValue: table.tableArn
    });

    // If we're running in production, create a backup plan for the DynamoDB table.
    if (props.stage === Stage.PROD) {
      const backupPlan = backup.BackupPlan.dailyMonthly1YearRetention(this, `${props.namespace}-backup-plan-${props.stage}`);
      backupPlan.addSelection(`${props.namespace}-selection-${props.stage}`, {
        resources: [
          backup.BackupResource.fromDynamoDbTable(table)
        ]
      });
    }
  }
}