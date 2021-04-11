import * as cdk from "@aws-cdk/core";
import * as certificatemanager from "@aws-cdk/aws-certificatemanager";
import * as cloudfront from "@aws-cdk/aws-cloudfront";
import * as route53 from "@aws-cdk/aws-route53";
import * as s3 from "@aws-cdk/aws-s3";
import * as s3deploy from "@aws-cdk/aws-s3-deployment";
import * as targets from "@aws-cdk/aws-route53-targets";
import * as path from "path";

export class MillhouseDevStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const hostedZone = route53.HostedZone.fromLookup(this, "hosted-zone", {
      domainName: "tothepoint.dev",
    });

    const bucket = new s3.Bucket(this, "bucket", {
      publicReadAccess: true,
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      autoDeleteObjects: true,
      websiteIndexDocument: "index.html",
    });

    const dnsValidatedCertificate = new certificatemanager.DnsValidatedCertificate(this, "dns-validated-certificate", {
      domainName: "example.tothepoint.dev",
      hostedZone: hostedZone,
      region: "us-east-1",
    });

    const webDistribution = new cloudfront.CloudFrontWebDistribution(this, "web-distribution", {
      // aliasConfiguration: {
      //   acmCertRef: dnsValidatedCertificate.certificateArn,
      //   names: [
      //     "example.tothepoint.dev",
      //   ],
      //   sslMethod: cloudfront.SSLMethod.SNI,
      //   securityPolicy: cloudfront.SecurityPolicyProtocol.TLS_V1_2_2019,
      // },
      originConfigs: [
        {
          customOriginSource: {
            domainName: bucket.bucketWebsiteDomainName,
            originProtocolPolicy: cloudfront.OriginProtocolPolicy.HTTP_ONLY,
          },
          behaviors: [
            {
              isDefaultBehavior: true,
            },
          ],
        },
      ],
      viewerCertificate: cloudfront.ViewerCertificate.fromAcmCertificate(dnsValidatedCertificate, {
        // sslMethod: cloudfront.SSLMethod.SNI,
        // securityPolicy: cloudfront.SecurityPolicyProtocol.TLS_V1_2_2019
        aliases: [
          "example.tothepoint.dev"
        ]
      }),
    });

    new route53.ARecord(this, "a-record", {
      recordName: "example.tothepoint.dev",
      zone: hostedZone,
      ttl: cdk.Duration.seconds(60),
      target: route53.RecordTarget.fromAlias(new targets.CloudFrontTarget(webDistribution)),
    });

    new s3deploy.BucketDeployment(this, 'bucket-deployment', {
      sources: [
        s3deploy.Source.asset(path.join(__dirname, "../frontend"), {
          bundling: {
            image: cdk.DockerImage.fromRegistry("node"),
            command: [
              "bash", "-c", [
                "npm ci",
                "npm run build",
                "mv dist/* /asset-output",
              ].join(" && "),
            ],
            user: "root",
          },
        }),
      ],
      destinationBucket: bucket,
      distribution: webDistribution,
      distributionPaths: [
        "/*",
      ],
    });
  }
}
