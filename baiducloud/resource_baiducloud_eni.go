/*
Provide a resource to create an ENI.

Example Usage

```hcl
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
    public_ip_address  = baiducloud_eip.eip1.eip
  }
  private_ip {
    primary            = false
    private_ip_address = "172.16.0.11"
    public_ip_address  = baiducloud_eip.eip2.eip
  }
  private_ip {
    primary            = false
    private_ip_address = "172.16.0.13"
  }
}
```

Import

ENI can be imported, e.g.

```hcl
$ terraform import baiducloud_eni.default eni_id
```
*/
package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/eni"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
	"log"
	"time"
)

const EniAvailable = "available"

func resourceBaiduCloudEni() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudEniCreate,
		Read:   resourceBaiduCloudEniRead,
		Update: resourceBaiduCloudEniUpdate,
		Delete: resourceBaiduCloudEniDelete,

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
				Type: schema.TypeString,
				Description: "Name of the ENI. Support for uppercase and lowercase letters, " +
					"numbers, Chinese and special characters, " +
					"such as \"-\",\"_\",\"/\",\".\", the value must start with a letter, length 1-65.",
				Required: true,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Description: "Subnet ID which ENI belong to",
				Required:    true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Description: "Instance ID the ENI bind",
				Computed:    true,
			},
			"security_group_ids": {
				Type:        schema.TypeList,
				Description: "Specifies the set of bound security group IDs",
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"enterprise_security_group_ids": {
				Type:        schema.TypeList,
				Description: "Specifies the set of bound enterprise security group IDs",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"private_ip": {
				Type:        schema.TypeList,
				Description: "Specified intranet IP information",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"public_ip_address": {
							Type:        schema.TypeString,
							Description: "The public IP address of the ENI, that is, the eip address",
							Optional:    true,
						},
						"primary": {
							Type:        schema.TypeBool,
							Description: "True or false, true mean it is primary IP, it's private IP address can not modify, only one primary IP in a ENI",
							Required:    true,
						},
						"private_ip_address": {
							Type:        schema.TypeString,
							Description: "Intranet IP address of the ENI",
							Required:    true,
						},
					},
				},
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description of the ENI",
				Optional:    true,
			},
			"zone_name": {
				Type:        schema.TypeString,
				Description: "Availability zone name which ENI belong to",
				Computed:    true,
			},
			"mac_address": {
				Type:        schema.TypeString,
				Description: "Mac address of the ENI",
				Computed:    true,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Description: "VPC id which the ENI belong to",
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Status of ENI, may be inuse, binding, unbinding, available",
				Computed:    true,
			},
			"created_time": {
				Type:        schema.TypeString,
				Description: "ENI create time",
				Computed:    true,
			},
		},
	}
}

func resourceBaiduCloudEniCreate(d *schema.ResourceData, meta interface{}) error {
	action := "Create Eni"
	client := meta.(*connectivity.BaiduClient)
	args := buildEniCreateArgs(d)
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithEniClient(func(eniClient *eni.Client) (interface{}, error) {
			return eniClient.CreateEni(args)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, raw)
		res := raw.(*eni.CreateEniResult)
		d.SetId(res.EniId)
		err = updateEniPrivateIP(d, meta)
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_eni", action, BCESDKGoERROR)
	}
	return resourceBaiduCloudEniRead(d, meta)
}

func resourceBaiduCloudEniRead(d *schema.ResourceData, meta interface{}) error {
	action := "Query Eni Detail"
	client := meta.(*connectivity.BaiduClient)
	err := resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
		raw, err := client.WithEniClient(func(eniClient *eni.Client) (interface{}, error) {
			return eniClient.GetEniDetail(d.Id())
		})
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, raw)
		eniDetail := raw.(*eni.Eni)
		d.Set("eni_id", eniDetail.EniId)
		d.Set("name", eniDetail.Name)
		d.Set("zone_name", eniDetail.ZoneName)
		d.Set("description", eniDetail.Description)
		d.Set("instance_id", eniDetail.InstanceId)
		d.Set("mac_address", eniDetail.MacAddress)
		d.Set("vpc_id", eniDetail.VpcId)
		d.Set("subnet_id", eniDetail.SubnetId)
		d.Set("status", eniDetail.Status)
		d.Set("security_group_ids", eniDetail.SecurityGroupIds)
		d.Set("enterprise_security_group_ids", eniDetail.EnterpriseSecurityGroupIds)
		d.Set("created_time", eniDetail.CreatedTime)
		privateIps := make([]map[string]interface{}, 0)
		for _, item := range eniDetail.PrivateIpSet {
			privateIps = append(privateIps, map[string]interface{}{
				"primary":            item.Primary,
				"private_ip_address": item.PrivateIpAddress,
				"public_ip_address":  item.PublicIpAddress,
			})
		}
		d.Set("private_ip", privateIps)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_eni", action, BCESDKGoERROR)
	}
	return nil
}

