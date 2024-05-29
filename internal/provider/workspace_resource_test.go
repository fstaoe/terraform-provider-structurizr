package provider

import (
	"github.com/fstaoe/terraform-provider-structurizr/internal/acctest"
	"github.com/fstaoe/terraform-provider-structurizr/internal/util"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"net/http"
	"testing"
)

func TestResourceWorkspace_Basic(t *testing.T) {
	endpoints := []*acctest.MockEndpoint{
		{
			Request: &acctest.MockRequest{Method: http.MethodPost, Uri: "/api/workspace", Body: util.StringPtr("")},
			Response: &acctest.MockResponse{
				StatusCode:  http.StatusOK,
				Body:        acctest.MockResourceWorkspaceBasicCreate,
				ContentType: "application/json",
			},
			Calls: 1,
		},
		{
			Request: &acctest.MockRequest{Method: http.MethodGet, Uri: "/api/workspace"},
			Response: &acctest.MockResponse{
				StatusCode:  http.StatusOK,
				Body:        acctest.MockResourceWorkspaceBasicGet,
				ContentType: "application/json",
			},
			Calls: 1,
		},
		{
			Request: &acctest.MockRequest{Method: http.MethodDelete, Uri: "/api/workspace/1"},
			Response: &acctest.MockResponse{
				StatusCode:  http.StatusOK,
				Body:        acctest.MockResourceWorkspaceBasicDelete,
				ContentType: "text/plain",
			},
			Calls: 1,
		},
	}

	mockServer := acctest.NewMockServer(t, "Workspace API", endpoints)
	defer mockServer.Close()

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		CheckDestroy: func(state *terraform.State) error {
			return acctest.AssertMockEndpointsCalls(endpoints)
		},
		Steps: []resource.TestStep{
			{
				Config:          testAccResourceWorkspaceConfigBasic(),
				ConfigVariables: config.Variables{"host": config.StringVariable(mockServer.URL)},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("structurizr_workspace.test", "id", "1"),
					resource.TestCheckResourceAttr("structurizr_workspace.test", "name", "Workspace 0001"),
					resource.TestCheckResourceAttr("structurizr_workspace.test", "description", "Description"),
					resource.TestCheckResourceAttr("structurizr_workspace.test", "api_key", "691e0542-5c4d-4f74-be4a-38134a0aa0bf"),
					resource.TestCheckResourceAttr("structurizr_workspace.test", "api_secret", "8497f68e-75b9-431b-b067-cf86a074205c"),
					resource.TestCheckResourceAttr("structurizr_workspace.test", "private_url", "/workspace/1"),
					resource.TestCheckResourceAttr("structurizr_workspace.test", "public_url", "/share/1"),
					resource.TestCheckResourceAttr("structurizr_workspace.test", "shareable_url", ""),
					resource.TestCheckResourceAttrSet("structurizr_workspace.test", "last_updated"),
				),
			},
		},
	})
}

func TestResourceWorkspace_Import(t *testing.T) {
	endpoints := []*acctest.MockEndpoint{
		{
			Request: &acctest.MockRequest{Method: http.MethodPost, Uri: "/api/workspace", Body: util.StringPtr("")},
			Response: &acctest.MockResponse{
				StatusCode:  http.StatusOK,
				Body:        acctest.MockResourceWorkspaceBasicCreate,
				ContentType: "application/json",
			},
			Calls: 1,
		},
		{
			Request: &acctest.MockRequest{Method: http.MethodGet, Uri: "/api/workspace"},
			Response: &acctest.MockResponse{
				StatusCode:  http.StatusOK,
				Body:        acctest.MockResourceWorkspaceBasicGet,
				ContentType: "application/json",
			},
			Calls: 2,
		},
		{
			Request: &acctest.MockRequest{Method: http.MethodDelete, Uri: "/api/workspace/1"},
			Response: &acctest.MockResponse{
				StatusCode:  http.StatusOK,
				Body:        acctest.MockResourceWorkspaceBasicDelete,
				ContentType: "text/plain",
			},
			Calls: 1,
		},
	}

	mockServer := acctest.NewMockServer(t, "Workspace API", endpoints)
	defer mockServer.Close()

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		CheckDestroy: func(state *terraform.State) error {
			return acctest.AssertMockEndpointsCalls(endpoints)
		},
		Steps: []resource.TestStep{
			{
				Config:          testAccResourceWorkspaceConfigBasic(),
				ConfigVariables: config.Variables{"host": config.StringVariable(mockServer.URL)},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("structurizr_workspace.test", "id", "1"),
					resource.TestCheckResourceAttr("structurizr_workspace.test", "name", "Workspace 0001"),
					resource.TestCheckResourceAttr("structurizr_workspace.test", "description", "Description"),
					resource.TestCheckResourceAttr("structurizr_workspace.test", "api_key", "691e0542-5c4d-4f74-be4a-38134a0aa0bf"),
					resource.TestCheckResourceAttr("structurizr_workspace.test", "api_secret", "8497f68e-75b9-431b-b067-cf86a074205c"),
					resource.TestCheckResourceAttr("structurizr_workspace.test", "private_url", "/workspace/1"),
					resource.TestCheckResourceAttr("structurizr_workspace.test", "public_url", "/share/1"),
					resource.TestCheckResourceAttr("structurizr_workspace.test", "shareable_url", ""),
					resource.TestCheckResourceAttrSet("structurizr_workspace.test", "last_updated"),
				),
			},
			{
				ConfigVariables:   config.Variables{"host": config.StringVariable(mockServer.URL)},
				ResourceName:      "structurizr_workspace.test",
				ImportStateId:     "1",
				ImportState:       true,
				ImportStateVerify: true,
				// The below attributes does not exist in the Structurizr API, therefore there is no value for it
				// during import.
				ImportStateVerifyIgnore: []string{"source", "source_checksum", "last_updated"},
			},
		},
	})
}

