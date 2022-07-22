/*
Provide a resource to Upload a cert.

Example Usage

```hcl
resource "baiducloud_cert" "cert" {
  cert_name         = "testCert"
  cert_server_data  = ""
  cert_private_data = ""
}
```
*/
package baiducloud

import (
	"time"

	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/cert"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func resourceBaiduCloudCert() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudCertCreate,
		Read:   resourceBaiduCloudCertRead,
		Update: resourceBaiduCloudCertUpdate,
		Delete: resourceBaiduCloudCertDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"cert_name": {
				Type:        schema.TypeString,
				Description: "Cert Name",
				Required:    true,
			},
			"cert_server_data": {
				Type:        schema.TypeString,
				Description: "Server Cert data, base64 encode",
				Required:    true,
			},
			"cert_private_data": {
				Type:        schema.TypeString,
				Description: "Cert private key data, base64 encode",
				Required:    true,
			},
			"cert_link_data": {
				Type:        schema.TypeString,
				Description: "Cert lint data, base64 encode",
				Optional:    true,
			},
			"cert_type": {
				Type:        schema.TypeInt,
				Description: "Cert type",
				Optional:    true,
				Computed:    true,
			},
			"cert_common_name": {
				Type:        schema.TypeString,
				Description: "Cert common name",
				Computed:    true,
			},
			"cert_start_time": {
				Type:        schema.TypeString,
				Description: "Cert start time",
				Computed:    true,
			},
			"cert_stop_time": {
				Type:        schema.TypeString,
				Description: "Cert stop time",
				Computed:    true,
			},
			"cert_create_time": {
				Type:        schema.TypeString,
				Description: "Cert create time",
				Computed:    true,
			},
			"cert_update_time": {
				Type:        schema.TypeString,
				Description: "Cert update time",
				Computed:    true,
			},
		},
	}
}

func resourceBaiduCloudCertCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	args := buildBaiduCloudCreateCertArgs(d)
	action := "Create Cert " + args.CertName

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithCertClient(func(client *cert.Client) (i interface{}, e error) {
			return client.CreateCert(args)
		})

		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		response := raw.(*cert.CreateCertResult)
		d.SetId(response.CertId)

		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cert", action, BCESDKGoERROR)
	}

	return resourceBaiduCloudCertRead(d, meta)
}

func resourceBaiduCloudCertRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	id := d.Id()
	action := "Get Cert " + id + " Meta"
	raw, err := client.WithCertClient(func(client *cert.Client) (i interface{}, e error) {
		return client.GetCertMeta(id)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cert", action, BCESDKGoERROR)
	}

	certMeta := raw.(*cert.CertificateMeta)
	d.Set("cert_name", certMeta.CertName)
	d.Set("cert_common_name", certMeta.CertCommonName)
	d.Set("cert_start_time", certMeta.CertStartTime)
	d.Set("cert_stop_time", certMeta.CertStopTime)
	d.Set("cert_create_time", certMeta.CertCreateTime)
	d.Set("cert_update_time", certMeta.CertUpdateTime)
	d.Set("cert_type", certMeta.CertType)

	return nil
}

func resourceBaiduCloudCertUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	d.Partial(true)
	id := d.Id()

	if d.HasChange("cert_name") {
		action := "Update Cert " + id + " Name"
		updateNameArgs := &cert.UpdateCertNameArgs{
			CertName: d.Get("cert_name").(string),
		}

		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			_, err := client.WithCertClient(func(client *cert.Client) (i interface{}, e error) {
				return nil, client.UpdateCertName(id, updateNameArgs)
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
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cert", action, BCESDKGoERROR)
		}

		d.SetPartial("cert_name")
	}

	if d.HasChange("cert_server_data") || d.HasChange("cert_private_data") || d.HasChange("cert_link_data") || d.HasChange("cert_type") {
		updateDataArgs := &cert.UpdateCertDataArgs{
			CertName:        d.Get("cert_name").(string),
			CertServerData:  d.Get("cert_server_data").(string),
			CertPrivateData: d.Get("cert_private_data").(string),
		}
		if v, ok := d.GetOk("cert_link_data"); ok && v.(string) != "" {
			updateDataArgs.CertLinkData = v.(string)
		}

		if v, ok := d.GetOk("cert_type"); ok {
			updateDataArgs.CertType = v.(int)
		}

		action := "Update Cert " + id + " Data"
		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			_, err := client.WithCertClient(func(client *cert.Client) (i interface{}, e error) {
				return nil, client.UpdateCertData(id, updateDataArgs)
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
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cert", action, BCESDKGoERROR)
		}

		d.SetPartial("cert_server_data")
		d.SetPartial("cert_private_data")
		d.SetPartial("cert_link_data")
		d.SetPartial("cert_type")
	}

	return resourceBaiduCloudCertRead(d, meta)
}

func resourceBaiduCloudCertDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	id := d.Id()
	action := "Delete Cert " + id

	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, err := client.WithCertClient(func(client *cert.Client) (i interface{}, e error) {
			return nil, client.DeleteCert(id)
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
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cert", action, BCESDKGoERROR)
	}

	return nil
}

func buildBaiduCloudCreateCertArgs(d *schema.ResourceData) *cert.CreateCertArgs {
	result := &cert.CreateCertArgs{}

	if v, ok := d.GetOk("cert_name"); ok && v.(string) != "" {
		result.CertName = v.(string)
	}

	if v, ok := d.GetOk("cert_server_data"); ok && v.(string) != "" {
		result.CertServerData = v.(string)
	}

	if v, ok := d.GetOk("cert_private_data"); ok && v.(string) != "" {
		result.CertPrivateData = v.(string)
	}

	if v, ok := d.GetOk("cert_link_data"); ok && v.(string) != "" {
		result.CertLinkData = v.(string)
	}

	if v, ok := d.GetOk("cert_type"); ok {
		result.CertType = v.(int)
	}

	return result
}
