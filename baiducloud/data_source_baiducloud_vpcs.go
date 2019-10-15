/*
Use this data source to query vpc list.

Example Usage

```hcl
data "baiducloud_vpcs" "default" {
    name="test-vpc"
}

output "cidr" {
  value = "${data.baiducloud_vpcs.default.vpcs.0.cidr}"
}
```
*/
package baiducloud

import (
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudVpcs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudVpcsRead,

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Description: "ID of the specific VPC to retrieve.",
				Optional:    true,
				ForceNew:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the specific VPC to retrieve.",
				Optional:    true,
				ForceNew:    true,
			},
			"output_file": {
				Type:        schema.TypeString,
				Description: "Output file for saving result.",
				Optional:    true,
				ForceNew:    true,
			},

			// Attributes used for result
			"vpcs": {
				Type:        schema.TypeList,
				Description: "Result of VPCs.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:        schema.TypeString,
							Description: "ID of the VPC.",
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "Name of the VPC.",
							Computed:    true,
						},
						"is_default": {
							Type:        schema.TypeBool,
							Description: "Specify if it is the default VPC.",
							Computed:    true,
						},
						"cidr": {
							Type:        schema.TypeString,
							Description: "CIDR block of the VPC.",
							Computed:    true,
						},
						"description": {
							Type:        schema.TypeString,
							Description: "Description of the VPC.",
							Computed:    true,
						},
						"route_table_id": {
							Type:        schema.TypeString,
							Description: "Route table ID of the VPC.",
							Computed:    true,
						},
						"secondary_cidrs": {
							Type:        schema.TypeList,
							Description: "The secondary cidr list of the VPC. They will not be repeated.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed: true,
						},
						"tags": tagsComputedSchema(),
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudVpcsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	vpcService := VpcService{client}

	var (
		vpcId      string
		name       string
		outputFile string
	)
	if v := d.Get("vpc_id").(string); v != "" {
		vpcId = v
	}
	if v := d.Get("name").(string); v != "" {
		name = v
	}
	if v := d.Get("output_file").(string); v != "" {
		outputFile = v
	}

	action := "Query VPCs " + vpcId + "_" + name

	vpcs, err := vpcService.ListAllVpcs()
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_vpcs", action, BCESDKGoERROR)
	}

	vpcsResult := make([]map[string]interface{}, 0)
	for _, vpc := range vpcs {
		if (vpcId != "" && vpcId != vpc.VPCID) ||
			(name != "" && name != vpc.Name) {
			continue
		}

		vpcMap := make(map[string]interface{})
		vpcMap["vpc_id"] = vpc.VPCID
		vpcMap["name"] = vpc.Name
		vpcMap["is_default"] = vpc.IsDefault
		vpcMap["cidr"] = vpc.Cidr
		vpcMap["description"] = vpc.Description
		vpcMap["secondary_cidrs"] = vpc.SecondaryCidr
		vpcMap["tags"] = flattenTagsToMap(vpc.Tags)

		res, err := vpcService.GetRouteTableDetail("", vpcId)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_vpc", action, BCESDKGoERROR)
		}
		vpcMap["route_table_id"] = res.RouteTableId

		vpcsResult = append(vpcsResult, vpcMap)
	}
	addDebug(action, vpcsResult)

	if err := d.Set("vpcs", vpcsResult); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_vpcs", action, BCESDKGoERROR)
	}

	d.SetId(resource.UniqueId())

	if outputFile != "" {
		if err := writeToFile(outputFile, vpcsResult); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_vpcs", action, BCESDKGoERROR)
		}
	}

	return nil
}
