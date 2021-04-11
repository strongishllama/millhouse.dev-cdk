import * as cdk from "@aws-cdk/core";
import { MillhouseDevStack } from "../lib/millhouse.dev-stack";

const app = new cdk.App();
new MillhouseDevStack(app, "millhouse-dev-stack", {
  env: {
    account: "320045747480",
    region: "ap-southeast-2"
  }
});
