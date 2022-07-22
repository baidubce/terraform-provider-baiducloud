/*
Provide a resource to create a VPC.

Example Usage

```hcl
resource "baiducloud_vpc" "default" {
    name = "my-vpc"
    description = "baiducloud vpc created by terraform"
	cidr = "192.168.0.0/24"
}
```

Import

VPC instance can be imported, e.g.

```hcl
$ terraform import baiducloud_vpc.default vpc_id
```
*/
package baiducloud

import (
	"time"

	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func resourceBaiduCloudVpc() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudVpcCreate,
		Read:   resourceBaiduCloudVpcRead,
		Update: resourceBaiduCloudVpcUpdate,
		Delete: resourceBaiduCloudVpcDelete,

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
				Description: "Name of the VPC, which cannot take the value \"default\", the length is no more than 65 characters, and the value can be composed of numbers, characters and underscores.",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description of the VPC. The value is no more than 200 characters.",
				Optional:    true,
			},
			"cidr": {
				Type:        schema.TypeString,
				Description: "CIDR block for the VPC.",
				Required:    true,
				ForceNew:    true,
			},
			"route_table_id": {
				Type:        schema.TypeString,
				Description: "Route table ID created by default on VPC creation.",
				Computed:    true,
			},
			"secondary_cidrs": {
				Type:        schema.TypeList,
				Description: "Secondary cidr list of the VPC. They will not be repeated.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceBaiduCloudVpcCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	createVpcArgs := buildBaiduCloudVpcArgs(d, meta)
	action := "Create VPC " + createVpcArgs.Name

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (i interface{}, e error) {
			return vpcClient.CreateVPC(createVpcArgs)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, raw)
		result, _ := raw.(*vpc.CreateVPCResult)
		d.SetId(result.VPCID)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_vpc", action, BCESDKGoERROR)
	}

	return resourceBaiduCloudVpcRead(d, meta)
}

func resourceBaiduCloudVpcRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	vpcService := &VpcService{client}

	vpcId := d.Id()
	action := "Query VPC " + vpcId

	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (i interface{}, e error) {
		return vpcClient.GetVPCDetail(vpcId)
	})
	addDebug(action, raw)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_vpc", action, BCESDKGoERROR)
	}

	result, _ := raw.(*vpc.GetVPCDetailResult)
	d.Set("name", result.VPC.Name)
	d.Set("description", result.VPC.Description)
	d.Set("cidr", result.VPC.Cidr)
	d.Set("tags", flattenTagsToMap(result.VPC.Tags))
	d.Set("secondary_cidrs", result.VPC.SecondaryCidr)

	//computed attribute
	res, err := vpcService.GetRouteTableDetail("", vpcId)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_vpc", action, BCESDKGoERROR)
	}
	d.Set("route_table_id", res.RouteTableId)

	return nil
}

func resourceBaiduCloudVpcUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	vpcId := d.Id()
	action := "Update VPC " + vpcId
	update := false

	updateVpcArgs := &vpc.UpdateVPCArgs{}
	if d.HasChange("name") || d.HasChange("description") {
		update = true
		updateVpcArgs.Name = d.Get("name").(string)
		updateVpcArgs.Description = d.Get("description").(string)
	}

	if update {
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (i interface{}, e error) {
			return nil, vpcClient.UpdateVPC(vpcId, updateVpcArgs)
		})
		addDebug(action, updateVpcArgs)
		if err != nil {
			if NotFoundError(err) {
				d.SetId("")
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_vpc", action, BCESDKGoERROR)
		}
	}

	return resourceBaiduCloudVpcRead(d, meta)
}

func resourceBaiduCloudVpcDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	vpcId := d.Id()
	action := "Delete VPC " + vpcId

	clientToken := buildClientToken()
	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (i interface{}, e error) {
			return vpcId, vpcClient.DeleteVPC(vpcId, clientToken)
		})
		addDebug(action, raw)
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR, SubnetInuseError, NotAllowDeleteVpc}) {
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
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_vpc", action, BCESDKGoERROR)
	}

	return nil
}

func buildBaiduCloudVpcArgs(d *schema.ResourceData, meta interface{}) *vpc.CreateVPCArgs {
	request := &vpc.CreateVPCArgs{
		ClientToken: buildClientToken(),
	}

	if v := d.Get("name").(string); v != "" {
		request.Name = v
	}

	if v := d.Get("description").(string); v != "" {
		request.Description = v
	}

	if v := d.Get("cidr").(string); v != "" {
		request.Cidr = v
	}

	if v, ok := d.GetOk("tags"); ok {
		request.Tags = tranceTagMapToModel(v.(map[string]interface{}))
	}

	return request
}
