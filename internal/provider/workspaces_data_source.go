package provider

import (
	"context"
	"fmt"
	"github.com/fstaoe/terraform-provider-structurizr/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// WorkspaceModel represents a workspace configured in the structurizr
type WorkspaceModel struct {
	ID           types.Int64  `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	APIKey       types.String `tfsdk:"api_key"`
	APISecret    types.String `tfsdk:"api_secret"`
	PublicURL    types.String `tfsdk:"public_url"`
	PrivateURL   types.String `tfsdk:"private_url"`
	ShareableURL types.String `tfsdk:"shareable_url"`
}

// WorkspacesModel is the response body for any CRU methods
type WorkspacesModel struct {
	Workspaces []WorkspaceModel `tfsdk:"workspaces"`
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &workspacesDataSource{}
	_ datasource.DataSourceWithConfigure = &workspacesDataSource{}
)

// NewWorkspacesDataSource is a helper function to simplify the provider implementation.
func NewWorkspacesDataSource() datasource.DataSource {
	return &workspacesDataSource{}
}

// workspacesDataSource is the data source implementation.
type workspacesDataSource struct {
	client *client.Manager
}

// Configure adds the provider configured client to the data source.
func (d *workspacesDataSource) Configure(
	_ context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Manager)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf(
				"Expected *client.Manager, got: %T. Please report this issue to the provider developers.",
				req.ProviderData,
			),
		)

		return
	}

	d.client = c
}

// Metadata returns the data source type name. It can be used to register other type of information
func (d *workspacesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workspaces"
}

// Schema defines the schema for the data source.
func (d *workspacesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"workspaces": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed:    true,
							Description: "The identifier of the Workspace used to perform further operations.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the Workspace",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "The description of the Workspace explaining roughly what it is about.",
						},
						"api_key": schema.StringAttribute{
							Computed:    true,
							Sensitive:   true,
							Description: "The API key specific to the Workspace used to perform operations such as update.",
						},
						"api_secret": schema.StringAttribute{
							Computed:    true,
							Sensitive:   true,
							Description: "The API secret key specific to the Workspace used to perform operations such as update.",
						},
						"public_url": schema.StringAttribute{
							Computed:    true,
							Description: "A public URL that does not require authentication to access the Workspace.",
						},
						"private_url": schema.StringAttribute{
							Computed:    true,
							Description: "A private URL that requires authentication to access the Workspace.",
						},
						"shareable_url": schema.StringAttribute{
							Computed:    true,
							Description: "A shareable URL that does not require authentication and it has randomly generated ID which can be deactivated.",
						},
					},
				},
			},
		},
	}
}

// Read fetches the Terraform state with the latest data.
func (d *workspacesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state WorkspacesModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("[READ] State: %s", state))

	res, err := d.client.GetWorkspaces(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read structurizr workspaces", err.Error())
		return
	}

	for _, workspace := range res.Workspaces {
		workspaceState := WorkspaceModel{
			ID:           types.Int64Value(workspace.ID),
			Name:         types.StringValue(workspace.Name),
			Description:  types.StringValue(workspace.Description),
			APIKey:       types.StringValue(workspace.APIKey),
			APISecret:    types.StringValue(workspace.APISecret),
			PublicURL:    types.StringValue(workspace.PublicURL),
			PrivateURL:   types.StringValue(workspace.PrivateURL),
			ShareableURL: types.StringValue(workspace.ShareableURL),
		}

		state.Workspaces = append(state.Workspaces, workspaceState)
	}

	tflog.Trace(ctx, fmt.Sprintf("[READ] Storing Workspaces: %+v", state))

	// Set refreshed state to see if there is a diff
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
