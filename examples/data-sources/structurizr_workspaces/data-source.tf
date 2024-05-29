data "structurizr_workspaces" "example" {}

output "ids" {
  value = data.structurizr_workspaces.example.workspaces.*.id
}