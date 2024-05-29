package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

//nolint:unparam
func protoV6ProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"structurizr": providerserver.NewProtocol6WithError(New("test")()),
	}
}

// testAccProvider is a shared configuration to combine with the actual
// test configuration so the Structurizr client is properly configured.
// It is also possible to use the STRUCTURIZR_ environment variables instead,
// such as updating the Makefile and running the testing through that tool.
func testAccProvider() string {
	return `variable "host" {}
provider "structurizr" {
    host = var.host
    admin_api_key = "test"
	tls_insecure = true
}
`
}
