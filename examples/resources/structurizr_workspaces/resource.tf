terraform {
  required_providers {
    structurizr = {
      source  = "fstaoe/structurizr"
      version = "0.1.0"
    }
  }
}
provider "structurizr" {
  host          = "http://localhost:8080"
  admin_api_key = "structurizr"
  tls_insecure  = true
}
resource "structurizr_workspace" "example" {}