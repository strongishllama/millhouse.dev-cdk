import * as cdk from '@aws-cdk/core';
import { ApiStack } from '../lib/api-stack';
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
const namespace = 'millhouse-dev';

const app = new cdk.App();

// Development stacks.
new BootstrapStack(app, `${namespace}-bootstrap-stack-${Stage.DEV}`, {
  env: environmentDev,
  namespace: namespace,
  stage: Stage.DEV,
  adminTo: checkEnv('ADMIN_TO_DEV'),
  adminFrom: checkEnv('ADMIN_FROM_DEV')
});
new ApiStack(app, `${namespace}-api-stack-${Stage.DEV}`, {
  env: environmentDev,
  namespace: namespace,
  stage: Stage.DEV,
  lambdasConfigArn: checkEnv('SUBSCRIBE_CONFIG_ARN_DEV'),
  adminTo: checkEnv('ADMIN_TO_DEV'),
  adminFrom: checkEnv('ADMIN_FROM_DEV')
});
new PipelineStack(app, `${namespace}-pipeline-stack-${Stage.DEV}`, {
  env: environmentDev,
  namespace: namespace,
  stage: Stage.DEV,
  oauthTokenSecretArn: checkEnv('OAUTH_TOKEN_SECRET_ARN'),
  approvalNotifyEmails: checkEnv('APPROVAL_NOTIFY_EMAILS').split(',')
});
new WebsiteStack(app, `${namespace}-website-stack-${Stage.DEV}`, {
  env: environmentDev,
  namespace: namespace,
  stage: Stage.DEV,
});

// Production stacks.
new BootstrapStack(app, `${namespace}-bootstrap-stack-${Stage.PROD}`, {
  env: environmentProd,
  namespace: namespace,
  stage: Stage.PROD,
  adminTo: checkEnv('ADMIN_TO_PROD'),
  adminFrom: checkEnv('ADMIN_FROM_PROD')
});
new ApiStack(app, `${namespace}-api-stack-${Stage.PROD}`, {
  env: environmentProd,
  namespace: namespace,
  stage: Stage.PROD,
  lambdasConfigArn: checkEnv('SUBSCRIBE_CONFIG_ARN_PROD'),
  adminTo: checkEnv('ADMIN_TO_PROD'),
  adminFrom: checkEnv('ADMIN_FROM_PROD')
});
new PipelineStack(app, `${namespace}-pipeline-stack-${Stage.PROD}`, {
  env: environmentProd,
  namespace: namespace,
  stage: Stage.PROD,
  oauthTokenSecretArn: checkEnv('OAUTH_TOKEN_SECRET_ARN'),
  approvalNotifyEmails: checkEnv('APPROVAL_NOTIFY_EMAILS').split(',')
});
new WebsiteStack(app, `${namespace}-website-stack-${Stage.PROD}`, {
  env: environmentProd,
  namespace: namespace,
  stage: Stage.PROD,
});

cdk.Tags.of(app).add('project', 'millhouse.dev');