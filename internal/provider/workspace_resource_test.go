package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResourceWorkspace_New(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		switch r.Method {
		case http.MethodGet:
			_, err = w.Write([]byte(workspacesJSON))
		case http.MethodPost:
			_, err = w.Write([]byte(workspaceJSON))
		case http.MethodDelete:
			_, err = w.Write([]byte(testAccGenericResponse("Workspace Deleted")))
		}
		if err != nil {
			t.Errorf("error writing body: %s", err)
		}
	}))
	defer testServer.Close()

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkspace(testServer.URL),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("structurizr_workspace.test", "id", "1"),
					resource.TestCheckResourceAttr("structurizr_workspace.test", "name", "workspace 1"),
					resource.TestCheckResourceAttr("structurizr_workspace.test", "description", "description"),
					resource.TestCheckResourceAttr("structurizr_workspace.test", "api_key", "api-key"),
					resource.TestCheckResourceAttr("structurizr_workspace.test", "api_secret", "api-secret"),
					resource.TestCheckResourceAttr("structurizr_workspace.test", "private_url", "/workspace/1"),
					resource.TestCheckResourceAttr("structurizr_workspace.test", "public_url", "/share/1"),
					resource.TestCheckResourceAttr("structurizr_workspace.test", "shareable_url", "/shareable/1"),
				),
			},
			{
				ResourceName:      "structurizr_workspace.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceWorkspace(url string) string {
	return fmt.Sprintf(`%s
resource "structurizr_workspace" "test" {}
`, testAccProvider(url))
}
