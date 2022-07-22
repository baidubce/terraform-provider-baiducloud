/*
Use this data source to get cce container network.

Example Usage

```hcl
data "baiducloud_cce_container_net" "default" {
	vpc_id   = "vpc-t6d16myuuqyu"
	vpc_cidr = "192.168.0.0/20"
}

output "net" {
  value = "${data.baiducloud_cce_container_net.default.container_net}"
}
```
*/
package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/cce"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudCceContainerNet() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudCceContainerNetRead,

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Description: "CCE used vpc id",
				Required:    true,
				ForceNew:    true,
			},
			"vpc_cidr": {
				Type:        schema.TypeString,
				Description: "CCE used vpc cidr",
				Required:    true,
				ForceNew:    true,
			},
			"size": {
				Type:        schema.TypeInt,
				Description: "CCE used max container count",
				Optional:    true,
				ForceNew:    true,
			},

			// Attributes used for result
			"container_net": {
				Type:        schema.TypeString,
				Description: "container net",
				Computed:    true,
			},
			"capacity": {
				Type:        schema.TypeInt,
				Description: "container net capacity",
				Computed:    true,
			},
		},
	}
}

func dataSourceBaiduCloudCceContainerNetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	vpcId := ""
	if v, ok := d.GetOk("vpc_id"); ok {
		vpcId = v.(string)
	}

	vpcCidr := ""
	if v, ok := d.GetOk("vpc_cidr"); ok {
		vpcCidr = v.(string)
	}

	size := 0
	if v, ok := d.GetOk("size"); ok {
		size = v.(int)
	}

	action := "Get container net in vpc " + vpcId + "[" + vpcCidr + "]"
	args := &cce.GetContainerNetArgs{
		VpcShortId: vpcId,
		VpcCidr:    vpcCidr,
	}
	if size > 0 {
		args.Size = size
	}
	raw, err := client.WithCCEClient(func(client *cce.Client) (i interface{}, e error) {
		return client.GetContainerNet(args)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cce_container_net", action, BCESDKGoERROR)
	}
	addDebug(action, raw)

	response := raw.(*cce.GetContainerNetResult)
	if err := d.Set("container_net", response.ContainerNet); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cce_container_net", action, BCESDKGoERROR)
	}
	if err := d.Set("capacity", response.Capacity); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cce_container_net", action, BCESDKGoERROR)
	}
	d.SetId(resource.UniqueId())

	return nil
}
