package bec

import (
	"time"

	"github.com/baidubce/bce-sdk-go/services/bec/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

const (
	VMInstanceAvailableTimeout = 10 * time.Minute
)

func waitVMInstanceAvailable(conn *connectivity.BaiduClient, vmID string) (*api.VmInstanceDetailsVo, error) {
	stateConf := &resource.StateChangeConf{
		Delay:   0,
		Pending: []string{VMInstanceStatusCreating, VMInstanceStatusRestarting},
		Target:  []string{VMInstanceStatusRunning},
		Refresh: StatusVMInstance(conn, vmID),
		Timeout: VMInstanceAvailableTimeout,
	}

	raw, err := stateConf.WaitForState()
	if v, ok := raw.(*api.VmInstanceDetailsVo); ok {
		return v, nil
	}
	return nil, err
}
