/*
Provide a resource to create an BLB Backend Server.

Example Usage

```hcl
resource "baiducloud_blb_backend_server" "default" {
  blb_id      = "lb-0d29xxx6"

  backend_server_list {
    instance_id = "i-VRxxxx1a"
    weight = 50
  }
}
```
*/
package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/blb"
	"time"

	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func resourceBaiduCloudBlbBackendServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudBlbBackendServerCreate,
		Read:   resourceBaiduCloudBlbBackendServerRead,
		Update: resourceBaiduCloudBlbBackendServerUpdate,
		Delete: resourceBaiduCloudBlbBackendServerDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"blb_id": {
				Type:        schema.TypeString,
				Description: "ID of the lication LoadBalance instance",
				Required:    true,
				ForceNew:    true,
			},
			"backend_server_list": {
				Type:        schema.TypeList,
				Description: "Server group bound backend server list",
				Required:    true,
				//Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Description: "Backend server instance ID",
							Required:    true,
							ForceNew:    true,
						},
						"weight": {
							Type:         schema.TypeInt,
							Description:  "Backend server instance weight in this group, range from 0-100",
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 100),
						},
						"private_ip": {
							Type:        schema.TypeString,
							Description: "Backend server instance bind private ip",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func resourceBaiduCloudBlbBackendServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	createArgs := buildBaiduCloudCreateBlbServerGroupArgs(d)
	blbId := d.Get("blb_id").(string)
	action := "Create Blb " + blbId + " BackendServer "
	addDebug(action, createArgs)
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithBLBClient(func(client *blb.Client) (i interface{}, e error) {
			return nil, client.AddBackendServers(blbId, createArgs)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		addDebug(action, raw)
		d.SetId(blbId)

		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_blb_backend_server", action, BCESDKGoERROR)
	}
	return resourceBaiduCloudBlbBackendServerRead(d, meta)
}

func resourceBaiduCloudBlbBackendServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	blbService := BLBService{client}

	blbId := d.Get("blb_id").(string)
	action := "Query BLB " + blbId + " BackendServer "

	servers, err := blbService.BackendServerList(blbId)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_blb_backend_server", action, BCESDKGoERROR)
	}
	addDebug(action, servers)

	return nil
}

func resourceBaiduCloudBlbBackendServerUpdate(d *schema.ResourceData, meta interface{}) error {

	if d.HasChange("backend_server_list") {

		if err := updateBackendServer(d, meta); err != nil {
			return err
		}
	}

	return resourceBaiduCloudBlbBackendServerRead(d, meta)
}

func resourceBaiduCloudBlbBackendServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	blbId := d.Get("blb_id").(string)
	id := d.Id()
	backendServerList := d.Get("backend_server_list").([]interface{})

	deleteArgs := buildRemoveBackendServersArgs(backendServerList)

	action := "Delete BLB " + blbId + "  Server " + id

	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, err := client.WithBLBClient(func(client *blb.Client) (i interface{}, e error) {
			return id, client.RemoveBackendServers(blbId, deleteArgs)
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
		if IsExceptedErrors(err, ObjectNotFound) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_blb_backend_server", action, BCESDKGoERROR)
	}

	return nil
}

func buildBaiduCloudCreateBlbServerGroupArgs(d *schema.ResourceData) *blb.AddBackendServersArgs {
	result := &blb.AddBackendServersArgs{
		ClientToken: buildClientToken(),
	}

	if v, ok := d.Get("backend_server_list").([]interface{}); ok && len(v) > 0 {
		for _, value := range v {
			m := value.(map[string]interface{})

			result.BackendServerList = append(result.BackendServerList, blb.BackendServerModel{
				InstanceId: m["instance_id"].(string),
				Weight:     m["weight"].(int),
			})
		}
	}

	return result
}

func updateBackendServer(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	updateArgs := &blb.UpdateBackendServersArgs{
		ClientToken: buildClientToken(),
	}

	blbId := d.Get("blb_id").(string)

	if v, ok := d.Get("backend_server_list").([]interface{}); ok && len(v) > 0 {
		for _, value := range v {
			m := value.(map[string]interface{})

			updateArgs.BackendServerList = append(updateArgs.BackendServerList, blb.BackendServerModel{
				InstanceId: m["instance_id"].(string),
				Weight:     m["weight"].(int),
			})
		}
	}

	_, err := client.WithBLBClient(func(client *blb.Client) (i interface{}, e error) {
		return nil, client.UpdateBackendServers(blbId, updateArgs)
	})

	return err
}

func buildRemoveBackendServersArgs(list []interface{}) *blb.RemoveBackendServersArgs {
	if len(list) == 0 {
		return nil
	}

	result := &blb.RemoveBackendServersArgs{
		ClientToken: buildClientToken(),
	}
	for _, v := range list {
		removeValue := v.(map[string]interface{})
		result.BackendServerList = append(result.BackendServerList, removeValue["instance_id"].(string))
	}

	return result
}
