/*
Use this data source to list instances of a instancegroup.

Example Usage

```hcl
data "baiducloud_ccev2_instance_group_instances" "default" {
  cluster_id = baiducloud_ccev2_cluster.default_custom.id
  instance_group_id = baiducloud_ccev2_instance_group.ccev2_instance_group_1.id
  page_no = 0
  page_size = 0
}
```
*/
package baiducloud

import (
	"errors"
	"log"

	ccev2 "github.com/baidubce/bce-sdk-go/services/cce/v2"
	"github.com/baidubce/bce-sdk-go/services/cce/v2/types"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudCCEv2InstanceGroupInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudCCEv2InstanceGroupInstancesRead,
		Schema: map[string]*schema.Schema{
			//Query Params
			"cluster_id": {
				Type:        schema.TypeString,
				Description: "CCEv2 Cluster ID",
				Required:    true,
				ForceNew:    true,
			},
			"instance_group_id": {
				Type:        schema.TypeString,
				Description: "CCEv2 instance group ID. instanceName/instanceID",
				Required:    true,
				ForceNew:    true,
			},
			"page_no": {
				Type:        schema.TypeInt,
				Description: "Page number of query result",
				Optional:    true,
				ForceNew:    true,
				Default:     0,
			},
			"page_size": {
				Type:        schema.TypeInt,
				Description: "The size of every page",
				Optional:    true,
				ForceNew:    true,
				Default:     0,
			},
			//Query Result
			"total_count": {
				Type:        schema.TypeInt,
				Description: "The total count of the result",
				Computed:    true,
			},
			"instance_list": {
				Type:        schema.TypeList,
				Description: "The search result",
				Computed:    true,
				Elem:        resourceCCEv2Instance(),
			},
		},
	}
}

func dataSourceBaiduCloudCCEv2InstanceGroupInstancesRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.BaiduClient)
	args := &ccev2.ListInstanceByInstanceGroupIDArgs{}

	if value, ok := d.GetOk("cluster_id"); ok && value.(string) != "" {
		args.ClusterID = value.(string)
	} else {
		err := errors.New("get cluster_id fail or cluster_id empty")
		log.Printf("Build ListInstanceByInstanceGroupIDArgs Error:" + err.Error())
		return WrapError(err)
	}
	if value, ok := d.GetOk("instance_group_id"); ok && value.(string) != "" {
		args.InstanceGroupID = value.(string)
	} else {
		err := errors.New("get instance_group_id fail or instance_group_id empty")
		log.Printf("Build ListInstanceByInstanceGroupIDArgs Error:" + err.Error())
		return WrapError(err)
	}
	if value, ok := d.GetOk("page_size"); ok {
		args.PageSize = value.(int)
	}
	if value, ok := d.GetOk("page_no"); ok {
		args.PageNo = value.(int)
	}

	action := "Get CCEv2 InstanceGroup Nodes Cluster ID:" + args.ClusterID + " InstanceGroup ID:" + args.InstanceGroupID
	raw, err := client.WithCCEv2Client(func(client *ccev2.Client) (i interface{}, e error) {
		return client.ListInstancesByInstanceGroupID(args)
	})
	if err != nil {
		log.Printf("List InstanceGroup Instances Error:" + err.Error())
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2_instance_group_instances", action, BCESDKGoERROR)
	}
	addDebug(action, raw)

	response := raw.(*ccev2.ListInstancesByInstanceGroupIDResponse)
	if response.Page.List == nil {
		err := errors.New("instance list is nil")
		log.Printf("List InstanceGroup Instances Error:" + err.Error())
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2_instance_group_instances", action, BCESDKGoERROR)
	}

	nodes, err := convertInstanceFromJsonToMap(response.Page.List, types.ClusterRoleNode)
	if err != nil {
		log.Printf("Get Instance Group Follower Nodes Fail" + err.Error())
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2_instance_group_instances", action, BCESDKGoERROR)
	}
	masters, err := convertInstanceFromJsonToMap(response.Page.List, types.ClusterRoleMaster)
	if err != nil {
		log.Printf("Get Instance Group Master Nodes Fail" + err.Error())
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2_instance_group_instances", action, BCESDKGoERROR)
	}
	total := make([]interface{}, len(masters)+len(nodes))
	copy(total, masters)
	copy(total[len(masters):], nodes)

	err = d.Set("instance_list", total)
	if err != nil {
		log.Printf("Set 'instance_list' to State Error:" + err.Error())
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2_instance_group_instances", action, BCESDKGoERROR)
	}

	err = d.Set("total_count", response.Page.TotalCount)
	if err != nil {
		log.Printf("Set 'total_count' to State Error:" + err.Error())
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2_instance_group_instances", action, BCESDKGoERROR)
	}

	d.SetId(resource.UniqueId())

	return nil
}
