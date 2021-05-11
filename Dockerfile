FROM golang:1.16 AS build

WORKDIR /go/src/github.com/ww24/pubsub-gateway
COPY . .
ENV CGO_ENABLED=0
RUN go build -o /usr/local/bin/pubsub-gateway ./main.go

FROM gcr.io/distroless/static

COPY --from=build /usr/local/bin/pubsub-gateway /usr/local/bin/pubsub-gateway
ENTRYPOINT [ "pubsub-gateway" ]
