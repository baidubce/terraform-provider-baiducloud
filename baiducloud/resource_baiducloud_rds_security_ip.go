/*
Use this resource to get information about a RDS Security Ip.

~> **NOTE:** The terminate operation of rds instance does NOT take effect immediately，maybe takes for several minites.

Example Usage

```hcl
resource "baiducloud_rds_security_ip" "default" {
    instance_id                    = "rds-ZuZd7s1l"
    security_ips                   = [192.168.0.8]
}
```

Import

RDS RDS Security Ip. can be imported, e.g.

```hcl
$ terraform import baiducloud_rds_security_ip.default id
```
*/
package baiducloud

import (
	"log"
	"time"

	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/rds"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func resourceBaiduCloudRdsSecurityIp() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudRdsSecurityIpCreate,
		Read:   resourceBaiduCloudRdsSecurityIpRead,
		Update: resourceBaiduCloudRdsSecurityIpUpdate,
		Delete: resourceBaiduCloudRdsSecurityIpDelete,

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
			"e_tag": {
				Type:        schema.TypeString,
				Description: "ETag of the instance.",
				Computed:    true,
			},
		},
	}
}

func resourceBaiduCloudRdsSecurityIpCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	instanceIdArg := d.Get("instance_id").(string)

	updateSecurityArgs, err := buildBaiduCloudRdsSecurityIpArgs(d, meta)

	if err != nil {
		return WrapError(err)
	}

	action := "Create RDS SecurityIp instance id is" + instanceIdArg
	addDebug(action, updateSecurityArgs)

	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		_, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {

			result, e := rdsClient.GetSecurityIps(instanceIdArg)
			log.Printf("GetSecurityIps Etag is:" + result.Etag)
			if e != nil {
				return nil, e
			}
			return nil, rdsClient.UpdateSecurityIps(instanceIdArg, result.Etag, updateSecurityArgs)
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
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_rds_security_ip", action, BCESDKGoERROR)
	}

	return resourceBaiduCloudRdsSecurityIpRead(d, meta)
}

func resourceBaiduCloudRdsSecurityIpRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	instanceID := d.Id()
	action := "Query RDS SecurityIp instanceID is " + instanceID

	raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
		return rdsClient.GetSecurityIps(instanceID)
	})

	addDebug(action, raw)

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_rds_instance", action, BCESDKGoERROR)
	}

	result, _ := raw.(*rds.GetSecurityIpsResult)

	d.Set("e_tag", result.Etag)
	d.Set("security_ips", result.SecurityIps)

	return nil
}

func resourceBaiduCloudRdsSecurityIpUpdate(d *schema.ResourceData, meta interface{}) error {

	if !d.HasChange("security_ips") {
		return resourceBaiduCloudRdsSecurityIpRead(d, meta)
	}

	client := meta.(*connectivity.BaiduClient)
	instanceID := d.Id()

	updateSecurityArgs, err := buildBaiduCloudRdsSecurityIpArgs(d, meta)

	if err != nil {
		return WrapError(err)
	}

	action := "Update RDS SecurityIp instance id is" + instanceID
	addDebug(action, updateSecurityArgs)

	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		_, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {

			result, e := rdsClient.GetSecurityIps(instanceID)
			log.Printf("GetSecurityIps Etag is:" + result.Etag)
			if e != nil {
				return nil, e
			}
			return nil, rdsClient.UpdateSecurityIps(instanceID, result.Etag, updateSecurityArgs)
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
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_rds_security_ip", action, BCESDKGoERROR)
	}

	return resourceBaiduCloudRdsSecurityIpRead(d, meta)
}

func resourceBaiduCloudRdsSecurityIpDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	instanceId := d.Id()

	action := "Delete RDS SecurityIp instance id is" + instanceId

	request := &rds.UpdateSecurityIpsArgs{}
	ips := make([]string, 0)
	//ips = append(ips, "")
	request.SecurityIps = ips

	addDebug(action, "")

	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {

			result, e := rdsClient.GetSecurityIps(instanceId)
			log.Printf("GetSecurityIps Etag is:" + result.Etag)
			if e != nil {
				return nil, e
			}
			return nil, rdsClient.UpdateSecurityIps(instanceId, result.Etag, request)
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
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_rds_security_ip", action, BCESDKGoERROR)
	}

	return nil
}

func buildBaiduCloudRdsSecurityIpArgs(d *schema.ResourceData, meta interface{}) (*rds.UpdateSecurityIpsArgs, error) {
	request := &rds.UpdateSecurityIpsArgs{}

	if securityIps, ok := d.GetOk("security_ips"); ok {

		ips := make([]string, 0)
		for _, ip := range securityIps.(*schema.Set).List() {
			ips = append(ips, ip.(string))
		}
		request.SecurityIps = ips
	}

	return request, nil

}
