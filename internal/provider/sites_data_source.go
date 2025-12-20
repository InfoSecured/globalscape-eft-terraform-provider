package provider

import (
	"context"

	"github.com/InfoSecured/globalscape-eft-terraform-provider/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &sitesDataSource{}

func NewSitesDataSource() datasource.DataSource {
	return &sitesDataSource{}
}

type sitesDataSource struct {
	client *client.Client
}

type sitesDataSourceModel struct {
	Sites []siteModel `tfsdk:"sites"`
}

type siteModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (d *sitesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sites"
}

func (d *sitesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List Globalscape EFT sites configured on the server.",
		Attributes: map[string]schema.Attribute{
			"sites": schema.ListNestedAttribute{
				MarkdownDescription: "Sites configured on the server.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Unique site identifier.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Site name/label.",
						},
					},
				},
			},
		},
	}
}

func (d *sitesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	if c, ok := req.ProviderData.(*client.Client); ok {
		d.client = c
	}
}

func (d *sitesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Unconfigured client", "the provider client was not initialized")
		return
	}

	var state sitesDataSourceModel

	sites, err := d.client.ListSites(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Unable to list sites", err.Error())
		return
	}

	for _, s := range sites {
		state.Sites = append(state.Sites, siteModel{
			ID:   types.StringValue(s.ID),
			Name: types.StringValue(s.Attributes.Name),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
