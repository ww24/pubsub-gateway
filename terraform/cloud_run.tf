data "google_cloud_run_service" "gateway" {
  name     = var.name
  location = var.location
}

locals {
  current_image = data.google_cloud_run_service.gateway.template != null ? data.google_cloud_run_service.gateway.template[0].spec[0].containers[0].image : null
  new_image     = "${var.location}-docker.pkg.dev/${var.project}/${var.gar_repository}/${var.image_name}:${var.image_tag}"
  image         = (local.current_image != null && var.image_tag == "latest") ? local.current_image : local.new_image
}

resource "google_cloud_run_service" "gateway" {
  name     = var.name
  location = var.location
  project  = var.project

  template {
    spec {
      service_account_name = var.service_account

      timeout_seconds = 15
      containers {
        image = local.image

        resources {
          limits = {
            cpu    = "1000m"
            memory = "128Mi"
          }
        }

        env {
          name  = "AUTHORIZED_USERS"
          value = var.authorized_users
        }

        env {
          name  = "TOPIC_NAME"
          value = google_pubsub_topic.remocon.name
        }

        env {
          name  = "DEFAULT_ORIGIN"
          value = var.default_origin
        }

        env {
          name  = "ALLOW_ORIGIN_SUFFIX"
          value = var.allow_origin_suffix
        }

        env {
          name  = "CONFIG_YAML"
          value = var.config_yaml
        }

        env {
          name  = "MODE"
          value = var.mode
        }
      }
    }

    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale" = "1"

        # not working (2021/05/12)
        "run.googleapis.com/ingress" = "internal-and-cloud-load-balancing"
      }

      labels = {
        service = var.name
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }

  autogenerate_revision_name = true

  lifecycle {
    ignore_changes = [
      template[0].metadata[0].annotations["run.googleapis.com/ingress"],
    ]
  }
}
