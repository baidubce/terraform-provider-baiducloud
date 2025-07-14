package hpas

import (
	"fmt"
	"log"
	"strings"

	"github.com/baidubce/bce-sdk-go/services/hpas"
	"github.com/baidubce/bce-sdk-go/services/hpas/api"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/flex"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func ResourceInstance() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to manage HPAS Instance. \n\n",

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Create: resourceInstanceCreate,
		Read:   resourceInstanceRead,
		Update: resourceInstanceUpdate,
		Delete: resourceInstanceDelete,

		Schema: map[string]*schema.Schema{
			"payment_timing":         flex.SchemaPaymentTiming(),
			"period":                 flex.SchemaReservationLength(),
			"auto_renew_period":      flex.SchemaAutoRenewLength(),
			"auto_renew_period_unit": flex.SchemaAutoRenewTimeUnit(),
			"tags":                   flex.UpdatableTagsSchema(),
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
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the instance.",
			},
			"application_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Name of the application.",
			},
			"auto_seq_suffix": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to automatically append a suffix to the application name. Defaults to `false`.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if old == "" {
						return true
					}
					return false
				},
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
			"image_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Image ID used for the application. Changing this value triggers a reinstallation of the OS.",
			},
			"internal_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Internal IP addresses. Must match the CIDR block of the specified subnet. Changing this value triggers a restart of the instance.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Subnet ID. Changing this value triggers a restart of the instance.",
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				Description: "Password of the instance. This value should be 8-16 characters, and letters, numbers and symbols must exist at the same time. " +
					"The symbols is limited to `!@#$%^*()`. Changing this value triggers a restart of the instance.",
			},
			"keypair_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the keypair to bind to the instance.",
			},
			"ehc_cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "EHC cluster ID. If not specified, the system will automatically select default EHC cluster.",
			},
			"security_group_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "normal",
				Description: "Security group type. Valid values: `normal`, `enterprise`. Defaults to `normal`.",
			},
			"security_group_ids": {
				Type:        schema.TypeSet,
				Required:    true,
				MinItems:    1,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of security group IDs, must be in the same VPC as the subnet",
			}, // computed
			"status": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "Status of the instance. Possible values: `Creating`, `Active`, `Expired`, `Error`, `Stopping`, `Starting`, " +
					"`Stopped`, `Reboot`, `Rebuild`, `Password`, `ChangeVpc`, `ChangeSubnet`, `Template`.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "VPC ID.",
			},
			"vpc_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the VPC.",
			},
			"vpc_cidr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CIDR block of the VPC.",
			},
			"subnet_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the subnet.",
			},
			"ehc_cluster_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the EHC cluster.",
			},
			"image_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the image.",
			},
			"keypair_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the keypair.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of the instance in ISO8601 format.",
			},
		},

		CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {
			_, passSet := diff.GetOk("password")
			_, keyPairSet := diff.GetOk("keypair_id")
			if !passSet && !keyPairSet {
				return fmt.Errorf("at least one of password or keypair_id must be set")
			}
			return nil
		},
	}
}

func resourceInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)

	raw, err := conn.WithHPASClient(func(client *hpas.Client) (interface{}, error) {
		args := buildCreationArgs(d, client)
		return client.CreateHpas(args)
	})
	log.Printf("[DEBUG] Create HPAS Instance result: %+v", raw)
	if err != nil {
		return fmt.Errorf("error creating HPAS Instance: %w", err)
	}
	response := raw.(*api.CreateHpasResp)
	if response.HpasIds == nil || len(response.HpasIds) == 0 {
		return fmt.Errorf("error creating HPAS Instance: %+v", raw)
	}

	d.SetId(response.HpasIds[0])

	if _, err = waitInstanceAvailable(conn, d.Id()); err != nil {
		return fmt.Errorf("error waiting HPAS Instance (%s) becoming available: %w", d.Id(), err)
	}
	return resourceInstanceRead(d, meta)
}

func resourceInstanceRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)

	detail, err := FindInstance(conn, d.Id())
	log.Printf("[DEBUG] Read HPAS Instance (%s) result: %+v", d.Id(), detail)
	if err != nil {
		return fmt.Errorf("error reading HPAS Instance (%s): %w", d.Id(), err)
	}

	d.SetId(detail.HpasId)

	if err := d.Set("payment_timing", detail.ChargeType); err != nil {
		return fmt.Errorf("error setting payment_timing: %w", err)
	}
	if err := d.Set("tags", flex.FlattenTagModelToMap(detail.Tags)); err != nil {
		return fmt.Errorf("error setting tags: %w", err)
	}
	if err := d.Set("app_type", detail.AppType); err != nil {
		return fmt.Errorf("error setting app_type: %w", err)
	}
	if err := d.Set("app_performance_level", detail.AppPerformanceLevel); err != nil {
		return fmt.Errorf("error setting app_performance_level: %w", err)
	}
	if err := d.Set("name", detail.Name); err != nil {
		return fmt.Errorf("error setting name: %w", err)
	}
	if err := d.Set("zone_name", detail.ZoneName); err != nil {
		return fmt.Errorf("error setting zone_name: %w", err)
	}
	if err := d.Set("image_id", detail.ImageId); err != nil {
		return fmt.Errorf("error setting image_id: %w", err)
	}
	if err := d.Set("internal_ip", detail.InternalIp); err != nil {
		return fmt.Errorf("error setting internal_ip: %w", err)
	}
	if err := d.Set("subnet_id", detail.SubnetId); err != nil {
		return fmt.Errorf("error setting subnet_id: %w", err)
	}
	if err := d.Set("ehc_cluster_id", detail.EhcClusterId); err != nil {
		return fmt.Errorf("error setting ehc_cluster_id: %w", err)
	}
	if err := d.Set("keypair_id", detail.KeypairId); err != nil {
		return fmt.Errorf("error setting keypair_id: %w", err)
	}
	if len(detail.NicInfo) > 0 {
		if err := d.Set("security_group_type", detail.NicInfo[0].SecurityGroupType); err != nil {
			return fmt.Errorf("error setting security_group_type: %w", err)
		}
		if err := d.Set("security_group_ids", flex.FlattenStringValueSet(detail.NicInfo[0].SecurityGroupIds)); err != nil {
			return fmt.Errorf("error setting security_group_ids: %w", err)
		}
	}

	// computed fields
	if err := d.Set("status", detail.Status); err != nil {
		return fmt.Errorf("error setting status: %w", err)
	}
	if err := d.Set("vpc_id", detail.VpcId); err != nil {
		return fmt.Errorf("error setting vpc_id: %w", err)
	}
	if err := d.Set("vpc_name", detail.VpcName); err != nil {
		return fmt.Errorf("error setting vpc_name: %w", err)
	}
	if err := d.Set("vpc_cidr", detail.VpcCidr); err != nil {
		return fmt.Errorf("error setting vpc_cidr: %w", err)
	}
	if err := d.Set("subnet_name", detail.SubnetName); err != nil {
		return fmt.Errorf("error setting subnet_name: %w", err)
	}
	if err := d.Set("ehc_cluster_name", detail.EhcClusterName); err != nil {
		return fmt.Errorf("error setting ehc_cluster_name: %w", err)
	}
	if err := d.Set("image_name", detail.ImageName); err != nil {
		return fmt.Errorf("error setting image_name: %w", err)
	}
	if err := d.Set("keypair_name", detail.KeypairName); err != nil {
		return fmt.Errorf("error setting keypair_name: %w", err)
	}
	if err := d.Set("create_time", detail.CreateTime); err != nil {
		return fmt.Errorf("error setting create_time: %w", err)
	}

	return nil
}

func resourceInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)

	if err := updateSecurityGroups(d, conn); err != nil {
		return fmt.Errorf("error updating HPAS Instance security groups: %w", err)
	}
	if err := updateTags(d, conn); err != nil {
		return fmt.Errorf("error updating HPAS Instance tags: %w", err)
	}
	if err := updateAttributes(d, conn); err != nil {
		return fmt.Errorf("error updating HPAS Instance attributes: %w", err)
	}

	// need reboot or rebuild
	if err := updateSubnetAndInternalIP(d, conn); err != nil {
		return err
	}

	if err := updateImageAndPassword(d, conn); err != nil {
		return err
	}

	return resourceInstanceRead(d, meta)
}

func resourceInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)

	_, err := conn.WithHPASClient(func(client *hpas.Client) (interface{}, error) {
		args := &api.DeleteHpasReq{
			HpasIds: []string{d.Id()},
		}
		return nil, client.DeleteHpas(args)
	})
	log.Printf("[DEBUG] Delete HPAS Instance (%s)", d.Id())

	if err != nil {
		return fmt.Errorf("error delete HPAS Instance (%s): %w", d.Id(), err)
	}

	return nil
}

func buildCreationArgs(d *schema.ResourceData, client *hpas.Client) *api.CreateHpasReq {
	billingModel := api.BillingModel{
		ChargeType: d.Get("payment_timing").(string),
		Period:     int32(d.Get("period").(int)),
		PeriodUnit: "month",
	}

	if _, ok := d.GetOk("auto_renew_period"); ok {
		billingModel.AutoRenew = true
		billingModel.AutoRenewPeriod = int32(d.Get("auto_renew_period").(int))
		billingModel.AutoRenewPeriodUnit = d.Get("auto_renew_period_unit").(string)
	}

	args := &api.CreateHpasReq{
		AppType:             d.Get("app_type").(string),
		AppPerformanceLevel: d.Get("app_performance_level").(string),
		Name:                d.Get("name").(string),
		ApplicationName:     d.Get("application_name").(string),
		AutoSeqSuffix:       d.Get("auto_seq_suffix").(bool),
		PurchaseNum:         1,
		ZoneName:            d.Get("zone_name").(string),
		ImageId:             d.Get("image_id").(string),
		SubnetId:            d.Get("subnet_id").(string),
		KeypairId:           d.Get("keypair_id").(string),
		EhcClusterId:        d.Get("ehc_cluster_id").(string),
		SecurityGroupType:   d.Get("security_group_type").(string),
		SecurityGroupIds:    flex.ExpandStringValueSet(d.Get("security_group_ids").(*schema.Set)),
		BillingModel:        billingModel,
		Tags:                flex.ExpandMapToTagModel[api.TagModel](d.Get("tags").(map[string]interface{})),
	}

	if _, ok := d.GetOk("password"); ok {
		args.Password = encryptPassword(d, client)
	}

	if v, ok := d.GetOk("internal_ip"); ok {
		args.InternalIps = []string{v.(string)}
	}

	return args
}

