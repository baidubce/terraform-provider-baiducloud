/*
Use this data source to query subnet list.

Example Usage

```hcl
data "baiducloud_subnets" "default" {
 vpc_id = "vpc-y4p102r3mz6m"
}

output "subnets" {
 value = "${data.baiducloud_subnets.default.subnets}"
}
```
*/
package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudSubnets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudSubnetsRead,

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Description: "VPC ID for subnets to retrieve.",
				Optional:    true,
				ForceNew:    true,
			},
			"zone_name": {
				Type:        schema.TypeString,
				Description: "Specify the zone name for subnets.",
				Optional:    true,
				ForceNew:    true,
			},
			"subnet_type": {
				Type:         schema.TypeString,
				Description:  "Specify the subnet type for subnets.",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateSubnetType(),
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Description: "ID of the subnet.",
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
			"subnets": {
				Type:        schema.TypeList,
				Description: "Result of the subnets.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "Name of the subnet.",
							Computed:    true,
						},
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
						"cidr": {
							Type:        schema.TypeString,
							Description: "CIDR block of the subnet.",
							Computed:    true,
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Description: "VPC ID of the subnet.",
							Computed:    true,
						},
						"subnet_type": {
							Type:        schema.TypeString,
							Description: "Type of the subnet.",
							Computed:    true,
						},
						"description": {
							Type:        schema.TypeString,
							Description: "Description of the subnet.",
							Computed:    true,
						},
						"available_ip": {
							Type:        schema.TypeInt,
							Description: "Available IP address of the subnet.",
							Computed:    true,
						},
						"tags": tagsComputedSchema(),
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudSubnetsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	vpcService := VpcService{client}

	var (
		vpcID      string
		zoneName   string
		subnetType string
		subnetID   string
		outputFile string
	)
	if v, ok := d.GetOk("vpc_id"); ok {
		vpcID = v.(string)
	}
	if v, ok := d.GetOk("zone_name"); ok {
		zoneName = v.(string)
	}
	if v, ok := d.GetOk("subnet_type"); ok {
		subnetType = v.(string)
	}
	if v, ok := d.GetOk("subnet_id"); ok {
		subnetID = v.(string)
	}
	if v, ok := d.GetOk("output_file"); ok {
		outputFile = v.(string)
	}

	action := "Query subnets " + vpcID + "_" + zoneName + "_" + subnetType + "_" + subnetID

	args := &vpc.ListSubnetArgs{
		VpcId:      vpcID,
		ZoneName:   zoneName,
		SubnetType: vpc.SubnetType(subnetType),
	}
	subnets, err := vpcService.ListAllSubnets(args)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_subnets", action, BCESDKGoERROR)
	}

	subnetsResult := make([]map[string]interface{}, 0, len(subnets))
	for _, subnet := range subnets {
		if subnetID != "" && subnetID != subnet.SubnetId {
			continue
		}

		subnetMap := make(map[string]interface{})
		subnetMap["name"] = subnet.Name
		subnetMap["subnet_id"] = subnet.SubnetId
		subnetMap["zone_name"] = subnet.ZoneName
		subnetMap["cidr"] = subnet.Cidr
		subnetMap["vpc_id"] = subnet.VPCId
		subnetMap["subnet_type"] = subnet.SubnetType
		subnetMap["description"] = subnet.Description
		subnetMap["available_ip"] = subnet.AvailableIp
		subnetMap["tags"] = flattenTagsToMap(subnet.Tags)

		subnetsResult = append(subnetsResult, subnetMap)
	}
	addDebug(action, subnetsResult)

	d.Set("subnets", subnetsResult)

	d.SetId(resource.UniqueId())

	if outputFile != "" {
		if err := writeToFile(outputFile, subnetsResult); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_subnets", action, BCESDKGoERROR)
		}
	}

	return nil
}
