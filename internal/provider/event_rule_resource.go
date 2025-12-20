package provider

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/InfoSecured/globalscape-eft-terraform-provider/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &eventRuleResource{}
var _ resource.ResourceWithConfigure = &eventRuleResource{}
var _ resource.ResourceWithImportState = &eventRuleResource{}

func NewEventRuleResource() resource.Resource {
	return &eventRuleResource{}
}

type eventRuleResource struct {
	client *client.Client
}

type eventRuleResourceModel struct {
	ID                types.String `tfsdk:"id"`
	SiteID            types.String `tfsdk:"site_id"`
	AttributesJSON    types.String `tfsdk:"attributes_json"`
	RelationshipsJSON types.String `tfsdk:"relationships_json"`
}

func (r *eventRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_event_rule"
}

func (r *eventRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Globalscape EFT event rules via the REST API. The attributes and relationships payloads are provided as JSON strings.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Event rule identifier assigned by EFT.",
			},
			"site_id": schema.StringAttribute{
				MarkdownDescription: "Site identifier that owns the event rule.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"attributes_json": schema.StringAttribute{
				MarkdownDescription: "JSON body for the event rule attributes block as documented by EFT.",
				Required:            true,
			},
			"relationships_json": schema.StringAttribute{
				MarkdownDescription: "Optional JSON body for the event rule relationships block.",
				Optional:            true,
			},
		},
	}
}

func (r *eventRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if c, ok := req.ProviderData.(*client.Client); ok {
		r.client = c
	}
}

func (r *eventRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured client", "the provider client was not initialized")
		return
	}

	var plan eventRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	attrRaw, diags := parseJSON(plan.AttributesJSON)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	relRaw, diags := parseOptionalJSON(plan.RelationshipsJSON)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := client.EventRuleRequestData{
		Type:       "eventRule",
		Attributes: attrRaw,
	}
	if len(relRaw) > 0 {
		request.Relationships = relRaw
	}

	rule, err := r.client.CreateEventRule(ctx, plan.SiteID.ValueString(), request)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create event rule", err.Error())
		return
	}

	resp.Diagnostics.Append(setEventRuleStateFromAPI(rule, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *eventRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured client", "the provider client was not initialized")
		return
	}

	var state eventRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rule, err := r.client.GetEventRule(ctx, state.SiteID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read event rule", err.Error())
		return
	}

	resp.Diagnostics.Append(setEventRuleStateFromAPI(rule, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *eventRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured client", "the provider client was not initialized")
		return
	}

	var plan eventRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	attrRaw, diags := parseJSON(plan.AttributesJSON)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	relRaw, diags := parseOptionalJSON(plan.RelationshipsJSON)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := client.EventRuleRequestData{
		Type:       "eventRule",
		ID:         plan.ID.ValueString(),
		Attributes: attrRaw,
	}
	if len(relRaw) > 0 {
		request.Relationships = relRaw
	}

	rule, err := r.client.UpdateEventRule(ctx, plan.SiteID.ValueString(), plan.ID.ValueString(), request)
	if err != nil {
		resp.Diagnostics.AddError("Failed to update event rule", err.Error())
		return
	}

	resp.Diagnostics.Append(setEventRuleStateFromAPI(rule, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *eventRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured client", "the provider client was not initialized")
		return
	}

	var state eventRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteEventRule(ctx, state.SiteID.ValueString(), state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Failed to delete event rule", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

func setEventRuleStateFromAPI(rule *client.EventRule, model *eventRuleResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	model.ID = types.StringValue(rule.ID)

	attrs, err := normalizeRawJSON(rule.Attributes)
	if err != nil {
		diags.AddError("Failed to normalize event rule attributes", err.Error())
		return diags
	}
	model.AttributesJSON = types.StringValue(attrs)

	if len(rule.Relationships) > 0 {
		relationships, err := normalizeRawJSON(rule.Relationships)
		if err != nil {
			diags.AddError("Failed to normalize event rule relationships", err.Error())
			return diags
		}
		model.RelationshipsJSON = types.StringValue(relationships)
	} else {
		model.RelationshipsJSON = types.StringNull()
	}

	return diags
}

func (r *eventRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Invalid import identifier", "Expected identifier in the form <site_id>/<rule_id>")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("site_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}

func parseJSON(value types.String) (json.RawMessage, diag.Diagnostics) {
	if value.IsUnknown() || value.IsNull() {
		return nil, diag.Diagnostics{diag.NewErrorDiagnostic("Invalid JSON", "Value must be provided")}
	}

	var raw json.RawMessage
	if err := json.Unmarshal([]byte(value.ValueString()), &raw); err != nil {
		return nil, diag.Diagnostics{diag.NewErrorDiagnostic("Invalid JSON", err.Error())}
	}
	return raw, nil
}

func parseOptionalJSON(value types.String) (json.RawMessage, diag.Diagnostics) {
	if value.IsUnknown() || value.IsNull() || value.ValueString() == "" {
		return nil, nil
	}
	var raw json.RawMessage
	if err := json.Unmarshal([]byte(value.ValueString()), &raw); err != nil {
		return nil, diag.Diagnostics{diag.NewErrorDiagnostic("Invalid JSON", err.Error())}
	}
	return raw, nil
}

func normalizeRawJSON(raw json.RawMessage) (string, error) {
	if len(raw) == 0 {
		return "", nil
	}
	var v interface{}
	if err := json.Unmarshal(raw, &v); err != nil {
		return "", err
	}
	normalized, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(normalized), nil
}
