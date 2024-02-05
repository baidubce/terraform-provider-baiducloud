/*
Use this resource to get information about a RDS instance.

~> **NOTE:** The terminate operation of rds instance does NOT take effect immediately，maybe takes for several minites.

# Example Usage

```hcl

	resource "baiducloud_rds_instance" "default" {
	    billing = {
	        payment_timing        = "Postpaid"
	    }
	    engine_version            = "5.6"
	    engine                    = "MySQL"
	    cpu_count                 = 1
	    memory_capacity           = 1
	    volume_capacity           = 5
	}

```

# Import

RDS instance can be imported, e.g.

```hcl
$ terraform import baiducloud_rds_instance.default id
```
*/
package baiducloud

import (
	"strconv"
	"time"

	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/rds"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func resourceBaiduCloudRdsInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudRdsInstanceCreate,
		Read:   resourceBaiduCloudRdsInstanceRead,
		Update: resourceBaiduCloudRdsInstanceUpdate,
		Delete: resourceBaiduCloudRdsInstanceDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"purchase_count": {
				Type:        schema.TypeInt,
				Description: "Count of the instance to buy",
				Default:     1,
				Optional:    true,
			},
			"instance_name": {
				Type:        schema.TypeString,
				Description: "Name of the instance. Support for uppercase and lowercase letters, numbers, Chinese and special characters, such as \"-\",\"_\",\"/\",\".\", the value must start with a letter, length 1-65.",
				Optional:    true,
				Computed:    true,
			},
			"engine_version": {
				Type:        schema.TypeString,
				Description: "Engine version of the instance. MySQL support 5.5、5.6、5.7, SQLServer support 2008r2、2012sp3、2016sp1, PostgreSQL support 9.4",
				Required:    true,
				ForceNew:    true,
			},
			"engine": {
				Type:         schema.TypeString,
				Description:  "Engine of the instance. Available values are MySQL、SQLServer、PostgreSQL.",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"MySQL", "SQLServer", "PostgreSQL"}, false),
			},
			"category": {
				Type:         schema.TypeString,
				Description:  "Category of the instance. Available values are Basic、Standard(Default), only SQLServer 2012sp3 support Basic.",
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Standard", "Basic"}, false),
			},
			"cpu_count": {
				Type:         schema.TypeInt,
				Description:  "The number of CPU",
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"memory_capacity": {
				Type:         schema.TypeFloat,
				Description:  "Memory capacity(GB) of the instance.",
				Required:     true,
				ValidateFunc: validation.FloatBetween(1, 480),
			},
			"volume_capacity": {
				Type:         schema.TypeInt,
				Description:  "Volume capacity(GB) of the instance",
				Required:     true,
				ValidateFunc: validation.IntAtLeast(5),
			},
			"disk_io_type": {
				Type:         schema.TypeString,
				Description:  "Type of disk, Available values are normal_io,cloud_high,cloud_nor,cloud_enha",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"normal_io", "cloud_high", "cloud_nor", "cloud_enha"}, false),
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Description: "ID of the specific VPC",
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
			},
			"subnets": {
				Type:        schema.TypeList,
				Description: "Subnets of the instance.",
				Optional:    true,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_id": {
							Type:        schema.TypeString,
							Description: "ID of the subnet.",
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
						},
						"zone_name": {
							Type:        schema.TypeString,
							Description: "Zone name of the subnet.",
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
						},
					},
				},
			},
			"zone_names": {
				Type:        schema.TypeList,
				Description: "Zone name list",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tags": tagsSchema(),
			"instance_id": {
				Type:        schema.TypeString,
				Description: "ID of the instance.",
				Computed:    true,
			},
			"instance_status": {
				Type:        schema.TypeString,
				Description: "Status of the instance.",
				Computed:    true,
			},
			"node_amount": {
				Type:        schema.TypeInt,
				Description: "Number of proxy node.",
				Computed:    true,
			},
			"used_storage": {
				Type:        schema.TypeFloat,
				Description: "Memory capacity(GB) of the instance to be used.",
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
			"address": {
				Type:        schema.TypeString,
				Description: "The domain used to access a instance.",
				Computed:    true,
			},
			"port": {
				Type:        schema.TypeInt,
				Description: "The port used to access a instance.",
				Computed:    true,
			},
			"v_net_ip": {
				Type:        schema.TypeString,
				Description: "The internal ip used to access a instance.",
				Computed:    true,
			},
			"region": {
				Type:        schema.TypeString,
				Description: "Region of the instance.",
				Computed:    true,
			},
			"instance_type": {
				Type:        schema.TypeString,
				Description: "Type of the instance,  Available values are Master, ReadReplica, RdsProxy.",
				Computed:    true,
			},
			"lower_case_table_names": {
				Type: schema.TypeInt,
				Description: "Whether the table name is case-sensitive. The default value is 0, " +
					"which means case-sensitive; passing 1 means case-insensitive.",
				Optional: true,
			},
			"parameter_template_id": {
				Type:        schema.TypeString,
				Description: "Parameter template id.",
				Optional:    true,
			},
			"replication_type": {
				Type: schema.TypeString,
				Description: "Data replication method. Asynchronous replication: async, " +
					"Semi-synchronous replication: semi_sync.",
				Optional: true,
			},
			"resource_group_id": {
				Type:        schema.TypeString,
				Description: "resource group id.",
				Optional:    true,
			},
			"public_access": {
				Type:        schema.TypeBool,
				Description: "public access.",
				Optional:    true,
			},
			"billing": {
				Type:        schema.TypeMap,
				Description: "Billing information of the Rds.",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"payment_timing": {
							Type:         schema.TypeString,
							Description:  "Payment timing of billing, which can be Prepaid or Postpaid. The default is Postpaid.",
							Required:     true,
							Default:      PaymentTimingPostpaid,
							ValidateFunc: validatePaymentTiming(),
						},
					},
				},
			},
			"reservation": {
				Type:             schema.TypeMap,
				Description:      "Reservation of the Rds.",
				Optional:         true,
				DiffSuppressFunc: postPaidDiffSuppressFunc,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"reservation_length": {
							Type: schema.TypeInt,
							Description: "The reservation length that you will pay for your resource. " +
								"It is valid when payment_timing is Prepaid. " +
								"Valid values: [1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36].",
							Required:         true,
							Default:          1,
							ValidateFunc:     validateReservationLength(),
							DiffSuppressFunc: postPaidDiffSuppressFunc,
						},
						"reservation_time_unit": {
							Type: schema.TypeString,
							Description: "The reservation time unit that you will pay for your resource. " +
								"It is valid when payment_timing is Prepaid. " +
								"The value can only be month currently, which is also the default value.",
							Required:         true,
							Default:          "Month",
							ValidateFunc:     validateReservationUnit(),
							DiffSuppressFunc: postPaidDiffSuppressFunc,
						},
					},
				},
			},
			"payment_timing": {
				Type:        schema.TypeString,
				Description: "RDS payment timing",
				Computed:    true,
			},
			"auto_renew_time_unit": {
				Type: schema.TypeString,
				Description: "Time unit of automatic renewal, the value can be month or year. " +
					"The default value is empty, indicating no automatic renewal. " +
					"It is valid only when the payment_timing is Prepaid.",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"month", "year"}, false),
			},
			"auto_renew_time_length": {
				Type: schema.TypeInt,
				Description: "The time length of automatic renewal. It is valid when payment_timing is Prepaid, " +
					"and the value should be 1-9 when the auto_renew_time_unit is month and 1-3 when the " +
					"auto_renew_time_unit is year. Default to 1.",
				Optional:     true,
				ForceNew:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 9),
			},
			"backup_days": {
				Type: schema.TypeString,
				Description: "Backup date and time separated by English half-width commas, " +
					"Sunday is the first day, the value is 0 Example: 0,1,2,3,5,6 ",
				Optional: true,
			},
			"backup_time": {
				Type:        schema.TypeString,
				Description: "Backup start time, the time here is UTC time",
				Optional:    true,
			},
			"expire_in_days": {
				Type:        schema.TypeInt,
				Description: "Number of persistence days, range 1-730 days; if not enabled, it is 0 or left blank",
				Optional:    true,
			},
		},
	}
}

func resourceBaiduCloudRdsInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	rdsService := RdsService{client}

	createRdsArgs, err := buildBaiduCloudRdsInstanceArgs(d, meta)
	if err != nil {
		return WrapError(err)
	}

	action := "Create RDS Instance " + createRdsArgs.InstanceName
	addDebug(action, createRdsArgs)

	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.CreateRds(createRdsArgs)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, raw)
		response, _ := raw.(*rds.CreateResult)
		d.SetId(response.InstanceIds[0])
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_rds_instance", action, BCESDKGoERROR)
	}

	stateConf := buildStateConf(
		[]string{RDSStatusCreating},
		[]string{RDSStatusRunning},
		d.Timeout(schema.TimeoutCreate),
		rdsService.InstanceStateRefresh(d.Id(), []string{}),
	)
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_rds_instance", action, BCESDKGoERROR)
	}
	// 开启公网访问
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			var publicAccess bool
			if v, ok := d.GetOk("public_access"); ok {
				publicAccess = v.(bool)
			}
			args := &rds.ModifyPublicAccessArgs{
				PublicAccess: publicAccess,
			}
			return nil, rdsClient.ModifyPublicAccess(d.Id(), args)
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
		addDebug(action, err)
	}
	// 开启自动续费
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			args := &rds.AutoRenewArgs{
				InstanceIds: []string{d.Id()},
			}
			// if the field is set, then auto-renewal is effective.
			if v, ok := d.GetOk("auto_renew_time_unit"); ok {
				args.AutoRenewTimeUnit = v.(string)
				if v, ok := d.GetOk("auto_renew_time_length"); ok {
					args.AutoRenewTime = v.(int)
				}
			} else {
				return nil, nil
			}
			return nil, rdsClient.AutoRenew(args)
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
		addDebug(action, err)
	}
	// 设置备份策略
	err = setBackupPolicy(d, meta, d.Id())
	if err != nil {
		addDebug(action, err)
	}
	return resourceBaiduCloudRdsInstanceRead(d, meta)
}

func resourceBaiduCloudRdsInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	instanceID := d.Id()
	action := "Query RDS Instance " + instanceID

	raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
		return rdsClient.GetDetail(instanceID)
	})

	addDebug(action, raw)

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_rds_instance", action, BCESDKGoERROR)
	}

	result, _ := raw.(*rds.Instance)

	d.Set("instance_id", result.InstanceId)
	d.Set("instance_name", result.InstanceName)
	d.Set("engine", result.Engine)
	d.Set("engine_version", result.EngineVersion)
	d.Set("category", result.Category)
	d.Set("instance_status", result.InstanceStatus)
	d.Set("cpu_count", result.CpuCount)
	d.Set("memory_capacity", result.MemoryCapacity)
	d.Set("volume_capacity", result.VolumeCapacity)
	d.Set("node_amount", result.NodeAmount)
	d.Set("used_storage", result.UsedStorage)
	d.Set("create_time", result.InstanceCreateTime)
	d.Set("expire_time", result.InstanceExpireTime)
	d.Set("region", result.Region)
	d.Set("instance_type", result.InstanceType)
	d.Set("payment_timing", result.PaymentTiming)
	d.Set("zone_names", result.ZoneNames)
	d.Set("vpc_id", result.VpcId)
	d.Set("port", result.Endpoint.Port)
	d.Set("address", result.Endpoint.Address)
	d.Set("v_net_ip", result.Endpoint.VnetIp)
	d.Set("volume_capacity", result.VolumeCapacity)
	d.Set("subnets", transRdsSubnetsToSchema(result.Subnets))
	d.Set("public_access", result.PublicAccessStatus)
	d.Set("backup_days", result.BackupPolicy.BackupDays)
	d.Set("backup_time", result.BackupPolicy.BackupTime)
	d.Set("expire_in_days", result.BackupPolicy.ExpireInDays)
	return nil
}

