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
