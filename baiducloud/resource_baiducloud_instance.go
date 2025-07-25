/*
Use this resource to create a BCC instance.

~> **NOTE:** The terminate operation of bcc does NOT take effect immediately，maybe takes for several minites.
Example Usage

```hcl

	resource "baiducloud_instance" "my-server" {
	  image_id = "m-A4jJpFzi"
	  name = "my-instance"
	  availability_zone = "cn-bj-a"
	  cpu_count = "2"
	  memory_capacity_in_gb = "8"
	  billing = {
	    payment_timing = "Postpaid"
	  }
	}

```

# Import

BCC instance can be imported, e.g.

```hcl
$ terraform import baiducloud_instance.my-server id
```
*/
package baiducloud

import (
	"fmt"
	"strconv"
	"time"

	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/bcc"
	"github.com/baidubce/bce-sdk-go/services/bcc/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/rateLimit"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/flex"
)

func resourceBaiduCloudInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudInstanceCreate,
		Read:   resourceBaiduCloudInstanceRead,
		Update: resourceBaiduCloudInstanceUpdate,
		Delete: resourceBaiduCloudInstanceDelete,

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
				Computed:    true,
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Description: "Availability zone to start the instance in.",
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
			},
			"payment_timing": {
				Type: schema.TypeString,
				Description: "Payment timing of billing, which can be Prepaid or Postpaid. The default is Postpaid. " +
					"When switching to Prepaid, reservation length must be set. " +
					"Switching to Postpaid takes effect immediately.",
				Optional:     true,
				Default:      api.PaymentTimingPostPaid,
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
							Type: schema.TypeInt,
							Description: "The reservation length that you will pay for your resource. " +
								"It is valid when payment_timing is Prepaid. Valid values: [1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36].",
							Required:         true,
							Default:          1,
							ValidateFunc:     validateReservationLength(),
							DiffSuppressFunc: postPaidDiffSuppressFunc,
						},
						"reservation_time_unit": {
							Type: schema.TypeString,
							Description: "The reservation time unit that you will pay for your resource." +
								" It is valid when payment_timing is Prepaid. The value can only be month currently, which is also the default value.",
							Required:         true,
							Default:          "Month",
							ValidateFunc:     validateReservationUnit(),
							DiffSuppressFunc: postPaidDiffSuppressFunc,
						},
					},
				},
			},
			"instance_type": {
				Type:             schema.TypeString,
				Description:      "Type of the instance to start. Available values are N1, N2, N3, N4, N5, C1, C2, S1, G1, F1. Default to N3.",
				Optional:         true,
				ForceNew:         true,
				Default:          api.InstanceTypeN3,
				ValidateFunc:     validateInstanceType(),
				DiffSuppressFunc: specDiffSuppressFunc,
			},
			"admin_pass": {
				Type:        schema.TypeString,
				Description: "Password of the instance to be started. This value should be 8-16 characters, and English, numbers and symbols must exist at the same time. The symbols is limited to \"!@#$%^*()\".",
				Optional:    true,
				Sensitive:   true,
			},
			"cpu_count": {
				Type:             schema.TypeInt,
				Description:      "Number of CPU cores to be created for the instance.",
				Optional:         true,
				ValidateFunc:     validation.IntAtLeast(1),
				DiffSuppressFunc: specDiffSuppressFunc,
			},
			"memory_capacity_in_gb": {
				Type:             schema.TypeInt,
				Description:      "Memory capacity(GB) of the instance to be created.",
				Optional:         true,
				ValidateFunc:     validation.IntAtLeast(1),
				DiffSuppressFunc: specDiffSuppressFunc,
			},
			"root_disk_size_in_gb": {
				Type: schema.TypeInt,
				Description: "System disk size(GB) of the instance to be created. The value range is [40,2048]GB," +
					" Default to 40GB, and more than 40GB is charged according to the cloud disk price. " +
					"Note that the specified system disk size needs to meet the minimum disk space limit of the mirror used.",
				Optional:     true,
				ForceNew:     true,
				Default:      40,
				ValidateFunc: validation.IntBetween(40, 2048),
			},
			"root_disk_storage_type": {
				Type: schema.TypeString,
				Description: "System disk storage type of the instance. " +
					"Available values are " +
					"enhanced_ssd_pl1, enhanced_ssd_pl2, cloud_hp1, premium_ssd, hp1, " +
					"ssd, sata, hdd, local, sata, local-ssd, local-hdd, local-nvme. " +
					"Default to cloud_hp1.",
				Optional: true,
				ForceNew: true,
				Default:  api.StorageTypeCloudHP1,
				//ValidateFunc: validateStorageType(),
			},
			"ephemeral_disks": {
				Type:        schema.TypeList,
				Description: "Ephemeral disks of the instance.",
				Optional:    true,
				MinItems:    1,
				MaxItems:    15,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size_in_gb": {
							Type:         schema.TypeInt,
							Description:  "Size(GB) of the ephemeral disk.",
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.IntAtLeast(0),
						},
						"storage_type": {
							Type:         schema.TypeString,
							Description:  "Storage type of the ephemeral disk. Available values are std1, hp1, cloud_hp1, local, sata, ssd. Default to cloud_hp1.",
							Optional:     true,
							ForceNew:     true,
							Default:      api.StorageTypeCloudHP1,
							ValidateFunc: validateStorageType(),
						},
					},
				},
			},
			"cds_disks": {
				Type:        schema.TypeList,
				Description: "CDS disks of the instance.",
				Optional:    true,
				MinItems:    1,
				MaxItems:    5,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cds_size_in_gb": {
							Type:         schema.TypeInt,
							Description:  "The size(GB) of CDS.",
							Optional:     true,
							ForceNew:     true,
							Default:      0,
							ValidateFunc: validation.IntAtLeast(0),
						},
						"storage_type": {
							Type:         schema.TypeString,
							Description:  "Storage type of the CDS.",
							Optional:     true,
							ForceNew:     true,
							Default:      api.StorageTypeCloudHP1,
							ValidateFunc: validateStorageType(),
						},
						"snapshot_id": {
							Type:        schema.TypeString,
							Description: "Snapshot ID of CDS.",
							Optional:    true,
							ForceNew:    true,
						},
					},
				},
			},
			"public_ip": {
				Type:        schema.TypeString,
				Description: "Public IP",
				Computed:    true,
			},
			"dedicate_host_id": {
				Type:        schema.TypeString,
				Description: "The ID of dedicated host.",
				Optional:    true,
				ForceNew:    true,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Description: "The subnet ID of VPC. The default subnet will be used when it is empty. The instance will restart after changing the subnet.",
				Optional:    true,
				Computed:    true,
			},
			"security_groups": {
				Type:          schema.TypeSet,
				Description:   "Security group ids of the instance.",
				ConflictsWith: []string{"enterprise_security_groups"},
				Optional:      true,
				Computed:      true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"enterprise_security_groups": {
				Type:          schema.TypeSet,
				Description:   "Enterprise security group ids of the instance. ",
				ConflictsWith: []string{"security_groups"},
				Optional:      true,
				Computed:      true,
				MinItems:      1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"gpu_card": {
				Type:          schema.TypeString,
				Description:   "GPU card of the instance.",
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"fpga_card"},
			},
			"fpga_card": {
				Type:          schema.TypeString,
				Description:   "FPGA card of the instance.",
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"gpu_card"},
			},
			"card_count": {
				Type:        schema.TypeString,
				Description: "Count of the GPU cards or FPGA cards to be carried for the instance to be created, it is valid only when the gpu_card or fpga_card field is not empty.",
				Optional:    true,
				Default:     0,
			},
			"auto_renew_time_unit": {
				Type:             schema.TypeString,
				Description:      "Time unit of automatic renewal, the value can be month or year. The default value is empty, indicating no automatic renewal. It is valid only when the payment_timing is Prepaid.",
				Optional:         true,
				ValidateFunc:     validation.StringInSlice([]string{"month", "year"}, false),
				DiffSuppressFunc: postPaidDiffSuppressFunc,
			},
			"auto_renew_time_length": {
				Type: schema.TypeInt,
				Description: "The time length of automatic renewal. Effective only when `payment_timing` is `Prepaid`. " +
					"Valid values are `1–9` when `auto_renew_time_unit` is `month` and `1–3` when it is `year`. " +
					"Defaults to `1`. Due to API limitations, modifying this parameter after the auto-renewal rule is created " +
					"will first delete the existing rule and then recreate it.",
				Optional:         true,
				Default:          1,
				ValidateFunc:     validation.IntBetween(1, 9),
				DiffSuppressFunc: autoRenewDiffSuppressFunc,
			},
			"cds_auto_renew": {
				Type: schema.TypeBool,
				Description: "[This parameter is deprecated as CDS auto-renewal now aligns with the BCC instance.] " +
					"Whether the cds is automatically renewed. It is valid when payment_timing is Prepaid. Default to false.",
				Deprecated:       "This parameter is deprecated as CDS auto-renewal now aligns with the BCC instance.",
				Optional:         true,
				Default:          false,
				DiffSuppressFunc: postPaidDiffSuppressFunc,
			},
			"sync_eip_auto_renew_rule": {
				Type: schema.TypeBool,
				Description: "Whether to synchronize the EIP's auto-renewal rule with that of the associated BCC instance. " +
					"This setting applies during both the creation and deletion of the BCC's auto-renewal rule. " +
					"Modifying this parameter alone does not trigger any change to the EIP's auto-renewal rule. " +
					"Effective only when `payment_timing` is `Prepaid`. Defaults to `true`.",
				Optional:         true,
				Default:          true,
				DiffSuppressFunc: postPaidDiffSuppressFunc,
			},
			"related_release_flag": {
				Type:        schema.TypeBool,
				Description: "Whether to release the eip and data disks mounted by the current instance. Can only be released uniformly or not. Default to false.",
				Optional:    true,
				Default:     false,
				ForceNew:    true,
			},
			"delete_cds_snapshot_flag": {
				Type:        schema.TypeBool,
				Description: "Whether to release the cds disk snapshots, default to false. It is effective only when the related_release_flag is true.",
				Optional:    true,
				Default:     false,
				ForceNew:    true,
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
			"placement_policy": {
				Type:        schema.TypeString,
				Description: "The placement policy of the instance, which can be default or dedicatedHost.",
				Computed:    true,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Description: "VPC ID of the instance.",
				Computed:    true,
			},
			"network_capacity_in_mbps": {
				Type:        schema.TypeInt,
				Description: "Public network bandwidth(Mbps) of the instance.",
				Computed:    true,
			},
			"auto_renew": {
				Type:        schema.TypeBool,
				Description: "Whether to automatically renew.",
				Computed:    true,
			},
			"keypair_id": {
				Type:        schema.TypeString,
				Description: "Key pair id of the instance.",
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
			},
			"keypair_name": {
				Type:        schema.TypeString,
				Description: "Key pair name of the instance.",
				Computed:    true,
			},
			"relation_tag": {
				Type:        schema.TypeBool,
				Description: "The new instance associated with existing Tags or not, default false. The Tags should already exit if set true",
				Optional:    true,
				ForceNew:    true,
			},
			"action": {
				Type:         schema.TypeString,
				Description:  "Start or stop the instance, which can only be start or stop, default start.",
				Optional:     true,
				Default:      INSTANCE_ACTION_START,
				ValidateFunc: validation.StringInSlice([]string{INSTANCE_ACTION_START, INSTANCE_ACTION_STOP}, false),
			},
			"user_data": {
				Type:        schema.TypeString,
				Description: "User Data",
				Optional:    true,
			},
			"instance_spec": {
				Type:        schema.TypeString,
				Description: "spec",
				Optional:    true,
				Computed:    true,
			},
			"deploy_set_ids": {
				Type:        schema.TypeSet,
				Description: "Deploy set ids the instance belong to",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"hostname": {
				Type:        schema.TypeString,
				Description: "Hostname of the instance.",
				Optional:    true,
				Computed:    true,
			},
			"is_open_hostname_domain": {
				Type:        schema.TypeBool,
				Description: "Whether to automatically generate hostname domain.",
				Optional:    true,
			},
			"is_open_ipv6": {
				Type: schema.TypeBool,
				Description: "Whether to enable IPv6 for the instance to be created. " +
					"It can be enabled only when both the image and the subnet support IPv6. " +
					"True means enabled, false means disabled, " +
					"undefined means automatically adapting to the IPv6 support of the image and subnet.",
				Optional: true,
			},
			"tags": tagsSchema(),
			"resource_group_id": {
				Type:        schema.TypeString,
				Description: "Resource group Id of the instance.",
				Optional:    true,
			},
			"stop_with_no_charge": {
				Type:        schema.TypeBool,
				Description: "Whether to enable stopping charging after shutdown for postpaid instance without local disks. Defaults to false.",
				Optional:    true,
				Default:     false,
			},
		},
	}
}

