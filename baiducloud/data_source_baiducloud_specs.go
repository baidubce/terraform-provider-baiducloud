/*
Use this data source to query spec list.

Example Usage

```hcl
data "baiducloud_specs" "default" {}

output "spec" {
  value = "${data.baiducloud_specs.default.specs}"
}
```
*/
package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/bcc"
	"github.com/baidubce/bce-sdk-go/services/bcc/api"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudSpecs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudSpecsRead,

		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:        schema.TypeString,
				Description: "Output file for saving result.",
				Optional:    true,
				ForceNew:    true,
			},

			// Attributes used for result
			"specs": {
				Type:        schema.TypeList,
				Description: "Useful spec list, when create a bcc instance, suggest use instance_type/cpu_count/memory_capacity_in_gb as bcc instance parameters",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "Spec name",
							Computed:    true,
						},
						"instance_type": {
							Type:        schema.TypeString,
							Description: "Useful instance type",
							Computed:    true,
						},
						"cpu_count": {
							Type:        schema.TypeInt,
							Description: "Useful cpu count",
							Computed:    true,
						},

						"memory_size_in_gb": {
							Type:        schema.TypeInt,
							Description: "Useful memory size in GB",
							Computed:    true,
						},
						"local_disk_size_in_gb": {
							Type:        schema.TypeInt,
							Description: "Useful local disk size in GB",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudSpecsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	action := "Query all Specs"
	raw, err := client.WithBccClient(func(client *bcc.Client) (i interface{}, e error) {
		return client.ListSpec()
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_specs", action, BCESDKGoERROR)
	}
	addDebug(action, raw)

	response := raw.(*api.ListSpecResult)
	specMap := make([]map[string]interface{}, 0, len(response.InstanceTypes))
	for _, spec := range response.InstanceTypes {
		specMap = append(specMap, map[string]interface{}{
			"name":                  spec.Name,
			"instance_type":         spec.Type,
			"cpu_count":             spec.CpuCount,
			"memory_size_in_gb":     spec.MemorySizeInGB,
			"local_disk_size_in_gb": spec.LocalDiskSizeInGB,
		})
	}

	if err := d.Set("specs", specMap); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_specs", action, BCESDKGoERROR)
	}
	d.SetId(resource.UniqueId())

	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := writeToFile(v.(string), specMap); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_specs", action, BCESDKGoERROR)
		}
	}

	return nil
}
