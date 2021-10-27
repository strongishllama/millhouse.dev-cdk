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
import { SQS } from '@strongishllama/aws-iam-constants';
import { bundling } from './lambda';

export interface BootstrapStackProps extends cdk.StackProps {
  readonly tableRemovalPolicy: cdk.RemovalPolicy;
  readonly enableBackups: boolean;
  readonly fromAddress: string;
  readonly apiDomainName: string;
  readonly websiteDomainName: string;
}

export class BootstrapStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props: BootstrapStackProps) {
    super(scope, id, props);

    if (props.env === undefined || props.env.account === undefined || props.env.region === undefined) {
      throw Error('BootstrapStackProps.env property must be fully defined.');
    }

    const emailService = new EmailService(this, 'email-service', {
      receiveMessageWaitTime: cdk.Duration.seconds(20)
    });

    const table = new dynamodb.Table(this, 'table', {
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
      removalPolicy: props.tableRemovalPolicy
    });

    const streamFunction = new go_lambda.GoFunction(this, 'stream-function', {
      entry: 'lambdas/stream',
      bundling: bundling,
      environment: {
        'FROM_ADDRESS': props.fromAddress,
        'EMAIL_QUEUE_URL': emailService.queue.queueUrl,
        'API_DOMAIN': props.apiDomainName,
        'WEBSITE_DOMAIN': props.websiteDomainName
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
      onFailure: new lambda_events.SqsDlq(new sqs.Queue(this, 'stream-dead-letter-queue', {
        receiveMessageWaitTime: cdk.Duration.seconds(20)
      })),
      retryAttempts: 3,
      startingPosition: lambda.StartingPosition.TRIM_HORIZON
    }));

    if (props.enableBackups) {
      const backupPlan = backup.BackupPlan.dailyMonthly1YearRetention(this, 'backup-plan');
      backupPlan.addSelection('selection', {
        resources: [
          backup.BackupResource.fromDynamoDbTable(table)
        ]
      });
    }

    new ssm.StringParameter(this, 'table-arn', {
      parameterName: 'table-arn',
      tier: ssm.ParameterTier.STANDARD,
      stringValue: table.tableArn
    });
    new ssm.StringParameter(this, 'queue-arn', {
      parameterName: 'email-queue-arn',
      tier: ssm.ParameterTier.STANDARD,
      stringValue: emailService.queue.queueArn
    });
  }
}
