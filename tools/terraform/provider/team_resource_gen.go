// Code generated by terraform-plugin-framework-generator DO NOT EDIT.

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func TeamResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"agent_options": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				Description:         "The desired team's ID.",
				MarkdownDescription: "The desired team's ID.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "Team name",
				MarkdownDescription: "Team name",
			},
			"secrets": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"created_at": schema.StringAttribute{
							Computed: true,
						},
						"secret": schema.StringAttribute{
							Computed: true,
						},
						"team_id": schema.Int64Attribute{
							Computed: true,
						},
					},
					CustomType: SecretsType{
						ObjectType: types.ObjectType{
							AttrTypes: SecretsValue{}.AttributeTypes(ctx),
						},
					},
				},
				Computed: true,
			},
		},
	}
}

type TeamModel struct {
	AgentOptions types.String `tfsdk:"agent_options"`
	Description  types.String `tfsdk:"description"`
	Id           types.Int64  `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Secrets      types.List   `tfsdk:"secrets"`
}

var _ basetypes.ObjectTypable = SecretsType{}

type SecretsType struct {
	basetypes.ObjectType
}

func (t SecretsType) Equal(o attr.Type) bool {
	other, ok := o.(SecretsType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t SecretsType) String() string {
	return "SecretsType"
}

func (t SecretsType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	createdAtAttribute, ok := attributes["created_at"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`created_at is missing from object`)

		return nil, diags
	}

	createdAtVal, ok := createdAtAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`created_at expected to be basetypes.StringValue, was: %T`, createdAtAttribute))
	}

	secretAttribute, ok := attributes["secret"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`secret is missing from object`)

		return nil, diags
	}

	secretVal, ok := secretAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`secret expected to be basetypes.StringValue, was: %T`, secretAttribute))
	}

	teamIdAttribute, ok := attributes["team_id"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`team_id is missing from object`)

		return nil, diags
	}

	teamIdVal, ok := teamIdAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`team_id expected to be basetypes.Int64Value, was: %T`, teamIdAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return SecretsValue{
		CreatedAt: createdAtVal,
		Secret:    secretVal,
		TeamId:    teamIdVal,
		state:     attr.ValueStateKnown,
	}, diags
}

func NewSecretsValueNull() SecretsValue {
	return SecretsValue{
		state: attr.ValueStateNull,
	}
}

func NewSecretsValueUnknown() SecretsValue {
	return SecretsValue{
		state: attr.ValueStateUnknown,
	}
}

func NewSecretsValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (SecretsValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing SecretsValue Attribute Value",
				"While creating a SecretsValue value, a missing attribute value was detected. "+
					"A SecretsValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("SecretsValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid SecretsValue Attribute Type",
				"While creating a SecretsValue value, an invalid attribute value was detected. "+
					"A SecretsValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("SecretsValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("SecretsValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra SecretsValue Attribute Value",
				"While creating a SecretsValue value, an extra attribute value was detected. "+
					"A SecretsValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra SecretsValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewSecretsValueUnknown(), diags
	}

	createdAtAttribute, ok := attributes["created_at"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`created_at is missing from object`)

		return NewSecretsValueUnknown(), diags
	}

	createdAtVal, ok := createdAtAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`created_at expected to be basetypes.StringValue, was: %T`, createdAtAttribute))
	}

	secretAttribute, ok := attributes["secret"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`secret is missing from object`)

		return NewSecretsValueUnknown(), diags
	}

	secretVal, ok := secretAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`secret expected to be basetypes.StringValue, was: %T`, secretAttribute))
	}

	teamIdAttribute, ok := attributes["team_id"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`team_id is missing from object`)

		return NewSecretsValueUnknown(), diags
	}

	teamIdVal, ok := teamIdAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`team_id expected to be basetypes.Int64Value, was: %T`, teamIdAttribute))
	}

	if diags.HasError() {
		return NewSecretsValueUnknown(), diags
	}

	return SecretsValue{
		CreatedAt: createdAtVal,
		Secret:    secretVal,
		TeamId:    teamIdVal,
		state:     attr.ValueStateKnown,
	}, diags
}

func NewSecretsValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) SecretsValue {
	object, diags := NewSecretsValue(attributeTypes, attributes)

	if diags.HasError() {
		// This could potentially be added to the diag package.
		diagsStrings := make([]string, 0, len(diags))

		for _, diagnostic := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf(
				"%s | %s | %s",
				diagnostic.Severity(),
				diagnostic.Summary(),
				diagnostic.Detail()))
		}

		panic("NewSecretsValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t SecretsType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewSecretsValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewSecretsValueUnknown(), nil
	}

	if in.IsNull() {
		return NewSecretsValueNull(), nil
	}

	attributes := map[string]attr.Value{}

	val := map[string]tftypes.Value{}

	err := in.As(&val)

	if err != nil {
		return nil, err
	}

	for k, v := range val {
		a, err := t.AttrTypes[k].ValueFromTerraform(ctx, v)

		if err != nil {
			return nil, err
		}

		attributes[k] = a
	}

	return NewSecretsValueMust(SecretsValue{}.AttributeTypes(ctx), attributes), nil
}

func (t SecretsType) ValueType(ctx context.Context) attr.Value {
	return SecretsValue{}
}

var _ basetypes.ObjectValuable = SecretsValue{}

type SecretsValue struct {
	CreatedAt basetypes.StringValue `tfsdk:"created_at"`
	Secret    basetypes.StringValue `tfsdk:"secret"`
	TeamId    basetypes.Int64Value  `tfsdk:"team_id"`
	state     attr.ValueState
}

func (v SecretsValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 3)

	var val tftypes.Value
	var err error

	attrTypes["created_at"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["secret"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["team_id"] = basetypes.Int64Type{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 3)

		val, err = v.CreatedAt.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["created_at"] = val

		val, err = v.Secret.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["secret"] = val

		val, err = v.TeamId.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["team_id"] = val

		if err := tftypes.ValidateValue(objectType, vals); err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(objectType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(objectType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(objectType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Object state in ToTerraformValue: %s", v.state))
	}
}

func (v SecretsValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v SecretsValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v SecretsValue) String() string {
	return "SecretsValue"
}

func (v SecretsValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	objVal, diags := types.ObjectValue(
		map[string]attr.Type{
			"created_at": basetypes.StringType{},
			"secret":     basetypes.StringType{},
			"team_id":    basetypes.Int64Type{},
		},
		map[string]attr.Value{
			"created_at": v.CreatedAt,
			"secret":     v.Secret,
			"team_id":    v.TeamId,
		})

	return objVal, diags
}

func (v SecretsValue) Equal(o attr.Value) bool {
	other, ok := o.(SecretsValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.CreatedAt.Equal(other.CreatedAt) {
		return false
	}

	if !v.Secret.Equal(other.Secret) {
		return false
	}

	if !v.TeamId.Equal(other.TeamId) {
		return false
	}

	return true
}

func (v SecretsValue) Type(ctx context.Context) attr.Type {
	return SecretsType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v SecretsValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"created_at": basetypes.StringType{},
		"secret":     basetypes.StringType{},
		"team_id":    basetypes.Int64Type{},
	}
}
