package mongodb

import (
	"fmt"
	"log"

	"github.com/baidubce/bce-sdk-go/services/mongodb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/flex"
)

func DataSourceInstances() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to query MongoDB instance list. \n\n",

		Read: dataSourceInstancesRead,

		Schema: map[string]*schema.Schema{
			"type": {
				Type:         schema.TypeString,
				Description:  "Type of the instance. Valid values: `replica`, `sharding`.",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{mongodb.S_REPLICA, mongodb.S_SHARDING}, false),
			},
			"storage_engine": {
				Type:         schema.TypeString,
				Description:  "Storage engine of the instance. Valid values: `WiredTiger`.",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"WiredTiger"}, false),
			},
			"engine_version": {
				Type:         schema.TypeString,
				Description:  "Database version of the instance. Valid values: `3.4`, `3.6`.",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"3.4", "3.6"}, false),
			},
			// Computed
			"instance_list": {
				Type:        schema.TypeList,
				Description: "Instance list.",
				Computed:    true,
				Elem:        schemaInstance(),
			},
		},
	}
}

func schemaInstance() *schema.Resource {
	instanceSchema := map[string]*schema.Schema{
		"instance_id": {
			Type:        schema.TypeString,
			Description: "ID of the instance.",
			Computed:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Name of the instance.",
			Computed:    true,
		},
		"payment_timing": flex.ComputedSchemaPaymentTiming(),
		"vpc_id":         flex.ComputedSchemaVpcID(),
		"subnets":        flex.ComputedSchemaSubnets(),
		"tags":           flex.ComputedSchemaTags(),
		"type": {
			Type:        schema.TypeString,
			Description: "Type of the instance. Possible values: `replica`, `sharding`.",
			Computed:    true,
		},
		"storage_engine": {
			Type:        schema.TypeString,
			Description: "Storage engine of the instance. Possible values: `WiredTiger`.",
			Computed:    true,
		},
		"engine_version": {
			Type:        schema.TypeString,
			Description: "Database version of the instance. Possible values: `3.4`, `3.6`.",
			Computed:    true,
		},
		// type replica
		"cpu_count": {
			Type:        schema.TypeInt,
			Description: "CPU core count of the instance.",
			Computed:    true,
		},
		"memory_capacity": {
			Type:        schema.TypeInt,
			Description: "Memory size (GB) of the instance.",
			Computed:    true,
		},
		"storage": {
			Type:        schema.TypeInt,
			Description: "Storage size (GB) of the instance.",
			Computed:    true,
		},
		"voting_member_num": {
			Type:        schema.TypeInt,
			Description: "Number of voting nodes in the instance. Possible values: `1`~`3`.",
			Computed:    true,
		},
		"readonly_node_num": {
			Type:        schema.TypeInt,
			Description: "Number of readonly nodes in the instance. Possible values: `0`~`5`.",
			Computed:    true,
		},
		"port": {
			Type:        schema.TypeString,
			Description: "Connection port of the instance.",
			Computed:    true,
		},
		// type sharding
		"mongos_count": {
			Type:        schema.TypeInt,
			Description: "Number of mongos nodes in the sharding instance.",
			Computed:    true,
		},
		"shard_count": {
			Type:        schema.TypeInt,
			Description: "Number of shard nodes in the sharding instance.",
			Computed:    true,
		},
	}
	flex.MergeSchema(instanceSchema, basicComputedOnlySchema())
	return &schema.Resource{
		Schema: instanceSchema,
	}
}

func dataSourceInstancesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)

	args := mongodb.ListMongodbArgs{
		DbInstanceType: d.Get("type").(string),
		StorageEngine:  d.Get("storage_engine").(string),
		EngineVersion:  d.Get("engine_version").(string),
	}

	instances, err := findAllInstance(conn, args)

	log.Printf("[DEBUG] Read MongoDB instance list result: %+v", instances)
	if err != nil {
		return fmt.Errorf("error reading MongoDB instance list: %w, %+v", err, args)
	}

	if err := d.Set("instance_list", flattenInstanceList(instances)); err != nil {
		return fmt.Errorf("error setting instance_list: %w", err)
	}

	d.SetId(resource.UniqueId())
	return nil
}
