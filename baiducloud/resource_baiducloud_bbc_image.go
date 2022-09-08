/*
Use this resource to create BBC custom image.

Example Usage
```hcl
resource "baiducloud_bbc_image" "test-image" {
  image_name = "terrform-bbc-image-test"
  instance_id = "i-qwIq4vKi"
}
```
*/
package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/bbc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
	"time"
)

func resourceBaiduCloudBbcImage() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudBbcImageCreate,
		Read:   resourceBaiduCloudBbcImageRead,
		Delete: resourceBaiduCloudBbcImageDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"image_name": {
				Type:        schema.TypeString,
				Description: "Image name.",
				Required:    true,
				ForceNew:    true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Description: "The id of the image source instance.",
				Required:    true,
				ForceNew:    true,
			},
			"image_id": {
				Type:        schema.TypeString,
				Description: "Image id.Computed after apply",
				Computed:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "Image type.System, Custom, Integration",
				Computed:    true,
			},
			"os_type": {
				Type:        schema.TypeString,
				Description: "Image os type.CentOS, Windows, etc.",
				Computed:    true,
			},
			"os_version": {
				Type:        schema.TypeString,
				Description: "Image os version.",
				Computed:    true,
			},
			"os_arch": {
				Type:        schema.TypeString,
				Description: "Image os arch.",
				Computed:    true,
			},
			"os_name": {
				Type:        schema.TypeString,
				Description: "Image os name.",
				Computed:    true,
			},
			"os_build": {
				Type:        schema.TypeString,
				Description: "Image os build.",
				Computed:    true,
			},
			"create_time": {
				Type:        schema.TypeString,
				Description: "The creation time of the image, in a date format that conforms to the BCE specification",
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Image status.Creating, CreatedFailed, Available, NotAvailable, Error",
				Computed:    true,
			},
			"desc": {
				Type:        schema.TypeString,
				Description: "Image description",
				Computed:    true,
			},
		},
	}
}
func resourceBaiduCloudBbcImageCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	bbcService := BbcService{client}
	action := "Create bbc custom image"

	createImageArgs := &bbc.CreateImageArgs{
		ClientToken: buildClientToken(),
	}
	if iamgeName, ok := d.GetOk("image_name"); ok {
		createImageArgs.ImageName = iamgeName.(string)
	}
	if instanceId, ok := d.GetOk("instance_id"); ok {
		createImageArgs.InstanceId = instanceId.(string)
	}
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		res, err := client.WithBbcClient(func(bbcClient *bbc.Client) (i interface{}, e error) {
			return bbcClient.CreateImageFromInstanceId(createImageArgs)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		d.SetId(res.(*bbc.CreateImageResult).ImageId)
		return nil
	})
	stateConf := buildStateConf(
		[]string{string(bbc.ImageStatusCreating)},
		[]string{string(bbc.ImageStatusAvailable)},
		d.Timeout(schema.TimeoutCreate),
		bbcService.BbcImageStateRefresh(d.Id()),
	)
	if _, err = stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_instance", action, BCESDKGoERROR)
	}
	return resourceBaiduCloudBbcImageRead(d, meta)
}
func resourceBaiduCloudBbcImageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	bbcService := BbcService{client}

	action := "Query Bbc Image detail."

	imageResult, err := bbcService.GetBbcImageDetails(d.Id())
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_images", action, BCESDKGoERROR)
	}
	addDebug(action, imageResult)
	d.Set("image_id", imageResult.Id)
	d.Set("image_name", imageResult.Name)
	d.Set("type", imageResult.Type)
	d.Set("os_type", imageResult.OsType)
	d.Set("os_version", imageResult.OsVersion)
	d.Set("os_arch", imageResult.OsArch)
	d.Set("os_name", imageResult.OsName)
	d.Set("os_build", imageResult.OsBuild)
	d.Set("create_time", imageResult.CreateTime)
	d.Set("status", imageResult.Status)
	d.Set("desc", imageResult.Desc)
	return nil
}
func resourceBaiduCloudBbcImageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	imageId := d.Id()
	action := "Delete BBC Image " + imageId

	// delete bbc Image
	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		raw, err := client.WithBbcClient(func(bbcClient *bbc.Client) (interface{}, error) {
			return imageId, bbcClient.DeleteImage(imageId)
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
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_image", action, BCESDKGoERROR)
	}

	return nil
}
