import * as cdk from "@aws-cdk/core";
import * as certificatemanager from "@aws-cdk/aws-certificatemanager";
import * as cloudfront from "@aws-cdk/aws-cloudfront";
import * as origins from "@aws-cdk/aws-cloudfront-origins";
import * as route53 from "@aws-cdk/aws-route53";
import * as targets from "@aws-cdk/aws-route53-targets";
import * as s3 from "@aws-cdk/aws-s3";
import * as s3deploy from "@aws-cdk/aws-s3-deployment";
import * as path from "path";

export class MillhouseDevStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // Fetch hosted zone via the domain name.
    const hostedZone = route53.HostedZone.fromLookup(this, "hosted-zone", {
      domainName: "millhouse.dev",
    });

    // Create an S3 bucket to store the compiled website code.
    const bucket = new s3.Bucket(this, "bucket", {
      publicReadAccess: true,
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      autoDeleteObjects: true,
      websiteIndexDocument: "index.html",
    });

    // Create a DNS validated certificate for HTTPS. The region has to be 'us-east-1'.
    const dnsValidatedCertificate = new certificatemanager.DnsValidatedCertificate(this, "dns-validated-certificate", {
      domainName: "millhouse.dev",
      hostedZone: hostedZone,
      region: "us-east-1",
    });

    // Create a distribution attached to the S3 bucket and DNS validated certificate.
    const distribution = new cloudfront.Distribution(this, "distribution", {
      defaultBehavior: {
        origin: new origins.S3Origin(bucket)
      },
      certificate: dnsValidatedCertificate,
      defaultRootObject: "index.html",
      domainNames: [
        "millhouse.dev"
      ],
      errorResponses: [
        {
          httpStatus: 403,
          responseHttpStatus: 200,
          responsePagePath: "/index.html"
        },
        {
          httpStatus: 404,
          responseHttpStatus: 200,
          responsePagePath: "/index.html"
        }
      ]
    });

    // Create an A record pointing at the web distribution.
    new route53.ARecord(this, "a-record", {
      zone: hostedZone,
      recordName: "millhouse.dev",
      ttl: cdk.Duration.seconds(60),
      target: route53.RecordTarget.fromAlias(new targets.CloudFrontTarget(distribution)),
    });

    // Create a bucket deployment. This will use Docker to compile the website
    // and deploy it to an S3 bucket.
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
      distribution: distribution,
      distributionPaths: [
        "/*",
      ],
    });
  }
}
