package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

type VpcService struct {
	client *connectivity.BaiduClient
}

func (s *VpcService) ListAllVpcs() ([]vpc.VPC, error) {
	listVpcArgs := &vpc.ListVPCArgs{}
	action := "List all VPCs"

	vpcs := make([]vpc.VPC, 0)
	for {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (i interface{}, e error) {
			return vpcClient.ListVPC(listVpcArgs)
		})
		if err != nil {
			return nil, WrapErrorf(err, DefaultErrorMsg, "baiducloud_vpcs", action, BCESDKGoERROR)
		}

		result, _ := raw.(*vpc.ListVPCResult)
		vpcs = append(vpcs, result.VPCs...)
		if !result.IsTruncated {
			break
		}
		listVpcArgs.Marker = result.NextMarker
		listVpcArgs.MaxKeys = result.MaxKeys
	}

	return vpcs, nil
}

func (s *VpcService) NatGatewayStateRefresh(natId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		action := "Query Nat Gateway " + natId
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (i interface{}, e error) {
			return vpcClient.GetNatGatewayDetail(natId)
		})
		addDebug(action, raw)
		if err != nil {
			if NotFoundError(err) {
				return 0, string(vpc.NAT_STATUS_DELETED), nil
			}
			return nil, "", WrapErrorf(err, DefaultErrorMsg, "baiducloud_nat_gateway", action, BCESDKGoERROR)
		}
		if raw == nil {
			return nil, "", nil
		}

		result, _ := raw.(*vpc.NAT)
		return result, string(result.Status), nil
	}
}

func (s *VpcService) PeerConnStateRefresh(peerConnId string, role vpc.PeerConnRoleType) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		action := "Query Peer Conn " + peerConnId
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (i interface{}, e error) {
			return vpcClient.GetPeerConnDetail(peerConnId, role)
		})
		addDebug(action, raw)
		if err != nil {
			if NotFoundError(err) || IsExceptedErrors(err, PeerConnNotFound) {
				return 0, string(vpc.PEERCONN_STATUS_DELETED), nil
			}
			return nil, "", WrapErrorf(err, DefaultErrorMsg, "baiducloud_peer_conn", action, BCESDKGoERROR)
		}
		if raw == nil {
			return nil, "", nil
		}

		result, _ := raw.(*vpc.PeerConn)
		return result, string(result.Status), nil
	}
}

func (s *VpcService) PeerConnDNSStatusRefresh(peerConnId string, role vpc.PeerConnRoleType) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		action := "Query Peer Conn DNS status " + peerConnId
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (i interface{}, e error) {
			return vpcClient.GetPeerConnDetail(peerConnId, role)
		})
		addDebug(action, err)
		if err != nil {
			return nil, "", WrapErrorf(err, DefaultErrorMsg, "baiducloud_peer_conn", action, BCESDKGoERROR)
		}

		result, _ := raw.(*vpc.PeerConn)
		return result, string(result.DnsStatus), nil
	}
}

func (s *VpcService) DescribeAclRules(subnetId string) ([]vpc.AclRule, error) {
	args := &vpc.ListAclRulesArgs{
		SubnetId: subnetId,
	}
	action := "Describe ACL Rules"

	aclRules := make([]vpc.AclRule, 0)
	for {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (i interface{}, e error) {
			return vpcClient.ListAclRules(args)
		})
		if err != nil {
			return nil, WrapErrorf(err, DefaultErrorMsg, "baiducloud_acl", action, BCESDKGoERROR)
		}

		result, _ := raw.(*vpc.ListAclRulesResult)
		aclRules = append(aclRules, result.AclRules...)
		if !result.IsTruncated {
			break
		}
		args.Marker = result.NextMarker
	}

	return aclRules, nil
}

func (s *VpcService) ListAllSubnets(args *vpc.ListSubnetArgs) ([]vpc.Subnet, error) {
	action := "List all subnets for VPC " + args.VpcId

	subnets := make([]vpc.Subnet, 0)
	for {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (i interface{}, e error) {
			return vpcClient.ListSubnets(args)
		})
		if err != nil {
			return nil, err
		}
		addDebug(action, raw)

		result, _ := raw.(*vpc.ListSubnetResult)
		subnets = append(subnets, result.Subnets...)

		if !result.IsTruncated {
			break
		}
		args.Marker = result.NextMarker
	}

	return subnets, nil
}

