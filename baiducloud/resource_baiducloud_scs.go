/*
Use this resource to get information about a SCS.

More information about SCS can be found in the [Developer Guide](https://cloud.baidu.com/doc/SCS/index.html).

~> **NOTE:** The terminate operation of scs does NOT take effect immediately，maybe takes for several minites.

# Example Usage

### Memcache
~> **NOTE:** Memcache currently does NOT support specifying `node_type`, set to `cache.n1.micro` directly.
```terraform

	resource "baiducloud_scs" "default" {
		payment_timing = "Postpaid"
		instance_name = "terraform-memcache"
		engine = "memcache"
		port = 11211
		node_type = "cache.n1.micro"
		cluster_type = "defalut"
		shard_num = 2
	}

```

### Redis
```terraform

	resource "baiducloud_scs" "default" {
		payment_timing = "Postpaid"
		instance_name = "terraform-redis"
		port = 6379
		engine_version = "3.2"
		node_type = "cache.n1.micro"
		cluster_type = "master_slave"
		replication_num = 1
		shard_num = 1
	}

```

### PegaDb
```terraform

	resource "baiducloud_scs" "default" {
		payment_timing = "Prepaid"
		reservation_length = 2
		reservation_time_unit = "month"
		instance_name = "terraform-pegadb"
		purchase_count = 1
		engine = "PegaDB"
		node_type = "pega.g4s1.micro"
		cluster_type = "cluster"
		store_type = 3
		disk_flavor = 60
		port = 6379
		replication_num = 2
		shard_num = 1
		proxy_num = 2
		vpc_id = "vpc-ne32rahkaceu"
		subnets {
			subnet_id = "sbn-vhnqd71mivjq"
			zone_name = "cn-bj-d"
		}
		replication_info {
			availability_zone = "cn-bj-d"
			is_master         = 1
			subnet_id         = "sbn-vhnqd71mivjq"
		}
		replication_info {
			availability_zone = "cn-bj-d"
			is_master         = 0
			subnet_id         = "sbn-vhnqd71mivjq"
		}
	}

```

# Import

SCS can be imported, e.g.

```hcl
$ terraform import baiducloud_scs.default id
```
*/
package baiducloud

