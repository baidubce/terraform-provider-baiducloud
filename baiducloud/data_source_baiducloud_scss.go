/*
Use this data source to query SCS list.

Example Usage

```hcl
data "baiducloud_scss" "default" {}

output "scss" {
 value = "${data.baiducloud_scss.default.scss}"
}
```
*/
package baiducloud

import (
	"regexp"

	"github.com/baidubce/bce-sdk-go/services/scs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudScss() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudScssRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Description:  "Regex pattern of the search name of scs instance",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateNameRegex,
			},
			"output_file": {
				Type:        schema.TypeString,
				Description: "Output file of the instances search result",
				Optional:    true,
				ForceNew:    true,
			},
			"filter": dataSourceFiltersSchema(),

			"scss": {
				Type:        schema.TypeList,
				Description: "The result of the instances list.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Description: "ID of the instance.",
							Computed:    true,
						},
						"instance_status": {
							Type:        schema.TypeString,
							Description: "Status of the instance.",
							Computed:    true,
						},
						"instance_name": {
							Type:        schema.TypeString,
							Description: "Name of the instance. Support for uppercase and lowercase letters, numbers, Chinese and special characters, such as \"-\",\"_\",\"/\",\".\", the value must start with a letter, length 1-65.",
							Computed:    true,
						},
						"node_type": {
							Type:        schema.TypeString,
							Description: "Type of the instance. Available values are cache.n1.micro, cache.n1.small, cache.n1.medium...cache.n1hs3.4xlarge.",
							Computed:    true,
						},
						"shard_num": {
							Type:        schema.TypeInt,
							Description: "The number of instance shard. IF cluster_type is cluster, support 2/4/6/8/12/16/24/32/48/64/96/128, if cluster_type is master_slave, support 1.",
							Computed:    true,
						},
						"proxy_num": {
							Type:        schema.TypeInt,
							Description: "The number of instance proxy.",
							Computed:    true,
						},
						"replication_num": {
							Type:        schema.TypeInt,
							Description: "The number of instance copies.",
							Computed:    true,
						},
						"cluster_type": {
							Type:        schema.TypeString,
							Description: "Type of the instance,  Available values are cluster, master_slave.",
							Computed:    true,
						},
						"engine": {
							Type:        schema.TypeString,
							Description: "Engine of the instance. Available values are redis, memcache.",
							Computed:    true,
						},
						"engine_version": {
							Type:        schema.TypeString,
							Description: "Engine version of the instance. Available values are 3.2, 4.0.",
							Computed:    true,
						},
						"v_net_ip": {
							Type:        schema.TypeString,
							Description: "ID of the specific vnet.",
							Computed:    true,
						},
						"domain": {
							Type:        schema.TypeString,
							Description: "Domain of the instance.",
							Computed:    true,
						},
						"port": {
							Type:        schema.TypeInt,
							Description: "The port used to access a instance.",
							Computed:    true,
						},
						"create_time": {
							Type:        schema.TypeString,
							Description: "Create time of the instance.",
							Computed:    true,
						},
						"expire_time": {
							Type:        schema.TypeString,
							Description: "Expire time of the instance.",
							Computed:    true,
						},
						"capacity": {
							Type:        schema.TypeInt,
							Description: "Memory capacity(GB) of the instance.",
							Computed:    true,
						},
						"used_capacity": {
							Type:        schema.TypeInt,
							Description: "Memory capacity(GB) of the instance to be used.",
							Computed:    true,
						},
						"payment_timing": {
							Type:        schema.TypeString,
							Description: "SCS payment timing",
							Computed:    true,
						},
						"zone_names": {
							Type:        schema.TypeList,
							Description: "Zone name list",
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"auto_renew": {
							Type:        schema.TypeBool,
							Description: "Whether to automatically renew.",
							Computed:    true,
						},
						"security_ips": {
							Type:        schema.TypeSet,
							Description: "Security ips of the scs.",
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"tags": tagsComputedSchema(),
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudScssRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	scsService := ScsService{client}

	action := "List all scs instances"
	listArgs := &scs.ListInstancesArgs{}
	scsList, err := scsService.ListAllInstances(listArgs)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_scss", action, BCESDKGoERROR)
	}

	var nameRegex string
	var specNameRegex *regexp.Regexp

	if value, ok := d.GetOk("name_regex"); ok {
		nameRegex = value.(string)
		if len(nameRegex) > 0 {
			specNameRegex = regexp.MustCompile(nameRegex)
		}
	}

	scsMap := make([]map[string]interface{}, 0, len(scsList))
	for _, e := range scsList {
		if len(nameRegex) > 0 && specNameRegex != nil {
			if !specNameRegex.MatchString(e.InstanceName) {
				continue
			}
		}
		ips, err := scsService.GetSecurityIPs(e.InstanceID)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_scss", action, BCESDKGoERROR)
		}
		scsMap = append(scsMap, map[string]interface{}{
			"instance_id":     e.InstanceID,
			"instance_name":   e.InstanceName,
			"instance_status": e.InstanceStatus,
			"cluster_type":    e.ClusterType,
			"engine":          e.Engine,
			"engine_version":  e.EngineVersion,
			"v_net_ip":        e.VnetIP,
			"domain":          e.Domain,
			"port":            e.Port,
			"create_time":     e.InstanceCreateTime,
			"capacity":        e.Capacity,
			"used_capacity":   e.UsedCapacity,
			"payment_timing":  e.PaymentTiming,
			"zone_names":      e.ZoneNames,
			"security_ips":    ips,
			"tags":            flattenTagsToMap(e.Tags),
		})
	}

	FilterDataSourceResult(d, &scsMap)

	addDebug("List filtered scs instances", scsMap)
	if err = d.Set("scss", scsMap); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_scss", action, BCESDKGoERROR)
	}
	d.SetId(resource.UniqueId())

	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := writeToFile(v.(string), scsMap); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_scss", action, BCESDKGoERROR)
		}
	}

	return nil
}
