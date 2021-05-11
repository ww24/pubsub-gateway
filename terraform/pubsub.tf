resource "google_pubsub_topic" "remocon" {
  name = "remocon"

  message_storage_policy {
    allowed_persistence_regions = [
      var.location,
    ]
  }
}

resource "google_pubsub_subscription" "remocon" {
  name  = "remocon"
  topic = google_pubsub_topic.remocon.name

  message_retention_duration = "600s"
  ack_deadline_seconds       = 10
  enable_message_ordering    = true
}
