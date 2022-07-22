/*
Use this resource to bind to an instance and modify some of its attributes.
Note that this resource will not create a real instance, it is just a way to bind to a remote instance and modify its attributes.
If you wish to create more instances, please use baiducloud_ccev2_instance_group.

Example Usage

```hcl
resource "baiducloud_ccev2_instance" "default" {
  cluster_id        = "your-cluster-id"
  instance_id       = "your-instance-id"
  spec {
    cce_instance_priority = 0
  }
}
```
*/
package baiducloud

import (
	"errors"
	"log"
	"time"

	ccev2 "github.com/baidubce/bce-sdk-go/services/cce/v2"
	"github.com/baidubce/bce-sdk-go/services/cce/v2/types"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func resourceBaiduCloudCCEv2Instance() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudCCEv2InstanceCreate,
		Read:   resourceBaiduCloudCCEv2InstanceRead,
		Delete: resourceBaiduCloudCCEv2InstanceDelete,
		Update: resourceBaiduCloudCCEv2InstanceUpdate,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Description: "Cluster ID of this Instance.",
				ForceNew:    true,
				Required:    true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Description: "Cluster ID of this Instance.",
				ForceNew:    true,
				Required:    true,
			},
			"spec": {
				Type:        schema.TypeList,
				Description: "Instance Spec",
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cce_instance_priority": {
							Type:        schema.TypeInt,
							Description: "Priority of this instance.",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func resourceBaiduCloudCCEv2InstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	//如果用户有设定priority的话 把设置的值记录下来 在本方法结尾执行一次Update Priority的操作
	var userSetPriority *int
	if _, ok := d.Get("spec.0").(map[string]interface{})["cce_instance_priority"]; ok {
		priority := d.Get("spec.0").(map[string]interface{})["cce_instance_priority"].(int)
		userSetPriority = &priority
	}

	args, err := buildGetInstanceArgs(d)
	if err != nil {
		log.Printf("Build GetInstanceArgs Error:" + err.Error())
		return WrapError(err)
	}

	action := "Create CCEv2 Instance ClusterID:" + args.ClusterID + " InstanceID:" + args.InstanceID
	raw, err := client.WithCCEv2Client(func(client *ccev2.Client) (interface{}, error) {
		return client.GetInstance(args)
	})
	if err != nil {
		log.Printf("Get Instance Error:" + err.Error())
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2_instance", action, BCESDKGoERROR)
	}
	resp := raw.(*ccev2.GetInstanceResponse)

	err = setInstanceToState(resp.Instance, d)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2_instance", action, BCESDKGoERROR)
	}

	d.SetId(resp.Instance.Spec.CCEInstanceID)

	//如果用户设定了节点的Priority的话 就执行一次Update Priority的操作
	if userSetPriority != nil {
		resp.Instance.Spec.CCEInstancePriority = *userSetPriority
		argsUpdate := &ccev2.UpdateInstanceArgs{
			ClusterID:    resp.Instance.Spec.ClusterID,
			InstanceID:   resp.Instance.Spec.CCEInstanceID,
			InstanceSpec: resp.Instance.Spec,
		}
		rawUpdate, err := doUpdate(client, argsUpdate)

		if err != nil {
			log.Printf("Update Instance Error:" + err.Error())
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2_instance", action, BCESDKGoERROR)
		}
		_, ok := rawUpdate.(*ccev2.UpdateInstancesResponse)
		if !ok {
			err = errors.New("response format illegal")
			return err
		}
		//执行完Update之后立刻读取一次最新的节点数据
		return resourceBaiduCloudCCEv2InstanceRead(d, meta)
	}

	return nil
}

func resourceBaiduCloudCCEv2InstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	args, err := buildGetInstanceArgs(d)
	if err != nil {
		log.Printf("Build GetInstanceArgs Error:" + err.Error())
		return WrapError(err)
	}

	action := "Read CCEv2 Instance ClusterID:" + args.ClusterID + " InstanceID:" + args.InstanceID
	raw, err := client.WithCCEv2Client(func(client *ccev2.Client) (interface{}, error) {
		return client.GetInstance(args)
	})
	if err != nil {
		if NotFoundError(err) {
			log.Printf("Instance Not Found. " + err.Error())
			d.SetId("") //Resource Not Found, make the ID of resource to empty to delete it in state file.
			return nil
		}
		log.Printf("Get Instance Error:" + err.Error())
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2_instance", action, BCESDKGoERROR)
	}
	resp := raw.(*ccev2.GetInstanceResponse)

	err = setInstanceToState(resp.Instance, d)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2_instance", action, BCESDKGoERROR)
	}
	return nil
}

