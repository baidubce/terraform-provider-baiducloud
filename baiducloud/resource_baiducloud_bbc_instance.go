/*
Use this resource to create a BBC instance.

~> **NOTE:** The terminate operation of bcc does NOT take effect immediately，maybe takes for several minites.

Example Usage

```hcl
data "baiducloud_bbc_images" "bbc_images" {
  image_type = "BbcSystem"
  os_name    = "CentOS"
}
data "baiducloud_security_groups" "sg" {
  filter {
    name   = "name"
    values = ["default"]
  }
}
data "baiducloud_subnets" "subnets" {
  filter {
    name   = "zone_name"
    values = ["cn-bj-d"]
  }
  filter {
    name   = "name"
    values = ["系统预定义子网D"]
  }
}
data "baiducloud_bbc_flavors" "bbc_flavors" {
  filter {
    name   = "flavor_id"
    values = ["BBC-I4-01S"]
  }
}
resource "baiducloud_bbc_instance" "bbc_instance2" {
  action         = "start"
  payment_timing = "Postpaid"
  flavor_id            = "${data.baiducloud_bbc_flavors.bbc_flavors.flavors.0.flavor_id}"
  image_id             = "${data.baiducloud_bbc_images.bbc_images.images.0.id}"
  name                 = "terraform_test1"
  purchase_count       = 1
  raid                 = "Raid5"
  zone_name            = "cn-bj-d"
  root_disk_size_in_gb = 40
  security_groups      = [
    "${data.baiducloud_security_groups.sg.security_groups.0.id}",
    "${data.baiducloud_security_groups.sg.security_groups.1.id}",
  ]
  tags = {
    "testKey" = "terraform_test"
  }
  description = "terraform_test"
}
```

Import

BBC instance can be imported, e.g.

```hcl
$ terraform import baiducloud_bbc_instance.my-server id
```
*/
package baiducloud

import (
	"encoding/json"
	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/bbc"
	"github.com/baidubce/bce-sdk-go/services/bcc"
	"github.com/baidubce/bce-sdk-go/services/bcc/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
	"log"
	"strconv"
	"time"
)

func resourceBaiduCloudBccInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudBccInstanceCreate,
		Read:   resourceBaiduCloudBbcInstanceRead,
		Update: resourceBaiduCloudBbcInstanceUpdate,
		Delete: resourceBaiduCloudBbcInstanceDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(40 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "BBC name.Support for uppercase and lowercase letters, numbers, Chinese and special characters, such as \"-\",\"_\",\"/\",\".\", the value must start with a letter, length 1-65.",
				Required:    true,
			},
			"flavor_id": {
				Type:        schema.TypeString,
				Description: "Id of the BBC Flavor.",
				Required:    true,
			},
			"image_id": {
				Type:        schema.TypeString,
				Description: "Id of the BBC Image.",
				Required:    true,
			},
			"raid_id": {
				Type:        schema.TypeString,
				Description: "Id of the raid.",
				Computed:    true,
			},
			"raid": {
				Type:         schema.TypeString,
				Description:  "Type of the raid to start. Available values are Raid5, NoRaid.",
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"Raid5", "NoRaid"}, false),
			},
			"root_disk_size_in_gb": {
				Type:         schema.TypeInt,
				Description:  "The system disk size of the BBC instance to be created.",
				Required:     true,
				ValidateFunc: validation.IntBetween(40, 500),
			},
			"purchase_count": {
				Type:         schema.TypeInt,
				Description:  "The number of BBC instances created (purchased) in batch. It must be an integer greater than 0. It is an optional parameter. The default value is 1.",
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 2),
			},
			"zone_name": {
				Type:        schema.TypeString,
				Description: "The naming convention of zonename is \"country-region-availability area\", in lowercase, for example, Beijing availability area A is \"cn-bj-a\"“",
				Required:    true,
				ForceNew:    true,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Description: "Id of bbc subnet.",
				Optional:    true,
			},
			"auto_renew_time_unit": {
				Type:         schema.TypeString,
				Description:  "Monthly payment or annual payment, month is \"month\" and year is \"year\".",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"month", "year"}, false),
			},
			"auto_renew_time": {
				Type:         schema.TypeInt,
				Description:  "The automatic renewal time is 1-9 per month and 1-3 per year.",
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 9),
			},
			"payment_timing": {
				Type:         schema.TypeString,
				Description:  "Payment timing of billing, which can be Prepaid or Postpaid. The default is Postpaid.",
				Optional:     true,
				Default:      bbc.PaymentTimingPostPaid,
				ValidateFunc: validatePaymentTiming(),
			},
			"reservation": {
				Type:             schema.TypeMap,
				Description:      "Reservation of the instance.",
				Optional:         true,
				DiffSuppressFunc: postPaidDiffSuppressFunc,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"reservation_length": {
							Type:             schema.TypeInt,
							Description:      "The reservation length that you will pay for your resource. It is valid when payment_timing is Prepaid. Valid values: [1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36].",
							Required:         true,
							Default:          1,
							ValidateFunc:     validateReservationLength(),
							DiffSuppressFunc: postPaidDiffSuppressFunc,
						},
						"reservation_time_unit": {
							Type:             schema.TypeString,
							Description:      "The reservation time unit that you will pay for your resource. It is valid when payment_timing is Prepaid. The value can only be month currently, which is also the default value.",
							Required:         true,
							Default:          "Month",
							ValidateFunc:     validateReservationUnit(),
							DiffSuppressFunc: postPaidDiffSuppressFunc,
						},
					},
				},
			},
			"hostname": {
				Type:        schema.TypeString,
				Description: "Hostname is not specified by default. Hostname only supports lowercase letters, numbers and -. Special characters. It must start with a letter. Special symbols cannot be used consecutively. It does not support starting or ending with special symbols. The length is 2-64.",
				Optional:    true,
			},
			"admin_pass": {
				Type:        schema.TypeString,
				Description: "admin password.",
				Optional:    true,
			},
			"deploy_set_id": {
				Type:        schema.TypeString,
				Description: "deploy set of bbc.",
				Optional:    true,
			},
			"security_groups": {
				Type:        schema.TypeSet,
				Description: "Security groups of the bbc instance.It can use",
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tags": tagsSchema(),
			"request_token": {
				Type:        schema.TypeString,
				Description: "request_token.",
				Optional:    true,
			},
			"enable_numa": {
				Type:        schema.TypeBool,
				Description: "enableNuma.",
				Optional:    true,
			},
			"enable_ht": {
				Type:        schema.TypeBool,
				Description: "enableHt",
				Optional:    true,
				Default:     true,
			},
			"root_partition_type": {
				Type:        schema.TypeString,
				Description: "namroot_partition_type.",
				Optional:    true,
			},
			"data_partition_type": {
				Type:        schema.TypeString,
				Description: "data_partition_type.",
				Optional:    true,
			},
			"host_name": {
				Type:        schema.TypeString,
				Description: "host_name.",
				Computed:    true,
			},
			"uuid": {
				Type:        schema.TypeString,
				Description: "uuid.",
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Status of the bbc instance.",
				Computed:    true,
			},
			"create_time": {
				Type:        schema.TypeString,
				Description: "BBC create time.",
				Computed:    true,
			},
			"expire_time": {
				Type:        schema.TypeString,
				Description: "expire time.",
				Computed:    true,
			},
			"public_ip": {
				Type:        schema.TypeString,
				Description: "public ip.",
				Computed:    true,
			},
			"internal_ip": {
				Type:        schema.TypeString,
				Description: "internal ip.",
				Computed:    true,
			},
			"rdma_ip": {
				Type:        schema.TypeString,
				Description: "rdma ip.",
				Computed:    true,
			},
			"region": {
				Type:        schema.TypeString,
				Description: "region.",
				Computed:    true,
			},
			"has_alive": {
				Type:        schema.TypeInt,
				Description: "hasAlive.",
				Computed:    true,
			},
			"switch_id": {
				Type:        schema.TypeString,
				Description: "switch id.",
				Computed:    true,
			},
			"host_id": {
				Type:        schema.TypeString,
				Description: "switch id.",
				Computed:    true,
			},
			"network_capacity_in_mbps": {
				Type:        schema.TypeString,
				Description: "network capacity in mbps.",
				Computed:    true,
			},
			"rack_id": {
				Type:        schema.TypeString,
				Description: "rack id.",
				Computed:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "description.",
				Optional:    true,
			},
			"action": {
				Type:         schema.TypeString,
				Description:  "action.",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"start", "stop"}, false),
			},
		},
	}
}

func resourceBaiduCloudBccInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	bbcService := BbcService{client}
	action := "Create bbc"

	securityGroups := make([]interface{}, 0)
	groups, ok := d.GetOk("security_groups")
	if ok {
		securityGroups = groups.(*schema.Set).List()
	}
	createInstanceArgs, err := buildBaiduCloudBbcInstanceArgs(d, meta)
	// init security group id, create only bind the first security group
	if len(securityGroups) > 0 {
		createInstanceArgs.SecurityGroupId = securityGroups[0].(string)
	}
	jsonData,_ := json.Marshal(createInstanceArgs)
	log.Print("BBC args is ", string(jsonData))
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		res, err := client.WithBbcClient(func(bbcClient *bbc.Client) (i interface{}, e error) {
			return bbcClient.CreateInstance(createInstanceArgs)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		d.SetId(res.(*bbc.CreateInstanceResult).InstanceIds[0])
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_instance", action, BCESDKGoERROR)
	}
	stateConf := buildStateConf(
		[]string{string(bbc.InstanceStatusStarting)},
		[]string{string(bbc.InstanceStatusRunning), string(bbc.InstanceStatusDeleted)},
		d.Timeout(schema.TimeoutCreate),
		bbcService.InstanceStateRefresh(d.Id()),
	)
	if _, err = stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_instance", action, BCESDKGoERROR)
	}
	// bind security groups args build
	if err := bbcService.updateBbcInstanceSecurityGroups(d, meta, d.Id()); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_instance", action, BCESDKGoERROR)
	}
	// update description
	if err := bbcService.updateBccInstanceDescription(d, meta, d.Id()); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_instance", action, BCESDKGoERROR)
	}
	// if action is stop,stop the instance
	if d.Get("action").(string) == INSTANCE_ACTION_STOP {
		if err := bbcService.StopBbcInstance(d.Id(), d.Timeout(schema.TimeoutUpdate)); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_instance", action, BCESDKGoERROR)
		}
	}

	return resourceBaiduCloudBbcInstanceRead(d, meta)
}
func resourceBaiduCloudBbcInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	instanceID := d.Id()
	action := "Query BBC Instance " + instanceID

	raw, err := client.WithBbcClient(func(bbcClient *bbc.Client) (interface{}, error) {
		return bbcClient.GetInstanceDetail(instanceID)
	})
	addDebug(action, raw)

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_instance", action, BCESDKGoERROR)
	}
	instance, _ := raw.(*bbc.InstanceModel)

	// write to result
	d.Set("name", instance.Name)
	d.Set("host_name", instance.Hostname)
	d.Set("uuid", instance.Uuid)
	d.Set("description", instance.Desc)
	d.Set("status", instance.Status)
	d.Set("create_time", instance.CreateTime)
	d.Set("expire_time", instance.ExpireTime)
	d.Set("public_ip", instance.PublicIp)
	d.Set("internal_ip", instance.InternalIp)
	d.Set("rdma_ip", instance.RdmaIp)
	d.Set("image_id", instance.ImageId)
	d.Set("flavor_id", instance.FlavorId)
	d.Set("zone_name", instance.Zone)
	d.Set("region", instance.Region)
	d.Set("has_alive", instance.HasAlive)
	d.Set("tags", flattenTagsToMap(instance.Tags))
	d.Set("switch_id", instance.SwitchId)
	d.Set("host_id", instance.HostId)
	d.Set("network_capacity_in_mbps", instance.NetworkCapacityInMbps)
	d.Set("rack_id", instance.RackId)
	d.Set("payment_timing", instance.PaymentTiming)

	// security groups
	raw, err = client.WithBccClient(func(bccClient *bcc.Client) (interface{}, error) {
		args := &api.ListSecurityGroupArgs{
			InstanceId: instanceID,
		}
		return bccClient.ListSecurityGroup(args)
	})
	addDebug(action, raw)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_instance", action, BCESDKGoERROR)
	}

	securityGroups, _ := raw.(*api.ListSecurityGroupResult)
	sgIDs := make([]string, len(securityGroups.SecurityGroups))
	for i, sg := range securityGroups.SecurityGroups {
		sgIDs[i] = sg.Id
	}
	addDebug(action, sgIDs)
	d.Set("security_groups", sgIDs)
	return nil
}
func resourceBaiduCloudBbcInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	bbcService := BbcService{client}

	instanceID := d.Id()
	action := "bbc instance update"
	d.Partial(true)

	// update bbc instance attribute
	if err := bbcService.updateBccInstanceDescription(d, meta, instanceID); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_instance", action, BCESDKGoERROR)
	}
	// update bbc name
	if err := bbcService.updateBbcInstanceName(d, meta, instanceID); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_instance", action, BCESDKGoERROR)
	}
	// update bbc instance security groups
	if err := bbcService.updateBbcInstanceSecurityGroups(d, meta, instanceID); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_instance", action, BCESDKGoERROR)
	}
	// update bbc instance admin pass
	if err := bbcService.updateBbcInstanceAdminPass(d, meta, instanceID); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_instance", action, BCESDKGoERROR)
	}
	// update bbc instance action
	if err := bbcService.updateBbcInstanceAction(d, meta, instanceID); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_instance", action, BCESDKGoERROR)
	}

	d.Partial(false)

	return resourceBaiduCloudBbcInstanceRead(d, meta)
}
func resourceBaiduCloudBbcInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	bbcService := BbcService{client}

	instanceId := d.Id()
	action := "Delete BBC Instance " + instanceId

	// delete bbc instance
	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		raw, err := client.WithBbcClient(func(bbcClient *bbc.Client) (interface{}, error) {
			return instanceId, bbcClient.DeleteInstance(instanceId)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{ReleaseWhileCreating, bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, raw)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_instance", action, BCESDKGoERROR)
	}

	stateConf := buildStateConf(
		[]string{string(bbc.InstanceStatusStopping), string(bbc.InstanceStatusStopped)},
		[]string{string(bbc.InstanceStatusDeleted)},
		d.Timeout(schema.TimeoutDelete),
		bbcService.InstanceStateRefresh(instanceId),
	)
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_instance", action, BCESDKGoERROR)
	}

	return nil
}
func buildBaiduCloudBbcInstanceArgs(d *schema.ResourceData, meta interface{}) (*bbc.CreateInstanceArgs, error) {
	client := meta.(*connectivity.BaiduClient)
	bbcService := BbcService{client}

	request := &bbc.CreateInstanceArgs{
		ClientToken: buildClientToken(),
	}
	if flavorId, ok := d.GetOk("flavor_id"); ok {
		request.FlavorId = flavorId.(string)
	}
	if imageId, ok := d.GetOk("image_id"); ok {
		request.ImageId = imageId.(string)
	}
	if raid, ok := d.GetOk("raid"); ok {
		raidId, err := bbcService.getRaidIdByFlavor(request.FlavorId, raid.(string))
		if err != nil {
			return nil, WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc", "no such raid", BCESDKGoERROR)
		}
		request.RaidId = raidId
	}
	if rootDiskSizeInGb, ok := d.GetOk("root_disk_size_in_gb"); ok {
		request.RootDiskSizeInGb = rootDiskSizeInGb.(int)
	}
	if purchaseCount, ok := d.GetOk("purchase_count"); ok {
		request.PurchaseCount = purchaseCount.(int)
	}
	if zoneName, ok := d.GetOk("zone_name"); ok {
		request.ZoneName = zoneName.(string)
	}
	// build billing
	billingRequest := bbc.Billing{
		PaymentTiming: bbc.PaymentTimingType(""),
		Reservation:   bbc.Reservation{},
	}
	if p, ok := d.GetOk("payment_timing"); ok {
		paymentTiming := bbc.PaymentTimingType(p.(string))
		billingRequest.PaymentTiming = paymentTiming
	}
	if billingRequest.PaymentTiming == bbc.PaymentTimingPrePaid {
		if r, ok := d.GetOk("reservation"); ok {
			reservation := r.(map[string]interface{})
			if reservationLength, ok := reservation["reservation_length"]; ok {
				reservationLengthInt, err := strconv.Atoi(reservationLength.(string))
				billingRequest.Reservation.Length = reservationLengthInt
				if err != nil {
					return nil, err
				}
			}
			if reservationTimeUnit, ok := reservation["reservation_time_unit"]; ok {
				billingRequest.Reservation.TimeUnit = reservationTimeUnit.(string)
			}
		}
		// if the field is set, then auto-renewal is effective.
		if v, ok := d.GetOk("auto_renew_time_unit"); ok {
			request.AutoRenewTimeUnit = v.(string)
			if v, ok := d.GetOk("auto_renew_time"); ok {
				request.AutoRenewTime = v.(int)
			}
		}
	}
	request.Billing = billingRequest
	if subnetId, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = subnetId.(string)
	}
	if name, ok := d.GetOk("name"); ok {
		request.Name = name.(string)
	}
	if hostname, ok := d.GetOk("hostname"); ok {
		request.Hostname = hostname.(string)
	}
	if adminPass, ok := d.GetOk("admin_pass"); ok {
		request.AdminPass = adminPass.(string)
	}
	if deploySetId, ok := d.GetOk("deploy_set_id"); ok {
		request.DeploySetId = deploySetId.(string)
	}
	if clientToken, ok := d.GetOk("client_token"); ok {
		request.ClientToken = clientToken.(string)
	}
	if tags, ok := d.GetOk("tags"); ok {
		request.Tags = tranceTagMapToModel(tags.(map[string]interface{}))
	}
	if internalIps, ok := d.GetOk("internal_ips"); ok {
		ips := make([]string, 0)
		for _, ip := range internalIps.(*schema.Set).List() {
			ips = append(ips, ip.(string))
		}
		request.InternalIps = ips
	}
	if requestToken, ok := d.GetOk("request_token"); ok {
		request.RequestToken = requestToken.(string)
	}
	if enableNuma, ok := d.GetOk("enable_numa"); ok {
		request.EnableNuma = enableNuma.(bool)
	}
	if rootPartitionType, ok := d.GetOk("root_partition_type"); ok {
		request.RootPartitionType = rootPartitionType.(string)
	}
	if dataPartitionType, ok := d.GetOk("data_partition_type"); ok {
		request.DataPartitionType = dataPartitionType.(string)
	}
	if v, ok := d.GetOk("tags"); ok {
		request.Tags = tranceTagMapToModel(v.(map[string]interface{}))
	}
	if v, ok := d.GetOk("enable_ht"); ok {
		request.EnableHt = v.(bool)
	}
	return request, nil
}
