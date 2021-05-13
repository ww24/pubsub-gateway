resource "google_compute_global_address" "gateway" {
  name = "${var.name}-address"
}

resource "google_compute_region_network_endpoint_group" "gateway" {
  name                  = "${var.name}-neg"
  network_endpoint_type = "SERVERLESS"
  region                = var.location
  cloud_run {
    service = google_cloud_run_service.gateway.name
  }
}

module "gateway_lb" {
  source  = "GoogleCloudPlatform/lb-http/google//modules/serverless_negs"
  version = "~> 5.0"

  project                         = var.project
  name                            = var.name
  address                         = google_compute_global_address.gateway.self_link
  create_address                  = false
  cdn                             = false
  ssl                             = true
  use_ssl_certificates            = false
  managed_ssl_certificate_domains = var.domains
  https_redirect                  = true
  quic                            = true

  backends = {
    default = {
      description            = null
      enable_cdn             = false
      security_policy        = null
      custom_request_headers = null

      log_config = {
        enable      = false
        sample_rate = 0
      }

      groups = [
        {
          group = google_compute_region_network_endpoint_group.gateway.id
        }
      ]

      iap_config = {
        enable               = true
        oauth2_client_id     = var.oauth_client_id
        oauth2_client_secret = var.oauth_client_secret
      }
    }
  }
}
