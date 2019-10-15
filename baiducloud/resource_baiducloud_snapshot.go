/*
Provide a resource to create a Snapshot.

Example Usage

```hcl
resource "baiducloud_snapshot" "my-snapshot" {
  name        = "${var.name}"
  description = "${var.description}"
  volume_id   = "v-Trb3rQXa"
}
```

Import

Snapshot can be imported, e.g.

```hcl
$ terraform import baiducloud_snapshot.default id
```
*/
package baiducloud

import (
	"time"

	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/bcc"
	"github.com/baidubce/bce-sdk-go/services/bcc/api"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func resourceBaiduCloudSnapshot() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudSnapshotCreate,
		Read:   resourceBaiduCloudSnapshotRead,
		Delete: resourceBaiduCloudSnapshotDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the snapshot, which supports uppercase and lowercase letters, numbers, Chinese and special characters, such as \"-\",\"_\",\"/\",\",\".\", and the value must start with a letter, length 1-65.",
				Required:    true,
				ForceNew:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description of the snapshot.",
				Optional:    true,
				ForceNew:    true,
			},
			"volume_id": {
				Type:        schema.TypeString,
				Description: "Volume id of the snapshot, this value will be nil if volume has been released.",
				ForceNew:    true,
				Required:    true,
			},
			"create_time": {
				Type:        schema.TypeString,
				Description: "Creation time of the snapshot.",
				Computed:    true,
			},
			"size_in_gb": {
				Type:        schema.TypeInt,
				Description: "Size of the snapshot in GB.",
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Status of the snapshot.",
				Computed:    true,
			},
			"create_method": {
				Type:        schema.TypeString,
				Description: "Creation method of the snapshot.",
				Computed:    true,
			},
		},
	}

}

func resourceBaiduCloudSnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	bccService := BccService{client}

	args := &api.CreateSnapshotArgs{
		ClientToken: buildClientToken(),
	}
	if v, ok := d.GetOk("volume_id"); ok && v.(string) != "" {
		args.VolumeId = v.(string)
	}

	if v, ok := d.GetOk("name"); ok && v.(string) != "" {
		args.SnapshotName = v.(string)
	}

	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		args.Description = v.(string)
	}

	// check volume status is available or inuse
	stateConf := buildStateConf(
		CDSProcessingStatus,
		[]string{string(api.VolumeStatusAVAILABLE), string(api.VolumeStatusINUSE)},
		d.Timeout(schema.TimeoutCreate),
		bccService.CDSVolumeStateRefreshFunc(args.VolumeId, CDSFailedStatus))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapError(err)
	}

	action := "Create snapshot " + args.SnapshotName

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithBccClient(func(client *bcc.Client) (i interface{}, e error) {
			return client.CreateSnapshot(args)
		})

		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, raw)

		response := raw.(*api.CreateSnapshotResult)
		d.SetId(response.SnapshotId)

		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_snapshot", action, BCESDKGoERROR)
	}

	stateConf = buildStateConf(
		[]string{"Creating"},
		[]string{"Available"},
		d.Timeout(schema.TimeoutCreate),
		bccService.SnapshotStateRefreshFunc(d.Id(), []string{"CreatedFailed", "NotAvailable"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapError(err)
	}

	return resourceBaiduCloudSnapshotRead(d, meta)
}

func resourceBaiduCloudSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	id := d.Id()
	action := "Query Snapshot " + id

	raw, err := client.WithBccClient(func(client *bcc.Client) (i interface{}, e error) {
		return client.GetSnapshotDetail(id)
	})

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_snapshot", action, BCESDKGoERROR)
	}

	snapshot := raw.(*api.GetSnapshotDetailResult).Snapshot
	d.Set("size_in_gb", snapshot.SizeInGB)
	d.Set("create_time", snapshot.CreateTime)
	d.Set("create_method", snapshot.CreateMethod)
	d.Set("status", snapshot.Status)
	d.Set("volume_id", snapshot.VolumeId)
	d.Set("name", snapshot.Name)
	d.Set("description", snapshot.Description)

	return nil
}

func resourceBaiduCloudSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	id := d.Id()
	action := "Delete Snapshot " + id

	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, err := client.WithBccClient(func(client *bcc.Client) (i interface{}, e error) {
			return nil, client.DeleteSnapshot(id)
		})
		addDebug(action, id)

		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
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
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_snapshot", action, BCESDKGoERROR)
	}

	return nil
}
