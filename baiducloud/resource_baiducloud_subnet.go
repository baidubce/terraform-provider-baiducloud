/*
Provide a resource to create a VPC subnet.

Example Usage

```hcl
resource "baiducloud_subnet" "default" {
  name = "my-subnet"
  zone_name = "cn-bj-a"
  cidr = "192.168.3.0/24"
  vpc_id = "${baiducloud_vpc.default.id}"
}

resource "baiducloud_vpc" "default" {
  name = "my-vpc"
  cidr = "192.168.0.0/16"
}
```

Import

VPC subnet instance can be imported, e.g.

```hcl
$ terraform import baiducloud_subnet.default subnet_id
```
*/
package baiducloud

import (
	"fmt"
	"net"
	"time"

	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func resourceBaiduCloudSubnet() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudSubnetCreate,
		Read:   resourceBaiduCloudSubnetRead,
		Update: resourceBaiduCloudSubnetUpdate,
		Delete: resourceBaiduCloudSubnetDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the subnet, which cannot take the value \"default\", the length is no more than 65 characters, and the value can be composed of numbers, characters and underscores.",
				Required:    true,
			},
			"zone_name": {
				Type:        schema.TypeString,
				Description: "The availability zone name within which the subnet should be created.",
				Required:    true,
				ForceNew:    true,
			},
			"cidr": {
				Type:        schema.TypeString,
				Description: "CIDR block of the subnet.",
				Required:    true,
				ForceNew:    true,
			},
			"ipv6_cidr": {
				Type:        schema.TypeString,
				Description: "CIDR block of the subnet.",
				Computed:    true,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Description: "The VPC ID that the subnet belongs to.",
				Required:    true,
				ForceNew:    true,
			},
			"subnet_type": {
				Type:         schema.TypeString,
				Description:  "Type of the subnet, valid values are BCC, BCC_NAT and BBC. Default to BCC.",
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validateSubnetType(),
			},
			"enable_ipv6": {
				Type:        schema.TypeBool,
				Description: "Whether to enable the IPv6 subnet. 'true' indicates enabled, by default it is 'false' which means not enabled.",
				Optional:    true,
				Computed:    true,
			},
			"vpc_secondary_cidr": {
				Type:        schema.TypeString,
				Description: "The CIDR of the secondary subnet belonging to the VPC.",
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description of the subnet, and the value must be no more than 200 characters.",
				Optional:    true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceBaiduCloudSubnetCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	createSubnetArgs := buildBaiduCloudSubnetArgs(d, meta)
	action := "Create Subnet " + createSubnetArgs.Name

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (i interface{}, e error) {
			return vpcClient.CreateSubnet(createSubnetArgs)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, raw)
		result, _ := raw.(*vpc.CreateSubnetResult)
		d.SetId(result.SubnetId)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_subnet", action, BCESDKGoERROR)
	}

	return resourceBaiduCloudSubnetRead(d, meta)
}

func resourceBaiduCloudSubnetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	subnetId := d.Id()
	action := "Query Subnet " + subnetId

	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (i interface{}, e error) {
		return vpcClient.GetSubnetDetail(subnetId)
	})
	addDebug(action, raw)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_subnet", action, BCESDKGoERROR)
	}

	result, _ := raw.(*vpc.GetSubnetDetailResult)
	d.Set("name", result.Subnet.Name)
	d.Set("zone_name", result.Subnet.ZoneName)
	d.Set("cidr", result.Subnet.Cidr)
	d.Set("ipv6_cidr", result.Subnet.Ipv6Cidr)
	d.Set("enable_ipv6", result.Subnet.Ipv6Cidr != "")
	d.Set("vpc_id", result.Subnet.VPCId)
	d.Set("subnet_type", result.Subnet.SubnetType)
	d.Set("description", result.Subnet.Description)
	d.Set("tags", flattenTagsToMap(result.Subnet.Tags))

	raw, err = client.WithVpcClient(func(vpcClient *vpc.Client) (i interface{}, e error) {
		return vpcClient.GetVPCDetail(result.Subnet.VPCId)
	})
	addDebug(action, raw)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_subnet", action, BCESDKGoERROR)
	}
	vpcResult, _ := raw.(*vpc.GetVPCDetailResult)
	cidr, err := findParentCidr(append(vpcResult.VPC.SecondaryCidr, vpcResult.VPC.Cidr), result.Subnet.Cidr)
	if err != nil {
		return WrapError(err)
	}
	d.Set("vpc_secondary_cidr", cidr)
	return nil
}

func resourceBaiduCloudSubnetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	subnetId := d.Id()
	action := "Update Subnet " + subnetId

	if d.HasChange("name") || d.HasChange("description") || d.HasChange("enable_ipv6") {
		updateSubnetArgs := &vpc.UpdateSubnetArgs{
			Name:        d.Get("name").(string),
			EnableIpv6:  d.Get("enable_ipv6").(bool),
			Description: d.Get("description").(string),
		}

		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (i interface{}, e error) {
			return nil, vpcClient.UpdateSubnet(subnetId, updateSubnetArgs)
		})
		addDebug(action, updateSubnetArgs)
		if err != nil {
			if NotFoundError(err) {
				d.SetId("")
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_subnet", action, BCESDKGoERROR)
		}
	}

	return resourceBaiduCloudSubnetRead(d, meta)
}

func resourceBaiduCloudSubnetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	subnetId := d.Id()
	action := "Delete Subnet " + subnetId

	clientToken := buildClientToken()
	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (i interface{}, e error) {
			return nil, vpcClient.DeleteSubnet(subnetId, clientToken)
		})
		addDebug(action, nil)
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR, SubnetInuseError, ResourceNeedRelease, NotAllowOperateSubnet}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_subnet", action, BCESDKGoERROR)
	}

	return nil
}

func buildBaiduCloudSubnetArgs(d *schema.ResourceData, meta interface{}) *vpc.CreateSubnetArgs {
	request := &vpc.CreateSubnetArgs{
		ClientToken: buildClientToken(),
	}

	if v := d.Get("name").(string); v != "" {
		request.Name = v
	}
	if v := d.Get("zone_name").(string); v != "" {
		request.ZoneName = v
	}
	if v := d.Get("cidr").(string); v != "" {
		request.Cidr = v
	}
	if v := d.Get("vpc_id").(string); v != "" {
		request.VpcId = v
	}
	if v, ok := d.Get("enable_ipv6").(bool); ok {
		request.EnableIpv6 = v
	}
	if v := d.Get("vpc_secondary_cidr").(string); v != "" {
		request.VpcSecondaryCidr = v
	}
	if v := d.Get("subnet_type").(string); v != "" {
		request.SubnetType = vpc.SubnetType(v)
	}
	if v := d.Get("description").(string); v != "" {
		request.Description = v
	}
	if v, ok := d.GetOk("tags"); ok {
		request.Tags = tranceTagMapToModel(v.(map[string]interface{}))
	}

	return request
}

func findParentCidr(cidrs []string, subnetCidr string) (string, error) {
	// 解析子网的 *net.IPNet
	_, subnet, err := net.ParseCIDR(subnetCidr)
	if err != nil {
		return "", fmt.Errorf("can not resolve subnet CIDR: %v", err)
	}

	// 遍历所有 CIDR，查找匹配的父网段
	for _, cidr := range cidrs {
		_, parentNet, err := net.ParseCIDR(cidr)
		if err != nil {
			return "", fmt.Errorf("can not resolve subnet CIDR '%s': %v", cidr, err)
		}

		// 检查子网的 IP 是否在父网段范围内
		if parentNet.Contains(subnet.IP) {
			// 检查子网掩码长度是否大于等于父网段掩码长度
			subnetOnes, _ := subnet.Mask.Size()
			parentOnes, _ := parentNet.Mask.Size()
			if subnetOnes >= parentOnes {
				return cidr, nil
			}
		}
	}

	return "", fmt.Errorf("can not find vpc cidr")
}