func resourceBaiduCloudInstanceCreate(d *schema.ResourceData, meta interface{}) error {

	action := "Create BCC Instance"

	client := meta.(*connectivity.BaiduClient)
	bccService := BccService{client}

	var createArgs interface{}
	createBySpec := false
	if value, ok := d.GetOk("instance_spec"); ok && value.(string) != "" {
		createBySpec = true
	}

	securityGroups := expandStringSet(d.Get("security_groups").(*schema.Set))
	enterpriseSecurityGroups := expandStringSet(d.Get("enterprise_security_groups").(*schema.Set))

	if createBySpec {
		createInstanceArgs, err := buildBaiduCloudInstanceBySpecArgs(d, meta)
		if err != nil {
			return WrapError(err)
		}

		if len(securityGroups) > 0 {
			createInstanceArgs.SecurityGroupIds = securityGroups
		}
		if len(enterpriseSecurityGroups) > 0 {
			createInstanceArgs.EnterpriseSecurityGroupIds = enterpriseSecurityGroups
		}

		createArgs = createInstanceArgs
	} else {
		createInstanceArgs, err := buildBaiduCloudInstanceArgs(d, meta)
		if err != nil {
			return WrapError(err)
		}

		if len(securityGroups) > 0 {
			createInstanceArgs.SecurityGroupIds = securityGroups
		}
		if len(enterpriseSecurityGroups) > 0 {
			createInstanceArgs.EnterpriseSecurityGroupIds = enterpriseSecurityGroups
		}

		createArgs = createInstanceArgs
	}

	err := ratelimit.Check(action)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_instance", action, BCESDKGoERROR)
	}

	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithBccClient(func(bccClient *bcc.Client) (interface{}, error) {
			if createBySpec {
				return bccClient.CreateInstanceBySpec(createArgs.(*api.CreateInstanceBySpecArgs))
			}
			return bccClient.CreateInstance(createArgs.(*api.CreateInstanceArgs))
		})
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, raw)
		if createBySpec {
			response, _ := raw.(*api.CreateInstanceBySpecResult)
			d.SetId(response.InstanceIds[0])
		} else {
			response, _ := raw.(*api.CreateInstanceResult)
			d.SetId(response.InstanceIds[0])
		}
		return nil
	})
	ratelimit.CheckEnd()

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_instance", action, BCESDKGoERROR)
	}

	stateConf := buildStateConf(
		[]string{string(api.InstanceStatusStarting)},
		[]string{string(api.InstanceStatusRunning), InstanceStatusDeleted},
		d.Timeout(schema.TimeoutCreate),
		bccService.InstanceStateRefresh(d.Id()),
	)
	if _, err = stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_instance", action, BCESDKGoERROR)
	}

	// set instance description
	if err := updateInstanceDescription(d, meta, d.Id()); err != nil {
		return err
	}

	// check tag bind
	if err := checkTagBind(d, meta); err != nil {
		return err
	}

	// stop the instance if the action field is stop.
	if d.Get("action").(string) == INSTANCE_ACTION_STOP {
		stopWithNoCharge := d.Get("stop_with_no_charge").(bool)
		if err := bccService.StopInstance(d.Id(), stopWithNoCharge, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return err
		}
	}

	return resourceBaiduCloudInstanceRead(d, meta)
}