func (s *VpcService) ListAllAclEntrysWithVPCID(vpcID string) ([]vpc.AclEntry, error) {
	action := "List all ACLs for VPC " + vpcID

	acls := make([]vpc.AclEntry, 0)

	raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (i interface{}, e error) {
		return vpcClient.ListAclEntrys(vpcID)
	})
	if err != nil {
		return nil, err
	}
	addDebug(action, raw)

	result, _ := raw.(*vpc.ListAclEntrysResult)
	acls = append(acls, result.AclEntrys...)

	return acls, nil
}

func (s *VpcService) ListAllAclRulesWithSubnetID(subnetID string) ([]vpc.AclRule, error) {
	action := "List all ACLs for subnet " + subnetID

	acls := make([]vpc.AclRule, 0)
	args := &vpc.ListAclRulesArgs{
		SubnetId: subnetID,
	}
	for {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (i interface{}, e error) {
			return vpcClient.ListAclRules(args)
		})
		if err != nil {
			return nil, err
		}
		addDebug(action, raw)

		result, _ := raw.(*vpc.ListAclRulesResult)
		acls = append(acls, result.AclRules...)

		if !result.IsTruncated {
			break
		}
		args.Marker = result.NextMarker
		args.MaxKeys = result.MaxKeys
	}

	return acls, nil
}

func (s *VpcService) ListAllNatGateways(args *vpc.ListNatGatewayArgs) ([]vpc.NAT, error) {
	action := "List all NAT gateways for vpc " + args.VpcId

	nats := make([]vpc.NAT, 0)
	for {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (i interface{}, e error) {
			return vpcClient.ListNatGateway(args)
		})
		if err != nil {
			return nil, err
		}
		addDebug(action, raw)

		result, _ := raw.(*vpc.ListNatGatewayResult)
		nats = append(nats, result.Nats...)

		if !result.IsTruncated {
			break
		}
		args.Marker = result.NextMarker
	}

	return nats, nil
}

func (s *VpcService) ListAllPeerConns(vpcID string) ([]vpc.PeerConn, error) {
	action := "List all Peer Conns for vpc " + vpcID

	args := &vpc.ListPeerConnsArgs{
		VpcId: vpcID,
	}
	peerConns := make([]vpc.PeerConn, 0)
	for {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (i interface{}, e error) {
			return vpcClient.ListPeerConn(args)
		})
		if err != nil {
			return nil, err
		}
		addDebug(action, raw)

		result, _ := raw.(*vpc.ListPeerConnsResult)
		peerConns = append(peerConns, result.PeerConns...)

		if !result.IsTruncated {
			break
		}
		args.Marker = result.NextMarker
	}

	return peerConns, nil
}

func (s *VpcService) GetVPCDetail(vpcID string) (*vpc.GetVPCDetailResult, error) {
	action := "Get VPC detail " + vpcID
	raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (i interface{}, e error) {
		return vpcClient.GetVPCDetail(vpcID)
	})
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, "baiducloud_vpc", action, BCESDKGoERROR)
	}

	result, _ := raw.(*vpc.GetVPCDetailResult)
	return result, nil
}

func (s *VpcService) GetSubnetDetail(subnetID string) (*vpc.GetSubnetDetailResult, error) {
	action := "Get Subnet detail " + subnetID
	raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (i interface{}, e error) {
		return vpcClient.GetSubnetDetail(subnetID)
	})
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, "baiducloud_subnet", action, BCESDKGoERROR)
	}

	result, _ := raw.(*vpc.GetSubnetDetailResult)
	return result, nil
}

func (s *VpcService) GetRouteTableDetail(routeTableID, vpcID string) (*vpc.GetRouteTableResult, error) {
	action := "Get route table detail " + routeTableID
	raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (i interface{}, e error) {
		return vpcClient.GetRouteTableDetail(routeTableID, vpcID)
	})
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, "baiducloud_route_rule", action, BCESDKGoERROR)
	}

	result, _ := raw.(*vpc.GetRouteTableResult)
	return result, nil
}

