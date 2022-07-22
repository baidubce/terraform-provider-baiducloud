/*
Provide a resource to create an EIP association, bind an EIP with instance.

Example Usage

```hcl
resource "baiducloud_eip_association" "default" {
  eip           = "1.1.1.1"
  instance_type = "BCC"
  instance_id   = "i-7xc9Q6KR"
}
```

Import

EIP association can be imported, e.g.

```hcl
$ terraform import baiducloud_eip_association.default eip
```
*/
package baiducloud

import (
	"time"

	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func resourceBaiduCloudEipAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudEipAssociationCreate,
		Read:   resourceBaiduCloudEipAssociationRead,
		Delete: resourceBaiduCloudEipAssociationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"eip": {
				Type:         schema.TypeString,
				Description:  "EIP which need to associate with instance",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.SingleIP(),
			},
			"instance_type": {
				Type:         schema.TypeString,
				Description:  "Instance type which need to associate with EIP, support BCC/BLB/NAT/VPN",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"BCC", "BLB", "NAT", "VPN"}, false),
			},
			"instance_id": {
				Type:        schema.TypeString,
				Description: "Instance ID which need to associate with EIP",
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceBaiduCloudEipAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	eipClient := EipService{client}

	eipAddress := d.Get("eip").(string)
	instanceId := d.Get("instance_id").(string)
	instanceType := d.Get("instance_type").(string)
	action := "Bind EIP " + eipAddress + " with " + instanceId

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		errDelete := eipClient.EipBind(eipAddress, instanceType, instanceId)
		addDebug(action, errDelete)
		if errDelete != nil {
			if IsExceptedErrors(errDelete, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(errDelete)
			}
			return resource.NonRetryableError(errDelete)
		}

		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_eip_association", action, BCESDKGoERROR)
	}

	d.SetId(eipAddress)

	stateConf := buildStateConf(EIPProcessingStatus,
		[]string{EIPStatusBinded},
		d.Timeout(schema.TimeoutCreate),
		eipClient.EipStateRefreshFunc(eipAddress, append(EIPFailedStatus, EIPStatusAvailable)))
	if _, err = stateConf.WaitForState(); err != nil {
		return WrapError(err)
	}

	return resourceBaiduCloudEipAssociationRead(d, meta)
}

func resourceBaiduCloudEipAssociationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	eipClient := EipService{client}

	eipAddress := d.Id()
	action := "Query EIP " + eipAddress + "association"
	result, err := eipClient.EipGetDetail(eipAddress)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_eip_association", action, BCESDKGoERROR)
	}

	d.Set("eip", result.Eip)
	d.Set("instance_id", result.InstanceId)
	d.Set("instance_type", result.InstanceType)

	return nil
}

func resourceBaiduCloudEipAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	eipClient := EipService{client}

	eipAddress := d.Id()
	action := "Unbind EIP " + eipAddress
	result, err := eipClient.EipGetDetail(eipAddress)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_eip_association", action, BCESDKGoERROR)
	}

	if result.Status != EIPStatusBinded {
		return nil
	}

	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		errDelete := eipClient.EipUnBind(eipAddress)
		addDebug(action, errDelete)
		if errDelete != nil {
			if IsExceptedErrors(errDelete, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(errDelete)
			}
			return resource.NonRetryableError(errDelete)
		}

		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_eip_association", action, BCESDKGoERROR)
	}

	stateConf := buildStateConf(EIPProcessingStatus,
		[]string{EIPStatusAvailable},
		d.Timeout(schema.TimeoutDelete),
		eipClient.EipStateRefreshFunc(eipAddress, append(EIPFailedStatus, EIPStatusBinded)))
	if _, err = stateConf.WaitForState(); err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	return nil
}
