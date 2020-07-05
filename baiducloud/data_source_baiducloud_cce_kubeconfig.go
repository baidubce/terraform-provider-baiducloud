/*
Use this data source to get cce kubeconfig.

Example Usage

```hcl
data "baiducloud_cce_kubeconfig" "default" {
	cluster_uuid = "c-NqYwWEhu"
}

output "kubeconfig" {
  value = "${data.baiducloud_cce_kubeconfig.default.data}"
}
```
*/
package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/cce"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudCceKubeConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudCceKubeConfigRead,

		Schema: map[string]*schema.Schema{
			"cluster_uuid": {
				Type:        schema.TypeString,
				Description: "UUID of the cce cluster.",
				Required:    true,
				ForceNew:    true,
			},
			"config_type": {
				Type:        schema.TypeString,
				Description: "Config type of the cce cluster.",
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
			"data": {
				Type:        schema.TypeString,
				Description: "Data of the cce kubeconfig.",
				Computed:    true,
			},
		},
	}
}

func dataSourceBaiduCloudCceKubeConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	clusterUuid := d.Get("cluster_uuid").(string)

	action := "Get CCE Cluster " + clusterUuid + " kubeConfig"
	args := &cce.GetKubeConfigArgs{ClusterUuid: clusterUuid}

	if v, ok := d.GetOk("config_type"); ok {
		args.Type = cce.KubeConfigType(v.(string))
	}

	raw, err := client.WithCCEClient(func(client *cce.Client) (i interface{}, e error) {
		return client.GetKubeConfig(args)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cce_kubeconfig", action, BCESDKGoERROR)
	}
	addDebug(action, raw)

	data := raw.(*cce.GetKubeConfigResult).Data

	if err := d.Set("data", data); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cce_kubeconfig", action, BCESDKGoERROR)
	}
	d.SetId(resource.UniqueId())

	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := writeStringToFile(v.(string), data); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cce_kubeconfig", action, BCESDKGoERROR)
		}
	}

	return nil
}