import (
	"strings"
	"time"

	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/scs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func resourceBaiduCloudScs() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudScsCreate,
		Read:   resourceBaiduCloudScsRead,
		Update: resourceBaiduCloudScsUpdate,
		Delete: resourceBaiduCloudScsDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Minute),
			Update: schema.DefaultTimeout(120 * time.Minute),
			Delete: schema.DefaultTimeout(120 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"purchase_count": {
				Type:         schema.TypeInt,
				Description:  "Count of the instance to buy. Must be between `1` and `10`. Defaults to `1`.",
				Default:      1,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 10),
			},
			"instance_name": {
				Type:        schema.TypeString,
				Description: "Name of the instance. Support for uppercase and lowercase letters, numbers, Chinese and special characters, such as `-`, `_`, `/`, `.`. Must start with a letter, length 1-65.",
				Required:    true,
			},
			"node_type": {
				Type:        schema.TypeString,
				Description: "Node type of the instance. e.g. `cache.n1.micro`. To learn about supported node type, see documentation on [Supported Node Types](https://cloud.baidu.com/doc/SCS/s/1jwvxtsh0#%E5%AE%9E%E4%BE%8B%E8%A7%84%E6%A0%BC)",
				Required:    true,
			},
			"shard_num": {
				Type:        schema.TypeInt,
				Description: "The number of instance shard. Defaults to `1`. To learn about supported shard number, see documentation on [Supported Node Types](https://cloud.baidu.com/doc/SCS/s/1jwvxtsh0#%E5%AE%9E%E4%BE%8B%E8%A7%84%E6%A0%BC)",
				Default:     1,
				Optional:    true,
			},
			"proxy_num": {
				Type:        schema.TypeInt,
				Description: "The number of instance proxy. If `cluster_type` is `cluster`, set to the value of `shard_num` (if `shard_num` equals `1`, set to `2`). If `cluster_type` is `master_slave`, set to `0`. Defaults to `0`.",
				Default:     0,
				Optional:    true,
				ForceNew:    true,
			},
			"replication_num": {
				Type:        schema.TypeInt,
				Description: "The number of instance replicas. If `cluster_type` is `cluster`, must be between `2` and `5`. If `cluster_type` is `master_slave`, must be between `1` and `5`. Defaults to `2`.",
				Default:     2,
				Optional:    true,
				ForceNew:    true,
			},
			"port": {
				Type:        schema.TypeInt,
				Description: "Port number used to access the instance. Must be between `1025` and `65534`. Defaults to `6379`.",
				Optional:    true,
				Default:     6379,
				ForceNew:    true,
			},
			"domain": {
				Type:        schema.TypeString,
				Description: "Domain of the instance.",
				Computed:    true,
			},
			"cluster_type": {
				Type:         schema.TypeString,
				Description:  "Type of the instance. If `engine` is `memcache`, must be `default`. Valid values for other engine type: `cluster`, `master_slave`.  Defaults to `master_slave`.",
				Optional:     true,
				ForceNew:     true,
				Default:      "master_slave",
				ValidateFunc: validation.StringInSlice([]string{"cluster", "master_slave", "default"}, false),
			},
			"engine_version": {
				Type:        schema.TypeString,
				Description: "Engine version of the instance. Must be set when `engine` is `redis`. Valid values: `3.2`, `4.0`, `5.0`, `6.0`.",
				Optional:    true,
				Computed:    true,
			},
			"engine": {
				Type:         schema.TypeString,
				Description:  "Engine of the instance. Valid values: `memcache`, `redis`, `PegaDB`. Defaults to `redis`.",
				Optional:     true,
				Default:      "redis",
				ValidateFunc: validation.StringInSlice([]string{"memcache", "redis", "PegaDB"}, false),
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Description: "ID of the specific VPC",
				Optional:    true,
				Computed:    true,
			},
			"v_net_ip": {
				Type:        schema.TypeString,
				Description: "The internal ip used to access a instance.",
				Computed:    true,
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
						},
						"zone_name": {
							Type:        schema.TypeString,
							Description: "Zone name of the subnet. e.g. `cn-bj-a`.",
							Optional:    true,
						},
					},
				},
			},
			"billing": {
				Type:        schema.TypeMap,
				Description: "**Deprecated**. Use `payment_timing`, `reservation_length`, `reservation_time_unit` instead. Billing information of the Scs.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"payment_timing": {
							Type:         schema.TypeString,
							Description:  "**Deprecated**. Payment timing of billing, which can be Prepaid or Postpaid. The default is Postpaid.",
							Required:     true,
							Default:      PaymentTimingPostpaid,
							ValidateFunc: validatePaymentTiming(),
						},
						"reservation": {
							Type:             schema.TypeMap,
							Description:      "**Deprecated**. Reservation of the Scs.",
							Optional:         true,
							DiffSuppressFunc: postPaidDiffSuppressFunc,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"reservation_length": {
										Type:             schema.TypeInt,
										Description:      "**Deprecated**. The reservation length that you will pay for your resource. It is valid when payment_timing is Prepaid. Valid values: [1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36].",
										Required:         true,
										Default:          1,
										ValidateFunc:     validateReservationLength(),
										DiffSuppressFunc: postPaidDiffSuppressFunc,
									},
									"reservation_time_unit": {
										Type:             schema.TypeString,
										Description:      "**Deprecated**. The reservation time unit that you will pay for your resource. It is valid when payment_timing is Prepaid. The value can only be month currently, which is also the default value.",
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
			"payment_timing": {
				Type:         schema.TypeString,
				Description:  "Payment timing of billing, Valid values: `Prepaid`, `Postpaid`.",
				Optional:     true,
				ValidateFunc: validatePaymentTiming(),
			},
			"reservation_length": {
				Type:             schema.TypeInt,
				Description:      "Prepaid billing reservation length, only useful when `payment_timing` is `Prepaid`. Valid values: `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `12`, `24`, `36`",
				Optional:         true,
				DiffSuppressFunc: postPaidDiffSuppressFunc,
				ValidateFunc:     validateReservationLength(),
			},
			"reservation_time_unit": {
				Type:             schema.TypeString,
				Description:      "Prepaid billing reservation time unit, only useful when `payment_timing` is `Prepaid`. Only support `month` now.",
				Optional:         true,
				DiffSuppressFunc: postPaidDiffSuppressFunc,
				ValidateFunc:     validateReservationUnit(),
			},
			"auto_renew_time_unit": {
				Type:        schema.TypeString,
				Description: "Time unit of automatic renewal, the value can be month or year. The default value is empty, indicating no automatic renewal. It is valid only when the payment_timing is Prepaid.",
				Computed:    true,
			},
			"auto_renew_time_length": {
				Type:        schema.TypeInt,
				Description: "The time length of automatic renewal. It is valid when payment_timing is Prepaid, and the value should be 1-9 when the auto_renew_time_unit is month and 1-3 when the auto_renew_time_unit is year.",
				Computed:    true,
			},
			"tags": tagsCreationSchema(),
			"auto_renew": {
				Type:        schema.TypeBool,
				Description: "Whether to automatically renew.",
				Computed:    true,
			},
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
			"capacity": {
				Type:        schema.TypeInt,
				Description: "Memory capacity(GB) of the instance.",
				Computed:    true,
			},
			"used_capacity": {
				Type:        schema.TypeInt,
				Description: "The amount of memory(GB) used by the instance.",
				Computed:    true,
			},
			"zone_names": {
				Type:        schema.TypeList,
				Description: "Zone name list",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"client_auth": {
				Type:        schema.TypeString,
				Description: "Access password of the instance. Should be 8-16 characters, and contains at least two types of letters, numbers and symbols. Allowed symbols include `$ ^ * ( ) _ + - =`.",
				Optional:    true,
				Sensitive:   true,
			},
			"store_type": {
				Type:        schema.TypeInt,
				Description: "Store type of the instance. Valid values: `0`(high performance memory), `1`(ssd local disk), `3`(capacity storage, only for PegaDB).",
				Optional:    true,
			},
			"enable_read_only": {
				Type:         schema.TypeInt,
				Description:  "Whether the copies are read only. Valid values: `1`(enabled), `2`(disabled). Defaults to `2`.",
				Optional:     true,
				Default:      2,
				ValidateFunc: validation.IntInSlice([]int{1, 2}),
			},
			"disk_flavor": {
				Type:         schema.TypeInt,
				Description:  "Storage size(GB) when use PegaDB. Must be between `50` and `160`",
				Optional:     true,
				ValidateFunc: validation.IntBetween(50, 160),
			},
			"disk_type": {
				Type:        schema.TypeString,
				Description: "Disk type of the instance. Valid values: `cloud_hp1`, `enhanced_ssd_pl1`.",
				Optional:    true,
			},
			"replication_info": {
				Type:        schema.TypeList,
				Description: "Replica info of the instance. Adding and removing replicas at same time in one operation is not supported.",
				Optional:    true,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"availability_zone": {
							Type:        schema.TypeString,
							Description: "Availability zone of the replica. e.g. `cn-bj-a`.",
							Required:    true,
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Description: "Subnet id of the replica.",
							Required:    true,
						},
						"is_master": {
							Type:        schema.TypeInt,
							Description: "Whether the replica is master node. Valid values: `1`(master node), `0`(slave node).",
							Required:    true,
						},
					},
				},
			},
			"replication_resize_type": {
				Type:         schema.TypeString,
				Description:  "Replica resize type. Must set when change `replication_info`. Valid values: `add`, `delete`.",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"add", "delete"}, false),
			},
			"security_ips": {
				Type:        schema.TypeSet,
				Description: "Security ips of the scs.",
				Optional:    true,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"backup_days": {
				Type: schema.TypeString,
				Description: "Identifies which days of the week the backup cycle is performed: Mon (Monday) " +
					"Tue (Tuesday) Wed (Wednesday) Thu (Thursday) Fri (Friday) Sat (Saturday) Sun (Sunday) " +
					"comma separated, the values are as follows: Sun,Mon,Tue,Wed,Thu,Fri,Sta. Note: Automatic backup is " +
					"only supported if the number of slave nodes is greater than 1",
				ValidateFunc: validation.StringInSlice([]string{"Mon", "Tue", "Wed",
					"Thu", "Fri", "Sat", "Sun"}, false),
				Optional: true,
			},
			"backup_time": {
				Type: schema.TypeString,
				Description: "Identifies when to perform backup in a day, UTC time (+8 is Beijing time) " +
					"value such as: 01:05:00",
				Optional: true,
			},
			"expire_day": {
				Type:        schema.TypeInt,
				Description: "Backup file expiration time, value such as: 3",
				Optional:    true,
			},
		},
	}
}

func resourceBaiduCloudScsCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	scsService := ScsService{client}

	createScsArgs, err := buildBaiduCloudScsArgs(d, meta)
	if err != nil {
		return WrapError(err)
	}

	action := "Create SCS Instance " + createScsArgs.InstanceName
	addDebug(action, createScsArgs)

	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithScsClient(func(scsClient *scs.Client) (interface{}, error) {
			return scsClient.CreateInstance(createScsArgs)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, raw)
		response, _ := raw.(*scs.CreateInstanceResult)
		d.SetId(response.InstanceIds[0])
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_scs", action, BCESDKGoERROR)
	}

	stateConf := buildStateConf(
		[]string{SCSStatusCreating, SCSStatusPrecreate},
		[]string{SCSStatusRunning},
		d.Timeout(schema.TimeoutCreate),
		scsService.InstanceStateRefresh(d.Id(), []string{
			SCSStatusPausing,
			SCSStatusPaused,
			SCSStatusDeleted,
			SCSStatusDeleting,
			SCSStatusFailed,
			SCSStatusModifying,
			SCSStatusModifyFailed,
			SCSStatusExpire,
		}),
	)
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_scs", action, BCESDKGoERROR)
	}
	err = updateInstanceSecurityIPs(d, meta, d.Id())
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_scs", action, BCESDKGoERROR)
	}

	err = setScsBackupPolicy(d, meta, d.Id())
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_scs", action, BCESDKGoERROR)
	}

	return resourceBaiduCloudScsRead(d, meta)
}

func resourceBaiduCloudScsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	scsService := ScsService{client: client}
	instanceID := d.Id()
	action := "Query SCS Instance " + instanceID

	raw, err := client.WithScsClient(func(scsClient *scs.Client) (interface{}, error) {
		return scsClient.GetInstanceDetail(instanceID)
	})

	addDebug(action, raw)

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			d.Set("scs", "")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_scs", action, BCESDKGoERROR)
	}

	result, _ := raw.(*scs.GetInstanceDetailResult)

	d.Set("instance_name", result.InstanceName)
	d.Set("instance_id", result.InstanceID)
	d.Set("cluster_type", result.ClusterType)
	d.Set("instance_status", result.InstanceStatus)
	d.Set("engine", result.Engine)
	d.Set("engine_version", result.EngineVersion)
	d.Set("v_net_ip", result.VnetIP)
	d.Set("domain", result.Domain)
	d.Set("port", result.Port)
	d.Set("create_time", result.InstanceCreateTime)
	d.Set("expire_time", result.InstanceExpireTime)
	d.Set("capacity", result.Capacity)
	d.Set("used_capacity", result.UsedCapacity)
	d.Set("payment_timing", result.PaymentTiming)
	d.Set("zone_names", result.ZoneNames)
	d.Set("vpc_id", result.VpcID)
	d.Set("subnets", transSubnetsToSchema(result.Subnets))
	d.Set("auto_renew", result.AutoRenew)
	d.Set("tags", flattenTagsToMap(result.Tags))
	d.Set("replication_info", transReplicationInfoToSchema(result.ReplicationInfo))
	d.Set("shard_num", result.ShardNum)
	ips, err := scsService.GetSecurityIPs(d.Id())
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_scs", action, BCESDKGoERROR)
	}
	d.Set("security_ips", ips)

	return nil
}

