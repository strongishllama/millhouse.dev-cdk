import * as cdk from '@aws-cdk/core';
import * as certificatemanager from '@aws-cdk/aws-certificatemanager';
import * as cloudfront from '@aws-cdk/aws-cloudfront';
import * as origins from '@aws-cdk/aws-cloudfront-origins';
import * as route53 from '@aws-cdk/aws-route53';
import * as targets from '@aws-cdk/aws-route53-targets';
import * as s3 from '@aws-cdk/aws-s3';
import * as ssm from '@aws-cdk/aws-ssm';
import { Stage } from './stage';

export interface WebsiteStackProps extends cdk.StackProps {
  prefix: string;
  stage: Stage;
}

export class WebsiteStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props: WebsiteStackProps) {
    super(scope, id, props);

    // Fetch hosted zone via the domain name.
    const hostedZone = route53.HostedZone.fromLookup(this, `${props.prefix}-hosted-zone-${props.stage}`, {
      domainName: 'millhouse.dev'
    });

    // Determine the full domain name based on the stage.
    const fullDomainName = props.stage === Stage.PROD ? 'millhouse.dev' : `${props.stage}.millhouse.dev`;

    // Create a DNS validated certificate for HTTPS. The region has to be 'us-east-1'.
    const dnsValidatedCertificate = new certificatemanager.DnsValidatedCertificate(this, `${props.prefix}-dns-validated-certificate-${props.stage}`, {
      domainName: fullDomainName,
      hostedZone: hostedZone,
      region: 'us-east-1',
    });

    // Create a distribution attached to the S3 bucket and DNS validated certificate.
    const distribution = new cloudfront.Distribution(this, `${props.prefix}-distribution-${props.stage}`, {
      defaultBehavior: {
        viewerProtocolPolicy: cloudfront.ViewerProtocolPolicy.REDIRECT_TO_HTTPS,
        origin: new origins.S3Origin(s3.Bucket.fromBucketArn(
          this,
          `${props.prefix}-origin-bucket-${props.stage}`,
          ssm.StringParameter.fromStringParameterName(this, `${props.prefix}-origin-bucket-arn-${props.stage}`, `${props.prefix}-bucket-arn-${props.stage}`).stringValue
        )),
      },
      certificate: dnsValidatedCertificate,
      defaultRootObject: 'index.html',
      domainNames: [
        fullDomainName
      ],
      errorResponses: [
        {
          httpStatus: 403,
          responseHttpStatus: 200,
          responsePagePath: '/index.html'
        },
        {
          httpStatus: 404,
          responseHttpStatus: 200,
          responsePagePath: '/index.html'
        }
      ]
    });

    // Create an A record pointing at the web distribution.
    new route53.ARecord(this, `${props.prefix}-a-record-${props.stage}`, {
      zone: hostedZone,
      recordName: fullDomainName,
      ttl: cdk.Duration.seconds(60),
      target: route53.RecordTarget.fromAlias(new targets.CloudFrontTarget(distribution))
    });
  }
}