func (s *VpcService) GetNatGatewayDetail(natID string) (*vpc.NAT, error) {
	action := "Get Nat Gateway detail " + natID
	raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (i interface{}, e error) {
		return vpcClient.GetNatGatewayDetail(natID)
	})
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, "baiducloud_nat_gateway", action, BCESDKGoERROR)
	}

	result, _ := raw.(*vpc.NAT)
	return result, nil
}

func (s *VpcService) GetPeerConnDetail(peerConnID string, role vpc.PeerConnRoleType) (*vpc.PeerConn, error) {
	action := "Get PeerConn detail " + peerConnID
	raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (i interface{}, e error) {
		return vpcClient.GetPeerConnDetail(peerConnID, role)
	})
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, "baiducloud_nat_gateway", action, BCESDKGoERROR)
	}

	result, _ := raw.(*vpc.PeerConn)
	return result, nil
}

func (s *VpcService) OpenPeerConnDNSSync(d *schema.ResourceData, peerConnId string, role vpc.PeerConnRoleType) error {
	args := &vpc.PeerConnSyncDNSArgs{
		Role:        role,
		ClientToken: buildClientToken(),
	}

	return resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		result, err := s.GetPeerConnDetail(peerConnId, role)
		if err != nil {
			return resource.NonRetryableError(err)
		}

		// If the dns status is open, then nothing to do.
		if result.DnsStatus == vpc.DNS_STATUS_OPEN {
			return nil
		}

		// If the dns status is one of following, then retry.
		for _, status := range []vpc.DnsStatusType{vpc.DNS_STATUS_CLOSING, vpc.DNS_STATUS_SYNCING, vpc.DNS_STATUS_WAIT} {
			if result.DnsStatus == status {
				return resource.RetryableError(err)
			}
		}

		// If the dns status is close, then open the dns sync directly.
		if _, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (i interface{}, e error) {
			return nil, vpcClient.OpenPeerConnSyncDNS(peerConnId, args)
		}); err != nil {
			return resource.NonRetryableError(err)
		}

		stateConf := buildStateConf(
			[]string{string(vpc.DNS_STATUS_WAIT), string(vpc.DNS_STATUS_SYNCING), string(vpc.DNS_STATUS_CLOSE)},
			[]string{string(vpc.DNS_STATUS_OPEN)},
			d.Timeout(schema.TimeoutUpdate),
			s.PeerConnDNSStatusRefresh(peerConnId, role))
		if _, err := stateConf.WaitForState(); err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})
}

func (s *VpcService) ClosePeerConnDNSSync(d *schema.ResourceData, peerConnId string, role vpc.PeerConnRoleType) error {
	args := &vpc.PeerConnSyncDNSArgs{
		Role:        role,
		ClientToken: buildClientToken(),
	}

	return resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		result, err := s.GetPeerConnDetail(peerConnId, role)
		if err != nil {
			return resource.NonRetryableError(err)
		}

		// If the dns status is close, then nothing to do.
		if result.DnsStatus == vpc.DNS_STATUS_CLOSE {
			return nil
		}

		// If the dns status is one of following, then retry.
		for _, status := range []vpc.DnsStatusType{vpc.DNS_STATUS_CLOSING, vpc.DNS_STATUS_SYNCING, vpc.DNS_STATUS_WAIT} {
			if result.DnsStatus == status {
				return resource.RetryableError(err)
			}
		}

		// If the dns status is open, then close the dns sync directly.
		if _, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (i interface{}, e error) {
			return nil, vpcClient.ClosePeerConnSyncDNS(peerConnId, args)
		}); err != nil {
			return resource.NonRetryableError(err)
		}

		stateConf := buildStateConf(
			[]string{string(vpc.DNS_STATUS_OPEN), string(vpc.DNS_STATUS_CLOSING)},
			[]string{string(vpc.DNS_STATUS_CLOSE)},
			d.Timeout(schema.TimeoutUpdate),
			s.PeerConnDNSStatusRefresh(d.Id(), role))
		if _, err := stateConf.WaitForState(); err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})
}
