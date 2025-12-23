package provider

import (
	"context"
	"net/url"
	"strings"

	"github.com/InfoSecured/globalscape-eft-terraform-provider/internal/client"
	"github.com/InfoSecured/globalscape-eft-terraform-provider/internal/version"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ provider.Provider = &globalscapeProvider{}

func New() provider.Provider {
	return &globalscapeProvider{}
}

type globalscapeProvider struct{}

func (p *globalscapeProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = version.ProviderTypeName
	resp.Version = version.ProviderVersion
}

type providerModel struct {
	Host               types.String `tfsdk:"host"`
	Username           types.String `tfsdk:"username"`
	Password           types.String `tfsdk:"password"`
	AuthType           types.String `tfsdk:"auth_type"`
	InsecureSkipVerify types.Bool   `tfsdk:"insecure_skip_verify"`
}

func (p *globalscapeProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Provider for managing Globalscape EFT resources via the REST API.",
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				MarkdownDescription: "Base URL for the EFT admin API, for example https://eft.example.com:4450/admin.",
				Required:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Admin username with access to the REST API.",
				Required:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Admin password for the REST API.",
				Required:            true,
				Sensitive:           true,
			},
			"auth_type": schema.StringAttribute{
				MarkdownDescription: "Authentication type accepted by EFT (EFT or AD). Defaults to EFT.",
				Optional:            true,
			},
			"insecure_skip_verify": schema.BoolAttribute{
				MarkdownDescription: "Skip TLS verification when communicating with EFT. Useful for lab systems with self-signed certificates.",
				Optional:            true,
			},
		},
	}
}

func (p *globalscapeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config providerModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	host := strings.TrimSpace(config.Host.ValueString())
	user := config.Username.ValueString()
	password := config.Password.ValueString()
	authType := config.AuthType.ValueString()
	insecure := config.InsecureSkipVerify.ValueBool()

	if authType == "" {
		authType = "EFT"
	}

	if host == "" || user == "" || password == "" {
		resp.Diagnostics.AddError(
			"Missing provider configuration",
			"host, username, and password must all be provided",
		)
		return
	}

	parsedURL, err := url.Parse(host)
	if err != nil {
		resp.Diagnostics.AddError("Invalid host", err.Error())
		return
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		resp.Diagnostics.AddError(
			"Invalid host URL scheme",
			"Host must use http:// or https:// scheme. Got: "+parsedURL.Scheme,
		)
		return
	}

	apiClient, err := client.NewClient(ctx, client.Config{
		BaseURL:            host,
		Username:           user,
		Password:           password,
		AuthType:           authType,
		InsecureSkipVerify: insecure,
	})
	if err != nil {
		resp.Diagnostics.AddError("Failed to initialize client", err.Error())
		return
	}

	tflog.Info(ctx, "configured Globalscape EFT provider", map[string]any{"host": host})

	resp.DataSourceData = apiClient
	resp.ResourceData = apiClient
}

func (p *globalscapeProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewServerSMTPResource,
		NewSiteUserResource,
		NewEventRuleResource,
	}
}

func (p *globalscapeProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewServerDataSource,
		NewSitesDataSource,
	}
}
