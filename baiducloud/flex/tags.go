package flex

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

func SchemaTagsOnlySupportCreation() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeMap,
		Description: "Tags of the resource. Effective upon creation, modifications are not supported currently.",
		Optional:    true,
		Computed:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
}

func UpdatableTagsSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeMap,
		Optional:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: "Tags of the resource.",
	}
}

func ComputedSchemaTags() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeMap,
		Description: "Tags of the resource.",
		Computed:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
}
