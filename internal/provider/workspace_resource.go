package provider

import (
	"context"
	"errors"
	"fmt"
	"github.com/fstaoe/terraform-provider-structurizr/internal/client"
	"github.com/fstaoe/terraform-provider-structurizr/internal/client/api/model"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strconv"
	"sync"
	"time"
)

// WorkspaceResourceModel represents a workspace in the structurizr
type WorkspaceResourceModel struct {
	ID               types.Int64  `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	APIKey           types.String `tfsdk:"api_key"`
	APISecret        types.String `tfsdk:"api_secret"`
	PublicURL        types.String `tfsdk:"public_url"`
	PrivateURL       types.String `tfsdk:"private_url"`
	ShareableURL     types.String `tfsdk:"shareable_url"`
	Source           types.String `tfsdk:"source"`
	SourceChecksum   types.String `tfsdk:"source_checksum"`
	SourcePassphrase types.String `tfsdk:"source_passphrase"`
	LastUpdated      types.String `tfsdk:"last_updated"`
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_     resource.Resource                     = &workspaceResource{}
	_     resource.ResourceWithConfigure        = &workspaceResource{}
	_     resource.ResourceWithImportState      = &workspaceResource{}
	_     resource.ResourceWithConfigValidators = &workspaceResource{}
	guard sync.Mutex
)

// NewWorkspaceResource is a helper function to simplify the provider implementation.
func NewWorkspaceResource() resource.Resource {
	return &workspaceResource{}
}

// workspaceResource is the resource implementation.
type workspaceResource struct {
	clientManager *client.Manager
}

// Configure adds the provider configured client to the resource.
func (r *workspaceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	m, ok := req.ProviderData.(*client.Manager)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected resource Configure Type",
			fmt.Sprintf(
				"Expected *client.Manager, got: %T. Please report this issue to the provider developers.",
				req.ProviderData,
			),
		)
		return
	}

	r.clientManager = m
}

// ConfigValidators returns a list of functions which will all be performed during validation.
func (r *workspaceResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		// Validate the schema defined attributes are either both null or both known values.
		resourcevalidator.RequiredTogether(
			path.MatchRoot("source"),
			path.MatchRoot("source_checksum"),
		),
	}
}

// Metadata returns the resource type name. It can be used to register other type of information.
func (r *workspaceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workspace"
}

// Schema defines the schema for the resource.
func (r *workspaceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
				Description:   "The identifier of the Workspace used to perform further operations.",
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
				Computed:      true,
				Sensitive:     true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Description:   "The API key specific to the Workspace used to perform operations such as update.",
			},
			"api_secret": schema.StringAttribute{
				Computed:      true,
				Sensitive:     true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Description:   "The API secret key specific to the Workspace used to perform operations such as update.",
			},
			"public_url": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Description:   "A public URL that does not require authentication to access the Workspace.",
			},
			"private_url": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Description:   "A private URL that requires authentication to access the Workspace.",
			},
			"shareable_url": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Description:   "A shareable URL that does not require authentication and it has randomly generated ID which can be deactivated.",
			},
			"source": schema.StringAttribute{
				Optional:    true,
				Description: "The DSL/JSON file representing a Workspace.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplaceIf(func(
					_ context.Context,
					req planmodifier.StringRequest,
					resp *stringplanmodifier.RequiresReplaceIfFuncResponse,
				) {
					if !req.ConfigValue.IsNull() && req.StateValue.IsNull() {
						return
					}

					resp.RequiresReplace = true
				},
					"If the value of this attribute is configured and removed, Terraform will destroy and recreate the resource.",
					"If the value of this attribute is configured and removed, Terraform will destroy and recreate the resource.",
				)},
			},
			"source_checksum": schema.StringAttribute{
				Optional:    true,
				Description: "The checksum of the source file.",
			},
			"source_passphrase": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "The passphrase to use when the client-side encryption is enabled on the workspace.",
			},
			"last_updated": schema.StringAttribute{
				Computed:    true,
				Description: "It provides the information when the Workspace was last updated.",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *workspaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Preventing race conditions when running Terraform with multiple resources
	guard.Lock()
	defer guard.Unlock()

	// Retrieve values from plan
	var state, plan WorkspaceResourceModel
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("[CREATE] State: %s Plan: %s", state, plan))

	workspace, err := r.clientManager.CreateWorkspace(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Workspace",
			fmt.Sprintf("Failed to create Workspace with error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("[CREATE] Setting Workspace %+v with State: %s Plan: %s", workspace, state, plan))

	state.ID = types.Int64Value(workspace.ID)
	state.Name = types.StringValue(workspace.Name)
	state.Description = types.StringValue(workspace.Description)
	state.APIKey = types.StringValue(workspace.APIKey)
	state.APISecret = types.StringValue(workspace.APISecret)
	state.PublicURL = types.StringValue(workspace.PublicURL)
	state.PrivateURL = types.StringValue(workspace.PrivateURL)
	state.ShareableURL = types.StringValue(workspace.ShareableURL)
	state.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	tflog.Trace(ctx, fmt.Sprintf("[CREATE] After Setting Workspace %+v with State: %s Plan: %s", workspace, state, plan))

	// The workspace will be updated on the remote server using it source when provided
	// and the state will be as well refreshed using the latest data from the remote server
	if plan.Source.ValueString() != "" {
		tflog.Trace(ctx, fmt.Sprintf("[CREATE] Updating Workspace %+v with State: %s Plan: %s", workspace, state, plan))

		err = r.clientManager.PushWorkspace(
			ctx,
			workspace.ID,
			workspace.APIKey,
			workspace.APISecret,
			plan.SourcePassphrase.ValueString(),
			plan.Source.ValueString(),
		)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating Workspace",
				fmt.Sprintf("Failed to update Workspace (id: %d) with error: %s", workspace.ID, err),
			)

			tflog.Trace(ctx, fmt.Sprintf("[CREATE] Rolling back Workspace %+v with State: %s Plan: %s", workspace, state, plan))

			// Rolling back workspace creation
			if _, err = r.clientManager.DeleteWorkspace(ctx, workspace.ID); err != nil {
				resp.Diagnostics.AddError(
					"Error rolling back Workspace",
					fmt.Sprintf("Failed to rollback Workspace (id: %d) creation with error: %s", workspace.ID, err),
				)
			}
			return
		}

		tflog.Trace(ctx, fmt.Sprintf("[CREATE] Refreshing Workspace %+v with State: %s Plan: %s", workspace, state, plan))

		updatedWorkspace, err := r.getWorkspaceByID(ctx, workspace.ID)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error retrieving Workspace",
				fmt.Sprintf("Failed to retrieve Workspace (id: %d) after updating with error: %s", workspace.ID, err),
			)
			return
		}

		tflog.Trace(ctx, fmt.Sprintf("[CREATE] Setting updated Workspace %+v with State: %s Plan: %s", workspace, state, plan))

		state.Name = types.StringValue(updatedWorkspace.Name)
		state.Description = types.StringValue(updatedWorkspace.Description)
		state.Source = plan.Source
		state.SourceChecksum = plan.SourceChecksum
		state.SourcePassphrase = plan.SourcePassphrase
	}

	tflog.Trace(ctx, fmt.Sprintf("[CREATE] Storing Workspace State: %+v", state))

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *workspaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state WorkspaceResourceModel
	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("[READ] State %s", state))

	workspace, err := r.getWorkspaceByID(ctx, state.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving Workspace",
			fmt.Sprintf("Failed to retrieve Workspace (id: %s) with error: %s", state.ID, err),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("[READ] Setting Workspace %+v to state %s", workspace, state))

	state.Name = types.StringValue(workspace.Name)
	state.Description = types.StringValue(workspace.Description)
	state.APIKey = types.StringValue(workspace.APIKey)
	state.APISecret = types.StringValue(workspace.APISecret)
	state.PublicURL = types.StringValue(workspace.PublicURL)
	state.PrivateURL = types.StringValue(workspace.PrivateURL)
	state.ShareableURL = types.StringValue(workspace.ShareableURL)
	state.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	tflog.Trace(ctx, fmt.Sprintf("[READ] Storing Workspace: %+v", state))

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *workspaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan WorkspaceResourceModel
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("[UPDATE] Plan %s", plan))

	err := r.clientManager.PushWorkspace(
		ctx,
		plan.ID.ValueInt64(),
		plan.APIKey.ValueString(),
		plan.APISecret.ValueString(),
		plan.SourcePassphrase.ValueString(),
		plan.Source.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Workspace",
			fmt.Sprintf("Failed to update Workspace (id: %s) with error: %s", plan.ID, err),
		)
		return
	}

	workspace, err := r.getWorkspaceByID(ctx, plan.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving Workspace",
			fmt.Sprintf("Failed to retrieve Workspace (id: %s) after updating with error: %s", plan.ID, err),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("[UPDATE] Setting Workspace %+v to state %s", workspace, plan))

	plan.Name = types.StringValue(workspace.Name)
	plan.Description = types.StringValue(workspace.Description)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	tflog.Trace(ctx, fmt.Sprintf("[UPDATE] Storing Workspace: %+v", plan))

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *workspaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state WorkspaceResourceModel
	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("[DELETE] State %s", state))

	if _, err := r.clientManager.DeleteWorkspace(ctx, state.ID.ValueInt64()); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Workspace",
			fmt.Sprintf("Failed to delete Workspace (id: %s) after updating with error: %s", state.ID, err),
		)
	}
}

// ImportState imports and store the current resource state from the remote server. This is ideal when migrating to Terraform.
func (r *workspaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	attrID := path.Root("id")
	if attrID.Equal(path.Empty()) {
		resp.Diagnostics.AddError(
			"Error importing Workspace",
			"Failed to import Workspace due to missing or invalid resource path",
		)
		return
	}

	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing Workspace ID",
			fmt.Sprintf("Failed to parse Workspace (id: %s) with error: %s", req.ID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, attrID, id)...)
}

func (r *workspaceResource) getWorkspaceByID(ctx context.Context, id int64) (*model.Workspace, error) {
	res, err := r.clientManager.GetWorkspaces(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list workspaces with error: %s", err)
	}

	if len(res.Workspaces) == 0 {
		return nil, errors.New("workspaces not found on remote server")
	}

	workspace := res.FindByID(id)
	if workspace == nil {
		return nil, errors.New("workspace not found")
	}

	return workspace, nil
}
