/*
 * Copyright 2023 Baidu, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
 * except in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the
 * License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions
 * and limitations under the License.
 */

// et.go - the et APIs definition supported by the et service

package et

import (
	"fmt"
	"strconv"

	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/http"
)

// GetEtChannel - get an et channel
//
// PARAMS:
//    - args: the arguments to get et channel
// RETURNS:
//    - *GetEtChannelResult: the info of the et channel
//    - error: nil if success otherwise the specific error
func (c *Client) GetEtChannel(args *GetEtChannelArgs) (*GetEtChannelsResult, error) {
	if args == nil {
		return nil, fmt.Errorf("The GetEtChannelArgs cannot be nil.")
	}

	result := &GetEtChannelsResult{}
	err := bce.NewRequestBuilder(c).
		WithURL(getURLForEtChannel(args.EtId)).
		WithMethod(http.GET).
		WithQueryParamFilter("clientToken", args.ClientToken).
		WithResult(result).
		Do()

	return result, err
}

// RecommitEtChannel - recommit et channel
//
// PARAMS:
//    - args: the arguments to recommit et channel
// RETURNS:
//    - error: nil if success otherwise the specific error
func (c *Client) RecommitEtChannel(args *RecommitEtChannelArgs) error {
	if args == nil {
		return fmt.Errorf("The RecommitEtChannelArgs cannot be nil.")
	}
	err := bce.NewRequestBuilder(c).
		WithURL(getURLForEtChannelId(args.EtId, args.EtChannelId)).
		WithMethod(http.PUT).
		WithBody(args.Result).
		WithQueryParam("reCreate", "").
		WithQueryParamFilter("clientToken", args.ClientToken).
		Do()
	return err
}

// UpdateEtChannel - update et channel
//
// PARAMS:
//    - args: the arguments to update et channel
// RETURNS:
//    - error: nil if success otherwise the specific error
func (c *Client) UpdateEtChannel(args *UpdateEtChannelArgs) error {
	if args == nil {
		return fmt.Errorf("The UpdateEtChannelArgs cannot be nil.")
	}
	err := bce.NewRequestBuilder(c).
		WithURL(getURLForEtChannelId(args.EtId, args.EtChannelId)).
		WithMethod(http.PUT).
		WithBody(args.Result).
		WithQueryParam("modifyAttribute", "").
		WithQueryParamFilter("clientToken", args.ClientToken).
		Do()
	return err
}

// DeleteEtChannel - delete et channel
//
// PARAMS:
//    - args: the arguments to delete et channel
// RETURNS:
//    - error: nil if success otherwise the specific error
func (c *Client) DeleteEtChannel(args *DeleteEtChannelArgs) error {
	if args == nil {
		return fmt.Errorf("The DeleteEtChannelArgs cannot be nil.")
	}

	err := bce.NewRequestBuilder(c).
		WithURL(getURLForEtChannelId(args.EtId, args.EtChannelId)).
		WithMethod(http.DELETE).
		WithQueryParamFilter("clientToken", args.ClientToken).
		Do()
	return err
}

// EnableEtChannelIPv6 - enable et channel ipv6
//
// PARAMS:
//    - args: the arguments to enable et channel ipv6
// RETURNS:
//    - error: nil if success otherwise the specific error
func (c *Client) EnableEtChannelIPv6(args *EnableEtChannelIPv6Args) error {
	if args == nil {
		return fmt.Errorf("The EnableEtChannelIPv6Args cannot be nil.")
	}
	err := bce.NewRequestBuilder(c).
		WithURL(getURLForEtChannelId(args.EtId, args.EtChannelId)).
		WithMethod(http.PUT).
		WithBody(args.Result).
		WithQueryParam("enableIpv6", "").
		WithQueryParamFilter("clientToken", args.ClientToken).
		Do()
	return err
}