func checkTagBind(d *schema.ResourceData, meta interface{}) error {
	if v, ok := d.GetOk("tags"); ok {
		client := meta.(*connectivity.BaiduClient)
		instanceID := d.Id()
		action := "Retry BCC Instance tags bind " + instanceID
		raw, err := client.WithBccClient(func(bccClient *bcc.Client) (interface{}, error) {
			return bccClient.GetInstanceDetail(instanceID)
		})
		addDebug(action, raw)

		if err != nil {
			if NotFoundError(err) {
				d.SetId("")
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_instance", action, BCESDKGoERROR)
		}
		response, _ := raw.(*api.GetInstanceDetailResult)
		if response.Instance.Tags == nil || len(response.Instance.Tags) == 0 {
			// bind tags failed, retry
			_, err := client.WithBccClient(func(bccClient *bcc.Client) (interface{}, error) {
				tagArgs := &api.BindTagsRequest{
					ChangeTags: tranceTagMapToModel(v.(map[string]interface{})),
				}
				return nil, bccClient.BindInstanceToTags(instanceID, tagArgs)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, "baiducloud_instance", action, BCESDKGoERROR)
			}
		}
	}
	return nil
}

func resourceBaiduCloudInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	bccService := BccService{client}

	instanceID := d.Id()
	action := "Query BCC Instance " + instanceID
	raw, err := client.WithBccClient(func(bccClient *bcc.Client) (interface{}, error) {
		return bccClient.GetInstanceDetail(instanceID)
	})
	addDebug(action, raw)

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_instance", action, BCESDKGoERROR)
	}
	response, _ := raw.(*api.GetInstanceDetailResult)

	// Required or Optional
	d.Set("image_id", response.Instance.ImageId)
	d.Set("name", response.Instance.InstanceName)
	d.Set("availability_zone", response.Instance.ZoneName)
	d.Set("instance_type", string(response.Instance.InstanceType))
	d.Set("cpu_count", response.Instance.CpuCount)
	d.Set("memory_capacity_in_gb", response.Instance.MemoryCapacityInGB)
	d.Set("subnet_id", response.Instance.SubnetId)
	d.Set("gpu_card", response.Instance.GpuCard)
	d.Set("fpga_card", response.Instance.FpgaCard)
	d.Set("card_count", response.Instance.CardCount)
	d.Set("dedicate_host_id", response.Instance.DedicatedHostId)
	d.Set("tags", flattenTagsToMap(response.Instance.Tags))
	d.Set("instance_spec", response.Instance.Spec)

	d.Set("payment_timing", response.Instance.PaymentTiming)
	d.Set("auto_renew_time_unit", response.Instance.AutoRenewPeriodUnit)
	d.Set("auto_renew_time_length", response.Instance.AutoRenewPeriod)

	// Computed
	d.Set("description", response.Instance.Description)
	d.Set("status", response.Instance.Status)
	d.Set("create_time", response.Instance.CreationTime)
	d.Set("expire_time", response.Instance.ExpireTime)
	d.Set("public_ip", response.Instance.PublicIP)
	d.Set("internal_ip", response.Instance.InternalIP)
	d.Set("placement_policy", response.Instance.PlacementPolicy)
	d.Set("vpc_id", response.Instance.VpcId)
	d.Set("network_capacity_in_mbps", response.Instance.NetworkCapacityInMbps)
	d.Set("keypair_id", response.Instance.KeypairId)
	d.Set("keypair_name", response.Instance.KeypairName)
	d.Set("auto_renew", response.Instance.AutoRenew)
	d.Set("hostname", response.Instance.Hostname)

	d.Set("security_groups", flex.FlattenStringValueSet(response.Instance.NicInfo.SecurityGroups))
	d.Set("enterprise_security_groups", flex.FlattenStringValueSet(response.Instance.NicInfo.EnterpriseSecurityGroups))

	deploysetIds := make([]string, 0)
	for _, value := range response.Instance.DeploySetList {
		deploysetIds = append(deploysetIds, value.DeploySetId)
	}
	d.Set("deploy_set_ids", deploysetIds)

	// read ephemeral disks
	ephVolumes, err := bccService.ListAllEphemeralVolumes(instanceID)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_instance", action, BCESDKGoERROR)
	}
	ephDisks := make([]interface{}, 0, len(ephVolumes))
	for _, eph := range ephVolumes {
		ephMap := make(map[string]interface{})
		ephMap["size_in_gb"] = eph.DiskSizeInGB
		ephMap["storage_type"] = eph.StorageType

		ephDisks = append(ephDisks, ephMap)
	}
	d.Set("ephemeral_disks", ephDisks)

	// read system disks
	sysVolume, err := bccService.GetSystemVolume(instanceID)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_instance", action, BCESDKGoERROR)
	}
	d.Set("root_disk_size_in_gb", sysVolume.DiskSizeInGB)
	d.Set("root_disk_storage_type", sysVolume.StorageType)

	return nil
}

func resourceBaiduCloudInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	instanceID := d.Id()

	d.Partial(true)

	// update instance attribute
	if err := updateInstanceAttribute(d, meta, instanceID); err != nil {
		return err
	}

	// update instance description
	if err := updateInstanceDescription(d, meta, instanceID); err != nil {
		return err
	}

	// update instance image id (rebuild)
	if err := updateInstanceImage(d, meta, instanceID); err != nil {
		return err
	}

	// update instance admin pass
	if err := updateInstanceAdminPass(d, meta, instanceID); err != nil {
		return err
	}

	// update instance capacity, include cpu count, memory capacity and ephemeral disks
	if err := updateInstanceCapacity(d, meta, instanceID); err != nil {
		return err
	}

	// update instance spec
	if err := updateInstanceSpec(d, meta, instanceID); err != nil {
		return err
	}

	// update instance security groups
	if err := updateInstanceSecurityGroups(d, meta, instanceID); err != nil {
		return err
	}

	// update instance enterprise security groups
	if err := updateInstanceEnterpriseSecurityGroups(d, meta, instanceID); err != nil {
		return err
	}

	// update instance subnet
	if err := updateInstanceSubnet(d, meta, instanceID); err != nil {
		return err
	}

	// update instance action
	if err := updateInstanceAction(d, meta, instanceID); err != nil {
		return err
	}

	// update instance deploy
	if err := updateInstanceDeploy(d, meta, instanceID); err != nil {
		return err
	}

	// update instance hostname
	if err := updateInstanceHostname(d, meta, instanceID); err != nil {
		return err
	}

	if d.HasChange("payment_timing") {
		// update payment timing
		if err := updateInstancePaymentTiming(d, meta, instanceID); err != nil {
			return err
		}
	} else if d.HasChanges("auto_renew_time_unit", "auto_renew_time_length") {
		// update auto renew rules
		if err := updateInstanceAutoRenew(d, meta, instanceID); err != nil {
			return err
		}
	}

	d.Partial(false)

	return resourceBaiduCloudInstanceRead(d, meta)
}

func resourceBaiduCloudInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	bccService := BccService{client}

	instanceId := d.Id()
	action := "Delete BCC Instance " + instanceId

	// delete instance
	paymentTiming := d.Get("payment_timing").(string)
	var err error
	if paymentTiming == "Postpaid" {
		args := &api.DeleteInstanceWithRelateResourceArgs{}
		if v, ok := d.GetOk("related_release_flag"); ok {
			args.RelatedReleaseFlag = v.(bool)
		}
		if v, ok := d.GetOk("delete_cds_snapshot_flag"); ok {
			args.DeleteCdsSnapshotFlag = v.(bool)
		}
		err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
			raw, err := client.WithBccClient(func(bccClient *bcc.Client) (interface{}, error) {
				return instanceId, bccClient.DeleteInstanceWithRelateResource(instanceId, args)
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
	} else {
		args := &api.DeletePrepaidInstanceWithRelateResourceArgs{
			InstanceId: instanceId,
		}
		if v, ok := d.GetOk("related_release_flag"); ok {
			args.RelatedReleaseFlag = v.(bool)
		}
		if v, ok := d.GetOk("delete_cds_snapshot_flag"); ok {
			args.DeleteCdsSnapshotFlag = v.(bool)
		}
		err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
			raw, err := client.WithBccClient(func(bccClient *bcc.Client) (interface{}, error) {
				return bccClient.DeletePrepaidInstanceWithRelateResource(args)
			})
			if err != nil {
				if IsExceptedErrors(err, []string{ReleaseWhileCreating, bce.EINTERNAL_ERROR}) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, raw)
			response := raw.(*api.ReleasePrepaidInstanceResponse)
			if !response.InstanceRefundFlag {
				return resource.NonRetryableError(fmt.Errorf("release prepaid instance failed: %+v", response))
			}
			return nil
		})
	}

	if err != nil {
		if IsExceptedErrors(err, BccNotFound) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_instance", action, BCESDKGoERROR)
	}

	stateConf := buildStateConf(
		[]string{string(api.InstanceStatusRunning), string(api.InstanceStatusStopping), string(api.InstanceStatusStopped)},
		[]string{string(api.InstanceStatusDeleted), string(api.InstanceStatusExpired), string(api.InstanceStatusRecycled)},
		d.Timeout(schema.TimeoutDelete),
		bccService.InstanceStateRefresh(instanceId),
	)
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_instance", action, BCESDKGoERROR)
	}

	return nil
}

