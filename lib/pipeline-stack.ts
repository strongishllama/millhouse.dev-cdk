import * as cdk from "@aws-cdk/core";
import * as codebuild from "@aws-cdk/aws-codebuild";
import * as codepipeline from "@aws-cdk/aws-codepipeline";
import * as codepipeline_actions from "@aws-cdk/aws-codepipeline-actions";
import * as s3 from "@aws-cdk/aws-s3";
import * as secretsmanager from "@aws-cdk/aws-secretsmanager";
import { Stage } from "./stage";

export interface PipelineStackProps extends cdk.StackProps {
  prefix: string;
  stage: Stage;
  oauthTokenSecretArn: string;
  deployBucketArn: string;
  approvalNotifyEmails: string[];
}

export class PipelineStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props: PipelineStackProps) {
    super(scope, id, props);

    // Create artifacts to pass data between stages.
    const sourceOutput = new codepipeline.Artifact();
    const buildOutput = new codepipeline.Artifact();

    // Create pipeline to source, build and deploy the frontend website.
    const pipeline = new codepipeline.Pipeline(this, `${props.prefix}-pipeline-${props.stage}`, {
      artifactBucket: new s3.Bucket(this, `${props.prefix}-bucket-${props.stage}`, {
        removalPolicy: cdk.RemovalPolicy.DESTROY,
        autoDeleteObjects: true
      }),
      stages: [
        {
          stageName: "Source",
          actions: [
            new codepipeline_actions.GitHubSourceAction({
              actionName: "Source",
              output: sourceOutput,
              owner: "strongishllama",
              repo: "millhouse.dev-frontend",
              branch: props.stage === Stage.PROD ? "release" : "main",
              oauthToken: secretsmanager.Secret.fromSecretCompleteArn(this, `${props.prefix}-secret-${props.stage}`, props.oauthTokenSecretArn).secretValue,
              trigger: codepipeline_actions.GitHubTrigger.POLL
            })
          ]
        },
        {
          stageName: "Build",
          actions: [
            new codepipeline_actions.CodeBuildAction({
              actionName: "Build",
              input: sourceOutput,
              outputs: [
                buildOutput
              ],
              project: new codebuild.PipelineProject(this, `${props.prefix}-project-${props.stage}`, {
                environment: {
                  buildImage: codebuild.LinuxBuildImage.STANDARD_5_0
                }
              })
            })
          ]
        }
      ]
    });

    // If we're running in production, add a manual approval step before deployment.
    if (props.stage === Stage.PROD) {
      pipeline.addStage({
        stageName: "Approval",
        actions: [
          new codepipeline_actions.ManualApprovalAction({
            actionName: "Approval",
            notifyEmails: props.approvalNotifyEmails
          })
        ]
      });
    }

    pipeline.addStage({
      stageName: "Deploy",
      actions: [
        new codepipeline_actions.S3DeployAction({
          actionName: "Deploy",
          input: buildOutput,
          bucket: s3.Bucket.fromBucketArn(this, `${props.prefix}-deploy-bucket-${props.stage}`, props.deployBucketArn)
        })
      ]
    });
  }
}