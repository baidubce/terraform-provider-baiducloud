/*
Use this resource to get information about a SCS Security Ip.

~> **NOTE:** The terminate operation of scs instance does NOT take effect immediately，maybe takes for several minites.

Example Usage

```hcl
resource "baiducloud_scs_security_ip" "default" {
    instance_id                    = "scs-xxxxx"
    security_ips                   = [192.168.0.8]
}
```

Import

SCS Security Ip. can be imported, e.g.

```hcl
$ terraform import baiducloud_scs_security_ip.default id
```
*/
package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/scs"
	"time"

	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func resourceBaiduCloudScsSecurityIp() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudScsSecurityIpCreate,
		Read:   resourceBaiduCloudScsSecurityIpRead,
		Update: resourceBaiduCloudScsSecurityIpUpdate,
		Delete: resourceBaiduCloudScsSecurityIpDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Description: "ID of the instance",
				Required:    true,
				ForceNew:    true,
			},
			"security_ips": {
				Type:        schema.TypeSet,
				Description: "securityIps",
				Optional:    true,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceBaiduCloudScsSecurityIpCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	instanceIdArg := d.Get("instance_id").(string)

	updateSecurityArgs, err := buildBaiduCloudScsSecurityIpArgs(d, meta)

	if err != nil {
		return WrapError(err)
	}

	action := "Create Scs SecurityIp instance id is" + instanceIdArg
	addDebug(action, updateSecurityArgs)

	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		_, err := client.WithScsClient(func(scsClient *scs.Client) (interface{}, error) {

			return nil, scsClient.AddSecurityIp(instanceIdArg, updateSecurityArgs)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, err)
		d.SetId(instanceIdArg)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_scs_security_ip", action, BCESDKGoERROR)
	}

	return resourceBaiduCloudScsSecurityIpRead(d, meta)
}

func resourceBaiduCloudScsSecurityIpRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	instanceID := d.Id()
	action := "Query Scs SecurityIp instanceID is " + instanceID

	raw, err := client.WithScsClient(func(scsClient *scs.Client) (interface{}, error) {
		return scsClient.GetSecurityIp(instanceID)
	})

	addDebug(action, raw)

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "bbaiducloud_scs_security_ip", action, BCESDKGoERROR)
	}

	result, _ := raw.(*scs.GetSecurityIpResult)

	d.Set("security_ips", result.SecurityIps)

	return nil
}

func resourceBaiduCloudScsSecurityIpUpdate(d *schema.ResourceData, meta interface{}) error {

	if !d.HasChange("security_ips") {
		return resourceBaiduCloudScsSecurityIpRead(d, meta)
	}

	client := meta.(*connectivity.BaiduClient)
	instanceID := d.Id()

	updateSecurityArgs, err := buildBaiduCloudScsSecurityIpArgs(d, meta)

	if err != nil {
		return WrapError(err)
	}

	action := "Update Scs SecurityIp instance id is" + instanceID
	addDebug(action, updateSecurityArgs)

	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		_, err := client.WithScsClient(func(scsClient *scs.Client) (interface{}, error) {
			oldIps, e := scsClient.GetSecurityIp(instanceID)
			if e != nil {
				addDebug("Update Scs SecurityIp : Get Old Ips instanceId is "+instanceID, e)
				return nil, e
			}
			deleteRequest := &scs.SecurityIpArgs{}
			deleteRequest.SecurityIps = oldIps.SecurityIps
			deleteRequest.ClientToken = buildClientToken()
			e = scsClient.DeleteSecurityIp(instanceID, deleteRequest)
			if e != nil {
				addDebug("Update Scs SecurityIp : Delete Old Ips instanceId is "+instanceID, e)
				return nil, e
			}
			return nil, scsClient.AddSecurityIp(instanceID, updateSecurityArgs)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, err)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_scs_security_ip", action, BCESDKGoERROR)
	}

	return resourceBaiduCloudScsSecurityIpRead(d, meta)
}

func resourceBaiduCloudScsSecurityIpDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	instanceId := d.Id()

	action := "Delete Scs SecurityIp instance id is" + instanceId

	request, err := buildBaiduCloudScsSecurityIpArgs(d, meta)
	addDebug(action, "")

	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, err := client.WithScsClient(func(scsClient *scs.Client) (interface{}, error) {

			return nil, scsClient.DeleteSecurityIp(instanceId, request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, err)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_scs_security_ip", action, BCESDKGoERROR)
	}

	return nil
}

func buildBaiduCloudScsSecurityIpArgs(d *schema.ResourceData, meta interface{}) (*scs.SecurityIpArgs, error) {
	request := &scs.SecurityIpArgs{}

	if securityIps, ok := d.GetOk("security_ips"); ok {

		ips := make([]string, 0)
		for _, ip := range securityIps.(*schema.Set).List() {
			ips = append(ips, ip.(string))
		}
		request.SecurityIps = ips
	}

	request.ClientToken = buildClientToken()

	return request, nil

}
