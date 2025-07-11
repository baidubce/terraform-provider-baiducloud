/*
Use this resource to attach instances to a CCE InstanceGroup.

~> **NOTE:** After creation, instances may take several minutes to reach the `running` state.
Destroying this resource **does not** remove instances from the instance group.

Example Usage

```hcl
resource "baiducloud_ccev2_instance_group_attachment" "example" {
  cluster_id = "cce-example"
  instance_group_id = "cce-ig-example"
  existed_instances = ["i-example"]

  existed_instances_config {
    rebuild = true
    image_id = "m-example"
    admin_password = "pass@word"
  }
}
```
*/
package baiducloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	ccev2 "github.com/baidubce/bce-sdk-go/services/cce/v2"
	"github.com/baidubce/bce-sdk-go/services/cce/v2/types"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func resourceBaiduCloudCCEv2InstanceGroupAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudCCEv2InstanceGroupAttachmentCreate,
		Read:   resourceBaiduCloudCCEv2InstanceGroupAttachmentRead,
		Delete: resourceBaiduCloudCCEv2InstanceGroupAttachmentDelete,
		Update: resourceBaiduCloudCCEv2InstanceGroupAttachmentUpdate,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Description: "The ID of the CCE cluster.",
				ForceNew:    true,
				Required:    true,
			},
			"instance_group_id": {
				Type:        schema.TypeString,
				Description: "The ID of the instance group.",
				ForceNew:    true,
				Required:    true,
			},
			"existed_instances": {
				Type:          schema.TypeSet,
				Description:   "IDs of instances outside the cluster to be added. Requires `existed_instances_config`.",
				Optional:      true,
				MinItems:      1,
				ConflictsWith: []string{"existed_instances_in_cluster"},
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"existed_instances_config": {
				Type:        schema.TypeList,
				Description: "Configuration for adding instances from outside the cluster. Required with `existed_instances`.",
				Optional:    true,
				MinItems:    1,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"machine_type": {
							Type:         schema.TypeString,
							Description:  "Machine type. Valid values: `BCC`, `BBC`, `EBC`, `HPAS`. Defaults to `BCC`.",
							Optional:     true,
							Default:      "BCC",
							ValidateFunc: validation.StringInSlice([]string{"BCC", "BBC", "EBC", "HPAS"}, false),
						},
						"rebuild": {
							Type: schema.TypeBool,
							Description: "Whether to reinstall the operating system. This will reinstall the OS on the selected instances, " +
								"clearing all data on the system disk (irrecoverable). Data on cloud disks will not be affected. " +
								"Only 'true' is supported currently.",
							Optional:     true,
							Default:      true,
							ValidateFunc: ValidateTrueOnly,
						},
						"use_instance_group_config": {
							Type:         schema.TypeBool,
							Description:  "Whether to apply the instance group’s configuration. Only 'true' is supported currently.",
							Optional:     true,
							Default:      true,
							ValidateFunc: ValidateTrueOnly,
						},
						"use_instance_group_config_with_disk_info": {
							Type:        schema.TypeBool,
							Description: "Whether to apply the instance group’s disk mount configuration. Defaults to `false`.",
							Optional:    true,
							Default:     false,
						},
						"image_id": {
							Type:        schema.TypeString,
							Description: "Image ID used for rebuild.",
							Optional:    true,
						},
						"admin_password": {
							Type:          schema.TypeString,
							Description:   "Admin password for login.",
							Optional:      true,
							ConflictsWith: []string{"existed_instances_config.ssh_key_id"},
						},
						"ssh_key_id": {
							Type:          schema.TypeString,
							Description:   "Key pair ID for login.",
							Optional:      true,
							ConflictsWith: []string{"existed_instances_config.admin_password"},
						},
					},
				},
			},
			"existed_instances_in_cluster": {
				Type:          schema.TypeSet,
				Description:   "IDs of instances already in the cluster to be added to the instance group.",
				Optional:      true,
				ConflictsWith: []string{"existed_instances"},
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},

		CustomizeDiff: func(diff *schema.ResourceDiff, i interface{}) error {
			instances, instancesExists := diff.GetOk("existed_instances")
			instancesInCluster, instancesInClusterExists := diff.GetOk("existed_instances_in_cluster")
			instancesConfig, instancesConfigExists := diff.GetOk("existed_instances_config")

			if (!instancesExists || instances.(*schema.Set).Len() == 0) && (!instancesInClusterExists || instancesInCluster.(*schema.Set).Len() == 0) {
				return fmt.Errorf("'existed_instances' and 'existed_instances_in_cluster' cannot both be empty")
			}
			if instancesExists && instances.(*schema.Set).Len() > 0 && (!instancesConfigExists || len(instancesConfig.([]interface{})) == 0) {
				return fmt.Errorf("'existed_instances_config' must be set when 'existed_instances' is set")
			}

			return nil
		},
	}
}