func buildBaiduCloudInstanceArgs(d *schema.ResourceData, meta interface{}) (*api.CreateInstanceArgs, error) {
	request := &api.CreateInstanceArgs{
		ClientToken: buildClientToken(),
	}

	if imageID, ok := d.GetOk("image_id"); ok {
		request.ImageId = imageID.(string)
	}

	if name, ok := d.GetOk("name"); ok {
		request.Name = name.(string)
	}

	if zoneName, ok := d.GetOk("availability_zone"); ok {
		request.ZoneName = zoneName.(string)
	}

	if instanceType, ok := d.GetOk("instance_type"); ok {
		it := instanceType.(string)
		request.InstanceType = api.InstanceType(it)
	}

	billingRequest := api.Billing{
		PaymentTiming: api.PaymentTimingType(""),
		Reservation:   &api.Reservation{},
	}
	if p, ok := d.GetOk("payment_timing"); ok {
		paymentTiming := api.PaymentTimingType(p.(string))
		billingRequest.PaymentTiming = paymentTiming
	}
	if billingRequest.PaymentTiming == api.PaymentTimingPrePaid {
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
		// if the field is set, then auto-renewal is effective.
		if v, ok := d.GetOk("auto_renew_time_unit"); ok {
			request.AutoRenewTimeUnit = v.(string)
			if v, ok := d.GetOk("auto_renew_time_length"); ok {
				request.AutoRenewTime = v.(int)
			}
		}
	}
	request.Billing = billingRequest

	if adminPass, ok := d.GetOk("admin_pass"); ok {
		request.AdminPass = adminPass.(string)
	}

	if cpuCount, ok := d.GetOk("cpu_count"); ok {
		request.CpuCount = cpuCount.(int)
	}

	if memoryCapacityInGB, ok := d.GetOk("memory_capacity_in_gb"); ok {
		request.MemoryCapacityInGB = memoryCapacityInGB.(int)
	}

	if rootDiskSizeInGb, ok := d.GetOk("root_disk_size_in_gb"); ok {
		request.RootDiskSizeInGb = rootDiskSizeInGb.(int)
	}

	if rootDiskStorageType, ok := d.GetOk("root_disk_storage_type"); ok {
		dst := rootDiskStorageType.(string)
		request.RootDiskStorageType = api.StorageType(dst)
	}

	if v, ok := d.GetOk("ephemeral_disks"); ok {
		disks := v.([]interface{})
		var ephemeralDiskRequests []api.EphemeralDisk
		for iDisk := range disks {
			disk := disks[iDisk].(map[string]interface{})

			ephemeralDiskRequest := api.EphemeralDisk{
				SizeInGB:    disk["size_in_gb"].(int),
				StorageType: api.StorageType(disk["storage_type"].(string)),
			}

			ephemeralDiskRequests = append(ephemeralDiskRequests, ephemeralDiskRequest)
		}
		request.EphemeralDisks = ephemeralDiskRequests
	}

	if v, ok := d.GetOk("cds_disks"); ok {
		cdsList := v.([]interface{})
		cdsRequests := make([]api.CreateCdsModel, len(cdsList))
		for iCds := range cdsList {
			cds := cdsList[iCds].(map[string]interface{})

			cdsRequest := api.CreateCdsModel{
				CdsSizeInGB: cds["cds_size_in_gb"].(int),
				StorageType: api.StorageType(cds["storage_type"].(string)),
				SnapShotId:  cds["snapshot_id"].(string),
			}

			cdsRequests[iCds] = cdsRequest
		}
		request.CreateCdsList = cdsRequests
	}

	if dedicateHostId, ok := d.GetOk("dedicate_host_id"); ok {
		request.DedicateHostId = dedicateHostId.(string)
	}

	if keypairId, ok := d.GetOk("keypair_id"); ok {
		request.KeypairId = keypairId.(string)
	}

	if subnetId, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = subnetId.(string)
	}

	if gpuCard, ok := d.GetOk("gpu_card"); ok {
		request.GpuCard = gpuCard.(string)
	}

	if fpgaCard, ok := d.GetOk("fpga_card"); ok {
		request.FpgaCard = fpgaCard.(string)
	}

	if cardCount, ok := d.GetOk("card_count"); ok {
		request.CardCount = cardCount.(string)
	}

	if relationTag, ok := d.GetOk("relation_tag"); ok && relationTag.(bool) {
		request.RelationTag = relationTag.(bool)
	}

	if userData, ok := d.GetOk("user_data"); ok {
		request.UserData = userData.(string)
	}

	if v, ok := d.GetOk("tags"); ok {
		request.Tags = tranceTagMapToModel(v.(map[string]interface{}))
	}

	deploysetIds := make([]string, 0)
	v, ok := d.GetOk("deploy_set_ids")
	if ok {
		for _, value := range v.(*schema.Set).List() {
			deploysetIds = append(deploysetIds, value.(string))
		}
	}
	request.DeployIdList = deploysetIds

	if v, ok := d.GetOk("hostname"); ok && v.(string) != "" {
		request.Hostname = v.(string)
	}

	if v, ok := d.GetOk("is_open_hostname_domain"); ok {
		request.IsOpenHostnameDomain = v.(bool)
	}

	if v, ok := d.GetOk("is_open_ipv6"); ok {
		request.IsOpenIpv6 = v.(bool)
	}

	if v, ok := d.GetOk("cds_auto_renew"); ok {
		request.CdsAutoRenew = v.(bool)
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request.ResGroupId = v.(string)
	}

	return request, nil
}

