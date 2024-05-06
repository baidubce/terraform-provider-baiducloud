package mongodb

import (
	"fmt"
	"log"

	"github.com/baidubce/bce-sdk-go/services/mongodb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/flex"
)

func basicResourceInstanceSchema() map[string]*schema.Schema {
	basicSchema := map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "Name of the instance. If not specified, it will be randomly generated.",
			Optional:    true,
			Computed:    true,
		},
		"payment_timing":     flex.SchemaPaymentTiming(),
		"reservation_length": flex.SchemaReservationLength(),
		"auto_renew_length":  flex.SchemaAutoRenewLength(),
		//"auto_renew_time_unit": flex.SchemaAutoRenewTimeUnit(),
		"vpc_id":            flex.SchemaVpcID(),
		"subnets":           flex.SchemaSubnets(),
		"tags":              flex.SchemaTagsOnlySupportCreation(),
		"resource_group_id": flex.SchemaResourceGroupID(),
		"storage_engine": {
			Type:         schema.TypeString,
			Description:  "Storage engine of the instance. Valid values: `WiredTiger`.",
			Optional:     true,
			Default:      "WiredTiger",
			ValidateFunc: validation.StringInSlice([]string{"WiredTiger"}, false),
		},
		"engine_version": {
			Type:         schema.TypeString,
			Description:  "Database version of the instance. Valid values: `3.4`, `3.6`.",
			Optional:     true,
			Default:      "3.4",
			ValidateFunc: validation.StringInSlice([]string{"3.4", "3.6"}, false),
		},
		"account_password": {
			Type: schema.TypeString,
			Description: "Password for root account. If not specified, it will be randomly generated. " +
				"Must be 8-32 characters, including letters, numbers, and symbols(`!#$%^*()`only).",
			Optional:  true,
			Sensitive: true,
		},
		"expire_time": {
			Type:        schema.TypeString,
			Description: "Expiration time of the prepaid instance.",
			Computed:    true,
		},
	}
	flex.MergeSchema(basicSchema, basicComputedOnlySchema())
	return basicSchema
}

func basicComputedOnlySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"status": {
			Type: schema.TypeString,
			Description: "Status of the instance. Possible values: `CREATING`, `RUNNING`, `STOPPING`, `EXPIRED`, `RESTARTING`, " +
				"`STARTING`, `CLASS_CHANGING`, `NODE_RESTARTING`, `NODE_CREATING`, `NODE_CLASS_CHANGING`.",
			Computed: true,
		},
		"connection_string": {
			Type:        schema.TypeString,
			Description: "Connection address of the instance.",
			Computed:    true,
		},
		"create_time": {
			Type:        schema.TypeString,
			Description: "Creation time of the instance.",
			Computed:    true,
		},
	}
}

func schemaCPUCount() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeInt,
		Description:  "CPU core count. At least 1 core.",
		Required:     true,
		ValidateFunc: validation.IntInSlice([]int{1, 2, 4, 8, 16}),
	}
}

func schemaMemoryCapacity() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeInt,
		Description:  "Memory size (GB). At least 2 GB.",
		Required:     true,
		ValidateFunc: validation.IntInSlice([]int{2, 4, 8, 16, 32, 64}),
	}
}

func schemaStorage() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeInt,
		Description:  "Storage size (GB). At least 5 GB.",
		Required:     true,
		ValidateFunc: validation.IntAtLeast(5),
	}
}

func schemaStorageType() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Description:  "Storage type of the instance. Valid values: `CDS_PREMIUM_SSD`, `CDS_ENHANCED_SSD`, `LOCAL_DISK`. Defaults to `CDS_PREMIUM_SSD`.",
		Optional:     true,
		Default:      StorageTypeSSD,
		ValidateFunc: validation.StringInSlice([]string{StorageTypeSSD, StorageTypeEnhancedSSD, StorageTypeLocal}, false),
	}
}

func resourceInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)

	_, err := conn.WithMongoDBClient(func(client *mongodb.Client) (interface{}, error) {
		return nil, client.ReleaseMongodb(d.Id())
	})
	log.Printf("[DEBUG] Delete MongoDB Instance (%s).", d.Id())

	if err != nil {
		return fmt.Errorf("error delete MongoDB Instance (%s): %w", d.Id(), err)
	}
	return nil
}

func basicInstanceSchemaRead(d *schema.ResourceData, meta interface{}) (*mongodb.InstanceDetail, error) {
	conn := meta.(*connectivity.BaiduClient)
	detail, err := findInstance(conn, d.Id())
	log.Printf("[DEBUG] Read MongoDB Instance (%s) result: %+v", d.Id(), detail)
	if err != nil {
		return nil, fmt.Errorf("error reading MongoDB Instance (%s): %w", d.Id(), err)
	}

	if err := d.Set("name", detail.DbInstanceName); err != nil {
		return nil, fmt.Errorf("error setting name: %w", err)
	}
	if err := d.Set("payment_timing", detail.PaymentTiming); err != nil {
		return nil, fmt.Errorf("error setting payment_timing: %w", err)
	}
	if err := d.Set("vpc_id", detail.VpcId); err != nil {
		return nil, fmt.Errorf("error setting vpc_id: %w", err)
	}
	if err := d.Set("subnets", flattenSubnets(detail.Subnets)); err != nil {
		return nil, fmt.Errorf("error setting subnets: %w", err)
	}
	if err := d.Set("tags", flattenTags(detail.Tags)); err != nil {
		return nil, fmt.Errorf("error setting tags: %w", err)
	}
	if err := d.Set("storage_engine", detail.StorageEngine); err != nil {
		return nil, fmt.Errorf("error setting storage_engine: %w", err)
	}
	if err := d.Set("engine_version", detail.EngineVersion); err != nil {
		return nil, fmt.Errorf("error setting engine_version: %w", err)
	}
	if err := d.Set("status", detail.DbInstanceStatus); err != nil {
		return nil, fmt.Errorf("error setting status: %w", err)
	}
	if err := d.Set("connection_string", detail.ConnectionString); err != nil {
		return nil, fmt.Errorf("error setting connection_string: %w", err)
	}

	layout := "2006-01-02 15:04:05"
	if err := d.Set("create_time", detail.CreateTime.Local().Format(layout)); err != nil {
		return nil, fmt.Errorf("error setting create_time: %w", err)
	}
	if detail.PaymentTiming == flex.PaymentTimingPrepaid {
		if err := d.Set("expire_time", detail.ExpireTime.Local().Format(layout)); err != nil {
			return nil, fmt.Errorf("error setting expire_time: %w", err)
		}
	}

	return detail, nil
}

func updateName(d *schema.ResourceData, conn *connectivity.BaiduClient) error {
	if d.HasChange("name") {
		args := &mongodb.UpdateInstanceNameArgs{
			DbInstanceName: d.Get("name").(string),
		}
		log.Printf("[DEBUG] Update MongoDB Instance (%s) name: %+v", d.Id(), args)

		_, err := conn.WithMongoDBClient(func(client *mongodb.Client) (interface{}, error) {
			return nil, client.UpdateInstanceName(d.Id(), args)
		})
		return err
	}
	return nil
}

func updatePassword(d *schema.ResourceData, conn *connectivity.BaiduClient) error {
	if d.HasChange("account_password") {
		args := &mongodb.UpdatePasswordArgs{
			AccountPassword: d.Get("account_password").(string),
		}
		log.Printf("[DEBUG] Update MongoDB Instance (%s) password: %+v", d.Id(), args)

		_, err := conn.WithMongoDBClient(func(client *mongodb.Client) (interface{}, error) {
			return nil, client.UpdateAccountPassword(d.Id(), args)
		})
		return err
	}
	return nil
}
