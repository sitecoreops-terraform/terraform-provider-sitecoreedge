package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sitecoreops-terraform/terraform-provider-sitecoreedge/pkg/apiclient"
)

func NewTokenResource() resource.Resource {
	return &tokenResource{}
}

type tokenResource struct {
	client *apiclient.Client
}

type tokenResourceModel struct {
	Hash      types.String `tfsdk:"hash"`
	IsRevoked types.Bool   `tfsdk:"is_revoked"`
	Label     types.String `tfsdk:"label"`
	Scopes    types.List   `tfsdk:"scopes"`
	CreatedBy types.String `tfsdk:"created_by"`
	Created   types.String `tfsdk:"created"`
	Token     types.String `tfsdk:"token"`
}

func (r *tokenResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "sitecoreedge_token"
}

func (r *tokenResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Sitecore Edge API token",
		Attributes: map[string]schema.Attribute{
			"hash": schema.StringAttribute{
				Description: "The unique hash identifier of the token",
				Computed:    true,
			},
			"is_revoked": schema.BoolAttribute{
				Description: "Whether the token is revoked",
				Computed:    true,
			},
			"label": schema.StringAttribute{
				Description: "The label/name of the token",
				Required:    true,
			},
			"scopes": schema.ListAttribute{
				Description: "The scopes/permissions for the token",
				ElementType: types.StringType,
				Required:    true,
			},
			"created_by": schema.StringAttribute{
				Description: "The user who created the token",
				Optional:    true,
			},
			"created": schema.StringAttribute{
				Description: "When the token was created",
				Computed:    true,
			},
			"token": schema.StringAttribute{
				Description: "The actual token value (sensitive - only available after creation)",
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
}

func (r *tokenResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *tokenResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan tokenResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert scopes from types.List to []string
	scopes := make([]string, 0)
	if !plan.Scopes.IsNull() {
		diags = plan.Scopes.ElementsAs(ctx, &scopes, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	createdBy := "terraform"
	if !plan.CreatedBy.IsNull() {
		createdBy = plan.CreatedBy.ValueString()
	}

	input := &apiclient.TokenInput{
		CreatedBy: createdBy,
		Label:     plan.Label.ValueString(),
		Scopes:    scopes,
	}

	// Create the token
	tokenValue, err := r.client.CreateToken(input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create token",
			err.Error(),
		)
		return
	}

	// Get token details using the token value
	tokenDetails, err := r.client.GetTokenByToken(tokenValue)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get token details",
			err.Error(),
		)
		return
	}

	// Set the computed values
	plan.Hash = types.StringValue(tokenDetails.Hash)
	plan.IsRevoked = types.BoolValue(tokenDetails.IsRevoked)
	plan.Created = types.StringValue(tokenDetails.Created)
	plan.Token = types.StringValue(tokenValue)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *tokenResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state tokenResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get token by hash
	token, err := r.client.GetTokenByHash(state.Hash.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read token",
			err.Error(),
		)
		return
	}

	// Update the state with current values
	state.Label = types.StringValue(token.Label)
	state.IsRevoked = types.BoolValue(token.IsRevoked)
	state.CreatedBy = types.StringValue(token.CreatedBy)
	state.Created = types.StringValue(token.Created)

	// Convert scopes to types.List
	scopes, diags := types.ListValueFrom(ctx, types.StringType, token.Scopes)
	state.Scopes = scopes
	resp.Diagnostics.Append(diags...)

	// Token value is not stored in state after creation for security
	if state.Token.IsUnknown() {
		state.Token = types.StringNull()
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *tokenResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan tokenResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state tokenResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Only label can be updated (renamed)
	if !plan.Label.Equal(state.Label) {
		err := r.client.RenameToken(state.Hash.ValueString(), plan.Label.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to rename token",
				err.Error(),
			)
			return
		}
	}

	// Read the updated token to get current state
	token, err := r.client.GetTokenByHash(state.Hash.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read updated token",
			err.Error(),
		)
		return
	}

	// Update the plan with current values
	plan.Hash = types.StringValue(token.Hash)
	plan.IsRevoked = types.BoolValue(token.IsRevoked)
	plan.Label = types.StringValue(token.Label)
	plan.CreatedBy = types.StringValue(token.CreatedBy)
	plan.Created = types.StringValue(token.Created)

	scopes, diags := types.ListValueFrom(ctx, types.StringType, token.Scopes)
	plan.Scopes = scopes
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *tokenResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state tokenResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Revoke the token by hash
	err := r.client.RevokeTokenByHash(state.Hash.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to delete token",
			err.Error(),
		)
		return
	}
}

func (r *tokenResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import by hash
	token, err := r.client.GetTokenByHash(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to import token",
			err.Error(),
		)
		return
	}

	// Convert scopes to types.List
	scopes, diags := types.ListValueFrom(ctx, types.StringType, token.Scopes)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	state := tokenResourceModel{
		Hash:      types.StringValue(token.Hash),
		IsRevoked: types.BoolValue(token.IsRevoked),
		Label:     types.StringValue(token.Label),
		Scopes:    scopes,
		CreatedBy: types.StringValue(token.CreatedBy),
		Created:   types.StringValue(token.Created),
		Token:     types.StringNull(), // Token value is not available on import
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}
