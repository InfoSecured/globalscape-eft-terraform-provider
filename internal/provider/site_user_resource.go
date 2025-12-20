package provider

import (
	"context"

	"github.com/InfoSecured/globalscape-eft-terraform-provider/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &siteUserResource{}
var _ resource.ResourceWithConfigure = &siteUserResource{}

func NewSiteUserResource() resource.Resource {
	return &siteUserResource{}
}

type siteUserResource struct {
	client *client.Client
}

type siteUserResourceModel struct {
	ID                types.String `tfsdk:"id"`
	SiteID            types.String `tfsdk:"site_id"`
	LoginName         types.String `tfsdk:"login_name"`
	Password          types.String `tfsdk:"password"`
	PasswordType      types.String `tfsdk:"password_type"`
	DisplayName       types.String `tfsdk:"display_name"`
	Email             types.String `tfsdk:"email"`
	AccountEnabled    types.String `tfsdk:"account_enabled"`
	HomeFolderPath    types.String `tfsdk:"home_folder_path"`
	HomeFolderEnabled types.String `tfsdk:"home_folder_enabled"`
	HomeFolderRoot    types.String `tfsdk:"home_folder_root"`
}

var yesNoInherit = []string{"inherit", "yes", "no"}

func (r *siteUserResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_site_user"
}

func (r *siteUserResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Globalscape EFT users for a specific site.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "User identifier assigned by EFT.",
			},
			"site_id": schema.StringAttribute{
				MarkdownDescription: "Site identifier that owns the user.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"login_name": schema.StringAttribute{
				MarkdownDescription: "Unique login name for the user.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Password for EFT local accounts.",
				Optional:            true,
				Sensitive:           true,
			},
			"password_type": schema.StringAttribute{
				MarkdownDescription: "Password type as expected by EFT (for example `Default` or `Disabled`).",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("Default"),
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Friendly display name.",
				Optional:            true,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "User email address.",
				Optional:            true,
			},
			"account_enabled": schema.StringAttribute{
				MarkdownDescription: "Account enablement flag (`yes`, `no`, or `inherit`).",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("inherit"),
				Validators: []validator.String{
					stringvalidator.OneOf(yesNoInherit...),
				},
			},
			"home_folder_path": schema.StringAttribute{
				MarkdownDescription: "Path for the user's home folder.",
				Optional:            true,
			},
			"home_folder_enabled": schema.StringAttribute{
				MarkdownDescription: "Whether the home folder is enabled (`yes`, `no`, or `inherit`).",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("inherit"),
				Validators: []validator.String{
					stringvalidator.OneOf(yesNoInherit...),
				},
			},
			"home_folder_root": schema.StringAttribute{
				MarkdownDescription: "Controls if the home folder is treated as root (`yes`, `no`, or `inherit`).",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("inherit"),
				Validators: []validator.String{
					stringvalidator.OneOf(yesNoInherit...),
				},
			},
		},
	}
}

func (r *siteUserResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	if c, ok := req.ProviderData.(*client.Client); ok {
		r.client = c
	}
}

func (r *siteUserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured client", "the provider client was not initialized")
		return
	}

	var plan siteUserResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	attrs := plan.toAPIModel()

	user, err := r.client.CreateSiteUser(ctx, plan.SiteID.ValueString(), attrs)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create user", err.Error())
		return
	}

	passwordType := plan.PasswordType
	siteID := plan.SiteID

	plan.fromAPI(user)
	plan.SiteID = siteID
	plan.PasswordType = passwordType
	plan.Password = types.StringNull()

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *siteUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured client", "the provider client was not initialized")
		return
	}

	var state siteUserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, err := r.client.GetSiteUser(ctx, state.SiteID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read user", err.Error())
		return
	}

	passwordType := state.PasswordType
	state.fromAPI(user)
	state.PasswordType = passwordType
	state.Password = types.StringNull()

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *siteUserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured client", "the provider client was not initialized")
		return
	}

	var plan siteUserResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, err := r.client.UpdateSiteUser(ctx, plan.SiteID.ValueString(), plan.ID.ValueString(), plan.toAPIModel())
	if err != nil {
		resp.Diagnostics.AddError("Failed to update user", err.Error())
		return
	}

	passwordType := plan.PasswordType
	plan.fromAPI(user)
	plan.PasswordType = passwordType
	plan.Password = types.StringNull()

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *siteUserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured client", "the provider client was not initialized")
		return
	}

	var state siteUserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteSiteUser(ctx, state.SiteID.ValueString(), state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Failed to delete user", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

func (m *siteUserResourceModel) toAPIModel() client.UserAttributes {
	attr := client.UserAttributes{
		LoginName: m.LoginName.ValueString(),
	}

	if v := stringValueOrEmpty(m.AccountEnabled); v != "" {
		attr.AccountEnabled = v
	}

	if v := stringValueOrEmpty(m.HomeFolderRoot); v != "" {
		attr.HasHomeFolderAsRoot = v
	}

	if v := stringValueOrEmpty(m.DisplayName); v != "" || stringValueOrEmpty(m.Email) != "" {
		attr.Personal = &client.UserPersonal{
			Name:  stringValueOrEmpty(m.DisplayName),
			Email: stringValueOrEmpty(m.Email),
		}
	}

	if v := stringValueOrEmpty(m.Password); v != "" {
		attr.Password = &client.UserPassword{
			Type:  stringValueOrEmpty(m.PasswordType),
			Value: v,
		}
	}

	if enabled := stringValueOrEmpty(m.HomeFolderEnabled); enabled != "" || stringValueOrEmpty(m.HomeFolderPath) != "" {
		attr.HomeFolder = &client.UserHomeFolder{
			Enabled: enabled,
		}
		if path := stringValueOrEmpty(m.HomeFolderPath); path != "" {
			attr.HomeFolder.Value = &client.UserHomeFolderValue{Path: path}
		}
	}

	return attr
}

func (m *siteUserResourceModel) fromAPI(user *client.User) {
	m.ID = types.StringValue(user.ID)
	m.LoginName = types.StringValue(user.Attributes.LoginName)
	if user.Attributes.Personal != nil {
		m.DisplayName = types.StringValue(user.Attributes.Personal.Name)
		m.Email = types.StringValue(user.Attributes.Personal.Email)
	} else {
		m.DisplayName = types.StringNull()
		m.Email = types.StringNull()
	}
	if user.Attributes.AccountEnabled != "" {
		m.AccountEnabled = types.StringValue(user.Attributes.AccountEnabled)
	} else {
		m.AccountEnabled = types.StringNull()
	}

	if user.Attributes.HomeFolder != nil {
		m.HomeFolderEnabled = types.StringValue(user.Attributes.HomeFolder.Enabled)
		if user.Attributes.HomeFolder.Value != nil {
			m.HomeFolderPath = types.StringValue(user.Attributes.HomeFolder.Value.Path)
		} else {
			m.HomeFolderPath = types.StringNull()
		}
	} else {
		m.HomeFolderEnabled = types.StringNull()
		m.HomeFolderPath = types.StringNull()
	}

	if user.Attributes.HasHomeFolderAsRoot != "" {
		m.HomeFolderRoot = types.StringValue(user.Attributes.HasHomeFolderAsRoot)
	} else {
		m.HomeFolderRoot = types.StringNull()
	}
}

func stringValueOrEmpty(value types.String) string {
	if value.IsNull() || value.IsUnknown() {
		return ""
	}
	return value.ValueString()
}
