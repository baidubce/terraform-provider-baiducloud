package baiducloud

import (
	"fmt"

	"github.com/baidubce/bce-sdk-go/services/appblb"
)

func (s *APPBLBService) DescribeListener(blbId, protocol string, port int) (interface{}, error) {
	args := &appblb.DescribeAppListenerArgs{
		ListenerPort: uint16(port),
	}
	action := fmt.Sprintf("Describe APPBLB %s Listener [%s.%d]", blbId, protocol, port)

	raw, err := s.client.WithAppBLBClient(func(client *appblb.Client) (i interface{}, e error) {
		switch protocol {
		case TCP:
			return client.DescribeAppTCPListeners(blbId, args)
		case UDP:
			return client.DescribeAppUDPListeners(blbId, args)
		case HTTP:
			return client.DescribeAppHTTPListeners(blbId, args)
		case HTTPS:
			return client.DescribeAppHTTPSListeners(blbId, args)
		case SSL:
			return client.DescribeAppSSLListeners(blbId, args)
		default:
			return nil, fmt.Errorf("unsupport listener type: %s", protocol)
		}
	})
	addDebug(action, raw)

	if err != nil {
		return nil, WrapError(err)
	}

	switch protocol {
	case TCP:
		response := raw.(*appblb.DescribeAppTCPListenersResult)
		if len(response.ListenerList) > 0 {
			return &response.ListenerList[0], nil
		}
	case UDP:
		response := raw.(*appblb.DescribeAppUDPListenersResult)
		if len(response.ListenerList) > 0 {
			return &response.ListenerList[0], nil
		}
	case HTTP:
		response := raw.(*appblb.DescribeAppHTTPListenersResult)
		if len(response.ListenerList) > 0 {
			return &response.ListenerList[0], nil
		}
	case HTTPS:
		response := raw.(*appblb.DescribeAppHTTPSListenersResult)
		if len(response.ListenerList) > 0 {
			return &response.ListenerList[0], nil
		}
	case SSL:
		response := raw.(*appblb.DescribeAppSSLListenersResult)
		if len(response.ListenerList) > 0 {
			return &response.ListenerList[0], nil
		}
	default:
	}

	return nil, WrapError(fmt.Errorf(ResourceNotFound))
}

func (s *APPBLBService) DescribePolicys(blbId, protocol string, port int) ([]appblb.AppPolicy, error) {
	args := &appblb.DescribePolicysArgs{
		Port: uint16(port),
	}
	action := fmt.Sprintf("Describe APPBLB %s Listener [%s.%d]'s policy", blbId, protocol, port)

	policyList := make([]appblb.AppPolicy, 0)
	for {
		raw, err := s.client.WithAppBLBClient(func(client *appblb.Client) (i interface{}, e error) {
			return client.DescribePolicys(blbId, args)
		})
		addDebug(action, raw)

		if err != nil {
			return nil, WrapError(err)
		}

		response := raw.(*appblb.DescribePolicysResult)
		policyList = append(policyList, response.PolicyList...)

		if protocol != TCP && response.IsTruncated {
			args.Marker = response.NextMarker
			args.MaxKeys = response.MaxKeys
		} else {
			return policyList, nil
		}
	}
}

func (s *APPBLBService) FlattenAppPolicysToMap(policys []appblb.AppPolicy) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(policys))

	for _, policy := range policys {
		pMap := map[string]interface{}{
			"id":                    policy.Id,
			"description":           policy.Description,
			"app_server_group_id":   policy.AppServerGroupId,
			"app_server_group_name": policy.AppServerGroupName,
			"frontend_port":         policy.FrontendPort,
			"backend_port":          policy.BackendPort,
			"priority":              policy.Priority,
			"port_type":             policy.PortType,
		}

		ruleList := make([]map[string]interface{}, 0)
		for _, r := range policy.RuleList {
			rule := map[string]interface{}{
				"key":   r.Key,
				"value": r.Value,
			}

			ruleList = append(ruleList, rule)
		}
		pMap["rule_list"] = ruleList

		result = append(result, pMap)
	}

	return result
}

