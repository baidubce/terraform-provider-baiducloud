/*
Provide a resource to create an APPBLB Server Group.

Example Usage

```hcl
resource "baiducloud_appblb_server_group" "default" {
  name        = "testServerGroup"
  description = "this is a test Server Group"
  blb_id      = "lb-0d29a3f6"

  backend_server_list {
    instance_id = "i-VRKxC21a"
    weight = 50
  }

  port_list {
    port = 66
    type = "TCP"
  }
}
```
*/
package baiducloud

import (
	"fmt"
	"reflect"
	"time"

	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/appblb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func resourceBaiduCloudAppBlbServerGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudAppBlbServerGroupCreate,
		Read:   resourceBaiduCloudAppBlbServerGroupRead,
		Update: resourceBaiduCloudAppBlbServerGroupUpdate,
		Delete: resourceBaiduCloudAppBlbServerGroupDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"blb_id": {
				Type:        schema.TypeString,
				Description: "ID of the Application LoadBalance instance",
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Type:         schema.TypeString,
				Description:  "Name of the Server Group, length must be between 1 and 65 bytes, and will be automatically generated if not set",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 65),
			},
			"description": {
				Type:         schema.TypeString,
				Description:  "Server Group's description, length must be between 0 and 450 bytes, and support Chinese",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(0, 450),
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Server Group's status, see https://cloud.baidu.com/doc/BLB/s/Pjwvxnxdm/#blbstatus for detail",
				Computed:    true,
			},
			"backend_server_list": {
				Type:        schema.TypeSet,
				Description: "Server group bound backend server list",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Description: "Backend server instance ID",
							Required:    true,
						},
						"weight": {
							Type:         schema.TypeInt,
							Description:  "Backend server instance weight in this group, range from 0-100",
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 100),
						},
						"private_ip": {
							Type:        schema.TypeString,
							Description: "Backend server instance bind private ip",
							Computed:    true,
						},
						"port_list": {
							Type:        schema.TypeSet,
							Description: "Backend server open port list",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"listener_port": {
										Type:        schema.TypeInt,
										Description: "Listener port",
										Computed:    true,
									},
									"backend_port": {
										Type:        schema.TypeInt,
										Description: "Backend open port",
										Computed:    true,
									},
									"port_type": {
										Type:        schema.TypeString,
										Description: "Port protocol type",
										Computed:    true,
									},
									"health_check_port_type": {
										Type:        schema.TypeString,
										Description: "Health check port protocol type",
										Computed:    true,
									},
									"status": {
										Type:        schema.TypeString,
										Description: "Port status, include Alive/Dead/Unknown",
										Computed:    true,
									},
									"port_id": {
										Type:        schema.TypeString,
										Description: "Port id",
										Computed:    true,
									},
									"policy_id": {
										Type:        schema.TypeString,
										Description: "Port bind policy id",
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
			"port_list": {
				Type:        schema.TypeList,
				Description: "Server Group backend port list",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Description: "Server Group port id",
							Computed:    true,
						},
						"port": {
							Type:         schema.TypeInt,
							Description:  "App Server Group port, range from 1-65535",
							Required:     true,
							ValidateFunc: validatePort(),
						},
						"type": {
							Type:         schema.TypeString,
							Description:  "Server Group port protocol type, support TCP/UDP/HTTP",
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{TCP, UDP, HTTP}, false),
						},
						"status": {
							Type:        schema.TypeString,
							Description: "Server Group port status",
							Computed:    true,
						},
						"health_check": {
							Type:        schema.TypeString,
							Description: "Server Group port health check protocol, support TCP/UDP/HTTP, default same as port protocol type",
							// should be optional, but may cause bug, so set Required
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{TCP, UDP, HTTP}, false),
						},
						"health_check_port": {
							Type:         schema.TypeInt,
							Description:  "Server Group port health check port, default same as Server Group port",
							Computed:     true,
							Optional:     true,
							ValidateFunc: validatePort(),
						},
						"health_check_timeout_in_second": {
							Type:         schema.TypeInt,
							Description:  "Server Group health check timeout(second), support in [1, 60], default 3",
							Computed:     true,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 60),
						},
						"health_check_interval_in_second": {
							Type:         schema.TypeInt,
							Description:  "Server Group health check interval time(second), support in [1, 10], default 3",
							Computed:     true,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 10),
						},
						"health_check_down_retry": {
							Type:         schema.TypeInt,
							Description:  "Server Group health check down retry time, support in [2, 5], default 3",
							Computed:     true,
							Optional:     true,
							ValidateFunc: validation.IntBetween(2, 5),
						},
						"health_check_up_retry": {
							Type:         schema.TypeInt,
							Description:  "Server Group health check up retry time, support in [2, 5], default 3",
							Computed:     true,
							Optional:     true,
							ValidateFunc: validation.IntBetween(2, 5),
						},
						"health_check_normal_status": {
							Type:             schema.TypeString,
							Description:      "Server Group health check normal http status code, only useful when health_check is HTTP",
							Computed:         true,
							Optional:         true,
							DiffSuppressFunc: appServerGroupPortHealthCheckHTTPSuppressFunc,
						},
						"health_check_url_path": {
							Type:             schema.TypeString,
							Description:      "Server Group health check url path, only useful when health_check is HTTP",
							Computed:         true,
							Optional:         true,
							DiffSuppressFunc: appServerGroupPortHealthCheckHTTPSuppressFunc,
						},
						"udp_health_check_string": {
							Type:             schema.TypeString,
							Description:      "Server Group udp health check string, if type is UDP, this parameter is required",
							Computed:         true,
							Optional:         true,
							DiffSuppressFunc: appServerGroupPortHealthCheckUDPSuppressFunc,
						},
					},
				},
			},
		},
	}
}

func resourceBaiduCloudAppBlbServerGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	appblbService := APPBLBService{client}

	createArgs := buildBaiduCloudCreateAppBlbAppServerGroupArgs(d)
	blbId := d.Get("blb_id").(string)
	action := "Create AppBlb " + blbId + " AppServerGroup " + createArgs.Name

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithAppBLBClient(func(client *appblb.Client) (i interface{}, e error) {
			return client.CreateAppServerGroup(blbId, createArgs)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		addDebug(action, raw)
		response, _ := raw.(*appblb.CreateAppServerGroupResult)
		d.SetId(response.Id)

		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_appservergroup", action, BCESDKGoERROR)
	}

	stateConf := buildStateConf(
		APPBLBProcessingStatus,
		APPBLBAvailableStatus,
		d.Timeout(schema.TimeoutCreate),
		appblbService.AppServerGroupStateRefreshFunc(blbId, d.Id(), APPBLBFailedStatus))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapError(err)
	}

	return resourceBaiduCloudAppBlbServerGroupUpdate(d, meta)
}

func resourceBaiduCloudAppBlbServerGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	appblbService := APPBLBService{client}

	blbId := d.Get("blb_id").(string)
	id := d.Id()
	action := "Query APPBLB " + blbId + " App Server Group " + id

	group, err := appblbService.AppServerGroupDetail(blbId, id)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_appservergroup", action, BCESDKGoERROR)
	}
	addDebug(action, group)
	if err := d.Set("port_list", appblbService.FlattenAppServerGroupPortsToMap(group.PortList)); err != nil {
		return WrapError(err)
	}

	d.Set("status", group.Status)
	d.Set("name", group.Name)
	d.Set("description", group.Description)
	d.Set("blb_id", blbId)

	action = "Query APPBLB " + blbId + " App Server Group " + id + " backend list"
	servers, err := appblbService.AppServerGroupBlbRsDetail(blbId, id)
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_appservergroup", action, BCESDKGoERROR)
	}
	addDebug(action, group)
	if err := d.Set("backend_server_list", appblbService.FlattenAppBackendServersToMap(servers)); err != nil {
		return WrapError(err)
	}

	return nil
}

func resourceBaiduCloudAppBlbServerGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)
	if d.HasChange("port_list") {
		d.SetPartial("port_list")
		if err := updateAppServerGroupPortList(d, meta); err != nil {
			return err
		}
	}

	if d.HasChange("backend_server_list") {
		d.SetPartial("backend_server_list")
		if err := updateAppServerGroupRs(d, meta); err != nil {
			return err
		}
	}

	d.Partial(false)
	return resourceBaiduCloudAppBlbServerGroupRead(d, meta)
}

func resourceBaiduCloudAppBlbServerGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	blbId := d.Get("blb_id").(string)
	id := d.Id()
	deleteArgs := &appblb.DeleteAppServerGroupArgs{
		SgId: id,
	}
	action := "Delete APPBLB " + blbId + " App Server Group " + id

	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, err := client.WithAppBLBClient(func(client *appblb.Client) (i interface{}, e error) {
			return id, client.DeleteAppServerGroup(blbId, deleteArgs)
		})
		addDebug(action, id)

		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})
	if err != nil {
		if IsExceptedErrors(err, ObjectNotFound) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_appservergroup", action, BCESDKGoERROR)
	}

	return nil
}

func buildBaiduCloudCreateAppBlbAppServerGroupArgs(d *schema.ResourceData) *appblb.CreateAppServerGroupArgs {
	result := &appblb.CreateAppServerGroupArgs{
		ClientToken: buildClientToken(),
	}

	if v, ok := d.GetOk("name"); ok && v.(string) != "" {
		result.Name = v.(string)
	}

	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		result.Description = v.(string)
	}

	if v, ok := d.Get("backend_server_list").([]interface{}); ok && len(v) > 0 {
		for _, value := range v {
			m := value.(map[string]interface{})

			weight := m["weight"].(int)
			result.BackendServerList = append(result.BackendServerList, appblb.AppBackendServer{
				InstanceId: m["instance_id"].(string),
				Weight:     &weight,
			})
		}
	}

	return result
}

func updateAppServerGroupRs(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	appblbService := APPBLBService{client}

	blbId := d.Get("blb_id").(string)
	id := d.Id()

	o, n := d.GetChange("backend_server_list")
	os := o.(*schema.Set)
	ns := n.(*schema.Set)

	remove := os.Difference(ns).List()
	add := ns.Difference(os).List()

	addArgs, updateArgs, removeArgs := buildAppServerGroupRsWriteOpArgs(add, remove)
	if removeArgs != nil {
		removeArgs.SgId = id
		removeArgs.ClientToken = buildClientToken()

		if err := appblbService.DeleteAppServerGroupRs(blbId, removeArgs); err != nil {
			return WrapError(err)
		}

		if err := appblbService.WaitForServerGroupUpdateFinish(d); err != nil {
			return WrapError(err)
		}
	}

	if updateArgs != nil {
		updateArgs.SgId = id
		updateArgs.ClientToken = buildClientToken()
		args := &appblb.UpdateBlbRsArgs{
			BlbRsWriteOpArgs: *updateArgs,
		}

		if err := appblbService.UpdateAppServerGroupRs(blbId, args); err != nil {
			return WrapError(err)
		}

		if err := appblbService.WaitForServerGroupUpdateFinish(d); err != nil {
			return WrapError(err)
		}
	}

	if addArgs != nil {
		addArgs.SgId = id
		addArgs.ClientToken = buildClientToken()
		args := &appblb.CreateBlbRsArgs{
			BlbRsWriteOpArgs: *addArgs,
		}

		if err := appblbService.CreateAppServerGroupRs(blbId, args); err != nil {
			return WrapError(err)
		}

		if err := appblbService.WaitForServerGroupUpdateFinish(d); err != nil {
			return WrapError(err)
		}
	}

	return nil
}

