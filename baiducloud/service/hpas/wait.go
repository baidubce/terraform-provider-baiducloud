package hpas

import (
	"time"

	"github.com/baidubce/bce-sdk-go/services/hpas/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

const (
	InstanceAvailableTimeout = 10 * time.Minute
)

func waitInstanceAvailable(conn *connectivity.BaiduClient, instanceID string) (*api.HpasResponse, error) {
	stateConf := &resource.StateChangeConf{
		Delay:   0,
		Pending: []string{InstanceStatusCreating, InstanceStatusPassword, InstanceStatusStarting, InstanceStatusReboot, InstanceStatusRebuild},
		Target:  []string{InstanceStatusActive},
		Refresh: StatusInstance(conn, instanceID),
		Timeout: InstanceAvailableTimeout,
	}

	raw, err := stateConf.WaitForState()
	if v, ok := raw.(*api.HpasResponse); ok {
		return v, nil
	}
	return nil, err
}