func (s *APPBLBService) ListAllTCPListeners(blbId string, port int) ([]interface{}, error) {
	args := &appblb.DescribeAppListenerArgs{
		ListenerPort: uint16(port),
	}

	action := fmt.Sprintf("Describe APPBLB %s Listeners [TCP.%d]", blbId, port)
	listeners := make([]appblb.AppTCPListenerModel, 0)
	for {
		raw, err := s.client.WithAppBLBClient(func(client *appblb.Client) (i interface{}, e error) {
			return client.DescribeAppTCPListeners(blbId, args)
		})
		if err != nil {
			return nil, WrapError(err)
		}
		addDebug(action, raw)

		response := raw.(*appblb.DescribeAppTCPListenersResult)
		listeners = append(listeners, response.ListenerList...)

		if response.IsTruncated {
			args.Marker = response.Marker
			args.MaxKeys = response.MaxKeys
		} else {
			break
		}
	}

	result := make([]interface{}, 0, len(listeners))
	for _, listener := range listeners {
		policys, err := s.DescribePolicys(blbId, TCP, int(listener.Port))
		if err != nil {
			return nil, WrapError(err)
		}

		result = append(result, map[string]interface{}{
			"listener_port":       listener.Port,
			"tcp_session_timeout": listener.TcpSessionTimeout,
			"protocol":            TCP,
			"scheduler":           listener.Scheduler,
			"policys":             s.FlattenAppPolicysToMap(policys),
		})
	}

	return result, nil
}

func (s *APPBLBService) ListAllUDPListeners(blbId string, port int) ([]interface{}, error) {
	args := &appblb.DescribeAppListenerArgs{
		ListenerPort: uint16(port),
	}

	action := fmt.Sprintf("Describe APPBLB %s Listeners [UDP.%d]", blbId, port)
	listeners := make([]appblb.AppUDPListenerModel, 0)
	for {
		raw, err := s.client.WithAppBLBClient(func(client *appblb.Client) (i interface{}, e error) {
			return client.DescribeAppUDPListeners(blbId, args)
		})
		if err != nil {
			return nil, WrapError(err)
		}
		addDebug(action, raw)

		response := raw.(*appblb.DescribeAppUDPListenersResult)
		listeners = append(listeners, response.ListenerList...)

		if response.IsTruncated {
			args.Marker = response.Marker
			args.MaxKeys = response.MaxKeys
		} else {
			break
		}
	}

	result := make([]interface{}, 0, len(listeners))
	for _, listener := range listeners {
		policys, err := s.DescribePolicys(blbId, UDP, int(listener.Port))
		if err != nil {
			return nil, WrapError(err)
		}

		result = append(result, map[string]interface{}{
			"listener_port": listener.Port,
			"protocol":      UDP,
			"scheduler":     listener.Scheduler,
			"policys":       s.FlattenAppPolicysToMap(policys),
		})
	}

	return result, nil
}

func (s *APPBLBService) ListAllHTTPListeners(blbId string, port int) ([]interface{}, error) {
	args := &appblb.DescribeAppListenerArgs{
		ListenerPort: uint16(port),
	}

	action := fmt.Sprintf("Describe APPBLB %s Listeners [HTTP.%d]", blbId, port)
	listeners := make([]appblb.AppHTTPListenerModel, 0)
	for {
		raw, err := s.client.WithAppBLBClient(func(client *appblb.Client) (i interface{}, e error) {
			return client.DescribeAppHTTPListeners(blbId, args)
		})
		if err != nil {
			return nil, WrapError(err)
		}
		addDebug(action, raw)

		response := raw.(*appblb.DescribeAppHTTPListenersResult)
		listeners = append(listeners, response.ListenerList...)

		if response.IsTruncated {
			args.Marker = response.Marker
			args.MaxKeys = response.MaxKeys
		} else {
			break
		}
	}

	result := make([]interface{}, 0, len(listeners))
	for _, listener := range listeners {
		policys, err := s.DescribePolicys(blbId, HTTP, int(listener.ListenerPort))
		if err != nil {
			return nil, WrapError(err)
		}

		result = append(result, map[string]interface{}{
			"listener_port":            listener.ListenerPort,
			"protocol":                 HTTP,
			"scheduler":                listener.Scheduler,
			"keep_session":             listener.KeepSession,
			"keep_session_type":        listener.KeepSessionType,
			"keep_session_timeout":     listener.KeepSessionTimeout,
			"keep_session_cookie_name": listener.KeepSessionCookieName,
			"x_forwarded_for":          listener.XForwardedFor,
			"server_timeout":           listener.ServerTimeout,
			"redirect_port":            listener.RedirectPort,
			"policys":                  s.FlattenAppPolicysToMap(policys),
		})
	}

	return result, nil
}

