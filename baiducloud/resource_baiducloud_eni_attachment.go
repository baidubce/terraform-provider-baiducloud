/*
Provide a resource to create an ENI association, bind an ENI with instance.

Example Usage

```hcl
data "baiducloud_images" "images" {
  image_type = "System"
  name_regex = "8.4 aarch"
  os_name    = "CentOS"
}

resource "baiducloud_vpc" "vpc" {
  name = "terraform_vpc"
  cidr = "172.16.0.0/20"
}
resource "baiducloud_subnet" "subnet" {
  name        = "terraform_subnet"
  zone_name   = "cn-bj-d"
  cidr        = "172.16.0.0/24"
  vpc_id      = baiducloud_vpc.vpc.id
  description = "terraform test subnet"
}
resource "baiducloud_security_group" "sg" {
  name        = "terraform-sg"
  description = "security group created by terraform"
  vpc_id      = baiducloud_vpc.vpc.id
}
resource "baiducloud_security_group_rule" "sgr_in" {
  security_group_id = baiducloud_security_group.sg.id
  remark            = "remark"
  protocol          = "icmp"
  port_range        = ""
  direction         = "ingress"
}
resource "baiducloud_security_group_rule" "sgr_out" {
  security_group_id = baiducloud_security_group.sg.id
  remark            = "remark"
  protocol          = "all"
  port_range        = ""
  direction         = "egress"
  dest_ip           = "all"
}

resource "baiducloud_instance" "server1" {
  availability_zone = "cn-bj-d"
  instance_spec     = "bcc.gr1.c1m4"
  image_id          = data.baiducloud_images.images.images.0.id
  billing           = {
    payment_timing = "Postpaid"
  }
  admin_pass      = "Eni12345"
  subnet_id       = baiducloud_subnet.subnet.id
  security_groups = [
    baiducloud_security_group.sg.id
  ]
}
resource "baiducloud_eip" "eip1" {
  bandwidth_in_mbps = 1
  billing_method    = "ByBandwidth"
  payment_timing    = "Postpaid"
}
resource "baiducloud_eip" "eip2" {
  bandwidth_in_mbps = 1
  billing_method    = "ByBandwidth"
  payment_timing    = "Postpaid"
}

resource "baiducloud_eni" "eni" {
  name      = "terraform-eni"
  subnet_id = baiducloud_subnet.subnet.id

  description        = "terraform test"
  security_group_ids = [
    baiducloud_security_group.sg.id
  ]
  private_ip {
    primary            = true
    private_ip_address = "172.16.0.10"
    public_ip_address  = baiducloud_eip.eip2.eip
  }
  private_ip {
    primary            = false
    private_ip_address = "172.16.0.11"
    public_ip_address  = baiducloud_eip.eip1.eip
  }
  private_ip {
    primary            = false
    private_ip_address = "172.16.0.13"
    #    public_ip_address  = baiducloud_eip.eip2.eip
  }
}
resource "time_sleep" "wait_30_seconds" {
  depends_on      = [baiducloud_instance.server1, baiducloud_eni.eni]
  create_duration = "60s"
}
resource "baiducloud_eni_attachment" "default" {
  depends_on  = [time_sleep.wait_30_seconds]
  eni_id      = baiducloud_eni.eni.id
  instance_id = baiducloud_instance.server1.id
}
```

*/
package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/eni"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
	"time"
)

const (
	EniStatusInuse     = "inuse"
	EniStatusDetaching = "detaching"
	EniStatusAvailable = "available"
	EniStatusAttaching = "attaching"
)

func resourceBaiduCloudEniInstanceAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudEniAttachmentCreate,
		Read:   resourceBaiduCloudEniAttachmentRead,
		Delete: resourceBaiduCloudEniAttachmentDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Description: "Instance ID",
				Required:    true,
				ForceNew:    true,
			},
			"eni_id": {
				Type:        schema.TypeString,
				Description: "Eni ID",
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceBaiduCloudEniAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	action := "Create Eni Attachment"
	client := meta.(*connectivity.BaiduClient)
	eniService := EniService{
		client: client,
	}
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithEniClient(func(eniClient *eni.Client) (interface{}, error) {
			return nil, eniClient.AttachEniInstance(&eni.EniInstance{
				EniId:       d.Get("eni_id").(string),
				InstanceId:  d.Get("instance_id").(string),
				ClientToken: buildClientToken(),
			})
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
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_eni_attachment", action, BCESDKGoERROR)
	}
	stateConf := buildStateConf(
		[]string{EniStatusAvailable, EniStatusAttaching},
		[]string{EniStatusInuse},
		d.Timeout(schema.TimeoutCreate),
		eniService.eniStateRefresh(d.Get("eni_id").(string)),
	)
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_instance", action, BCESDKGoERROR)
	}
	return nil
}

func resourceBaiduCloudEniAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	action := "Query Eni Attachment"
	client := meta.(*connectivity.BaiduClient)

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithEniClient(func(eniClient *eni.Client) (interface{}, error) {
			return eniClient.GetEniDetail(d.Get("eni_id").(string))
		})
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, raw)
		res := raw.(*eni.Eni)
		d.Set("instance_id", res.InstanceId)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_eni_attachment", action, BCESDKGoERROR)
	}
	return nil
}

func resourceBaiduCloudEniAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	action := "Delete Eni Attachment"
	client := meta.(*connectivity.BaiduClient)
	eniService := &EniService{
		client: client,
	}
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithEniClient(func(eniClient *eni.Client) (interface{}, error) {
			return nil, eniClient.DetachEniInstance(&eni.EniInstance{
				EniId:       d.Get("eni_id").(string),
				InstanceId:  d.Get("instance_id").(string),
				ClientToken: buildClientToken(),
			})
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
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_eni_attachment", action, BCESDKGoERROR)
	}
	stateConf := buildStateConf(
		[]string{EniStatusInuse, EniStatusDetaching},
		[]string{EniStatusAvailable},
		d.Timeout(schema.TimeoutCreate),
		eniService.eniStateRefresh(d.Get("eni_id").(string)),
	)
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_eni_attachment", action, BCESDKGoERROR)
	}
	return nil
}
