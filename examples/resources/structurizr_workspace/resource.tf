// Example of a managed workspace without source
resource "structurizr_workspace" "example" {}
// Example of a managed workspace with its source in a DSL file
resource "structurizr_workspace" "example_with_dsl" {
  source          = abspath("source/workspace.dsl")
  source_checksum = md5(file("source/workspace.dsl"))
}
// Example of a managed workspace with its source in a JSON file
resource "structurizr_workspace" "example_with_json" {
  source          = abspath("source/workspace.json")
  source_checksum = md5(file("source/workspace.json"))
}
// Example of a managed encrypted workspace
resource "structurizr_workspace" "example_with_encryption" {
  source            = abspath("source/workspace2encrypt.dsl")
  source_checksum   = md5(file("source/workspace2encrypt.dsl"))
  source_passphrase = "structurizr"
}