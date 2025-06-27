package eip

import (
	"time"

	"github.com/baidubce/bce-sdk-go/services/eip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func waitEipGroupAvailable(conn *connectivity.BaiduClient, timeout time.Duration, eipGroupID string) (*eip.EipGroupModel, error) {
	stateConf := &resource.StateChangeConf{
		Delay:   5 * time.Second,
		Pending: []string{EIPStatusCreating, EIPStatusBinding, EIPStatusUnBinding, EIPStatusUpdating},
		Target:  []string{EIPStatusAvailable},
		Refresh: StatusEipGroup(conn, eipGroupID),
		Timeout: timeout,
	}

	raw, err := stateConf.WaitForState()
	if v, ok := raw.(*eip.EipGroupModel); ok {
		return v, nil
	}
	return nil, err
}

func waitEipGroupDeleted(conn *connectivity.BaiduClient, timeout time.Duration, eipGroupID string) error {
	stateConf := &resource.StateChangeConf{
		Delay:   5 * time.Second,
		Pending: []string{EIPStatusAvailable, EIPStatusExpired},
		Target:  []string{},
		Refresh: StatusEipGroup(conn, eipGroupID),
		Timeout: timeout,
	}
	_, err := stateConf.WaitForState()
	return err
}
