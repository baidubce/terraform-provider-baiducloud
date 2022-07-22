package baiducloud

import (
	"fmt"
	"strconv"

	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/appblb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func (s *APPBLBService) AppServerGroupDetail(blbId, sgId string) (*appblb.AppServerGroup, error) {
	describeArgs := &appblb.DescribeAppServerGroupArgs{}

	for {
		raw, err := s.client.WithAppBLBClient(func(client *appblb.Client) (i interface{}, e error) {
			return client.DescribeAppServerGroup(blbId, describeArgs)
		})
		if err != nil {
			return nil, WrapError(err)
		}

		result := raw.(*appblb.DescribeAppServerGroupResult)
		for _, group := range result.AppServerGroupList {
			if group.Id == sgId {
				return &group, nil
			}
		}

		if result.IsTruncated {
			describeArgs.Marker = result.NextMarker
			describeArgs.MaxKeys = result.MaxKeys
		} else {
			return nil, WrapError(fmt.Errorf(ResourceNotFound))
		}
	}
}

func (s *APPBLBService) AppServerGroupBlbRsDetail(blbId, sgId string) ([]appblb.AppBackendServer, error) {
	describeArgs := &appblb.DescribeBlbRsArgs{
		SgId: sgId,
	}

	result := make([]appblb.AppBackendServer, 0)
	for {
		raw, err := s.client.WithAppBLBClient(func(client *appblb.Client) (i interface{}, e error) {
			return client.DescribeBlbRs(blbId, describeArgs)
		})
		if err != nil {
			return nil, WrapError(err)
		}

		response := raw.(*appblb.DescribeBlbRsResult)
		result = append(result, response.BackendServerList...)

		if response.IsTruncated {
			describeArgs.Marker = response.NextMarker
			describeArgs.MaxKeys = response.MaxKeys
		} else {
			break
		}
	}

	return result, nil
}

func (s *APPBLBService) AppServerGroupStateRefreshFunc(blbId, sgId string, failState []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		result, err := s.AppServerGroupDetail(blbId, sgId)

		if err != nil {
			return nil, "", WrapError(err)
		}

		for _, state := range failState {
			if string(result.Status) == state {
				return result, string(result.Status), WrapError(Error(GetFailTargetStatus, result.Status))
			}
		}

		return result, string(result.Status), nil
	}
}

func (s *APPBLBService) FlattenAppServerGroupPortsToMap(group []appblb.AppServerGroupPort) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(group))

	for _, g := range group {
		gMap := map[string]interface{}{
			"id":                              g.Id,
			"port":                            g.Port,
			"type":                            g.Type,
			"status":                          g.Status,
			"health_check":                    g.HealthCheck,
			"health_check_port":               g.HealthCheckPort,
			"health_check_timeout_in_second":  g.HealthCheckTimeoutInSecond,
			"health_check_interval_in_second": g.HealthCheckIntervalInSecond,
			"health_check_down_retry":         g.HealthCheckDownRetry,
			"health_check_up_retry":           g.HealthCheckUpRetry,
			"health_check_normal_status":      g.HealthCheckNormalStatus,
			"health_check_url_path":           g.HealthCheckUrlPath,
			"udp_health_check_string":         g.UdpHealthCheckString,
		}

		result = append(result, gMap)
	}

	return result
}

func (s *APPBLBService) FlattenAppBackendServersToMap(servers []appblb.AppBackendServer) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(servers))

	for _, server := range servers {
		sMap := map[string]interface{}{
			"instance_id": server.InstanceId,
			"weight":      server.Weight,
			"private_ip":  server.PrivateIp,
			"port_list":   s.FlattenAppRsPortsToMap(server.PortList),
		}

		result = append(result, sMap)
	}

	return result
}

func (s *APPBLBService) FlattenAppRsPortsToMap(ports []appblb.AppRsPortModel) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(ports))

	for _, port := range ports {
		backendPort, _ := strconv.Atoi(port.BackendPort)
		pMap := map[string]interface{}{
			"listener_port":          port.ListenerPort,
			"backend_port":           backendPort,
			"port_type":              port.PortType,
			"health_check_port_type": port.HealthCheckPortType,
			"status":                 port.Status,
			"port_id":                port.PortId,
			"policy_id":              port.PolicyId,
		}

		result = append(result, pMap)
	}

	return result
}

func (s *APPBLBService) CreateAppServerGroupPort(blbId string, args *appblb.CreateAppServerGroupPortArgs) error {
	if args == nil {
		return nil
	}

	action := fmt.Sprintf("Create App Server Group Port %s: %d: %+v", args.Type, args.Port, args)

	raw, err := s.client.WithAppBLBClient(func(client *appblb.Client) (i interface{}, e error) {
		return client.CreateAppServerGroupPort(blbId, args)
	})
	addDebug(action, raw)

	if err != nil {
		if !IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_appservergroup", action, BCESDKGoERROR)
		}
	}

	return nil
}