func transSubnetsToSchema(subnets []scs.Subnet) []map[string]string {
	subnetList := []map[string]string{}
	for _, subnet := range subnets {
		subnetMap := make(map[string]string)
		subnetMap["subnet_id"] = subnet.SubnetID
		subnetMap["zone_name"] = subnet.ZoneName
		subnetList = append(subnetList, subnetMap)
	}
	return subnetList
}

func transReplicationInfoToSchema(replicationInfo []scs.Replication) []map[string]interface{} {
	var schemaList []map[string]interface{}
	for _, replication := range replicationInfo {
		replicationMap := make(map[string]interface{})
		replicationMap["availability_zone"] = replication.AvailabilityZone
		replicationMap["subnet_id"] = replication.SubnetId
		replicationMap["is_master"] = replication.IsMaster
		schemaList = append(schemaList, replicationMap)
	}
	return schemaList
}

func transSchemaToReplicationInfo(schema []interface{}) []scs.Replication {
	replicationInfo := make([]scs.Replication, len(schema))
	for id := range schema {
		input := schema[id].(map[string]interface{})
		replication := scs.Replication{
			AvailabilityZone: input["availability_zone"].(string),
			SubnetId:         input["subnet_id"].(string),
			IsMaster:         input["is_master"].(int),
		}
		replicationInfo[id] = replication
	}
	return replicationInfo
}

func resourceBaiduCloudScsUpdate(d *schema.ResourceData, meta interface{}) error {
	instanceID := d.Id()

	d.Partial(true)

	// update instance name
	if err := updateScsInstanceName(d, meta, instanceID); err != nil {
		return err
	}

	// update instance nodeType/diskFlavor
	if err := updateInstanceNodeTypeAndDiskFlavor(d, meta, instanceID); err != nil {
		return err
	}

	// update instance shardNum
	if err := updateInstanceShardNum(d, meta, instanceID); err != nil {
		return err
	}

	// update instance replicationInfo
	if err := updateInstanceReplicationInfo(d, meta, instanceID); err != nil {
		return err
	}

	if err := updateInstanceSecurityIPs(d, meta, instanceID); err != nil {
		return err
	}

	// update back policy
	if err := setScsBackupPolicy(d, meta, instanceID); err != nil {
		return err
	}

	d.Partial(false)

	return resourceBaiduCloudScsRead(d, meta)
}

func resourceBaiduCloudScsDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	scsService := ScsService{client}

	instanceId := d.Id()
	action := "Delete SCS Instance " + instanceId

	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		raw, err := client.WithScsClient(func(scsClient *scs.Client) (interface{}, error) {
			return instanceId, scsClient.DeleteInstance(instanceId, buildClientToken())
		})
		if err != nil {
			if IsExceptedErrors(err, []string{InvalidInstanceStatus, bce.EINTERNAL_ERROR, ReleaseInstanceFailed}) {
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
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_scs", action, BCESDKGoERROR)
	}

	stateConf := buildStateConf(
		[]string{SCSStatusRunning,
			SCSStatusDeleting,
			SCSStatusPausing},
		[]string{SCSStatusPaused,
			SCSStatusDeleted,
			SCSStatusIsolated},
		d.Timeout(schema.TimeoutDelete),
		scsService.InstanceStateRefresh(instanceId, []string{}),
	)
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_scs", action, BCESDKGoERROR)
	}

	return nil
}

