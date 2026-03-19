package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sitecoreops-terraform/terraform-provider-sitecoreedge/pkg/apiclient"
)

var _ datasource.DataSource = &tenantIDDataSource{}

func NewTenantIDDataSource() datasource.DataSource {
	return &tenantIDDataSource{}
}

type tenantIDDataSource struct {
	client *apiclient.Client
}

type tenantIDDataSourceModel struct {
	ID       types.String `tfsdk:"id"`
	TenantID types.String `tfsdk:"tenant_id"`
}

func (d *tenantIDDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "sitecoreedge_tenant"
}

func (d *tenantIDDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Gets the tenant ID for Sitecore Edge",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Dummy ID for the datasource",
				Computed:    true,
			},
			"tenant_id": schema.StringAttribute{
				Description: "The tenant ID",
				Computed:    true,
			},
		},
	}
}

func (d *tenantIDDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*apiclient.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *apiclient.Client, got: %T", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *tenantIDDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model tenantIDDataSourceModel

	tenantID, err := d.client.GetTenantID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get tenant ID",
			err.Error(),
		)
		return
	}

	model.ID = types.StringValue(tenantID)
	model.TenantID = types.StringValue(tenantID)

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
