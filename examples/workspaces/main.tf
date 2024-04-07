terraform {
  required_providers {
    structurizr = {
      source  = "fstaoe/structurizr"
      version = "0.1.0"
    }
  }
}
// Initialise the provider
provider "structurizr" {
  host          = "http://localhost:8080"
  admin_api_key = "structurizr"
  tls_insecure  = true
}
// Create a new Workspace
resource "structurizr_workspace" "example" {}
// Output the id of the first workspace
data "structurizr_workspaces" "example" {}
output "structurizr_workspaces" {
  value = data.structurizr_workspaces.example.workspaces.0.id
}