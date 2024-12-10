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

// client.go - define the client for BLS service

// Package bls defines the BLS services of BCE. The supported APIs are all defined in sub-package
// model with five types: 5 LogStore APIs, 1 LogStream API, 3 logRecord APIs, 5 FastQuery APIs
// and 3 index APIs.

package bls

import (
	"encoding/json"

	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/model"
	"github.com/baidubce/bce-sdk-go/services/bls/api"
)

const (
	DEFAULT_SERVICE_DOMAIN = "bls-log.bj.baidubce.com"
)

type Client struct {
	*bce.BceClient
}

type BlsClientConfiguration struct {
	Ak       string
	Sk       string
	Endpoint string
}

func NewClient(ak, sk, endpoint string) (*Client, error) {
	return NewClientWithConfig(&BlsClientConfiguration{
		Ak:       ak,
		Sk:       sk,
		Endpoint: endpoint,
	})
}

func NewClientWithConfig(config *BlsClientConfiguration) (*Client, error) {
	ak, sk, endpoint := config.Ak, config.Sk, config.Endpoint
	if len(endpoint) == 0 {
		endpoint = DEFAULT_SERVICE_DOMAIN
	}
	client, _ := bce.NewBceClientWithAkSk(ak, sk, endpoint)
	return &Client{client}, nil
}

// LogStore opts
func (c *Client) CreateLogStore(logStore string, retention int) error {
	params, jsonErr := json.Marshal(&api.LogStore{
		LogStoreName: logStore,
		Retention:    retention,
	})
	if jsonErr != nil {
		return jsonErr
	}
	body, err := bce.NewBodyFromString(string(params))
	if err != nil {
		return err
	}
	return api.CreateLogStore(c, body)
}

func (c *Client) CreateLogStoreWithTags(logStore string, retention int, tags []model.TagModel) error {
	params, jsonErr := json.Marshal(&api.LogStore{
		LogStoreName: logStore,
		Retention:    retention,
		Tags:         tags,
	})
	if jsonErr != nil {
		return jsonErr
	}
	body, err := bce.NewBodyFromString(string(params))
	if err != nil {
		return err
	}
	return api.CreateLogStore(c, body)
}

func (c *Client) UpdateLogStore(logStore string, retention int) error {
	param, jsonErr := json.Marshal(&api.LogStore{
		Retention: retention,
	})
	if jsonErr != nil {
		return jsonErr
	}
	body, err := bce.NewBodyFromString(string(param))
	if err != nil {
		return err
	}
	return api.UpdateLogStore(c, DefaultProject, logStore, body)
}

func (c *Client) DescribeLogStore(logStore string) (*api.LogStore, error) {
	return api.DescribeLogStore(c, DefaultProject, logStore)
}

func (c *Client) DeleteLogStore(logStore string) error {
	return api.DeleteLogStore(c, DefaultProject, logStore)
}

func (c *Client) ListLogStore(args *api.QueryConditions) (*api.ListLogStoreResult, error) {
	return api.ListLogStore(c, DefaultProject, args)
}

// LogStream opt
func (c *Client) ListLogStream(logStore string, args *api.QueryConditions) (*api.ListLogStreamResult, error) {
	return api.ListLogStream(c, DefaultProject, logStore, args)
}

// LogRecord opts
func (c *Client) PushLogRecord(logStore string, logStream string, logType string, logRecords []api.LogRecord) error {
	params, jsonErr := json.Marshal(&api.PushLogRecordBody{
		LogStreamName: logStream,
		Type:          logType,
		LogRecords:    logRecords,
	})
	if jsonErr != nil {
		return jsonErr
	}
	body, err := bce.NewBodyFromString(string(params))
	if err != nil {
		return err
	}
	return api.PushLogRecord(c, DefaultProject, logStore, body)
}

func (c *Client) PullLogRecord(logStore string, args *api.PullLogRecordArgs) (*api.PullLogRecordResult, error) {
	return api.PullLogRecord(c, DefaultProject, logStore, args)
}

func (c *Client) QueryLogRecord(logStore string, args *api.QueryLogRecordArgs) (*api.QueryLogResult, error) {
	return api.QueryLogRecord(c, DefaultProject, logStore, args)
}

// FastQuery opts
func (c *Client) CreateFastQuery(args *api.CreateFastQueryBody) error {
	params, jsonErr := json.Marshal(args)
	if jsonErr != nil {
		return jsonErr
	}
	body, err := bce.NewBodyFromString(string(params))
	if err != nil {
		return nil
	}
	return api.CreateFastQuery(c, body)
}

