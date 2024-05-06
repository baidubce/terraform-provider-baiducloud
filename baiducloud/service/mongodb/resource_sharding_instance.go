package mongodb

import (
	"fmt"
	"log"
	"time"

	"github.com/baidubce/bce-sdk-go/services/mongodb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/flex"
)

func ResourceShardingInstance() *schema.Resource {
	fullSchema := map[string]*schema.Schema{
		"mongos_count": {
			Type:         schema.TypeInt,
			Description:  "Mongos nodes count of the instance. Valid values: `2`~`32`.",
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(2, 32),
		},
		"shard_count": {
			Type:         schema.TypeInt,
			Description:  "Shard nodes count of the instance. Valid values: `2`~`500`.",
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(2, 500),
		},
		"mongos_cpu_count":       schemaCPUCount(),
		"mongos_memory_capacity": schemaMemoryCapacity(),
		"shard_cpu_count":        schemaCPUCount(),
		"shard_memory_capacity":  schemaMemoryCapacity(),
		"shard_storage":          schemaStorage(),
		"shard_storage_type":     schemaStorageType(),
		// computed
		"mongos_list": {
			Type:        schema.TypeList,
			Description: "Mongos node list of the instance.",
			Computed:    true,
			Elem:        nodeSchema(),
		},
		"shard_list": {
			Type:        schema.TypeList,
			Description: "Shard node list of the instance.",
			Computed:    true,
			Elem:        nodeSchema(),
		},
	}
	flex.MergeSchema(fullSchema, basicResourceInstanceSchema())

	return &schema.Resource{
		Description: "Use this resource to manage MongoDB Sharding Instance. \n\n" +
			"More information can be found in the [Developer Guide](https://cloud.baidu.com/doc/MONGODB/s/ikdgsphbp). \n\n",

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Create: resourceShardingInstanceCreate,
		Read:   resourceShardingInstanceRead,
		Update: resourceShardingInstanceUpdate,
		Delete: resourceInstanceDelete,

		Schema: fullSchema,
	}
}

func nodeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"node_id": {
				Type:        schema.TypeString,
				Description: "ID of the node.",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the node.",
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Status of the node. Possible values: `CREATING`, `RUNNING`, `RESTARTING`, `CLASS_CHANGING`.",
				Computed:    true,
			},
			"cpu_count": {
				Type:        schema.TypeInt,
				Description: "CPU core count of the node.",
				Computed:    true,
			},
			"memory_capacity": {
				Type:        schema.TypeInt,
				Description: "Memory size (GB) of the node.",
				Computed:    true,
			},
			"storage": {
				Type:        schema.TypeInt,
				Description: "Storage size (GB) of the node.",
				Computed:    true,
			},
			"storage_type": {
				Type:        schema.TypeString,
				Description: "Storage type of the node. Possible values: `CDS_PREMIUM_SSD`, `CDS_ENHANCED_SSD`, `LOCAL_DISK`.",
				Computed:    true,
			},
			"connection_string": {
				Type:        schema.TypeString,
				Description: "Connection address of the node.",
				Computed:    true,
			},
		},
	}
}

func resourceShardingInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)
	raw, err := conn.WithMongoDBClient(func(client *mongodb.Client) (interface{}, error) {
		return client.CreateSharding(buildShardingCreationArgs(d))
	})
	log.Printf("[DEBUG] Create MongoDB Sharding Instance result: %+v", raw)
	if err != nil {
		return fmt.Errorf("error creating MongoDB Sharding Instance: %w", err)
	}
	response := raw.(*mongodb.CreateResult)
	if response.DbInstanceSimpleModels == nil || len(response.DbInstanceSimpleModels) == 0 {
		return fmt.Errorf("error creating MongoDB Sharding Instance: %+v", raw)
	}

	instance := response.DbInstanceSimpleModels[0]
	d.SetId(instance.DbInstanceId)

	time.Sleep(60 * time.Second)
	if _, err = waitInstanceAvailable(conn, d.Id()); err != nil {
		return fmt.Errorf("error waiting MongoDB Instance (%s) becoming available: %w", d.Id(), err)
	}

	return resourceShardingInstanceRead(d, meta)
}

func resourceShardingInstanceRead(d *schema.ResourceData, meta interface{}) error {
	detail, err := basicInstanceSchemaRead(d, meta)
	if err != nil {
		return err
	}
	if err := d.Set("mongos_count", detail.MongosCount); err != nil {
		return fmt.Errorf("error setting mongos_count: %w", err)
	}
	if err := d.Set("shard_count", detail.ShardCount); err != nil {
		return fmt.Errorf("error setting shard_count: %w", err)
	}
	if err := d.Set("mongos_list", flattenNodeList(detail.MongosList)); err != nil {
		return fmt.Errorf("error setting mongos_list: %w", err)
	}
	if err := d.Set("shard_list", flattenNodeList(detail.ShardList)); err != nil {
		return fmt.Errorf("error setting shard_list: %w", err)
	}
	return nil
}

func resourceShardingInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)

	if err := updateName(d, conn); err != nil {
		return fmt.Errorf("error updating MongoDB Sharding Instance (%s) name: %w", d.Id(), err)
	}
	if err := updatePassword(d, conn); err != nil {
		return fmt.Errorf("error updating MongoDB Sharding Instance (%s) password: %w", d.Id(), err)
	}
	return resourceShardingInstanceRead(d, meta)
}

func buildShardingCreationArgs(d *schema.ResourceData) *mongodb.CreateShardingArgs {
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
	args := &mongodb.CreateShardingArgs{
		Billing:         billing,
		DbInstanceName:  d.Get("name").(string),
		StorageEngine:   d.Get("storage_engine").(string),
		EngineVersion:   d.Get("engine_version").(string),
		DbInstanceType:  mongodb.S_SHARDING,
		AccountPassword: d.Get("account_password").(string),
		VpcId:           d.Get("vpc_id").(string),
		Subnets:         expandSubnets(d.Get("subnets").([]interface{})),
		Tags:            expandTags(d.Get("tags").(map[string]interface{})),
		ResGroupId:      d.Get("resource_group_id").(string),

		MongosCount:          d.Get("mongos_count").(int),
		MongosCpuCount:       d.Get("mongos_cpu_count").(int),
		MongosMemoryCapacity: d.Get("mongos_memory_capacity").(int),
		ShardCount:           d.Get("shard_count").(int),
		ShardCpuCount:        d.Get("shard_cpu_count").(int),
		ShardMemoryCapacity:  d.Get("shard_memory_capacity").(int),
		ShardStorage:         d.Get("shard_storage").(int),
		ShardStorageType:     d.Get("shard_storage_type").(string),
	}
	return args
}
