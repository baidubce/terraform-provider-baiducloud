/*
Use this resource to get information about a BBC instance.

~> **NOTE:** The terminate operation of bbc does NOT take effect immediately，maybe takes for several minites.

Example Usage

```hcl
resource "baiducloud_bbc_instance" "my-server" {
  image_id = "m-A4jJpFzi"
  hostname = "hostname"
  name = "my-instance"
  raid_id = "raidId"
  subnet_id = ""
  security_group = ""
  availability_zone = "cn-bj-a"
  billing = {
    payment_timing = "Postpaid"
  }
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
	"time"

	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/bbc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func resourceBaiduCloudBbcInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudBbcInstanceCreate,
		Read:   resourceBaiduCloudBbcInstanceRead,
		Update: resourceBaiduCloudBbcInstanceUpdate,
		Delete: resourceBaiduCloudBbcInstanceDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"image_id": {
				Type:        schema.TypeString,
				Description: "ID of the image to be used for the instance.",
				Required:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the instance. Support for uppercase and lowercase letters, numbers, Chinese and special characters, such as \"-\",\"_\",\"/\",\".\", the value must start with a letter, length 1-65.",
				Optional:    true,
			},
			"hostname": {
				Type:        schema.TypeString,
				Description: "Host Name of the instance. Support for uppercase and lowercase letters, numbers, Chinese and special characters, such as \"-\",\"_\",\"/\",\".\", the value must start with a letter, length 1-65.",
				Optional:    true,
			},
			"flavor_id": {
				Type:        schema.TypeString,
				Description: " Physical machine package ID",
				Required:    true,
			},
			"raid_id": {
				Type:        schema.TypeString,
				Description: "raid configration id",
				Required:    true,
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Description: "Availability zone to start the instance in.",
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
			},
			"billing": {
				Type:        schema.TypeMap,
				Description: "Billing information of the instance.",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"payment_timing": {
							Type:         schema.TypeString,
							Description:  "Payment timing of billing, which can be Prepaid or Postpaid. The default is Postpaid.",
							Required:     true,
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
					},
				},
			},
			"cds_disks": {
				Type:        schema.TypeList,
				Description: "CDS disks of the instance.",
				Computed:    true,
				MinItems:    1,
				MaxItems:    10,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_size_in_gb": {
							Type:         schema.TypeInt,
							Description:  "The size(GB) of CDS.",
							Optional:     true,
							Default:      0,
							ValidateFunc: validation.IntAtLeast(0),
						},
						"storage_type": {
							Type:         schema.TypeString,
							Description:  "Storage type of the CDS.",
							Optional:     true,
							Default:      bbc.StorageTypeCloudHP1,
							ValidateFunc: validateStorageType(),
						},
						"is_system_volume": {
							Type:        schema.TypeBool,
							Description: "Snapshot ID of CDS.",
							Optional:    true,
						},
					},
				},
			},
			"admin_pass": {
				Type:        schema.TypeString,
				Description: "Password of the instance to be started. This value should be 8-16 characters, and English, numbers and symbols must exist at the same time. The symbols is limited to \"!@#$%^*()\".",
				Optional:    true,
				Sensitive:   true,
			},
			"cpu_count": {
				Type:        schema.TypeInt,
				Description: "Number of CPU cores to be created for the instance.",
				Computed:    true,
			},
			"memory_capacity_in_gb": {
				Type:        schema.TypeInt,
				Description: "Memory capacity(GB) of the instance to be created.",
				Computed:    true,
			},
			"root_disk_size_in_gb": {
				Type:         schema.TypeInt,
				Description:  "System disk size(GB) of the instance to be created. The value range is [20,100]GB, Default to 20GB, and more than 20GB is charged according to the cloud disk price. Note that the specified system disk size needs to meet the minimum disk space limit of the mirror used.",
				Optional:     true,
				ForceNew:     true,
				Default:      20,
				ValidateFunc: validation.IntBetween(20, 100),
			},
			"public_ip": {
				Type:        schema.TypeString,
				Description: "Public IP",
				Computed:    true,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Description: "The subnet ID of VPC. The default subnet will be used when it is empty. The instance will restart after changing the subnet.",
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
			},
			"security_group": {
				Type:        schema.TypeString,
				Description: "Security groups of the instance.",
				Optional:    true,
				Computed:    true,
			},
			"auto_renew_time_unit": {
				Type:         schema.TypeString,
				Description:  "Time unit of automatic renewal, the value can be month or year. The default value is empty, indicating no automatic renewal. It is valid only when the payment_timing is Prepaid.",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"month", "year"}, false),
			},
			"auto_renew_time_length": {
				Type:         schema.TypeInt,
				Description:  "The time length of automatic renewal. It is valid when payment_timing is Prepaid, and the value should be 1-9 when the auto_renew_time_unit is month and 1-3 when the auto_renew_time_unit is year. Default to 1.",
				Optional:     true,
				ForceNew:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 9),
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description of the instance.",
				Optional:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Status of the instance.",
				Computed:    true,
			},
			"create_time": {
				Type:        schema.TypeString,
				Description: "Create time of the instance.",
				Computed:    true,
			},
			"expire_time": {
				Type:        schema.TypeString,
				Description: "Expire time of the instance.",
				Computed:    true,
			},
			"internal_ip": {
				Type:        schema.TypeString,
				Description: "Internal IP assigned to the instance.",
				Computed:    true,
			},
			"network_capacity_in_mbps": {
				Type:        schema.TypeString,
				Description: "The placement policy of the instance, which can be default or dedicatedHost.",
				Computed:    true,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Description: "VPC ID of the instance.",
				Computed:    true,
			},
			"tags": normalTagsSchema(),
		},
	}
}

func resourceBaiduCloudBbcInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	bbcService := BbcService{client}

	securityGroups := ""
	groups, ok := d.GetOk("security_group")
	if ok {
		securityGroups = groups.(string)
	}

	var err error

	createInstanceArgs, err := buildBaiduCloudBbcInstanceArgs(d, meta)
	if err != nil {
		return WrapError(err)
	}

	if len(securityGroups) > 0 {
		createInstanceArgs.SecurityGroupId = securityGroups
	}

	action := "Create Bbc Instance"
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithBbcClient(func(bbcClient *bbc.Client) (interface{}, error) {
			return bbcClient.CreateInstance(createInstanceArgs)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, raw)
		response, _ := raw.(*bbc.CreateInstanceResult)
		d.SetId(response.InstanceIds[0])
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_instance", action, BCESDKGoERROR)
	}

	stateConf := buildStateConf(
		[]string{string(bbc.InstanceStatusStarting)},
		[]string{string(bbc.InstanceStatusRunning)},
		d.Timeout(schema.TimeoutCreate),
		bbcService.InstanceBbcStateRefresh(d.Id()),
	)
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_instance", action, BCESDKGoERROR)
	}

	// set instance description
	if err := updateInstanceDescription(d, meta, d.Id()); err != nil {
		return err
	}

	return resourceBaiduCloudBbcInstanceRead(d, meta)
}

func resourceBaiduCloudBbcInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	bbcService := BbcService{client}

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
	response, _ := raw.(*bbc.InstanceModel)

	// Required or Optional
	d.Set("image_id", response.ImageId)
	d.Set("name", response.Name)
	d.Set("availability_zone", response.Zone)
	//TODO: unsupport to import
	// d.Set("raid_id", string(response.RdmaIp))
	//if res, err := bbcService.client.WithBbcClient(func(client *bbc.Client) (interface{}, error){
	//	return client.GetFlavorRaid(response.FlavorId)
	//})

	d.Set("flavor_id", response.FlavorId)
	flavor, err := bbcService.GetFlavorDetail(response.FlavorId)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_instance", action, BCESDKGoERROR)
	}

	d.Set("cpu_count", flavor.CpuCount)
	d.Set("memory_capacity_in_gb", flavor.MemoryCapacityInGB)

	d.Set("tags", flattenTagsToMap(response.Tags))

	billingMap := map[string]interface{}{"payment_timing": response.PaymentTiming}
	d.Set("billing", billingMap)

	// Computed
	d.Set("description", response.Desc)
	d.Set("status", response.Status)
	d.Set("create_time", response.CreateTime)
	d.Set("expire_time", response.ExpireTime)
	d.Set("public_ip", response.PublicIp)
	d.Set("internal_ip", response.InternalIp)

	d.Set("network_capacity_in_mbps", response.NetworkCapacityInMbps)
	// d.Set("deploy_set_id", response.DeploysetId)

	args := &bbc.GetVpcSubnetArgs{
		BbcIds: []string{instanceID},
	}
	raw, err = client.WithBbcClient(func(bbcClient *bbc.Client) (interface{}, error) {
		return bbcClient.GetVpcSubnet(args)
	})
	addDebug(action, raw)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_instance", action, BCESDKGoERROR)
	}
	vpcs, _ := raw.(*bbc.GetVpcSubnetResult)
	for _, sg := range vpcs.NetworkInfo {
		d.Set("subnet_id", sg.Subnet.SubnetId)
		d.Set("vpc_id", sg.Vpc.VpcId)
	}
	//TODO: unsupported read
	//if sg, ok := d.GetOk("security_group"); !ok {
	//	d.Set("security_group", sg)
	//}

	// read all disks
	volumes, err := bbcService.ListAllVolumes(instanceID)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_instance", action, BCESDKGoERROR)
	}

	cdsDisks := make([]interface{}, 0, len(volumes))
	for _, vol := range volumes {
		cdsMap := make(map[string]interface{})
		cdsMap["disk_size_in_gb"] = vol.DiskSizeInGB
		cdsMap["storage_type"] = vol.StorageType
		cdsMap["is_system_volume"] = vol.IsSystemVolume
		cdsDisks = append(cdsDisks, cdsMap)
	}
	d.Set("cds_disks", cdsDisks)

	return nil
}

func resourceBaiduCloudBbcInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	instanceID := d.Id()

	d.Partial(true)

	// update instance attribute
	if err := updateBbcInstanceName(d, meta, instanceID); err != nil {
		return err
	}

	// update instance description
	if err := updateBbcInstanceDescription(d, meta, instanceID); err != nil {
		return err
	}

	// update instance image id (rebuild)
	if err := updateBbcInstanceImage(d, meta, instanceID); err != nil {
		return err
	}

	// update instance admin pass
	if err := updateBbcInstanceAdminPass(d, meta, instanceID); err != nil {
		return err
	}

	// update instance security groups
	if err := updateBbcInstanceSecurityGroups(d, meta, instanceID); err != nil {
		return err
	}

	// update tags
	if err := updateBbcInstanceTags(d, meta, instanceID); err != nil {
		return err
	}
	// update instance subnet
	if err := updateBbcInstanceSubnet(d, meta, instanceID); err != nil {
		return err
	}

	d.Partial(false)

	return resourceBaiduCloudInstanceRead(d, meta)
}

func resourceBaiduCloudBbcInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	bbcService := BbcService{client}

	instanceId := d.Id()
	action := "Delete BBC Instance " + instanceId
	args := &bbc.DeleteInstanceIngorePaymentArgs{
		InstanceId:         instanceId,
		RelatedReleaseFlag: true, //true or false
	}
	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		raw, err := client.WithBbcClient(func(bbcClient *bbc.Client) (interface{}, error) {
			return bbcClient.DeleteInstanceIngorePayment(args)
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
		if IsExceptedErrors(err, BbcNotFound) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_instance", action, BCESDKGoERROR)
	}

	stateConf := buildStateConf(
		[]string{string(bbc.InstanceStatusStopping), string(bbc.InstanceStatusStopped)},
		[]string{string(bbc.InstanceStatusDeleted)},
		d.Timeout(schema.TimeoutDelete),
		bbcService.InstanceBbcStateRefresh(instanceId),
	)
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_instance", action, BCESDKGoERROR)
	}

	return nil
}

func buildBaiduCloudBbcInstanceArgs(d *schema.ResourceData, meta interface{}) (*bbc.CreateInstanceArgs, error) {
	request := &bbc.CreateInstanceArgs{
		ClientToken:   buildClientToken(),
		PurchaseCount: 1,
	}

	if imageID, ok := d.GetOk("image_id"); ok {
		request.ImageId = imageID.(string)
	}

	if name, ok := d.GetOk("name"); ok {
		request.Name = name.(string)
	}
	if hostname, ok := d.GetOk("hostname"); ok {
		request.Hostname = hostname.(string)
	}

	if zoneName, ok := d.GetOk("availability_zone"); ok {
		request.ZoneName = zoneName.(string)
	}

	if flavorId, ok := d.GetOk("flavor_id"); ok {
		request.FlavorId = flavorId.(string)
	}

	if raidId, ok := d.GetOk("raid_id"); ok {
		request.RaidId = raidId.(string)
	}

	if v, ok := d.GetOk("billing"); ok {
		billing := v.(map[string]interface{})
		billingRequest := bbc.Billing{
			PaymentTiming: bbc.PaymentTimingType(""),
			Reservation:   bbc.Reservation{},
		}
		if p, ok := billing["payment_timing"]; ok {
			paymentTiming := bbc.PaymentTimingType(p.(string))
			billingRequest.PaymentTiming = paymentTiming
		}
		if billingRequest.PaymentTiming == bbc.PaymentTimingPrePaid {
			if r, ok := billing["reservation"]; ok {
				reservation := r.(map[string]interface{})
				if reservationLength, ok := reservation["reservation_length"]; ok {
					billingRequest.Reservation.Length = reservationLength.(int)
				}
				if reservationTimeUnit, ok := reservation["reservation_time_unit"]; ok {
					billingRequest.Reservation.TimeUnit = reservationTimeUnit.(string)
				}
			}
			// if the field is set, then auto-renewal is effective.
			if v, ok := d.GetOk("auto_renew_time_unit"); ok {
				request.AutoRenewTimeUnit = v.(string)
				if v, ok := d.GetOk("auto_renew_time_length"); ok {
					request.AutoRenewTime = v.(int)
				}
			}
		}

		request.Billing = billingRequest
	}

	if adminPass, ok := d.GetOk("admin_pass"); ok {
		request.AdminPass = adminPass.(string)
	}
	if rootDiskSizeInGb, ok := d.GetOk("root_disk_size_in_gb"); ok {
		request.RootDiskSizeInGb = rootDiskSizeInGb.(int)
	}
	// not supported
	if deploySetId, ok := d.GetOk("deploy_set_id"); ok {
		request.DeploySetId = deploySetId.(string)
	}
	// not supported
	if requestToken, ok := d.GetOk("request_token"); ok {
		request.RequestToken = requestToken.(string)
	}

	if subnetId, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = subnetId.(string)
	}

	if v, ok := d.GetOk("tags"); ok {
		request.Tags = tranceTagMapToModel(v.(map[string]interface{}))
	}

	return request, nil
}

func updateBbcInstanceName(d *schema.ResourceData, meta interface{}, instanceID string) error {
	action := "Update BBC Instance attribute " + instanceID
	client := meta.(*connectivity.BaiduClient)

	if d.HasChange("name") {
		modifyInstanceNameArgs := &bbc.ModifyInstanceNameArgs{}
		modifyInstanceNameArgs.Name = d.Get("name").(string)
		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			_, err := client.WithBbcClient(func(bbcClient *bbc.Client) (interface{}, error) {
				return nil, bbcClient.ModifyInstanceName(instanceID, modifyInstanceNameArgs)
			})
			if err != nil {
				if IsExceptedErrors(err, []string{OperationDenied, bce.EINTERNAL_ERROR}) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, modifyInstanceNameArgs)
			return nil
		})

		if err != nil {
			if NotFoundError(err) {
				d.SetId("")
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_instance", action, BCESDKGoERROR)
		}

		d.SetPartial("name")
	}

	return nil
}

func updateBbcInstanceDescription(d *schema.ResourceData, meta interface{}, instanceID string) error {
	action := "Update BBC Instance Description " + instanceID
	client := meta.(*connectivity.BaiduClient)

	if d.HasChange("description") {
		modifyInstanceDescArgs := &bbc.ModifyInstanceDescArgs{}
		modifyInstanceDescArgs.Description = d.Get("description").(string)

		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			_, err := client.WithBbcClient(func(bbcClient *bbc.Client) (interface{}, error) {
				return nil, bbcClient.ModifyInstanceDesc(instanceID, modifyInstanceDescArgs)
			})
			if err != nil {
				if IsExceptedErrors(err, []string{OperationDenied, bce.EINTERNAL_ERROR}) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, modifyInstanceDescArgs)
			return nil
		})

		if err != nil {
			if NotFoundError(err) {
				d.SetId("")
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_instance", action, BCESDKGoERROR)
		}

		d.SetPartial("description")
	}

	return nil
}

func updateBbcInstanceImage(d *schema.ResourceData, meta interface{}, instanceID string) error {
	action := "Update instance image " + instanceID
	client := meta.(*connectivity.BaiduClient)
	bbcService := &BbcService{client}

	if d.HasChange("image_id") {
		args := &bbc.RebuildInstanceArgs{
			ImageId: d.Get("image_id").(string),
		}
		if adminPass, ok := d.GetOk("admin_pass"); ok {
			args.AdminPass = adminPass.(string)
		}
		if d.HasChange("raid_id") || d.HasChange("root_disk_size_in_gb") {
			args.IsPreserveData = false
			args.RaidId = d.Get("raid_id").(string)
			args.SysRootSize = d.Get("root_disk_size_in_gb").(int)
		} else {
			args.IsPreserveData = true
		}
		if _, err := client.WithBbcClient(func(bbcClient *bbc.Client) (i interface{}, e error) {
			return nil, bbcClient.RebuildInstance(instanceID, args.IsPreserveData, args)
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_instance", action, BCESDKGoERROR)
		}

		stateConf := buildStateConf(
			[]string{string(bbc.InstanceStatusStarting), string(bbc.InstanceStatusImageProcessing)},
			[]string{string(bbc.InstanceStatusRunning)},
			d.Timeout(schema.TimeoutUpdate),
			bbcService.InstanceBbcStateRefresh(instanceID),
		)
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_instance", action, BCESDKGoERROR)
		}

		d.SetPartial("image_id")
		d.SetPartial("admin_pass")
		if !args.IsPreserveData {
			d.SetPartial("raid_id")
			d.SetPartial("root_disk_size_in_gb")
		}

	}

	return nil
}

func updateBbcInstanceAdminPass(d *schema.ResourceData, meta interface{}, instanceID string) error {
	action := "Update Instance admin pass " + instanceID
	client := meta.(*connectivity.BaiduClient)
	bbcService := &BbcService{client}

	if d.HasChange("admin_pass") {
		args := &bbc.ModifyInstancePasswordArgs{
			AdminPass: d.Get("admin_pass").(string),
		}

		if _, err := client.WithBbcClient(func(bbcClient *bbc.Client) (i interface{}, e error) {
			return nil, bbcClient.ModifyInstancePassword(instanceID, args)
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_instance", action, BCESDKGoERROR)
		}

		stateConf := buildStateConf(
			[]string{string(bbc.InstanceStatusStarting)},
			[]string{string(bbc.InstanceStatusRunning)},
			d.Timeout(schema.TimeoutUpdate),
			bbcService.InstanceBbcStateRefresh(instanceID),
		)
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_instance", action, BCESDKGoERROR)
		}

		d.SetPartial("admin_pass")
	}

	return nil
}

func updateBbcInstanceSecurityGroups(d *schema.ResourceData, meta interface{}, instanceID string) error {
	action := "Update instance security groups " + instanceID
	client := meta.(*connectivity.BaiduClient)

	if d.HasChange("security_group") {
		o, n := d.GetChange("security_group")

		os := o.(string)
		ns := n.(string)
		if len(os) > 0 {
			args := &bbc.UnBindSecurityGroupsArgs{
				InstanceId:      instanceID,
				SecurityGroupId: os,
			}
			if _, err := client.WithBbcClient(func(bbcClient *bbc.Client) (i interface{}, e error) {
				return nil, bbcClient.UnBindSecurityGroups(args)
			}); err != nil {
				return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_instance", action, BCESDKGoERROR)
			}
		}
		if len(ns) > 0 {
			args := &bbc.BindSecurityGroupsArgs{
				InstanceIds:      []string{instanceID},
				SecurityGroupIds: []string{ns},
			}
			if _, err := client.WithBbcClient(func(bbcClient *bbc.Client) (i interface{}, e error) {
				return nil, bbcClient.BindSecurityGroups(args)
			}); err != nil {
				return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_instance", action, BCESDKGoERROR)
			}
		}
		d.SetPartial("security_group")
	}

	return nil
}

func updateBbcInstanceTags(d *schema.ResourceData, meta interface{}, instanceID string) error {
	action := "Update instance security groups " + instanceID
	client := meta.(*connectivity.BaiduClient)

	if d.HasChange("tags") {
		o, n := d.GetChange("tags")
		if o != nil {
			tagModes := tranceTagMapToModel(o.(map[string]interface{}))
			unbindTagsArgs := &bbc.UnbindTagsArgs{
				ChangeTags: tagModes,
			}
			if _, err := client.WithBbcClient(func(bbcClient *bbc.Client) (data interface{}, e error) {
				return nil, bbcClient.UnbindTags(instanceID, unbindTagsArgs)
			}); err != nil {
				return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_instance", action, BCESDKGoERROR)
			}
		}
		if n != nil {
			tagModes := tranceTagMapToModel(o.(map[string]interface{}))
			bindTagsArgs := &bbc.BindTagsArgs{
				ChangeTags: tagModes,
			}
			if _, err := client.WithBbcClient(func(bbcClient *bbc.Client) (data interface{}, e error) {
				return nil, bbcClient.BindTags(instanceID, bindTagsArgs)
			}); err != nil {
				return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_instance", action, BCESDKGoERROR)
			}
		}
		d.SetPartial("tags")
	}

	return nil
}

func updateBbcInstanceSubnet(d *schema.ResourceData, meta interface{}, instanceID string) error {
	action := "Update instance subnet " + instanceID
	client := meta.(*connectivity.BaiduClient)
	bbcService := BbcService{client}

	if d.HasChange("subnet_id") {
		args := &bbc.InstanceChangeSubnetArgs{
			InstanceId: instanceID,
			Reboot:     true,
		}
		if v, ok := d.GetOk("subnet_id"); ok {
			args.SubnetId = v.(string)
		}

		_, err := client.WithBbcClient(func(bbcClient *bbc.Client) (i interface{}, e error) {
			return nil, bbcClient.InstanceChangeSubnet(args)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_instance", action, BCESDKGoERROR)
		}

		stateConf := buildStateConf(
			[]string{string(bbc.InstanceStatusStopping), string(bbc.InstanceStatusStopped), string(bbc.InstanceStatusStarting)},
			[]string{string(bbc.InstanceStatusRunning)},
			d.Timeout(schema.TimeoutUpdate),
			bbcService.InstanceBbcStateRefresh(instanceID),
		)
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_instance", action, BCESDKGoERROR)
		}

		d.SetPartial("subnet_id")
	}

	return nil
}
