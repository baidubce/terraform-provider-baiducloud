/*
Use this data source to recommend ccev2 cluster IP CIDR.

Example Usage

```hcl
data "baiducloud_ccev2_clusterip_cidr" "default" {
  vpc_cidr = var.vpc_cidr
  container_cidr = var.container_cidr
  cluster_max_service_num = 32
  private_net_cidrs = ["172.16.0.0/12",]
  ip_version = "ipv4"
  output_file = "${path.cwd}/recommendClusterIPCidr.txt"
}
```
*/
package baiducloud

import (
	"encoding/json"
	"fmt"

	ccev2 "github.com/baidubce/bce-sdk-go/services/cce/v2"
	"github.com/baidubce/bce-sdk-go/services/cce/v2/types"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudCCEv2ClusterIPCidrs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudCCEv2ClusterIPCidrsRead,

		Schema: map[string]*schema.Schema{
			//以下是资源的请求参数
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
			"container_cidr": {
				Type:        schema.TypeString,
				Description: "Container CIDR",
				Optional:    true,
			},
			"container_cidr_ipv6": {
				Type:        schema.TypeString,
				Description: "Container CIDR IPv6",
				Optional:    true,
			},
			"cluster_max_service_num": {
				Type:        schema.TypeInt,
				Description: "Max service number in the cluster",
				Optional:    true,
			},
			"private_net_cidrs": {
				Type:        schema.TypeList,
				Description: "Private Net CIDRs",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"private_net_cidrs_ipv6": {
				Type:        schema.TypeList,
				Description: "Private Net CIDRs IPv6",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ip_version": {
				Type:        schema.TypeString,
				Description: "IP Version",
				Optional:    true,
			},
			//输出结果写到哪
			"output_file": {
				Type:        schema.TypeString,
				Description: "Result output file",
				Optional:    true,
			},
			//以下是资源的返回结果
			"is_success": {
				Type:        schema.TypeBool,
				Description: "Is the recommendation request success",
				Computed:    true,
			},
			"err_msg": {
				Type:        schema.TypeString,
				Description: "Error message if an error occurs",
				Computed:    true,
			},
			"request_id": {
				Type:        schema.TypeString,
				Description: "Request ID",
				Computed:    true,
			},
			"recommended_clusterip_cidrs": {
				Type:        schema.TypeList,
				Description: "Recommend Cluster IP CIDR List",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"recommended_clusterip_cidrs_ipv6": {
				Type:        schema.TypeList,
				Description: "Recommend Cluster IP CIDR List IPv6",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceBaiduCloudCCEv2ClusterIPCidrsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	args := &ccev2.RecommendClusterIPCIDRArgs{}

	//构建请求参数
	if v := d.Get("vpc_cidr"); v.(string) != "" {
		args.VPCCIDR = v.(string)
	}
	if v := d.Get("vpc_cidr_ipv6"); v.(string) != "" {
		args.VPCCIDRIPv6 = v.(string)
	}
	if v := d.Get("container_cidr"); v.(string) != "" {
		args.ContainerCIDR = v.(string)
	}
	if v := d.Get("container_cidr_ipv6"); v.(string) != "" {
		args.ContainerCIDRIPv6 = v.(string)
	}
	if v := d.Get("cluster_max_service_num"); v != nil {
		args.ClusterMaxServiceNum = v.(int)
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
	if v := d.Get("ip_version"); v != nil {
		args.IPVersion = types.ContainerNetworkIPType(v.(string))
	}

	action := "Recommend CCEv2 Cluster IP CIDR in vpc"
	raw, err := client.WithCCEv2Client(func(client *ccev2.Client) (i interface{}, e error) {
		return client.RecommendClusterIPCIDR(args)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2_clusterip_cidr", action, BCESDKGoERROR)
	}
	addDebug(action, raw)

	response := raw.(*ccev2.RecommendClusterIPCIDRResponse)

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
	if err := d.Set("recommended_clusterip_cidrs", response.RecommendedClusterIPCIDRs); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2_container_cidr", action, BCESDKGoERROR)
	}
	if err := d.Set("recommended_clusterip_cidrs_ipv6", response.RecommendedClusterIPCIDRIPv6s); err != nil {
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
