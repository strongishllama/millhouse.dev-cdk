import * as cdk from "@aws-cdk/core";
import { checkEnv } from "../lib/env";
import { MillhouseDevStack } from "../lib/millhouse.dev-stack";

const app = new cdk.App();

cdk.Tags.of(app).add("project", "millhouse.dev");

new MillhouseDevStack(app, "millhouse-dev-stack", {
  env: {
    account: checkEnv("AWS_ACCOUNT"),
    region: checkEnv("AWS_REGION")
  }
});