func TestResourceWorkspace_Update(t *testing.T) {
	endpoints := []*acctest.MockEndpoint{
		{
			Request: &acctest.MockRequest{Method: http.MethodPost, Uri: "/api/workspace", Body: util.StringPtr("")},
			Response: &acctest.MockResponse{
				StatusCode:  http.StatusOK,
				Body:        acctest.MockResourceWorkspaceBasicCreate,
				ContentType: "application/json",
			},
			Calls: 1,
		},
		{
			Request: &acctest.MockRequest{Method: http.MethodPut, Uri: "/api/workspace/1"},
			Response: &acctest.MockResponse{
				StatusCode:  http.StatusOK,
				Body:        acctest.MockResourceWorkspaceBasicGet,
				ContentType: "application/json",
			},
			Calls: 1,
		},
		{
			Request: &acctest.MockRequest{Method: http.MethodGet, Uri: "/api/workspace"},
			Response: &acctest.MockResponse{
				StatusCode:  http.StatusOK,
				Body:        acctest.MockResourceWorkspaceWithSourceGet,
				ContentType: "application/json",
			},
			Calls: 4,
		},
		{
			Request: &acctest.MockRequest{Method: http.MethodDelete, Uri: "/api/workspace/1"},
			Response: &acctest.MockResponse{
				StatusCode:  http.StatusOK,
				Body:        acctest.MockResourceWorkspaceBasicDelete,
				ContentType: "text/plain",
			},
			Calls: 1,
		},
	}

	mockServer := acctest.NewMockServer(t, "Workspace API", endpoints)
	defer mockServer.Close()

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		CheckDestroy: func(state *terraform.State) error {
			return acctest.AssertMockEndpointsCalls(endpoints)
		},
		Steps: []resource.TestStep{
			{
				Config:          testAccResourceWorkspaceConfigBasic(),
				ConfigVariables: config.Variables{"host": config.StringVariable(mockServer.URL)},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("structurizr_workspace.test", "id", "1"),
					resource.TestCheckResourceAttr("structurizr_workspace.test", "name", "Workspace 0001"),
					resource.TestCheckResourceAttr("structurizr_workspace.test", "description", "Description"),
					resource.TestCheckResourceAttr("structurizr_workspace.test", "api_key", "691e0542-5c4d-4f74-be4a-38134a0aa0bf"),
					resource.TestCheckResourceAttr("structurizr_workspace.test", "api_secret", "8497f68e-75b9-431b-b067-cf86a074205c"),
					resource.TestCheckResourceAttr("structurizr_workspace.test", "private_url", "/workspace/1"),
					resource.TestCheckResourceAttr("structurizr_workspace.test", "public_url", "/share/1"),
					resource.TestCheckResourceAttr("structurizr_workspace.test", "shareable_url", ""),
					resource.TestCheckResourceAttrSet("structurizr_workspace.test", "last_updated"),
				),
			},
			// Update
			{
				Config:          testAccResourceWorkspaceConfigBasicUpdate(),
				ConfigVariables: config.Variables{"host": config.StringVariable(mockServer.URL)},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("structurizr_workspace.test", "id", "1"),
					resource.TestCheckResourceAttr("structurizr_workspace.test", "name", "Workspace DSL"),
					resource.TestCheckResourceAttr("structurizr_workspace.test", "description", "Managed Workspace by DSL"),
					resource.TestCheckResourceAttr("structurizr_workspace.test", "api_key", "691e0542-5c4d-4f74-be4a-38134a0aa0bf"),
					resource.TestCheckResourceAttr("structurizr_workspace.test", "api_secret", "8497f68e-75b9-431b-b067-cf86a074205c"),
					resource.TestCheckResourceAttr("structurizr_workspace.test", "private_url", "/workspace/1"),
					resource.TestCheckResourceAttr("structurizr_workspace.test", "public_url", "/share/1"),
					resource.TestCheckResourceAttr("structurizr_workspace.test", "shareable_url", ""),
					resource.TestCheckResourceAttr("structurizr_workspace.test", "source", "testdata/workspace.dsl"),
					resource.TestCheckResourceAttr("structurizr_workspace.test", "source_checksum", "ba47f1dae6946adbad62496b6dd6b7a3"),
					resource.TestCheckResourceAttrSet("structurizr_workspace.test", "last_updated"),
				),
			},
		},
	})
}

func testAccResourceWorkspaceConfigBasic() string {
	return util.ConfigCompose(testAccProvider(), `resource "structurizr_workspace" "test" {}`)
}

func testAccResourceWorkspaceConfigBasicUpdate() string {
	return util.ConfigCompose(testAccProvider(), `
resource "structurizr_workspace" "test" {
    source          = "testdata/workspace.dsl"
    source_checksum = "ba47f1dae6946adbad62496b6dd6b7a3"
}
`)
}
