import * as cdk from "@aws-cdk/core";
import { BootstrapStack } from "../lib/bootstrap-stack";
import { checkEnv } from "../lib/env";
import { PipelineStack } from "../lib/pipeline-stack";
import { Stage } from "../lib/stage";
import { WebsiteStack } from "../lib/website-stack";

const app = new cdk.App();
cdk.Tags.of(app).add("project", "millhouse.dev");

// Development stacks.
const environmentDev = {
  account: checkEnv("AWS_ACCOUNT_DEV"),
  region: checkEnv("AWS_REGION_DEV")
};

const bootstrapStackDev = new BootstrapStack(app, `bootstrap-stack-${Stage.DEV}`, {
  env: environmentDev,
  stage: Stage.DEV
});

new PipelineStack(app, `pipeline-stack-${Stage.DEV}`, {
  env: environmentDev,
  stage: Stage.DEV,
  oauthTokenSecretArn: checkEnv("OAUTH_TOKEN_SECRET_ARN"),
  deployBucket: bootstrapStackDev.bucket,
  approvalNotifyEmails: checkEnv("APPROVAL_NOTIFY_EMAILS").split(",")
});

new WebsiteStack(app, `website-stack-${Stage.DEV}`, {
  env: environmentDev,
  stage: Stage.DEV,
  originBucket: bootstrapStackDev.bucket
});

// Production stacks.
const environmentProd = {
  account: checkEnv("AWS_ACCOUNT_PROD"),
  region: checkEnv("AWS_REGION_PROD")
};

const bootstrapStackProd = new BootstrapStack(app, `bootstrap-stack-${Stage.PROD}`, {
  env: environmentProd,
  stage: Stage.PROD
});

new PipelineStack(app, `pipeline-stack-${Stage.PROD}`, {
  env: environmentProd,
  stage: Stage.PROD,
  oauthTokenSecretArn: checkEnv("OAUTH_TOKEN_SECRET_ARN"),
  deployBucket: bootstrapStackProd.bucket,
  approvalNotifyEmails: checkEnv("APPROVAL_NOTIFY_EMAILS").split(",")
});

new WebsiteStack(app, `website-stack-${Stage.PROD}`, {
  env: environmentProd,
  stage: Stage.PROD,
  originBucket: bootstrapStackProd.bucket
});