pubsub-gateway
===

[![Go Reference][go-dev-img]][go-dev-url]
![Test][github-actions-img]

Gateway to publish events into Cloud Pub/Sub topic.

[![dockeri.co][dockeri-img]][dockeri-url]

## Use case

- Web client want to publish events.
  - Cloud Pub/Sub do not support CORS.
- Do not want to permit Cloud Pub/Sub resources to end-user.

[github-actions-img]: https://github.com/ww24/pubsub-gateway/workflows/Test/badge.svg?branch=master
[dockeri-img]: https://dockeri.co/image/ww24/pubsub-gateway
[dockeri-url]: https://hub.docker.com/r/ww24/pubsub-gateway
[go-dev-img]: https://pkg.go.dev/badge/github.com/ww24/pubsub-gateway.svg
[go-dev-url]: https://pkg.go.dev/github.com/ww24/pubsub-gateway
