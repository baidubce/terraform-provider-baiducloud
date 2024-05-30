package mongodb

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/baidubce/bce-sdk-go/services/mongodb"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/flex"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func ResourceInstance() *schema.Resource {
	fullSchema := map[string]*schema.Schema{
		"cpu_count":       schemaCPUCount(),
		"memory_capacity": schemaMemoryCapacity(),
		"storage":         schemaStorage(),
		"storage_type":    schemaStorageType(),
		"voting_member_num": {
			Type:         schema.TypeInt,
			Description:  "Number of voting nodes in the instance. Valid values: `1`~`3`. Defaults to `3`.",
			Optional:     true,
			Default:      3,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1, 3),
		},
		"readonly_node_num": {
			Type: schema.TypeInt,
			Description: "Number of readonly nodes in the instance. Only effective when `voting_member_num` is set to `2` or `3`. " +
				"Valid values: `0`~`5`. Defaults to `0`.",
			Optional:     true,
			Default:      0,
			ValidateFunc: validation.IntBetween(0, 5),
		},
		"port": {
			Type:        schema.TypeString,
			Description: "Connection port of the instance.",
			Computed:    true,
		},
		"auto_backup_enable": {
			Type:         schema.TypeString,
			Default:      "OFF",
			Description:  "Auto backup status.Valid values: 'ON','OFF'. ",
			ValidateFunc: validation.StringInSlice([]string{"ON", "OFF"}, false),
			Optional:     true,
		},
		"preferred_backup_period": {
			Type: schema.TypeSet,
			Description: "Backup period. Value reference: Monday, Tuesday, Wednesday," +
				" Thursday, Friday, Saturday, Sunday",
			DiffSuppressFunc: backupPolicyDiffSuppressFunc,
			Optional:         true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"Monday", "Tuesday", "Wednesday",
					"Thursday", "Friday", "Saturday", "Sunday"}, false),
			},
			Set: schema.HashString,
		},
		"preferred_backup_time": {
			Type:             schema.TypeString,
			Description:      "Backup time. The format is HH:mmZ-HH:mmZ. The time range is limited to 1 hour.",
			DiffSuppressFunc: backupPolicyDiffSuppressFunc,
			Optional:         true,
		},
		"enable_increment_backup": {
			Type:             schema.TypeInt,
			Description:      "Whether incremental backup is enabled. 0: disabled; 1: enabled",
			DiffSuppressFunc: backupPolicyDiffSuppressFunc,
			ValidateFunc:     validation.IntBetween(0, 1),
			Optional:         true,
		},
		//"backup_method": {
		//	Type:     schema.TypeString,
		//	Computed: true,
		//},
		//"incr_backup_retention_period": {
		//	Type:     schema.TypeInt,
		//	Computed: true,
		//},
		"backup_retention_period": {
			Type:        schema.TypeInt,
			Description: "Number of days to keep backups",
			Computed:    true,
		},
	}
	flex.MergeSchema(fullSchema, basicResourceInstanceSchema())

	return &schema.Resource{
		Description: "Use this resource to manage MongoDB Instance. \n\n" +
			"More information can be found in the [Developer Guide](https://cloud.baidu.com/doc/MONGODB/s/Ekdgskkrk). \n\n",

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Create: resourceInstanceCreate,
		Read:   resourceInstanceRead,
		Update: resourceInstanceUpdate,
		Delete: resourceInstanceDelete,

		Schema: fullSchema,
	}
}

func resourceInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)
	raw, err := conn.WithMongoDBClient(func(client *mongodb.Client) (interface{}, error) {
		return client.CreateReplica(buildCreationArgs(d))
	})
	log.Printf("[DEBUG] Create MongoDB Instance result: %+v", raw)
	if err != nil {
		return fmt.Errorf("error creating MongoDB Instance: %w", err)
	}
	response := raw.(*mongodb.CreateResult)
	if response.DbInstanceSimpleModels == nil || len(response.DbInstanceSimpleModels) == 0 {
		return fmt.Errorf("error creating MongoDB Instance: %+v", raw)
	}

	instance := response.DbInstanceSimpleModels[0]
	d.SetId(instance.DbInstanceId)

	time.Sleep(60 * time.Second)
	if _, err = waitInstanceAvailable(conn, d.Id()); err != nil {
		return fmt.Errorf("error waiting MongoDB Instance (%s) becoming available: %w", d.Id(), err)
	}
	if err = updateSecurityIps(d, meta); err != nil {
		return fmt.Errorf("error set MongoDB Instance (%s) Security ips : %w", d.Id(), err)
	}
	if err = updateBackupPolicy(d, meta); err != nil {
		return fmt.Errorf("error set MongoDB Instance (%s) Backup Policy : %w", d.Id(), err)
	}
	return resourceInstanceRead(d, meta)
}

