package eip

import (
	"fmt"
	"log"

	"github.com/baidubce/bce-sdk-go/services/eip"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/flex"
)

func FindEipGroup(conn *connectivity.BaiduClient, eipGroupID string) (*eip.EipGroupModel, error) {
	raw, err := conn.WithEipClient(func(client *eip.Client) (interface{}, error) {
		return client.EipGroupDetail(eipGroupID)
	})

	log.Printf("[DEBUG] Read EIP Group detail result: %+v", raw)
	if err != nil {
		log.Printf("[DEBUG] Read EIP Group detail error: %s", err)
		if flex.IsResourceNotFound(err) {
			return nil, nil
		}

		return nil, fmt.Errorf("error reading EIP Group (%s) detail: %w", eipGroupID, err)
	}
	response := raw.(*eip.EipGroupModel)

	return response, nil
}

func FindEipDDosProtection(conn *connectivity.BaiduClient, ip string) (*eip.DdosModel, error) {
	args := eip.ListDdosRequest{Ips: ip}
	raw, err := conn.WithEipClient(func(client *eip.Client) (interface{}, error) {
		return client.ListDdos(&args)
	})

	log.Printf("[DEBUG] Read EIP DDoS Protection detail result: %+v", raw)
	if err != nil {
		return nil, fmt.Errorf("error reading DDoS protection threshold for eip (%s): %w", ip, err)
	}
	response, _ := raw.(*eip.ListDdosResponse)
	if response.DdosList == nil || len(*response.DdosList) == 0 {
		return nil, nil
	}
	return &(*response.DdosList)[0], nil
}

func FindEip(conn *connectivity.BaiduClient, ip string) (*eip.EipModel, error) {
	args := &eip.ListEipArgs{Eip: ip}
	raw, err := conn.WithEipClient(func(client *eip.Client) (interface{}, error) {
		return client.ListEip(args)
	})

	log.Printf("[DEBUG] Read EIP detail result: %+v", raw)
	if err != nil {
		return nil, fmt.Errorf("error reading EIP detail for eip (%s): %w", ip, err)
	}
	response, _ := raw.(*eip.ListEipResult)
	if response.EipList == nil || len(response.EipList) == 0 {
		return nil, nil
	}
	return &response.EipList[0], nil
}
