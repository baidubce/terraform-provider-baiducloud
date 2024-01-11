/*
Use this data source to query EIP group list.

Example Usage

```hcl
data "baiducloud_eipgroups" "default" {}

output "eip_groups" {
 value = "${data.baiducloud_eipgroups.default.eip_groups}"
}
```
*/
package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/eip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudEipgroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudEipgroupsRead,

		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:        schema.TypeString,
				Description: "Id of Eip group",
				Optional:    true,
				ForceNew:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "name of Eip group",
				Optional:    true,
				ForceNew:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Eip group status",
				Optional:    true,
				ForceNew:    true,
			},
			"output_file": {
				Type:        schema.TypeString,
				Description: "Eipgroups search result output file",
				Optional:    true,
				ForceNew:    true,
			},
			"filter": dataSourceFiltersSchema(),

			"eip_groups": {
				Type:        schema.TypeList,
				Description: "Eip group list",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "Eip group name",
							Computed:    true,
						},
						"status": {
							Type:        schema.TypeString,
							Description: "Eip group status",
							Computed:    true,
						},
						"id": {
							Type:        schema.TypeString,
							Description: "Eip group id",
							Computed:    true,
						},
						"band_width_in_mbps": {
							Type:        schema.TypeInt,
							Description: "Eip group band width in mbps",
							Computed:    true,
						},
						"default_domestic_bandwidth": {
							Type:        schema.TypeInt,
							Description: "Eip group default domestic bandwidth",
							Computed:    true,
						},
						"bw_short_id": {
							Type:        schema.TypeString,
							Description: "Eip group bw short id",
							Computed:    true,
						},
						"bw_bandwidth_in_mbps": {
							Type:        schema.TypeInt,
							Description: "Eip group bw bandwidth in mbps",
							Computed:    true,
						},
						"domestic_bw_short_id": {
							Type:        schema.TypeString,
							Description: "Eip group domestic bw short id",
							Computed:    true,
						},
						"domestic_bw_bandwidth_in_mbps": {
							Type:        schema.TypeInt,
							Description: "Eip group domestic bw bandwidth in mbps",
							Computed:    true,
						},
						"payment_timing": {
							Type:        schema.TypeString,
							Description: "Eip group payment timing",
							Computed:    true,
						},
						"billing_method": {
							Type:        schema.TypeString,
							Description: "Eip group billing method",
							Computed:    true,
						},
						"create_time": {
							Type:        schema.TypeString,
							Description: "Eip group create time",
							Computed:    true,
						},
						"expire_time": {
							Type:        schema.TypeString,
							Description: "Eip group expire time",
							Computed:    true,
						},
						"region": {
							Type:        schema.TypeString,
							Description: "Eip group region",
							Computed:    true,
						},
						"route_type": {
							Type:        schema.TypeString,
							Description: "Eip group route type",
							Computed:    true,
						},
						"tags": tagsComputedSchema(),
						"eips": {
							Type:        schema.TypeList,
							Description: "Eip list",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"eip": {
										Type:        schema.TypeString,
										Description: "Eip address",
										Computed:    true,
									},
									"name": {
										Type:        schema.TypeString,
										Description: "Eip name",
										Computed:    true,
									},
									"bandwidth_in_mbps": {
										Type:        schema.TypeInt,
										Description: "Eip bandwidth(Mbps)",
										Computed:    true,
									},
									"status": {
										Type:        schema.TypeString,
										Description: "Eip status",
										Computed:    true,
									},
									"eip_instance_type": {
										Type:        schema.TypeString,
										Description: "Eip instance type",
										Computed:    true,
									},
									"share_group_id": {
										Type:        schema.TypeString,
										Description: "Eip share group id",
										Computed:    true,
									},
									"payment_timing": {
										Type:        schema.TypeString,
										Description: "Eip payment timing",
										Computed:    true,
									},
									"billing_method": {
										Type:        schema.TypeString,
										Description: "Eip billing method",
										Computed:    true,
									},
									"create_time": {
										Type:        schema.TypeString,
										Description: "Eip create time",
										Computed:    true,
									},
									"expire_time": {
										Type:        schema.TypeString,
										Description: "Eip expire time",
										Computed:    true,
									},
									"tags": tagsComputedSchema(),
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudEipgroupsRead(d *schema.ResourceData, meta interface{}) error {

	groupId := d.Get("group_id").(string)

	name := d.Get("name").(string)

	action := "List all eip group id is " + groupId + " name is " + name

	eipGroupArgs, err := buildBaiduCloudEipGroupListArgs(d, meta)

	if err != nil {
		return WrapError(err)
	}

	eipGroups, err := listAllEipGroups(eipGroupArgs, meta)

	addDebug(action, err)

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_eipgroups", action, BCESDKGoERROR)
	}

	eipGroupsResult := make([]map[string]interface{}, 0)

	for _, eipGroup := range eipGroups {

		innerMap := make(map[string]interface{})
		innerMap["name"] = eipGroup.Name
		innerMap["status"] = eipGroup.Status
		innerMap["id"] = eipGroup.Id
		innerMap["band_width_in_mbps"] = eipGroup.BandWidthInMbps
		innerMap["default_domestic_bandwidth"] = eipGroup.DefaultDomesticBandwidth
		innerMap["bw_short_id"] = eipGroup.BwShortId
		innerMap["bw_bandwidth_in_mbps"] = eipGroup.BwBandwidthInMbps
		innerMap["domestic_bw_short_id"] = eipGroup.DomesticBwShortId
		innerMap["domestic_bw_bandwidth_in_mbps"] = eipGroup.DomesticBwBandwidthInMbps
		innerMap["payment_timing"] = eipGroup.PaymentTiming
		innerMap["billing_method"] = eipGroup.BillingMethod
		innerMap["create_time"] = eipGroup.CreateTime
		innerMap["expire_time"] = eipGroup.ExpireTime
		innerMap["region"] = eipGroup.Region
		innerMap["route_type"] = eipGroup.RouteType

		innerEips := make([]map[string]interface{}, 0, len(eipGroup.Eips))

		for _, e := range eipGroup.Eips {
			innerEips = append(innerEips, map[string]interface{}{
				"eip":               e.Eip,
				"name":              e.Name,
				"status":            e.Status,
				"eip_instance_type": e.EipInstanceType,
				"share_group_id":    e.ShareGroupId,
				"bandwidth_in_mbps": e.BandWidthInMbps,
				"payment_timing":    e.PaymentTiming,
				"billing_method":    e.BillingMethod,
				"create_time":       e.CreateTime,
				"expire_time":       e.ExpireTime,
				"tags":              flattenTagsToMap(e.Tags),
			})
		}

		innerMap["eips"] = innerEips

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_eipgroups", action, BCESDKGoERROR)
		}

		eipGroupsResult = append(eipGroupsResult, innerMap)
	}

	addDebug(action, eipGroupsResult)

	FilterDataSourceResult(d, &eipGroupsResult)

	if err := d.Set("eip_groups", eipGroupsResult); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_eipgroups", action, BCESDKGoERROR)
	}

	d.SetId(resource.UniqueId())

	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := writeToFile(v.(string), eipGroupsResult); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_eipgroups", action, BCESDKGoERROR)
		}
	}
	return nil
}

