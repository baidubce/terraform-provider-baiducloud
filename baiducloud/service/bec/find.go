package bec

import (
	"log"

	"github.com/baidubce/bce-sdk-go/services/bec"
	"github.com/baidubce/bce-sdk-go/services/bec/api"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func FindVMInstance(conn *connectivity.BaiduClient, vmID string) (*api.VmInstanceDetailsVo, error) {
	raw, err := conn.WithBECClient(func(client *bec.Client) (interface{}, error) {
		return client.GetVirtualMachine(vmID)
	})
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] Query BEC VM Instance (%s) detail: %+v", vmID, raw)
	return raw.(*api.VmInstanceDetailsVo), nil
}