func (s *APPBLBService) UpdateAppServerGroupPort(blbId string, args *appblb.UpdateAppServerGroupPortArgs) error {
	if args == nil {
		return nil
	}

	action := fmt.Sprintf("Update App Server Group Port %s", args.PortId)

	raw, err := s.client.WithAppBLBClient(func(client *appblb.Client) (i interface{}, e error) {
		return args.PortId, client.UpdateAppServerGroupPort(blbId, args)
	})
	addDebug(action, raw)

	if err != nil {
		if !IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_appservergroup", action, BCESDKGoERROR)
		}
	}

	return nil
}

func (s *APPBLBService) DeleteAppServerGroupPort(blbId string, args *appblb.DeleteAppServerGroupPortArgs) error {
	if args == nil {
		return nil
	}

	action := fmt.Sprintf("Create App Server Group Port %v", args.PortIdList)

	raw, err := s.client.WithAppBLBClient(func(client *appblb.Client) (i interface{}, e error) {
		return nil, client.DeleteAppServerGroupPort(blbId, args)
	})
	addDebug(action, raw)

	if err != nil {
		if !IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_appservergroup", action, BCESDKGoERROR)
		}
	}

	return nil
}

func (s *APPBLBService) CreateAppServerGroupRs(blbId string, args *appblb.CreateBlbRsArgs) error {
	if args == nil {
		return nil
	}

	action := fmt.Sprintf("Create App Server Group %v Rs", args.SgId)

	raw, err := s.client.WithAppBLBClient(func(client *appblb.Client) (i interface{}, e error) {
		return nil, client.CreateBlbRs(blbId, args)
	})
	addDebug(action, raw)

	if err != nil {
		if !IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_appservergroup", action, BCESDKGoERROR)
		}
	}

	return nil
}

func (s *APPBLBService) UpdateAppServerGroupRs(blbId string, args *appblb.UpdateBlbRsArgs) error {
	if args == nil {
		return nil
	}

	action := fmt.Sprintf("Update App Server Group %v Rs", args.SgId)

	raw, err := s.client.WithAppBLBClient(func(client *appblb.Client) (i interface{}, e error) {
		return nil, client.UpdateBlbRs(blbId, args)
	})
	addDebug(action, raw)

	if err != nil {
		if !IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_appservergroup", action, BCESDKGoERROR)
		}
	}

	return nil
}

func (s *APPBLBService) DeleteAppServerGroupRs(blbId string, args *appblb.DeleteBlbRsArgs) error {
	if args == nil {
		return nil
	}

	action := fmt.Sprintf("Delete App Server Group %v Rs", args.SgId)

	raw, err := s.client.WithAppBLBClient(func(client *appblb.Client) (i interface{}, e error) {
		return nil, client.DeleteBlbRs(blbId, args)
	})
	addDebug(action, raw)

	if err != nil {
		if !IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_appservergroup", action, BCESDKGoERROR)
		}
	}

	return nil
}

func (s *APPBLBService) WaitForServerGroupUpdateFinish(d *schema.ResourceData) error {
	stateConf := buildStateConf(
		APPBLBProcessingStatus,
		APPBLBAvailableStatus,
		d.Timeout(schema.TimeoutCreate),
		s.AppServerGroupStateRefreshFunc(d.Get("blb_id").(string), d.Id(), APPBLBFailedStatus))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapError(err)
	}

	return nil
}

func (s *APPBLBService) ListAllServerGroups(blbId string, args *appblb.DescribeAppServerGroupArgs) ([]map[string]interface{}, error) {
	serverGroupList := make([]appblb.AppServerGroup, 0)
	for {
		raw, err := s.client.WithAppBLBClient(func(client *appblb.Client) (i interface{}, e error) {
			return client.DescribeAppServerGroup(blbId, args)
		})
		if err != nil {
			return nil, WrapError(err)
		}

		result := raw.(*appblb.DescribeAppServerGroupResult)
		serverGroupList = append(serverGroupList, result.AppServerGroupList...)

		if result.IsTruncated {
			args.Marker = result.NextMarker
			args.MaxKeys = result.MaxKeys
		} else {
			break
		}
	}

	result := make([]map[string]interface{}, 0, len(serverGroupList))
	for _, sg := range serverGroupList {
		backendServers, err := s.AppServerGroupBlbRsDetail(blbId, sg.Id)
		if err != nil {
			return nil, WrapError(err)
		}
		result = append(result, map[string]interface{}{
			"sg_id":               sg.Id,
			"name":                sg.Name,
			"description":         sg.Description,
			"status":              sg.Status,
			"port_list":           s.FlattenAppServerGroupPortsToMap(sg.PortList),
			"backend_server_list": s.FlattenAppBackendServersToMap(backendServers),
		})
	}

	return result, nil
}
