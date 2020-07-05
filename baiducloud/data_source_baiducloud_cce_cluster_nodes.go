/*
Use this data source to get cce cluster nodes.

Example Usage

```hcl
data "baiducloud_cce_cluster_nodes" "default" {
   cluster_uuid	 = "c-NqYwWEhu"
}

output "nodes" {
 value = "${data.baiducloud_cce_cluster_nodes.default.nodes}"
}
```
*/
package baiducloud

import (
	"regexp"

	"github.com/baidubce/bce-sdk-go/services/cce"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudCCEClusterNodes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudCCEClusterNodesRead,

		Schema: map[string]*schema.Schema{
			"cluster_uuid": {
				Type:        schema.TypeString,
				Description: "UUID of the cce cluster.",
				Required:    true,
				ForceNew:    true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Description: "ID of the search instance.",
				Optional:    true,
				ForceNew:    true,
			},
			"instance_name_regex": {
				Type:         schema.TypeString,
				Description:  "Regex pattern of the search spec name.",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateNameRegex,
			},
			"instance_type": {
				Type:        schema.TypeString,
				Description: "Type of the search instance.",
				Optional:    true,
				ForceNew:    true,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Description: "ID of the subnet.",
				Optional:    true,
				ForceNew:    true,
			},
			"available_zone": {
				Type:        schema.TypeString,
				Description: "Available zone of the cluster node.",
				Optional:    true,
				ForceNew:    true,
			},
			"nodes": {
				Type:        schema.TypeList,
				Description: "Result of the cluster nodes list.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Description: "ID of the instance.",
							Computed:    true,
						},
						"instance_name": {
							Type:        schema.TypeString,
							Description: "Name of the instance.",
							Computed:    true,
						},
						"instance_uuid": {
							Type:        schema.TypeString,
							Description: "UUID of the instance.",
							Computed:    true,
						},
						"available_zone": {
							Type:        schema.TypeString,
							Description: "CDS disk size, should in [1, 32765], when snapshot_id not set, this parameter is required.",
							Computed:    true,
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Description: "VPC id of the instance.",
							Computed:    true,
						},
						"vpc_cidr": {
							Type:        schema.TypeString,
							Description: "VPC cidr of the instance.",
							Computed:    true,
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Description: "Subnet id of the instance.",
							Computed:    true,
						},
						"subnet_type": {
							Type:        schema.TypeString,
							Description: "Subnet type of the instance.",
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
						"delete_time": {
							Type:        schema.TypeString,
							Description: "Delete time of the instance.",
							Computed:    true,
						},
						"status": {
							Type:        schema.TypeString,
							Description: "Status of the instance.",
							Computed:    true,
						},
						"eip": {
							Type:        schema.TypeString,
							Description: "Eip of the instance.",
							Computed:    true,
						},
						"eip_bandwidth": {
							Type:        schema.TypeInt,
							Description: "Eip bandwidth(Mbps) of the instance.",
							Computed:    true,
						},
						"cpu": {
							Type:        schema.TypeInt,
							Description: "Number of cpu cores.",
							Computed:    true,
						},
						"memory": {
							Type:        schema.TypeInt,
							Description: "Memory capacity(GB) of the instance.",
							Computed:    true,
						},
						"disk_size": {
							Type:        schema.TypeInt,
							Description: "Local disk size of the node.",
							Computed:    true,
						},
						"sys_disk": {
							Type:        schema.TypeInt,
							Description: "System disk size of the node.",
							Computed:    true,
						},
						"instance_type": {
							Type:        schema.TypeString,
							Description: "Type of the instance.",
							Computed:    true,
						},
						"blb": {
							Type:        schema.TypeString,
							Description: "BLB address of the node.",
							Computed:    true,
						},
						"floating_ip": {
							Type:        schema.TypeString,
							Description: "Floating ip of the node.",
							Computed:    true,
						},
						"fix_ip": {
							Type:        schema.TypeString,
							Description: "Fix ip of the node, which is assigned in VPC.",
							Computed:    true,
						},
						"payment_method": {
							Type:        schema.TypeString,
							Description: "Payment method of the node.",
							Computed:    true,
						},
						"runtime_version": {
							Type:        schema.TypeString,
							Description: "Version of the instance runtime.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudCCEClusterNodesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	clusterUuid := d.Get("cluster_uuid").(string)
	instanceId := ""
	if value, ok := d.GetOk("instance_id"); ok {
		instanceId = value.(string)
	}

	args := &cce.ListNodeArgs{
		ClusterUuid: clusterUuid,
	}

	action := "Get CCE Cluster " + clusterUuid
	raw, err := client.WithCCEClient(func(client *cce.Client) (i interface{}, e error) {
		return client.ListNodes(args)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cce_cluster_nodes", action, BCESDKGoERROR)
	}
	addDebug(action, raw)

	response := raw.(*cce.ListNodeResult)
	nodeListWithInstanceId := make([]cce.Node, 0, len(response.Nodes))
	if instanceId == "" {
		nodeListWithInstanceId = append(nodeListWithInstanceId, response.Nodes...)
	} else {
		for _, node := range response.Nodes {
			if node.InstanceShortId == instanceId {
				nodeListWithInstanceId = append(nodeListWithInstanceId, node)
			}
		}
	}

	instanceType := ""
	if value, ok := d.GetOk("instance_type"); ok {
		instanceType = value.(string)
	}
	nodeListWithInstanceType := make([]cce.Node, 0, len(nodeListWithInstanceId))
	if instanceType == "" {
		nodeListWithInstanceType = append(nodeListWithInstanceType, nodeListWithInstanceId...)
	} else {
		for _, node := range nodeListWithInstanceId {
			if node.InstanceType == instanceType {
				nodeListWithInstanceId = append(nodeListWithInstanceId, node)
			}
		}
	}

	subnetId := ""
	if value, ok := d.GetOk("subnet_id"); ok {
		subnetId = value.(string)
	}
	nodeListWithSubnetId := make([]cce.Node, 0, len(nodeListWithInstanceType))
	if subnetId == "" {
		nodeListWithSubnetId = append(nodeListWithSubnetId, nodeListWithInstanceType...)
	} else {
		for _, node := range nodeListWithInstanceType {
			if node.InstanceType == instanceType {
				nodeListWithSubnetId = append(nodeListWithSubnetId, node)
			}
		}
	}

	availableZone := ""
	if value, ok := d.GetOk("available_zone"); ok {
		availableZone = value.(string)
	}
	nodeListWithZone := make([]cce.Node, 0, len(nodeListWithSubnetId))
	if availableZone == "" {
		nodeListWithZone = append(nodeListWithZone, nodeListWithSubnetId...)
	} else {
		for _, node := range nodeListWithSubnetId {
			if node.InstanceType == instanceType {
				nodeListWithZone = append(nodeListWithZone, node)
			}
		}
	}

	var instanceNameRegexStr string
	var instanceNameRegex *regexp.Regexp
	if value, ok := d.GetOk("instance_name_regex"); ok {
		instanceNameRegexStr = value.(string)
		if len(instanceNameRegexStr) > 0 {
			instanceNameRegex = regexp.MustCompile(instanceNameRegexStr)
		}
	}
	resultNodeList := make([]cce.Node, 0, len(nodeListWithZone))
	if len(instanceNameRegexStr) > 0 && instanceNameRegex != nil {
		for _, node := range nodeListWithZone {
			if !instanceNameRegex.MatchString(node.InstanceName) {
				continue
			}
		}
	} else {
		resultNodeList = append(resultNodeList, nodeListWithZone...)
	}

	nodesMap := make([]map[string]interface{}, 0, len(resultNodeList))
	for _, node := range resultNodeList {
		nodesMap = append(nodesMap, map[string]interface{}{
			"instance_id":     node.InstanceShortId,
			"instance_name":   node.InstanceName,
			"instance_uuid":   node.InstanceUuid,
			"available_zone":  node.AvailableZone,
			"vpc_id":          node.VpcId,
			"vpc_cidr":        node.VpcCidr,
			"subnet_id":       node.SubnetId,
			"subnet_type":     node.SubnetType,
			"eip":             node.Eip,
			"eip_bandwidth":   node.EipBandwidth,
			"cpu":             node.Cpu,
			"memory":          node.Memory,
			"disk_size":       node.DiskSize,
			"sys_disk":        node.SysDisk,
			"instance_type":   node.InstanceType,
			"blb":             node.Blb,
			"floating_ip":     node.FloatingIp,
			"fix_ip":          node.FixIp,
			"create_time":     node.CreateTime.String(),
			"delete_time":     node.DeleteTime.String(),
			"status":          node.Status,
			"expire_time":     node.ExpireTime.String(),
			"payment_method":  node.PaymentMethod,
			"runtime_version": node.RuntimeVersion,
		})
	}

	if err := d.Set("nodes", nodesMap); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cce_cluster_nodes", action, BCESDKGoERROR)
	}

	d.SetId(clusterUuid + instanceType + instanceId + subnetId)

	return nil
}