func buildBaiduCloudInstanceBySpecArgs(d *schema.ResourceData, meta interface{}) (*api.CreateInstanceBySpecArgs, error) {
	request := &api.CreateInstanceBySpecArgs{
		ClientToken: buildClientToken(),
	}

	if imageID, ok := d.GetOk("image_id"); ok {
		request.ImageId = imageID.(string)
	}

	if name, ok := d.GetOk("name"); ok {
		request.Name = name.(string)
	}

	if zoneName, ok := d.GetOk("availability_zone"); ok {
		request.ZoneName = zoneName.(string)
	}

	if instanceSpec, ok := d.GetOk("instance_spec"); ok {
		request.Spec = instanceSpec.(string)
	}

	billingRequest := api.Billing{
		PaymentTiming: api.PaymentTimingType(""),
		Reservation:   &api.Reservation{},
	}
	if p, ok := d.GetOk("payment_timing"); ok {
		paymentTiming := api.PaymentTimingType(p.(string))
		billingRequest.PaymentTiming = paymentTiming
	}
	if billingRequest.PaymentTiming == api.PaymentTimingPrePaid {
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
		// if the field is set, then auto-renewal is effective.
		if v, ok := d.GetOk("auto_renew_time_unit"); ok {
			request.AutoRenewTimeUnit = v.(string)
			if v, ok := d.GetOk("auto_renew_time_length"); ok {
				request.AutoRenewTime = v.(int)
			}
		}
	}
	request.Billing = billingRequest

	if adminPass, ok := d.GetOk("admin_pass"); ok {
		request.AdminPass = adminPass.(string)
	}

	if rootDiskSizeInGb, ok := d.GetOk("root_disk_size_in_gb"); ok {
		request.RootDiskSizeInGb = rootDiskSizeInGb.(int)
	}

	if keypairId, ok := d.GetOk("keypair_id"); ok {
		request.KeypairId = keypairId.(string)
	}

	if rootDiskStorageType, ok := d.GetOk("root_disk_storage_type"); ok {
		dst := rootDiskStorageType.(string)
		request.RootDiskStorageType = api.StorageType(dst)
	}

	if v, ok := d.GetOk("ephemeral_disks"); ok {
		disks := v.([]interface{})
		var ephemeralDiskRequests []api.EphemeralDisk
		for iDisk := range disks {
			disk := disks[iDisk].(map[string]interface{})

			ephemeralDiskRequest := api.EphemeralDisk{
				SizeInGB:    disk["size_in_gb"].(int),
				StorageType: api.StorageType(disk["storage_type"].(string)),
			}

			ephemeralDiskRequests = append(ephemeralDiskRequests, ephemeralDiskRequest)
		}
		request.EphemeralDisks = ephemeralDiskRequests
	}

	if v, ok := d.GetOk("cds_disks"); ok {
		cdsList := v.([]interface{})
		cdsRequests := make([]api.CreateCdsModel, len(cdsList))
		for iCds := range cdsList {
			cds := cdsList[iCds].(map[string]interface{})

			cdsRequest := api.CreateCdsModel{
				CdsSizeInGB: cds["cds_size_in_gb"].(int),
				StorageType: api.StorageType(cds["storage_type"].(string)),
				SnapShotId:  cds["snapshot_id"].(string),
			}

			cdsRequests[iCds] = cdsRequest
		}
		request.CreateCdsList = cdsRequests
	}

	if subnetId, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = subnetId.(string)
	}

	if relationTag, ok := d.GetOk("relation_tag"); ok && relationTag.(bool) {
		request.RelationTag = relationTag.(bool)
	}

	if userData, ok := d.GetOk("user_data"); ok {
		request.UserData = userData.(string)
	}

	if v, ok := d.GetOk("tags"); ok {
		request.Tags = tranceTagMapToModel(v.(map[string]interface{}))
	}
	deploysetIds := make([]string, 0)
	v, ok := d.GetOk("deploy_set_ids")
	if ok {
		for _, value := range v.(*schema.Set).List() {
			deploysetIds = append(deploysetIds, value.(string))
		}
	}
	request.DeployIdList = deploysetIds

	if v, ok := d.GetOk("hostname"); ok && v.(string) != "" {
		request.Hostname = v.(string)
	}

	if v, ok := d.GetOk("is_open_hostname_domain"); ok {
		request.IsOpenHostnameDomain = v.(bool)
	}

	if v, ok := d.GetOk("is_open_ipv6"); ok {
		request.IsOpenIpv6 = v.(bool)
	}

	if v, ok := d.GetOk("cds_auto_renew"); ok {
		request.CdsAutoRenew = v.(bool)
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request.ResGroupId = v.(string)
	}

	return request, nil
}

func updateInstanceAttribute(d *schema.ResourceData, meta interface{}, instanceID string) error {
	action := "Update Instance attribute " + instanceID
	client := meta.(*connectivity.BaiduClient)

	if d.HasChange("name") {
		modifyInstanceAttributeArgs := &api.ModifyInstanceAttributeArgs{}
		modifyInstanceAttributeArgs.Name = d.Get("name").(string)

		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			_, err := client.WithBccClient(func(bccClient *bcc.Client) (interface{}, error) {
				return nil, bccClient.ModifyInstanceAttribute(instanceID, modifyInstanceAttributeArgs)
			})
			if err != nil {
				if IsExceptedErrors(err, []string{OperationDenied, bce.EINTERNAL_ERROR}) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, modifyInstanceAttributeArgs)
			return nil
		})

		if err != nil {
			if NotFoundError(err) {
				d.SetId("")
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_instance", action, BCESDKGoERROR)
		}

		d.SetPartial("name")
	}

	return nil
}

