import * as cdk from '@aws-cdk/core';
import * as events from '@aws-cdk/aws-events';
import * as events_targets from '@aws-cdk/aws-events-targets';
import * as go_lambda from '@aws-cdk/aws-lambda-go';
import { bundling } from './lambda';
import { Stage } from './stage';

export interface AnalyticsStackProps extends cdk.StackProps {
  namespace: string;
  stage: Stage;
}

export class AnalyticsStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props: AnalyticsStackProps) {
    super(scope, id, props);

    const rule = new events.Rule(this, `${props.namespace}-rule-${props.stage}`, {
      schedule: events.Schedule.cron({
        minute: '30',
        hour: '6'
      }),
      targets: [
        new events_targets.LambdaFunction(new go_lambda.GoFunction(this, `${props.namespace}-new-subscribers-function-${props.stage}`, {
          entry: 'lambdas/analytics/new-subscribers',
          bundling: bundling
        }))
      ]
    });
  }
}