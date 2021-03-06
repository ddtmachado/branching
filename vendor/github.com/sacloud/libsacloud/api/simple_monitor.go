// Copyright 2016-2020 The Libsacloud Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	"encoding/json" //	"strings"
	"fmt"

	"github.com/sacloud/libsacloud/sacloud"
)

//HACK: さくらのAPI側仕様: CommonServiceItemsの内容によってJSONフォーマットが異なるため
//      DNS/GSLB/シンプル監視それぞれでリクエスト/レスポンスデータ型を定義する。

// SearchSimpleMonitorResponse シンプル監視検索レスポンス
type SearchSimpleMonitorResponse struct {
	// Total 総件数
	Total int `json:",omitempty"`
	// From ページング開始位置
	From int `json:",omitempty"`
	// Count 件数
	Count int `json:",omitempty"`
	// SimpleMonitors シンプル監視 リスト
	SimpleMonitors []sacloud.SimpleMonitor `json:"CommonServiceItems,omitempty"`
}

type simpleMonitorRequest struct {
	SimpleMonitor *sacloud.SimpleMonitor `json:"CommonServiceItem,omitempty"`
	From          int                    `json:",omitempty"`
	Count         int                    `json:",omitempty"`
	Sort          []string               `json:",omitempty"`
	Filter        map[string]interface{} `json:",omitempty"`
	Exclude       []string               `json:",omitempty"`
	Include       []string               `json:",omitempty"`
}

type simpleMonitorResponse struct {
	*sacloud.ResultFlagValue
	*sacloud.SimpleMonitor `json:"CommonServiceItem,omitempty"`
}

// SimpleMonitorAPI シンプル監視API
type SimpleMonitorAPI struct {
	*baseAPI
}

// NewSimpleMonitorAPI シンプル監視API作成
func NewSimpleMonitorAPI(client *Client) *SimpleMonitorAPI {
	return &SimpleMonitorAPI{
		&baseAPI{
			client: client,
			FuncGetResourceURL: func() string {
				return "commonserviceitem"
			},
			FuncBaseSearchCondition: func() *sacloud.Request {
				res := &sacloud.Request{}
				res.AddFilter("Provider.Class", "simplemon")
				return res
			},
		},
	}
}

// Find 検索
func (api *SimpleMonitorAPI) Find() (*SearchSimpleMonitorResponse, error) {
	data, err := api.client.newRequest("GET", api.getResourceURL(), api.getSearchState())
	if err != nil {
		return nil, err
	}
	var res SearchSimpleMonitorResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (api *SimpleMonitorAPI) request(f func(*simpleMonitorResponse) error) (*sacloud.SimpleMonitor, error) {
	res := &simpleMonitorResponse{}
	err := f(res)
	if err != nil {
		return nil, err
	}
	return res.SimpleMonitor, nil
}

func (api *SimpleMonitorAPI) createRequest(value *sacloud.SimpleMonitor) *simpleMonitorResponse {
	return &simpleMonitorResponse{SimpleMonitor: value}
}

// New 新規作成用パラメーター作成
func (api *SimpleMonitorAPI) New(target string) *sacloud.SimpleMonitor {
	return sacloud.CreateNewSimpleMonitor(target)
}

// Create 新規作成
func (api *SimpleMonitorAPI) Create(value *sacloud.SimpleMonitor) (*sacloud.SimpleMonitor, error) {
	return api.request(func(res *simpleMonitorResponse) error {
		return api.create(api.createRequest(value), res)
	})
}

// Read 読み取り
func (api *SimpleMonitorAPI) Read(id sacloud.ID) (*sacloud.SimpleMonitor, error) {
	return api.request(func(res *simpleMonitorResponse) error {
		return api.read(id, nil, res)
	})
}

// Update 更新
func (api *SimpleMonitorAPI) Update(id sacloud.ID, value *sacloud.SimpleMonitor) (*sacloud.SimpleMonitor, error) {
	return api.request(func(res *simpleMonitorResponse) error {
		return api.update(id, api.createRequest(value), res)
	})
}

// Delete 削除
func (api *SimpleMonitorAPI) Delete(id sacloud.ID) (*sacloud.SimpleMonitor, error) {
	return api.request(func(res *simpleMonitorResponse) error {
		return api.delete(id, nil, res)
	})
}

// Health ヘルスチェック
//
// まだチェックが行われていない場合nilを返す
func (api *SimpleMonitorAPI) Health(id sacloud.ID) (*sacloud.SimpleMonitorHealthCheckStatus, error) {
	var (
		method = "GET"
		uri    = fmt.Sprintf("%s/%d/health", api.getResourceURL(), id)
	)
	res := struct {
		SimpleMonitor *sacloud.SimpleMonitorHealthCheckStatus `json:",omitempty"`
	}{}

	err := api.baseAPI.request(method, uri, nil, &res)
	if err != nil {
		return nil, err
	}
	return res.SimpleMonitor, nil
}

// MonitorResponseTimeSec アクティビティーモニター(レスポンスタイム)取得
func (api *SimpleMonitorAPI) MonitorResponseTimeSec(id sacloud.ID, body *sacloud.ResourceMonitorRequest) (*sacloud.MonitorValues, error) {
	var (
		method = "GET"
		uri    = fmt.Sprintf("%s/%d/activity/responsetimesec/monitor", api.getResourceURL(), id)
	)
	res := &sacloud.ResourceMonitorResponse{}
	err := api.baseAPI.request(method, uri, body, res)
	if err != nil {
		return nil, err
	}
	return res.Data, nil
}
