import * as cdk from '@aws-cdk/core';
import * as codebuild from '@aws-cdk/aws-codebuild';
import * as secretsmanager from '@aws-cdk/aws-secretsmanager';
import * as s3 from '@aws-cdk/aws-s3';
import { StaticWebsiteDeployment, StaticWebsitePipeline } from '@strongishllama/static-website-cdk';

export interface WebsiteStackProps extends cdk.StackProps {
  readonly baseDomainName: string;
  readonly fullDomainName: string;
  readonly githubOAuthTokenArn: string;
  readonly apiBaseUrl: string;
  readonly approvalNotifyEmails: string[];
}

export class WebsiteStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props: WebsiteStackProps) {
    super(scope, id, props);

    if (props.env === undefined || props.env.account === undefined) {
      throw new Error('WebsiteStackProps.env property must be fully defined');
    }

    const bucket = new s3.Bucket(this, 'bucket', {
      publicReadAccess: true,
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      autoDeleteObjects: true,
      websiteIndexDocument: 'index.html'
    });

    const deployment = new StaticWebsiteDeployment(this, 'deployment', {
      baseDomainName: props.baseDomainName,
      fullDomainName: props.fullDomainName,
      originBucketArn: bucket.bucketArn
    });

    new StaticWebsitePipeline(this, 'pipeline', {
      sourceOwner: 'strongishllama',
      sourceRepo: 'millhouse.dev-frontend',
      sourceBranch: 'main',
      githubOAuthToken: secretsmanager.Secret.fromSecretCompleteArn(this, 'secret', props.githubOAuthTokenArn).secretValue,
      buildEnvironmentVariables: {
        'VUE_APP_API_BASE_URL': {
          type: codebuild.BuildEnvironmentVariableType.PLAINTEXT,
          value: props.apiBaseUrl
        }
      },
      approvalNotifyEmails: props.approvalNotifyEmails,
      deployBucketArn: bucket.bucketArn,
      distributionId: deployment.distribution.distributionId,
      account: props.env.account
    });
  }
}