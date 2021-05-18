import * as cdk from '@aws-cdk/core';
import * as apigateway from '@aws-cdk/aws-apigateway';
import * as certificatemanager from '@aws-cdk/aws-certificatemanager';
import * as iam from '@aws-cdk/aws-iam';
import * as lambda from '@aws-cdk/aws-lambda-go';
import * as route53 from '@aws-cdk/aws-route53';
import * as targets from '@aws-cdk/aws-route53-targets';
import { EmailService } from '@strongishllama/email-service-cdk';
import { Stage } from './stage';
import { Method } from './method';
import { Certificate } from 'crypto';

export interface ApiStackProps extends cdk.StackProps {
  prefix: string;
  stage: Stage;
  subscribeConfigArn: string;
}

export class ApiStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props: ApiStackProps) {
    super(scope, id, props);

    // Create a REST API for the website to interact with.
    const api = new apigateway.RestApi(this, `${props.prefix}-rest-api-${props.stage}`, {
      deployOptions: {
        stageName: props.stage
      }
    });

    // Add ping method - /
    const subscribeFunction = new lambda.GoFunction(this, `${props.prefix}-ping-function-${props.stage}`, {
      entry: 'lambdas/api/ping'
    });
    api.root.addMethod(Method.GET, new apigateway.LambdaIntegration(subscribeFunction));

    // Add subscribe method - /subscribe
    api.root.addResource('subscribe').addMethod(Method.PUT, new apigateway.LambdaIntegration(new lambda.GoFunction(this, `${props.prefix}-subscribe-function-${props.stage}`, {
      entry: 'lambdas/api/subscribe',
      environment: {
        "CONFIG_SECRET_ARN": props.subscribeConfigArn
      },
      initialPolicy: [
        new iam.PolicyStatement({
          actions: [
            "secretsmanager:GetSecretValue"
          ],
          resources: [
            props.subscribeConfigArn
          ]
        })
      ]
    })));

    // Fetch hosted zone via the domain name.
    const hostedZone = route53.HostedZone.fromLookup(this, `${props.prefix}-hosted-zone-${props.stage}`, {
      domainName: 'millhouse.dev'
    });

    // Determine the full domain name based on the stage.
    const fullDomainName = props.stage === Stage.PROD ? 'api.millhouse.dev' : `${props.stage}.api.millhouse.dev`;

    // Create a DNS validated certificate for HTTPS
    const certificate = new certificatemanager.DnsValidatedCertificate(this, `${props.prefix}-api-certificate-${props.stage}`, {
      domainName: fullDomainName,
      hostedZone: route53.HostedZone.fromLookup(this, `${props.prefix}-api-hosted-zone-${props.stage}`, {
        domainName: 'millhouse.dev'
      })
    });

    // Create a domain name for the API and map it.
    const domain = new apigateway.DomainName(this, `${props.prefix}-api-domain-name-${props.stage}`, {
      domainName: fullDomainName,
      certificate: certificate,
    });
    domain.addBasePathMapping(api);

    // Create an A record pointing at the web distribution.
    new route53.ARecord(this, `${props.prefix}-a-record-${props.stage}`, {
      zone: hostedZone,
      recordName: fullDomainName,
      ttl: cdk.Duration.seconds(60),
      target: route53.RecordTarget.fromAlias(new targets.ApiGatewayDomain(domain))
    });

    // Create the email service.
    new EmailService(this, `${props.prefix}-email-service-${props.stage}`, {
      prefix: props.prefix,
      suffix: props.stage,
      sendMessageArns: [
        subscribeFunction.functionArn
      ]
    });
  }

}