func resourceBaiduCloudCCEv2InstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	//1. 先获取节点最新数据 在最新数据上修改tf的设定
	argsGet, err := buildGetInstanceArgs(d)
	if err != nil {
		log.Printf("Build GetInstanceArgs Error:" + err.Error())
		return WrapError(err)
	}
	actionGet := "Read CCEv2 Instance ClusterID:" + argsGet.ClusterID + " InstanceID:" + argsGet.InstanceID
	rawGet, err := client.WithCCEv2Client(func(client *ccev2.Client) (interface{}, error) {
		return client.GetInstance(argsGet)
	})
	if err != nil {
		if NotFoundError(err) {
			log.Printf("Instance Not Found. " + err.Error())
			d.SetId("") //Resource Not Found, make the ID of resource to empty to delete it in state file.
			return nil
		}
		log.Printf("Get Instance Error:" + err.Error())
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2_instance", actionGet, BCESDKGoERROR)
	}
	respGet := rawGet.(*ccev2.GetInstanceResponse)

	//2. 把Update的信息写回原spec 然后执行Update
	argsUpdate, err := buildUpdateInstanceArgs(d, respGet.Instance.Spec)
	if err != nil {
		log.Printf("Build UpdateInstanceArgs Error:" + err.Error())
		return WrapError(err)
	}

	actionUpdate := "Update CCEv2 Instance ClusterID:" + argsUpdate.ClusterID + " InstanceID:" + argsUpdate.InstanceID
	rawUpdate, err := doUpdate(client, argsUpdate)

	if err != nil {
		log.Printf("Update Instance Error:" + err.Error())
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2_instance", actionUpdate, BCESDKGoERROR)
	}
	addDebug(actionUpdate, rawUpdate)
	_, ok := rawUpdate.(*ccev2.UpdateInstancesResponse)
	if !ok {
		err = errors.New("response format illegal")
		return err
	}

	return resourceBaiduCloudCCEv2InstanceRead(d, meta)
}

func doUpdate(client *connectivity.BaiduClient, argsUpdate *ccev2.UpdateInstanceArgs) (interface{}, error) {
	rawUpdate, err := client.WithCCEv2Client(func(client *ccev2.Client) (interface{}, error) {
		return client.UpdateInstance(argsUpdate)
	})

	return rawUpdate, err
}

func resourceBaiduCloudCCEv2InstanceDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func setInstanceToState(instance *ccev2.Instance, d *schema.ResourceData) error {
	if instance == nil {
		return errors.New("instance is nil")
	}

	//set spec
	specMaps, err := convertInstanceSpecToMaps(instance.Spec)
	if err != nil {
		return err
	}
	err = d.Set("spec", specMaps)
	if err != nil {
		return err
	}

	return nil
}

func convertInstanceSpecToMaps(spec *types.InstanceSpec) ([]interface{}, error) {
	result := make([]interface{}, 0)
	if spec == nil {
		return result, nil
	}

	specMap := make(map[string]interface{})
	specMap["cce_instance_priority"] = spec.CCEInstancePriority

	result = append(result, specMap)
	return result, nil
}

func buildGetInstanceArgs(d *schema.ResourceData) (*ccev2.GetInstanceArgs, error) {
	clusterID := d.Get("cluster_id").(string)
	instanceID := d.Get("instance_id").(string)

	if clusterID == "" || instanceID == "" {
		return nil, errors.New("cluster_id or instance_id empty")
	}

	args := &ccev2.GetInstanceArgs{
		ClusterID:  clusterID,
		InstanceID: instanceID,
	}
	return args, nil
}

func buildUpdateInstanceArgs(d *schema.ResourceData, oldInstanceSpec *types.InstanceSpec) (*ccev2.UpdateInstanceArgs, error) {
	clusterID := d.Get("cluster_id").(string)
	instanceID := d.Get("instance_id").(string)
	if clusterID == "" || instanceID == "" {
		return nil, errors.New("cluster_id or instance_id empty")
	}

	instanceSpecMaps := d.Get("spec.0").(map[string]interface{})
	instanceSpec, err := buildInstanceSpec(instanceSpecMaps)
	if err != nil {
		return nil, err
	}

	//update了哪些字段
	oldInstanceSpec.CCEInstancePriority = instanceSpec.CCEInstancePriority

	args := &ccev2.UpdateInstanceArgs{
		ClusterID:    clusterID,
		InstanceID:   instanceID,
		InstanceSpec: oldInstanceSpec,
	}
	return args, nil
}
