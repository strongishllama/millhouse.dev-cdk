import * as cdk from '@aws-cdk/core';
import { ApiStack } from '../lib/api-stack';
import { BootstrapStack } from '../lib/bootstrap-stack';
import { WebsiteStack } from '../lib/website-stack';

const app = new cdk.App();


// Development stacks.
{
  const env: cdk.Environment = {
    account: '250096756762',
    region: 'ap-southeast-2'
  };
  const namespace = 'dev';
  new BootstrapStack(app, `${namespace}-bootstrap-stack`, {
    env: env,
    tableRemovalPolicy: cdk.RemovalPolicy.DESTROY,
    enableBackups: false,
    fromAddress: 'no-reply@dev.millhouse.dev',
    apiDomainName: 'api.dev.millhouse.dev',
    websiteDomainName: 'dev.millhouse.dev'
  });
  new ApiStack(app, `${namespace}-api-stack`, {
    env: env,
    accessControlAllowOrigin: 'https://dev.millhouse.dev',
    recaptchaSecretArn: 'arn:aws:secretsmanager:ap-southeast-2:250096756762:secret:recaptcha-secret-arn-LQz25E',
    baseDomainName: 'dev.millhouse.dev',
    fullDomainName: 'api.dev.millhouse.dev'
  });
  new WebsiteStack(app, `${namespace}-website-stack`, {
    env: env,
    baseDomainName: 'dev.millhouse.dev',
    fullDomainName: 'dev.millhouse.dev',
    githubOAuthTokenArn: 'arn:aws:secretsmanager:ap-southeast-2:250096756762:secret:github-personal-access-token-6eAvoW',
    apiBaseUrl: 'https://api.dev.millhouse.dev',
    approvalNotifyEmails: []
  });
}

// Production stacks.
// new BootstrapStack(app, `${namespace}-bootstrap-stack-${Stage.PROD}`, {
//   env: env,
//   namespace: namespace,
//   stage: Stage.PROD,
//   adminTo: checkEnv('ADMIN_TO'),
//   adminFrom: checkEnv('ADMIN_FROM')
// });
// new ApiStack(app, `${namespace}-api-stack-${Stage.PROD}`, {
//   env: env,
//   namespace: namespace,
//   stage: Stage.PROD,
//   lambdasConfigArn: 'arn:aws:secretsmanager:ap-southeast-2:320045747480:secret:millhouse-dev-lambda-config-8VhyY9',
//   adminTo: checkEnv('ADMIN_TO'),
//   adminFrom: checkEnv('ADMIN_FROM')
// });
// new WebsiteStack(app, `${namespace}-website-stack-${Stage.PROD}`, {
//   env: env,
//   namespace: namespace,
//   stage: Stage.PROD,
//   approvalNotifyEmails: checkEnv('APPROVAL_NOTIFY_EMAILS').split(',')
// });

// cdk.Tags.of(app).add('project', 'millhouse.dev');