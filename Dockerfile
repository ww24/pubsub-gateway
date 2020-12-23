FROM golang:1.14-alpine3.12 AS build

WORKDIR /go/src/github.com/ww24/pubsub-gateway
COPY . /go/src/github.com/ww24/pubsub-gateway
ENV CGO_ENABLED=0
RUN go build -o /usr/local/bin/pubsub-gateway ./main.go


FROM alpine:3.12

RUN apk add --no-cache tzdata ca-certificates

COPY --from=build /usr/local/bin/pubsub-gateway /usr/local/bin/pubsub-gateway
COPY --from=build /go/src/github.com/ww24/pubsub-gateway/entrypoint.sh /usr/local/bin/entrypoint.sh
COPY --from=build /go/src/github.com/ww24/pubsub-gateway/receiver-config.sample.yml /usr/local/etc/pubsub-gateway/receiver-config.yml

ENTRYPOINT [ "entrypoint.sh" ]
