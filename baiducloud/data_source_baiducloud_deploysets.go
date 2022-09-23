/*
Use this data source to query deploy set list.

Example Usage

```hcl
data "baiducloud_deploysets" "default" {}

output "deploysets" {
 value = "${data.baiducloud_deploysets.default.deploysets}"
}
```
*/
package baiducloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudDeploySets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudDeploySetRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"output_file": {
				Type:        schema.TypeString,
				Description: "deployset search result output file",
				Optional:    true,
				ForceNew:    true,
			},
			"deploy_sets": {
				Type:        schema.TypeList,
				Description: "Image list",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "Name of the deployset.",
							Computed:    true,
						},
						"strategy": {
							Type:        schema.TypeString,
							Description: "Strategy of deployset.Available values are HOST_HA, RACK_HA and TOR_HA",
							Computed:    true,
						},
						"concurrency": {
							Type:        schema.TypeInt,
							Description: "concurrency of deployset.",
							Computed:    true,
						},
						"deployset_id": {
							Type:        schema.TypeString,
							Description: "Id of deployset.",
							Computed:    true,
						},
						"desc": {
							Type:        schema.TypeString,
							Description: "Description of the deployset.",
							Computed:    true,
						},

						"az_intstance_statis_list": {
							Type:        schema.TypeList,
							Description: "Availability Zone Instance Statistics List.",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_count": {
										Type:        schema.TypeInt,
										Description: "Count of instance which is in the deployset.",
										Computed:    true,
									},
									"bcc_instance_cnt": {
										Type:        schema.TypeInt,
										Description: "Count of BCC instance which is in the deployset.",
										Computed:    true,
									},
									"bbc_instance_cnt": {
										Type:        schema.TypeInt,
										Description: "Count of BBC instance which is in the deployset.",
										Computed:    true,
									},
									"instance_total": {
										Type:        schema.TypeInt,
										Description: "Total of instance which is in the deployset.",
										Computed:    true,
									},
									"zone_name": {
										Type:        schema.TypeString,
										Description: "Zone name of deployset.",
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudDeploySetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	bccService := BccService{client: client}
	action := "Query deploy set List."

	deploySetList, err := bccService.ListAllDeploySets()
	deploySetsMaps := make([]map[string]interface{}, 0, len(deploySetList))
	for _, deploySet := range deploySetList {
		intstanceStatisMap := make([]map[string]interface{}, 0, len(deploySet.InstanceList))
		for _, ins := range deploySet.InstanceList {
			intstanceStatisMap = append(intstanceStatisMap, map[string]interface{}{
				"bcc_instance_cnt": ins.BccCount,
				"bbc_instance_cnt": ins.BbcCount,
				"instance_count":   ins.Count,
				"instance_total":   ins.Total,
				"zone_name":        ins.ZoneName,
			})
		}
		deploySetsMaps = append(deploySetsMaps, map[string]interface{}{
			"name":                     deploySet.Name,
			"desc":                     deploySet.Desc,
			"strategy":                 deploySet.Strategy,
			"concurrency":              deploySet.Concurrency,
			"deployset_id":             deploySet.DeploySetId,
			"az_intstance_statis_list": intstanceStatisMap,
		})
	}
	FilterDataSourceResult(d, &deploySetsMaps)
	d.Set("deploy_sets", deploySetsMaps)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_deployset", action, BCESDKGoERROR)
	}
	d.SetId(resource.UniqueId())

	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := writeToFile(v.(string), deploySetsMaps); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_deployset", action, BCESDKGoERROR)
		}
	}
	return nil
}
