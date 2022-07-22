/*
Use this data source to recommend ccev2 container CIDR.

Example Usage

```hcl
data "baiducloud_ccev2_container_cidr" "default" {
  vpc_id = var.vpc_id
  vpc_cidr = var.vpc_cidr
  cluster_max_node_num = 16
  max_pods_per_node = 32
  private_net_cidrs = ["172.16.0.0/12",]
  k8s_version = "1.16.8"
  ip_version = "ipv4"
  output_file = "${path.cwd}/recommendContainerCidr.txt"
}
```
*/
package baiducloud

import (
	"encoding/json"
	"fmt"
	ccev2 "github.com/baidubce/bce-sdk-go/services/cce/v2"
	"github.com/baidubce/bce-sdk-go/services/cce/v2/types"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudCCEv2ContainerCIDRs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudCCEv2ContainerCIDRsRead,

		Schema: map[string]*schema.Schema{
			//以下是资源的请求参数
			"vpc_id": {
				Type:        schema.TypeString,
				Description: "VPC ID",
				Optional:    true,
			},
			"vpc_cidr": {
				Type:        schema.TypeString,
				Description: "VPC CIDR",
				Optional:    true,
			},
			"vpc_cidr_ipv6": {
				Type:        schema.TypeString,
				Description: "VPC CIDR IPv6",
				Optional:    true,
			},
			"cluster_max_node_num": {
				Type:        schema.TypeInt,
				Description: "Max node number in a cluster",
				Optional:    true,
			},
			"max_pods_per_node": {
				Type:        schema.TypeInt,
				Description: "Max pod number in a node",
				Optional:    true,
			},
			"private_net_cidrs": {
				Type:        schema.TypeList,
				Description: "Private Net CIDR List",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"private_net_cidrs_ipv6": {
				Type:        schema.TypeList,
				Description: "Private Net CIDR List IPv6",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"k8s_version": {
				Type:        schema.TypeString,
				Description: "K8s Version",
				Optional:    true,
			},
			"ip_version": {
				Type:        schema.TypeString,
				Description: "IP version",
				Optional:    true,
			},
			//输出结果写到哪
			"output_file": {
				Type:        schema.TypeString,
				Description: "Eips search result output file",
				Optional:    true,
			},
			//以下是资源的返回结果
			"is_success": {
				Type:        schema.TypeBool,
				Description: "Whether the recommendation success",
				Computed:    true,
			},
			"err_msg": {
				Type:        schema.TypeString,
				Description: "Error message if an error occures",
				Computed:    true,
			},
			"request_id": {
				Type:        schema.TypeString,
				Description: "Request ID",
				Computed:    true,
			},
			"recommended_container_cidrs": {
				Type:        schema.TypeList,
				Description: "Recomment Container CIDR",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"recommended_container_cidrs_ipv6": {
				Type:        schema.TypeList,
				Description: "Recomment Container CIDRs IPv6",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceBaiduCloudCCEv2ContainerCIDRsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	args := &ccev2.RecommendContainerCIDRArgs{}

	//构建请求参数
	if v := d.Get("vpc_id"); v.(string) != "" {
		args.VPCID = v.(string)
	}
	if v := d.Get("vpc_cidr"); v.(string) != "" {
		args.VPCCIDR = v.(string)
	}
	if v := d.Get("vpc_cidr_ipv6"); v.(string) != "" {
		args.VPCCIDRIPv6 = v.(string)
	}
	if v := d.Get("cluster_max_node_num"); v != nil {
		args.ClusterMaxNodeNum = v.(int)
	}
	if v := d.Get("max_pods_per_node"); v != nil {
		args.MaxPodsPerNode = v.(int)
	}
	if v := d.Get("private_net_cidrs"); v != nil {
		cidrs := make([]ccev2.PrivateNetString, 0)
		for _, cidrRaw := range v.([]interface{}) {
			cidrs = append(cidrs, ccev2.PrivateNetString(cidrRaw.(string)))
		}
		args.PrivateNetCIDRs = cidrs
	}
	if v := d.Get("private_net_cidrs_ipv6"); v != nil {
		cidrsipv6 := make([]ccev2.PrivateNetString, 0)
		for _, cidrRaw := range v.([]interface{}) {
			cidrsipv6 = append(cidrsipv6, ccev2.PrivateNetString(cidrRaw.(string)))
		}
		args.PrivateNetCIDRIPv6s = cidrsipv6
	}
	if v := d.Get("k8s_version"); v != nil {
		args.K8SVersion = types.K8SVersion(v.(string))
	}
	if v := d.Get("ip_version"); v != nil {
		args.IPVersion = types.ContainerNetworkIPType(v.(string))
	}

	action := "Recommend CCEv2 Container CIDR in vpc"
	raw, err := client.WithCCEv2Client(func(client *ccev2.Client) (i interface{}, e error) {
		return client.RecommendContainerCIDR(args)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2_container_cidr", action, BCESDKGoERROR)
	}
	addDebug(action, raw)

	response := raw.(*ccev2.RecommendContainerCIDRResponse)

	//设置返回结果
	d.SetId(resource.UniqueId())
	if err := d.Set("is_success", response.IsSuccess); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2_container_cidr", action, BCESDKGoERROR)
	}
	if err := d.Set("err_msg", response.ErrMsg); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2_container_cidr", action, BCESDKGoERROR)
	}
	if err := d.Set("request_id", response.RequestID); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2_container_cidr", action, BCESDKGoERROR)
	}
	if err := d.Set("recommended_container_cidrs", response.RecommendedContainerCIDRs); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2_container_cidr", action, BCESDKGoERROR)
	}
	if err := d.Set("recommended_container_cidrs_ipv6", response.RecommendedContainerCIDRIPv6s); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2_container_cidr", action, BCESDKGoERROR)
	}

	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		s, _ := json.MarshalIndent(response, "", "\t")
		str := fmt.Sprintf("%s", s)
		if err := writeStringToFile(v.(string), str); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2_container_cidr", action, BCESDKGoERROR)
		}
	}

	return nil
}
