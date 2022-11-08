/*
Use this resource to create a CFS mount target.

Example Usage

```hcl
resource "baiducloud_cfs_mount_target" "default" {
  fs_id = "cfs-xxxxxx"
  subnet_id = "sbn-xxxxxxx"
  vpc_id = "vpc-xxxxxxx"
}

Import

CFS mount target can be imported, e.g.

```hcl
$ terraform import baiducloud_cfs_mount_target.default id
```
*/
package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/cfs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
	"time"
)

func resourceBaiduCloudCfsMountTarget() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudCfsMountTargetCreate,
		Read:   resourceBaiduCloudCfsMountTargetRead,
		Delete: resourceBaiduCloudCfsMountTargetDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"fs_id": {
				Type:        schema.TypeString,
				Description: "CFS id which mount target belong to.",
				Required:    true,
				ForceNew:    true,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Description: "Subnet ID which mount target belong to.",
				Required:    true,
				ForceNew:    true,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Description: "VPC ID which mount target belong to.",
				Required:    true,
				ForceNew:    true,
			},
			"domain": {
				Type:        schema.TypeString,
				Description: "Domain of mount target.",
				Computed:    true,
			},
		},
	}
}

func resourceBaiduCloudCfsMountTargetCreate(d *schema.ResourceData, meta interface{}) error {
	action := "Create CFS mount target"
	client := meta.(*connectivity.BaiduClient)
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithCfsClient(func(client *cfs.Client) (i interface{}, e error) {
			return client.CreateMountTarget(buildCfsMountTargetArgs(d))
		})

		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		addDebug(action, raw)
		response, _ := raw.(*cfs.CreateMountTargetResult)
		d.SetId(response.MountID)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cfs_mount_target", action, BCESDKGoERROR)
	}
	return resourceBaiduCloudCfsMountTargetRead(d, meta)
}

func resourceBaiduCloudCfsMountTargetRead(d *schema.ResourceData, meta interface{}) error {
	action := "query cfs mount target"
	client := meta.(*connectivity.BaiduClient)
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithCfsClient(func(client *cfs.Client) (i interface{}, e error) {
			var fsId string
			if v, ok := d.GetOk("fs_id"); ok {
				fsId = v.(string)
			}
			return client.DescribeMountTarget(&cfs.DescribeMountTargetArgs{
				FSID:    fsId,
				MountID: d.Id(),
			})
		})

		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response, _ := raw.(*cfs.DescribeMountTargetResult)
		mountTarget := response.MountTargetList[0]
		if err := d.Set("domain", mountTarget.Domain); err != nil {
			return resource.NonRetryableError(err)
		}
		if err := d.Set("subnet_id", mountTarget.SubnetID); err != nil {
			return resource.NonRetryableError(err)
		}
		addDebug(action, raw)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cfs_mount_target", action, BCESDKGoERROR)
	}
	return nil
}
func resourceBaiduCloudCfsMountTargetDelete(d *schema.ResourceData, meta interface{}) error {
	action := "Delete cfs mount target"
	client := meta.(*connectivity.BaiduClient)
	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, err := client.WithCfsClient(func(client *cfs.Client) (i interface{}, e error) {
			var fsId string
			if v, ok := d.GetOk("fs_id"); ok {
				fsId = v.(string)
			}
			return nil, client.DropMountTarget(&cfs.DropMountTargetArgs{
				FSID:    fsId,
				MountId: d.Id(),
			})
		})

		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cfs_mount_target", action, BCESDKGoERROR)
	}
	return nil
}

func buildCfsMountTargetArgs(d *schema.ResourceData) *cfs.CreateMountTargetArgs {
	res := &cfs.CreateMountTargetArgs{}
	if v, ok := d.GetOk("fs_id"); ok {
		res.FSID = v.(string)
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		res.VpcID = v.(string)
	}
	if v, ok := d.GetOk("subnet_id"); ok {
		res.SubnetId = v.(string)
	}
	return res
}
