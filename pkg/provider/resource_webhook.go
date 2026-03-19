package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sitecoreops-terraform/terraform-provider-sitecoreedge/pkg/apiclient"
)

func NewWebhookResource() resource.Resource {
	return &webhookResource{}
}

type webhookResource struct {
	client *apiclient.Client
}

type webhookResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Label         types.String `tfsdk:"label"`
	URI           types.String `tfsdk:"uri"`
	Method        types.String `tfsdk:"method"`
	Headers       types.Map    `tfsdk:"headers"`
	Body          types.String `tfsdk:"body"`
	CreatedBy     types.String `tfsdk:"created_by"`
	ExecutionMode types.String `tfsdk:"execution_mode"`
	BodyInclude   types.String `tfsdk:"body_include"`
	Disabled      types.Bool   `tfsdk:"disabled"`
}

func (r *webhookResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "sitecoreedge_webhook"
}

func (r *webhookResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Sitecore Edge webhook",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the webhook",
				Computed:    true,
			},
			"label": schema.StringAttribute{
				Description: "Describes the purpose of the webhook.",
				Required:    true,
			},
			"uri": schema.StringAttribute{
				Description: "The URI of the webhook endpoint",
				Required:    true,
			},
			"method": schema.StringAttribute{
				Description: "The HTTP method to use when making the web request. Must be GET or POST.",
				Required:    true,
			},
			"headers": schema.MapAttribute{
				Description: "Custom headers to send when making the web request. Commonly used for authentication.",
				ElementType: types.StringType,
				Optional:    true,
			},
			"body": schema.StringAttribute{
				Description: "The body to post when making the web request. Only populated if ExecutionMode is OnEnd.",
				Optional:    true,
			},
			"created_by": schema.StringAttribute{
				Description: "The name of the user who created the webhook.",
				Optional:    true,
			},
			"execution_mode": schema.StringAttribute{
				Description: "Determines how the webhook is executed. Must be one of the following options: OnEnd, OnUpdate.",
				Required:    true,
			},
			"body_include": schema.StringAttribute{
				Description: "Additional JSON object to include in the body of the webhook request. Only populated if ExecutionMode is OnUpdate. This field must be a valid JSON object. Use jsonencode({}) to create the string.",
				Optional:    true,
			},
			"disabled": schema.BoolAttribute{
				Description: "Whether the webhook is disabled",
				Optional:    true,
			},
		},
	}
}

func (r *webhookResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *webhookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan webhookResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	headers := make(map[string]string)
	if !plan.Headers.IsNull() {
		plan.Headers.ElementsAs(ctx, &headers, false)
	}

	var bodyInclude *string
	if !plan.BodyInclude.IsNull() {
		bodyIncludeStr := plan.BodyInclude.ValueString()
		bodyInclude = &bodyIncludeStr
	}

	var body string
	if !plan.Body.IsNull() {
		body = plan.Body.ValueString()
	}

	input := &apiclient.WebhookInput{
		Label:         plan.Label.ValueString(),
		URI:           plan.URI.ValueString(),
		Method:        plan.Method.ValueString(),
		Headers:       headers,
		Body:          body,
		CreatedBy:     plan.CreatedBy.ValueString(),
		ExecutionMode: plan.ExecutionMode.ValueString(),
		BodyInclude:   bodyInclude,
	}

	webhook, err := r.client.CreateWebhook(input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create webhook",
			err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(webhook.ID)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *webhookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state webhookResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	webhook, err := r.client.GetWebhook(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read webhook",
			err.Error(),
		)
		return
	}

	state.Label = types.StringValue(webhook.Label)
	state.URI = types.StringValue(webhook.URI)
	state.Method = types.StringValue(webhook.Method)
	state.Body = types.StringValue(webhook.Body)
	state.CreatedBy = types.StringValue(webhook.CreatedBy)
	state.ExecutionMode = types.StringValue(webhook.ExecutionMode)

	if len(webhook.BodyInclude) > 0 {
		bodyIncludeBytes, _ := json.Marshal(webhook.BodyInclude)
		state.BodyInclude = types.StringValue(string(bodyIncludeBytes))
	}

	if webhook.Disabled != nil {
		state.Disabled = types.BoolValue(*webhook.Disabled)
	}

	headers, diags := types.MapValueFrom(ctx, types.StringType, webhook.Headers)
	state.Headers = headers
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *webhookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan webhookResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state webhookResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	headers := make(map[string]string)
	if !plan.Headers.IsNull() {
		plan.Headers.ElementsAs(ctx, &headers, false)
	}

	var bodyInclude *string
	if !plan.BodyInclude.IsNull() {
		bodyIncludeStr := plan.BodyInclude.ValueString()
		bodyInclude = &bodyIncludeStr
	}

	var body string
	if !plan.Body.IsNull() {
		body = plan.Body.ValueString()
	}

	input := &apiclient.WebhookInput{
		Label:         plan.Label.ValueString(),
		URI:           plan.URI.ValueString(),
		Method:        plan.Method.ValueString(),
		Headers:       headers,
		Body:          body,
		CreatedBy:     plan.CreatedBy.ValueString(),
		ExecutionMode: plan.ExecutionMode.ValueString(),
		BodyInclude:   bodyInclude,
	}

	webhook, err := r.client.UpdateWebhook(state.ID.ValueString(), input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update webhook",
			err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(webhook.ID)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *webhookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state webhookResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteWebhook(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to delete webhook",
			err.Error(),
		)
		return
	}
}

func (r *webhookResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	webhook, err := r.client.GetWebhook(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to import webhook",
			err.Error(),
		)
		return
	}

	state := webhookResourceModel{
		ID:            types.StringValue(webhook.ID),
		Label:         types.StringValue(webhook.Label),
		URI:           types.StringValue(webhook.URI),
		Method:        types.StringValue(webhook.Method),
		Body:          types.StringValue(webhook.Body),
		CreatedBy:     types.StringValue(webhook.CreatedBy),
		ExecutionMode: types.StringValue(webhook.ExecutionMode),
	}

	if len(webhook.Headers) > 0 {
		headers, diags := types.MapValueFrom(ctx, types.StringType, webhook.Headers)
		state.Headers = headers
		resp.Diagnostics.Append(diags...)
	}

	if len(webhook.BodyInclude) > 0 {
		bodyIncludeBytes, _ := json.Marshal(webhook.BodyInclude)
		state.BodyInclude = types.StringValue(string(bodyIncludeBytes))
	}

	if webhook.Disabled != nil {
		state.Disabled = types.BoolValue(*webhook.Disabled)
	}

	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}
