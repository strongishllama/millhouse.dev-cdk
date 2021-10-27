import * as lambda from '@aws-cdk/aws-lambda-go';

export const bundling: lambda.BundlingOptions = {
  goBuildFlags: [
    '-ldflags="-s -w"'
  ]
};
