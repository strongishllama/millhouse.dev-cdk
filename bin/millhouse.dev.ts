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
  const baseDomainName =  'dev.millhouse.dev';
  const apiDomainName = `api.${baseDomainName}`;

  new BootstrapStack(app, `${namespace}-bootstrap-stack`, {
    env: env,
    tableRemovalPolicy: cdk.RemovalPolicy.DESTROY,
    enableBackups: false,
    fromAddress: `no-reply@${baseDomainName}`,
    apiDomainName: apiDomainName,
    websiteDomainName: baseDomainName
  });
  new ApiStack(app, `${namespace}-api-stack`, {
    env: env,
    accessControlAllowOrigin: `https://${baseDomainName}`,
    recaptchaSecretArn: 'arn:aws:secretsmanager:ap-southeast-2:250096756762:secret:recaptcha-secret-arn-LQz25E',
    baseDomainName: baseDomainName,
    fullDomainName: apiDomainName
  });
  new WebsiteStack(app, `${namespace}-website-stack`, {
    env: env,
    baseDomainName: baseDomainName,
    fullDomainName: baseDomainName,
    githubOAuthTokenArn: 'arn:aws:secretsmanager:ap-southeast-2:250096756762:secret:github-personal-access-token-6eAvoW',
    apiBaseUrl: `https://${apiDomainName}`,
    approvalNotifyEmails: []
  });
}

// Production stacks.
{
  const env: cdk.Environment = {
    account: '535766190525',
    region: 'ap-southeast-2'
  };
  const namespace = 'prod';
  const baseDomainName = 'millhouse.dev';
  const apiDomainName = 'api.millhouse.dev';

  new BootstrapStack(app, `${namespace}-bootstrap-stack`, {
    env: env,
    tableRemovalPolicy: cdk.RemovalPolicy.DESTROY,
    enableBackups: false,
    fromAddress: 'no-reply@millhouse.dev',
    apiDomainName: apiDomainName,
    websiteDomainName: baseDomainName
  });
  new ApiStack(app, `${namespace}-api-stack`, {
    env: env,
    accessControlAllowOrigin: 'https://millhouse.dev',
    recaptchaSecretArn: 'arn:aws:secretsmanager:ap-southeast-2:535766190525:secret:recaptcha-secret-arn-ZO4Wfp',
    baseDomainName: baseDomainName,
    fullDomainName: apiDomainName
  });
  new WebsiteStack(app, `${namespace}-website-stack`, {
    env: env,
    baseDomainName: baseDomainName,
    fullDomainName: baseDomainName,
    githubOAuthTokenArn: 'arn:aws:secretsmanager:ap-southeast-2:535766190525:secret:github-personal-access-token-aDuSLr',
    apiBaseUrl: `https://${apiDomainName}`,
    approvalNotifyEmails: []
  });
}
