/*
Provide a resource to create a CDS attachment, can attach a CDS volume with instance.

Example Usage

```hcl
resource "baiducloud_cds_attachment" "default" {
  cds_id      = "v-FJjJeTiG"
  instance_id = "i-tgZhS50C"
}
```

Import

CDS attachment can be imported, e.g.

```hcl
$ terraform import baiducloud_cds_attachment.default id
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

func resourceBaiduCloudCDSAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudCDSAttachmentCreate,
		Read:   resourceBaiduCloudCDSAttachmentRead,
		Delete: resourceBaiduCloudCDSAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"cds_id": {
				Type:        schema.TypeString,
				Description: "CDS volume ID",
				Required:    true,
				ForceNew:    true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Description: "The ID of Instance which will attach CDS volume",
				Required:    true,
				ForceNew:    true,
			},
			"attachment_device": {
				Type:        schema.TypeString,
				Description: "CDS mount device path",
				Computed:    true,
			},
			"attachment_serial": {
				Type:        schema.TypeString,
				Description: "CDS serial",
				Computed:    true,
			},
		},
	}
}

func resourceBaiduCloudCDSAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	bccService := BccService{client}

	cdsId := d.Get("cds_id").(string)
	instanceId := d.Get("instance_id").(string)

	action := "Attach CDS volume " + cdsId + " with Instance " + instanceId
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		err := bccService.AttachCDSVolume(cdsId, instanceId)
		addDebug(action, err)
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cds_attachment", action, BCESDKGoERROR)
	}
	d.SetId(cdsId)

	return resourceBaiduCloudCDSAttachmentRead(d, meta)
}

func resourceBaiduCloudCDSAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	id := d.Id()
	action := "Query CDS volume " + id + " attachment"

	raw, err := client.WithBccClient(func(client *bcc.Client) (i interface{}, e error) {
		return client.GetCDSVolumeDetail(id)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cds_attachment", action, BCESDKGoERROR)
	}

	volume := raw.(*api.GetVolumeDetailResult).Volume
	if len(volume.Attachments) > 0 {
		d.Set("instance_id", volume.Attachments[0].InstanceId)
		d.Set("attachment_device", volume.Attachments[0].Device)
		d.Set("attachment_serial", volume.Attachments[0].Serial)
	}
	d.Set("cds_id", volume.Id)

	return nil
}

func resourceBaiduCloudCDSAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	bccService := BccService{client}

	id := d.Id()

	action := "Get CDS " + id + " detail"
	raw, err := client.WithBccClient(func(client *bcc.Client) (i interface{}, e error) {
		return client.GetCDSVolumeDetail(id)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cds_attachment", action, BCESDKGoERROR)
	}

	volume := raw.(*api.GetVolumeDetailResult).Volume
	if volume.Status == api.VolumeStatusINUSE {
		instanceId := volume.Attachments[0].InstanceId
		err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
			err := bccService.DetachCDSVolume(id, instanceId)

			if err != nil {
				if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}

			return nil
		})

		if err != nil {
			if IsExceptedErrors(err, ObjectNotFound) {
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cds_attachment", action, BCESDKGoERROR)
		}
	}

	return nil
}
