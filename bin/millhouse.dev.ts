import * as cdk from '@aws-cdk/core';
import { BootstrapStack } from '../lib/bootstrap-stack';
import { checkEnv } from '../lib/env';
import { PipelineStack } from '../lib/pipeline-stack';
import { Stage } from '../lib/stage';
import { WebsiteStack } from '../lib/website-stack';

const environmentDev = {
  account: checkEnv('AWS_ACCOUNT_DEV'),
  region: checkEnv('AWS_REGION_DEV')
};
const environmentProd = {
  account: checkEnv('AWS_ACCOUNT_PROD'),
  region: checkEnv('AWS_REGION_PROD')
};
const prefix = 'millhouse-dev';

const app = new cdk.App();

// Development stacks.
new BootstrapStack(app, `${prefix}-bootstrap-stack-${Stage.DEV}`, {
  prefix: prefix,
  env: environmentDev,
  stage: Stage.DEV
});
new PipelineStack(app, `${prefix}-pipeline-stack-${Stage.DEV}`, {
  prefix: prefix,
  env: environmentDev,
  stage: Stage.DEV,
  oauthTokenSecretArn: checkEnv('OAUTH_TOKEN_SECRET_ARN'),
  approvalNotifyEmails: checkEnv('APPROVAL_NOTIFY_EMAILS').split(',')
});
new WebsiteStack(app, `${prefix}-website-stack-${Stage.DEV}`, {
  prefix: prefix,
  env: environmentDev,
  stage: Stage.DEV,
});

// Production stacks.
new BootstrapStack(app, `${prefix}-bootstrap-stack-${Stage.PROD}`, {
  prefix: prefix,
  env: environmentProd,
  stage: Stage.PROD
});
new PipelineStack(app, `${prefix}-pipeline-stack-${Stage.PROD}`, {
  prefix: prefix,
  env: environmentProd,
  stage: Stage.PROD,
  oauthTokenSecretArn: checkEnv('OAUTH_TOKEN_SECRET_ARN'),
  approvalNotifyEmails: checkEnv('APPROVAL_NOTIFY_EMAILS').split(',')
});
new WebsiteStack(app, `${prefix}-website-stack-${Stage.PROD}`, {
  prefix: prefix,
  env: environmentProd,
  stage: Stage.PROD,
});

cdk.Tags.of(app).add('project', 'millhouse.dev');