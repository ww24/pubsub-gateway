#!/bin/ash

set -eo pipefail

if [ -n "$SERVICE_ACCOUNT" ]; then
    GOOGLE_APPLICATION_CREDENTIALS=/usr/local/etc/pubsub-gateway/credential.json
    echo -n "$SERVICE_ACCOUNT" | base64 -d > "$GOOGLE_APPLICATION_CREDENTIALS"
fi

if [ -n "$CONFIG" ]; then
    echo -n "$CONFIG" | base64 -d > /usr/local/etc/pubsub-gateway/receiver-config.yml
fi

export GOOGLE_APPLICATION_CREDENTIALS
pubsub-gateway -config /usr/local/etc/pubsub-gateway/receiver-config.yml