func transRdsSubnetsToSchema(subnets []rds.Subnet) []map[string]string {
	subnetList := []map[string]string{}
	for _, subnet := range subnets {
		subnetMap := make(map[string]string)
		subnetMap["subnet_id"] = subnet.SubnetId
		subnetMap["zone_name"] = subnet.ZoneName
		subnetList = append(subnetList, subnetMap)
	}
	return subnetList
}

func resourceBaiduCloudRdsInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	instanceID := d.Id()
	action := "Update RDS Instance " + instanceID
	client := meta.(*connectivity.BaiduClient)

	d.Partial(true)

	// resize instance
	if err := resizeRds(d, meta, instanceID); err != nil {
		return err
	}
	// 公网访问权限
	if d.HasChange("public_access") {
		err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
				var publicAccess bool
				if v, ok := d.GetOk("public_access"); ok {
					publicAccess = v.(bool)
				}
				args := &rds.ModifyPublicAccessArgs{
					PublicAccess: publicAccess,
				}
				return nil, rdsClient.ModifyPublicAccess(d.Id(), args)
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
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_rds_instance", action, BCESDKGoERROR)
		}
	}

	// 设置备份策略
	err := setBackupPolicy(d, meta, d.Id())
	if err != nil {
		addDebug(action, err)
	}
	d.Partial(false)

	return resourceBaiduCloudRdsInstanceRead(d, meta)
}

func resourceBaiduCloudRdsInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	instanceId := d.Id()
	action := "Delete RDS Instance " + instanceId

	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return instanceId, rdsClient.DeleteRds(instanceId)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{InvalidInstanceStatus, bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, raw)
		return nil
	})
	if err != nil {
		if IsExceptedErrors(err, []string{InvalidInstanceStatus, InstanceNotExist, bce.EINTERNAL_ERROR}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_rds_instance", action, BCESDKGoERROR)
	}

	return nil
}

func buildBaiduCloudRdsInstanceArgs(d *schema.ResourceData, meta interface{}) (*rds.CreateRdsArgs, error) {
	request := &rds.CreateRdsArgs{
		ClientToken: buildClientToken(),
		IsDirectPay: true,
	}

	if v, ok := d.GetOk("billing"); ok {
		billing := v.(map[string]interface{})
		billingRequest := rds.Billing{
			PaymentTiming: "",
			Reservation:   rds.Reservation{},
		}
		if p, ok := billing["payment_timing"]; ok {
			paymentTiming := p.(string)
			billingRequest.PaymentTiming = paymentTiming
		}
		if billingRequest.PaymentTiming == PaymentTimingPrepaid {
			if r, ok := d.GetOk("reservation"); ok {
				reservation := r.(map[string]interface{})
				if reservationLength, ok := reservation["reservation_length"]; ok {
					reservationLengthInt, err := strconv.Atoi(reservationLength.(string))
					billingRequest.Reservation.ReservationLength = reservationLengthInt
					if err != nil {
						return nil, err
					}
				}
				if reservationTimeUnit, ok := reservation["reservation_time_unit"]; ok {
					billingRequest.Reservation.ReservationTimeUnit = reservationTimeUnit.(string)
				}
			}
		}
		request.Billing = billingRequest
	}

	if purchaseCount, ok := d.GetOk("purchase_count"); ok {
		request.PurchaseCount = purchaseCount.(int)
	}

	if diskIoType, ok := d.GetOk("disk_io_type"); ok {
		request.DiskIoType = diskIoType.(string)
	}

	if instanceName, ok := d.GetOk("instance_name"); ok {
		request.InstanceName = instanceName.(string)
	}

	if engineVersion, ok := d.GetOk("engine_version"); ok {
		request.EngineVersion = engineVersion.(string)
	}

	if engine, ok := d.GetOk("engine"); ok {
		request.Engine = engine.(string)
	}

	if category, ok := d.GetOk("category"); ok {
		request.Category = category.(string)
	}

	if cpuCount, ok := d.GetOk("cpu_count"); ok {
		request.CpuCount = cpuCount.(int)
	}

	if memoryCapacity, ok := d.GetOk("memory_capacity"); ok {
		request.MemoryCapacity = memoryCapacity.(float64)
	}

	if lowerCaseTableNames, ok := d.GetOk("lower_case_table_names"); ok {
		request.LowerCaseTableNames = lowerCaseTableNames.(int)
	}

	if parameterTemplateId, ok := d.GetOk("parameter_template_id"); ok {
		request.ParameterTemplateId = parameterTemplateId.(string)
	}

	if replicationType, ok := d.GetOk("replication_type"); ok {
		request.ReplicationType = replicationType.(string)
	}

	if resourceGroupId, ok := d.GetOk("resource_group_id"); ok {
		request.ResourceGroupId = resourceGroupId.(string)
	}

	if volumeCapacity, ok := d.GetOk("volume_capacity"); ok {
		request.VolumeCapacity = volumeCapacity.(int)
	}

	if vpcID, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = vpcID.(string)
	}

	if v, ok := d.GetOk("subnets"); ok {
		subnetList := v.([]interface{})
		subnetRequests := make([]rds.SubnetMap, len(subnetList))
		for id := range subnetList {
			subnet := subnetList[id].(map[string]interface{})

			subnetRequest := rds.SubnetMap{
				SubnetId: subnet["subnet_id"].(string),
				ZoneName: subnet["zone_name"].(string),
			}

			subnetRequests[id] = subnetRequest
		}
		request.Subnets = subnetRequests
	}

	return request, nil

}

