package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDataSourceWorkspaces(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(workspacesJSON))
		if err != nil {
			t.Errorf("error writing body: %s", err)
		}
	}))
	defer testServer.Close()

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccProvider(testServer.URL) + testAccDataSourceWorkspacesConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.structurizr_workspaces.test", "workspaces.#", "1"),
					resource.TestCheckResourceAttr("data.structurizr_workspaces.test", "workspaces.0.id", "1"),
					resource.TestCheckResourceAttr("data.structurizr_workspaces.test", "workspaces.0.name", "workspace 1"),
					resource.TestCheckResourceAttr("data.structurizr_workspaces.test", "workspaces.0.description", "description"),
					resource.TestCheckResourceAttr("data.structurizr_workspaces.test", "workspaces.0.api_key", "api-key"),
					resource.TestCheckResourceAttr("data.structurizr_workspaces.test", "workspaces.0.api_secret", "api-secret"),
					resource.TestCheckResourceAttr("data.structurizr_workspaces.test", "workspaces.0.private_url", "/workspace/1"),
					resource.TestCheckResourceAttr("data.structurizr_workspaces.test", "workspaces.0.public_url", "/share/1"),
					resource.TestCheckResourceAttr("data.structurizr_workspaces.test", "workspaces.0.shareable_url", "/shareable/1"),
				),
			},
		},
	})
}

var testAccDataSourceWorkspacesConfig = `
data "structurizr_workspaces" "test" {}
`
