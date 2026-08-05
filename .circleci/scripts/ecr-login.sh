#!/usr/bin/env bash
set -euo pipefail

# Derive ECRPowerUserRole ARN from ECR_REGISTRY (e.g. 123456789.dkr.ecr.us-east-1.amazonaws.com)
account_id=$(echo "$ECR_REGISTRY" | cut -d. -f1)
role_arn="arn:aws:iam::${account_id}:role/ECRPowerUserRole"

# Assume ECRPowerUserRole and export temp credentials for subsequent steps.
# Credentials are never printed to stdout.
creds=$(aws sts assume-role \
  --role-arn "$role_arn" \
  --role-session-name bb-portal-build)

export AWS_ACCESS_KEY_ID
export AWS_SECRET_ACCESS_KEY
export AWS_SESSION_TOKEN
AWS_ACCESS_KEY_ID=$(jq -r '.Credentials.AccessKeyId' <<< "$creds")
AWS_SECRET_ACCESS_KEY=$(jq -r '.Credentials.SecretAccessKey' <<< "$creds")
AWS_SESSION_TOKEN=$(jq -r '.Credentials.SessionToken' <<< "$creds")

printf 'export AWS_ACCESS_KEY_ID=%s\nexport AWS_SECRET_ACCESS_KEY=%s\nexport AWS_SESSION_TOKEN=%s\n' \
  "$AWS_ACCESS_KEY_ID" "$AWS_SECRET_ACCESS_KEY" "$AWS_SESSION_TOKEN" >> "$BASH_ENV"

aws ecr get-login-password --region "$AWS_DEFAULT_REGION" | \
  crane auth login "$ECR_REGISTRY" --username AWS --password-stdin
