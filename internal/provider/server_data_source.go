package provider

import (
	"context"

	"github.com/InfoSecured/globalscape-eft-terraform-provider/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &serverDataSource{}

func NewServerDataSource() datasource.DataSource {
	return &serverDataSource{}
}

type serverDataSource struct {
	client *client.Client
}

type serverDataSourceModel struct {
	ID               types.String                `tfsdk:"id"`
	Version          types.String                `tfsdk:"version"`
	General          serverGeneralModel          `tfsdk:"general"`
	ListenerSettings serverListenerSettingsModel `tfsdk:"listener_settings"`
	SMTP             serverSMTPModel             `tfsdk:"smtp"`
}

type serverGeneralModel struct {
	ConfigFilePath      types.String `tfsdk:"config_file_path"`
	EnableUtcInListings types.Bool   `tfsdk:"enable_utc_in_listings"`
	LastModifiedBy      types.String `tfsdk:"last_modified_by"`
	LastModifiedTime    types.Int64  `tfsdk:"last_modified_time"`
}

type serverListenerSettingsModel struct {
	AdminPort                  types.Int64 `tfsdk:"admin_port"`
	EnableRemoteAdministration types.Bool  `tfsdk:"enable_remote_administration"`
	ListenIPs                  types.List  `tfsdk:"listen_ips"`
}

type serverSMTPModel struct {
	Login             types.String `tfsdk:"login"`
	Password          types.String `tfsdk:"password"`
	Port              types.Int64  `tfsdk:"port"`
	SenderAddress     types.String `tfsdk:"sender_address"`
	SenderName        types.String `tfsdk:"sender_name"`
	Server            types.String `tfsdk:"server"`
	UseAuthentication types.Bool   `tfsdk:"use_authentication"`
	UseImplicitTLS    types.Bool   `tfsdk:"use_implicit_tls"`
}

func (d *serverDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server"
}

func (d *serverDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetch current Globalscape EFT server settings.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Server identifier provided by the API.",
				Computed:            true,
			},
			"version": schema.StringAttribute{
				MarkdownDescription: "Server version string.",
				Computed:            true,
			},
			"general": schema.SingleNestedAttribute{
				MarkdownDescription: "General configuration details.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"config_file_path":       schema.StringAttribute{Computed: true},
					"enable_utc_in_listings": schema.BoolAttribute{Computed: true},
					"last_modified_by":       schema.StringAttribute{Computed: true},
					"last_modified_time":     schema.Int64Attribute{Computed: true},
				},
			},
			"listener_settings": schema.SingleNestedAttribute{
				MarkdownDescription: "Administrative listener configuration.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"admin_port":                   schema.Int64Attribute{Computed: true},
					"enable_remote_administration": schema.BoolAttribute{Computed: true},
					"listen_ips":                   schema.ListAttribute{ElementType: types.StringType, Computed: true},
				},
			},
			"smtp": schema.SingleNestedAttribute{
				MarkdownDescription: "Outgoing SMTP settings.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"login":              schema.StringAttribute{Computed: true},
					"password":           schema.StringAttribute{Computed: true, Sensitive: true},
					"port":               schema.Int64Attribute{Computed: true},
					"sender_address":     schema.StringAttribute{Computed: true},
					"sender_name":        schema.StringAttribute{Computed: true},
					"server":             schema.StringAttribute{Computed: true},
					"use_authentication": schema.BoolAttribute{Computed: true},
					"use_implicit_tls":   schema.BoolAttribute{Computed: true},
				},
			},
		},
	}
}

func (d *serverDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	if c, ok := req.ProviderData.(*client.Client); ok {
		d.client = c
	}
}

func (d *serverDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Unconfigured client", "the provider client was not initialized")
		return
	}

	var data serverDataSourceModel

	server, err := d.client.GetServer(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Unable to query server", err.Error())
		return
	}

	data.ID = types.StringValue(server.ID)
	data.Version = types.StringValue(server.Attributes.Version)
	data.General = serverGeneralModel{
		ConfigFilePath:      types.StringValue(server.Attributes.General.ConfigFilePath),
		EnableUtcInListings: types.BoolValue(server.Attributes.General.EnableUtcInListings),
		LastModifiedBy:      types.StringValue(server.Attributes.General.LastModifiedBy),
		LastModifiedTime:    types.Int64Value(server.Attributes.General.LastModifiedUnixTime),
	}

	listenIPs, diag := types.ListValueFrom(ctx, types.StringType, server.Attributes.ListenerSettings.ListenIPs)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.ListenerSettings = serverListenerSettingsModel{
		AdminPort:                  types.Int64Value(server.Attributes.ListenerSettings.AdminPort),
		EnableRemoteAdministration: types.BoolValue(server.Attributes.ListenerSettings.EnableRemoteAdministration),
		ListenIPs:                  listenIPs,
	}

	data.SMTP = serverSMTPModel{
		Login:             types.StringValue(server.Attributes.SMTP.Login),
		Password:          types.StringValue(server.Attributes.SMTP.Password),
		Port:              types.Int64Value(server.Attributes.SMTP.Port),
		SenderAddress:     types.StringValue(server.Attributes.SMTP.SenderAddress),
		SenderName:        types.StringValue(server.Attributes.SMTP.SenderName),
		Server:            types.StringValue(server.Attributes.SMTP.Server),
		UseAuthentication: types.BoolValue(server.Attributes.SMTP.UseAuthentication),
		UseImplicitTLS:    types.BoolValue(server.Attributes.SMTP.UseImplicitTLS),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
