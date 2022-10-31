package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/sms"
	"github.com/baidubce/bce-sdk-go/services/sms/api"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

type SMSService struct {
	client *connectivity.BaiduClient
}

func (s *SMSService) GetSMSSignatureDetail(signatureId string) (*api.GetSignatureResult, error) {
	action := "Query SMS Signature " + signatureId

	detailArgs := &api.GetSignatureArgs{
		SignatureId: signatureId,
	}
	raw, err := s.client.WithSMSClient(func(smsClient *sms.Client) (i interface{}, e error) {
		return smsClient.GetSignature(detailArgs)
	})
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, "baiducloud_sms_signature", action, BCESDKGoERROR)
	}

	return raw.(*api.GetSignatureResult), nil
}
