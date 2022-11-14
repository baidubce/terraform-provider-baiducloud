/*
Provide a resource to create an BLS LogStore.

Example Usage

```hcl
resource "baiducloud_bls_log_store" "default" {
  log_store_name   = "MyTest"
  retention        = 10

}
```

Import

BLS LogStore can be imported, e.g.

```hcl
$ terraform import baiducloud_bls_log_store.default id
```
*/
package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/bls"
	"time"

	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func resourceBaiduCloudBLSLogStore() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudBLSLogStoreCreate,
		Read:   resourceBaiduCloudBLSLogStoreRead,
		Update: resourceBaiduCloudBLSLogStoreUpdate,
		Delete: resourceBaiduCloudBLSLogStoreDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"log_store_name": {
				Type:        schema.TypeString,
				Description: "name of log store",
				Required:    true,
				ForceNew:    true,
			},
			"retention": {
				Type:        schema.TypeInt,
				Description: "retention days of log store",
				Required:    true,
			},
			"creation_date_time": {
				Type:        schema.TypeString,
				Description: "log store create date time",
				Computed:    true,
			},
			"last_modified_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceBaiduCloudBLSLogStoreCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	logStoreName := d.Get("log_store_name").(string)
	retention := d.Get("retention").(int)

	action := "Create BLS LogStore "

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithBLSClient(func(client *bls.Client) (i interface{}, e error) {
			return nil, client.CreateLogStore(logStoreName, retention)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		addDebug(action, raw)
		d.SetId(resource.UniqueId())

		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bls_log_store", action, BCESDKGoERROR)
	}

	return resourceBaiduCloudBLSLogStoreRead(d, meta)
}
func resourceBaiduCloudBLSLogStoreRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	blsService := BLSService{client}

	logStoreName := d.Get("log_store_name").(string)
	action := "Query BLS LogStore " + logStoreName

	logStore, err := blsService.GetBLSLogStoreDetail(logStoreName)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bls_log_store", action, BCESDKGoERROR)
	}
	d.Set("creation_date_time", logStore.CreationDateTime)
	d.Set("last_modified_time", logStore.LastModifiedTime)

	return nil
}

func resourceBaiduCloudBLSLogStoreUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	action := "Update BLS LogStore " + d.Get("log_store_name").(string)

	update := false
	if d.HasChange("retention") {
		update = true
	}

	if update {
		d.Partial(true)
		logStoreName := d.Get("log_store_name").(string)
		retention := d.Get("retention").(int)

		_, err := client.WithBLSClient(func(client *bls.Client) (i interface{}, e error) {
			return nil, client.UpdateLogStore(logStoreName, retention)
		})

		if err != nil {
			if NotFoundError(err) {
				d.SetId("")
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bls_log_store", action, BCESDKGoERROR)
		}
	}

	d.Partial(false)
	return resourceBaiduCloudBLSLogStoreRead(d, meta)
}

func resourceBaiduCloudBLSLogStoreDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	logStoreName := d.Get("log_store_name").(string)

	action := "Delete BLS LogStore " + logStoreName
	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, err := client.WithBLSClient(func(client *bls.Client) (i interface{}, e error) {
			return nil, client.DeleteLogStore(logStoreName)
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
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bls_log_store", action, BCESDKGoERROR)
	}

	return nil
}