func buildBaiduCloudScsArgs(d *schema.ResourceData, meta interface{}) (*scs.CreateInstanceArgs, error) {
	request := &scs.CreateInstanceArgs{
		ClientToken: buildClientToken(),
	}

	// billing is deprecated
	if v, ok := d.GetOk("billing"); ok {
		billing := v.(map[string]interface{})
		billingRequest := scs.Billing{
			PaymentTiming: "",
			Reservation:   &scs.Reservation{},
		}
		if p, ok := billing["payment_timing"]; ok {
			paymentTiming := p.(string)
			billingRequest.PaymentTiming = paymentTiming
		}
		if billingRequest.PaymentTiming == PaymentTimingPrepaid {
			if r, ok := billing["reservation"]; ok {
				reservation := r.(map[string]interface{})
				if reservationLength, ok := reservation["reservation_length"]; ok {
					billingRequest.Reservation.ReservationLength = reservationLength.(int)
				}
				if reservationTimeUnit, ok := reservation["reservation_time_unit"]; ok {
					billingRequest.Reservation.ReservationTimeUnit = reservationTimeUnit.(string)
				}
			}
		}
		request.Billing = billingRequest
	}

	if paymentTiming, ok := d.GetOk("payment_timing"); ok {
		billingRequest := scs.Billing{
			PaymentTiming: paymentTiming.(string),
			Reservation:   &scs.Reservation{},
		}
		if billingRequest.PaymentTiming == PaymentTimingPrepaid {
			if length, ok := d.GetOk("reservation_length"); ok {
				billingRequest.Reservation.ReservationLength = length.(int)
			}
			if timeUnit, ok := d.GetOk("reservation_time_unit"); ok {
				billingRequest.Reservation.ReservationTimeUnit = timeUnit.(string)
			}
		}
		request.Billing = billingRequest
	}

	if request.Billing.PaymentTiming == "" {
		return nil, Error(InvalidInputField, "payment_timing")
	}

	if request.Billing.PaymentTiming == PaymentTimingPrepaid {
		// if the field is set, then auto-renewal is effective.
		if v, ok := d.GetOk("auto_renew_time_unit"); ok {
			request.AutoRenewTimeUnit = v.(string)
			if v, ok := d.GetOk("auto_renew_time_length"); ok {
				request.AutoRenewTime = v.(int)
			}
		}
	}

	if purchaseCount, ok := d.GetOk("purchase_count"); ok {
		request.PurchaseCount = purchaseCount.(int)
	}

	if instanceName, ok := d.GetOk("instance_name"); ok {
		request.InstanceName = instanceName.(string)
	}

	if nodeType, ok := d.GetOk("node_type"); ok {
		request.NodeType = nodeType.(string)
	}

	if shardNum, ok := d.GetOk("shard_num"); ok {
		request.ShardNum = shardNum.(int)
	}

	if proxyNum, ok := d.GetOk("proxy_num"); ok {
		request.ProxyNum = proxyNum.(int)
	}

	if clusterType, ok := d.GetOk("cluster_type"); ok {
		request.ClusterType = clusterType.(string)
	}

	if replicationNum, ok := d.GetOk("replication_num"); ok {
		request.ReplicationNum = replicationNum.(int)
	}

	if port, ok := d.GetOk("port"); ok {
		request.Port = port.(int)
	}

	if engineVersion, ok := d.GetOk("engine_version"); ok {
		request.EngineVersion = engineVersion.(string)
	}

	if vpcID, ok := d.GetOk("vpc_id"); ok {
		request.VpcID = vpcID.(string)
	}

	if v, ok := d.GetOk("subnets"); ok {
		subnetList := v.([]interface{})
		subnetRequests := make([]scs.Subnet, len(subnetList))
		for id := range subnetList {
			subnet := subnetList[id].(map[string]interface{})
			subnetRequest := scs.Subnet{
				SubnetID: subnet["subnet_id"].(string),
				ZoneName: subnet["zone_name"].(string),
			}
			subnetRequests[id] = subnetRequest
		}
		request.Subnets = subnetRequests
	}

	if engine, ok := d.GetOk("engine"); ok {
		request.Engine = SCSEngineIntegers()[engine.(string)]
	}

	if diskFlavor, ok := d.GetOk("disk_flavor"); ok {
		request.DiskFlavor = diskFlavor.(int)
	}

	if diskType, ok := d.GetOk("disk_type"); ok {
		request.DiskType = diskType.(string)
	}

	if clientAuth, ok := d.GetOk("client_auth"); ok {
		request.ClientAuth = clientAuth.(string)
	}

	if storeType, ok := d.GetOk("store_type"); ok {
		request.StoreType = storeType.(int)
	}

	if enableReadOnly, ok := d.GetOk("enable_read_only"); ok {
		request.EnableReadOnly = enableReadOnly.(int)
	}

	if info, ok := d.GetOk("replication_info"); ok {
		inputList := info.([]interface{})
		request.ReplicationInfo = transSchemaToReplicationInfo(inputList)
	}

	if tags, ok := d.GetOk("tags"); ok {
		request.Tags = tranceTagMapToModel(tags.(map[string]interface{}))
	}

	return request, nil
}

