package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/localDns"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

type PrivateZoneDnsService struct {
	client *connectivity.BaiduClient
}

const (
	RECORDENABLE  = "enable"
	RECORDDISABLE = "pause"
)

func (s *PrivateZoneDnsService) ListAlDnsRecords(zoneId string) ([]localDns.Record, error) {
	result := make([]localDns.Record, 0)

	action := "Get private zone records list " + zoneId
	raw, err := s.client.WithLocalDnsClient(func(localDnsClient *localDns.Client) (i interface{}, e error) {
		return localDnsClient.ListRecord(zoneId)
	})
	if err != nil {
		return nil, err
	}
	addDebug(action, raw)

	response := raw.(*localDns.ListRecordResponse)
	result = append(result, response.Records...)

	return result, nil
}
func UpdateStatus(localDnsClient *localDns.Client, recordId string, status string) error {
	if status == RECORDENABLE {
		return localDnsClient.EnableRecord(recordId, buildClientToken())
	} else if status == RECORDDISABLE {
		return localDnsClient.DisableRecord(recordId, buildClientToken())
	} else {
		return Error(InvalidInputField, "status")
	}
}
func inputFieldCheck(d *schema.ResourceData) error {
	var (
		recordType string
		priority   int
	)
	if v, ok := d.GetOk("type"); ok {
		recordType = v.(string)
	}
	if v, ok := d.GetOk("priority"); ok {
		priority = v.(int)
	}
	if recordType != MXType {
		if priority != 0 {
			return Error(InvalidInputField, "priority")
		}
	}
	return nil
}
