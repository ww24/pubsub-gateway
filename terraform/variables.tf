variable "location" {
  type    = string
  default = "asia-northeast1"
}

variable "project" {
  type = string
}

// credentials json value
variable "google_credentials" {
  type = string
}

variable "name" {
  type    = string
  default = "pubsub-gateway"
}

variable "gar_repository" {
  type    = string
  default = "ww24"
}

variable "image_name" {
  type    = string
  default = "pubsub-gateway"
}

variable "image_tag" {
  type    = string
  default = "latest"
}

// cloud run service account
variable "service_account" {
  type = string
}

variable "timeout" {
  type    = number
  default = 15
}

// iap
variable "oauth_client_id" {
  type = string
}

variable "oauth_client_secret" {
  type = string
}

variable "domains" {
  type = list(string)
}

// application environments
variable "mode" {
  type    = string
  default = "gateway"
}

variable "config_yaml" {
  type    = string
  default = ""
}

variable "default_origin" {
  type = string
}

variable "allow_origin_suffix" {
  type = string
}

variable "authorized_users" {
  type = string
}