func resourceInstanceRead(d *schema.ResourceData, meta interface{}) error {
	detail, err := basicInstanceSchemaRead(d, meta)
	if err != nil {
		return err
	}
	if err := d.Set("cpu_count", detail.DbInstanceCpuCount); err != nil {
		return fmt.Errorf("error setting cpu_count: %w", err)
	}
	if err := d.Set("memory_capacity", detail.DbInstanceMemoryCapacity); err != nil {
		return fmt.Errorf("error setting memory_capacity: %w", err)
	}
	if err := d.Set("storage", detail.DbInstanceStorage); err != nil {
		return fmt.Errorf("error setting storage: %w", err)
	}
	if err := d.Set("storage_type", detail.DbInstanceStorageType); err != nil {
		return fmt.Errorf("error setting storage_type: %w", err)
	}
	if err := d.Set("voting_member_num", detail.VotingMemberNum); err != nil {
		return fmt.Errorf("error setting voting_member_num: %w", err)
	}
	if err := d.Set("readonly_node_num", detail.ReadonlyNodeNum); err != nil {
		return fmt.Errorf("error setting readonly_node_num: %w", err)
	}
	if err := d.Set("port", detail.Port); err != nil {
		return fmt.Errorf("error setting port: %w", err)
	}
	backupPolicy, err := basicInstanceBackupPolicyRead(d, meta)
	if err != nil {
		return err
	}
	if err := d.Set("auto_backup_enable", backupPolicy.AutoBackupEnable); err != nil {
		return fmt.Errorf("error setting auto_backup_enable: %w", err)
	}
	if err := d.Set("backup_retention_period", backupPolicy.BackupRetentionPeriod); err != nil {
		return fmt.Errorf("error setting backup_retention_period: %w", err)
	}
	if err := d.Set("preferred_backup_period", strings.Split(backupPolicy.PreferredBackupPeriod, ",")); err != nil {
		return fmt.Errorf("error setting preferred_backup_period: %w", err)
	}
	if err := d.Set("preferred_backup_time", backupPolicy.PreferredBackupTime); err != nil {
		return fmt.Errorf("error setting preferred_backup_time: %w", err)
	}
	if err := d.Set("enable_increment_backup", backupPolicy.EnableIncrementBackup); err != nil {
		return fmt.Errorf("error setting enable_increment_backup: %w", err)
	}
	//if err := d.Set("backup_method", backupPolicy.BackupMethod); err != nil {
	//	return fmt.Errorf("error setting backup_method: %w", err)
	//}
	//if err := d.Set("incr_backup_retention_period", backupPolicy.IncrBackupRetentionPeriod); err != nil {
	//	return fmt.Errorf("error setting incr_backup_retention_period: %w", err)
	//}
	return nil
}

func resourceInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)

	if err := updateName(d, conn); err != nil {
		return fmt.Errorf("error updating MongoDB Instance (%s) name: %w", d.Id(), err)
	}
	if err := updatePassword(d, conn); err != nil {
		return fmt.Errorf("error updating MongoDB Instance (%s) password: %w", d.Id(), err)
	}
	if err := updateSecurityIps(d, meta); err != nil {
		return fmt.Errorf("error updating MongoDB Instance (%s) security ips: %w", d.Id(), err)
	}
	if err := updateBackupPolicy(d, meta); err != nil {
		return fmt.Errorf("error updating MongoDB Instance (%s) Backup Policy : %w", d.Id(), err)
	}
	if err := resizeInstance(d, conn); err != nil {
		return err
	}
	return resourceInstanceRead(d, meta)
}

func buildCreationArgs(d *schema.ResourceData) *mongodb.CreateReplicaArgs {
	billing := mongodb.BillingModel{
		PaymentTiming: d.Get("payment_timing").(string),
	}
	if billing.PaymentTiming == flex.PaymentTimingPrepaid {
		reservation := mongodb.Reservation{
			ReservationLength:   d.Get("reservation_length").(int),
			ReservationTimeUnit: "month",
		}
		billing.Reservation = reservation
		if v, ok := d.GetOk("auto_renew_length"); ok {
			autoRenew := mongodb.AutoRenewModel{
				AutoRenewLength:   v.(int),
				AutoRenewTimeUnit: "month",
			}
			billing.AutoRenew = autoRenew
		}
	}
	args := &mongodb.CreateReplicaArgs{
		Billing:         billing,
		DbInstanceName:  d.Get("name").(string),
		StorageEngine:   d.Get("storage_engine").(string),
		EngineVersion:   d.Get("engine_version").(string),
		DbInstanceType:  mongodb.S_REPLICA,
		AccountPassword: d.Get("account_password").(string),
		VpcId:           d.Get("vpc_id").(string),
		Subnets:         expandSubnets(d.Get("subnets").([]interface{})),
		Tags:            expandTags(d.Get("tags").(map[string]interface{})),
		ResGroupId:      d.Get("resource_group_id").(string),

		DbInstanceCpuCount:       d.Get("cpu_count").(int),
		DbInstanceMemoryCapacity: d.Get("memory_capacity").(int),
		DbInstanceStorage:        d.Get("storage").(int),
		DbInstanceStorageType:    d.Get("storage_type").(string),
		VotingMemberNum:          d.Get("voting_member_num").(int),
		ReadonlyNodeNum:          d.Get("readonly_node_num").(int),
	}
	return args
}

func resizeInstance(d *schema.ResourceData, conn *connectivity.BaiduClient) error {
	if d.HasChanges("cpu_count", "memory_capacity", "storage") {
		args := &mongodb.ReplicaResizeArgs{
			DbInstanceCpuCount:       d.Get("cpu_count").(int),
			DbInstanceMemoryCapacity: d.Get("memory_capacity").(int),
			DbInstanceStorage:        d.Get("storage").(int),
		}
		log.Printf("[DEBUG] Resize MongoDB Instance (%s): %+v", d.Id(), args)
		_, err := conn.WithMongoDBClient(func(client *mongodb.Client) (interface{}, error) {
			return nil, client.ReplicaResize(d.Id(), args)
		})
		if err != nil {
			return fmt.Errorf("error resing MongoDB Instance (%s): %w", d.Id(), err)
		}

		time.Sleep(60 * time.Second)
		_, err = waitInstanceAvailable(conn, d.Id())
		if err != nil {
			return fmt.Errorf("error waiting MongoDB Instance (%s) becoming available after resizing: %w", d.Id(), err)
		}
	}
	return nil
}
