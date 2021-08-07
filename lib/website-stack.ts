import * as cdk from '@aws-cdk/core';
import * as codebuild from '@aws-cdk/aws-codebuild';
import * as secretsmanager from '@aws-cdk/aws-secretsmanager';
import * as s3 from '@aws-cdk/aws-s3';
import { Stage } from './stage';
import { StaticWebsiteDeployment, StaticWebsitePipeline } from '@strongishllama/static-website-cdk';

export interface WebsiteStackProps extends cdk.StackProps {
  namespace: string;
  stage: Stage;
  approvalNotifyEmails: string[];
}

export class WebsiteStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props: WebsiteStackProps) {
    super(scope, id, props);

    // Create a bucket to store the built website.
    const bucket = new s3.Bucket(this, `${props.namespace}-bucket-${props.stage}`, {
      publicReadAccess: true,
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      autoDeleteObjects: true,
      websiteIndexDocument: 'index.html'
    });

    // Create the pipeline for the website.
    new StaticWebsitePipeline(this, `${props.namespace}-pipeline-${props.stage}`, {
      namespace: props.namespace,
      stage: props.stage,
      sourceOwner: 'strongishllama',
      sourceRepo: 'millhouse.dev-frontend',
      sourceBranch: 'main',
      sourceOAuthToken: secretsmanager.Secret.fromSecretCompleteArn(
        this,
        `${props.namespace}-secret-${props.stage}`,
        'arn:aws:secretsmanager:ap-southeast-2:320045747480:secret:GithubPersonalAccessToken-ko08in'
      ).secretValue,
      buildEnvironmentVariables: {
        'VUE_APP_API_BASE_URL': {
          type: codebuild.BuildEnvironmentVariableType.PLAINTEXT,
          value: props.stage === Stage.PROD ? 'https://api.millhouse.dev' : 'https://dev.api.millhouse.dev'
        }
      },
      approvalNotifyEmails: props.approvalNotifyEmails,
      deployBucketArn: bucket.bucketArn
    });

    // Create deployment for the website.
    new StaticWebsiteDeployment(this, `${props.namespace}-deployment-${props.stage}`, {
      namespace: props.namespace,
      stage: props.stage,
      baseDomainName: 'millhouse.dev',
      fullDomainName: props.stage === Stage.PROD ? 'millhouse.dev' : 'dev.millhouse.dev',
      originBucketArn: bucket.bucketArn
    });
  }
}