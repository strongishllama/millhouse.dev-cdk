import * as cdk from '@aws-cdk/core';
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

    // Create an S3 bucket to store the compiled website code.
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
  }
}