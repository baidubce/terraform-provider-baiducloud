/*
Use this data source to list instances of a cluster.

Example Usage

```hcl
data "baiducloud_ccev2_cluster_instances" "default" {
  cluster_id = baiducloud_ccev2_cluster.default_custom.id
  keyword_type = "instanceName"
  keyword = ""
  order_by = "instanceName"
  order = "ASC"
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
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudCCEv2ClusterInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudCCEv2ClusterInstancesRead,
		Schema: map[string]*schema.Schema{
			//Query Params
			"cluster_id": {
				Type:        schema.TypeString,
				Description: "CCEv2 Cluster ID",
				Required:    true,
				ForceNew:    true,
			},
			"keyword_type": {
				Type:         schema.TypeString,
				Description:  "Keyword type. Available Value: [instanceName, instanceID].",
				Optional:     true,
				ForceNew:     true,
				Default:      "instanceName",
				ValidateFunc: validation.StringInSlice(InstanceQueryKeywordTypePermitted, false),
			},
			"keyword": {
				Type:        schema.TypeString,
				Description: "The search keyword",
				Optional:    true,
				ForceNew:    true,
				Default:     "",
			},
			"order_by": {
				Type:         schema.TypeString,
				Description:  "The field that used to order the list. Available Value: [instanceName, instanceID, createdAt].",
				Optional:     true,
				ForceNew:     true,
				Default:      "createdAt",
				ValidateFunc: validation.StringInSlice(InstanceQueryOrderByPermitted, false),
			},
			"order": {
				Type:         schema.TypeString,
				Description:  "Ascendant or descendant order. Available Value: [ASC, DESC].",
				Optional:     true,
				ForceNew:     true,
				Default:      "ASC",
				ValidateFunc: validation.StringInSlice(QueryOrderPermitted, false),
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
			"master_list": {
				Type:        schema.TypeList,
				Description: "The search result",
				Computed:    true,
				Elem:        resourceCCEv2Instance(),
			},
			"nodes_list": {
				Type:        schema.TypeList,
				Description: "The search result",
				Computed:    true,
				Elem:        resourceCCEv2Instance(),
			},
		},
	}
}

func dataSourceBaiduCloudCCEv2ClusterInstancesRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.BaiduClient)

	args := &ccev2.ListInstancesByPageArgs{}
	if value, ok := d.GetOk("cluster_id"); ok && value.(string) != "" {
		args.ClusterID = value.(string)
	} else {
		err := errors.New("get cluster_id fail or cluster_id empty")
		log.Printf("Build ListInstancesByPageParams Error:" + err.Error())
		return WrapError(err)
	}
	listParams := &ccev2.ListInstancesByPageParams{}
	if value, ok := d.GetOk("keyword"); ok && value.(string) != "" {
		listParams.Keyword = value.(string)
	}
	if value, ok := d.GetOk("keyword_type"); ok && value.(string) != "" {
		listParams.KeywordType = ccev2.InstanceKeywordType(value.(string))
	}
	if value, ok := d.GetOk("order_by"); ok && value.(string) != "" {
		listParams.OrderBy = ccev2.InstanceOrderBy(value.(string))
	}
	if value, ok := d.GetOk("order"); ok && value.(string) != "" {
		listParams.Order = ccev2.Order(value.(string))
	}
	if value, ok := d.GetOk("page_size"); ok {
		listParams.PageSize = value.(int)
	}
	if value, ok := d.GetOk("page_no"); ok {
		listParams.PageNo = value.(int)
	}
	args.Params = listParams

	action := "Get CCEv2 Cluster Nodes " + args.ClusterID
	raw, err := client.WithCCEv2Client(func(client *ccev2.Client) (i interface{}, e error) {
		return client.ListInstancesByPage(args)
	})
	if err != nil {
		log.Printf("List Cluster Instances Error:" + err.Error())
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2_cluster_instances", action, BCESDKGoERROR)
	}
	addDebug(action, raw)

	response := raw.(*ccev2.ListInstancesResponse)
	if response.InstancePage == nil {
		err := errors.New("InstancePage is nil")
		log.Printf("List Cluster Instances Error:" + err.Error())
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2_cluster_instances", action, BCESDKGoERROR)
	}

	d.SetId(resource.UniqueId())

	nodes, err := convertInstanceFromJsonToMap(response.InstancePage.InstanceList, types.ClusterRoleNode)
	if err != nil {
		log.Printf("Get Cluster Follower Nodes Error:" + err.Error())
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2_cluster_instances", action, BCESDKGoERROR)
	}
	masters, err := convertInstanceFromJsonToMap(response.InstancePage.InstanceList, types.ClusterRoleMaster)
	if err != nil {
		log.Printf("Get Cluster Master Nodes Error:" + err.Error())
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2_cluster_instances", action, BCESDKGoERROR)
	}
	err = d.Set("nodes_list", nodes)
	if err != nil {
		log.Printf("Set 'nodes_list' to State Error:" + err.Error())
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2_cluster_instances", action, BCESDKGoERROR)
	}
	err = d.Set("master_list", masters)
	if err != nil {
		log.Printf("Set 'master_list' to State Error:" + err.Error())
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2_cluster_instances", action, BCESDKGoERROR)
	}

	err = d.Set("total_count", response.InstancePage.TotalCount)
	if err != nil {
		log.Printf("Set 'total_count' to State Error:" + err.Error())
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2_cluster_instances", action, BCESDKGoERROR)
	}

	return nil
}