func (c *Client) DescribeFastQuery(fastQueryName string) (*api.FastQuery, error) {
	return api.DescribeFastQuery(c, fastQueryName)
}

func (c *Client) UpdateFastQuery(fastQueryName string, args *api.UpdateFastQueryBody) error {
	params, jsonErr := json.Marshal(args)
	if jsonErr != nil {
		return jsonErr
	}
	body, err := bce.NewBodyFromString(string(params))
	if err != nil {
		return nil
	}
	return api.UpdateFastQuery(c, body, fastQueryName)
}

func (c *Client) DeleteFastQuery(fastQueryName string) error {
	return api.DeleteFastQuery(c, fastQueryName)
}

func (c *Client) ListFastQuery(args *api.QueryConditions) (*api.ListFastQueryResult, error) {
	return api.ListFastQuery(c, "", "", args)
}

// Index opts
func (c *Client) CreateIndex(logStore string, fulltext bool, fields map[string]api.LogField) error {
	params, jsonErr := json.Marshal(&api.IndexFields{
		FullText: fulltext,
		Fields:   fields,
	})
	if jsonErr != nil {
		return jsonErr
	}
	body, err := bce.NewBodyFromString(string(params))
	if err != nil {
		return err
	}
	return api.CreateIndex(c, DefaultProject, logStore, body)
}

func (c *Client) UpdateIndex(logStore string, fulltext bool, fields map[string]api.LogField) error {
	params, jsonErr := json.Marshal(&api.IndexFields{
		FullText: fulltext,
		Fields:   fields,
	})
	if jsonErr != nil {
		return jsonErr
	}
	body, err := bce.NewBodyFromString(string(params))
	if err != nil {
		return err
	}
	return api.UpdateIndex(c, DefaultProject, logStore, body)
}

func (c *Client) DeleteIndex(logStore string) error {
	return api.DeleteIndex(c, DefaultProject, logStore)
}

func (c *Client) DescribeIndex(logStore string) (*api.IndexFields, error) {
	return api.DescribeIndex(c, DefaultProject, logStore)
}

// LogShipper opts
func (c *Client) ListLogShipper(args *api.ListLogShipperCondition) (*api.ListShipperResult, error) {
	return api.ListLogShipper(c, args)
}

func (c *Client) CreateLogShipper(args *api.CreateLogShipperBody) (string, error) {
	params, jsonErr := json.Marshal(args)
	if jsonErr != nil {

		return "", jsonErr
	}
	body, err := bce.NewBodyFromString(string(params))
	if err != nil {

		return "", nil
	}
	return api.CreateLogShipper(c, body)
}

func (c *Client) UpdateLogShipper(logShipperID string, args *api.UpdateLogShipperBody) error {
	params, jsonErr := json.Marshal(args)
	if jsonErr != nil {
		return jsonErr
	}
	body, err := bce.NewBodyFromString(string(params))
	if err != nil {
		return nil
	}
	return api.UpdateLogShipper(c, body, logShipperID)
}

func (c *Client) GetLogShipper(logShipperID string) (*api.LogShipper, error) {
	return api.GetLogShipper(c, logShipperID)
}

func (c *Client) ListLogShipperRecord(logShipperID string, args *api.ListShipperRecordCondition) (*api.ListShipperRecordResult, error) {
	return api.ListLogShipperRecord(c, logShipperID, args)
}

func (c *Client) DeleteSingleLogShipper(logShipperID string) error {
	return api.DeleteSingleLogShipper(c, logShipperID)
}

func (c *Client) BulkDeleteLogShipper(args *api.BulkDeleteShipperCondition) error {
	params, jsonErr := json.Marshal(args)
	if jsonErr != nil {
		return jsonErr
	}
	body, err := bce.NewBodyFromString(string(params))
	if err != nil {
		return nil
	}
	return api.BulkDeleteLogShipper(c, body)
}

func (c *Client) SetSingleLogShipperStatus(logShipperID string, args *api.SetSingleShipperStatusCondition) error {
	params, jsonErr := json.Marshal(args)
	if jsonErr != nil {
		return jsonErr
	}
	body, err := bce.NewBodyFromString(string(params))
	if err != nil {
		return nil
	}
	return api.SetSingleLogShipperStatus(c, logShipperID, body)
}

func (c *Client) BulkSetLogShipperStatus(args *api.BulkSetShipperStatusCondition) error {
	params, jsonErr := json.Marshal(args)
	if jsonErr != nil {
		return jsonErr
	}
	body, err := bce.NewBodyFromString(string(params))
	if err != nil {
		return nil
	}
	return api.BulkSetLogShipperStatus(c, body)
}
