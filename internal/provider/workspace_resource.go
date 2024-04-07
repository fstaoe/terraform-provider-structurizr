package provider

import (
	"context"
	"fmt"
	"github.com/fstaoe/terraform-provider-structurizr/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
)

// WorkspaceResourceModel represents a workspace in the structurizr
type WorkspaceResourceModel struct {
	ID           types.Int64  `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	APIKey       types.String `tfsdk:"api_key"`
	APISecret    types.String `tfsdk:"api_secret"`
	PublicURL    types.String `tfsdk:"public_url"`
	PrivateURL   types.String `tfsdk:"private_url"`
	ShareableURL types.String `tfsdk:"shareable_url"`
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &workspaceResource{}
	_ resource.ResourceWithConfigure   = &workspaceResource{}
	_ resource.ResourceWithImportState = &workspaceResource{}
)

// NewWorkspaceResource is a helper function to simplify the provider implementation.
func NewWorkspaceResource() resource.Resource {
	return &workspaceResource{}
}

// workspaceResource is the resource implementation.
type workspaceResource struct {
	client *client.Structurizr
}

// Configure adds the provider configured client to the resource.
func (r *workspaceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Structurizr)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected resource Configure Type",
			fmt.Sprintf("Expected *client.Structurizr, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = c
}

// Metadata returns the resource type name.
func (r *workspaceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workspace"
}

// Schema defines the schema for the resource.
func (r *workspaceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Computed: true,
			},
			"description": schema.StringAttribute{
				Computed: true,
			},
			"api_key": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
			},
			"api_secret": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
			},
			"public_url": schema.StringAttribute{
				Computed: true,
			},
			"private_url": schema.StringAttribute{
				Computed: true,
			},
			"shareable_url": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *workspaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan WorkspaceResourceModel
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	// Create new workspace
	workspace, err := r.client.WithAdminAuth().CreateWorkspace(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Workspace",
			"Could not create workspace, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.Int64Value(workspace.ID)
	plan.Name = types.StringValue(workspace.Name)
	plan.Description = types.StringValue(workspace.Description)
	plan.APIKey = types.StringValue(workspace.APIKey)
	plan.APISecret = types.StringValue(workspace.APISecret)
	plan.PublicURL = types.StringValue(workspace.PublicURL)
	plan.PrivateURL = types.StringValue(workspace.PrivateURL)
	plan.ShareableURL = types.StringValue(workspace.ShareableURL)

	// Set state to fully populated data
	if resp.Diagnostics.Append(resp.State.Set(ctx, plan)...); resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *workspaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state WorkspaceResourceModel
	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed workspace value from Structurizr
	res, err := r.client.WithAdminAuth().GetWorkspaces(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Workspaces",
			"Could not read workspace ID "+strconv.FormatInt(state.ID.ValueInt64(), 10)+": "+err.Error(),
		)
		return
	}

	if len(res.Workspaces) == 0 {
		resp.Diagnostics.AddError(
			"There were no workspaces found",
			"The server did not return workspaces.",
		)
		return
	}

	workspace := res.FindByID(state.ID.ValueInt64())

	// Overwrite items with refreshed state
	state.ID = types.Int64Value(workspace.ID)
	state.Name = types.StringValue(workspace.Name)
	state.Description = types.StringValue(workspace.Description)
	state.APIKey = types.StringValue(workspace.APIKey)
	state.APISecret = types.StringValue(workspace.APISecret)
	state.PublicURL = types.StringValue(workspace.PublicURL)
	state.PrivateURL = types.StringValue(workspace.PrivateURL)
	state.ShareableURL = types.StringValue(workspace.ShareableURL)

	// Set refreshed state
	if resp.Diagnostics.Append(resp.State.Set(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *workspaceResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *workspaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state WorkspaceResourceModel

	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	// Delete existing workspace
	_, err := r.client.WithAdminAuth().DeleteWorkspace(ctx, state.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Structurizr Workspace",
			"Could not delete workspace, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *workspaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	attrID := path.Root("id")
	if attrID.Equal(path.Empty()) {
		resp.Diagnostics.AddError(
			"Resource Import Passthrough Missing Attribute Path",
			"Resource ImportState method call to ImportStatePassthroughID path must be set to a valid attribute path that can accept a string value.",
		)
	}

	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Resource Import Failed to Parse ID",
			"The import identifier must be an Integer.",
		)
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, attrID, id)...)
}
