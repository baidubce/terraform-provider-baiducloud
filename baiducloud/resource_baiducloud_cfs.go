/*
Use this resource to create a CFS.

Example Usage

```hcl
resource "baiducloud_cfs" "default" {
  name = "terraform_test"
  zone = "zoneD"
}

Import

CFS can be imported, e.g.

```hcl
$ terraform import baiducloud_cfs.default cfs_id
```
*/
package baiducloud

import (
	"fmt"
	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/cfs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
	"time"
)

func resourceBaiduCloudCfs() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudCfsCreate,
		Read:   resourceBaiduCloudCfsRead,
		Update: resourceBaiduCloudCfsUpdate,
		Delete: resourceBaiduCloudCfsDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "cfs name, length must be between 1 and 64 bytes",
				Required:    true,
			},
			"zone": {
				Type:        schema.TypeString,
				Description: "cfs zone",
				Optional:    true,
				ForceNew:    true,
			},
			"protocol": {
				Type:         schema.TypeString,
				Description:  "CFS protocol, available value is nfs and smb, default is nfs",
				Default:      "nfs",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"nfs", "smb"}, false),
			},
			"type": {
				Type:         schema.TypeString,
				Description:  "CFS type, default is cap",
				Default:      "cap",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"cap"}, false),
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Description: "VPC ID",
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "CFS status, available value is available,updating,paused and unavailable",
				Computed:    true,
			},
		},
	}
}

func resourceBaiduCloudCfsCreate(d *schema.ResourceData, meta interface{}) error {
	action := "create cfs"
	client := meta.(*connectivity.BaiduClient)
	cfsService := CfsService{
		Client: client,
	}
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithCfsClient(func(client *cfs.Client) (i interface{}, e error) {
			return client.CreateFS(buildCreateCfsArgs(d))
		})

		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		addDebug(action, raw)
		response, _ := raw.(*cfs.CreateFSResult)
		d.SetId(response.FSID)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cfs", action, BCESDKGoERROR)
	}
	stateConf := buildStateConf(
		[]string{string(cfs.FSStatusUnavailable), string(cfs.FSStatusPaused)},
		[]string{string(cfs.FSStatusAvailable)},
		d.Timeout(schema.TimeoutCreate),
		cfsService.cfsStateRefresh(d.Id()),
	)
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cfs", action, BCESDKGoERROR)
	}
	return resourceBaiduCloudCfsRead(d, meta)
}

func resourceBaiduCloudCfsUpdate(d *schema.ResourceData, meta interface{}) error {
	action := "Update CFS " + d.Id()
	client := meta.(*connectivity.BaiduClient)
	err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		raw, err := client.WithCfsClient(func(client *cfs.Client) (i interface{}, e error) {
			return nil, client.UpdateFS(&cfs.UpdateFSArgs{
				FSID:   d.Id(),
				FSName: d.Get("name").(string),
			})
		})

		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		addDebug(action, raw)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cfs", action, BCESDKGoERROR)
	}
	return resourceBaiduCloudCfsRead(d, meta)
}

func resourceBaiduCloudCfsRead(d *schema.ResourceData, meta interface{}) error {
	action := "Query CFS detail " + d.Id()
	cfsService := CfsService{
		Client: meta.(*connectivity.BaiduClient),
	}
	cfsDetail, err := cfsService.GetCfsDetail(d.Id())
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cfs", action, BCESDKGoERROR)
	}
	if err := d.Set("name", cfsDetail.Name); err != nil {
		return fmt.Errorf("error setting name: %w", err)
	}
	if err := d.Set("type", cfsDetail.Type); err != nil {
		return fmt.Errorf("error setting type: %w", err)
	}
	if err := d.Set("protocol", cfsDetail.Protocol); err != nil {
		return fmt.Errorf("error setting protocol: %w", err)
	}
	if err := d.Set("vpc_id", cfsDetail.VpcID); err != nil {
		return fmt.Errorf("error setting vpc_id: %w", err)
	}
	if err := d.Set("status", cfsDetail.Status); err != nil {
		return fmt.Errorf("error setting status: %w", err)
	}
	return nil
}
func resourceBaiduCloudCfsDelete(d *schema.ResourceData, meta interface{}) error {
	action := "Delete CFS " + d.Id()
	client := meta.(*connectivity.BaiduClient)
	cfsService := CfsService{
		Client: meta.(*connectivity.BaiduClient),
	}

	// wait for all mount targets were deleted fully
	stateConf := buildStateConf(
		[]string{CfsMountTargetDeleting},
		[]string{CfsMountTargetDeleted},
		d.Timeout(schema.TimeoutDelete),
		cfsService.cfsMountTargetCountRefresh(d.Id()),
	)
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cfs", action, BCESDKGoERROR)
	}

	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		raw, err := client.WithCfsClient(func(client *cfs.Client) (i interface{}, e error) {
			return nil, client.DropFS(&cfs.DropFSArgs{
				FSID: d.Id(),
			})
		})

		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		addDebug(action, raw)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cfs", action, BCESDKGoERROR)
	}
	return nil
}

func buildCreateCfsArgs(d *schema.ResourceData) *cfs.CreateFSArgs {
	res := &cfs.CreateFSArgs{
		ClientToken: buildClientToken(),
	}
	if v, ok := d.GetOk("name"); ok {
		res.Name = v.(string)
	}
	if v, ok := d.GetOk("zone"); ok {
		res.Zone = v.(string)
	}
	if v, ok := d.GetOk("protocol"); ok {
		res.Protocol = v.(string)
	}
	if v, ok := d.GetOk("type"); ok {
		res.Type = v.(string)
	}
	return res
}
