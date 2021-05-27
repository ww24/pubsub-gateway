terraform {
  required_version = "~> 0.15.3"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 3.69.0"
    }

    google-beta = {
      source  = "hashicorp/google-beta"
      version = "~> 3.69.0"
    }
  }

  backend "remote" {
    organization = "ww24"

    workspaces {
      name = "pubsub-gateway"
    }
  }
}
