/*
Provide a resource to create an SMS template.

Example Usage

```hcl
resource "baiducloud_sms_template" "default" {
  name	         = "My test template"
  content        = "Test content"
  sms_type       = "CommonNotice"
  country_type   = "GLOBAL"
  description    = "this is a test sms template"

}
```

Import

SMS template can be imported, e.g.

```hcl
$ terraform import baiducloud_sms_template.default id
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

func resourceBaiduCloudSMSTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudSMSTemplateCreate,
		Read:   resourceBaiduCloudSMSTemplateRead,
		Update: resourceBaiduCloudSMSTemplateUpdate,
		Delete: resourceBaiduCloudSMSTemplateDelete,

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
				Description: "Template name of sms",
				Required:    true,
				ForceNew:    true,
			},
			"content": {
				Type:        schema.TypeString,
				Description: "Template content of sms",
				Required:    true,
				ForceNew:    true,
			},
			"sms_type": {
				Type:        schema.TypeString,
				Description: "Type of sms",
				Required:    true,
				ForceNew:    true,
			},
			"country_type": {
				Type:        schema.TypeString,
				Description: "Template type of country",
				Required:    true,
				ForceNew:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "LoadBalance instance's status",
				Required:    true,
				ForceNew:    true,
			},
			"template_id": {
				Type:        schema.TypeString,
				Description: "base64 of Template file",
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

func resourceBaiduCloudSMSTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	createArgs := buildBaiduCloudCreateSMSTemplateArgs(d)
	action := "Create SMS Template "

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithSMSClient(func(client *sms.Client) (i interface{}, e error) {
			return client.CreateTemplate(createArgs)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		addDebug(action, raw)
		response, _ := raw.(*api.CreateTemplateResult)
		d.SetId(response.TemplateId)
		d.Set("status", response.Status)

		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_sms_Template", action, BCESDKGoERROR)
	}

	return resourceBaiduCloudSMSTemplateRead(d, meta)
}
func resourceBaiduCloudSMSTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	smsService := SMSService{client}

	smsTemplateId := d.Id()
	action := "Query SMS Template " + smsTemplateId

	template, err := smsService.GetSMSTemplateDetail(smsTemplateId)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_sms_template", action, BCESDKGoERROR)
	}
	d.Set("template_id", template.TemplateId)
	d.Set("user_id", template.UserId)
	d.Set("name", template.Name)
	d.Set("content", template.Content)
	d.Set("country_type", template.CountryType)
	d.Set("sms_type", template.SmsType)
	d.Set("status", template.Status)
	d.Set("description", template.Description)
	d.Set("review", template.Review)

	return nil
}

func resourceBaiduCloudSMSTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	smsTemplateId := d.Id()
	action := "Update SMS Template " + smsTemplateId

	update := false
	updateArgs := &api.ModifyTemplateArgs{
		TemplateId: smsTemplateId,
	}

	if d.HasChange("name") || d.HasChange("content") || d.HasChange("country_type") ||
		d.HasChange("sms_type") || d.HasChange("description") {
		update = true
	}

	if update {
		d.Partial(true)

		if v, ok := d.GetOk("name"); ok && v.(string) != "" {
			updateArgs.Name = v.(string)
		}

		if v, ok := d.GetOk("content"); ok && v.(string) != "" {
			updateArgs.Content = v.(string)
		}

		if v, ok := d.GetOk("country_type"); ok && v.(string) != "" {
			updateArgs.CountryType = v.(string)
		}

		if v, ok := d.GetOk("sms_type"); ok && v.(string) != "" {
			updateArgs.SmsType = v.(string)
		}

		if v, ok := d.GetOk("description"); ok {
			updateArgs.Description = v.(string)
		}

		_, err := client.WithSMSClient(func(client *sms.Client) (i interface{}, e error) {
			return smsTemplateId, client.ModifyTemplate(updateArgs)
		})

		if err != nil {
			if NotFoundError(err) {
				d.SetId("")
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_sms_template", action, BCESDKGoERROR)
		}
	}

	d.Partial(false)
	return resourceBaiduCloudSMSTemplateRead(d, meta)
}

func resourceBaiduCloudSMSTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	smsTemplateId := d.Id()
	action := "Delete SMS TemplateId " + smsTemplateId
	deleteArgs := &api.DeleteTemplateArgs{
		TemplateId: smsTemplateId,
	}
	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, err := client.WithSMSClient(func(client *sms.Client) (i interface{}, e error) {
			return smsTemplateId, client.DeleteTemplate(deleteArgs)
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
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_sms_template", action, BCESDKGoERROR)
	}

	return nil
}

func buildBaiduCloudCreateSMSTemplateArgs(d *schema.ResourceData) *api.CreateTemplateArgs {
	result := &api.CreateTemplateArgs{}

	if v, ok := d.GetOk("name"); ok && v.(string) != "" {
		result.Name = v.(string)
	}

	if v, ok := d.GetOk("content"); ok && v.(string) != "" {
		result.Content = v.(string)
	}

	if v, ok := d.GetOk("country_type"); ok && v.(string) != "" {
		result.CountryType = v.(string)
	}

	if v, ok := d.GetOk("sms_type"); ok && v.(string) != "" {
		result.SmsType = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		result.Description = v.(string)
	}

	return result
}
