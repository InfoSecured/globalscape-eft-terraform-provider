package provider

import (
	"context"

	"github.com/InfoSecured/globalscape-eft-terraform-provider/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &serverSMTPResource{}
var _ resource.ResourceWithConfigure = &serverSMTPResource{}

func NewServerSMTPResource() resource.Resource {
	return &serverSMTPResource{}
}

type serverSMTPResource struct {
	client *client.Client
}

type serverSMTPResourceModel struct {
	ID                types.String `tfsdk:"id"`
	Login             types.String `tfsdk:"login"`
	Password          types.String `tfsdk:"password"`
	Port              types.Int64  `tfsdk:"port"`
	SenderAddress     types.String `tfsdk:"sender_address"`
	SenderName        types.String `tfsdk:"sender_name"`
	Server            types.String `tfsdk:"server"`
	UseAuthentication types.Bool   `tfsdk:"use_authentication"`
	UseImplicitTLS    types.Bool   `tfsdk:"use_implicit_tls"`
}

func (r *serverSMTPResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_smtp"
}

func (r *serverSMTPResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Globalscape EFT SMTP configuration.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Server identifier.",
				Computed:            true,
			},
			"login": schema.StringAttribute{Optional: true},
			"password": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"port":               schema.Int64Attribute{Required: true},
			"sender_address":     schema.StringAttribute{Required: true},
			"sender_name":        schema.StringAttribute{Required: true},
			"server":             schema.StringAttribute{Required: true},
			"use_authentication": schema.BoolAttribute{Optional: true},
			"use_implicit_tls":   schema.BoolAttribute{Optional: true},
		},
	}
}

func (r *serverSMTPResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	if c, ok := req.ProviderData.(*client.Client); ok {
		r.client = c
	}
}

func (r *serverSMTPResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured client", "the provider client was not initialized")
		return
	}

	var plan serverSMTPResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	server, err := r.client.UpdateServerSMTP(ctx, plan.toAPIModel())
	if err != nil {
		resp.Diagnostics.AddError("Failed to update SMTP settings", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, fromServerToSMTPModel(server))...)
}

func (r *serverSMTPResource) Read(ctx context.Context, _ resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured client", "the provider client was not initialized")
		return
	}

	server, err := r.client.GetServer(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read server", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, fromServerToSMTPModel(server))...)
}

func (r *serverSMTPResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured client", "the provider client was not initialized")
		return
	}

	var plan serverSMTPResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	server, err := r.client.UpdateServerSMTP(ctx, plan.toAPIModel())
	if err != nil {
		resp.Diagnostics.AddError("Failed to update SMTP settings", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, fromServerToSMTPModel(server))...)
}

func (r *serverSMTPResource) Delete(ctx context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.State.RemoveResource(ctx)
}

func (m *serverSMTPResourceModel) toAPIModel() client.SMTPSettings {
	model := client.SMTPSettings{}

	if !m.Login.IsNull() && !m.Login.IsUnknown() {
		model.Login = m.Login.ValueString()
	}
	if !m.Password.IsNull() && !m.Password.IsUnknown() {
		model.Password = m.Password.ValueString()
	}
	model.Port = m.Port.ValueInt64()
	model.SenderAddress = m.SenderAddress.ValueString()
	model.SenderName = m.SenderName.ValueString()
	model.Server = m.Server.ValueString()
	if !m.UseAuthentication.IsNull() && !m.UseAuthentication.IsUnknown() {
		model.UseAuthentication = m.UseAuthentication.ValueBool()
	}
	if !m.UseImplicitTLS.IsNull() && !m.UseImplicitTLS.IsUnknown() {
		model.UseImplicitTLS = m.UseImplicitTLS.ValueBool()
	}

	return model
}

func fromServerToSMTPModel(server *client.Server) *serverSMTPResourceModel {
	return &serverSMTPResourceModel{
		ID:                types.StringValue(server.ID),
		Login:             types.StringValue(server.Attributes.SMTP.Login),
		Password:          types.StringValue(server.Attributes.SMTP.Password),
		Port:              types.Int64Value(server.Attributes.SMTP.Port),
		SenderAddress:     types.StringValue(server.Attributes.SMTP.SenderAddress),
		SenderName:        types.StringValue(server.Attributes.SMTP.SenderName),
		Server:            types.StringValue(server.Attributes.SMTP.Server),
		UseAuthentication: types.BoolValue(server.Attributes.SMTP.UseAuthentication),
		UseImplicitTLS:    types.BoolValue(server.Attributes.SMTP.UseImplicitTLS),
	}
}
