import * as cdk from "@aws-cdk/core";
import * as s3 from "@aws-cdk/aws-s3";
import { Stage } from "./stage";

export interface BootstrapStackProps extends cdk.StackProps {
  stage: Stage;
}

export class BootstrapStack extends cdk.Stack {
  public bucket: s3.Bucket;

  constructor(scope: cdk.Construct, id: string, props: BootstrapStackProps) {
    super(scope, id, props);

    // Create an S3 bucket to store the compiled website code.
    this.bucket = new s3.Bucket(this, "bucket", {
      publicReadAccess: true,
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      autoDeleteObjects: true,
      websiteIndexDocument: "index.html"
    });
  }
}