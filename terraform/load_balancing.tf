resource "google_compute_global_address" "gateway" {
  # disable
  count = 0

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

# Create Cert manually
# https://console.cloud.google.com/net-services/loadbalancing/advanced/sslCertificates/list
data "google_compute_ssl_certificate" "default" {
  name = var.cert
}

module "gateway_lb" {
  # disable
  count = 0
  #address = google_compute_global_address.gateway.self_link

  source  = "GoogleCloudPlatform/lb-http/google//modules/serverless_negs"
  version = "~> 5.0"

  project              = var.project
  name                 = var.name
  create_address       = false
  cdn                  = false
  ssl                  = true
  use_ssl_certificates = true
  ssl_certificates     = [data.google_compute_ssl_certificate.default.self_link]
  https_redirect       = true
  quic                 = true

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
