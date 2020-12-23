#!/bin/bash

set -eo pipefail

# deploy before custom build
gcloud run deploy "$CLOUD_RUN_SERVICE_NAME" \
    --image="gcr.io/endpoints-release/endpoints-runtime-serverless:2" \
    --allow-unauthenticated \
    --platform=managed \
    --project="$ESP_PROJECT_ID"