func updateInstanceDescription(d *schema.ResourceData, meta interface{}, instanceID string) error {
	action := "Update Instance Description " + instanceID
	client := meta.(*connectivity.BaiduClient)

	if d.HasChange("description") {
		modifyInstanceDescArgs := &api.ModifyInstanceDescArgs{}
		modifyInstanceDescArgs.Description = d.Get("description").(string)

		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			_, err := client.WithBccClient(func(bccClient *bcc.Client) (interface{}, error) {
				return nil, bccClient.ModifyInstanceDesc(instanceID, modifyInstanceDescArgs)
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
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_instance", action, BCESDKGoERROR)
		}

		d.SetPartial("description")
	}

	return nil
}

func updateInstanceImage(d *schema.ResourceData, meta interface{}, instanceID string) error {
	action := "Update instance image " + instanceID
	client := meta.(*connectivity.BaiduClient)
	bccService := &BccService{client}

	if d.HasChange("image_id") {
		args := &api.RebuildInstanceArgs{
			ImageId: d.Get("image_id").(string),
		}
		if adminPass, ok := d.GetOk("admin_pass"); ok {
			args.AdminPass = adminPass.(string)
		}

		if _, err := client.WithBccClient(func(bccClient *bcc.Client) (i interface{}, e error) {
			return nil, bccClient.RebuildInstance(instanceID, args)
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_instance", action, BCESDKGoERROR)
		}

		stateConf := buildStateConf(
			[]string{string(api.InstanceStatusStarting), string(api.InstanceStatusImageProcessing), string(api.InstanceStatusSnapshotProcessing)},
			[]string{string(api.InstanceStatusRunning)},
			d.Timeout(schema.TimeoutUpdate),
			bccService.InstanceStateRefresh(instanceID),
		)
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_instance", action, BCESDKGoERROR)
		}

		d.SetPartial("image_id")
		d.SetPartial("admin_pass")
	}

	return nil
}

func updateInstanceAdminPass(d *schema.ResourceData, meta interface{}, instanceID string) error {
	action := "Update Instance admin pass " + instanceID
	client := meta.(*connectivity.BaiduClient)
	bccService := &BccService{client}

	if d.HasChange("admin_pass") {
		args := &api.ChangeInstancePassArgs{
			AdminPass: d.Get("admin_pass").(string),
		}

		if _, err := client.WithBccClient(func(bccClient *bcc.Client) (i interface{}, e error) {
			return nil, bccClient.ChangeInstancePass(instanceID, args)
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_instance", action, BCESDKGoERROR)
		}

		stateConf := buildStateConf(
			[]string{string(api.InstanceStatusStarting)},
			[]string{string(api.InstanceStatusRunning)},
			d.Timeout(schema.TimeoutUpdate),
			bccService.InstanceStateRefresh(instanceID),
		)
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_instance", action, BCESDKGoERROR)
		}

		d.SetPartial("admin_pass")
	}

	return nil
}

func updateInstanceCapacity(d *schema.ResourceData, meta interface{}, instanceID string) error {
	action := "Update instance capacity " + instanceID
	client := meta.(*connectivity.BaiduClient)
	bccService := BccService{client}

	if d.HasChange("cpu_count") || d.HasChange("memory_capacity_in_gb") || d.HasChange("ephemeral_disks") {
		args := &api.ResizeInstanceArgs{
			ClientToken: buildClientToken(),
		}

		cpuCount := d.Get("cpu_count").(int)
		args.CpuCount = cpuCount

		memoryCapacityInGB := d.Get("memory_capacity_in_gb").(int)
		args.MemoryCapacityInGB = memoryCapacityInGB

		if v, ok := d.GetOk("ephemeral_disks"); ok {
			disks := v.([]interface{})
			ephemeralDiskRequests := make([]api.EphemeralDisk, 0, len(disks))
			for iDisk := range disks {
				disk := disks[iDisk].(map[string]interface{})

				ephemeralDiskRequest := api.EphemeralDisk{
					SizeInGB:    disk["size_in_gb"].(int),
					StorageType: api.StorageType(disk["storage_type"].(string)),
				}

				ephemeralDiskRequests = append(ephemeralDiskRequests, ephemeralDiskRequest)
			}
			args.EphemeralDisks = ephemeralDiskRequests
		}

		if _, err := client.WithBccClient(func(bccClient *bcc.Client) (i interface{}, e error) {
			return nil, bccClient.ResizeInstance(instanceID, args)
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_instance", action, BCESDKGoERROR)
		}

		stateConf := buildStateConf(
			[]string{string(api.InstanceStatusScaling)},
			[]string{string(api.InstanceStatusRunning)},
			d.Timeout(schema.TimeoutUpdate),
			bccService.InstanceStateRefresh(instanceID),
		)
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_instance", action, BCESDKGoERROR)
		}

		d.SetPartial("cpu_count")
		d.SetPartial("memory_capacity_in_gb")
		d.SetPartial("ephemeral_disks")
	}

	return nil
}

func updateInstanceSpec(d *schema.ResourceData, meta interface{}, instanceID string) error {
	action := "Update instance spec " + instanceID
	client := meta.(*connectivity.BaiduClient)
	bccService := BccService{client}

	if d.HasChange("instance_spec") {
		args := &api.ResizeInstanceArgs{
			ClientToken: buildClientToken(),
			Spec:        d.Get("instance_spec").(string),
		}

		if _, err := client.WithBccClient(func(bccClient *bcc.Client) (i interface{}, e error) {
			return nil, bccClient.ResizeInstanceBySpec(instanceID, args)
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_instance", action, BCESDKGoERROR)
		}

		stateConf := buildStateConf(
			[]string{string(api.InstanceStatusScaling)},
			[]string{string(api.InstanceStatusRunning)},
			d.Timeout(schema.TimeoutUpdate),
			bccService.InstanceStateRefresh(instanceID),
		)
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_instance", action, BCESDKGoERROR)
		}

		d.SetPartial("instance_spec")
	}

	return nil
}

func updateInstanceSecurityGroups(d *schema.ResourceData, meta interface{}, instanceID string) error {
	action := "Update instance security groups " + instanceID
	client := meta.(*connectivity.BaiduClient)

	if d.HasChange("security_groups") {
		o, n := d.GetChange("security_groups")

		os := o.(*schema.Set)
		ns := n.(*schema.Set)

		bindSGs := ns.Difference(os).List()
		unbindSGs := os.Difference(ns).List()

		// Each instance can be associated with 10 security groups at most and 1 security groups at least.
		for _, sg := range bindSGs {
			// bind security groups
			if _, err := client.WithBccClient(func(bccClient *bcc.Client) (i interface{}, e error) {
				return nil, bccClient.BindSecurityGroup(instanceID, sg.(string))
			}); err != nil {
				return WrapErrorf(err, DefaultErrorMsg, "baiducloud_instance", action, BCESDKGoERROR)
			}
		}
		for _, sg := range unbindSGs {
			// unbind security groups
			if _, err := client.WithBccClient(func(bccClient *bcc.Client) (i interface{}, e error) {
				return nil, bccClient.UnBindSecurityGroup(instanceID, sg.(string))
			}); err != nil {
				return WrapErrorf(err, DefaultErrorMsg, "baiducloud_instance", action, BCESDKGoERROR)
			}
		}

		d.SetPartial("security_groups")
	}

	return nil
}

func updateInstanceEnterpriseSecurityGroups(d *schema.ResourceData, meta interface{}, instanceID string) error {
	action := "Update instance enterprise security groups " + instanceID
	client := meta.(*connectivity.BaiduClient)

	if d.HasChange("enterprise_security_groups") {
		newGroupIds := expandStringSet(d.Get("enterprise_security_groups").(*schema.Set))

		request := &api.ReplaceSgV2Req{}
		request.SecurityGroupType = "enterprise"
		request.InstanceIds = []string{instanceID}
		request.SecurityGroupIds = newGroupIds
		if _, err := client.WithBccClient(func(bccClient *bcc.Client) (i interface{}, e error) {
			return bccClient.InstanceReplaceSecurityGroup(request)
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_instance", action, BCESDKGoERROR)
		}

		d.SetPartial("enterprise_security_groups")
	}

	return nil
}

func updateInstanceSubnet(d *schema.ResourceData, meta interface{}, instanceID string) error {
	action := "Update instance subnet " + instanceID
	client := meta.(*connectivity.BaiduClient)
	bccService := BccService{client}

	if d.HasChange("subnet_id") {
		args := &api.InstanceChangeSubnetArgs{
			InstanceId: instanceID,
			Reboot:     true,
		}
		if v, ok := d.GetOk("subnet_id"); ok {
			args.SubnetId = v.(string)
		}

		_, err := client.WithBccClient(func(bccClient *bcc.Client) (i interface{}, e error) {
			return nil, bccClient.InstanceChangeSubnet(args)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_instance", action, BCESDKGoERROR)
		}

		stateConf := buildStateConf(
			[]string{string(api.InstanceStatusStopping), string(api.InstanceStatusStopped), string(api.InstanceStatusStarting), InstanceStateChangeSubnet},
			[]string{string(api.InstanceStatusRunning)},
			d.Timeout(schema.TimeoutUpdate),
			bccService.InstanceStateRefresh(instanceID),
		)
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_instance", action, BCESDKGoERROR)
		}

		d.SetPartial("subnet_id")
	}

	return nil
}

func updateInstanceAction(d *schema.ResourceData, meta interface{}, instanceID string) error {
	action := "Update instance action " + instanceID
	client := meta.(*connectivity.BaiduClient)
	bccService := BccService{client}

	if d.HasChange("action") {
		act := d.Get("action").(string)
		addDebug(action, act)

		if act == INSTANCE_ACTION_START {
			if err := bccService.StartInstance(instanceID, d.Timeout(schema.TimeoutUpdate)); err != nil {
				return err
			}
		} else if act == INSTANCE_ACTION_STOP {
			stopWithNoCharge := d.Get("stop_with_no_charge").(bool)
			if err := bccService.StopInstance(instanceID, stopWithNoCharge, d.Timeout(schema.TimeoutUpdate)); err != nil {
				return err
			}
		}

		d.SetPartial("action")
	}

	return nil
}

func updateInstanceDeploy(d *schema.ResourceData, meta interface{}, instanceID string) error {
	action := "Update instance deploy sets " + instanceID
	client := meta.(*connectivity.BaiduClient)
	if d.HasChange("deploy_set_ids") {
		v, ok := d.GetOk("deploy_set_ids")
		deps := make([]string, 0)
		if ok {
			deploySets := v.(*schema.Set).List()
			for _, dep := range deploySets {
				// update deploy sets
				deps = append(deps, dep.(string))
			}

		}
		if _, err := client.WithBccClient(func(bccClient *bcc.Client) (i interface{}, e error) {
			req := &api.UpdateInstanceDeployArgs{
				ClientToken: buildClientToken(),
			}
			req.DeploySetIds = deps
			req.InstanceId = instanceID
			err, _ := bccClient.UpdateInstanceDeploySet(req)
			return nil, err
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_instance", action, BCESDKGoERROR)
		}
		d.SetPartial("deploy_set_ids")
	}
	return nil
}

func updateInstanceHostname(d *schema.ResourceData, meta interface{}, instanceID string) error {
	action := "Update instance hostname " + instanceID
	client := meta.(*connectivity.BaiduClient)
	if d.HasChanges("hostname", "is_open_hostname_domain") {
		modifyInstanceHostnameArgs := &api.ModifyInstanceHostnameArgs{}
		modifyInstanceHostnameArgs.Hostname = d.Get("hostname").(string)
		modifyInstanceHostnameArgs.IsOpenHostnameDomain = d.Get("is_open_hostname_domain").(bool)
		modifyInstanceHostnameArgs.Reboot = true

		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			_, err := client.WithBccClient(func(bccClient *bcc.Client) (interface{}, error) {
				return nil, bccClient.ModifyInstanceHostname(instanceID, modifyInstanceHostnameArgs)
			})
			if err != nil {
				if IsExceptedErrors(err, []string{OperationDenied, bce.EINTERNAL_ERROR}) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, modifyInstanceHostnameArgs)
			return nil
		})

		if err != nil {
			if NotFoundError(err) {
				d.SetId("")
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_instance", action, BCESDKGoERROR)
		}

		d.SetPartial("hostname")
		d.SetPartial("is_open_hostname_domain")
	}
	return nil

}

func updateInstancePaymentTiming(d *schema.ResourceData, meta interface{}, instanceID string) error {
	if d.HasChange("payment_timing") {
		action := "Update instance payment timing " + instanceID
		client := meta.(*connectivity.BaiduClient)
		newValue := d.Get("payment_timing").(string)
		if newValue == "Prepaid" {
			if _, ok := d.GetOk("reservation.reservation_length"); !ok {
				return fmt.Errorf("please set 'reservation.reservation_length' before changing payment_timing to 'Prepaid'")
			}
			reservationLength, _ := strconv.Atoi(d.Get("reservation.reservation_length").(string))
			prepayConfig := api.PrepayConfig{
				InstanceId: instanceID,
				Duration:   reservationLength,
			}
			if v, ok := d.GetOk("auto_renew_time_unit"); ok {
				prepayConfig.AutoRenew = true
				autoRenewTimeLength := d.Get("auto_renew_time_length").(int)
				autoRenewTimeUnit := v.(string)
				if autoRenewTimeUnit == "year" {
					autoRenewTimeLength *= 12
				}
				prepayConfig.AutoRenewPeriod = autoRenewTimeLength
			}
			args := &api.BatchChangeInstanceToPrepayArgs{
				Config: []api.PrepayConfig{prepayConfig},
			}
			addDebug(action, args)
			_, err := client.WithBccClient(func(bccClient *bcc.Client) (interface{}, error) {
				return bccClient.BatchChangeInstanceToPrepay(args)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, "baiducloud_instance", action, BCESDKGoERROR)
			}

		} else if newValue == "Postpaid" {
			postpayConfig := api.PostpayConfig{
				InstanceId:    instanceID,
				EffectiveType: "AtOnce",
			}
			args := &api.BatchChangeInstanceToPostpayArgs{
				Config: []api.PostpayConfig{postpayConfig},
			}
			addDebug(action, args)
			_, err := client.WithBccClient(func(bccClient *bcc.Client) (interface{}, error) {
				return bccClient.BatchChangeInstanceToPostpay(args)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, "baiducloud_instance", action, BCESDKGoERROR)
			}
		}

		d.SetPartial("payment_timing")
	}
	return nil

}

func updateInstanceAutoRenew(d *schema.ResourceData, meta interface{}, instanceID string) error {
	if d.HasChanges("auto_renew_time_unit", "auto_renew_time_length") {
		client := meta.(*connectivity.BaiduClient)

		deleteRule := func() error {
			_, err := client.WithBccClient(func(bccClient *bcc.Client) (interface{}, error) {
				args := &api.BccDeleteAutoRenewArgs{
					InstanceId: instanceID,
					RenewEip:   d.Get("sync_eip_auto_renew_rule").(bool),
				}
				return nil, bccClient.BatchDeleteAutoRenewRules(args)
			})
			return err

		}

		createRule := func() error {
			args := &api.BccCreateAutoRenewArgs{
				InstanceId:    instanceID,
				RenewTimeUnit: d.Get("auto_renew_time_unit").(string),
				RenewTime:     d.Get("auto_renew_time_length").(int),
				RenewEip:      d.Get("sync_eip_auto_renew_rule").(bool),
			}
			_, err := client.WithBccClient(func(bccClient *bcc.Client) (interface{}, error) {
				return nil, bccClient.BatchCreateAutoRenewRules(args)
			})
			return err
		}

		_, timeUnitSet := d.GetOk("auto_renew_time_unit")
		if timeUnitSet {
			autoRenewEnabled := d.Get("auto_renew").(bool)
			if autoRenewEnabled {
				err := deleteRule()
				if err != nil {
					return fmt.Errorf("delete auto renew rule failed: %s", err)
				}
				err = createRule()
				if err != nil {
					return fmt.Errorf("create auto renew rule failed: %s", err)
				}
			}
		} else {
			err := deleteRule()
			if err != nil {
				return fmt.Errorf("delete auto renew rule failed: %s", err)
			}
		}

		d.SetPartial("auto_renew_time_unit")
		d.SetPartial("auto_renew_time_length")
		d.SetPartial("sync_eip_auto_renew_rule")
	}
	return nil
}
