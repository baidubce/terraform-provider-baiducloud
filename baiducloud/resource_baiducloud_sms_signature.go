/*
Provide a resource to create an SMS signature.

Example Usage

```hcl
resource "baiducloud_sms_signature" "default" {
  content        = "baidu"
  description    = "this is a test sms signature"
  content_type   = "Enterprise"
  country_type   = "DOMESTIC"

}
```

Import

SMS signature can be imported, e.g.

```hcl
$ terraform import baiducloud_sms_signature.default id
```
*/
package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/sms"
	"github.com/baidubce/bce-sdk-go/services/sms/api"
	"time"

	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func resourceBaiduCloudSMSSignature() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudSMSSignatureCreate,
		Read:   resourceBaiduCloudSMSSignatureRead,
		Update: resourceBaiduCloudSMSSignatureUpdate,
		Delete: resourceBaiduCloudSMSSignatureDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"content": {
				Type:        schema.TypeString,
				Description: "signature content of sms",
				Required:    true,
				ForceNew:    true,
			},
			"content_type": {
				Type:        schema.TypeString,
				Description: "type of content",
				Required:    true,
				ForceNew:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "LoadBalance instance's status",
				Optional:    true,
			},
			"country_type": {
				Type:        schema.TypeString,
				Description: "signature type of country",
				Optional:    true,
				ForceNew:    true,
			},
			"signature_file_base64": {
				Type:        schema.TypeString,
				Description: "base64 of signature file",
				Optional:    true,
			},
			"signature_file_format": {
				Type:        schema.TypeString,
				Description: "Format of signature file",
				Optional:    true,
			},
			"user_id": {
				Type:        schema.TypeString,
				Description: "User id",
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "status",
				Computed:    true,
			},
			"review": {
				Type:        schema.TypeString,
				Description: "commit review",
				Computed:    true,
			},
		},
	}
}

func resourceBaiduCloudSMSSignatureCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	createArgs := buildBaiduCloudCreateSMSSignatureArgs(d)
	action := "Create SMS Signature "

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithSMSClient(func(client *sms.Client) (i interface{}, e error) {
			return client.CreateSignature(createArgs)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		addDebug(action, raw)
		response, _ := raw.(*api.CreateSignatureResult)
		d.SetId(response.SignatureId)
		d.Set("status", response.Status)

		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_sms_signature", action, BCESDKGoERROR)
	}

	return resourceBaiduCloudSMSSignatureRead(d, meta)
}
func resourceBaiduCloudSMSSignatureRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	smsService := SMSService{client}

	smsSignatureId := d.Id()
	action := "Query SMS signature " + smsSignatureId

	signature, err := smsService.GetSMSSignatureDetail(smsSignatureId)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_sms_signature", action, BCESDKGoERROR)
	}
	d.Set("content", signature.Content)
	d.Set("content_type", signature.ContentType)
	d.Set("country_type", signature.CountryType)
	d.Set("review", signature.Review)
	d.Set("user_id", signature.UserId)
	d.Set("status", signature.Status)

	return nil
}

func resourceBaiduCloudSMSSignatureUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	smsSignatureId := d.Id()
	action := "Update SMS Signature " + smsSignatureId

	update := false
	updateArgs := &api.ModifySignatureArgs{
		SignatureId: smsSignatureId,
	}

	if d.HasChange("content") || d.HasChange("content_type") || d.HasChange("country_type") {
		update = true
	}

	if update {
		d.Partial(true)

		if v, ok := d.GetOk("content"); ok && v.(string) != "" {
			updateArgs.Content = v.(string)
		}

		if v, ok := d.GetOk("content_type"); ok && v.(string) != "" {
			updateArgs.ContentType = v.(string)
		}

		if v, ok := d.GetOk("description"); ok && v.(string) != "" {
			updateArgs.Description = v.(string)
		}

		if v, ok := d.GetOk("country_type"); ok && v.(string) != "" {
			updateArgs.CountryType = v.(string)
		}

		if v, ok := d.GetOk("signature_file_base64"); ok {
			updateArgs.SignatureFileBase64 = v.(string)
		}

		if v, ok := d.GetOk("signature_file_format"); ok {
			updateArgs.SignatureFileFormat = v.(string)
		}

		_, err := client.WithSMSClient(func(client *sms.Client) (i interface{}, e error) {
			return smsSignatureId, client.ModifySignature(updateArgs)
		})

		if err != nil {
			if NotFoundError(err) {
				d.SetId("")
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_sms_signature", action, BCESDKGoERROR)
		}
	}

	d.Partial(false)
	return resourceBaiduCloudSMSSignatureRead(d, meta)
}

func resourceBaiduCloudSMSSignatureDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	smsSignatureId := d.Id()
	action := "Delete SMS SignatureId " + smsSignatureId
	deleteArgs := &api.DeleteSignatureArgs{
		SignatureId: smsSignatureId,
	}
	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, err := client.WithSMSClient(func(client *sms.Client) (i interface{}, e error) {
			return smsSignatureId, client.DeleteSignature(deleteArgs)
		})
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
		if IsExceptedErrors(err, ObjectNotFound) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_sms_signature", action, BCESDKGoERROR)
	}

	return nil
}

func buildBaiduCloudCreateSMSSignatureArgs(d *schema.ResourceData) *api.CreateSignatureArgs {
	result := &api.CreateSignatureArgs{}

	if v, ok := d.GetOk("content"); ok && v.(string) != "" {
		result.Content = v.(string)
	}

	if v, ok := d.GetOk("content_type"); ok && v.(string) != "" {
		result.ContentType = v.(string)
	}

	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		result.Description = v.(string)
	}

	if v, ok := d.GetOk("country_type"); ok && v.(string) != "" {
		result.CountryType = v.(string)
	}

	if v, ok := d.GetOk("signature_file_base64"); ok {
		result.SignatureFileBase64 = v.(string)
	}

	if v, ok := d.GetOk("signature_file_format"); ok {
		result.SignatureFileFormat = v.(string)
	}

	return result
}
