package snic

import (
	"time"

	"github.com/baidubce/bce-sdk-go/services/endpoint"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

const (
	SNICAvailableTimeout = 5 * time.Minute
)

func waitSNICAvailable(conn *connectivity.BaiduClient, snicId string) (*endpoint.Endpoint, error) {
	stateConf := &resource.StateChangeConf{
		Delay:   0,
		Pending: []string{SNICStatusUnavailable},
		Target:  []string{SNICStatusAvailable},
		Refresh: StatusSNIC(conn, snicId),
		Timeout: SNICAvailableTimeout,
	}

	raw, err := stateConf.WaitForState()
	if v, ok := raw.(*endpoint.Endpoint); ok {
		return v, nil
	}
	return nil, err
}
