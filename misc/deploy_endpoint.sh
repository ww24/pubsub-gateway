#!/bin/bash

set -eo pipefail

envsubst < openapi.yaml > .tmp_openapi.yaml
gcloud endpoints services deploy .tmp_openapi.yaml
