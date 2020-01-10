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
	"regexp"
	"strings"

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
			"instance_type": {
				Type:        schema.TypeString,
				Description: "Instance type of the search spec",
				Optional:    true,
				ForceNew:    true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Description:  "Regex pattern of the search spec name",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateNameRegex,
			},
			"cpu_count": {
				Type:        schema.TypeInt,
				Description: "Useful cpu count of the search spec",
				Optional:    true,
				ForceNew:    true,
			},
			"memory_size_in_gb": {
				Type:        schema.TypeInt,
				Description: "Useful memory size in GB of the search spec",
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

	var specName, instanceType string
	var cpuCount, memorySizeInGb int
	var specNameRegex *regexp.Regexp

	if value, ok := d.GetOk("spec_name"); ok {
		specName = value.(string)
		if len(specName) > 0 {
			specNameRegex = regexp.MustCompile(specName)
		}
	}

	if value, ok := d.GetOk("instance_type"); ok {
		instanceType = strings.TrimSpace(value.(string))
	}

	if value, ok := d.GetOk("cpu_count"); ok {
		cpuCount = value.(int)
	}

	if value, ok := d.GetOk("memory_size_in_gb"); ok {
		memorySizeInGb = value.(int)
	}

	response := raw.(*api.ListSpecResult)
	specMap := make([]map[string]interface{}, 0, len(response.InstanceTypes))
	for _, spec := range response.InstanceTypes {
		if len(specName) > 0 && specNameRegex != nil {
			if !specNameRegex.MatchString(spec.Name) {
				continue
			}
		}

		if len(instanceType) > 0 && spec.Type != instanceType {
			continue
		}

		if cpuCount > 0 && spec.CpuCount != cpuCount {
			continue
		}

		if memorySizeInGb > 0 && spec.MemorySizeInGB != memorySizeInGb {
			continue
		}

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
