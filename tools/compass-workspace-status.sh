#!/bin/sh -e

# Run upstream workspace status keys first
bash tools/workspace-status.sh

# Compass-specific keys for ECR image publishing.
# Values are supplied via CircleCI project env vars — nothing is hardcoded here.
echo "STABLE_COMPASS_ECR_REGISTRY ${ECR_REGISTRY:-}"
echo "STABLE_COMPASS_ECR_REPO ${ECR_REPO:-}"