func resizeRds(d *schema.ResourceData, meta interface{}, instanceID string) error {
	action := "Update rds nodeType " + instanceID
	client := meta.(*connectivity.BaiduClient)
	rdsService := RdsService{client}

	if d.HasChange("cpu_count") || d.HasChange("memory_capacity") || d.HasChange("volume_capacity") {
		args := &rds.ResizeRdsArgs{
			CpuCount:       d.Get("cpu_count").(int),
			MemoryCapacity: d.Get("memory_capacity").(float64),
			VolumeCapacity: d.Get("volume_capacity").(int),
		}

		addDebug(action, args)
		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			_, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
				return nil, rdsClient.ResizeRds(instanceID, args)
			})
			if err != nil {
				if IsExceptedErrors(err, []string{InvalidInstanceStatus, OperationException, bce.EINTERNAL_ERROR}) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_rds_instance", action, BCESDKGoERROR)
		}

		stateConf := buildStateConf(
			[]string{RDSStatusModifying},
			[]string{RDSStatusRunning},
			d.Timeout(schema.TimeoutUpdate),
			rdsService.InstanceStateRefresh(d.Id(), []string{}),
		)
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_rds_instance", action, BCESDKGoERROR)
		}

		d.SetPartial("cpu_count")
		d.SetPartial("memory_capacity")
		d.SetPartial("volume_capacity")
	}

	return nil
}

func setBackupPolicy(d *schema.ResourceData, meta interface{}, instanceID string) error {
	action := "Set rds backup policy " + instanceID
	client := meta.(*connectivity.BaiduClient)

	if d.HasChange("backup_days") || d.HasChange("backup_time") || d.HasChange("expire_in_days") {
		args := &rds.ModifyBackupPolicyArgs{
			BackupDays:   d.Get("backup_days").(string),
			BackupTime:   d.Get("backup_time").(string),
			ExpireInDays: d.Get("expire_in_days").(int),
		}

		addDebug(action, args)
		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			_, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
				return nil, rdsClient.ModifyBackupPolicy(instanceID, args)
			})
			if err != nil {
				if IsExceptedErrors(err, []string{InvalidInstanceStatus, OperationException, bce.EINTERNAL_ERROR}) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_rds_instance", action, BCESDKGoERROR)
		}

		d.SetPartial("backup_days")
		d.SetPartial("backup_time")
		d.SetPartial("expire_in_days")
	}

	return nil
}