func resourceBaiduCloudEniUpdate(d *schema.ResourceData, meta interface{}) error {
	if err := updateNameAndDescription(d, meta); err != nil {
		return err
	}

	if err := updateEniPrivateIP(d, meta); err != nil {
		return err
	}

	if err := updateEniSecurityGroup(d, meta); err != nil {
		return err
	}

	return nil
}

func resourceBaiduCloudEniDelete(d *schema.ResourceData, meta interface{}) error {
	action := "Delete Eni"
	client := meta.(*connectivity.BaiduClient)
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithEniClient(func(eniClient *eni.Client) (interface{}, error) {
			return nil, eniClient.DeleteEni(&eni.DeleteEniArgs{
				ClientToken: buildClientToken(),
				EniId:       d.Id(),
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
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_eni", action, BCESDKGoERROR)
	}
	return nil
}

func buildEniCreateArgs(d *schema.ResourceData) *eni.CreateEniArgs {
	res := &eni.CreateEniArgs{
		ClientToken: buildClientToken(),
	}
	if v, ok := d.GetOk("name"); ok {
		res.Name = v.(string)
	}
	if v, ok := d.GetOk("subnet_id"); ok {
		res.SubnetId = v.(string)
	}
	if v, ok := d.GetOk("description"); ok {
		res.Description = v.(string)
	}
	if v, ok := d.GetOk("security_group_ids"); ok {
		res.SecurityGroupIds = interfaceSlice2StringSlice(v.([]interface{}))
	}
	if v, ok := d.GetOk("enterprise_security_group_ids"); ok {
		res.EnterpriseSecurityGroupIds = v.([]string)
	}
	if v, ok := d.GetOk("private_ip"); ok {
		res.PrivateIpSet = interfaceSlice2PrivateIpSlice(v.([]interface{}))
	}
	return res
}

func interfaceSlice2StringSlice(v []interface{}) []string {
	res := make([]string, 0)
	for _, i := range v {
		res = append(res, i.(string))
	}
	return res
}

func interfaceSlice2PrivateIpSlice(v []interface{}) []eni.PrivateIp {
	res := make([]eni.PrivateIp, 0)
	for _, i := range v {
		ipMap := i.(map[string]interface{})
		item := eni.PrivateIp{}
		if ipMap["primary"] != nil {
			item.Primary = ipMap["primary"].(bool)
		}
		if ipMap["private_ip_address"] != nil {
			item.PrivateIpAddress = ipMap["private_ip_address"].(string)
		}
		if ipMap["public_ip_address"] != nil {
			item.PublicIpAddress = ipMap["public_ip_address"].(string)
		}
		res = append(res, item)
	}
	return res
}

func getIPAddress(key string, data []interface{}) []string {
	res := make([]string, 0)
	for _, datum := range data {
		temp := datum.(map[string]interface{})
		res = append(res, temp[key].(string))
	}
	return res
}

func updateNameAndDescription(d *schema.ResourceData, meta interface{}) error {
	action := "Update Eni Name And Description"
	client := meta.(*connectivity.BaiduClient)

	if d.HasChanges("name", "description") {
		err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			raw, err := client.WithEniClient(func(eniClient *eni.Client) (interface{}, error) {
				return nil, eniClient.UpdateEni(&eni.UpdateEniArgs{
					EniId:       d.Id(),
					ClientToken: buildClientToken(),
					Name:        d.Get("name").(string),
					Description: d.Get("description").(string),
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
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_eni", action, BCESDKGoERROR)
		}
	}
	return nil
}

func updateEniPrivateIP(d *schema.ResourceData, meta interface{}) error {
	action := "Update Eni Private IPs"
	client := meta.(*connectivity.BaiduClient)
	eipService := &EipService{
		client: client,
	}
	if d.HasChange("private_ip") {
		unbindIps := make([]string, 0)
		bindIps := make([]string, 0)
		o, n := d.GetChange("private_ip")
		os := o.([]interface{})
		ns := n.([]interface{})
		osSlice := getIPAddress("private_ip_address", os)
		nsSlice := getIPAddress("private_ip_address", ns)
		// 1.unbind EIP
		for _, item := range os {
			temp := item.(map[string]interface{})
			if temp["public_ip_address"].(string) != "" {
				_, err := client.WithEniClient(func(eniClient *eni.Client) (interface{}, error) {
					// unbind
					return nil, eniClient.UnBindEniPublicIp(&eni.UnBindEniPublicIpArgs{
						EniId:           d.Id(),
						ClientToken:     buildClientToken(),
						PublicIpAddress: temp["public_ip_address"].(string),
					})
				})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, "baiducloud_eni", action, BCESDKGoERROR)
				}
			}
		}
		// 2.compute the private IP which need unbind
		for _, item := range os {
			temp := item.(map[string]interface{})
			if !stringInSlice(nsSlice, temp["private_ip_address"].(string)) {
				unbindIps = append(unbindIps, temp["private_ip_address"].(string))
			}
		}
		// 3.compute the private IP which need bind
		for _, item := range ns {
			temp := item.(map[string]interface{})
			if !stringInSlice(osSlice, temp["private_ip_address"].(string)) {
				bindIps = append(bindIps, temp["private_ip_address"].(string))
			}
		}
		// 4.bind and unbind the private IP by computed result
		err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			raw, err := client.WithEniClient(func(eniClient *eni.Client) (interface{}, error) {
				// unbind
				if len(unbindIps) != 0 {
					err := eniClient.BatchDeletePrivateIp(&eni.EniBatchPrivateIpArgs{
						EniId:                 d.Id(),
						ClientToken:           buildClientToken(),
						PrivateIpAddresses:    unbindIps,
						PrivateIpAddressCount: len(unbindIps),
					})
					if err != nil {
						return nil, err
					}
				}
				if len(bindIps) != 0 {
					// bind
					_, err := eniClient.BatchAddPrivateIp(&eni.EniBatchPrivateIpArgs{
						EniId:              d.Id(),
						ClientToken:        buildClientToken(),
						PrivateIpAddresses: bindIps,
					})
					return nil, err
				}
				return nil, nil
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
		// 5.bind EIP
		for _, item := range ns {
			temp := item.(map[string]interface{})
			ip := temp["public_ip_address"].(string)
			if ip == "" {
				continue
			}
			for {
				res, err := eipService.EipGetDetail(ip)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, "baiducloud_eni", action, BCESDKGoERROR)
				}
				log.Print(res.Status)
				if res.Status != EniAvailable {
					time.Sleep(1 * time.Second)
					continue
				}
				_, updateErr := client.WithEniClient(func(eniClient *eni.Client) (interface{}, error) {
					// bind
					return nil, eniClient.BindEniPublicIp(&eni.BindEniPublicIpArgs{
						EniId:            d.Id(),
						ClientToken:      buildClientToken(),
						PrivateIpAddress: temp["private_ip_address"].(string),
						PublicIpAddress:  temp["public_ip_address"].(string),
					})
				})
				if updateErr != nil {
					return WrapErrorf(updateErr, DefaultErrorMsg, "baiducloud_eni", action, BCESDKGoERROR)
				}
				break
			}
		}
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_eni", action, BCESDKGoERROR)
		}
	}
	return nil
}

func updateEniSecurityGroup(d *schema.ResourceData, meta interface{}) error {
	action := "Update Eni Security Group"
	client := meta.(*connectivity.BaiduClient)
	if d.HasChange("security_group_ids") {
		sgs := interfaceSlice2StringSlice(d.Get("security_group_ids").([]interface{}))
		if len(sgs) == 0 {
			return WrapErrorf(Error("security group ids can not be nil"), DefaultErrorMsg, "baiducloud_eni", action, BCESDKGoERROR)
		}
		err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			raw, err := client.WithEniClient(func(eniClient *eni.Client) (interface{}, error) {
				return nil, eniClient.UpdateEniSecurityGroup(&eni.UpdateEniSecurityGroupArgs{
					EniId:            d.Id(),
					ClientToken:      buildClientToken(),
					SecurityGroupIds: sgs,
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
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_eni", action, BCESDKGoERROR)
		}
	}
	return nil
}