func buildAppServerGroupRsWriteOpArgs(addList, removeList []interface{}) (add, update *appblb.BlbRsWriteOpArgs, remove *appblb.DeleteBlbRsArgs) {
	if len(addList) == 0 && len(removeList) == 0 {
		return nil, nil, nil
	}

	// only need to remove
	if len(addList) == 0 {
		return nil, nil, buildDeleteRsArgs(removeList)
	}

	// only need to add
	if len(removeList) == 0 {
		return buildRsWriteOpArgs(addList), nil, nil
	}

	// some add, some remove, other update
	addMap := make(map[string]interface{})
	for _, v := range addList {
		value := v.(map[string]interface{})

		addMap[value["instance_id"].(string)] = v
	}
	removeMap := make(map[string]interface{})
	for _, v := range removeList {
		value := v.(map[string]interface{})

		removeMap[value["instance_id"].(string)] = v
	}

	addPartList := make([]interface{}, 0)
	updatePartList := make([]interface{}, 0)
	removePartList := make([]interface{}, 0)
	for key, value := range addMap {
		if v, ok := removeMap[key]; ok {
			if !reflect.DeepEqual(value, v) {
				updatePartList = append(updatePartList, value)
			}
			delete(removeMap, key)
		} else {
			addPartList = append(addPartList, value)
		}
	}

	for _, value := range removeMap {
		removePartList = append(removePartList, value)
	}

	return buildRsWriteOpArgs(addPartList),
		buildRsWriteOpArgs(updatePartList),
		buildDeleteRsArgs(removePartList)
}

func buildDeleteRsArgs(list []interface{}) *appblb.DeleteBlbRsArgs {
	if len(list) == 0 {
		return nil
	}

	result := &appblb.DeleteBlbRsArgs{}
	for _, v := range list {
		removeValue := v.(map[string]interface{})
		result.BackendServerIdList = append(result.BackendServerIdList, removeValue["instance_id"].(string))
	}

	return result
}

func buildRsWriteOpArgs(list []interface{}) *appblb.BlbRsWriteOpArgs {
	if len(list) == 0 {
		return nil
	}

	result := &appblb.BlbRsWriteOpArgs{}
	for _, v := range list {
		writeOpValue := v.(map[string]interface{})

		weight := writeOpValue["weight"].(int)
		result.BackendServerList = append(result.BackendServerList, appblb.AppBackendServer{
			InstanceId: writeOpValue["instance_id"].(string),
			Weight:     &weight,
		})
	}

	return result
}

func updateAppServerGroupPortList(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	appblbService := APPBLBService{client}

	blbId := d.Get("blb_id").(string)
	id := d.Id()

	o, n := d.GetChange("port_list")
	os := o.([]interface{})
	ns := n.([]interface{})

	remove := portListConfigDifference(os, ns)
	add := portListConfigDifference(ns, os)

	addArgs, updateArgs, deleteArgs, err := buildBaiduCloudCreateAppBlbAppServerGroupPortArgs(add, remove)
	if err != nil {
		return err
	}

	if deleteArgs != nil {
		deleteArgs.SgId = id
		deleteArgs.ClientToken = buildClientToken()

		if err := appblbService.DeleteAppServerGroupPort(blbId, deleteArgs); err != nil {
			return WrapError(err)
		}

		if err := appblbService.WaitForServerGroupUpdateFinish(d); err != nil {
			return WrapError(err)
		}
	}

	for _, args := range addArgs {
		args.SgId = id
		args.ClientToken = buildClientToken()

		if err := appblbService.CreateAppServerGroupPort(blbId, &args); err != nil {
			return WrapError(err)
		}

		if err := appblbService.WaitForServerGroupUpdateFinish(d); err != nil {
			return WrapError(err)
		}
	}

	for _, args := range updateArgs {
		args.SgId = id
		args.ClientToken = buildClientToken()

		if err := appblbService.UpdateAppServerGroupPort(blbId, &args); err != nil {
			return WrapError(err)
		}

		if err := appblbService.WaitForServerGroupUpdateFinish(d); err != nil {
			return WrapError(err)
		}
	}

	return nil
}