func updateScsInstanceName(d *schema.ResourceData, meta interface{}, instanceID string) error {
	action := "Update scs instanceName " + instanceID
	client := meta.(*connectivity.BaiduClient)

	if d.HasChange("instance_name") {
		args := &scs.UpdateInstanceNameArgs{
			InstanceName: d.Get("instance_name").(string),
			ClientToken:  buildClientToken(),
		}

		addDebug(action, args)
		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			_, err := client.WithScsClient(func(scsClient *scs.Client) (interface{}, error) {
				return nil, scsClient.UpdateInstanceName(instanceID, args)
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
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_scs", action, BCESDKGoERROR)
		}
		d.SetPartial("instance_name")
	}

	return nil
}

func updateInstanceNodeTypeAndDiskFlavor(d *schema.ResourceData, meta interface{}, instanceID string) error {
	action := "Update scs nodeType/diskFlavor " + instanceID
	client := meta.(*connectivity.BaiduClient)
	scsService := ScsService{client}

	if d.HasChange("node_type") || (d.HasChange("disk_flavor") && "PegaDB" == d.Get("engine").(string)) {
		shardNum := d.Get("shard_num")
		if d.HasChange("shard_num") {
			shardNum, _ = d.GetChange("shard_num")
		}

		args := &scs.ResizeInstanceArgs{
			ClientToken: buildClientToken(),
			NodeType:    d.Get("node_type").(string),
			DiskFlavor:  d.Get("disk_flavor").(int),
			ShardNum:    shardNum.(int),
		}

		addDebug(action, args)
		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			_, err := client.WithScsClient(func(scsClient *scs.Client) (interface{}, error) {
				return nil, scsClient.ResizeInstance(instanceID, args)
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
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_scs", action, BCESDKGoERROR)
		}

		stateConf := buildStateConf(
			[]string{SCSStatusModifying},
			[]string{SCSStatusRunning},
			d.Timeout(schema.TimeoutUpdate),
			scsService.InstanceStateRefresh(d.Id(), []string{}),
		)
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_scs", action, BCESDKGoERROR)
		}

		d.SetPartial("node_type")
		d.SetPartial("disk_flavor")
	}

	return nil
}

func updateInstanceShardNum(d *schema.ResourceData, meta interface{}, instanceID string) error {
	action := "Update scs shardNum " + instanceID
	client := meta.(*connectivity.BaiduClient)
	scsService := ScsService{client}

	if d.HasChange("shard_num") && "cluster" == d.Get("cluster_type").(string) {
		args := &scs.ResizeInstanceArgs{
			ClientToken: buildClientToken(),
			NodeType:    d.Get("node_type").(string),
			DiskFlavor:  d.Get("disk_flavor").(int),
			ShardNum:    d.Get("shard_num").(int),
		}

		addDebug(action, args)
		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			_, err := client.WithScsClient(func(scsClient *scs.Client) (interface{}, error) {
				return nil, scsClient.ResizeInstance(instanceID, args)
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
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_scs", action, BCESDKGoERROR)
		}

		stateConf := buildStateConf(
			[]string{SCSStatusModifying},
			[]string{SCSStatusRunning},
			d.Timeout(schema.TimeoutCreate),
			scsService.InstanceStateRefresh(d.Id(), []string{}),
		)
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_scs", action, BCESDKGoERROR)
		}
		d.SetPartial("shard_num")
	}

	return nil
}

func updateInstanceReplicationInfo(d *schema.ResourceData, meta interface{}, instanceID string) error {
	action := "Update scs replicationInfo " + instanceID
	client := meta.(*connectivity.BaiduClient)
	scsService := ScsService{client}

	if d.HasChange("replication_info") {
		resizeType, ok := d.GetOk("replication_resize_type")
		if !ok {
			return Error(InvalidInputField, "replication_resize_type")
		}
		isAddReplication := strings.HasPrefix(resizeType.(string), "add")

		args := &scs.ReplicationArgs{
			ResizeType:      d.Get("replication_resize_type").(string),
			ReplicationInfo: transSchemaToReplicationInfo(d.Get("replication_info").([]interface{})),
		}

		addDebug(action, args)
		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			_, err := client.WithScsClient(func(scsClient *scs.Client) (interface{}, error) {
				if isAddReplication {
					return nil, scsClient.AddReplication(instanceID, args)
				}
				return nil, scsClient.DeleteReplication(instanceID, args)
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
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_scs", action, BCESDKGoERROR)
		}

		stateConf := buildStateConf([]string{SCSStatusModifying}, []string{SCSStatusRunning}, d.Timeout(schema.TimeoutCreate), scsService.InstanceStateRefresh(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_scs", action, BCESDKGoERROR)
		}
		d.SetPartial("replication_info")
	}

	return nil
}

func updateInstanceSecurityIPs(d *schema.ResourceData, meta interface{}, instanceID string) error {
	action := "Update scs security ips " + instanceID
	client := meta.(*connectivity.BaiduClient)
	scsService := ScsService{
		client: client,
	}
	ips, err := scsService.GetSecurityIPs(d.Id())
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_scs", action, BCESDKGoERROR)
	}
	os := &schema.Set{
		F: schema.HashString,
	}
	for _, ip := range ips {
		os.Add(ip)
	}
	ns := d.Get("security_ips").(*schema.Set)
	addIPs := ns.Difference(os).List()
	deleteIPs := os.Difference(ns).List()

	addIPsArg := make([]string, 0)
	for _, ips := range addIPs {
		addIPsArg = append(addIPsArg, ips.(string))
	}
	// Add security IPs
	if _, err := client.WithScsClient(func(scsClient *scs.Client) (i interface{}, e error) {
		return nil, scsClient.AddSecurityIp(instanceID, &scs.SecurityIpArgs{
			SecurityIps: addIPsArg,
		})
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_scs", action, BCESDKGoERROR)
	}

	deleteIPsArg := make([]string, 0)
	for _, ips := range deleteIPs {
		deleteIPsArg = append(deleteIPsArg, ips.(string))
	}
	// Delete security IPs
	if _, err := client.WithScsClient(func(scsClient *scs.Client) (i interface{}, e error) {
		return nil, scsClient.DeleteSecurityIp(instanceID, &scs.SecurityIpArgs{
			SecurityIps: deleteIPsArg,
		})
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_scs", action, BCESDKGoERROR)
	}
	return nil
}

func setScsBackupPolicy(d *schema.ResourceData, meta interface{}, instanceID string) error {
	action := "Set scs backup policy " + instanceID
	client := meta.(*connectivity.BaiduClient)

	if d.HasChange("backup_days") || d.HasChange("backup_time") || d.HasChange("expire_in_days") {
		args := &scs.ModifyBackupPolicyArgs{
			BackupDays: d.Get("backup_days").(string),
			BackupTime: d.Get("backup_time").(string),
			ExpireDay:  d.Get("expire_day").(int),
		}

		addDebug(action, args)
		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			_, err := client.WithScsClient(func(scsClient *scs.Client) (interface{}, error) {
				return nil, scsClient.ModifyBackupPolicy(instanceID, args)
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
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_scs_instance", action, BCESDKGoERROR)
		}

		d.SetPartial("backup_days")
		d.SetPartial("backup_time")
		d.SetPartial("expire_day")
	}

	return nil
}