// DisableEtChannelIPv6 - disable EtChannelIPv6 with the specified parameters
//
// PARAMS:
//   - args: the arguments to disable EtChannelIPv6
//
// RETURNS:
//   - error: nil if success otherwise the specific error
func (c *Client) DisableEtChannelIPv6(args *DisableEtChannelIPv6Args) error {
	if args == nil {
		return fmt.Errorf("the createEtChannelRouteRuleArgs cannot be nil")
	}
	return bce.NewRequestBuilder(c).
		WithURL(getURLForEtChannelId(args.EtId, args.EtChannelId)).
		WithMethod(http.PUT).
		WithBody(args).
		WithQueryParam("disableIpv6", "").
		WithQueryParamFilter("clientToken", args.ClientToken).
		Do()
}


// CreateEtDcphy - init a new Et
//
// PARAMS:
//     - args: the arguments to init et dcphy
// RETURNS:
//     - CreateEtDcphyResult: the id of et dcphy newly created
//     - error: nil if success otherwise the specific error
func (c *Client) CreateEtDcphy(args *CreateEtDcphyArgs) (*CreateEtDcphyResult, error) {
	if args == nil {
		return nil, fmt.Errorf("The CreateEtDcphyArgs can not be nil")
	}

	result := &CreateEtDcphyResult{}
	err := bce.NewRequestBuilder(c).
		WithURL(getURLForEt() + "/init").
		WithMethod(http.POST).
		WithBody(args).
		WithQueryParamFilter("clientToken", args.ClientToken).
		WithResult(result).
		Do()

	return result, err
}

// UpdateEtDcphy - update an existed Et
//
// PARAMS:
//     - edId: the id of et dcphy
//     - args: the arguments to update et dcphy
// RETURNS:
//     - error: nil if success otherwise the specific error
func (c *Client) UpdateEtDcphy(dcphyId string, args *UpdateEtDcphyArgs) error {
	if len(dcphyId) == 0 {
		return fmt.Errorf("please set et dcphy id")
	}

	err := bce.NewRequestBuilder(c).
		WithURL(getURLForEtId(dcphyId)).
		WithMethod(http.PUT).
		WithBody(args).
		WithQueryParamFilter("clientToken", args.ClientToken).
		Do()

	return err
}

// ListEtDcphy - List ets
//
// PARAMS:
//     - args: the arguments to list et
// RETURNS:
//     - ListEtDcphyResult: list result
//     - error: nil if success otherwise the specific error
func (c *Client) ListEtDcphy(args *ListEtDcphyArgs) (*ListEtDcphyResult, error) {
	if args == nil {
		args = &ListEtDcphyArgs{}
	}

	if args.MaxKeys <= 0 || args.MaxKeys > 1000 {
		args.MaxKeys = 1000
	}

	result := &ListEtDcphyResult{}
	err := bce.NewRequestBuilder(c).
		WithMethod(http.GET).
		WithURL(getURLForEt()).
		WithQueryParamFilter("marker", args.Marker).
		WithQueryParamFilter("maxKeys", strconv.Itoa(args.MaxKeys)).
		WithQueryParamFilter("status", args.Status).
		WithResult(result).
		Do()

	return result, err
}

// ListEtDcphyDetail - List specific et detail
//
// PARAMS:
//     - dcphyId: the id of etDcphy
// RETURNS:
//     - EtDcphyDetail: etDcphy detail
//     - error: nil if success otherwise the specific error
func (c *Client) ListEtDcphyDetail(dcphyId string) (*EtDcphyDetail, error) {
	if len(dcphyId) == 0 {
		return nil, fmt.Errorf("please set et dcphy id")
	}

	result := &EtDcphyDetail{}
	err := bce.NewRequestBuilder(c).
		WithMethod(http.GET).
		WithURL(getURLForEtId(dcphyId)).
		WithResult(result).
		Do()

	return result, err
}

// CreateEtChannel - create an Et channel with the specific parameters
//
// PARAMS:
//     - args: the arguments to create an eip
// RETURNS:
//     - CreateEipResult: the result of create EIP, contains new EIP's address
//     - error: nil if success otherwise the specific error
func (c *Client) CreateEtChannel(args *CreateEtChannelArgs) (*CreateEtChannelResult, error) {
	if args == nil {
		return nil, fmt.Errorf("please set create etChannel argments")
	}

	if len(args.EtId) == 0 {
		return nil, fmt.Errorf("please set et id")
	}

	result := &CreateEtChannelResult{}
	err := bce.NewRequestBuilder(c).
		WithMethod(http.POST).
		WithURL(getURLForEtChannel(args.EtId)).
		WithQueryParamFilter("clientToken", args.ClientToken).
		WithBody(args).
		WithResult(result).
		Do()

	return result, err
}

