package provider

import (
	"github.com/fstaoe/terraform-provider-structurizr/internal/acctest"
	"github.com/fstaoe/terraform-provider-structurizr/internal/util"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"net/http"
	"testing"
)

func TestDataSourceWorkspaces_Basic(t *testing.T) {
	endpoints := []*acctest.MockEndpoint{
		{
			Request: &acctest.MockRequest{Method: http.MethodGet, Uri: "/api/workspace"},
			Response: &acctest.MockResponse{
				StatusCode:  http.StatusOK,
				Body:        acctest.MockDataSourceWorkspacesBasic,
				ContentType: "application/json",
			},
		},
	}

	mockServer := acctest.NewMockServer(t, "Workspaces", endpoints)
	defer mockServer.Close()

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:          testAccDataSourceWorkspacesConfig(),
				ConfigVariables: config.Variables{"host": config.StringVariable(mockServer.URL)},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.structurizr_workspaces.test", "workspaces.#", "1"),
					resource.TestCheckResourceAttr("data.structurizr_workspaces.test", "workspaces.0.id", "1"),
					resource.TestCheckResourceAttr("data.structurizr_workspaces.test", "workspaces.0.name", "Workspace 0001"),
					resource.TestCheckResourceAttr("data.structurizr_workspaces.test", "workspaces.0.description", "Description"),
					resource.TestCheckResourceAttr("data.structurizr_workspaces.test", "workspaces.0.api_key", "691e0542-5c4d-4f74-be4a-38134a0aa0bf"),
					resource.TestCheckResourceAttr("data.structurizr_workspaces.test", "workspaces.0.api_secret", "8497f68e-75b9-431b-b067-cf86a074205c"),
					resource.TestCheckResourceAttr("data.structurizr_workspaces.test", "workspaces.0.private_url", "/workspace/1"),
					resource.TestCheckResourceAttr("data.structurizr_workspaces.test", "workspaces.0.public_url", "/share/1"),
					resource.TestCheckResourceAttr("data.structurizr_workspaces.test", "workspaces.0.shareable_url", ""),
				),
			},
		},
	})
}

func testAccDataSourceWorkspacesConfig() string {
	return util.ConfigCompose(testAccProvider(), `
data "structurizr_workspaces" "test" {}
`)
}
