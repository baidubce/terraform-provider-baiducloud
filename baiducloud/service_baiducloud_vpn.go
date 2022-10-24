package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/bcc/api"
	"github.com/baidubce/bce-sdk-go/services/vpn"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

type VpnService struct {
	client *connectivity.BaiduClient
}

func (s *VpnService) VpnGatewayStateRefresh(vpnId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		action := "Query VPN gateway " + vpnId
		raw, err := s.client.WithVPNClient(func(vpnClient *vpn.Client) (i interface{}, e error) {
			return vpnClient.GetVpnGatewayDetail(vpnId)
		})
		if err != nil {
			if NotFoundError(err) {
				return 0, string(api.InstanceStatusDeleted), nil
			}
			return nil, "", WrapErrorf(err, DefaultErrorMsg, "baiducloud_vpn_gateway", action, BCESDKGoERROR)
		}
		if raw == nil {
			return nil, "", nil
		}
		addDebug(action, raw)
		result, _ := raw.(*vpn.VPN)
		return result, string(result.Status), nil
	}
}

func (s *VpnService) VpnGatewayDetail(vpnId string) (*vpn.VPN, error) {
	action := "Query VPN gateway " + vpnId
	raw, err := s.client.WithVPNClient(func(vpnClient *vpn.Client) (i interface{}, e error) {
		return vpnClient.GetVpnGatewayDetail(vpnId)
	})
	if err != nil {
		return nil, err
	}
	addDebug(action, raw)
	return raw.(*vpn.VPN), nil
}

func (s *VpnService) ListVpnGateways(vpcId string, eip string) ([]vpn.VPN, error) {
	action := "Query VPN Gateways "
	args := &vpn.ListVpnGatewayArgs{
		VpcId: vpcId,
		Eip:   eip,
	}

	raw, err := s.client.WithVPNClient(func(vpnClient *vpn.Client) (i interface{}, e error) {
		return vpnClient.ListVpnGateway(args)
	})
	addDebug(action, raw)

	if err != nil {
		return nil, err
	}
	result, _ := raw.(*vpn.ListVpnGatewayResult)
	return result.Vpns, nil
}

func (s *VpnService) ListVpnConns(vpnId string) ([]vpn.VpnConn, error) {
	action := "List all VPN conns "

	raw, err := s.client.WithVPNClient(func(vpnClient *vpn.Client) (i interface{}, e error) {
		return vpnClient.ListVpnConn(vpnId)
	})
	addDebug(action, raw)

	if err != nil {
		return nil, err
	}
	result, _ := raw.(*vpn.ListVpnConnResult)
	return result.VpnConns, nil
}

func (s *VpnService) connToMap(conn vpn.VpnConn) map[string]interface{} {
	res := map[string]interface{}{
		"vpn_id":         conn.VpnId,
		"vpn_conn_id":    conn.VpnConnId,
		"vpn_conn_name":  conn.VpnConnName,
		"secret_key":     conn.SecretKey,
		"local_subnets":  conn.LocalSubnets,
		"remote_ip":      conn.RemoteIp,
		"remote_subnets": conn.RemoteSubnets,
		"description":    conn.Description,
		"status":         conn.Status,
		"created_time":   conn.CreatedTime,
		"health_status":  conn.HealthStatus,
		"ike_config": map[string]interface{}{
			"ike_version":   conn.IkeConfig.IkeVersion,
			"ike_pfs":       conn.IkeConfig.IkePfs,
			"ike_mode":      conn.IkeConfig.IkeMode,
			"ike_enc_alg":   conn.IkeConfig.IkeEncAlg,
			"ike_auth_alg":  conn.IkeConfig.IkeAuthAlg,
			"ike_life_time": conn.IkeConfig.IkeLifeTime,
		},
		"ipsec_config": map[string]interface{}{
			"ipsec_enc_alg":   conn.IpsecConfig.IpsecEncAlg,
			"ipsec_auth_alg":  conn.IpsecConfig.IpsecAuthAlg,
			"ipsec_pfs":       conn.IpsecConfig.IpsecPfs,
			"ipsec_life_time": conn.IpsecConfig.IpsecLifetime,
		},
	}
	return res
}
