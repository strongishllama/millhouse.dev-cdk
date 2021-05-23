import * as cdk from '@aws-cdk/core';
import * as dynamodb from '@aws-cdk/aws-dynamodb';
import * as iam from '@aws-cdk/aws-iam';
import * as resources from '@aws-cdk/custom-resources';
import * as route53 from '@aws-cdk/aws-route53';
import * as s3 from '@aws-cdk/aws-s3';
import * as ssm from '@aws-cdk/aws-ssm';

import { Stage } from './stage';

export interface BootstrapStackProps extends cdk.StackProps {
  prefix: string;
  stage: Stage;
}

export class BootstrapStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props: BootstrapStackProps) {
    super(scope, id, props);

    if (props.env === undefined || props.env.account === undefined || props.env.region === undefined) {
      throw Error('Error: env property is undefined and is required to deploy this stack');
    }

    // Create an bucket to store the compiled website code.
    const bucket = new s3.Bucket(this, `${props.prefix}-bucket-${props.stage}`, {
      publicReadAccess: true,
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      autoDeleteObjects: true,
      websiteIndexDocument: 'index.html'
    });

    // Create a string parameter for the bucket ARN so other stacks can reference it.
    new ssm.StringParameter(this, `${props.prefix}-bucket-arn-${props.stage}`, {
      parameterName: `${props.prefix}-bucket-arn-${props.stage}`,
      tier: ssm.ParameterTier.STANDARD,
      stringValue: bucket.bucketArn
    });

    // Create a table to store subscription emails.
    const table = new dynamodb.Table(this, `${props.prefix}-table-${props.stage}`, {
      partitionKey: {
        name: 'pk',
        type: dynamodb.AttributeType.STRING
      },
      sortKey: {
        name: 'sk',
        type: dynamodb.AttributeType.STRING
      },
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
      removalPolicy: props.stage === Stage.PROD ? cdk.RemovalPolicy.RETAIN : cdk.RemovalPolicy.DESTROY
    });

    // Create a string parameter for the table ARN so other stacks can reference it.
    new ssm.StringParameter(this, `${props.prefix}-table-arn-${props.stage}`, {
      parameterName: `${props.prefix}-table-arn-${props.stage}`,
      tier: ssm.ParameterTier.STANDARD,
      stringValue: table.tableArn
    });

    // Fetch hosted zone via the domain name.
    const hostedZone = route53.HostedZone.fromLookup(this, `${props.prefix}-hosted-zone-${props.stage}`, {
      domainName: 'millhouse.dev'
    });

    // Create a custom resource to verify the domain.
    const verifyDomainIdentity = new resources.AwsCustomResource(this, `${props.prefix}-verify-domain-identity-${props.stage}`, {
      onCreate: {
        service: 'SES',
        action: 'verifyDomainIdentity',
        parameters: {
          Domain: 'millhouse.dev'
        },
        physicalResourceId: resources.PhysicalResourceId.fromResponse('VerificationToken')
      },
      onDelete: {
        service: 'SES',
        action: 'deleteIdentity',
        parameters: {
          Identity: 'millhouse.dev'
        }
      },
      policy: resources.AwsCustomResourcePolicy.fromStatements([
        new iam.PolicyStatement({
          actions: [
            'ses:VerifyDomainIdentity'
          ],
          resources: [
            '*'
          ]
        }),
        new iam.PolicyStatement({
          actions: [
            'ses:DeleteIdentity'
          ],
          resources: [
            `arn:aws:ses:${props.env.region}:${props.env.account}:identity/millhouse.dev`
          ]
        })
      ])
    });

    // Create a TXT record with the SES verification token.
    new route53.TxtRecord(this, `${props.prefix}-ses-verification-record-${props.stage}`, {
      zone: hostedZone,
      recordName: '_amazonses.millhouse.dev',
      values: [
        verifyDomainIdentity.getResponseField('VerificationToken')
      ]
    });
  }
}