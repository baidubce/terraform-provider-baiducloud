package flex

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

func SchemaResourceGroupID() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Description: "Resource group id of the resource. Effective upon creation, modifications are not supported currently.",
		Optional:    true,
		ForceNew:    true,
	}
}