func buildBaiduCloudCreateAppBlbAppServerGroupPortArgs(addList, removeList []interface{}) (
	add []appblb.CreateAppServerGroupPortArgs,
	update []appblb.UpdateAppServerGroupPortArgs,
	remove *appblb.DeleteAppServerGroupPortArgs,
	err error) {
	if len(addList) == 0 && len(removeList) == 0 {
		return nil, nil, nil, nil
	}

	// only need to remove
	if len(addList) == 0 {
		return nil, nil, buildDeleteAppServerGroupPortArgs(removeList), nil
	}

	// only need to add
	if len(removeList) == 0 {
		addArgList, err := buildCreateAppServerGroupPortArgs(addList)
		return addArgList, nil, nil, err
	}

	// some add, some remove, other update
	addMap := make(map[string]interface{})
	for _, v := range addList {
		value := v.(map[string]interface{})
		key := fmt.Sprintf("%d.%s", value["port"].(int), value["type"].(string))

		addMap[key] = v
	}
	removeMap := make(map[string]interface{})
	for _, v := range removeList {
		value := v.(map[string]interface{})
		key := fmt.Sprintf("%d.%s", value["port"].(int), value["type"].(string))

		removeMap[key] = v
	}

	addPartList := make([]interface{}, 0)
	removePartList := make([]interface{}, 0)
	updatePartList := make([]interface{}, 0)

	for key, value := range addMap {
		if v, ok := removeMap[key]; ok {
			if !reflect.DeepEqual(value, v) {
				vOldMap := v.(map[string]interface{})
				vNewMap := value.(map[string]interface{})

				vNewMap["id"] = vOldMap["id"]
				updatePartList = append(updatePartList, vNewMap)
			}
			delete(removeMap, key)
		} else {
			addPartList = append(addPartList, value)
		}
	}

	for _, value := range removeMap {
		removePartList = append(removePartList, value)
	}

	addArgList, addErr := buildCreateAppServerGroupPortArgs(addPartList)
	if addErr != nil {
		return nil, nil, nil, addErr
	}

	updateArgList, updateErr := buildUpdateAppServerGroupPortArgs(updatePartList)
	if updateErr != nil {
		return nil, nil, nil, updateErr
	}

	removeArgs := buildDeleteAppServerGroupPortArgs(removePartList)

	return addArgList, updateArgList, removeArgs, nil
}

func buildCreateAppServerGroupPortArgs(list []interface{}) ([]appblb.CreateAppServerGroupPortArgs, error) {
	if len(list) == 0 {
		return nil, nil
	}

	result := make([]appblb.CreateAppServerGroupPortArgs, 0)
	for _, v := range list {
		addValue := v.(map[string]interface{})
		addArgs := appblb.CreateAppServerGroupPortArgs{}

		addArgs.Port = uint16(addValue["port"].(int))
		addArgs.Type = addValue["type"].(string)

		if value, ok := addValue["health_check"]; ok {
			addArgs.HealthCheck = value.(string)
		}

		if addArgs.HealthCheck == "" {
			addArgs.HealthCheck = addArgs.Type
		}

		switch addArgs.Type {
		case TCP, UDP:
			if addArgs.HealthCheck != addArgs.Type {
				return nil, fmt.Errorf("%s port type should set %s healthcheck, but now is %s", addArgs.Type, addArgs.Type, addArgs.HealthCheck)
			}
		case HTTP:
			if !stringInSlice([]string{HTTP, TCP}, addArgs.HealthCheck) {
				return nil, fmt.Errorf("HTTP port type should set HTTP/TCP healthcheck, but now is %s", addArgs.HealthCheck)
			}
		default:
			return nil, fmt.Errorf("unsupport port type %s", addArgs.Type)
		}

		if value, ok := addValue["health_check_port"]; ok {
			addArgs.HealthCheckPort = value.(int)
		}

		if value, ok := addValue["health_check_url_path"]; ok {
			addArgs.HealthCheckUrlPath = value.(string)
		}

		if value, ok := addValue["health_check_timeout_in_second"]; ok {
			addArgs.HealthCheckTimeoutInSecond = value.(int)
		}

		if value, ok := addValue["health_check_interval_in_second"]; ok {
			addArgs.HealthCheckIntervalInSecond = value.(int)
		}

		if value, ok := addValue["health_check_down_retry"]; ok {
			addArgs.HealthCheckDownRetry = value.(int)
		}

		if value, ok := addValue["health_check_up_retry"]; ok {
			addArgs.HealthCheckUpRetry = value.(int)
		}

		if value, ok := addValue["health_check_normal_status"]; ok {
			addArgs.HealthCheckNormalStatus = value.(string)
		}

		if addArgs.HealthCheck == UDP {
			if value, ok := addValue["udp_health_check_string"]; ok && value.(string) != "" {
				addArgs.UdpHealthCheckString = value.(string)
			} else {
				return nil, fmt.Errorf("udp_health_check_string is required if type is UDP")
			}
		}
		result = append(result, addArgs)
	}

	return result, nil
}