func (s *APPBLBService) ListAllHTTPSListeners(blbId string, port int) ([]interface{}, error) {
	args := &appblb.DescribeAppListenerArgs{
		ListenerPort: uint16(port),
	}

	action := fmt.Sprintf("Describe APPBLB %s Listeners [HTTPS.%d]", blbId, port)
	listeners := make([]appblb.AppHTTPSListenerModel, 0)
	for {
		raw, err := s.client.WithAppBLBClient(func(client *appblb.Client) (i interface{}, e error) {
			return client.DescribeAppHTTPSListeners(blbId, args)
		})
		if err != nil {
			return nil, WrapError(err)
		}
		addDebug(action, raw)

		response := raw.(*appblb.DescribeAppHTTPSListenersResult)
		listeners = append(listeners, response.ListenerList...)

		if response.IsTruncated {
			args.Marker = response.Marker
			args.MaxKeys = response.MaxKeys
		} else {
			break
		}
	}

	result := make([]interface{}, 0, len(listeners))
	for _, listener := range listeners {
		policys, err := s.DescribePolicys(blbId, HTTPS, int(listener.ListenerPort))
		if err != nil {
			return nil, WrapError(err)
		}

		result = append(result, map[string]interface{}{
			"listener_port":            listener.ListenerPort,
			"protocol":                 HTTPS,
			"scheduler":                listener.Scheduler,
			"keep_session":             listener.KeepSession,
			"keep_session_type":        listener.KeepSessionType,
			"keep_session_timeout":     listener.KeepSessionTimeout,
			"keep_session_cookie_name": listener.KeepSessionCookieName,
			"x_forwarded_for":          listener.XForwardedFor,
			"server_timeout":           listener.ServerTimeout,
			"cert_ids":                 listener.CertIds,
			"encryption_type":          listener.EncryptionType,
			"encryption_protocols":     listener.EncryptionProtocols,
			"dual_auth":                listener.DualAuth,
			"client_cert_ids":          listener.ClientCertIds,
			"policys":                  s.FlattenAppPolicysToMap(policys),
		})
	}

	return result, nil
}

func (s *APPBLBService) ListAllSSLListeners(blbId string, port int) ([]interface{}, error) {
	args := &appblb.DescribeAppListenerArgs{
		ListenerPort: uint16(port),
	}

	action := fmt.Sprintf("Describe APPBLB %s Listeners [SSL.%d]", blbId, port)
	listeners := make([]appblb.AppSSLListenerModel, 0)
	for {
		raw, err := s.client.WithAppBLBClient(func(client *appblb.Client) (i interface{}, e error) {
			return client.DescribeAppSSLListeners(blbId, args)
		})
		if err != nil {
			return nil, WrapError(err)
		}
		addDebug(action, raw)

		response := raw.(*appblb.DescribeAppSSLListenersResult)
		listeners = append(listeners, response.ListenerList...)

		if response.IsTruncated {
			args.Marker = response.Marker
			args.MaxKeys = response.MaxKeys
		} else {
			break
		}
	}

	result := make([]interface{}, 0, len(listeners))
	for _, listener := range listeners {
		policys, err := s.DescribePolicys(blbId, SSL, int(listener.ListenerPort))
		if err != nil {
			return nil, WrapError(err)
		}

		result = append(result, map[string]interface{}{
			"listener_port":        listener.ListenerPort,
			"protocol":             SSL,
			"scheduler":            listener.Scheduler,
			"cert_ids":             listener.CertIds,
			"encryption_type":      listener.EncryptionType,
			"encryption_protocols": listener.EncryptionProtocols,
			"dual_auth":            listener.DualAuth,
			"client_cert_ids":      listener.ClientCertIds,
			"policys":              s.FlattenAppPolicysToMap(policys),
		})
	}

	return result, nil
}

func (s *APPBLBService) ListAllListeners(blbId, protocol string, port int) ([]interface{}, error) {
	result := make([]interface{}, 0)

	if protocol == TCP || protocol == "" {
		tcpListeners, err := s.ListAllTCPListeners(blbId, port)
		if err != nil {
			return nil, WrapError(err)
		}
		result = append(result, tcpListeners...)
	}

	if protocol == UDP || protocol == "" {
		udpListeners, err := s.ListAllUDPListeners(blbId, port)
		if err != nil {
			return nil, WrapError(err)
		}
		result = append(result, udpListeners...)
	}

	if protocol == HTTP || protocol == "" {
		httpListeners, err := s.ListAllHTTPListeners(blbId, port)
		if err != nil {
			return nil, WrapError(err)
		}
		result = append(result, httpListeners...)
	}

	if protocol == HTTPS || protocol == "" {
		httpsListeners, err := s.ListAllHTTPSListeners(blbId, port)
		if err != nil {
			return nil, WrapError(err)
		}
		result = append(result, httpsListeners...)
	}

	if protocol == SSL || protocol == "" {
		sslListeners, err := s.ListAllSSLListeners(blbId, port)
		if err != nil {
			return nil, WrapError(err)
		}
		result = append(result, sslListeners...)
	}

	return result, nil
}
