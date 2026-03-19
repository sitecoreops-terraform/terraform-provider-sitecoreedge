package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sitecoreops-terraform/terraform-provider-sitecoreedge/pkg/apiclient"
)

var _ resource.Resource = &settingsResource{}

func NewSettingsResource() resource.Resource {
	return &settingsResource{}
}

type settingsResource struct {
	client *apiclient.Client
}

type settingsResourceModel struct {
	ContentCacheTtl       types.String `tfsdk:"content_cache_ttl"`
	ContentCacheAutoClear types.Bool   `tfsdk:"content_cache_auto_clear"`
	MediaCacheTtl         types.String `tfsdk:"media_cache_ttl"`
	MediaCacheAutoClear   types.Bool   `tfsdk:"media_cache_auto_clear"`
	TenantCacheAutoClear  types.Bool   `tfsdk:"tenant_cache_auto_clear"`
}

func (r *settingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "sitecoreedge_settings"
}

func (r *settingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages Sitecore Edge settings",
		Attributes: map[string]schema.Attribute{
			"content_cache_ttl": schema.StringAttribute{
				Description: "Content cache TTL (e.g., 04:00:00)",
				Optional:    true,
			},
			"content_cache_auto_clear": schema.BoolAttribute{
				Description: "Whether to auto-clear content cache",
				Optional:    true,
			},
			"media_cache_ttl": schema.StringAttribute{
				Description: "Media cache TTL (e.g., 04:00:00)",
				Optional:    true,
			},
			"media_cache_auto_clear": schema.BoolAttribute{
				Description: "Whether to auto-clear media cache",
				Optional:    true,
			},
			"tenant_cache_auto_clear": schema.BoolAttribute{
				Description: "Whether to auto-clear tenant cache",
				Optional:    true,
			},
		},
	}
}

func (r *settingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*apiclient.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *apiclient.Client, got: %T", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *settingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan settingsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	settings, err := r.patchSettings(plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create settings",
			err.Error(),
		)
		return
	}

	updateModelFromSettings(&plan, settings)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *settingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state settingsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	settings, err := r.client.GetSettings()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read settings",
			err.Error(),
		)
		return
	}

	updateModelFromSettings(&state, settings)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *settingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan settingsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	settings, err := r.patchSettings(plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update settings",
			err.Error(),
		)
		return
	}

	updateModelFromSettings(&plan, settings)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *settingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state settingsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *settingsResource) patchSettings(model settingsResourceModel) (*apiclient.Settings, error) {
	nullable := apiclient.NullableSettings{}

	if !model.ContentCacheTtl.IsNull() && !model.ContentCacheTtl.IsUnknown() {
		nullable.ContentCacheTtl = model.ContentCacheTtl.ValueStringPointer()
	}
	if !model.ContentCacheAutoClear.IsNull() && !model.ContentCacheAutoClear.IsUnknown() {
		nullable.ContentCacheAutoClear = model.ContentCacheAutoClear.ValueBoolPointer()
	}
	if !model.MediaCacheTtl.IsNull() && !model.MediaCacheTtl.IsUnknown() {
		nullable.MediaCacheTtl = model.MediaCacheTtl.ValueStringPointer()
	}
	if !model.MediaCacheAutoClear.IsNull() && !model.MediaCacheAutoClear.IsUnknown() {
		nullable.MediaCacheAutoClear = model.MediaCacheAutoClear.ValueBoolPointer()
	}
	if !model.TenantCacheAutoClear.IsNull() && !model.TenantCacheAutoClear.IsUnknown() {
		nullable.TenantCacheAutoClear = model.TenantCacheAutoClear.ValueBoolPointer()
	}

	patches := nullable.ToPatchOperations()
	if len(patches) == 0 {
		return r.client.GetSettings()
	}

	return r.client.PatchSettings(patches)
}

func updateModelFromSettings(model *settingsResourceModel, settings *apiclient.Settings) {
	model.ContentCacheTtl = types.StringValue(settings.ContentCacheTtl)
	model.ContentCacheAutoClear = types.BoolValue(settings.ContentCacheAutoClear)
	model.MediaCacheTtl = types.StringValue(settings.MediaCacheTtl)
	model.MediaCacheAutoClear = types.BoolValue(settings.MediaCacheAutoClear)
	model.TenantCacheAutoClear = types.BoolValue(settings.TenantCacheAutoClear)
}
