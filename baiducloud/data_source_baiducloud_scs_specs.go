/*
Use this data source to query scs specs list.

Example Usage

```hcl
data "data.baiducloud_scs_specs" "default" {}

output "specs" {
  value = "${data.baiducloud_scs_specs.default.specs}"
}
```
*/
package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/scs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudScsSpecs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudScsSpecsRead,

		Schema: map[string]*schema.Schema{
			"cluster_type": {
				Type:         schema.TypeString,
				Description:  "Type of the instance,  Available values are cluster, master_slave.",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"cluster", "master_slave"}, false),
			},
			"node_capacity": {
				Type:        schema.TypeInt,
				Description: "Memory capacity(GB) of the instance node.",
				Required:    true,
				ForceNew:    true,
			},
			"output_file": {
				Type:        schema.TypeString,
				Description: "Output file for saving result.",
				Optional:    true,
				ForceNew:    true,
			},
			"filter": dataSourceFiltersSchema(),

			// Attributes used for result
			"specs": {
				Type:        schema.TypeList,
				Description: "Useful spec list, when create a scs instance, suggest use node_type/cpu_num/instance_flavor/allowed_nodeNum_list as scs instance parameters",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_capacity": {
							Type:        schema.TypeInt,
							Description: "Memory capacity(GB) of the instance node.",
							Computed:    true,
						},
						"node_type": {
							Type:        schema.TypeString,
							Description: "Useful node type",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudScsSpecsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	action := "Query all SCS Specs"
	raw, err := client.WithScsClient(func(client *scs.Client) (i interface{}, e error) {
		return client.GetNodeTypeList()
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_scs_specs", action, BCESDKGoERROR)
	}
	addDebug(action, raw)

	var clusterType string
	if value, ok := d.GetOk("cluster_type"); ok {
		clusterType = value.(string)
	}

	var nodeTypeList []scs.NodeType
	getNodeTypeListResult := raw.(*scs.GetNodeTypeListResult)
	if len(clusterType) > 0 && clusterType == "cluster" {
		nodeTypeList = getNodeTypeListResult.ClusterNodeTypeList
	} else if clusterType == "master_slave" {
		nodeTypeList = getNodeTypeListResult.DefaultNodeTypeList
	}

	var nodeCapacity int
	if value, ok := d.GetOk("node_capacity"); ok {
		nodeCapacity = value.(int)
	}

	specMapOrigin := make([]map[string]interface{}, 0, len(nodeTypeList))
	for _, spec := range nodeTypeList {

		if nodeCapacity > 0 && spec.InstanceFlavor != nodeCapacity {
			continue
		}

		specMapOrigin = append(specMapOrigin, map[string]interface{}{
			"node_capacity": spec.InstanceFlavor,
			"node_type":     spec.NodeType,
		})
	}

	//FilterDataSourceResult(d, &specMap)
	specMap := make([]map[string]interface{}, 0, len(nodeTypeList))
	filter := NewDataSourceFilter(d)
	for _, data := range specMapOrigin {
		if filter.checkFilter(data) {
			specMap = append(specMap, data)
		}
	}

	if err := d.Set("specs", specMap); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_scs_specs", action, BCESDKGoERROR)
	}
	d.SetId(resource.UniqueId())

	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := writeToFile(v.(string), specMap); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_scs_specs", action, BCESDKGoERROR)
		}
	}

	return nil
}