func buildUpdateAppServerGroupPortArgs(list []interface{}) ([]appblb.UpdateAppServerGroupPortArgs, error) {
	if len(list) == 0 {
		return nil, nil
	}

	result := make([]appblb.UpdateAppServerGroupPortArgs, 0)
	for _, v := range list {
		updateValue := v.(map[string]interface{})
		updateArgs := appblb.UpdateAppServerGroupPortArgs{}
		updateArgs.PortId = updateValue["id"].(string)
		updateType := updateValue["type"].(string)

		if value, ok := updateValue["health_check"]; ok {
			updateArgs.HealthCheck = value.(string)
		}

		switch updateType {
		case TCP, UDP:
			if updateArgs.HealthCheck != updateType {
				return nil, fmt.Errorf("%s port type should set %s healthcheck, but now is %s", updateType, updateType, updateArgs.HealthCheck)
			}
		case HTTP:
			if !stringInSlice([]string{HTTP, TCP}, updateArgs.HealthCheck) {
				return nil, fmt.Errorf("HTTP port type should set HTTP/TCP healthcheck, but now is %s", updateArgs.HealthCheck)
			}
		default:
			return nil, fmt.Errorf("unsupport port type %s", updateType)
		}

		if value, ok := updateValue["health_check_port"]; ok {
			updateArgs.HealthCheckPort = value.(int)
		}

		if value, ok := updateValue["health_check_url_path"]; ok {
			updateArgs.HealthCheckUrlPath = value.(string)
		}

		if value, ok := updateValue["health_check_timeout_in_second"]; ok {
			updateArgs.HealthCheckTimeoutInSecond = value.(int)
		}

		if value, ok := updateValue["health_check_interval_in_second"]; ok {
			updateArgs.HealthCheckIntervalInSecond = value.(int)
		}

		if value, ok := updateValue["health_check_down_retry"]; ok {
			updateArgs.HealthCheckDownRetry = value.(int)
		}

		if value, ok := updateValue["health_check_up_retry"]; ok {
			updateArgs.HealthCheckUpRetry = value.(int)
		}

		if value, ok := updateValue["health_check_normal_status"]; ok {
			updateArgs.HealthCheckNormalStatus = value.(string)
		}

		if updateArgs.HealthCheck == UDP {
			if value, ok := updateValue["udp_health_check_string"]; ok && value.(string) != "" {
				updateArgs.UdpHealthCheckString = value.(string)
			} else {
				return nil, fmt.Errorf("udp_health_check_string is required if type is UDP")
			}
		}
		result = append(result, updateArgs)
	}

	return result, nil
}

func buildDeleteAppServerGroupPortArgs(list []interface{}) *appblb.DeleteAppServerGroupPortArgs {
	if len(list) == 0 {
		return nil
	}

	result := &appblb.DeleteAppServerGroupPortArgs{}
	for _, v := range list {
		removeValue := v.(map[string]interface{})
		result.PortIdList = append(result.PortIdList, removeValue["id"].(string))
	}

	return result
}

func portListConfigDifference(rootList []interface{}, compareList []interface{}) []interface{} {
	result := make([]interface{}, 0)

	compareMap := make(map[string]interface{})
	for _, compare := range compareList {
		cMap := compare.(map[string]interface{})

		key := fmt.Sprintf("%v_%v", cMap["port"], cMap["type"])
		compareMap[key] = compare
	}

	for _, root := range rootList {
		rMap := root.(map[string]interface{})
		key := fmt.Sprintf("%v_%v", rMap["port"], rMap["type"])

		if compare, ok := compareMap[key]; !ok || !reflect.DeepEqual(root, compare) {
			result = append(result, root)
		}
	}

	return result
}
