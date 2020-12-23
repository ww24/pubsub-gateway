#!/bin/bash

set -eo pipefail

# build custom ESPv2 image
chmod +x gcloud_build_image
./gcloud_build_image -s "$CLOUD_RUN_HOSTNAME" \
    -c "$ENDPOINTS_CONFIG_ID" -p "$ESP_PROJECT_ID"

# deploy after custom build
gcloud run deploy "$CLOUD_RUN_SERVICE_NAME" \
    --image="gcr.io/${ESP_PROJECT_ID}/endpoints-runtime-serverless:${CLOUD_RUN_HOSTNAME}-${ENDPOINTS_CONFIG_ID}" \
    --allow-unauthenticated \
    --platform=managed \
    --project="$ESP_PROJECT_ID"

# Client ID を作成、削除する度に再 deploy が必要