func listAllEipGroups(args *eip.ListEipGroupArgs, meta interface{}) ([]eip.EipGroupModel, error) {
	client := meta.(*connectivity.BaiduClient)
	action := "List all eipgroups "

	eipGroups := make([]eip.EipGroupModel, 0)
	for {
		raw, err := client.WithEipClient(func(eipClient *eip.Client) (interface{}, error) {
			return eipClient.ListEipGroup(args)
		})

		if err != nil {
			return nil, WrapErrorf(err, DefaultErrorMsg, "baiducloud_eipgroups", action, BCESDKGoERROR)
		}

		result, _ := raw.(*eip.ListEipGroupResult)
		eipGroups = append(eipGroups, result.EipGroup...)
		if !result.IsTruncated {
			break
		}
		args.Marker = result.NextMarker
		args.MaxKeys = result.MaxKeys
	}

	return eipGroups, nil
}

func buildBaiduCloudEipGroupListArgs(d *schema.ResourceData, meta interface{}) (*eip.ListEipGroupArgs, error) {
	request := &eip.ListEipGroupArgs{}

	if v := d.Get("name").(string); v != "" {
		request.Name = v
	}

	if v := d.Get("group_id").(string); v != "" {
		request.Id = v
	}

	if v := d.Get("status").(string); v != "" {
		request.Status = v
	}

	return request, nil

}
