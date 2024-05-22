package provider

import (
	"context"
	"github.com/fstaoe/terraform-provider-structurizr/internal/client"
	"github.com/fstaoe/terraform-provider-structurizr/internal/client/api"
	"github.com/fstaoe/terraform-provider-structurizr/internal/client/cli"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"net/url"
	"os"
	"strconv"
)

// Ensure Structurizr satisfies various provider interfaces.
var _ provider.Provider = (*Structurizr)(nil)

func New(version string) func() provider.Provider {
	return func() provider.Provider { return &Structurizr{version: version} }
}

// Structurizr defines the provider implementation.
type Structurizr struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// StructurizrProviderModel describes the provider data model.
type StructurizrProviderModel struct {
	Host        types.String `tfsdk:"host"`
	AdminAPIKey types.String `tfsdk:"admin_api_key"`
	TLSInsecure types.Bool   `tfsdk:"tls_insecure"`
}

// Metadata returns the provider type name and version. It can be used to register other type of information
func (p *Structurizr) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "structurizr"
	resp.Version = p.version
}

// Schema defines the schema for the resource
func (p *Structurizr) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Required:    true,
				Description: "A fully qualified hostname (e.g. https://my.structurizr.instance)",
			},
			"admin_api_key": schema.StringAttribute{
				Sensitive:   true,
				Optional:    true,
				Description: "An API Key to be authorised against Structurizr API",
			},
			"tls_insecure": schema.BoolAttribute{
				Optional:    true,
				Description: "Disable TLS verification checks for self-hosted structurizr with self-signed certificates",
			},
		},
	}
}

func (p *Structurizr) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config StructurizrProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the attributes, it must be a known value.
	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Structurizr Host",
			"The provider cannot create a Structurizr client as there is an unknown configuration value for the host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the STRUCTURIZR_HOST environment variable.",
		)
	}

	if config.AdminAPIKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("admin_api_key"),
			"Unknown Structurizr Admin API Key",
			"The provider cannot create a Structurizr client as there is an unknown configuration value for the admin_api_key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the STRUCTURIZR_ADMIN_API_KEY environment variable.",
		)
	}

	if config.TLSInsecure.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("tls_insecure"),
			"Unknown Structurizr TLS Insecure",
			"The provider cannot create a Structurizr client as there is an unknown configuration value for the tls_insecure. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the STRUCTURIZR_TLS_INSECURE environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("STRUCTURIZR_HOST")
	adminApiKey := os.Getenv("STRUCTURIZR_ADMIN_API_KEY")

	var (
		tlsInsecure bool
		err         error
	)
	if v := os.Getenv("STRUCTURIZR_TLS_INSECURE"); v != "" {
		tlsInsecure, err = strconv.ParseBool(v)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to parse STRUCTURIZR_TLS_INSECURE environment variable",
				"An unexpected error occurred when parsing the STRUCTURIZR_TLS_INSECURE environment variable. "+
					"If the error is not clear, please contact the provider developers.\n\n"+
					"Error: "+err.Error(),
			)
		}
	}

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.AdminAPIKey.IsNull() {
		adminApiKey = config.AdminAPIKey.ValueString()
	}

	if !config.TLSInsecure.IsNull() {
		v, _ := config.TLSInsecure.ToBoolValue(ctx)
		tlsInsecure = v.ValueBool()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing Structurizr Host",
			"The provider cannot create a Structurizr client as there is a missing or empty value for the host. "+
				"Set the host value in the configuration or use the STRUCTURIZR_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	baseURL, err := url.Parse(host)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to parse Structurizr Host",
			"An unexpected error occurred when parsing the Structurizr Host. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Error: "+err.Error(),
		)
	}

	if adminApiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("admin_api_key"),
			"Missing Structurizr Admin API Key",
			"The provider cannot create a Structurizr client as there is a missing or empty value for the Structurizr Admin API Key. "+
				"Set the admin_api_key value in the configuration or use the STRUCTURIZR_ADMIN_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	cliWorkingDir, err := cli.WorkingDir(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to setting up the Structurizr CLI working directory",
			"An unexpected error occurred when setting up the Structurizr CLI working directory. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Error: "+err.Error(),
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new Structurizr client using the configuration values
	m := client.NewManager(
		api.NewClient(&api.Config{
			AdminAPIKey: adminApiKey,
			BaseURL:     baseURL,
			TLSInsecure: tlsInsecure,
			UserAgent:   api.DefaultUserAgent,
		}),
		cli.NewClient(&cli.Config{BaseURL: baseURL, WorkingDir: cliWorkingDir}, cli.DefaultCmdExec),
	)

	resp.DataSourceData = m
	resp.ResourceData = m
}

// Resources registers all available resource to be managed by Terraform
func (p *Structurizr) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewWorkspaceResource,
	}
}

// DataSources registers all available data sources that can be used to retrieve data
func (p *Structurizr) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewWorkspacesDataSource,
	}
}
