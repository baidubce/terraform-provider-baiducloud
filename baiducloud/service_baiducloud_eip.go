package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/eip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

type EipService struct {
	client *connectivity.BaiduClient
}

func (e *EipService) EipResizeBandwidth(ip string, new int) error {
	action := "Resize Eip Bandwidth " + ip

	resizeEipArgs := &eip.ResizeEipArgs{
		NewBandWidthInMbps: new,
		ClientToken:        buildClientToken(),
	}

	_, err := e.client.WithEipClient(func(client *eip.Client) (i interface{}, e error) {
		return nil, client.ResizeEip(ip, resizeEipArgs)
	})

	if err != nil {
		if !IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_eip", "", BCESDKGoERROR)
		}
	}

	addDebug(action, resizeEipArgs)
	return nil
}

func (e *EipService) StartAutoRenew(ip string, args *eip.StartAutoRenewArgs) error {
	action := "Start Eip Auto Renew " + ip

	_, err := e.client.WithEipClient(func(client *eip.Client) (i interface{}, e error) {
		return nil, client.StartAutoRenew(ip, args)
	})

	if err != nil {
		if !IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_eip", "", BCESDKGoERROR)
		}
	}

	addDebug(action, action)
	return nil
}

func (e *EipService) StopAutoRenew(ip string) error {
	action := "Stop Eip Auto Renew " + ip

	_, err := e.client.WithEipClient(func(client *eip.Client) (i interface{}, e error) {
		return nil, client.StopAutoRenew(ip, buildClientToken())
	})

	if err != nil {
		if !IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_eip", "", BCESDKGoERROR)
		}
	}

	addDebug(action, action)
	return nil
}

func (e *EipService) EipBind(ip, instanceType, instanceId string) error {
	action := "Bind Eip " + ip

	bindEipArgs := &eip.BindEipArgs{
		InstanceType: instanceType,
		InstanceId:   instanceId,
		ClientToken:  buildClientToken(),
	}

	_, err := e.client.WithEipClient(func(client *eip.Client) (i interface{}, e error) {
		return nil, client.BindEip(ip, bindEipArgs)
	})

	if err != nil {
		if !IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_eip", "", BCESDKGoERROR)
		}
	}

	addDebug(action, bindEipArgs)
	return nil
}

func (e *EipService) EipUnBind(ip string) error {
	action := "UnBind Eip " + ip

	_, err := e.client.WithEipClient(func(client *eip.Client) (i interface{}, e error) {
		return nil, client.UnBindEip(ip, buildClientToken())
	})

	if err != nil && !NotFoundError(err) {
		// if before unbind, relate resource like blb and instance has been deleted,
		// may return UnsupportedOperation error
		// so we check eip status again
		eipDetail, errDetail := e.EipGetDetail(ip)
		if errDetail != nil {
			if NotFoundError(errDetail) {
				addDebug(action, "")
				return nil
			}

			// return unbind err
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_eip", "", BCESDKGoERROR)
		}

		if stringInSlice([]string{EIPStatusBinded, EIPStatusBinding}, eipDetail.Status) {
			// eip still in use
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_eip", "", BCESDKGoERROR)
		}
		// eip not in use, return unbind success
		addDebug(action, "")
		return nil
	}

	addDebug(action, "")
	return nil
}

func (e *EipService) EipGetDetail(ip string) (*eip.EipModel, error) {
	if ip == "" {
		return nil, WrapError(Error("eip can not be empty if get eip detail"))
	}

	listArgs := &eip.ListEipArgs{
		Eip: ip,
	}

	raw, err := e.client.WithEipClient(func(client *eip.Client) (i interface{}, e error) {
		return client.ListEip(listArgs)
	})

	if err != nil {
		return nil, WrapError(err)
	}

	response := raw.(*eip.ListEipResult)
	if len(response.EipList) == 0 {
		return nil, WrapError(Error(ResourceNotFound))
	}

	return &response.EipList[0], nil
}

func (e *EipService) EipStateRefreshFunc(ip string, failState []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		result, err := e.EipGetDetail(ip)
		if err != nil {
			return nil, "", WrapError(err)
		}

		for _, statue := range failState {
			if result.Status == statue {
				return result, result.Status, WrapError(Error(GetFailTargetStatus, result.Status))
			}
		}

		return result, result.Status, nil
	}
}

func (e *EipService) ListAllEips(listArgs *eip.ListEipArgs) ([]eip.EipModel, error) {
	result := make([]eip.EipModel, 0)
	for {
		raw, err := e.client.WithEipClient(func(client *eip.Client) (interface{}, error) {
			return client.ListEip(listArgs)
		})

		if err != nil {
			return nil, WrapError(err)
		}

		response := raw.(*eip.ListEipResult)
		result = append(result, response.EipList...)

		if response.IsTruncated {
			listArgs.MaxKeys = response.MaxKeys
			listArgs.Marker = response.Marker
		} else {
			return result, nil
		}
	}
}

func (e *EipService) FlattenEipModelsToMap(eips []eip.EipModel) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(eips))

	for _, e := range eips {
		result = append(result, map[string]interface{}{
			"eip":               e.Eip,
			"name":              e.Name,
			"status":            e.Status,
			"eip_instance_type": e.EipInstanceType,
			"share_group_id":    e.ShareGroupId,
			"bandwidth_in_mbps": e.BandWidthInMbps,
			"payment_timing":    e.PaymentTiming,
			"billing_method":    e.BillingMethod,
			"create_time":       e.CreateTime,
			"expire_time":       e.ExpireTime,
			"tags":              flattenTagsToMap(e.Tags),
		})
	}

	return result
}