func updateSecurityGroups(d *schema.ResourceData, conn *connectivity.BaiduClient) error {
	if d.HasChanges("security_group_type", "security_group_ids") {
		_, err := conn.WithHPASClient(func(client *hpas.Client) (interface{}, error) {
			args := &api.SecurityGroupsReq{
				HpasIds:           []string{d.Id()},
				SecurityGroupIds:  flex.ExpandStringValueSet(d.Get("security_group_ids").(*schema.Set)),
				SecurityGroupType: d.Get("security_group_type").(string),
			}
			return client.ReplaceSecurityGroups(args)
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func updateTags(d *schema.ResourceData, conn *connectivity.BaiduClient) error {
	if d.HasChange("tags") {
		o, n := d.GetChange("tags")
		added, removed := flex.DiffMaps(o.(map[string]interface{}), n.(map[string]interface{}))

		if len(removed) > 0 {
			_, err := conn.WithHPASClient(func(client *hpas.Client) (interface{}, error) {
				args := &api.TagsOperationRequest{
					ResourceType: "hpas",
					ResourceIds:  []string{d.Id()},
					Tags:         flex.ExpandMapToTagModel[api.TagModel](removed),
				}
				return nil, client.DetachTags(args)
			})

			if err != nil {
				return err
			}
		}

		if len(added) > 0 {
			_, err := conn.WithHPASClient(func(client *hpas.Client) (interface{}, error) {
				args := &api.TagsOperationRequest{
					ResourceType: "hpas",
					ResourceIds:  []string{d.Id()},
					Tags:         flex.ExpandMapToTagModel[api.TagModel](added),
				}
				return nil, client.AttachTags(args)
			})
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func updateAttributes(d *schema.ResourceData, conn *connectivity.BaiduClient) error {
	if d.HasChanges("name", "application_name") {
		_, err := conn.WithHPASClient(func(client *hpas.Client) (interface{}, error) {
			args := &api.ModifyInstancesAttributeReq{
				HpasIds:         []string{d.Id()},
				Name:            d.Get("name").(string),
				ApplicationName: d.Get("application_name").(string),
			}
			return nil, client.ModifyInstancesAttribute(args)
		})

		if err != nil {
			return err
		}
	}
	return nil
}

func updateSubnetAndInternalIP(d *schema.ResourceData, conn *connectivity.BaiduClient) error {
	if d.HasChanges("subnet_id", "internal_ip") {
		_, err := conn.WithHPASClient(func(client *hpas.Client) (interface{}, error) {
			args := &api.ModifyInstancesSubnetRequest{
				HpasIds:   []string{d.Id()},
				SubnetId:  d.Get("subnet_id").(string),
				PrivateIp: d.Get("internal_ip").(string),
			}
			return client.ModifyInstancesSubnet(args)
		})

		if err != nil {
			return fmt.Errorf("error updating HPAS Instance subnet and internal IP: %w", err)
		}

		_, err = waitInstanceAvailable(conn, d.Id())
		if err != nil {
			return fmt.Errorf("error waiting HPAS Instance (%s) becoming available: %w", d.Id(), err)
		}
	}
	return nil
}

func updateImageAndPassword(d *schema.ResourceData, conn *connectivity.BaiduClient) error {
	if d.HasChanges("image_id", "password") {
		if d.HasChange("image_id") {
			err := resetInstance(d, conn)
			if err != nil {
				return fmt.Errorf("error resetting HPAS instance (%s): %w", d.Id(), err)
			}
		} else if d.HasChange("password") {
			err := modifyPassword(d, conn)
			if err != nil {
				return fmt.Errorf("error modifying password for HPAS instance (%s): %w", d.Id(), err)
			}
		}

		_, err := waitInstanceAvailable(conn, d.Id())
		if err != nil {
			return fmt.Errorf("error waiting HPAS Instance (%s) becoming available: %w", d.Id(), err)
		}
	}

	return nil
}

func resetInstance(d *schema.ResourceData, conn *connectivity.BaiduClient) error {
	_, err := conn.WithHPASClient(func(client *hpas.Client) (interface{}, error) {
		args := api.ResetHpasReq{
			HpasIds:   []string{d.Id()},
			ImageId:   d.Get("image_id").(string),
			KeypairId: d.Get("keypair_id").(string),
		}

		if _, ok := d.GetOk("password"); ok {
			args.Password = encryptPassword(d, client)
		}
		return nil, client.ResetHpas(&args)
	})
	if err != nil {
		return err
	}
	return nil
}

func modifyPassword(d *schema.ResourceData, conn *connectivity.BaiduClient) error {
	_, err := conn.WithHPASClient(func(client *hpas.Client) (interface{}, error) {
		args := api.ModifyPasswordHpasReq{
			HpasId:   d.Id(),
			Password: encryptPassword(d, client),
		}
		return nil, client.ModifyPasswordHpas(&args)
	})
	if err != nil {
		return err
	}
	return nil
}
