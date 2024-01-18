/*
Use this data source to query EIP bp list.

Example Usage

```hcl
data "baiducloud_eipbps" "default" {}

output "eip_bps" {
 value = "${data.baiducloud_eipbps.default.eip_bps}"
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

func dataSourceBaiduCloudEipbps() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudEipbpsRead,

		Schema: map[string]*schema.Schema{
			"bp_id": {
				Type:        schema.TypeString,
				Description: "Id of Eip bp",
				Optional:    true,
				ForceNew:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "name of Eip bp",
				Optional:    true,
				ForceNew:    true,
			},
			"bind_type": {
				Type:        schema.TypeString,
				Description: "Eip bp bind type",
				Optional:    true,
				ForceNew:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "Eip bp type",
				Optional:    true,
				ForceNew:    true,
			},
			"output_file": {
				Type:        schema.TypeString,
				Description: "Eipbps search result output file",
				Optional:    true,
				ForceNew:    true,
			},
			"filter": dataSourceFiltersSchema(),

			"eip_bps": {
				Type:        schema.TypeList,
				Description: "Eip bp list",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "Eip bp name",
							Computed:    true,
						},
						"id": {
							Type:        schema.TypeString,
							Description: "Eip bp id",
							Computed:    true,
						},
						"band_width_in_mbps": {
							Type:        schema.TypeInt,
							Description: "Eip bp band width in mbps",
							Computed:    true,
						},
						"instance_id": {
							Type:        schema.TypeString,
							Description: "Eip bp instance id",
							Computed:    true,
						},
						"eips": {
							Type:        schema.TypeSet,
							Description: "Eip bp eips",
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"create_time": {
							Type:        schema.TypeString,
							Description: "Eip bp create time",
							Computed:    true,
						},
						"auto_release_time": {
							Type:        schema.TypeString,
							Description: "Eip bp auto release time",
							Computed:    true,
						},
						"type": {
							Type:        schema.TypeString,
							Description: "Eip bp type",
							Computed:    true,
						},
						"region": {
							Type:        schema.TypeString,
							Description: "Eip bp region",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudEipbpsRead(d *schema.ResourceData, meta interface{}) error {

	bpId := d.Get("bp_id").(string)

	name := d.Get("name").(string)

	action := "List all eip bp id is " + bpId + " name is " + name

	eipbpArgs, err := buildBaiduCloudEipbpListArgs(d, meta)

	if err != nil {
		return WrapError(err)
	}

	eipbps, err := listAllEipbps(eipbpArgs, meta)

	addDebug(action, eipbps)

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_eipbps", action, BCESDKGoERROR)
	}

	eipbpsResult := make([]map[string]interface{}, 0)

	for _, eipbp := range eipbps {

		innerMap := make(map[string]interface{})
		innerMap["name"] = eipbp.Name
		innerMap["id"] = eipbp.Id
		innerMap["band_width_in_mbps"] = eipbp.BandwidthInMbps
		innerMap["instance_id"] = eipbp.InstanceId
		innerMap["eips"] = eipbp.Eips
		innerMap["create_time"] = eipbp.CreateTime
		innerMap["auto_release_time"] = eipbp.AutoReleaseTime
		innerMap["type"] = eipbp.Type
		innerMap["region"] = eipbp.Region

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_eipbps", action, BCESDKGoERROR)
		}

		eipbpsResult = append(eipbpsResult, innerMap)
	}

	addDebug(action, eipbpsResult)

	FilterDataSourceResult(d, &eipbpsResult)

	if err := d.Set("eip_bps", eipbpsResult); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_eipbps", action, BCESDKGoERROR)
	}

	d.SetId(resource.UniqueId())

	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := writeToFile(v.(string), eipbpsResult); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_eipbps", action, BCESDKGoERROR)
		}
	}
	return nil
}

func listAllEipbps(args *eip.ListEipBpArgs, meta interface{}) ([]eip.EipBpList, error) {
	client := meta.(*connectivity.BaiduClient)
	action := "List all eipbps "

	eipbps := make([]eip.EipBpList, 0)
	for {
		raw, err := client.WithEipClient(func(eipClient *eip.Client) (interface{}, error) {
			return eipClient.ListEipBp(args)
		})

		if err != nil {
			return nil, WrapErrorf(err, DefaultErrorMsg, "baiducloud_eipbps", action, BCESDKGoERROR)
		}

		addDebug(action, raw)

		result, _ := raw.(*eip.ListEipBpResult)

		eipbps = append(eipbps, result.EipGroup...)

		if !result.IsTruncated {
			break
		}
		args.Marker = result.NextMarker
		args.MaxKeys = result.MaxKeys
	}

	return eipbps, nil
}

func buildBaiduCloudEipbpListArgs(d *schema.ResourceData, meta interface{}) (*eip.ListEipBpArgs, error) {
	request := &eip.ListEipBpArgs{
		MaxKeys: 1000,
	}

	if v := d.Get("name").(string); v != "" {
		request.Name = v
	}

	if v := d.Get("bp_id").(string); v != "" {
		request.Id = v
	}

	if v := d.Get("bind_type").(string); v != "" {
		request.BindType = v
	}

	if v := d.Get("type").(string); v != "" {
		request.Type = v
	}

	return request, nil

}