func resourceBaiduCloudCCEv2InstanceGroupAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	args, err := buildInstanceGroupAttachmentArgs(d, meta)
	if err != nil {
		log.Printf("BuildInstanceGroupAttachmentArgs Error:" + err.Error())
		return WrapError(err)
	}

	argsStr, _ := json.Marshal(args)
	addDebug("BuildInstanceGroupAttachmentArgs", string(argsStr))

	raw, err := client.WithCCEv2Client(func(client *ccev2.Client) (interface{}, error) {
		return client.AttachInstancesToInstanceGroup(args)
	})
	if err != nil {
		return err
	}

	resp := raw.(*ccev2.AttachInstancesToInstanceGroupResponse)
	d.SetId(resp.TaskID)

	if v, ok := d.GetOk("existed_instances"); ok && v.(*schema.Set).Len() > 0 {
		ccev2Service := Ccev2Service{client}
		instanceIds := expandStringSet(v.(*schema.Set))

		err := ccev2Service.waitForInstancesOperation([]string{EventStatusCreating}, []string{EventStatusCreated}, d.Timeout(schema.TimeoutCreate), instanceIds)
		if err != nil {
			return err
		}
	}

	return resourceBaiduCloudCCEv2InstanceGroupAttachmentRead(d, meta)
}

func resourceBaiduCloudCCEv2InstanceGroupAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceBaiduCloudCCEv2InstanceGroupAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceBaiduCloudCCEv2InstanceGroupAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func buildInstanceGroupAttachmentArgs(d *schema.ResourceData, meta interface{}) (*ccev2.AttachInstancesToInstanceGroupArgs, error) {
	args := &ccev2.AttachInstancesToInstanceGroupArgs{}
	args.ClusterID = d.Get("cluster_id").(string)
	args.InstanceGroupID = d.Get("instance_group_id").(string)

	request := &ccev2.AttachInstancesToInstanceGroupRequest{}

	if instancesIds, ok := d.GetOk("existed_instances"); ok && instancesIds.(*schema.Set).Len() > 0 {
		configRaw := d.Get("existed_instances_config")
		config := configRaw.([]interface{})[0].(map[string]interface{})

		instanceGroupConfig, err := getInstanceGroupConfig(d, meta)
		if err != nil {
			return nil, WrapErrorf(err, "GetInstanceGroupConfig Error")
		}

		imageId := ""
		if v, ok := d.GetOk("existed_instances_config.image_id"); ok {
			imageId = v.(string)
		}
		if imageId == "" {
			imageId = instanceGroupConfig.Spec.InstanceTemplate.ImageID
		}

		useInstanceGroupConfigWithDiskInfo := config["use_instance_group_config_with_disk_info"].(bool)
		var instances []*ccev2.InstanceSet
		for _, instanceId := range instancesIds.(*schema.Set).List() {
			instanceSpec := buildAttachmentInstanceSpec(config, instanceId.(string), imageId)
			if useInstanceGroupConfigWithDiskInfo {
				instanceSpec.InstanceResource = types.InstanceResource{CDSList: instanceGroupConfig.Spec.InstanceTemplate.InstanceResource.CDSList}
			}
			instance := &ccev2.InstanceSet{
				InstanceSpec: *instanceSpec,
			}
			instances = append(instances, instance)
		}

		request.Incluster = false
		request.ExistedInstances = instances
		request.UseInstanceGroupConfig = config["use_instance_group_config"].(bool)
		request.UseInstanceGroupConfigWithDiskInfo = useInstanceGroupConfigWithDiskInfo
	}

	if instancesIds, ok := d.GetOk("existed_instances_in_cluster"); ok && instancesIds.(*schema.Set).Len() > 0 {
		var instances []*ccev2.ExistedInstanceInCluster
		for _, instanceId := range instancesIds.(*schema.Set).List() {
			instance := &ccev2.ExistedInstanceInCluster{
				ExistedInstanceID: instanceId.(string),
			}
			instances = append(instances, instance)
		}

		request.Incluster = true
		request.ExistedInstancesInCluster = instances
	}

	args.Request = request
	return args, nil
}

func buildAttachmentInstanceSpec(config map[string]interface{}, instanceId, imageId string) *types.InstanceSpec {
	rebuild := config["rebuild"].(bool)

	spec := &types.InstanceSpec{
		Existed:     true,
		MachineType: types.MachineType(config["machine_type"].(string)),
		ClusterRole: types.ClusterRoleNode,
		ExistedOption: types.ExistedOption{
			ExistedInstanceID: instanceId,
			Rebuild:           &rebuild,
		},
		AdminPassword: config["admin_password"].(string),
		SSHKeyID:      config["ssh_key_id"].(string),
	}

	if rebuild {
		spec.ImageID = imageId
	}

	return spec
}

func getInstanceGroupConfig(d *schema.ResourceData, meta interface{}) (*ccev2.InstanceGroup, error) {
	client := meta.(*connectivity.BaiduClient)

	args := &ccev2.GetInstanceGroupArgs{
		ClusterID:       d.Get("cluster_id").(string),
		InstanceGroupID: d.Get("instance_group_id").(string),
	}
	raw, err := client.WithCCEv2Client(func(client *ccev2.Client) (interface{}, error) {
		return client.GetInstanceGroup(args)
	})
	if err != nil {
		return nil, err
	}
	return raw.(*ccev2.GetInstanceGroupResponse).InstanceGroup, nil
}
