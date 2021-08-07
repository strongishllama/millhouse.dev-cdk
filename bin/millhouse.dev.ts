import * as cdk from '@aws-cdk/core';
import { ApiStack } from '../lib/api-stack';
import { BootstrapStack } from '../lib/bootstrap-stack';
import { checkEnv } from '../lib/env';
import { Stage } from '../lib/stage';
import { WebsiteStack } from '../lib/website-stack';

const env = {
  account: '320045747480',
  region: 'ap-southeast-2'
};
const namespace = 'millhouse-dev';

const app = new cdk.App();

// Development stacks.
new BootstrapStack(app, `${namespace}-bootstrap-stack-${Stage.DEV}`, {
  env: env,
  namespace: namespace,
  stage: Stage.DEV,
  adminTo: checkEnv('ADMIN_TO'),
  adminFrom: checkEnv('ADMIN_FROM')
});
new ApiStack(app, `${namespace}-api-stack-${Stage.DEV}`, {
  env: env,
  namespace: namespace,
  stage: Stage.DEV,
  lambdasConfigArn: 'arn:aws:secretsmanager:ap-southeast-2:320045747480:secret:millhouse-dev-lambda-config-8VhyY9',
  adminTo: checkEnv('ADMIN_TO'),
  adminFrom: checkEnv('ADMIN_FROM')
});
new WebsiteStack(app, `${namespace}-website-stack-${Stage.DEV}`, {
  env: env,
  namespace: namespace,
  stage: Stage.DEV,
  approvalNotifyEmails: []
});

// Production stacks.
new BootstrapStack(app, `${namespace}-bootstrap-stack-${Stage.PROD}`, {
  env: env,
  namespace: namespace,
  stage: Stage.PROD,
  adminTo: checkEnv('ADMIN_TO'),
  adminFrom: checkEnv('ADMIN_FROM')
});
new ApiStack(app, `${namespace}-api-stack-${Stage.PROD}`, {
  env: env,
  namespace: namespace,
  stage: Stage.PROD,
  lambdasConfigArn: 'arn:aws:secretsmanager:ap-southeast-2:320045747480:secret:millhouse-dev-lambda-config-8VhyY9',
  adminTo: checkEnv('ADMIN_TO'),
  adminFrom: checkEnv('ADMIN_FROM')
});
new WebsiteStack(app, `${namespace}-website-stack-${Stage.PROD}`, {
  env: env,
  namespace: namespace,
  stage: Stage.PROD,
  approvalNotifyEmails: checkEnv('APPROVAL_NOTIFY_EMAILS').split(',')
});

cdk.Tags.of(app).add('project', 'millhouse.dev');