package hpas

import (
	"fmt"
	"log"
	"strings"

	"github.com/baidubce/bce-sdk-go/services/hpas"
	"github.com/baidubce/bce-sdk-go/services/hpas/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/flex"
)

func ResourceReservedInstance() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to manage HPAS Reserved Instance. \n\n" +
			"~> **NOTE:** Destroying a reserved instance is not supported via API. " +
			"Running `terraform destroy` will only remove the resource from Terraform state; " +
			"the reserved instance itself will continue to run until it expires naturally.",

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Create: resourceReservedInstanceCreate,
		Read:   resourceReservedInstanceRead,
		Update: resourceReservedInstanceUpdate,
		Delete: resourceReservedInstanceDelete,

		Schema: map[string]*schema.Schema{
			"payment_timing": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "NoPrepay",
				Description:  "Payment timing of billing. Currently only `NoPrepay` is supported. Defaults to `NoPrepay`.",
				ValidateFunc: validation.StringInSlice([]string{"NoPrepay"}, false),
			},
			"period": {
				Type: schema.TypeInt,
				Description: "The reservation length (month). Effective when `payment_timing` is `NoPrepay`. " +
					"Valid values: `1`, `3`, `6`, `9`, `12`, `24`, `36`. Defaults to `1`.",
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntInSlice([]int{1, 3, 6, 9, 12, 24, 36}),
			},
			"auto_renew_period": {
				Type: schema.TypeInt,
				Description: "The automatic renewal time (month). Effective when `payment_timing` is `NoPrepay`. " +
					"Valid values: `1`, `3`, `6`, `9`, `12`, `24`, `36`.",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 3, 6, 9, 12, 24, 36}),
			},
			"auto_renew_period_unit": flex.SchemaAutoRenewTimeUnit(),
			"tags":                   flex.SchemaTagsOnlySupportCreation(),
			"name": {
				Type:     schema.TypeString,
				Required: true,
				Description: "Name of the reserved instance. Supports uppercase and lowercase letters, digits, " +
					"Chinese characters, and the special characters `-`, `_`, ` `, `/`, `.`. " +
					"Must start with a letter. Length: 1–65.",
			},
			"zone_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Zone information, e.g., `cn-bj-a`.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if len(old) == 0 || len(new) == 0 {
						return false
					}
					return strings.ToLower(old[len(old)-1:]) == strings.ToLower(new[len(new)-1:])
				},
			},
			"app_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Application type. e.g., `llama2_7B_train`.",
			},
			"app_performance_level": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Performance level of the application. e.g., `10k`.",
			},
			"ehc_cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "EHC cluster ID. If not specified, the system will automatically select the default EHC cluster.",
			},

			// computed fields
			"status": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "Status of the reserved instance. Possible values: `Creating`, `Active`, `Pending`, " +
					"`Expired`, `Recharge`, `Deleted`.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of the reserved instance in ISO8601 format.",
			},
			"expire_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Expiration time of the reserved instance in ISO8601 format.",
			},
			"hpas_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the HPAS instance that was deducted in the previous billing period.",
			},
			"hpas_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the HPAS instance that was deducted in the previous billing period.",
			},
			"deduct_instance": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether there is a deducting instance.",
			},
			"ehc_cluster_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the EHC cluster.",
			},
		},
	}
}

func resourceReservedInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)

	raw, err := conn.WithHPASClient(func(client *hpas.Client) (interface{}, error) {
		args := buildReservedInstanceCreationArgs(d)
		return client.CreateReservedHpas(args)
	})
	log.Printf("[DEBUG] Create HPAS Reserved Instance result: %+v", raw)
	if err != nil {
		return fmt.Errorf("error creating HPAS Reserved Instance: %w", err)
	}

	response := raw.(*api.CreateReservedHpasResp)
	if len(response.ReservedHpasIds) == 0 {
		return fmt.Errorf("error creating HPAS Reserved Instance: empty ID in response")
	}

	d.SetId(response.ReservedHpasIds[0])

	if _, err = waitReservedInstanceAvailable(conn, d.Id()); err != nil {
		return fmt.Errorf("error waiting for HPAS Reserved Instance (%s) to become available: %w", d.Id(), err)
	}

	return resourceReservedInstanceRead(d, meta)
}

func resourceReservedInstanceRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)

	detail, err := FindReservedInstance(conn, d.Id())
	log.Printf("[DEBUG] Read HPAS Reserved Instance (%s) result: %+v", d.Id(), detail)
	if err != nil {
		return fmt.Errorf("error reading HPAS Reserved Instance (%s): %w", d.Id(), err)
	}

	d.SetId(detail.ReservedHpasId)

	if err := d.Set("name", detail.Name); err != nil {
		return fmt.Errorf("error setting name: %w", err)
	}
	if err := d.Set("zone_name", detail.ZoneName); err != nil {
		return fmt.Errorf("error setting zone_name: %w", err)
	}
	if err := d.Set("app_type", detail.AppType); err != nil {
		return fmt.Errorf("error setting app_type: %w", err)
	}
	if err := d.Set("app_performance_level", detail.AppPerformanceLevel); err != nil {
		return fmt.Errorf("error setting app_performance_level: %w", err)
	}
	if err := d.Set("ehc_cluster_id", detail.EhcClusterId); err != nil {
		return fmt.Errorf("error setting ehc_cluster_id: %w", err)
	}
	if err := d.Set("tags", flex.FlattenTagModelToMap(detail.Tags)); err != nil {
		return fmt.Errorf("error setting tags: %w", err)
	}
	//if err := d.Set("period", detail.ReservedHpasPeriod); err != nil {
	//	return fmt.Errorf("error setting period: %w", err)
	//}
	if err := d.Set("status", detail.Status); err != nil {
		return fmt.Errorf("error setting status: %w", err)
	}
	if err := d.Set("create_time", detail.CreateTime); err != nil {
		return fmt.Errorf("error setting create_time: %w", err)
	}
	if err := d.Set("expire_time", detail.ExpireTime); err != nil {
		return fmt.Errorf("error setting expire_time: %w", err)
	}
	if err := d.Set("hpas_id", detail.HpasId); err != nil {
		return fmt.Errorf("error setting hpas_id: %w", err)
	}
	if err := d.Set("hpas_name", detail.HpasName); err != nil {
		return fmt.Errorf("error setting hpas_name: %w", err)
	}
	if err := d.Set("deduct_instance", detail.DeductInstance); err != nil {
		return fmt.Errorf("error setting deduct_instance: %w", err)
	}
	if err := d.Set("ehc_cluster_name", detail.EhcClusterName); err != nil {
		return fmt.Errorf("error setting ehc_cluster_name: %w", err)
	}

	return nil
}

func resourceReservedInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)

	if d.HasChange("name") {
		if err := updateReservedInstanceName(d, conn); err != nil {
			return fmt.Errorf("error updating HPAS Reserved Instance name: %w", err)
		}
	}

	if d.HasChanges("zone_name", "ehc_cluster_id") {
		if err := updateReservedInstanceZoneAndCluster(d, conn); err != nil {
			return fmt.Errorf("error updating HPAS Reserved Instance zone/cluster: %w", err)
		}
	}

	return resourceReservedInstanceRead(d, meta)
}

func resourceReservedInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	// Reserved instances cannot be deleted via API; they expire naturally.
	// Only the Terraform state record is removed.
	log.Printf("[WARN] HPAS Reserved Instance (%s) cannot be deleted via API. Removing from Terraform state only.", d.Id())
	return nil
}

func buildReservedInstanceCreationArgs(d *schema.ResourceData) *api.CreateReservedHpasReq {
	billingModel := api.BillingModel{
		ChargeType: d.Get("payment_timing").(string),
	}

	if v, ok := d.GetOk("period"); ok {
		billingModel.Period = int32(v.(int))
		billingModel.PeriodUnit = "month"
	}

	if _, ok := d.GetOk("auto_renew_period"); ok {
		billingModel.AutoRenew = true
		billingModel.AutoRenewPeriod = int32(d.Get("auto_renew_period").(int))
		billingModel.AutoRenewPeriodUnit = d.Get("auto_renew_period_unit").(string)
	}

	args := &api.CreateReservedHpasReq{
		Name:                d.Get("name").(string),
		ZoneName:            d.Get("zone_name").(string),
		AppType:             d.Get("app_type").(string),
		AppPerformanceLevel: d.Get("app_performance_level").(string),
		EhcClusterId:        d.Get("ehc_cluster_id").(string),
		BillingModel:        billingModel,
		PurchaseNum:         1,
		Tags:                flex.ExpandMapToTagModel[api.TagModel](d.Get("tags").(map[string]interface{})),
	}

	return args
}

func updateReservedInstanceName(d *schema.ResourceData, conn *connectivity.BaiduClient) error {
	_, err := conn.WithHPASClient(func(client *hpas.Client) (interface{}, error) {
		args := &api.ModifyReservedHpasNameReq{
			ReservedInstanceIds: []string{d.Id()},
			Name:                d.Get("name").(string),
		}
		return client.ModifyReservedHpasName(args)
	})
	return err
}

func updateReservedInstanceZoneAndCluster(d *schema.ResourceData, conn *connectivity.BaiduClient) error {
	_, err := conn.WithHPASClient(func(client *hpas.Client) (interface{}, error) {
		args := &api.ModifyReservedHpasReq{
			ModifyReservedHpasList: []api.ModifyReservedHpasModel{
				{
					ReservedInstanceId: d.Id(),
					ZoneName:           d.Get("zone_name").(string),
					EhcClusterId:       d.Get("ehc_cluster_id").(string),
				},
			},
		}
		return client.ModifyReservedHpas(args)
	})
	return err
}
