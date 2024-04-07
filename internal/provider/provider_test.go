package provider

import (
	"fmt"
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
func testAccProvider(url string) string {
	return fmt.Sprintf(`provider "structurizr" {
	host = "%s"
  	admin_api_key = "test"
  	tls_insecure = true
}
`, url)
}

func testAccGenericResponse(msg string) string {
	return fmt.Sprintf(`{
	"success": true,
	"message": "%s",
	"revision": 1
}
`, msg)
}

var workspacesJSON = fmt.Sprintf(`{"workspaces": [%s]}`, workspaceJSON)
var workspaceJSON = `{
	"id": 1,
	"name": "workspace 1",
	"description": "description",
	"apiKey": "api-key",
	"apiSecret": "api-secret",
	"privateUrl": "/workspace/1",
	"publicUrl": "/share/1",
	"shareableUrl": "/shareable/1"
}
`
