/*
 * Copyright 2021 Baidu, Inc.
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

// logstore.go - the logStore APIs definition supported by the BLS service

package api

import (
	"strconv"

	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/http"
)

// CreateLogStore - create logStore
//
// PARAMS:
//   - cli: the client agent which can perform sending request
//   - body: logStore parameters body
//
// RETURNS:
//   - error: nil if success otherwise the specific error
func CreateLogStore(cli bce.Client, body *bce.Body) error {
	req := &bce.BceRequest{}
	req.SetUri(LOGSTORE_PREFIX)
	req.SetMethod(http.POST)
	if body != nil {
		req.SetBody(body)
	}
	resp := &bce.BceResponse{}
	if err := cli.SendRequest(req, resp); err != nil {
		return err
	}
	if resp.IsFail() {
		return resp.ServiceError()
	}
	defer func() { resp.Body().Close() }()
	return nil
}

// UpdateLogStore - update logStore retention
//
// PARAMS:
//   - cli: the client agent which can perform sending request
//   - project: logstore project
//   - logStore: logStore to update
//   - body: logStore parameters body
//
// RETURNS:
//   - error: nil if success otherwise the specific error
func UpdateLogStore(cli bce.Client, project string, logStore string, body *bce.Body) error {
	req := &bce.BceRequest{}
	req.SetUri(getLogStoreUri(logStore))
	req.SetMethod(http.PUT)
	req.SetParam("project", project)
	req.SetBody(body)
	resp := &bce.BceResponse{}
	if err := cli.SendRequest(req, resp); err != nil {
		return err
	}
	if resp.IsFail() {
		return resp.ServiceError()
	}
	defer func() { resp.Body().Close() }()
	return nil
}

// DescribeLogStore - get logStore info
//
// PARAMS:
//   - cli: the client agent which can perform sending request
//   - project: logstore project
//   - logStore: logStore to get
//
// RETURNS:
//   - *LogStore: logStore info
//   - error: nil if success otherwise the specific error
func DescribeLogStore(cli bce.Client, project string, logStore string) (*LogStore, error) {
	req := &bce.BceRequest{}
	req.SetUri(getLogStoreUri(logStore))
	req.SetParam("project", project)
	req.SetMethod(http.GET)
	resp := &bce.BceResponse{}
	if err := cli.SendRequest(req, resp); err != nil {
		return nil, err
	}
	if resp.IsFail() {
		return nil, resp.ServiceError()
	}
	result := &LogStore{}
	if err := resp.ParseJsonBody(result); err != nil {
		return nil, err
	}
	return result, nil
}

// DeleteLogStore - delete logStore
//
// PARAMS:
//   - cli: the client agent which can perform sending request
//   - project: logstore project
//   - logStore: logStore to delete
//
// RETURNS:
//   - error: nil if success otherwise the specific error
func DeleteLogStore(cli bce.Client, project, logStore string) error {
	req := &bce.BceRequest{}
	req.SetUri(getLogStoreUri(logStore))
	req.SetParam("project", project)
	req.SetMethod(http.DELETE)
	resp := &bce.BceResponse{}
	if err := cli.SendRequest(req, resp); err != nil {
		return err
	}
	if resp.IsFail() {
		return resp.ServiceError()
	}
	defer func() { resp.Body().Close() }()
	return nil
}

// ListLogStore - get all pattern-match logStore info
//
// PARAMS:
//   - cli: the client agent which can perform sending request
//   - project: logstore project
//   - args: conditions logStore should match
//
// RETURNS:
//   - *ListLogStoreResult: logStore result set
//   - error: nil if success otherwise the specific error
func ListLogStore(cli bce.Client, project string, args *QueryConditions) (*ListLogStoreResult, error) {
	req := &bce.BceRequest{}
	req.SetUri(LOGSTORE_PREFIX)
	req.SetParam("project", project)
	req.SetMethod(http.GET)
	// Set optional args
	if args != nil {
		if args.NamePattern != "" {
			req.SetParam("namePattern", args.NamePattern)
		}
		if args.Order != "" {
			req.SetParam("order", args.Order)
		}
		if args.OrderBy != "" {
			req.SetParam("orderBy", args.OrderBy)
		}
		if args.PageNo > 0 {
			req.SetParam("pageNo", strconv.Itoa(args.PageNo))
		}
		if args.PageSize > 0 {
			req.SetParam("pageSize", strconv.Itoa(args.PageSize))
		}
	}
	resp := &bce.BceResponse{}
	if err := cli.SendRequest(req, resp); err != nil {
		return nil, err
	}
	if resp.IsFail() {
		return nil, resp.ServiceError()
	}
	result := &ListLogStoreResult{}
	if err := resp.ParseJsonBody(result); err != nil {
		return nil, err
	}
	return result, nil
}

func ListLogStoreV2(cli bce.Client, body *bce.Body) (*ListLogStoreResult, error) {
	req := &bce.BceRequest{}
	req.SetUri(LIST_LOGSTORE_PREFIX)
	req.SetMethod(http.POST)
	if body != nil {
		req.SetBody(body)
	}
	resp := &bce.BceResponse{}
	if err := cli.SendRequest(req, resp); err != nil {
		return nil, err
	}
	if resp.IsFail() {
		return nil, resp.ServiceError()
	}
	result := &ListLogStoreResult{}
	if err := resp.ParseJsonBody(result); err != nil {
		return nil, err
	}
	return result, nil
}

func GetLogStoreByProjects(cli bce.Client, body *bce.Body) (*BatchLogStoreResult, error) {
	req := &bce.BceRequest{}
	req.SetUri(BATCH_PREFIX)
	req.SetMethod(http.POST)
	if body != nil {
		req.SetBody(body)
	}
	resp := &bce.BceResponse{}
	if err := cli.SendRequest(req, resp); err != nil {
		return nil, err
	}
	if resp.IsFail() {
		return nil, resp.ServiceError()
	}
	result := &BatchLogStoreResult{}
	if err := resp.ParseJsonBody(result); err != nil {
		return nil, err
	}
	return result, nil
}

func BindResource(cli bce.Client, body *bce.Body) error {
	req := &bce.BceRequest{}
	req.SetUri(BIND_PREFIX)
	req.SetMethod(http.POST)
	if body != nil {
		req.SetBody(body)
	}
	resp := &bce.BceResponse{}
	if err := cli.SendRequest(req, resp); err != nil {
		return err
	}
	if resp.IsFail() {
		return resp.ServiceError()
	}
	defer func() { resp.Body().Close() }()
	return nil
}

func UnBindResource(cli bce.Client, body *bce.Body) error {
	req := &bce.BceRequest{}
	req.SetUri(UNBIND_PREFIX)
	req.SetMethod(http.POST)
	if body != nil {
		req.SetBody(body)
	}
	resp := &bce.BceResponse{}
	if err := cli.SendRequest(req, resp); err != nil {
		return err
	}
	if resp.IsFail() {
		return resp.ServiceError()
	}
	defer func() { resp.Body().Close() }()
	return nil
}
