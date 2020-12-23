#!/bin/bash

set -eo pipefail

gcloud run services add-iam-policy-binding "$BACKEND_SERVICE_NAME" \
  --member "serviceAccount:$ESP_PROJECT_NUMBER-compute@developer.gserviceaccount.com" \
  --role "roles/run.invoker" \
  --platform managed \
  --project "$BACKEND_PROJECT_ID"
