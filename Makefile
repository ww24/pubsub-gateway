REGION ?= asia-northeast1
PROJECT_ID ?=
CLOUD_RUN_SERVICE_ACCOUNT_EMAIL ?=
TOPIC_NAME ?=
DEFAULT_ORIGIN ?=
ALLOW_ORIGIN_SUFFIX ?=
IMAGE := asia.gcr.io/${PROJECT_ID}/pubsub-gateway:latest

BIN := $(abspath ./bin)
GO ?= go
GO_ENV ?= GOBIN=$(BIN)

.PHONY: build
build:
	docker build -t ${IMAGE} .

define FLAGS
--set-env-vars:
  TOPIC_NAME: "${TOPIC_NAME}"
  DEFAULT_ORIGIN: "${DEFAULT_ORIGIN}"
  ALLOW_ORIGIN_SUFFIX: "${ALLOW_ORIGIN_SUFFIX}"
  AUTHORIZED_USERS: "${AUTHORIZED_USERS}"
endef
export FLAGS

.PHONY: flags
flags:
	echo "$${FLAGS}" > .flags.yml

.PHONY: deploy
deploy: flags
	docker push ${IMAGE}
	gcloud run deploy pubsub-gateway \
	--image="${IMAGE}" \
	--region="${REGION}" \
	--platform=managed \
	--max-instances=1 \
	--memory=128Mi \
	--service-account=${CLOUD_RUN_SERVICE_ACCOUNT_EMAIL} \
	--no-allow-unauthenticated \
	--flags-file .flags.yml

.PHONY: generate
generate:
	PATH=$(BIN):${PATH} $(GO_ENV) $(GO) generate ./...

.PHONY: lint
lint:
	golangci-lint run