// CreateEtChannelRouteRule - create a new EtChannelRouteRule with the specified parameters
//
// PARAMS:
//   - args: the arguments to create EtChannelRouteRule
//
// RETURNS:
//   - *CreateEtChannelRouteRuleResult: the id of the EtChannelRouteRule newly created
//   - error: nil if success otherwise the specific error
func (c *Client) CreateEtChannelRouteRule(args *CreateEtChannelRouteRuleArgs) (*CreateEtChannelRouteRuleResult, error) {
	if args == nil {
		return nil, fmt.Errorf("the createEtChannelRouteRuleArgs cannot be nil")
	}
	result := &CreateEtChannelRouteRuleResult{}
	err := bce.NewRequestBuilder(c).
		WithURL(getURLForEtChannelRouteRule(args.EtId, args.EtChannelId)).
		WithMethod(http.POST).
		WithBody(args).
		WithQueryParamFilter("clientToken", args.ClientToken).
		WithResult(result).
		Do()
	return result, err
}

// ListEtChannelRouteRule - list all EtChannelRouteRules with the specified parameters
//
// PARAMS:
//   - args: the arguments to list EtChannelRouteRules
//
// RETURNS:
//   - *EtChannelRouteRuleResult: the result of all EtChannelRouteRules
//   - error: nil if success otherwise the specific error
func (c *Client) ListEtChannelRouteRule(args *ListEtChannelRouteRuleArgs) (*ListEtChannelRouteRuleResult, error) {
	if args == nil {
		return nil, fmt.Errorf("the listEtChannelRouteRuleArgs cannot be nil")
	}
	if args.MaxKeys < 0 || args.MaxKeys > 1000 {
		return nil, fmt.Errorf("the field maxKeys is out of range [0, 1000]")
	}
	result := &ListEtChannelRouteRuleResult{}
	builder := bce.NewRequestBuilder(c).
		WithURL(getURLForEtChannelRouteRule(args.EtId, args.EtChannelId)).
		WithMethod(http.GET).
		WithResult(result).
		WithQueryParamFilter("marker", args.Marker)
	if args.MaxKeys != 0 {
		builder.WithQueryParamFilter("maxKeys", strconv.Itoa(args.MaxKeys))
	}
	err := builder.Do()
	return result, err
}

// UpdateEtChannelRouteRule - update a specified EtChannelRouteRule
//
// PARAMS:
//   - args: the arguments to update EtChannelRouteRule
//
// RETURNS:
//   - error: nil if success otherwise the specific error
func (c *Client) UpdateEtChannelRouteRule(args *UpdateEtChannelRouteRuleArgs) error {
	if args == nil {
		return fmt.Errorf("the updateEtChannelRouteRuleArgs cannot be nil")
	}
	return bce.NewRequestBuilder(c).
		WithURL(getURLForEtChannelRouteRuleId(args.EtId, args.EtChannelId, args.RouteRuleId)).
		WithMethod(http.PUT).
		WithBody(args).
		WithQueryParamFilter("clientToken", args.ClientToken).
		Do()
}

// DeleteEtChannelRouteRule - delete a specified EtChannelRouteRule
//
// PARAMS:
//   - params: the arguments to delete EtChannelRouteRule
//
// RETURNS:
//   - error: nil if success otherwise the specific error
func (c *Client) DeleteEtChannelRouteRule(args *DeleteEtChannelRouteRuleArgs) error {
	if args == nil {
		return fmt.Errorf("the deleteEtChannelRouteRuleArgs cannot be nil")
	}
	return bce.NewRequestBuilder(c).
		WithURL(getURLForEtChannelRouteRuleId(args.EtId, args.EtChannelId, args.RouteRuleId)).
		WithMethod(http.DELETE).
		WithBody(args).
		WithQueryParamFilter("clientToken", args.ClientToken).
		Do()

}
