package flex

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

func SchemaVpcID() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Description: "VPC ID of the resource.",
		Optional:    true,
		Computed:    true,
		ForceNew:    true,
	}
}

func ComputedSchemaVpcID() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Description: "VPC ID of the resource.",
		Computed:    true,
	}
}

func SchemaSubnets() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "Subnets of the resource.",
		Optional:    true,
		Computed:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"subnet_id": {
					Type:        schema.TypeString,
					Description: "ID of the subnet.",
					Optional:    true,
					Computed:    true,
					ForceNew:    true,
				},
				"zone_name": {
					Type:        schema.TypeString,
					Description: "Zone name of the subnet.",
					Optional:    true,
					Computed:    true,
					ForceNew:    true,
				},
			},
		},
	}
}

func ComputedSchemaSubnets() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "Subnets of the resource.",
		Computed:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"subnet_id": {
					Type:        schema.TypeString,
					Description: "ID of the subnet.",
					Computed:    true,
				},
				"zone_name": {
					Type:        schema.TypeString,
					Description: "Zone name of the subnet.",
					Computed:    true,
				},
			},
		},
	}
}
